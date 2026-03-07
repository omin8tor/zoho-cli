package zohttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/omin8tor/zoho-cli/internal"
	"github.com/omin8tor/zoho-cli/internal/auth"
	"github.com/omin8tor/zoho-cli/internal/dc"
)

type Client struct {
	Config              *auth.AuthConfig
	HTTP                *http.Client
	accessToken         string
	AnalyticsBase       string
	AssistBase          string
	BackstageBase       string
	BiginBase           string
	BillingBase         string
	BookingsBase        string
	BooksBase           string
	CampaignsBase       string
	CliqBase            string
	CRMBase             string
	CreatorBase         string
	DeskBase            string
	DownloadBase        string
	ExpenseBase         string
	InventoryBase       string
	InvoiceBase         string
	LearnBase           string
	MailBase            string
	MarketingAutoBase   string
	MeetingBase         string
	PageSenseBase       string
	PeopleBase          string
	ProjectsBase        string
	RecruitBase         string
	SalesIQBase         string
	SheetBase           string
	ShowtimeBase        string
	SignBase            string
	SprintsBase         string
	VaultBase           string
	VoiceBase           string
	WorkDriveBase       string
	WriterBase          string
}

func NewClient(config *auth.AuthConfig) (*Client, error) {
	token, err := auth.EnsureAccessToken(config, false)
	if err != nil {
		return nil, err
	}

	d := config.DC
	api := dc.APIURL(d)
	return &Client{
		Config:            config,
		HTTP:              &http.Client{Timeout: 60 * time.Second},
		accessToken:       token,
		AnalyticsBase:     api + "/analytics/v2",
		AssistBase:        api + "/assist/v1",
		BackstageBase:     api + "/backstage/v1",
		BiginBase:         api + "/bigin/v2",
		BillingBase:       api + "/billing/v1",
		BookingsBase:      api + "/bookings/v1",
		BooksBase:         api + "/books/v3",
		CampaignsBase:     api + "/campaigns/v1",
		CliqBase:          dc.CliqURL(d),
		CRMBase:           dc.CRMURL(d) + "/crm/v8",
		CreatorBase:       api + "/creator/v2.1",
		DeskBase:          dc.DeskURL(d) + "/api/v1",
		DownloadBase:      dc.DownloadURL(d),
		ExpenseBase:       api + "/expense/v1",
		InventoryBase:     api + "/inventory/v1",
		InvoiceBase:       api + "/invoice/v3",
		LearnBase:         api + "/learn/v1",
		MailBase:          dc.MailURL(d),
		MarketingAutoBase: api + "/marketingautomation/v1",
		MeetingBase:       api + "/meeting/v1",
		PageSenseBase:     api + "/pagesense/v1",
		PeopleBase:        dc.PeopleURL(d) + "/people/api",
		ProjectsBase:      dc.ProjectsURL(d) + "/api/v3",
		RecruitBase:       api + "/recruit/v2",
		SalesIQBase:       api + "/salesiq/v2",
		SheetBase:         dc.SheetURL(d) + "/api/v2",
		ShowtimeBase:      api + "/showtime/v1",
		SignBase:          dc.SignURL(d) + "/api/v1",
		SprintsBase:       dc.SprintsURL(d) + "/zsapi",
		VaultBase:         api + "/vault/v1",
		VoiceBase:         api + "/voice/v1",
		WorkDriveBase:     dc.WorkDriveURL(d) + "/api/v1",
		WriterBase:        dc.WriterURL(d) + "/api/v1",
	}, nil
}

type RequestOpts struct {
	Params  map[string]string
	JSON    any
	Form    map[string]string
	Files   map[string]FileUpload
	Headers map[string]string
}

type FileUpload struct {
	Filename string
	Data     []byte
}

func (c *Client) Request(method, rawURL string, opts *RequestOpts) (json.RawMessage, error) {
	body, err := c.doRequest(method, rawURL, opts)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return json.RawMessage("{}"), nil
	}
	return json.RawMessage(body), nil
}

func (c *Client) RequestRaw(method, rawURL string, params map[string]string) ([]byte, http.Header, int, error) {
	req, err := c.buildRequest(method, rawURL, &RequestOpts{Params: params})
	if err != nil {
		return nil, nil, 0, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		if err := c.handleRetry(resp); err != nil {
			return nil, nil, 0, err
		}
		req, _ = c.buildRequest(method, rawURL, &RequestOpts{Params: params})
		resp, err = c.HTTP.Do(req)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == 401 {
			body, _ := io.ReadAll(resp.Body)
			return nil, nil, resp.StatusCode, internal.NewAuthError(fmt.Sprintf("Access token expired or invalid after refresh: %s", body))
		}
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, resp.StatusCode, internal.NewAPIError(resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	return body, resp.Header, resp.StatusCode, err
}

func (c *Client) doRequest(method, rawURL string, opts *RequestOpts) ([]byte, error) {
	req, err := c.buildRequest(method, rawURL, opts)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		if err := c.handleRetry(resp); err != nil {
			return nil, err
		}
		req, _ = c.buildRequest(method, rawURL, opts)
		resp, err = c.HTTP.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == 401 {
			body, _ := io.ReadAll(resp.Body)
			return nil, internal.NewAuthError(fmt.Sprintf("Access token expired or invalid after refresh: %s", body))
		}
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, internal.NewAPIError(resp.StatusCode, string(body))
	}

	if resp.StatusCode == 204 {
		return nil, nil
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) handleRetry(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if strings.Contains(bodyStr, "scope_invalid") || strings.Contains(bodyStr, "scope_mismatch") ||
		strings.Contains(bodyStr, "OAUTH_SCOPE_MISMATCH") {
		return internal.NewAuthError(fmt.Sprintf("OAuth scope insufficient — re-authorize with correct scopes: %s", bodyStr))
	}
	token, err := auth.EnsureAccessToken(c.Config, true)
	if err != nil {
		return err
	}
	c.accessToken = token
	return nil
}

func (c *Client) buildRequest(method, rawURL string, opts *RequestOpts) (*http.Request, error) {
	if opts == nil {
		opts = &RequestOpts{}
	}

	if len(opts.Params) > 0 {
		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		for k, v := range opts.Params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		rawURL = u.String()
	}

	var bodyReader io.Reader
	var contentType string

	if len(opts.Files) > 0 {
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		for key, file := range opts.Files {
			part, err := w.CreateFormFile(key, file.Filename)
			if err != nil {
				return nil, err
			}
			part.Write(file.Data)
		}
		for k, v := range opts.Form {
			w.WriteField(k, v)
		}
		w.Close()
		bodyReader = &buf
		contentType = w.FormDataContentType()
	} else if opts.JSON != nil {
		data, err := json.Marshal(opts.JSON)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
		if contentType == "" {
			contentType = "application/json"
		}
	} else if len(opts.Form) > 0 {
		vals := url.Values{}
		for k, v := range opts.Form {
			vals.Set(k, v)
		}
		bodyReader = strings.NewReader(vals.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	req, err := http.NewRequest(method, rawURL, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Zoho-oauthtoken "+c.accessToken)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "*/*")
	}

	return req, nil
}
