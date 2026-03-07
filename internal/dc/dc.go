package dc

type DCConfig struct {
	Accounts  string
	API       string
	Cliq      string
	CRM       string
	Desk      string
	Mail      string
	People    string
	Projects  string
	Recruit   string
	Sheet     string
	Sign      string
	Sprints   string
	WorkDrive string
	Writer    string
	Download  string
}

var dcMap = map[string]DCConfig{
	"com": {
		Accounts:  "https://accounts.zoho.com",
		API:       "https://www.zohoapis.com",
		Cliq:      "https://cliq.zoho.com",
		CRM:       "https://zohoapis.com",
		Desk:      "https://desk.zoho.com",
		Mail:      "https://mail.zoho.com",
		People:    "https://people.zoho.com",
		Projects:  "https://projectsapi.zoho.com",
		Recruit:   "https://recruit.zoho.com",
		Sheet:     "https://sheet.zoho.com",
		Sign:      "https://sign.zoho.com",
		Sprints:   "https://sprintsapi.zoho.com",
		WorkDrive: "https://workdrive.zoho.com",
		Writer:    "https://www.zohoapis.com/writer",
		Download:  "https://download.zoho.com",
	},
	"eu": {
		Accounts:  "https://accounts.zoho.eu",
		API:       "https://www.zohoapis.eu",
		Cliq:      "https://cliq.zoho.eu",
		CRM:       "https://zohoapis.eu",
		Desk:      "https://desk.zoho.eu",
		Mail:      "https://mail.zoho.eu",
		People:    "https://people.zoho.eu",
		Projects:  "https://projectsapi.zoho.eu",
		Recruit:   "https://recruit.zoho.eu",
		Sheet:     "https://sheet.zoho.eu",
		Sign:      "https://sign.zoho.eu",
		Sprints:   "https://sprintsapi.zoho.eu",
		WorkDrive: "https://workdrive.zoho.eu",
		Writer:    "https://www.zohoapis.eu/writer",
		Download:  "https://download.zoho.eu",
	},
	"in": {
		Accounts:  "https://accounts.zoho.in",
		API:       "https://www.zohoapis.in",
		Cliq:      "https://cliq.zoho.in",
		CRM:       "https://zohoapis.in",
		Desk:      "https://desk.zoho.in",
		Mail:      "https://mail.zoho.in",
		People:    "https://people.zoho.in",
		Projects:  "https://projectsapi.zoho.in",
		Recruit:   "https://recruit.zoho.in",
		Sheet:     "https://sheet.zoho.in",
		Sign:      "https://sign.zoho.in",
		Sprints:   "https://sprintsapi.zoho.in",
		WorkDrive: "https://workdrive.zoho.in",
		Writer:    "https://www.zohoapis.in/writer",
		Download:  "https://download.zoho.in",
	},
	"com.au": {
		Accounts:  "https://accounts.zoho.com.au",
		API:       "https://www.zohoapis.com.au",
		Cliq:      "https://cliq.zoho.com.au",
		CRM:       "https://zohoapis.com.au",
		Desk:      "https://desk.zoho.com.au",
		Mail:      "https://mail.zoho.com.au",
		People:    "https://people.zoho.com.au",
		Projects:  "https://projectsapi.zoho.com.au",
		Recruit:   "https://recruit.zoho.com.au",
		Sheet:     "https://sheet.zoho.com.au",
		Sign:      "https://sign.zoho.com.au",
		Sprints:   "https://sprintsapi.zoho.com.au",
		WorkDrive: "https://workdrive.zoho.com.au",
		Writer:    "https://www.zohoapis.com.au/writer",
		Download:  "https://download.zoho.com.au",
	},
	"jp": {
		Accounts:  "https://accounts.zoho.jp",
		API:       "https://www.zohoapis.jp",
		Cliq:      "https://cliq.zoho.jp",
		CRM:       "https://zohoapis.jp",
		Desk:      "https://desk.zoho.jp",
		Mail:      "https://mail.zoho.jp",
		People:    "https://people.zoho.jp",
		Projects:  "https://projectsapi.zoho.jp",
		Recruit:   "https://recruit.zoho.jp",
		Sheet:     "https://sheet.zoho.jp",
		Sign:      "https://sign.zoho.jp",
		Sprints:   "https://sprintsapi.zoho.jp",
		WorkDrive: "https://workdrive.zoho.jp",
		Writer:    "https://www.zohoapis.jp/writer",
		Download:  "https://download.zoho.jp",
	},
	"ca": {
		Accounts:  "https://accounts.zohocloud.ca",
		API:       "https://www.zohoapis.ca",
		Cliq:      "https://cliq.zohocloud.ca",
		CRM:       "https://zohoapis.ca",
		Desk:      "https://desk.zohocloud.ca",
		Mail:      "https://mail.zohocloud.ca",
		People:    "https://people.zohocloud.ca",
		Projects:  "https://projectsapi.zohocloud.ca",
		Recruit:   "https://recruit.zohocloud.ca",
		Sheet:     "https://sheet.zohocloud.ca",
		Sign:      "https://sign.zohocloud.ca",
		Sprints:   "https://sprintsapi.zohocloud.ca",
		WorkDrive: "https://workdrive.zohocloud.ca",
		Writer:    "https://www.zohoapis.ca/writer",
		Download:  "https://download.zohocloud.ca",
	},
	"sa": {
		Accounts:  "https://accounts.zoho.sa",
		API:       "https://www.zohoapis.sa",
		Cliq:      "https://cliq.zoho.sa",
		CRM:       "https://zohoapis.sa",
		Desk:      "https://desk.zoho.sa",
		Mail:      "https://mail.zoho.sa",
		People:    "https://people.zoho.sa",
		Projects:  "https://projectsapi.zoho.sa",
		Recruit:   "https://recruit.zoho.sa",
		Sheet:     "https://sheet.zoho.sa",
		Sign:      "https://sign.zoho.sa",
		Sprints:   "https://sprintsapi.zoho.sa",
		WorkDrive: "https://workdrive.zoho.sa",
		Writer:    "https://www.zohoapis.sa/writer",
		Download:  "https://download.zoho.sa",
	},
	"uk": {
		Accounts:  "https://accounts.zoho.uk",
		API:       "https://www.zohoapis.uk",
		Cliq:      "https://cliq.zoho.uk",
		CRM:       "https://zohoapis.uk",
		Desk:      "https://desk.zoho.uk",
		Mail:      "https://mail.zoho.uk",
		People:    "https://people.zoho.uk",
		Projects:  "https://projectsapi.zoho.uk",
		Recruit:   "https://recruit.zoho.uk",
		Sheet:     "https://sheet.zoho.uk",
		Sign:      "https://sign.zoho.uk",
		Sprints:   "https://sprintsapi.zoho.uk",
		WorkDrive: "https://workdrive.zoho.uk",
		Writer:    "https://www.zohoapis.uk/writer",
		Download:  "https://download.zoho.uk",
	},
	"com.cn": {
		Accounts:  "https://accounts.zoho.com.cn",
		API:       "https://www.zohoapis.com.cn",
		Cliq:      "https://cliq.zoho.com.cn",
		CRM:       "https://zohoapis.com.cn",
		Desk:      "https://desk.zoho.com.cn",
		Mail:      "https://mail.zoho.com.cn",
		People:    "https://people.zoho.com.cn",
		Projects:  "https://projectsapi.zoho.com.cn",
		Recruit:   "https://recruit.zoho.com.cn",
		Sheet:     "https://sheet.zoho.com.cn",
		Sign:      "https://sign.zoho.com.cn",
		Sprints:   "https://sprintsapi.zoho.com.cn",
		WorkDrive: "https://workdrive.zoho.com.cn",
		Writer:    "https://www.zohoapis.com.cn/writer",
		Download:  "https://download.zoho.com.cn",
	},
}

