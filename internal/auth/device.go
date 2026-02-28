package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/omin8tor/zoho-cli/internal"
	zohodc "github.com/omin8tor/zoho-cli/internal/dc"
)

const DefaultScopes = "ZohoCliq.Webhooks.CREATE,ZohoCliq.Channels.ALL,ZohoCliq.Messages.ALL," +
	"ZohoCliq.Chats.ALL,ZohoCliq.Users.ALL,ZohoCliq.Bots.ALL," +
	"WorkDrive.workspace.ALL,WorkDrive.files.ALL,WorkDrive.files.sharing.ALL," +
	"WorkDrive.links.ALL,WorkDrive.team.ALL,WorkDrive.teamfolders.ALL," +
	"ZohoSearch.securesearch.READ," +
	"ZohoWriter.documentEditor.ALL,ZohoPC.files.ALL," +
	"ZohoProjects.portals.ALL,ZohoProjects.projects.ALL,ZohoProjects.tasks.ALL," +
	"ZohoProjects.tasklists.ALL,ZohoProjects.timesheets.ALL,ZohoProjects.bugs.ALL," +
	"ZohoProjects.events.ALL,ZohoProjects.forums.ALL,ZohoProjects.milestones.ALL," +
	"ZohoProjects.documents.ALL,ZohoProjects.users.ALL," +
	"ZohoCRM.modules.ALL,ZohoCRM.settings.ALL,ZohoCRM.users.ALL,ZohoCRM.org.ALL," +
	"ZohoCRM.coql.READ,ZohoCRM.bulk.ALL,ZohoCRM.notifications.ALL," +
	"ZohoCRM.change_owner.CREATE," +
	"ZohoExpense.fullaccess.ALL," +
	"ZohoSheet.dataAPI.READ,ZohoSheet.dataAPI.UPDATE"

var locationToDC = map[string]string{
	"us": "com",
	"eu": "eu",
	"in": "in",
	"au": "com.au",
	"jp": "jp",
	"ca": "ca",
	"sa": "sa",
	"uk": "uk",
	"cn": "com.cn",
}

