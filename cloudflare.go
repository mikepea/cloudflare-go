/*
Package cloudflare implements the CloudFlare v4 API.

New API requests created like:

    api := cloudflare.New(apikey, apiemail)

*/
package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const apiURL = "https://api.cloudflare.com/client/v4"

type API struct {
	APIKey   string
	APIEmail string
}

// Initializes the API configuration.
func New(key, email string) *API {
	return &API{key, email}
}

// Initializes a new zone.
func NewZone() *Zone {
	return &Zone{}
}

// Params can be turned into a URL query string or a body
// TODO: Give this func a better name
func (api *API) makeRequest(method, uri string, params interface{}) ([]byte, error) {
	// Replace nil with a JSON object if needed
	var reqBody io.Reader
	switch method {
	case "GET", "DELETE":
		reqBody = nil
	default:
		json, err := json.Marshal(params)
		if err != nil {
			fmt.Println("Error marshalling params to JSON:", err)
			return nil, err
		}
		reqBody = bytes.NewReader(json)
	}
	req, err := http.NewRequest(method, apiURL+uri, reqBody)
	if err != nil {
		fmt.Println("Herp derp making request:", err)
		return nil, err
	}
	req.Header.Add("X-Auth-Key", api.APIKey)
	req.Header.Add("X-Auth-Email", api.APIEmail)
	// Could be application/json or multipart/form-data
	// req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	// TODO(jamesog): Check resp.StatusCode != http.StatusOK
	// if resp.StatusCode != http.StatusOK
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	return resBody, nil
}

// The Response struct is a template.  There will also be a result struct.
// There will be a unique response type for each response, which will include
// this type.
type Response struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
}

type ResultInfo struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Count   int `json:"count"`
	Total   int `json:"total_count"`
}

// An Organization describes a multi-user organization. (Enterprise only.)
type Organization struct {
	ID          string
	Name        string
	Status      string
	Permissions []string
	Roles       []string
}

