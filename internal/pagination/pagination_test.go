package pagination

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
)

func testClient(url string) *zohttp.Client {
	return &zohttp.Client{
		HTTP: http.DefaultClient,
	}
}

func TestHasNextPage(t *testing.T) {
	tests := []struct {
		input map[string]any
		want  bool
	}{
		{map[string]any{"has_next_page": true}, true},
		{map[string]any{"has_next_page": false}, false},
		{map[string]any{"has_next_page": "true"}, true},
		{map[string]any{"has_next_page": "false"}, false},
		{map[string]any{}, false},
		{nil, false},
	}
	for _, tt := range tests {
		got := hasNextPage(tt.input)
		if got != tt.want {
			t.Errorf("hasNextPage(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestExtractItems(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		key   string
		count int
	}{
		{
			name:  "top-level array",
			raw:   `[{"id":1},{"id":2}]`,
			key:   "",
			count: 2,
		},
		{
			name:  "simple key",
			raw:   `{"contacts":[{"id":1},{"id":2},{"id":3}]}`,
			key:   "contacts",
			count: 3,
		},
		{
			name:  "data key",
			raw:   `{"data":[{"id":1}],"info":{}}`,
			key:   "data",
			count: 1,
		},
		{
			name:  "dot-path nested",
			raw:   `{"response":{"result":[{"id":1},{"id":2}],"count":2}}`,
			key:   "response.result",
			count: 2,
		},
		{
			name:  "dot-path deeper",
			raw:   `{"a":{"b":{"c":[{"x":1}]}}}`,
			key:   "a.b.c",
			count: 1,
		},
		{
			name:  "missing key returns nil",
			raw:   `{"contacts":[{"id":1}]}`,
			key:   "invoices",
			count: 0,
		},
		{
			name:  "empty array",
			raw:   `{"data":[]}`,
			key:   "data",
			count: 0,
		},
		{
			name:  "not an array returns nil",
			raw:   `{"data":"hello"}`,
			key:   "data",
			count: 0,
		},
		{
			name:  "invalid json returns nil",
			raw:   `not json`,
			key:   "data",
			count: 0,
		},
		{
			name:  "empty key with non-array returns nil",
			raw:   `{"data":[1]}`,
			key:   "",
			count: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := ExtractItems(json.RawMessage(tt.raw), tt.key)
			if len(items) != tt.count {
				t.Errorf("ExtractItems(%s, %q) got %d items, want %d", tt.raw, tt.key, len(items), tt.count)
			}
		})
	}
}

func TestHasMoreBooks(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    bool
		wantPgN int
	}{
		{
			name:    "has more",
			raw:     `{"contacts":[],"page_context":{"has_more_page":true,"page":1}}`,
			want:    true,
			wantPgN: 2,
		},
		{
			name: "no more",
			raw:  `{"contacts":[],"page_context":{"has_more_page":false,"page":3}}`,
			want: false,
		},
		{
			name: "missing page_context",
			raw:  `{"contacts":[]}`,
			want: false,
		},
		{
			name: "invalid json",
			raw:  `not json`,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, state := HasMoreBooks(json.RawMessage(tt.raw), 0, 0)
			if got != tt.want {
				t.Errorf("HasMoreBooks got %v, want %v", got, tt.want)
			}
			if tt.want && state != nil && state.PageNum != tt.wantPgN {
				t.Errorf("HasMoreBooks PageNum = %d, want %d", state.PageNum, tt.wantPgN)
			}
		})
	}
}

func TestHasMoreCRM(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		want      bool
		wantToken string
	}{
		{
			name:      "more with token",
			raw:       `{"data":[],"info":{"more_records":true,"next_page_token":"abc123"}}`,
			want:      true,
			wantToken: "abc123",
		},
		{
			name: "more without token",
			raw:  `{"data":[],"info":{"more_records":true,"next_page_token":""}}`,
			want: true,
		},
		{
			name: "no more",
			raw:  `{"data":[],"info":{"more_records":false}}`,
			want: false,
		},
		{
			name: "invalid json",
			raw:  `not json`,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, state := HasMoreCRM(json.RawMessage(tt.raw), 0, 0)
			if got != tt.want {
				t.Errorf("HasMoreCRM got %v, want %v", got, tt.want)
			}
			if tt.wantToken != "" {
				if state == nil || state.Token != tt.wantToken {
					token := ""
					if state != nil {
						token = state.Token
					}
					t.Errorf("HasMoreCRM token = %q, want %q", token, tt.wantToken)
				}
			}
		})
	}
}

