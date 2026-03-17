package koyfin

import (
	"bytes"
	"errors"
	"fmt"

	"entext-applications/internal/utils"
	"entext-applications/internal/validator"

	"github.com/goccy/go-json"
	"github.com/google/jsonschema-go/jsonschema"
)

const (
	MaxFacetValue float64 = 9007199254740991
	MinFacetValue float64 = -9007199254740991

	// known filter keys
	FilterKeyISO2                = "iso2"
	FilterKeySector              = "t_sec"
	FilterKeyStyleClassification = "style_classification"
	FilterKeySizeClassification  = "size_classification"
)

// Allowed values for specific filter keys
var (
	AllowedSectors = []string{
		"Industrials",
		"Health Care",
		"Communication Services",
		"Information Technology",
		"Financials",
		"Consumer Discretionary",
		"Materials",
		"Consumer Staples",
		"Real Estate",
		"Energy",
		"Utilities",
	}

	AllowedStyleClassifications = []string{"Growth", "Core", "Value"}
	AllowedSizeClassifications  = []string{"Large Cap", "Small Cap", "Mid Cap"}
)

// base filter schema
var FilterSchema = &jsonschema.Schema{
	Schema: "http://json-schema.org/draft-07/schema#",
	Type:   "object",
	Properties: map[string]*jsonschema.Schema{
		"key": &jsonschema.Schema{
			Type: "string",
		},
		"clientResourceKey": &jsonschema.Schema{
			Type: "string",
		},
		"params": &jsonschema.Schema{
			Type: "object",
		},
		"max": &jsonschema.Schema{
			Description: "maximum range for numeric filter",
			Type:        "number",
		},
		"min": &jsonschema.Schema{
			Description: "maximum range for numeric filter",
			Type:        "number",
		},
		"values": &jsonschema.Schema{
			Description: "maximum range for numeric filter",
			Type:        "number",
		},
	},
	Required: []string{"key"},
}

