package koyfin

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"entext-applications/internal/validator"
	"github.com/goccy/go-json"
)

var (
	ErrBadRecordCount = errors.New("error: mismatched number of data and values records")
	ErrBadResponse    = errors.New("internal error: malformatted response")
	KoyfinAPIErr      = errors.New("error: Koyfin API error")
	ErrNoData         = errors.New("error: no data")
)

type SearchRequest struct {
	SearchString string   `json:"searchString" jsonschema:"Ticker or ETF name to search for"`            // ticker or ETF name
	Categories   []string `json:"categories" jsonschema:"Search categories, defaults to Equity and ETF"` // search categories, use default Equity and ETF if not set
	PrimaryOnly  bool     `json:"primaryOnly" jsonschema:"Using primary exchange, default false"`        // Using primary exchange - default false
}

// validate search req
func (s *SearchRequest) Validate(v *validator.Validator) {
	v.Check(s.SearchString != "", "searchString", "search string must not be empty")
}

type SearchResult struct {
	KID        string `json:"KID" jsonschema:"Koyfin ID"`                                    // series KID - this can be used for scraping
	Category   string `json:"category" jsonschema:"Data category"`                           // data category
	Name       string `json:"name" jsonschema:"Regular name"`                                // Regular name
	Identifier string `json:"identifier" jsonschema:"Identifier in ticker format"`           // Identifier in ticker format, i.e APPL:US
	Currency   string `json:"financialCurrency" jsonschema:"Security's financial currency"`  // Security's currency
	Country    string `json:"country" jsonschema:"Domicile, ISO 3166, alpha-2 country code"` // Security's domicile
	Exchange   string `json:"exchange,omitempty" jsonschema:"exchange where the security is traded"`
	Status     string `json:"status" jsonschema:"instrument's status"`
}

// LookUpByName performs series search by name and categories
func (c *Client) LookUpByName(session *Session, req SearchRequest) ([]SearchResult, error) {
	// if not set assuming equity and ETF tickers
	if len(req.Categories) == 0 {
		req.Categories = []string{"Equity", "ETF"}
	}

	// get bytes buffer
	b := bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer func() {
		bufferPool.Put(b)
	}()

	// create response body
	err := json.NewEncoder(b).Encode(req)
	if err != nil {
		return nil, err
	}

	// ensure access token validity
	err = c.ensureValidToken(session)
	if err != nil {
		return nil, err
	}

	// once access token validity is ensure - we can use directly without worrying about race condition
	headers := map[string][]string{
		"User-Agent":   {"Mozilla/5.0 (X11; Linux x86_64; rv:146.0) Gecko/20100101 Firefox/146.0"},
		"Accept":       {"application/json, text/plain, */*'"},
		"Content-Type": {"application/json"},
		"Origin":       {"https://app.koyfin.com"},
	}

	// populate authorization headers
	session.AuthorizeHeader(headers)

	body, err := c.getResponse("POST", "https://app.koyfin.com/api/v1/bfc/tickers/search", b, headers)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// parse response
	response := struct {
		Total int            `json:"total"`
		Data  []SearchResult `json:"data"`
	}{}

	err = json.NewDecoder(body).Decode(&response)
	if err != nil {
		return nil, err
	}

	// TODO: filter only live ticker
	return response.Data, nil
}

type TickerDataRequest struct {
	ID string `json:"id" jsonschema:"Identifier for the ticker"`
	// Indicator to search for
	Key string `json:"key" jsonschema:"Indicator to search for"`
	// data currency, default USD
	Currency string `json:"currency,omitempty" jsonschema:"Data currency, default USD"`
	// Start of data date
	DateFrom string `json:"dateFrom" jsonschema:"Start of data date, required"`
	// End of data date
	DateTo string `json:"dateTo" jsonschema:"End of data date, required"`
	// Series granularity - day, monthly, quarterly or annually
	AggPeriod string `json:"candleAggregationPeriod,omitempty" jsonschema:"Series granularity - day, monthly, quarterly or annually"`
	// price format - values: both, standard, adj
	PriceFormat string `json:"priceFormat,omitempty" jsonschema:"Price format - values: both, standard, adj"`
	// Financial report period . Accepted values: quarterly, annual, LTM
	FinPeriod string `json:"financialPeriodType,omitempty" jsonschema:"Financial report period - quarterly, annual, LTM"`
	// ForwardPeriod applies only to series on future estimate
	ForwardPeriod string `json:"financialForwardPeriod,omitempty" jsonschema:"Forward period for future estimates"`
	// Use financial report's end date instead of CY conversion. Set to false by default
	UseReportDate bool `json:"-" jsonschema:"Use financial report's end date instead of CY conversion"`
	// UseRawPrice set rerturn data to raw price
	// UseRawPrice bool	`json:"-"`
}

// show response data in packed schema
type GraphData struct {
	// observation date in YYYY-MM-DD format
	Date []string `json:"date" jsonschema:"Observation dates in YYYY-MM-DD format"`
	// Close price - this applies to equity price only
	Close []float64 `json:"close,omitzero" jsonschema:"Close prices for equity"`
	// Close price adjusted for dividend
	AdjClose []float64 `json:"adjClose,omitzero" jsonschema:"Adjusted close prices for dividends"`
	// Observation value
	Value []float64 `json:"value,omitzero" jsonschema:"Observation values"`
	// Graph metadata. used for converting financial reporting date to corresponding CY date
	Meta []Metadata `json:"_meta" jsonschema:"Metadata for converting financial reporting dates"`
	// reporing error
	Error *KoyfinError `json:"error,omitempty" jsonschema:"Error information"`
}

