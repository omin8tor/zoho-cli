package pagination

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
)

type PageState struct {
	Offset  int
	PageNum int
	Token   string
}

type PaginationConfig struct {
	Client   *zohttp.Client
	Method   string
	URL      string
	Opts     *zohttp.RequestOpts
	ItemsKey string
	PageSize int
	Limit    int
	SetPage  func(state *PageState, params map[string]string)
	HasMore  func(raw json.RawMessage, fetched int, pageSize int) (bool, *PageState)
}

func Paginate(ctx context.Context, cfg PaginationConfig) ([]json.RawMessage, error) {
	if cfg.Method == "" {
		cfg.Method = "GET"
	}
	if cfg.PageSize <= 0 {
		cfg.PageSize = 200
	}

	params := make(map[string]string)
	if cfg.Opts != nil {
		for k, v := range cfg.Opts.Params {
			params[k] = v
		}
	}

	var all []json.RawMessage
	state := &PageState{Offset: 0, PageNum: 1}

	for range 500 {
		cfg.SetPage(state, params)

		opts := &zohttp.RequestOpts{Params: params}
		if cfg.Opts != nil {
			opts.Headers = cfg.Opts.Headers
			opts.JSON = cfg.Opts.JSON
			opts.Form = cfg.Opts.Form
			opts.Files = cfg.Opts.Files
		}

		raw, err := cfg.Client.Request(ctx, cfg.Method, cfg.URL, opts)
		if err != nil {
			return nil, err
		}

		items := ExtractItems(raw, cfg.ItemsKey)
		if len(items) == 0 {
			break
		}
		all = append(all, items...)

		if cfg.Limit > 0 && len(all) >= cfg.Limit {
			all = all[:cfg.Limit]
			break
		}

		hasMore, nextState := cfg.HasMore(raw, len(items), cfg.PageSize)
		if !hasMore {
			break
		}
		if nextState != nil {
			if nextState.PageNum == 0 {
				nextState.PageNum = state.PageNum + 1
			}
			if nextState.Offset == 0 && nextState.Token == "" {
				nextState.Offset = len(all)
			}
			state = nextState
		} else {
			state.Offset = len(all)
			state.PageNum++
			state.Token = ""
		}
	}
	return all, nil
}

func ExtractItems(raw json.RawMessage, key string) []json.RawMessage {
	if key == "" {
		var items []json.RawMessage
		if json.Unmarshal(raw, &items) == nil {
			return items
		}
		return nil
	}

	parts := strings.SplitN(key, ".", 2)
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil
	}
	itemsRaw, ok := envelope[parts[0]]
	if !ok {
		return nil
	}
	if len(parts) == 2 {
		return ExtractItems(itemsRaw, parts[1])
	}
	var items []json.RawMessage
	if err := json.Unmarshal(itemsRaw, &items); err != nil {
		return nil
	}
	return items
}

func hasNextPage(pageInfo map[string]any) bool {
	v, ok := pageInfo["has_next_page"]
	if !ok {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true"
	default:
		return false
	}
}

func PagePerPage(pageSize int) func(*PageState, map[string]string) {
	return func(state *PageState, params map[string]string) {
		params["page"] = strconv.Itoa(state.PageNum)
		params["per_page"] = strconv.Itoa(pageSize)
	}
}

func SetPageCRM(state *PageState, params map[string]string) {
	if state.Token != "" {
		params["page_token"] = state.Token
		delete(params, "page")
	} else {
		params["page"] = strconv.Itoa(state.PageNum)
	}
	if _, ok := params["per_page"]; !ok {
		params["per_page"] = "200"
	}
}

func FromLimit(pageSize int) func(*PageState, map[string]string) {
	return func(state *PageState, params map[string]string) {
		params["from"] = strconv.Itoa(state.Offset)
		params["limit"] = strconv.Itoa(pageSize)
	}
}

func PageOffsetLimit(pageSize int) func(*PageState, map[string]string) {
	return func(state *PageState, params map[string]string) {
		params["page[offset]"] = strconv.Itoa(state.Offset)
		params["page[limit]"] = strconv.Itoa(pageSize)
	}
}

func IndexRange(pageSize int) func(*PageState, map[string]string) {
	return func(state *PageState, params map[string]string) {
		params["index"] = strconv.Itoa(state.Offset + 1)
		params["range"] = strconv.Itoa(pageSize)
	}
}

func SIndexLimit(pageSize int) func(*PageState, map[string]string) {
	return func(state *PageState, params map[string]string) {
		params["sIndex"] = strconv.Itoa(state.Offset + 1)
		params["limit"] = strconv.Itoa(pageSize)
	}
}

