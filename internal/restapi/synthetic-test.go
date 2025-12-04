package restapi

import "encoding/json"

// ApiTag represents an RBAC tag
type ApiTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MultipleScriptsConfiguration for Jest-based scripts
type MultipleScriptsConfiguration struct {
	Bundle     *string `json:"bundle,omitempty"`
	ScriptFile *string `json:"scriptFile,omitempty"`
}

// DNSFilterQueryTime represents DNS query time filter
type DNSFilterQueryTime struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    int64  `json:"value"`
}

// DNSFilterTargetValue represents DNS target value filter
type DNSFilterTargetValue struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// SSLCertificateValidation represents SSL certificate validation rule
type SSLCertificateValidation struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type SyntheticTestConfig struct {
	MarkSyntheticCall bool    `json:"markSyntheticCall"`
	Retries           int32   `json:"retries,omitempty"`
	RetryInterval     int32   `json:"retryInterval,omitempty"`
	SyntheticType     string  `json:"syntheticType"`
	Timeout           *string `json:"timeout,omitempty"`

	// HttpAction fields
	URL              *string           `json:"url,omitempty"`
	Operation        *string           `json:"operation,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
	Body             *string           `json:"body,omitempty"`
	ValidationString *string           `json:"validationString,omitempty"`
	FollowRedirect   *bool             `json:"followRedirect,omitempty"`
	AllowInsecure    *bool             `json:"allowInsecure,omitempty"`
	ExpectStatus     *int32            `json:"expectStatus,omitempty"`
	ExpectMatch      *string           `json:"expectMatch,omitempty"`
	ExpectExists     []string          `json:"expectExists,omitempty"`
	ExpectNotEmpty   []string          `json:"expectNotEmpty,omitempty"`
	ExpectJson       json.RawMessage   `json:"expectJson,omitempty"`

	// HttpScript fields
	Script     *string                       `json:"script,omitempty"`
	ScriptType *string                       `json:"scriptType,omitempty"`
	FileName   *string                       `json:"fileName,omitempty"`
	Scripts    *MultipleScriptsConfiguration `json:"scripts,omitempty"`

	// BrowserScript fields (shares Script, ScriptType, FileName, Scripts from HttpScript)
	Browser     *string `json:"browser,omitempty"`
	RecordVideo *bool   `json:"recordVideo,omitempty"`

	// DNS fields
	Lookup           *string                `json:"lookup,omitempty"`
	Server           *string                `json:"server,omitempty"`
	QueryType        *string                `json:"queryType,omitempty"`
	Port             *int32                 `json:"port,omitempty"`
	Transport        *string                `json:"transport,omitempty"`
	AcceptCNAME      *bool                  `json:"acceptCNAME,omitempty"`
	LookupServerName *bool                  `json:"lookupServerName,omitempty"`
	RecursiveLookups *bool                  `json:"recursiveLookups,omitempty"`
	ServerRetries    *int32                 `json:"serverRetries,omitempty"`
	QueryTime        *DNSFilterQueryTime    `json:"queryTime,omitempty"`
	TargetValues     []DNSFilterTargetValue `json:"targetValues,omitempty"`

	// SSLCertificate fields
	Hostname             *string                    `json:"hostname,omitempty"`
	DaysRemainingCheck   *int32                     `json:"daysRemainingCheck,omitempty"`
	AcceptSelfSignedCert *bool                      `json:"acceptSelfSignedCertificate,omitempty"`
	SSLPort              *int32                     `json:"port,omitempty"`
	ValidationRules      []SSLCertificateValidation `json:"validationRules,omitempty"`

	// WebpageAction fields (shares URL from HttpAction, Browser and RecordVideo from BrowserScript)
	// No additional fields needed

	// WebpageScript fields (shares Script, Browser, RecordVideo, FileName from BrowserScript)
	// No additional fields needed
}

type SyntheticTest struct {
	ID               string              `json:"id"`
	Label            string              `json:"label"`
	Description      *string             `json:"description,omitempty"`
	Active           bool                `json:"active"`
	ApplicationID    *string             `json:"applicationId,omitempty"`
	Applications     []string            `json:"applications,omitempty"`
	MobileApps       []string            `json:"mobileApps,omitempty"`
	Websites         []string            `json:"websites,omitempty"`
	Configuration    SyntheticTestConfig `json:"configuration"`
	CustomProperties map[string]string   `json:"customProperties,omitempty"`
	Locations        []string            `json:"locations"`
	PlaybackMode     string              `json:"playbackMode"`
	TestFrequency    *int32              `json:"testFrequency,omitempty"`
	RbacTags         []ApiTag            `json:"rbacTags,omitempty"`
	TenantId         *string             `json:"tenantId,omitempty"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject for SyntheticTest
func (s *SyntheticTest) GetIDForResourcePath() string {
	return s.ID
}