func DeviceFlowLogin(clientID, clientSecret, dc, scopes string) error {
	if scopes == "" {
		scopes = DefaultScopes
	}
	baseURL := zohodc.AccountsURL(dc)
	httpClient := &http.Client{Timeout: 30 * time.Second}

	params := url.Values{
		"client_id":   {clientID},
		"grant_type":  {"device_request"},
		"scope":       {scopes},
		"access_type": {"offline"},
	}
	resp, err := httpClient.PostForm(baseURL+"/oauth/v3/device/code", params)
	if err != nil {
		return internal.NewAuthError(fmt.Sprintf("Device code request failed: %v", err))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return internal.NewAuthError(fmt.Sprintf("Unexpected response: %s", body))
	}
	if errMsg, ok := data["error"]; ok {
		return internal.NewAuthError(fmt.Sprintf("Device code request failed: %v", errMsg))
	}

	userCode, _ := data["user_code"].(string)
	deviceCode, _ := data["device_code"].(string)
	verificationURL, _ := data["verification_url"].(string)
	if verificationURL == "" {
		verificationURL = baseURL + "/oauth/v3/device"
	}
	interval := 5
	if v, ok := data["interval"].(float64); ok && int(v) > 5 {
		interval = int(v)
	}
	expiresIn := 300
	if v, ok := data["expires_in"].(float64); ok {
		expiresIn = int(v)
	}

	fmt.Fprintf(os.Stderr, "\nTo authenticate, visit: %s\n", verificationURL)
	fmt.Fprintf(os.Stderr, "Enter code: %s\n", userCode)
	if vc, ok := data["verification_uri_complete"].(string); ok && vc != "" {
		fmt.Fprintf(os.Stderr, "\nOr open directly: %s\n", vc)
	}
	fmt.Fprintf(os.Stderr, "\nWaiting for authorization (expires in %ds)...\n", expiresIn)

	pollURL := baseURL
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)

	for time.Now().Before(deadline) {
		time.Sleep(time.Duration(interval) * time.Second)

		pollParams := url.Values{
			"client_id":     {clientID},
			"client_secret": {clientSecret},
			"grant_type":    {"device_token"},
			"code":          {deviceCode},
		}
		pollResp, err := httpClient.PostForm(pollURL+"/oauth/v3/device/token", pollParams)
		if err != nil {
			continue
		}
		pollBody, _ := io.ReadAll(pollResp.Body)
		pollResp.Body.Close()

		var pollData map[string]any
		if err := json.Unmarshal(pollBody, &pollData); err != nil {
			continue
		}

		if errStr, ok := pollData["error"].(string); ok {
			switch errStr {
			case "authorization_pending":
				continue
			case "slow_down":
				interval += 2
				if interval > 30 {
					interval = 30
				}
				continue
			case "other_dc":
				if newDC, ok := pollData["dc"].(string); ok && newDC != "" {
					if mapped, ok := locationToDC[newDC]; ok {
						dc = mapped
					} else {
						dc = newDC
					}
					pollURL = zohodc.AccountsURL(dc)
					fmt.Fprintf(os.Stderr, "Redirecting to %s data center...\n", dc)
				}
				continue
			case "access_denied":
				return internal.NewAuthError("Authorization denied by user.")
			case "expired":
				return internal.NewAuthError("Authorization code expired. Please try again.")
			default:
				return internal.NewAuthError(fmt.Sprintf("Device flow error: %s", errStr))
			}
		}

		accessToken, _ := pollData["access_token"].(string)
		refreshToken, _ := pollData["refresh_token"].(string)
		if accessToken == "" || refreshToken == "" {
			return internal.NewAuthError(fmt.Sprintf("Unexpected response: %s", pollBody))
		}

		apiDomain, _ := pollData["api_domain"].(string)
		expires := normalizeExpiresIn(pollData["expires_in"])

		if err := SaveClientConfig(clientID, clientSecret); err != nil {
			return err
		}
		if err := SaveTokens(refreshToken, accessToken, expires, dc, zohodc.AccountsURL(dc), apiDomain, scopes); err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "\nAuthenticated successfully!")
		return nil
	}

	return internal.NewAuthError("Authorization timed out. Please try again.")
}

func SelfClientExchange(clientID, clientSecret, code, dc, accountsServer string) error {
	baseURL := accountsServer
	if baseURL == "" {
		baseURL = zohodc.AccountsURL(dc)
	}

	params := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	}
	resp, err := http.PostForm(baseURL+"/oauth/v2/token", params)
	if err != nil {
		return internal.NewAuthError(fmt.Sprintf("Token exchange failed: %v", err))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return internal.NewAuthError(fmt.Sprintf("Unexpected response: %s", body))
	}
	if errMsg, ok := data["error"]; ok {
		return internal.NewAuthError(fmt.Sprintf("Token exchange failed: %v", errMsg))
	}

	accessToken, _ := data["access_token"].(string)
	refreshToken, _ := data["refresh_token"].(string)
	if accessToken == "" {
		return internal.NewAuthError(fmt.Sprintf("No access token in response: %s", body))
	}
	if refreshToken == "" {
		return internal.NewAuthError("No refresh token returned. Make sure you generated a Self Client code with access_type=offline.")
	}

	apiDomain, _ := data["api_domain"].(string)
	expires := normalizeExpiresIn(data["expires_in"])

	if err := SaveClientConfig(clientID, clientSecret); err != nil {
		return err
	}
	if err := SaveTokens(refreshToken, accessToken, expires, dc, baseURL, apiDomain, ""); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Authenticated successfully!")
	return nil
}

func SaveClientConfig(clientID, clientSecret string) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	cfg := map[string]any{
		"auth": map[string]string{
			"client_id":     clientID,
			"client_secret": clientSecret,
		},
	}
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(ConfigFile(), data, 0600)
}