type Metadata struct {
	EndDate string `json:"periodenddate" jsonschema:"Period end date"`
}

type TimeSeries struct {
	Labels []string  `json:"labels" jsonschema:"data series labels, array of date value or financial quarter"`
	Values []float64 `json:"values" jsonschema:"data series values"`
}

type TickerDataResponse struct {
	Graph *GraphData   `json:"graph" jsonschema:"Graph data response"`
	Error *KoyfinError `json:"error" jsonschema:"Error message if any"`
}

// validate ticker data request
func (t *TickerDataRequest) Validate(v *validator.Validator) {
	v.Check(t.ID != "", "id", "id must not be empty")
	v.Check(t.Key != "", "key", "key must not be empty")

	// Validate date format if DateFrom is provided
	if t.DateFrom != "" {
		_, err := time.Parse("2006-01-02", t.DateFrom)
		v.Check(err == nil, "dateFrom", "must be a valid date in YYYY-MM-DD format")
	}

	// Validate date format if DateTo is provided
	if t.DateTo != "" {
		_, err := time.Parse("2006-01-02", t.DateTo)
		v.Check(err == nil, "dateTo", "must be a valid date in YYYY-MM-DD format")
	}

	// Validate that DateFrom is not after DateTo if both are provided
	if t.DateFrom != "" && t.DateTo != "" {
		fromDate, fromErr := time.Parse("2006-01-02", t.DateFrom)
		toDate, toErr := time.Parse("2006-01-02", t.DateTo)

		if fromErr == nil && toErr == nil {
			v.Check(!fromDate.After(toDate), "dateFrom", "must not be after dateTo")
		}
	}
}

// GetDataSeries retrieves data for a given koyfin ID and series key
func (c *Client) GetDataSeries(session *Session, req TickerDataRequest) (*TimeSeries, error) {
	// get bytes buffer
	b := bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer func() {
		bufferPool.Put(b)
	}()

	// set default value
	if req.DateFrom == "" {
		req.DateFrom = time.Now().Add(-time.Hour * 24 * 365).Format("2006-01-02")
	}
	if req.DateTo == "" {
		req.DateTo = time.Now().Format("2006-01-02")
	}
	if req.AggPeriod == "" {
		req.AggPeriod = "day"
	}

	key := req.Key
	// decide price format using type of data series key
	switch {

	case strings.HasPrefix(key, "fest") || (req.FinPeriod != "" && req.PriceFormat == ""):
		req.PriceFormat = "standard"

	case key == "p_candle_range":
		req.PriceFormat = "adj"
	}

	dividendAdjPrice := key == "p_candle_range"

	// create response body
	err := json.NewEncoder(b).Encode(req)
	if err != nil {
		return nil, err
	}

	// ensure access token validity
	err = c.ensureValidToken(session)
	if err != nil {
		return nil, err
	}

	// once access token validity is ensure - we can use directly without worrying about race condition
	headers := map[string][]string{
		"User-Agent":   {"Mozilla/5.0 (X11; Linux x86_64; rv:146.0) Gecko/20100101 Firefox/146.0"},
		"Accept":       {"application/json, text/plain, */*'"},
		"Content-Type": {"application/json"},
		"Origin":       {"https://app.koyfin.com"},
	}

	// populate authorization headers
	session.AuthorizeHeader(headers)

	body, err := c.getResponse("POST", "https://app.koyfin.com/api/v3/data/graph?schema=packed", b, headers)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// parse response
	response := TickerDataResponse{}
	err = json.NewDecoder(body).Decode(&response)
	if err != nil {
		return nil, errors.Join(ErrBadResponse, err)
	}

	// API response error
	if response.Error != nil {
		return nil, errors.Join(KoyfinAPIErr, errors.New(response.Error.Message))
	}

	// graph data exists or not
	if response.Graph == nil {
		return nil, ErrNoData
	}

	graph := response.Graph
	// count date values
	iteratorCount := len(graph.Date)
	metaDataCount := len(graph.Meta)

	// convert quarter to calendar year equivalent
	hasMetaData := metaDataCount > 0
	convertToCY := hasMetaData && !req.UseReportDate && req.FinPeriod == "quarterly"

	if hasMetaData && iteratorCount < metaDataCount {
		iteratorCount = metaDataCount
	}

	var values []float64
	if dividendAdjPrice {
		values = graph.AdjClose
	} else if req.Key == "p_candle_range" {
		values = graph.Close
	} else {
		values = graph.Value
	}
	// count values in series
	valueCount := len(values)

	// date and values count should equals each other
	if iteratorCount != valueCount {
		return nil, ErrBadRecordCount
	}

	var dateLabels []string
	if convertToCY {
		dateLabels = make([]string, iteratorCount)
		for i := 0; i < iteratorCount; i++ {
			dateLabels[i] = graph.Meta[i].EndDate
		}
	} else {
		dateLabels = graph.Date
	}

	return &TimeSeries{Labels: dateLabels, Values: values}, nil
}
