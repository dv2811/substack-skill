package koyfin

import (
	"entext-applications/internal/validator"
)

type SchemaOption func(sc map[string]any)

func keyValueDescription(key, description string, value any) SchemaOption {
	return func(sc map[string]any) {
		// define description if needed
		var data map[string]any
		switch v := value.(type) {
		case string:
			data = map[string]any{
				"type":  "string",
				"const": v,
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			data = map[string]any{
				"type":  "number",
				"const": v,
			}
		case bool:
			data = map[string]any{
				"type":  "boolean",
				"const": v,
			}
		case []string:
			data = map[string]any{
				"type":  "string",
				"anyOf": v,
			}
		case map[string]any:
			data = v
		}
		if data != nil {
			if description != "" {
				data["description"] = description
			}
			sc[key] = data
		}
	}
}

func keyVal(key string, value any) SchemaOption {
	return func(sc map[string]any) {
		sc[key] = value
	}
}

func IndicatorSchema(opts ...SchemaOption) map[string]any {
	schema := map[string]any{}
	for _, opt := range opts {
		opt(schema)
	}
	return schema
}

var EquityIndicatorSchema = map[string][]map[string]any{
	"financials": {
		IndicatorSchema(
			keyVal("description", "Return on Invested Capital (ROIC), trailing 12m"),
			keyVal("key", "pf_roic_fin_mthd"),
			keyValueDescription("financialPeriodType", "", "LTM"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Total Revenues, CAGR"),
			keyVal("key", "f_revg1y"),
			keyValueDescription("financialPeriodType", "", "LTM"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Capital Expenditure: financial period options: LTM (last twelve months), quarterly (quarterly period)"),
			keyVal("key", "f_capex"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "quarterly"}),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "EBITDA: financial period options: LTM (last twelve months), quarterly (quarterly period)"),
			keyVal("key", "f_ebitda"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "quarterly"}),
		),
		IndicatorSchema(
			keyVal("description", "Operating Income: financial period options: LTM (last twelve months), quarterly (quarterly period)"),
			keyVal("key", "f_opinc"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "quarterly"}),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Total Revenues: financial period options: LTM (last twelve months), quarterly (quarterly period)"),
			keyVal("key", "f_r"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "quarterly"}),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "R&D Expenses: financial period options: LTM (last twelve months), quarterly (quarterly period)"),
			keyVal("key", "f_rd"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "quarterly"}),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Stock-Based Compensation (CF): financial period options: LTM (last twelve months), quarterly (quarterly period)"),
			keyVal("key", "f_stkcomp"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "quarterly"}),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Free Cash Flow: financial period options: LTM (last twelve months), quarterly (quarterly period)"),
			keyVal("key", "pf_fcf"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "quarterly"}),
			keyVal("priceFormat", "standard"),
		),
	},
	"forward_estimates": []map[string]any{
		IndicatorSchema(
			keyVal("description", "EPS Estimate 5Y Growth Rate NTM"),
			keyVal("key", "fest_est_eps_growth_5y"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "CAPEX Consensus Average 1st Unreported Financial Year"),
			keyVal("key", "fest_estcapex_avg"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyValueDescription("financialForwardPeriod", "", "fy0"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Estimate: Revenues Consensus Median"),
			keyVal("key", "fest_estsales_median"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyValueDescription("financialForwardPeriod", "forwards period 1st, 2nd unreported FY", []string{"fy0", "fy1"}),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Estimate: Free Cash Flow Median"),
			keyVal("key", "fest_estfcf_median"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyValueDescription("financialForwardPeriod", "forwards period 1st, 2nd unreported FY", []string{"fy0", "fy1"}),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "EBIT Consensus Average - Next 12m"),
			keyVal("key", "fest_estebit_ntm"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Estimate: Revenues Consensus Median (NTM)"),
			keyVal("key", "fest_estsales_median_ntm"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "EBITDA Consensus Average"),
			keyVal("key", "fest_estebitda_ntm"),
			keyVal("currency", IndicatorSchema(
				keyVal("format", "3-letter currency code, default=USD"),
				keyVal("type", "string"),
			)),
			keyVal("priceFormat", "standard"),
		),
	},
	"market_data": []map[string]any{
		IndicatorSchema(
			keyVal("description", "Short Interest"),
			keyVal("key", "f_sip"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Historical Price"),
			keyVal("key", "p_candle_range"),
			keyValueDescription("priceFormat", "", "adj"),
		),
	},
	"ratios": []map[string]any{
		IndicatorSchema(
			keyVal("description", "Capex as % of Revenues: financial period options: LTM (last twelve months), annual (fiscal year), quarterly (quarterly period)"),
			keyVal("key", "f_capexrev"),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "annual", "quarterly"}),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "EBIT Margin %: financial period options: LTM (last twelve months), annual (fiscal year), quarterly (quarterly period)"),
			keyVal("key", "f_em"),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "annual", "quarterly"}),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "EV / Gross Profit"),
			keyVal("key", "f_ev_gp"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "EV / EBITDA Next 12m"),
			keyVal("key", "f_evebitda"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "EV / EBITDA Trailing 12m"),
			keyVal("key", "f_evebitdaltm"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Free Cash Flow / EV Yield"),
			keyVal("key", "f_evfcf_yld"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "EV / Sales Next 12m"),
			keyVal("key", "f_evs"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "EV / Sales Trailing 12m"),
			keyVal("key", "f_evsltm"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Gross Profit Margin %: financial period options: LTM (last twelve months), quarterly (quarterly period)"),
			keyVal("key", "f_gma"),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "quarterly"}),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Price/Earnings Ratio, next 12m"),
			keyVal("key", "f_pe"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Price/Earnings to Growth"),
			keyVal("key", "f_peg"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Price/Earnings Ratio, last 12m"),
			keyVal("key", "f_peltm"),
			keyVal("priceFormat", "standard"),
		),
		IndicatorSchema(
			keyVal("description", "Free Cash Flow Margin %: financial period options: LTM (last twelve months), quarterly (quarterly period)"),
			keyVal("key", "pf_fcf_margin"),
			keyValueDescription("financialPeriodType", "", []string{"LTM", "quarterly"}),
			keyVal("priceFormat", "standard"),
		),
	},
}

type SchemaRequest struct {
	AssetType     string `json:"asset_type" jsonschema:"Financial assets types"`
	IndicatorType string `json:"indicator_type" jsonschema:"subset of indicator for the given asset type, must be one of financials | ratios | forward_estimates | market_data for Equity"`
}

// consumer of this package must create a validator and call this method of SchemaRequest by themselves
// *SchemaRequest passed to GetAvailableSchema is assumed to have been checked for validation
func (scr *SchemaRequest) Validate(v *validator.Validator) {
	// validation logics where
	// switch case for Equity type
	v.Check(validator.In(scr.AssetType, "Equity"), "asset_type", "must be 'Equity'")

	if scr.AssetType == "Equity" {
		// Validate indicator type for Equity
		v.Check(validator.In(scr.IndicatorType, "financials", "ratios", "forward_estimates", "market_data"),
			"indicator_type", "must be one of financials | ratios | forward_estimates | market_data for Equity")
	}
}

// GetAvailableSchema retrieve available schema for an asset type
func GetAvailableSchema(scr *SchemaRequest) []map[string]any {
	// switch case for equity asset type here
	switch scr.AssetType {
	case "Equity":
		return EquityIndicatorSchema[scr.IndicatorType]
	default:
		return nil // Return nil for unsupported asset types
	}
}