var ScreenerMetrics = map[string][]map[string]any{
	"equity_filters": {
		// ISO2 country filter - enum type with array values
		IndicatorSchema(
			keyVal("description", "countries to filter security for, values: array of ISO 3166-1 alpha-2 country codes"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "iso2"),
				keyVal("values", IndicatorSchema(
					keyVal("type", "array"),
					keyVal("items", IndicatorSchema(
						keyVal("type", "string"),
						keyVal("pattern", "^[A-Z]{2}$"),
					)),
					keyVal("uniqueItems", true),
					keyVal("minItems", 1),
				)),
			)),
			keyVal("required", []string{"key", "values"}),
		),
		// Sector filter - enum type with predefined values
		IndicatorSchema(
			keyVal("description", "Equity sectors filter, values: array of defined values"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "t_sec"),
				keyVal("values", IndicatorSchema(
					keyVal("type", "array"),
					keyVal("items", IndicatorSchema(
						keyVal("type", "string"),
						keyVal("enum", AllowedSectors),
					)),
					keyVal("uniqueItems", true),
					keyVal("minItems", 1),
				)),
			)),
			keyVal("required", []string{"key", "values"}),
		),
		// Size classification filter - enum type with predefined values
		IndicatorSchema(
			keyVal("description", "size classification filter, values: array of defined values"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "size_classification"),
				keyVal("values", IndicatorSchema(
					keyVal("type", "array"),
					keyVal("items", IndicatorSchema(
						keyVal("type", "string"),
						keyVal("enum", AllowedSizeClassifications),
					)),
					keyVal("uniqueItems", true),
					keyVal("minItems", 1),
				)),
			)),
			keyVal("required", []string{"key", "values"}),
		),
		// Style classification filter - enum type with predefined values
		IndicatorSchema(
			keyVal("description", "Equity style filter, values: array of defined values"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "style_classification"),
				keyVal("values", IndicatorSchema(
					keyVal("type", "array"),
					keyVal("items", IndicatorSchema(
						keyVal("type", "string"),
						keyVal("enum", AllowedStyleClassifications),
					)),
					keyVal("uniqueItems", true),
					keyVal("minItems", 1),
				)),
			)),
			keyVal("required", []string{"key", "values"}),
		),
		// Market cap filter - numeric type with min/max
		IndicatorSchema(
			keyVal("description", "Market cap filter, in million currency unit"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "mkt"),
				keyVal("min", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "minimum market cap in million currency unit"),
				)),
				keyVal("max", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "maximum market cap in million currency unit"),
				)),
				keyValueDescription("currency", "3-letter currency code, default=USD", "USD"),
			)),
			keyVal("required", []string{"key", "min", "max"}),
		),
		// 1Y revenue growth filter - numeric type with min/max
		IndicatorSchema(
			keyVal("description", "1Y revenue growth filter, float format. Example: 10% - 0.1"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "REVENUE_EST_FY0_YOY_PCT"),
				keyVal("min", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "minimum revenue growth rate"),
				)),
				keyVal("max", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "maximum revenue growth rate"),
				)),
			)),
			keyVal("required", []string{"key", "min", "max"}),
		),
		// EV/EBITDA Ratio filter - numeric type with min/max and key enum
		IndicatorSchema(
			keyVal("description", "EV/EBITDA Ratio filter, trailing (LTM) or forward (NTM)"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyVal("key", IndicatorSchema(
					keyVal("type", "string"),
					keyVal("enum", []string{"evebitdaltm", "evebitda"}),
					keyVal("description", "evebitdaltm: last 12m EB/EBITDA ratio, evebitda: next 12m EB/EBITDA ratio"),
				)),
				keyVal("min", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "minimum EV/EBITDA ratio"),
				)),
				keyVal("max", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "maximum EV/EBITDA ratio"),
				)),
			)),
			keyVal("required", []string{"key", "min", "max"}),
		),
		// FCF margin filter - numeric type with min/max
		IndicatorSchema(
			keyVal("description", "FCF margin filter (LTM), unit (%)"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "pf_fcf_margin-LTM"),
				keyVal("min", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "minimum FCF margin percentage"),
				)),
				keyVal("max", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "maximum FCF margin percentage"),
				)),
			)),
			keyVal("required", []string{"key", "min", "max"}),
		),
		// EBITDA Margin filter - numeric type with min/max and optional params
		IndicatorSchema(
			keyVal("description", "EBITDA Margin filter (LTM), unit (%)"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "f_ebm"),
				keyVal("min", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "minimum EBITDA margin percentage"),
				)),
				keyVal("max", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "maximum EBITDA margin percentage"),
				)),
				keyVal("params", IndicatorSchema(
					keyVal("type", "object"),
					keyVal("properties", IndicatorSchema(
						keyValueDescription("period", "financial period", "LTM"),
					)),
				)),
			)),
			keyVal("required", []string{"key", "min", "max"}),
		),
		// Net income margin filter - numeric type with min/max and optional params
		IndicatorSchema(
			keyVal("description", "Net income margin filter, unit (%)"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "f_nm"),
				keyVal("min", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "minimum net income margin percentage"),
				)),
				keyVal("max", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "maximum net income margin percentage"),
				)),
				keyVal("params", IndicatorSchema(
					keyVal("type", "object"),
					keyVal("properties", IndicatorSchema(
						keyValueDescription("period", "financial period", "LTM"),
					)),
				)),
			)),
			keyVal("required", []string{"key", "min", "max"}),
		),
		// 1Y Weekly Beta filter - numeric type with min/max
		IndicatorSchema(
			keyVal("description", "1Y Weekly Beta, range (-1.0,1.0) considered for low sensitivity"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "f_bet1yr"),
				keyVal("min", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "minimum beta value"),
				)),
				keyVal("max", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "maximum beta value"),
				)),
				keyVal("params", IndicatorSchema(
					keyVal("type", "object"),
					keyVal("properties", IndicatorSchema(
						keyValueDescription("period", "financial period", "LTM"),
					)),
				)),
			)),
			keyVal("required", []string{"key", "min", "max"}),
		),
		// Percentage below 52-week high - numeric type with min/max
		IndicatorSchema(
			keyVal("description", "Percentage below 52-week high"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyValueDescription("key", "", "chg52wHighPct-adj"),
				keyVal("min", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "minimum distance below 52-week high, percentage in decimal form"),
				)),
				keyVal("max", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "maximum percentage below 52-week high, percentage in decimal form"),
				)),
			)),
			keyVal("required", []string{"key", "min", "max"}),
		),
		// Price change filter - numeric type with enum key and single description
		IndicatorSchema(
			keyVal("description", "Price change filter, float format. Example: 10% - 0.1"),
			keyVal("type", "object"),
			keyVal("properties", IndicatorSchema(
				keyVal("key", IndicatorSchema(
					keyVal("type", "string"),
					keyVal("enum", []string{"chgYTDPct_L", "chg1yPct_L", "chg1mPct_L", "chg1wPct_L"}),
					keyVal("description", "chgYTDPct_L: YTD price change, chg1yPct_L: 1Y price change, chg1mPct_L: 1 month price change, chg1wPct_L: 1 week price change"),
				)),
				keyVal("min", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "minimum price change, percentage in decimal form"),
				)),
				keyVal("max", IndicatorSchema(
					keyVal("type", "number"),
					keyVal("description", "maximum price change, percentage in decimal form"),
				)),
			)),
			keyVal("required", []string{"key", "min", "max"}),
		),
	},
}

