package koyfin

import (
	"fmt"
	"net/url"
	"strconv"

	"entext-applications/internal/validator"
	"github.com/goccy/go-json"
)

// TranscriptListRequest represents the request for listing transcripts
type TranscriptListRequest struct {
	KID   string `json:"kid" description:"Koyfin identifier for the stock"`
	Limit int    `json:"limit" description:"Maximum number of results to return,maximum 64"`
}

func (req *TranscriptListRequest) Validate(v *validator.Validator) {
	v.Check(req.KID != "", "kid", "koyfin identifier must not be empty")
	v.Check(req.Limit < 64 && req.Limit > 0, "limit", "search limit must be between 0 and 64")
}

// Transcript represents a single transcript entry with only the fields we care about
type Transcript struct {
	AnnouncedDate   string `json:"announcedDate" description:"Date when event was announced"`
	EventDateTime   string `json:"eventDateTime" description:"Date and time of the event"`
	TranscriptTitle string `json:"transcriptTitle" description:"Title of the transcript"`
	CompanyName     string `json:"companyName,omitempty" description:"Name of the company"`
	Reference       string `json:"reference,omitempty" description:"link to the transcript on Koyfin"`
	// CompanyId cannot be traced to KID or used in other API so we ignore it in ouput for now
	// CompanyId       int    `json:"companyId" description:"Company identifier"`
	KeyDevId      int  `json:"keyDevId" description:"Key development identifier"`
	FiscalQuarter *int `json:"fiscalQuarter,omitempty" description:"Fiscal quarter of the event"`
	FiscalYear    *int `json:"fiscalYear,omitempty" description:"Fiscal year of the event"`
}

// TranscriptDetail represents the detailed transcript response
type TranscriptDetail struct {
	Header     Transcript `json:"header" description:"Header information for the transcript"`
	Components []Dialogue `json:"components" description:"List of transcript dialogues"`
}

// TranscriptComponent represents a single component of a transcript
type Dialogue struct {
	ComponentOrder int    `json:"componentOrder" description:"Order of the component in the transcript"`
	SpeakerName    string `json:"speakerName" description:"Name of the speaker"`
	SpeakerType    string `json:"speakerType" description:"Type of the speaker"`
	Text           string `json:"text" description:"Text content of the component"`
}

// TranscriptSummary represents the AI summary response
type TranscriptSummary struct {
	KeyDevId      int     `json:"keyDevId" description:"Key development identifier"`
	TranscriptId  int     `json:"transcriptId" description:"Transcript identifier"`
	EventDateTime string  `json:"eventDateTime" description:"Date and time of the event"`
	SummaryMd     string  `json:"summaryMd" description:"Markdown formatted summary"`
	ComparisonMd  *string `json:"comparisonMd" description:"Markdown formatted comparison"`
}

// ListTranscripts retrieves a list of recent transcripts for a given ticker
func (c *Client) ListTranscripts(session *Session, req *TranscriptListRequest) ([]Transcript, error) {
	// ensure access token validity
	err := c.ensureValidToken(session)
	if err != nil {
		return nil, err
	}

	// prepare headers
	headers := map[string][]string{
		"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64; rv:147.0) Gecko/20100101 Firefox/147.0"},
		"Accept":     {"application/json, text/plain, */*"},
	}

	// populate authorization headers
	session.AuthorizeHeader(headers)

	// construct URL with proper encoding
	baseURL := fmt.Sprintf("https://app.koyfin.com/api/v1/pubhub/transcript/list/%s", url.PathEscape(req.KID))
	params := url.Values{}
	params.Add("limit", strconv.Itoa(req.Limit))
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// make request
	body, err := c.getResponse("GET", fullURL, nil, headers)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// decode response
	var transcripts []Transcript
	err = json.NewDecoder(body).Decode(&transcripts)
	if err != nil {
		return nil, err
	}

	// populate reference so that we can link end user to display page
	for indx := range transcripts {
		tID := transcripts[indx].KeyDevId
		transcripts[indx].Reference = fmt.Sprintf("https://app.koyfin.com/news/ts/%s/all/%d?sourceType=transcript", req.KID, tID)
	}
	return transcripts, nil
}

// GetTranscript retrieves a specific transcript by keyDevId
func (c *Client) GetTranscript(session *Session, keyDevId int) (*TranscriptDetail, error) {
	// ensure access token validity
	err := c.ensureValidToken(session)
	if err != nil {
		return nil, err
	}

	// prepare headers
	headers := map[string][]string{
		"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64; rv:147.0) Gecko/20100101 Firefox/147.0"},
		"Accept":     {"application/json, text/plain, */*"},
	}

	// populate authorization headers
	session.AuthorizeHeader(headers)

	// construct URL with proper encoding
	transcriptURL := fmt.Sprintf("https://app.koyfin.com/api/v1/pubhub/v2/transcript/%d", keyDevId)

	// make request
	body, err := c.getResponse("GET", transcriptURL, nil, headers)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// decode response
	var transcriptDetail TranscriptDetail
	err = json.NewDecoder(body).Decode(&transcriptDetail)
	if err != nil {
		return nil, err
	}

	return &transcriptDetail, nil
}

// GetTranscriptSummary retrieves the AI summary for a specific transcript by keyDevId
func (c *Client) GetTranscriptSummary(session *Session, keyDevId int) (*TranscriptSummary, error) {
	// ensure access token validity
	err := c.ensureValidToken(session)
	if err != nil {
		return nil, err
	}

	// prepare headers
	headers := map[string][]string{
		"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64; rv:147.0) Gecko/20100101 Firefox/147.0"},
		"Accept":     {"application/json, text/plain, */*"},
	}

	// populate authorization headers
	session.AuthorizeHeader(headers)

	// construct URL with proper encoding
	summaryURL := fmt.Sprintf("https://app.koyfin.com/api/doc-ai/v1/transcript/summary/%d", keyDevId)

	// make request
	body, err := c.getResponse("GET", summaryURL, nil, headers)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// decode response
	var summary TranscriptSummary
	err = json.NewDecoder(body).Decode(&summary)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
