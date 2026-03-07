package internal

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

func MergeJSON(cmd *cli.Command, body map[string]any) error {
	j := cmd.String("json")
	if j == "" {
		return nil
	}
	var extra map[string]any
	if err := json.Unmarshal([]byte(j), &extra); err != nil {
		return NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
	}
	for k, v := range extra {
		if _, exists := body[k]; !exists {
			body[k] = v
		}
	}
	return nil
}

func MergeJSONForm(cmd *cli.Command, form map[string]string) error {
	j := cmd.String("json")
	if j == "" {
		return nil
	}
	var extra map[string]any
	if err := json.Unmarshal([]byte(j), &extra); err != nil {
		return NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
	}
	for k, v := range extra {
		if _, exists := form[k]; !exists {
			switch val := v.(type) {
			case string:
				form[k] = val
			default:
				b, err := json.Marshal(val)
				if err != nil {
					form[k] = fmt.Sprintf("%v", val)
				} else {
					form[k] = string(b)
				}
			}
		}
	}
	return nil
}
