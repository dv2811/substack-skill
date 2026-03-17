package koyfin

import (
	"bytes"
	"errors"

	"entext-applications/internal/validator"
	"github.com/goccy/go-json"
)

var (
	ErrAssetTypeNotSupported = errors.New("asset type not supported")
	ErrNoSnapshotData        = errors.New("no snapshot data")
	SnapShotIDTypes          = []string{"KID", "holdings"}
)

type Params map[string]any

func NewParameter(keyValues ...string) Params {
	prms := Params{}
	for i := 0; i < len(keyValues)-1; i += 2 {
		prms[keyValues[i]] = keyValues[i+1]
	}
	return prms
}

type IndicatorKey struct {
	Key    string `json:"key" jsonschema:"Indicator key"`
	Alias  string `json:"alias,omitempty" jsonschema:"Indicator alias"`
	Params Params `json:"params,omitempty" jsonschema:"Parameters for the indicator"`
}

func CreateIndicator(key, alias string, params Params) IndicatorKey {
	return IndicatorKey{Key: key, Alias: alias, Params: params}
}

type SnapshotID struct {
	KID    string `json:"id"`
	IDType string `json:"type"`
}

type SnapshotDataRequest struct {
	KIDs []SnapshotID   `json:"ids"`
	Keys []IndicatorKey `json:"keys,omitempty"`
}

// set default keys for all asset types instead of
func DefaultKeysForAssetType(req SnapshotRequest) []IndicatorKey {
	switch {
	case req.Category == "Equity":
		return []IndicatorKey{
			CreateIndicator("p_l", "current_price_adj", NewParameter("priceFormat", "adj")),
			// CreateIndicator("p_l", "unadj_price", nil),
			CreateIndicator("f_r", "revenue_last_12m", NewParameter("period", "LTM")),
			CreateIndicator("f_opinc", "op_income_last_12m", NewParameter("period", "LTM")),
			CreateIndicator("f_ni", "net_income_ltm", NewParameter("period", "LTM")),
			CreateIndicator("f_ebitda_incl", "ebitda_last_12m", NewParameter("period", "LTM")),
			CreateIndicator("p_c3y", "prev_3y_price_adj", NewParameter("priceFormat", "adj")),
			CreateIndicator("p_c1y", "prev_1y_price_adj", NewParameter("priceFormat", "adj")),
			CreateIndicator("p_cytd", "start_year_price_adj", NewParameter("priceFormat", "adj")),
			CreateIndicator("f_sip", "short_interest_pct", NewParameter("priceFormat", "standard")),
			CreateIndicator("f_bet1yr", "1y_wkly_beta", NewParameter("period", "LTM")),
			CreateIndicator("f_mkt", "market_cap_mil", nil),
			CreateIndicator("f_cf", "cashflow_last_12m", NewParameter("period", "LTM")),
			CreateIndicator("f_eps", "eps_last_12m", nil), // EPS /
			CreateIndicator("fest_est_eps_growth_5y", "est_5y_cagr_eps", nil),
			CreateIndicator("fest_estepsgaap_ntm", "ntm_eps_gaap", nil),
			CreateIndicator("f_ev", "current_ev", nil),

			// CreateIndicator("fest_estcapex_avg", "CAPEX Consensus Average 1st Unreported Financial Year", nil),
			CreateIndicator("fest_estebitda_ntm", "ebitda_next_12m", nil),
			CreateIndicator("fest_estsales_ntm", "revenue_next_12m", nil),
			CreateIndicator("fest_estebit_ntm", "ebit_next_12m", nil),
			CreateIndicator("fest_esteps_ntm", "eps_next_12m", nil),
			CreateIndicator("f_compSet", "competitors", nil),
			// metadata
			CreateIndicator("t_n", "name", nil),
			CreateIndicator("t_id", "ticker", nil),
			CreateIndicator("t_sec", "sector", nil),
			CreateIndicator("t_ind", "industry", nil),
		}
	case req.Category == "ETF":
		return []IndicatorKey{
			CreateIndicator("p_l", "current_price_adj", NewParameter("priceFormat", "adj")),
			CreateIndicator("p_c3y", "prev_3y_price_adj", NewParameter("priceFormat", "adj")),
			CreateIndicator("p_c1y", "prev_1y_price_adj", NewParameter("priceFormat", "adj")),
			CreateIndicator("p_cytd", "start_year_price_adj", NewParameter("priceFormat", "adj")),
			CreateIndicator("f_sip", "short_interest_pct", NewParameter("priceFormat", "standard")),
			CreateIndicator("f_mkt", "market_cap_mil", nil),
			CreateIndicator("f_bet1yr", "1y_wkly_beta", NewParameter("period", "LTM")),
			CreateIndicator("fest_esteps_ntm", "eps_next_12m", nil),
			// CreateIndicator("p_l", "unadj_price", nil),
			CreateIndicator("fest_estebitda_ntm", "ebitda_next_12m", nil),
			CreateIndicator("f_eps", "eps_last_12m", nil), // EPS /
			CreateIndicator("f_ebit", "ebit_last_12m", NewParameter("period", "LTM")),
			CreateIndicator("f_ebm", "ebita_margin_ltm", NewParameter("period", "LTM")),
			CreateIndicator("f_cf", "cashflow_last_12m", NewParameter("period", "LTM")),
			CreateIndicator("f_r", "revenue_last_12m", NewParameter("period", "LTM")),
			CreateIndicator("fest_estsales_ntm", "revenue_next_12m", nil),
			CreateIndicator("fest_estepsgaap_ntm", "ntm_eps_gaap", nil),
			CreateIndicator("f_ev", "current_ev", nil),
			// metadata
			CreateIndicator("t_n", "name", nil),
			CreateIndicator("t_id", "ticker", nil),
			CreateIndicator("t_sec", "sector", nil),
			CreateIndicator("t_ind", "industry", nil),
		}
	}
	return nil
}