// A User describes a user account.
type User struct {
	ID            string         `json:"id"`
	Email         string         `json:"email"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	Username      string         `json:"username"`
	Telephone     string         `json:"telephone"`
	Country       string         `json:"country"`
	Zipcode       string         `json:"zipcode"`
	CreatedOn     string         `json:"created_on"` // Should this be a time.Date?
	ModifiedOn    string         `json:"modified_on"`
	APIKey        string         `json:"api_key"`
	TwoFA         bool           `json:"two_factor_authentication_enabled"`
	Betas         []string       `json:"betas"`
	Organizations []Organization `json:"organizations"`
}

type UserResponse struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	Result   User     `json:"result"`
}

type Owner struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	OwnerType string `json:"owner_type"`
}

// A Zone describes a CloudFlare zone.
type Zone struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	DevMode           int      `json:"development_mode"`
	OriginalNS        []string `json:"original_name_servers"`
	OriginalRegistrar string   `json:"original_registrar"`
	OriginalDNSHost   string   `json:"original_dnshost"`
	CreatedOn         string   `json:"created_on"`
	ModifiedOn        string   `json:"modified_on"`
	NameServers       []string `json:"name_servers"`
	Owner             Owner    `json:"owner"`
	Permissions       []string `json:"permissions"`
	Plan              ZonePlan `json:"plan"`
	Status            string   `json:"success"`
	Paused            bool     `json:"paused"`
	Type              string   `json:"type"`
	Host              struct {
		Name    string
		Website string
	} `json:"host"`
	VanityNS    []string `json:"vanity_name_servers"`
	Betas       []string `json:"betas"`
	DeactReason string   `json:"deactivation_reason"`
	Meta        ZoneMeta `json:"meta"`
}

// Contains metadata about a zone.
type ZoneMeta struct {
	// custom_certificate_quota is broken - sometimes it's a string, sometimes a number!
	// CustCertQuota     int    `json:"custom_certificate_quota"`
	PageRuleQuota     string `json:"page_rule_quota"`
	WildcardProxiable bool   `json:"wildcard_proxiable"`
	PhishingDetected  bool   `json:"phishing_detected"`
}

// Contains the plan information for a zone.
type ZonePlan struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Currency     string `json:"currency"`
	Frequency    string `json:"frequency"`
	LegacyID     string `json:"legacy_id"`
	IsSubscribed bool   `json:"is_subscribed"`
	CanSubscribe bool   `json:"can_subscribe"`
}

type ZoneResponse struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	Result   []Zone   `json:"result"`
}

type ZonePlanResponse struct {
	Success  bool       `json:"success"`
	Errors   []string   `json:"errors"`
	Messages []string   `json:"messages"`
	Result   []ZonePlan `json:"result"`
}

// type zoneSetting struct {
// 	ID         string `json:"id"`
// 	Editable   bool   `json:"editable"`
// 	ModifiedOn string `json:"modified_on"`
// }
// type zoneSettingStringVal struct {
// 	zoneSetting
// 	Value string `json:"value"`
// }
// type zoneSettingIntVal struct {
// 	zoneSetting
// 	Value int64 `json:"value"`
// }

type ZoneSetting struct {
	ID            string      `json:"id"`
	Editable      bool        `json:"editable"`
	ModifiedOn    string      `json:"modified_on"`
	Value         interface{} `json:"value"`
	TimeRemaining int         `json:"time_remaining"`
}

type ZoneSettingResponse struct {
	Success  bool          `json:"success"`
	Errors   []string      `json:"errors"`
	Messages []string      `json:"messages"`
	Result   []ZoneSetting `json:"result"`
}

// Describes a DNS record for a zone.
type DNSRecord struct {
	ID         string      `json:"id,omitempty"`
	Type       string      `json:"type,omitempty"`
	Name       string      `json:"name,omitempty"`
	Content    string      `json:"content,omitempty"`
	Proxiable  bool        `json:"proxiable,omitempty"`
	Proxied    bool        `json:"proxied,omitempty"`
	TTL        int         `json:"ttl,omitempty"`
	Locked     bool        `json:"locked,omitempty"`
	ZoneID     string      `json:"zone_id,omitempty"`
	ZoneName   string      `json:"zone_name,omitempty"`
	CreatedOn  string      `json:"created_on,omitempty"`
	ModifiedOn string      `json:"modified_on,omitempty"`
	Data       interface{} `json:"data,omitempty"` // data returned by: SRV, LOC
	Meta       interface{} `json:"meta,omitempty"`
	Priority   int         `json:"priority,omitempty"`
}

// The response for creating or updating a DNS record.
type DNSRecordResponse struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []string      `json:"messages"`
	Result   DNSRecord     `json:"result"`
}

// The response for listing DNS records.
type DNSListResponse struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []string      `json:"messages"`
	Result   []DNSRecord   `json:"result"`
}

// Railgun status for a zone.
type ZoneRailgun struct {
	ID        string `json:"id"`
	Name      string `json:"string"`
	Enabled   bool   `json:"enabled"`
	Connected bool   `json:"connected"`
}

type ZoneRailgunResponse struct {
	Success  bool          `json:"success"`
	Errors   []string      `json:"errors"`
	Messages []string      `json:"messages"`
	Result   []ZoneRailgun `json:"result"`
}

// Custom SSL certificates for a zone.
type ZoneCustomSSL struct {
	ID            string     `json:"id"`
	Hosts         []string   `json:"hosts"`
	Issuer        string     `json:"issuer"`
	Priority      int        `json:"priority"`
	Status        string     `json:"success"`
	BundleMethod  string     `json:"bundle_method"`
	ZoneID        string     `json:"zone_id"`
	Permissions   []string   `json:"permissions"`
	UploadedOn    string     `json:"uploaded_on"`
	ModifiedOn    string     `json:"modified_on"`
	ExpiresOn     string     `json:"expires_on"`
	KeylessServer KeylessSSL `json:"keyless_server"`
}

type ZoneCustomSSLResponse struct {
	Success  bool            `json:"success"`
	Errors   []string        `json:"errors"`
	Messages []string        `json:"messages"`
	Result   []ZoneCustomSSL `json:"result"`
}

type KeylessSSL struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	Status      string   `json:"success"`
	Enabled     bool     `json:"enabled"`
	Permissions []string `json:"permissions"`
	CreatedOn   string   `json:"created_on"`
	ModifiedOn  string   `json:"modifed_on"`
}

type KeylessSSLResponse struct {
	Success  bool         `json:"success"`
	Errors   []string     `json:"errors"`
	Messages []string     `json:"messages"`
	Result   []KeylessSSL `json:"result"`
}

type Railgun struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Status         string `json:"success"`
	Enabled        bool   `json:"enabled"`
	ZonesConnected int    `json:"zones_connected"`
	Build          string `json:"build"`
	Version        string `json:"version"`
	Revision       string `json:"revision"`
	ActivationKey  string `json:"activation_key"`
	ActivatedOn    string `json:"activated_on"`
	CreatedOn      string `json:"created_on"`
	ModifiedOn     string `json:"modified_on"`
	// XXX: UpgradeInfo struct {
	// version string
	// url string
	// } `json:"upgrade_info"`
}

type RailgunResponse struct {
	Success  bool      `json:"success"`
	Errors   []string  `json:"errors"`
	Messages []string  `json:"messages"`
	Result   []Railgun `json:"result"`
}

// Custom error pages.
type CustomPage struct {
	CreatedOn      string   `json:"created_on"`
	ModifiedOn     string   `json:"modified_on"`
	URL            string   `json:"url"`
	State          string   `json:"state"`
	RequiredTokens []string `json:"required_tokens"`
	PreviewTarget  string   `json:"preview_target"`
	Description    string   `json:"description"`
}

type CustomPageResponse struct {
	Success  bool         `json:"success"`
	Errors   []string     `json:"errors"`
	Messages []string     `json:"messages"`
	Result   []CustomPage `json:"result"`
}
