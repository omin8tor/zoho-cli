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

const DefaultScopes = "" +
	"ZohoAnalytics.fullaccess.all," +
	"ZohoAssist.userapi.READ,ZohoAssist.sessionapi.CREATE,ZohoAssist.unattended.computer.READ,ZohoAssist.unattended.computer.UPDATE,ZohoAssist.unattended.computer.DELETE,ZohoAssist.unattended.group.READ,ZohoAssist.unattended.group.CREATE,ZohoAssist.unattended.group.UPDATE,ZohoAssist.unattended.group.DELETE,ZohoAssist.reportapi.READ," +
	"zohobackstage.portal.READ,zohobackstage.event.CREATE,zohobackstage.event.UPDATE,zohobackstage.event.DELETE,zohobackstage.agenda.CREATE,zohobackstage.speaker.CREATE,zohobackstage.sponsor.CREATE,zohobackstage.eventticket.CREATE,zohobackstage.order.CREATE,zohobackstage.webhook.CREATE," +
	"ZohoSubscriptions.fullaccess.all," +
	"zohobookings.data.CREATE," +
	"ZohoBooks.fullaccess.all," +
	"ZohoCampaigns.campaign.ALL,ZohoCampaigns.contact.ALL," +
	"ZohoCliq.Webhooks.CREATE,ZohoCliq.Channels.ALL,ZohoCliq.Messages.ALL,ZohoCliq.Chats.ALL,ZohoCliq.Users.ALL,ZohoCliq.Bots.ALL," +
	"ZohoCreator.form.CREATE,ZohoCreator.report.CREATE,ZohoCreator.report.READ,ZohoCreator.report.UPDATE,ZohoCreator.report.DELETE,ZohoCreator.meta.form.READ,ZohoCreator.meta.application.READ,ZohoCreator.dashboard.READ," +
	"Desk.tickets.ALL,Desk.contacts.ALL,Desk.basic.READ,Desk.settings.READ,Desk.search.READ,Desk.accounts.ALL," +
	"ZohoExpense.fullaccess.ALL," +
	"ZohoInventory.FullAccess.all," +
	"ZohoInvoice.fullaccess.all," +
	"ZohoLearn.customportal.ALL,ZohoLearn.manual.ALL,ZohoLearn.space.ALL,ZohoLearn.article.ALL,ZohoLearn.attachment.ALL,ZohoLearn.comment.ALL,ZohoLearn.template.ALL,ZohoLearn.favorite.ALL,ZohoLearn.activity.ALL,ZohoLearn.hubMember.ALL,ZohoLearn.questionbank.ALL,ZohoLearn.tag.ALL,ZohoLearn.member.ALL,ZohoLearn.network.ALL,ZohoLearn.commentlike.ALL,ZohoLearn.articleimage.ALL,ZohoLearn.profile.ALL,ZohoLearn.notification.ALL,ZohoLearn.course.ALL,ZohoLearn.lessondiscussion.ALL,ZohoLearn.quiz.ALL," +
	"ZohoMail.accounts.READ,ZohoMail.messages.ALL,ZohoMail.folders.ALL,ZohoMail.tags.ALL,ZohoMail.tasks.ALL,ZohoMail.links.ALL,ZohoMail.notes.ALL,ZohoMail.organization.accounts.ALL,ZohoMail.organization.domains.ALL,ZohoMail.organization.groups.ALL,ZohoMail.organization.policy.ALL,ZohoMail.organization.subscriptions.ALL,ZohoMail.organization.spam.ALL,ZohoMail.organization.audit.READ,ZohoMail.partner.organization.ALL," +
	"ZohoMarketingAutomation.campaign.ALL,ZohoMarketingAutomation.lead.ALL,ZohoMarketingAutomation.journey.READ,ZohoMarketingAutomation.journey.CREATE,ZohoMarketingAutomation.wa.READ," +
	"ZohoMeeting.meeting.ALL," +
	"PageSense.experiments.CREATE,PageSense.experiments.READ,PageSense.experiments.UPDATE,PageSense.experiments.DELETE,PageSense.goals.CREATE,PageSense.goals.READ,PageSense.goals.UPDATE,PageSense.goals.DELETE,PageSense.reports.all,PageSense.customevents.CREATE,PageSense.customevents.READ," +
	"ZOHOPEOPLE.forms.ALL,ZOHOPEOPLE.employee.ALL,ZOHOPEOPLE.dashboard.ALL,ZOHOPEOPLE.automation.ALL,ZOHOPEOPLE.timetracker.ALL,ZOHOPEOPLE.attendance.ALL,ZOHOPEOPLE.leave.ALL,ZOHOPEOPLE.timesheet.READ," +
	"ZohoProjects.portals.ALL,ZohoProjects.projects.ALL,ZohoProjects.tasks.ALL,ZohoProjects.tasklists.ALL,ZohoProjects.timesheets.ALL,ZohoProjects.bugs.ALL,ZohoProjects.events.ALL,ZohoProjects.forums.ALL,ZohoProjects.milestones.ALL,ZohoProjects.documents.ALL,ZohoProjects.users.ALL,ZohoProjects.projectgroups.ALL,ZohoProjects.tags.ALL,ZohoProjects.leave.ALL,ZohoProjects.teams.ALL,ZohoProjects.status.ALL," +
	"ZohoRecruit.modules.ALL,ZohoRecruit.settings.ALL,ZohoRecruit.users.all,ZohoRecruit.org.all,ZohoRecruit.bulk.all,ZohoRecruit.notifications.all," +
	"SalesIQ.operators.ALL,SalesIQ.portals.ALL,SalesIQ.departments.ALL,SalesIQ.leadscorerules.ALL,SalesIQ.leadscoreconfigs.ALL,SalesIQ.criteriafields.READ,SalesIQ.visitorroutingrules.ALL,SalesIQ.chatroutingrules.ALL,SalesIQ.cannedresponses.ALL,SalesIQ.blockedips.ALL,SalesIQ.chatmonitors.ALL,SalesIQ.counts.READ,SalesIQ.visitors.READ,SalesIQ.feedbacks.READ,SalesIQ.conversations.ALL,SalesIQ.trackingpresets.ALL,SalesIQ.userpreferences.ALL,SalesIQ.visitorhistoryviews.ALL,SalesIQ.triggerrules.ALL,SalesIQ.webhooks.ALL,SalesIQ.callbacks.UPDATE,SalesIQ.Apps.ALL,SalesIQ.articles.ALL,SalesIQ.encryptions.CREATE," +
	"ZohoSheet.dataAPI.READ,ZohoSheet.dataAPI.UPDATE," +
	"ZohoShowtime.sessionapi.CREATE,ZohoShowtime.sessionapi.DELETE,ZohoShowtime.talkapi.READ,ZohoShowtime.portalapi.READ," +
	"ZohoSign.documents.ALL," +
	"ZohoSprints.teams.READ,ZohoSprints.teams.CREATE,ZohoSprints.teams.UPDATE,ZohoSprints.teams.DELETE,ZohoSprints.projects.READ,ZohoSprints.projects.CREATE,ZohoSprints.projects.UPDATE,ZohoSprints.projects.DELETE,ZohoSprints.projectsgroups.READ,ZohoSprints.projectsgroups.CREATE,ZohoSprints.epic.READ,ZohoSprints.epic.CREATE,ZohoSprints.epic.UPDATE,ZohoSprints.epic.DELETE,ZohoSprints.sprints.READ,ZohoSprints.sprints.CREATE,ZohoSprints.sprints.UPDATE,ZohoSprints.sprints.DELETE,ZohoSprints.items.READ,ZohoSprints.items.CREATE,ZohoSprints.items.UPDATE,ZohoSprints.items.DELETE,ZohoSprints.projectsettings.READ,ZohoSprints.projectsettings.CREATE,ZohoSprints.projectsettings.UPDATE,ZohoSprints.projectsettings.DELETE,ZohoSprints.webhook.READ,ZohoSprints.webhook.CREATE,ZohoSprints.webhook.UPDATE,ZohoSprints.webhook.DELETE,ZohoSprints.meetings.READ,ZohoSprints.meetings.CREATE,ZohoSprints.meetings.UPDATE,ZohoSprints.meetings.DELETE,ZohoSprints.timesheets.READ,ZohoSprints.timesheets.CREATE,ZohoSprints.timesheets.UPDATE,ZohoSprints.timesheets.DELETE,ZohoSprints.release.READ,ZohoSprints.release.CREATE,ZohoSprints.release.UPDATE,ZohoSprints.release.DELETE,ZohoSprints.teamusers.READ,ZohoSprints.teamusers.CREATE,ZohoSprints.teamusers.UPDATE,ZohoSprints.teamusers.DELETE,ZohoSprints.extensions.READ,ZohoSprints.extensions.UPDATE,ZohoSprints.extensions.WRITE,ZohoSprints.expense.READ,ZohoSprints.expense.CREATE,ZohoSprints.expense.UPDATE,ZohoSprints.expense.DELETE,ZohoSprints.custommodulerecords.READ,ZohoSprints.custommodulerecords.CREATE,ZohoSprints.custommodulerecords.UPDATE,ZohoSprints.custommodulerecords.DELETE," +
	"ZohoVault.user.READ,ZohoVault.user.CREATE,ZohoVault.user.UPDATE,ZohoVault.secrets.CREATE,ZohoVault.secrets.UPDATE,ZohoVault.secrets.READ,ZohoVault.secrets.DELETE,ZohoVault.passwords.READ,ZohoVault.passwords.UPDATE," +
	"ZohoVoice.agents.CREATE,ZohoVoice.agents.UPDATE,ZohoVoice.agents.DELETE,ZohoVoice.agents.READ,ZohoVoice.telephony.CREATE,ZohoVoice.telephony.UPDATE,ZohoVoice.telephony.DELETE,ZohoVoice.telephony.READ,ZohoVoice.call.CREATE,ZohoVoice.call.READ,ZohoVoice.call.DELETE,ZohoVoice.powerdialer.CREATE,ZohoVoice.powerdialer.UPDATE,ZohoVoice.powerdialer.DELETE,ZohoVoice.powerdialer.READ,ZohoVoice.sms.CREATE,ZohoVoice.sms.READ," +
	"ZohoWriter.documentEditor.ALL,ZohoPC.files.ALL," +
	"ZohoSearch.securesearch.READ," +
	"ZohoCRM.modules.ALL,ZohoCRM.settings.ALL,ZohoCRM.users.ALL,ZohoCRM.org.ALL,ZohoCRM.coql.READ,ZohoCRM.bulk.ALL,ZohoCRM.notifications.ALL,ZohoCRM.change_owner.CREATE"

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
