package koyfin

import (
	"bytes"

	"github.com/goccy/go-json"
)

// Note: Using set of most useful indicator instead of leaving choice to the user
type ETFHoldingComponent struct {
	// KID extracted from map key
	KID                string         `json:"KID" jsonschema:"Koyfin identifier for the holding"`
	Context            Context        `json:"__context" jsonschema:"Additional context data for the holding"`
	CurrentPrice       *IndicatorData `json:"current_price_adj,omitempty" jsonschema:"Current adjusted price of the holding"`
	PrevFiveYearPrice  *IndicatorData `json:"prev_5y_price_adj,omitempty" jsonschema:"Price 5 years ago adjusted"`
	PrevThreeYearPrice *IndicatorData `json:"prev_3y_price_adj,omitempty" jsonschema:"Price 3 years ago adjusted"`
	PrevYearPrice      *IndicatorData `json:"prev_1y_price_adj,omitempty" jsonschema:"Price 1 year ago adjusted"`
	BeginYearPrice     *IndicatorData `json:"start_year_price_adj,omitempty" jsonschema:"Price at the beginning of the current year"`
	// financial metrics
	ShortInterest   *IndicatorData `json:"short_interest_pct,omitempty" jsonschema:"Percentage of shares sold short"`
	MktCap          *IndicatorData `json:"market_cap_mil,omitempty" jsonschema:"Market capitalization in millions"`
	EBIT            *IndicatorData `json:"ebit_last_12m,omitempty" jsonschema:"Earnings before interest and taxes margin percentage"`
	EBITDANTM       *IndicatorData `json:"ebitda_next_12m,omitempty" jsonschema:"EBITDA ratio for next 12 months"`
	EBITDALTM       *IndicatorData `json:"ebitda_last_12m,omitempty" jsonschema:"EBITDA ratio for past 12 months"`
	RevLTM          *IndicatorData `json:"revenue_last_12m,omitempty" jsonschema:"last 12m revenue"`
	RevNTM          *IndicatorData `json:"revenue_next_12m,omitempty" jsonschema:"next 12m revenue"`
	EBITDMargin     *IndicatorData `json:"ebita_margin_ltm,omitempty" jsonschema:"last 12m EBITDA margin"`
	EnterpriseValue *IndicatorData `json:"current_ev,omitempty" jsonschema:"enterprice value"`
	// CurrUnAdjPrice  *IndicatorData `json:"unadj_price,omitempty" jsonschema:"current unadjusted price"`
	EstEPSNTM *IndicatorData `json:"eps_next_12m,omitempty" jsonschema:"earning-per-share next 12m"`
	EPSLTM    *IndicatorData `json:"eps_last_12m,omitempty" jsonschema:"earning-per-share last 12m"`
	Beta1Y    *IndicatorData `json:"1y_wkly_beta,omitempty" jsonschema:"1Y weekly beta"`
	// security data
	Ticker   string `json:"ticker,omitempty" jsonschema:"Stock ticker symbol"`
	Name     string `json:"name,omitempty" jsonschema:"Name of the holding"`
	Sector   string `json:"sector,omitempty" jsonschema:"Sector classification of the holding"`
	Industry string `json:"industry,omitempty" jsonschema:"Industry classification of the holding"`
	Country  string `json:"trading_country,omitempty" jsonschema:"trading country"`
}

type Context struct {
	Weight      float64 `json:"weight" jsonschema:"Weight of the holding in the ETF"`
	Country     string  `json:"rawCountry" jsonschema:"Country of the holding"`
	HoldingType string  `json:"rawType" jsonschema:"Type classification of the holding"`
	Name        string  `json:"rawName" jsonschema:"Raw name of the holding"`
}

type ETFHoldingsInfo map[string]map[string]*ETFHoldingComponent

type ETFInformation struct {
	ETFKID   string                 `json:"etf_kid" jsonschema:"Koyfin ID the given ETF"`
	Holdings []*ETFHoldingComponent `json:"holdings" jsonschema:"list of ETF holdings with given financial metrics"`
}

func (c *Client) ListETFHoldings(session *Session, req SnapshotRequest) ([]ETFInformation, error) {
	snapshotIDs := make([]SnapshotID, len(req.KIDs))
	for indx, kid := range req.KIDs {
		snapshotIDs[indx] = SnapshotID{KID: kid, IDType: "holdings"}
	}

	// asset type not supported
	keys := DefaultKeysForAssetType(req)
	if keys == nil {
		return nil, ErrAssetTypeNotSupported
	}

	// ensure access token validity
	err := c.ensureValidToken(session)
	if err != nil {
		return nil, err
	}

	// get bytes buffer
	b := bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer func() {
		bufferPool.Put(b)
	}()

	requestData := SnapshotDataRequest{
		KIDs: snapshotIDs,
		Keys: keys,
	}

	// create response body
	err = json.NewEncoder(b).Encode(requestData)
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

	body, err := c.getResponse("POST", "https://app.koyfin.com/api/v3/data/keys", b, headers)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// define anonymous struct due to variable schema and we don't use this outside the function
	output := struct {
		Holdings ETFHoldingsInfo `json:"holdings"`
	}{}

	err = json.NewDecoder(body).Decode(&output)
	if err != nil {
		return nil, err
	}

	if len(output.Holdings) == 0 {
		return nil, ErrNoSnapshotData
	}

	etfInfoContainer := make([]ETFInformation, 0, len(output.Holdings))
	// return only one series data
	for etfKID, etfHoldingTable := range output.Holdings {
		holdings := make([]*ETFHoldingComponent, 0, len(etfHoldingTable))
		for kid, component := range etfHoldingTable {
			component.KID = kid
			holdings = append(holdings, component)
		}
		// add to result
		etfInfoContainer = append(etfInfoContainer, ETFInformation{ETFKID: etfKID, Holdings: holdings})
	}

	return etfInfoContainer, nil
}