type Facet struct {
	Key    string            `json:"key" jsonschema:"filter key for financial metrics or non-numeric characteristics"`
	Params map[string]string `json:"params,omitempty" jsonschema:"extra paramters for filter key"`
}

type Filter struct {
	Key               string            `json:"key" jsonschema:"filter key for financial metrics or non-numeric characteristics"`
	Params            map[string]string `json:"params,omitempty" jsonschema:"extra paramters for filter key"`
	ClientResourceKey string            `json:"clientResourceKey,omitempty" jsonschema:"client resouce key, use same value for key, optional"`
	Currency          string            `json:"currency,omitempty" jsonschema:"filter metric's currency"`
	Values            []string          `json:"values,omitempty" jsonschema:"characteristics filter for non-numeric key"`
	Max               *float64          `json:"max,omitempty" jsonschema:"maximum numeric range"`
	Min               *float64          `json:"min,omitempty" jsonschema:"minimum numeric range"`
}

// custom MarshalJSON method
// this is done to accommodate AI/MCP clients' nest JSON key bug
// flatten JSON accepting AI geneated tool calll -> uses custom marshaller to marshal to correct format
func (f *Filter) MarshalJSON() ([]byte, error) {
	convertedType := struct {
		Key               Facet    `json:"key"`
		ClientResourceKey string   `json:"clientResourceKey,omitempty"`
		Currency          string   `json:"currency,omitempty"`
		Values            []string `json:"values,omitempty"`
		Max               *float64 `json:"max,omitempty"`
		Min               *float64 `json:"min,omitempty"`
	}{
		Key:               Facet{Key: f.Key, Params: f.Params},
		ClientResourceKey: f.ClientResourceKey,
		Currency:          f.Currency,
		Values:            f.Values,
		Max:               f.Max,
		Min:               f.Min,
	}
	return json.Marshal(convertedType)
}

var DefaultPrimaryFilter = Filter{
	Key: "parent",
	Max: utils.Ptr(1.0),
	Min: utils.Ptr(1.0),
}

type Metric struct {
	Key  string `json:"key" jsonschema:"sort key"`
	Type string `json:"type" jsonschema:"sort key value type, default: numeric"`
}

type Order struct {
	Metric    Metric `json:"metric" jsonschema:"sorting metric"`
	Direction string `json:"direction" jsonschema:"sort direction, anyOf: DESC|ASC"`
}

func CreateOrder(key, direction string) Order {
	return Order{
		Direction: direction,
		Metric: Metric{
			Key:  key,
			Type: "numeric",
		},
	}
}

type ScreenCriteria struct {
	Conditions []Filter `json:"filters" jsonschema:"list of filter conditions"`
	OrderBy    Order    `json:"orderBy" jsonschema:"order by metric, only mkt metric allowed"`
	PageSize   uint32   `json:"pageSize" jsonschema:"page size, maximum 300"`
	// Cursor     string   `json:"cursor,omitempty" jsonschema:"pagination cursor from previous response"`
}