type SnapshotRequest struct {
	KIDs     []string `json:"kids" jsonschema:"list of Koyfin IDs, maximum: 32 KIDs for equity, 2 for ETF holdings request"`
	Category string   `json:"category" jsonschema:"category of the instrument"`
}

// validate snapshot request
func (s *SnapshotRequest) Validate(v *validator.Validator) {
	v.Check(len(s.KIDs) > 0, "ids", "ID list must not be empty")
	if s.Category == "Equity" {
		v.Check(len(s.KIDs) < 33, "ids", "only a maximum of 32 Koyfin IDs per request")
	} else {
		v.Check(len(s.KIDs) < 3, "ids", "only a maximum of 2 Koyfin IDs per request")
	}
	v.Check(s.Category != "", "category", "category must not be empty")
}

type IndicatorData struct {
	Date     string       `json:"date,omitempty" jsonschema:"Date of the indicator"`
	Label    string       `json:"label,omitempty" jsonschema:"Data label"`
	Value    float64      `json:"value,omitempty" jsonschema:"Value of the indicator"`
	Currency string       `json:"currency,omitempty" jsonschema:"Currency of the indicator"`
	Error    *KoyfinError `json:"error,omitempty" jsonschema:"Error information"`
}

type StockMetrics map[string]any

type StockSnapshotResp struct {
	KID     string       `json:"kids" jsonschema:"koyfin identifier for the instrument"`
	Metrics StockMetrics `json:"metrics" jsonschema:"financial metrics for a given security"`
}

// GetSnapshotData shows data snapshot for a given instrument using pre-defined schema
func (c *Client) GetSnapshotData(session *Session, req SnapshotRequest) ([]StockSnapshotResp, error) {
	snapshotIDs := make([]SnapshotID, len(req.KIDs))
	for indx, kid := range req.KIDs {
		snapshotIDs[indx] = SnapshotID{KID: kid, IDType: "KID"}
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

	output := struct {
		KID map[string]StockMetrics `json:"KID"`
	}{}

	err = json.NewDecoder(body).Decode(&output)
	if err != nil {
		return nil, err
	}

	if len(output.KID) == 0 {
		return nil, ErrNoSnapshotData
	}

	var snapshots []StockSnapshotResp

	// return only one series data
	for kid, data := range output.KID {
		snapshots = append(snapshots, StockSnapshotResp{KID: kid, Metrics: data})
	}
	return snapshots, nil
}
