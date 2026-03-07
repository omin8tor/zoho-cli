//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func probe(clientID, accountsURL, scope string) (bool, string) {
	params := url.Values{
		"client_id":  {clientID},
		"grant_type": {"device_request"},
		"scope":      {scope},
	}
	resp, err := (&http.Client{Timeout: 15 * time.Second}).PostForm(accountsURL+"/oauth/v3/device/code", params)
	if err != nil {
		return false, err.Error()
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data map[string]any
	json.Unmarshal(body, &data)
	if _, ok := data["device_code"]; ok {
		return true, ""
	}
	if e, ok := data["error"].(string); ok {
		return false, e
	}
	return false, string(body)
}

func main() {
	clientID := os.Getenv("ZOHO_CLIENT_ID")
	if clientID == "" {
		fmt.Fprintln(os.Stderr, "ZOHO_CLIENT_ID required")
		os.Exit(1)
	}
	dc := os.Getenv("ZOHO_DC")
	if dc == "" {
		dc = "com"
	}
	accountsURL := "https://accounts.zoho." + dc

	scopes := os.Args[1:]
	if len(scopes) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: go run scripts/probe_scopes.go <scope1> [scope2] ...")
		fmt.Fprintln(os.Stderr, "  or pipe: echo 'scope1,scope2' | go run scripts/probe_scopes.go -")
		fmt.Fprintln(os.Stderr, "  or all defaults: go run scripts/probe_scopes.go --defaults")
		os.Exit(1)
	}

	if scopes[0] == "-" {
		raw, _ := io.ReadAll(os.Stdin)
		scopes = strings.Split(strings.TrimSpace(string(raw)), ",")
	}

	if scopes[0] == "--defaults" {
		raw := strings.TrimRight(defaultScopes, ",")
		scopes = strings.Split(raw, ",")
	}

	fmt.Fprintf(os.Stderr, "Probing %d scopes against %s...\n\n", len(scopes), accountsURL)

	// First try all at once
	allValid, _ := probe(clientID, accountsURL, strings.Join(scopes, ","))
	if allValid {
		fmt.Fprintf(os.Stderr, "All %d scopes valid! ✓\n", len(scopes))
		for _, s := range scopes {
			fmt.Printf("✓ %s\n", s)
		}
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Batch failed, testing individually (1s delay between)...\n\n")
	time.Sleep(1 * time.Second)

	var invalid []string
	var valid []string
	for i, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		ok, errMsg := probe(clientID, accountsURL, scope)
		if ok {
			fmt.Printf("✓ %s\n", scope)
			valid = append(valid, scope)
		} else {
			fmt.Printf("✗ %s  (%s)\n", scope, errMsg)
			invalid = append(invalid, scope)
		}
		if i < len(scopes)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Fprintf(os.Stderr, "\n--- Summary ---\n")
	fmt.Fprintf(os.Stderr, "Valid:   %d\n", len(valid))
	fmt.Fprintf(os.Stderr, "Invalid: %d\n", len(invalid))
	if len(invalid) > 0 {
		fmt.Fprintf(os.Stderr, "\nInvalid scopes:\n")
		for _, s := range invalid {
			fmt.Fprintf(os.Stderr, "  %s\n", s)
		}
	}
	fmt.Fprintf(os.Stderr, "\nValid scope string:\n")
	fmt.Fprintln(os.Stderr, strings.Join(valid, ","))

	if len(invalid) > 0 {
		os.Exit(1)
	}
}

const defaultScopes = "use --defaults flag to read from device.go directly"