func TestHasMoreProjects(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want bool
	}{
		{
			name: "has next true bool",
			raw:  `{"tasks":[],"page_info":{"has_next_page":true}}`,
			want: true,
		},
		{
			name: "has next true string",
			raw:  `{"tasks":[],"page_info":{"has_next_page":"true"}}`,
			want: true,
		},
		{
			name: "has next false",
			raw:  `{"tasks":[],"page_info":{"has_next_page":false}}`,
			want: false,
		},
		{
			name: "no page_info",
			raw:  `{"tasks":[]}`,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := HasMoreProjects(json.RawMessage(tt.raw), 0, 0)
			if got != tt.want {
				t.Errorf("HasMoreProjects got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasMoreWorkDrive(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want bool
	}{
		{
			name: "has next",
			raw:  `{"data":[],"meta":{"has_next":true}}`,
			want: true,
		},
		{
			name: "no next",
			raw:  `{"data":[],"meta":{"has_next":false}}`,
			want: false,
		},
		{
			name: "missing meta",
			raw:  `{"data":[]}`,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := HasMoreWorkDrive(json.RawMessage(tt.raw), 0, 0)
			if got != tt.want {
				t.Errorf("HasMoreWorkDrive got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasMoreByCount(t *testing.T) {
	tests := []struct {
		fetched  int
		pageSize int
		want     bool
	}{
		{100, 100, true},
		{99, 100, false},
		{0, 100, false},
		{200, 100, true},
	}
	for _, tt := range tests {
		got, _ := HasMoreByCount(nil, tt.fetched, tt.pageSize)
		if got != tt.want {
			t.Errorf("HasMoreByCount(%d, %d) = %v, want %v", tt.fetched, tt.pageSize, got, tt.want)
		}
	}
}

func TestHasMoreSign(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want bool
	}{
		{
			name: "has more rows",
			raw:  `{"requests":[],"page_context":{"has_more_rows":true}}`,
			want: true,
		},
		{
			name: "no more rows",
			raw:  `{"requests":[],"page_context":{"has_more_rows":false}}`,
			want: false,
		},
		{
			name: "missing page_context",
			raw:  `{"requests":[]}`,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := HasMoreSign(json.RawMessage(tt.raw), 0, 0)
			if got != tt.want {
				t.Errorf("HasMoreSign got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetPageFunctions(t *testing.T) {
	t.Run("PagePerPage", func(t *testing.T) {
		params := map[string]string{}
		state := &PageState{PageNum: 3}
		PagePerPage(50)(state, params)
		if params["page"] != "3" || params["per_page"] != "50" {
			t.Errorf("PagePerPage: got page=%s per_page=%s", params["page"], params["per_page"])
		}
	})

	t.Run("SetPageCRM with token", func(t *testing.T) {
		params := map[string]string{"page": "1"}
		state := &PageState{Token: "tok123"}
		SetPageCRM(state, params)
		if params["page_token"] != "tok123" {
			t.Errorf("SetPageCRM: page_token = %q", params["page_token"])
		}
		if _, ok := params["page"]; ok {
			t.Error("SetPageCRM: page should be deleted when using token")
		}
	})

	t.Run("SetPageCRM without token", func(t *testing.T) {
		params := map[string]string{}
		state := &PageState{PageNum: 2}
		SetPageCRM(state, params)
		if params["page"] != "2" {
			t.Errorf("SetPageCRM: page = %q", params["page"])
		}
		if params["per_page"] != "200" {
			t.Errorf("SetPageCRM: per_page = %q", params["per_page"])
		}
	})

	t.Run("SetPageCRM preserves existing per_page", func(t *testing.T) {
		params := map[string]string{"per_page": "50"}
		state := &PageState{PageNum: 1}
		SetPageCRM(state, params)
		if params["per_page"] != "50" {
			t.Errorf("SetPageCRM: per_page = %q, want 50", params["per_page"])
		}
	})

	t.Run("FromLimit", func(t *testing.T) {
		params := map[string]string{}
		state := &PageState{Offset: 100}
		FromLimit(50)(state, params)
		if params["from"] != "100" || params["limit"] != "50" {
			t.Errorf("FromLimit: from=%s limit=%s", params["from"], params["limit"])
		}
	})

	t.Run("PageOffsetLimit", func(t *testing.T) {
		params := map[string]string{}
		state := &PageState{Offset: 50}
		PageOffsetLimit(50)(state, params)
		if params["page[offset]"] != "50" || params["page[limit]"] != "50" {
			t.Errorf("PageOffsetLimit: offset=%s limit=%s", params["page[offset]"], params["page[limit]"])
		}
	})

	t.Run("IndexRange", func(t *testing.T) {
		params := map[string]string{}
		state := &PageState{Offset: 0}
		IndexRange(100)(state, params)
		if params["index"] != "1" || params["range"] != "100" {
			t.Errorf("IndexRange: index=%s range=%s", params["index"], params["range"])
		}
		state.Offset = 100
		IndexRange(100)(state, params)
		if params["index"] != "101" {
			t.Errorf("IndexRange second page: index=%s", params["index"])
		}
	})

	t.Run("SIndexLimit", func(t *testing.T) {
		params := map[string]string{}
		state := &PageState{Offset: 0}
		SIndexLimit(200)(state, params)
		if params["sIndex"] != "1" || params["limit"] != "200" {
			t.Errorf("SIndexLimit: sIndex=%s limit=%s", params["sIndex"], params["limit"])
		}
	})

	t.Run("SignPageContext", func(t *testing.T) {
		params := map[string]string{}
		state := &PageState{Offset: 50}
		SignPageContext(25)(state, params)
		var parsed struct {
			PageContext struct {
				StartIndex int `json:"start_index"`
				RowCount   int `json:"row_count"`
			} `json:"page_context"`
		}
		if err := json.Unmarshal([]byte(params["data"]), &parsed); err != nil {
			t.Fatalf("SignPageContext: invalid JSON: %v", err)
		}
		if parsed.PageContext.StartIndex != 50 || parsed.PageContext.RowCount != 25 {
			t.Errorf("SignPageContext: start_index=%d row_count=%d", parsed.PageContext.StartIndex, parsed.PageContext.RowCount)
		}
	})
}

func TestPaginateMultiPage(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		page := r.URL.Query().Get("page")
		switch page {
		case "1":
			fmt.Fprint(w, `{"items":[{"id":1},{"id":2}],"page_context":{"has_more_page":true,"page":1}}`)
		case "2":
			fmt.Fprint(w, `{"items":[{"id":3},{"id":4}],"page_context":{"has_more_page":true,"page":2}}`)
		case "3":
			fmt.Fprint(w, `{"items":[{"id":5}],"page_context":{"has_more_page":false,"page":3}}`)
		default:
			fmt.Fprint(w, `{"items":[],"page_context":{"has_more_page":false}}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "items",
		PageSize: 2,
		SetPage:  PagePerPage(2),
		HasMore:  HasMoreBooks,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 5 {
		t.Errorf("Paginate got %d items, want 5", len(items))
	}
	if callCount != 3 {
		t.Errorf("Paginate made %d calls, want 3", callCount)
	}
}

func TestPaginateWithLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		switch page {
		case "1":
			fmt.Fprint(w, `{"items":[{"id":1},{"id":2},{"id":3}],"page_context":{"has_more_page":true,"page":1}}`)
		case "2":
			fmt.Fprint(w, `{"items":[{"id":4},{"id":5},{"id":6}],"page_context":{"has_more_page":true,"page":2}}`)
		default:
			fmt.Fprint(w, `{"items":[],"page_context":{"has_more_page":false}}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "items",
		PageSize: 3,
		Limit:    4,
		SetPage:  PagePerPage(3),
		HasMore:  HasMoreBooks,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 4 {
		t.Errorf("Paginate got %d items, want 4", len(items))
	}
}

func TestPaginateLimitExactPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"items":[{"id":1},{"id":2},{"id":3}],"page_context":{"has_more_page":true,"page":1}}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "items",
		PageSize: 3,
		Limit:    3,
		SetPage:  PagePerPage(3),
		HasMore:  HasMoreBooks,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("Paginate got %d items, want 3", len(items))
	}
}

func TestPaginateEmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"items":[],"page_context":{"has_more_page":false}}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "items",
		PageSize: 10,
		SetPage:  PagePerPage(10),
		HasMore:  HasMoreBooks,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("Paginate got %d items, want 0", len(items))
	}
}

func TestPaginatePassesParams(t *testing.T) {
	var gotStatus string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotStatus = r.URL.Query().Get("status")
		fmt.Fprint(w, `{"items":[{"id":1}],"page_context":{"has_more_page":false}}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	_, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		Opts:     &zohttp.RequestOpts{Params: map[string]string{"status": "active"}},
		ItemsKey: "items",
		PageSize: 10,
		SetPage:  PagePerPage(10),
		HasMore:  HasMoreBooks,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if gotStatus != "active" {
		t.Errorf("status param = %q, want active", gotStatus)
	}
}

func TestPaginatePassesHeaders(t *testing.T) {
	var gotOrgID string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotOrgID = r.Header.Get("orgId")
		fmt.Fprint(w, `{"data":[{"id":1}]}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	_, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		Opts:     &zohttp.RequestOpts{Headers: map[string]string{"orgId": "12345"}},
		ItemsKey: "data",
		PageSize: 100,
		SetPage:  FromLimit(100),
		HasMore:  HasMoreByCount,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if gotOrgID != "12345" {
		t.Errorf("orgId header = %q, want 12345", gotOrgID)
	}
}

func TestPaginateCRMWithToken(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		token := r.URL.Query().Get("page_token")
		if callCount == 1 {
			fmt.Fprint(w, `{"data":[{"id":1}],"info":{"more_records":true,"next_page_token":"tok1"}}`)
		} else if token == "tok1" {
			fmt.Fprint(w, `{"data":[{"id":2}],"info":{"more_records":false}}`)
		} else {
			t.Errorf("unexpected call %d with token %q", callCount, token)
			fmt.Fprint(w, `{"data":[],"info":{"more_records":false}}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "data",
		PageSize: 200,
		SetPage:  SetPageCRM,
		HasMore:  HasMoreCRM,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("got %d items, want 2", len(items))
	}
	if callCount != 2 {
		t.Errorf("made %d calls, want 2", callCount)
	}
}

func TestPaginateOffsetBased(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		switch from {
		case "0":
			fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
		case "2":
			fmt.Fprint(w, `{"data":[{"id":3}]}`)
		default:
			fmt.Fprint(w, `{"data":[]}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "data",
		PageSize: 2,
		SetPage:  FromLimit(2),
		HasMore:  HasMoreByCount,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("got %d items, want 3", len(items))
	}
}

func TestPaginateWorkDrivePattern(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset := r.URL.Query().Get("page[offset]")
		switch offset {
		case "0":
			fmt.Fprint(w, `{"data":[{"id":1},{"id":2}],"meta":{"has_next":true}}`)
		case "2":
			fmt.Fprint(w, `{"data":[{"id":3}],"meta":{"has_next":false}}`)
		default:
			fmt.Fprint(w, `{"data":[],"meta":{"has_next":false}}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "data",
		PageSize: 2,
		SetPage:  PageOffsetLimit(2),
		HasMore:  HasMoreWorkDrive,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("got %d items, want 3", len(items))
	}
}

func TestPaginateSignPattern(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dataParam := r.URL.Query().Get("data")
		var parsed struct {
			PageContext struct {
				StartIndex int `json:"start_index"`
				RowCount   int `json:"row_count"`
			} `json:"page_context"`
		}
		json.Unmarshal([]byte(dataParam), &parsed)
		if parsed.PageContext.StartIndex == 0 {
			fmt.Fprint(w, `{"requests":[{"id":1},{"id":2}],"page_context":{"has_more_rows":true}}`)
		} else {
			fmt.Fprint(w, `{"requests":[{"id":3}],"page_context":{"has_more_rows":false}}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "requests",
		PageSize: 2,
		SetPage:  SignPageContext(2),
		HasMore:  HasMoreSign,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("got %d items, want 3", len(items))
	}
}

func TestPaginateProjectsRawArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "",
		PageSize: 100,
		SetPage:  PagePerPage(100),
		HasMore:  HasMoreProjects,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("got %d items, want 2", len(items))
	}
}

func TestPaginatePeopleNestedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sindex := r.URL.Query().Get("sIndex")
		switch sindex {
		case "1":
			fmt.Fprint(w, `{"response":{"result":[{"id":1},{"id":2}],"count":5}}`)
		case "3":
			fmt.Fprint(w, `{"response":{"result":[{"id":3},{"id":4}],"count":5}}`)
		default:
			fmt.Fprint(w, `{"response":{"result":[{"id":5}],"count":5}}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "response.result",
		PageSize: 2,
		SetPage:  SIndexLimit(2),
		HasMore:  HasMoreByCount,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 5 {
		t.Errorf("got %d items, want 5", len(items))
	}
}

func TestPaginateDefaultMethod(t *testing.T) {
	var gotMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		fmt.Fprint(w, `{"data":[{"id":1}]}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "data",
		PageSize: 10,
		SetPage:  FromLimit(10),
		HasMore:  HasMoreByCount,
	})
	if gotMethod != "GET" {
		t.Errorf("default method = %q, want GET", gotMethod)
	}
}

func TestPaginateCustomMethod(t *testing.T) {
	var gotMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		fmt.Fprint(w, `{"data":[{"id":1}]}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	Paginate(PaginationConfig{
		Client:   c,
		Method:   "POST",
		URL:      server.URL,
		ItemsKey: "data",
		PageSize: 10,
		SetPage:  FromLimit(10),
		HasMore:  HasMoreByCount,
	})
	if gotMethod != "POST" {
		t.Errorf("method = %q, want POST", gotMethod)
	}
}

func TestPaginateDefaultPageSize(t *testing.T) {
	var gotPerPage string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPerPage = r.URL.Query().Get("per_page")
		fmt.Fprint(w, `{"data":[{"id":1}]}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "data",
		SetPage:  PagePerPage(0),
		HasMore:  HasMoreByCount,
	})
	if gotPerPage != "0" {
	}
}

func TestPaginateHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error":"internal"}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "data",
		PageSize: 10,
		SetPage:  FromLimit(10),
		HasMore:  HasMoreByCount,
	})
	if err == nil {
		t.Error("expected error on HTTP 500")
	}
	if items != nil {
		t.Errorf("expected nil items on error, got %d", len(items))
	}
}

func TestPaginateErrorMidPagination(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
		} else {
			w.WriteHeader(500)
			fmt.Fprint(w, `{"error":"boom"}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "data",
		PageSize: 2,
		SetPage:  FromLimit(2),
		HasMore:  HasMoreByCount,
	})
	if err == nil {
		t.Error("expected error on second page")
	}
	if items != nil {
		t.Errorf("expected nil items on mid-pagination error, got %d", len(items))
	}
}

func TestPaginateCRMWrapper(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if r.URL.Query().Get("per_page") == "" {
			t.Error("PaginateCRM should set per_page")
		}
		if callCount == 1 {
			fmt.Fprint(w, `{"data":[{"id":1}],"info":{"more_records":true,"next_page_token":"t1"}}`)
		} else {
			fmt.Fprint(w, `{"data":[{"id":2}],"info":{"more_records":false}}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := PaginateCRM(c, server.URL, nil, 0)
	if err != nil {
		t.Fatalf("PaginateCRM: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("PaginateCRM got %d items, want 2", len(items))
	}
}

func TestPaginateProjectsWrapper(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("per_page") != "100" {
			t.Errorf("PaginateProjects per_page = %q, want 100", r.URL.Query().Get("per_page"))
		}
		fmt.Fprint(w, `{"tasks":[{"id":1},{"id":2}],"page_info":{"has_next_page":false}}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := PaginateProjects(c, server.URL, "tasks", nil, 0)
	if err != nil {
		t.Fatalf("PaginateProjects: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("PaginateProjects got %d items, want 2", len(items))
	}
}

func TestPaginateProjectsWrapperRawArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1}]`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := PaginateProjects(c, server.URL, "", nil, 0)
	if err != nil {
		t.Fatalf("PaginateProjects: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("PaginateProjects got %d items, want 1", len(items))
	}
}

func TestPaginateWorkDriveWrapper(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page[limit]") != "50" {
			t.Errorf("PaginateWorkDrive page[limit] = %q, want 50", r.URL.Query().Get("page[limit]"))
		}
		fmt.Fprint(w, `{"data":[{"id":1}],"meta":{"has_next":false}}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := PaginateWorkDrive(c, server.URL, nil, 0)
	if err != nil {
		t.Fatalf("PaginateWorkDrive: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("PaginateWorkDrive got %d items, want 1", len(items))
	}
}

func TestPaginateWorkDriveWrapperMultiPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset := r.URL.Query().Get("page[offset]")
		switch offset {
		case "0":
			fmt.Fprint(w, `{"data":[{"id":1},{"id":2}],"meta":{"has_next":true}}`)
		case "2":
			fmt.Fprint(w, `{"data":[{"id":3}],"meta":{"has_next":false}}`)
		default:
			fmt.Fprint(w, `{"data":[],"meta":{"has_next":false}}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := PaginateWorkDrive(c, server.URL, nil, 0)
	if err != nil {
		t.Fatalf("PaginateWorkDrive: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("PaginateWorkDrive got %d items, want 3", len(items))
	}
}

func TestPaginateCRMWrapperWithMaxPages(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		fmt.Fprint(w, `{"data":[{"id":1}],"info":{"more_records":true,"next_page_token":"tok"}}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := PaginateCRM(c, server.URL, nil, 2)
	if err != nil {
		t.Fatalf("PaginateCRM: %v", err)
	}
	if len(items) > 400 {
		t.Errorf("PaginateCRM with maxPages=2 got %d items, expected <= 400", len(items))
	}
}

func TestPaginateIndexRangePattern(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idx := r.URL.Query().Get("index")
		rng := r.URL.Query().Get("range")
		if rng != "2" {
			t.Errorf("range = %q, want 2", rng)
		}
		switch idx {
		case "1":
			fmt.Fprint(w, `{"items":[{"id":1},{"id":2}]}`)
		case "3":
			fmt.Fprint(w, `{"items":[{"id":3}]}`)
		default:
			fmt.Fprint(w, `{"items":[]}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "items",
		PageSize: 2,
		SetPage:  IndexRange(2),
		HasMore:  HasMoreByCount,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("got %d items, want 3", len(items))
	}
}

func TestPaginateLimitZeroMeansNoLimit(t *testing.T) {
	pageCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageCount++
		if pageCount <= 3 {
			fmt.Fprint(w, `{"items":[{"id":1}],"page_context":{"has_more_page":true,"page":`+strconv.Itoa(pageCount)+`}}`)
		} else {
			fmt.Fprint(w, `{"items":[{"id":1}],"page_context":{"has_more_page":false,"page":`+strconv.Itoa(pageCount)+`}}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		ItemsKey: "items",
		PageSize: 1,
		Limit:    0,
		SetPage:  PagePerPage(1),
		HasMore:  HasMoreBooks,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 4 {
		t.Errorf("got %d items, want 4", len(items))
	}
}

func TestPaginateNilOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data":[{"id":1}]}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := Paginate(PaginationConfig{
		Client:   c,
		URL:      server.URL,
		Opts:     nil,
		ItemsKey: "data",
		PageSize: 10,
		SetPage:  FromLimit(10),
		HasMore:  HasMoreByCount,
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("got %d items, want 1", len(items))
	}
}

func TestPaginateProjectsWrapperMultiPage(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		page := r.URL.Query().Get("page")
		switch page {
		case "1":
			fmt.Fprint(w, `{"tasks":[{"id":1}],"page_info":{"has_next_page":true}}`)
		case "2":
			fmt.Fprint(w, `{"tasks":[{"id":2}],"page_info":{"has_next_page":false}}`)
		default:
			fmt.Fprint(w, `{"tasks":[],"page_info":{"has_next_page":false}}`)
		}
	}))
	defer server.Close()

	c := testClient(server.URL)
	items, err := PaginateProjects(c, server.URL, "tasks", nil, 0)
	if err != nil {
		t.Fatalf("PaginateProjects: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("PaginateProjects got %d items, want 2", len(items))
	}
	if callCount != 2 {
		t.Errorf("PaginateProjects made %d calls, want 2", callCount)
	}
}

func TestPaginateProjectsWrapperWithParams(t *testing.T) {
	var gotStatus string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotStatus = r.URL.Query().Get("status")
		fmt.Fprint(w, `{"tasks":[{"id":1}],"page_info":{"has_next_page":false}}`)
	}))
	defer server.Close()

	c := testClient(server.URL)
	params := map[string]string{"status": "open"}
	items, err := PaginateProjects(c, server.URL, "tasks", params, 0)
	if err != nil {
		t.Fatalf("PaginateProjects: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("PaginateProjects got %d items, want 1", len(items))
	}
	if gotStatus != "open" {
		t.Errorf("status = %q, want open", gotStatus)
	}
}