func (sc *ScreenCriteria) Validate(v *validator.Validator) {
	v.Check(len(sc.Conditions) > 0, "filters", "at least one filter condition is required")
	v.Check(sc.PageSize > 0, "pageSize", "page size must be greater than 0")
	v.Check(sc.PageSize <= 300, "pageSize", "page size must not exceed 300")

	// Validate order by direction
	v.Check(validator.In(sc.OrderBy.Direction, "ASC", "DESC"), "orderBy.direction", "must be one of ASC or DESC")

	// Validate order by metric key
	v.Check(sc.OrderBy.Metric.Key != "", "orderBy.metric.key", "metric key must not be empty")

	// Validate each filter condition
	for i, condition := range sc.Conditions {
		key := condition.Key
		v.Check(key != "", fmt.Sprintf("filters[%d].key", i), "filter key must not be empty")

		// At least one of Max, Min, or Values must be set
		hasRange := condition.Max != nil || condition.Min != nil
		hasValues := len(condition.Values) > 0
		v.Check(hasRange || hasValues, fmt.Sprintf("filters[%d]", i), "filter must have at least one of max, min, or values")

		// Validate Min <= Max if both are set
		if condition.Min != nil && condition.Max != nil {
			v.Check(*condition.Min <= *condition.Max, fmt.Sprintf("filters[%d].min", i), "must not be greater than max")
		}

		// Validate filter values are within allowed range
		for _, val := range condition.Values {
			v.Check(val != "", fmt.Sprintf("filters[%d].values", i), "filter values must not contain empty strings")
		}

		// Validate key-specific constraints
		validateFilterKeyValues(v, key, condition.Values, i)
	}
}

// validateFilterKeyValues validates values for specific filter keys
func validateFilterKeyValues(v *validator.Validator, key string, values []string, index int) {
	if len(values) == 0 {
		return
	}

	switch key {
	case FilterKeyISO2:
		for j, val := range values {
			v.Check(len(val) == 2, fmt.Sprintf("filters[%d].values[%d]", index, j), "iso2 country code must be 2 characters")
		}

	case FilterKeySector:
		v.Check(validator.IsSubset(values, AllowedSectors), fmt.Sprintf("filters[%d].values", index), "sector values must be valid sector values")

	case FilterKeyStyleClassification:
		v.Check(validator.IsSubset(values, AllowedStyleClassifications), fmt.Sprintf("filters[%d].values", index), "style classification values must be one of Growth, Core, Value")

	case FilterKeySizeClassification:
		v.Check(validator.IsSubset(values, AllowedSizeClassifications), fmt.Sprintf("filters[%d].values", index), "size classification values must be one of Large Cap, Small Cap, Mid Cap")
	}
}

type KIDString string

// Ideas: deserialize struct with single field -> value for direct extraction at client
func (k *KIDString) UnmarshalJSON(b []byte) error {
	tempStruct := struct {
		KID string `json:"kid"`
	}{}
	err := json.Unmarshal(b, &tempStruct)
	if err != nil {
		return err
	}
	*k = KIDString(tempStruct.KID)
	return nil
}

type ScreenResult struct {
	Continue string      `json:"continuation" jsonschema:"pagination cursor"`
	Count    uint32      `json:"count" jsonschema:"result count"`
	Kids     []KIDString `json:"hits" jsonschema:"list of KIDs for screened stocks"`
}

func (c *Client) ScreenForStocks(session *Session, req ScreenCriteria) (*ScreenResult, error) {
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
	// create response body
	err = json.NewEncoder(b).Encode(req)
	if err != nil {
		return nil, err
	}
	// println(b.String())

	// once access token validity is ensure - we can use directly without worrying about race condition
	headers := map[string][]string{
		"User-Agent":   {"Mozilla/5.0 (X11; Linux x86_64; rv:146.0) Gecko/20100101 Firefox/146.0"},
		"Accept":       {"application/json, text/plain, */*'"},
		"Content-Type": {"application/json"},
		"Origin":       {"https://app.koyfin.com"},
	}

	// populate authorization headers
	session.AuthorizeHeader(headers)

	body, err := c.getResponse("POST", "https://app.koyfin.com/api/v1/screener/query/", b, headers)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// parse response
	result := ScreenResult{}
	err = json.NewDecoder(body).Decode(&result)
	if err != nil {
		return nil, errors.Join(ErrBadResponse, err)
	}
	return &result, nil
}

// ScreenerSchemaRequest for retrieving screener filter schemas
type ScreenerSchemaRequest struct {
	AssetType string `json:"asset_type" jsonschema:"Financial asset type, currently only 'Equity' is supported"`
}

func (ssr *ScreenerSchemaRequest) Validate(v *validator.Validator) {
	v.Check(ssr.AssetType == "Equity", "asset_type", "must be 'Equity'")
}

// GetScreenerSchema returns available screener filters for the given asset type
func GetScreenerSchema(assetType string) []map[string]any {
	if assetType == "Equity" {
		return ScreenerMetrics["equity_filters"]
	}
	return nil
}