var ValidDCs = []string{"com", "eu", "in", "com.au", "jp", "ca", "sa", "uk", "com.cn"}

func GetDC(dc string) DCConfig {
	if cfg, ok := dcMap[dc]; ok {
		return cfg
	}
	return dcMap["com"]
}

func AccountsURL(dc string) string  { return GetDC(dc).Accounts }
func APIURL(dc string) string       { return GetDC(dc).API }
func CliqURL(dc string) string      { return GetDC(dc).Cliq }
func CRMURL(dc string) string       { return GetDC(dc).CRM }
func DeskURL(dc string) string      { return GetDC(dc).Desk }
func ExpenseURL(dc string) string   { return APIURL(dc) }
func MailURL(dc string) string      { return GetDC(dc).Mail }
func PeopleURL(dc string) string    { return GetDC(dc).People }
func ProjectsURL(dc string) string  { return GetDC(dc).Projects }
func RecruitURL(dc string) string   { return GetDC(dc).Recruit }
func SheetURL(dc string) string     { return GetDC(dc).Sheet }
func SignURL(dc string) string      { return GetDC(dc).Sign }
func SprintsURL(dc string) string   { return GetDC(dc).Sprints }
func WorkDriveURL(dc string) string { return GetDC(dc).WorkDrive }
func WriterURL(dc string) string    { return GetDC(dc).Writer }
func DownloadURL(dc string) string  { return GetDC(dc).Download }