func SignPageContext(pageSize int) func(*PageState, map[string]string) {
	return func(state *PageState, params map[string]string) {
		pc := map[string]any{
			"start_index": state.Offset,
			"row_count":   pageSize,
		}
		j, _ := json.Marshal(map[string]any{"page_context": pc})
		params["data"] = string(j)
	}
}

func HasMoreBooks(raw json.RawMessage, _ int, _ int) (bool, *PageState) {
	var env struct {
		PageContext struct {
			HasMorePage bool `json:"has_more_page"`
			Page        int  `json:"page"`
		} `json:"page_context"`
	}
	if json.Unmarshal(raw, &env) != nil {
		return false, nil
	}
	if !env.PageContext.HasMorePage {
		return false, nil
	}
	return true, &PageState{PageNum: env.PageContext.Page + 1}
}

func HasMoreCRM(raw json.RawMessage, _ int, _ int) (bool, *PageState) {
	var env struct {
		Info struct {
			MoreRecords   bool   `json:"more_records"`
			NextPageToken string `json:"next_page_token"`
		} `json:"info"`
	}
	if json.Unmarshal(raw, &env) != nil {
		return false, nil
	}
	if !env.Info.MoreRecords {
		return false, nil
	}
	if env.Info.NextPageToken != "" {
		return true, &PageState{Token: env.Info.NextPageToken}
	}
	return true, nil
}

func HasMoreProjects(raw json.RawMessage, _ int, _ int) (bool, *PageState) {
	var env map[string]json.RawMessage
	if json.Unmarshal(raw, &env) != nil {
		return false, nil
	}
	piRaw, ok := env["page_info"]
	if !ok {
		return false, nil
	}
	var pi map[string]any
	if json.Unmarshal(piRaw, &pi) != nil {
		return false, nil
	}
	return hasNextPage(pi), nil
}

func HasMoreWorkDrive(raw json.RawMessage, _ int, _ int) (bool, *PageState) {
	var env struct {
		Meta struct {
			HasNext bool `json:"has_next"`
		} `json:"meta"`
	}
	if json.Unmarshal(raw, &env) != nil {
		return false, nil
	}
	return env.Meta.HasNext, nil
}

func HasMoreByCount(_ json.RawMessage, fetched int, pageSize int) (bool, *PageState) {
	return fetched >= pageSize, nil
}

func HasMoreSign(raw json.RawMessage, _ int, _ int) (bool, *PageState) {
	var env struct {
		PageContext struct {
			HasMoreRows bool `json:"has_more_rows"`
		} `json:"page_context"`
	}
	if json.Unmarshal(raw, &env) != nil {
		return false, nil
	}
	return env.PageContext.HasMoreRows, nil
}

func PaginateCRM(ctx context.Context, client *zohttp.Client, url string, params map[string]string, maxPages int) ([]json.RawMessage, error) {
	limit := 0
	if maxPages > 0 {
		limit = maxPages * 200
	}
	return Paginate(ctx, PaginationConfig{
		Client:   client,
		URL:      url,
		Opts:     &zohttp.RequestOpts{Params: params},
		ItemsKey: "data",
		PageSize: 200,
		Limit:    limit,
		SetPage:  SetPageCRM,
		HasMore:  HasMoreCRM,
	})
}

func PaginateProjects(ctx context.Context, client *zohttp.Client, url string, itemsKey string, params map[string]string, maxPages int) ([]json.RawMessage, error) {
	limit := 0
	if maxPages > 0 {
		limit = maxPages * 100
	}
	return Paginate(ctx, PaginationConfig{
		Client:   client,
		URL:      url,
		Opts:     &zohttp.RequestOpts{Params: params},
		ItemsKey: itemsKey,
		PageSize: 100,
		Limit:    limit,
		SetPage:  PagePerPage(100),
		HasMore:  HasMoreProjects,
	})
}

func PaginateWorkDrive(ctx context.Context, client *zohttp.Client, url string, params map[string]string, maxPages int) ([]json.RawMessage, error) {
	limit := 0
	if maxPages > 0 {
		limit = maxPages * 50
	}
	return Paginate(ctx, PaginationConfig{
		Client:   client,
		URL:      url,
		Opts:     &zohttp.RequestOpts{Params: params},
		ItemsKey: "data",
		PageSize: 50,
		Limit:    limit,
		SetPage:  PageOffsetLimit(50),
		HasMore:  HasMoreWorkDrive,
	})
}
