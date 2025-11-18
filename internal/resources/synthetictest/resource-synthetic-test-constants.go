package synthetictest

// ResourceInstanaSyntheticTestFramework the name of the terraform-provider-instana resource to manage synthetic tests
const ResourceInstanaSyntheticTestFramework = "synthetic_test"

const (
	// Field name constants
	SyntheticTestFieldID               = "id"
	SyntheticTestFieldLabel            = "label"
	SyntheticTestFieldDescription      = "description"
	SyntheticTestFieldActive           = "active"
	SyntheticTestFieldApplicationID    = "application_id"
	SyntheticTestFieldApplications     = "applications"
	SyntheticTestFieldMobileApps       = "mobile_apps"
	SyntheticTestFieldWebsites         = "websites"
	SyntheticTestFieldCustomProperties = "custom_properties"
	SyntheticTestFieldLocations        = "locations"
	SyntheticTestFieldRbacTags         = "rbac_tags"
	SyntheticTestFieldPlaybackMode     = "playback_mode"
	SyntheticTestFieldTestFrequency    = "test_frequency"
	SyntheticTestFieldHttpAction       = "http_action"
	SyntheticTestFieldHttpScript       = "http_script"
	SyntheticTestFieldBrowserScript    = "browser_script"
	SyntheticTestFieldDNS              = "dns"
	SyntheticTestFieldSSLCertificate   = "ssl_certificate"
	SyntheticTestFieldWebpageAction    = "webpage_action"
	SyntheticTestFieldWebpageScript    = "webpage_script"

	// Common configuration field names
	SyntheticTestFieldMarkSyntheticCall           = "mark_synthetic_call"
	SyntheticTestFieldRetries                     = "retries"
	SyntheticTestFieldRetryInterval               = "retry_interval"
	SyntheticTestFieldTimeout                     = "timeout"
	SyntheticTestFieldURL                         = "url"
	SyntheticTestFieldOperation                   = "operation"
	SyntheticTestFieldHeaders                     = "headers"
	SyntheticTestFieldBody                        = "body"
	SyntheticTestFieldValidationString            = "validation_string"
	SyntheticTestFieldFollowRedirect              = "follow_redirect"
	SyntheticTestFieldAllowInsecure               = "allow_insecure"
	SyntheticTestFieldExpectStatus                = "expect_status"
	SyntheticTestFieldExpectMatch                 = "expect_match"
	SyntheticTestFieldExpectExists                = "expect_exists"
	SyntheticTestFieldExpectNotEmpty              = "expect_not_empty"
	SyntheticTestFieldExpectJson                  = "expect_json"
	SyntheticTestFieldScript                      = "script"
	SyntheticTestFieldScriptType                  = "script_type"
	SyntheticTestFieldFileName                    = "file_name"
	SyntheticTestFieldScripts                     = "scripts"
	SyntheticTestFieldBundle                      = "bundle"
	SyntheticTestFieldScriptFile                  = "script_file"
	SyntheticTestFieldBrowser                     = "browser"
	SyntheticTestFieldRecordVideo                 = "record_video"
	SyntheticTestFieldLookup                      = "lookup"
	SyntheticTestFieldServer                      = "server"
	SyntheticTestFieldQueryType                   = "query_type"
	SyntheticTestFieldPort                        = "port"
	SyntheticTestFieldTransport                   = "transport"
	SyntheticTestFieldAcceptCNAME                 = "accept_cname"
	SyntheticTestFieldLookupServerName            = "lookup_server_name"
	SyntheticTestFieldRecursiveLookups            = "recursive_lookups"
	SyntheticTestFieldServerRetries               = "server_retries"
	SyntheticTestFieldQueryTime                   = "query_time"
	SyntheticTestFieldTargetValues                = "target_values"
	SyntheticTestFieldKey                         = "key"
	SyntheticTestFieldOperator                    = "operator"
	SyntheticTestFieldValue                       = "value"
	SyntheticTestFieldHostname                    = "hostname"
	SyntheticTestFieldDaysRemainingCheck          = "days_remaining_check"
	SyntheticTestFieldAcceptSelfSignedCertificate = "accept_self_signed_certificate"
	SyntheticTestFieldValidationRules             = "validation_rules"
	SyntheticTestFieldName                        = "name"

	// Playback mode constants
	SyntheticTestPlaybackModeSimultaneous = "Simultaneous"
	SyntheticTestPlaybackModeStaggered    = "Staggered"

	// HTTP operation constants
	SyntheticTestOperationGET     = "GET"
	SyntheticTestOperationHEAD    = "HEAD"
	SyntheticTestOperationOPTIONS = "OPTIONS"
	SyntheticTestOperationPATCH   = "PATCH"
	SyntheticTestOperationPOST    = "POST"
	SyntheticTestOperationPUT     = "PUT"
	SyntheticTestOperationDELETE  = "DELETE"

	// Script type constants
	SyntheticTestScriptTypeBasic = "Basic"
	SyntheticTestScriptTypeJest  = "Jest"

	// Browser type constants
	SyntheticTestBrowserChrome  = "chrome"
	SyntheticTestBrowserFirefox = "firefox"

	// DNS query type constants
	SyntheticTestDNSQueryTypeALL            = "ALL"
	SyntheticTestDNSQueryTypeALL_CONDITIONS = "ALL_CONDITIONS"
	SyntheticTestDNSQueryTypeANY            = "ANY"
	SyntheticTestDNSQueryTypeA              = "A"
	SyntheticTestDNSQueryTypeAAAA           = "AAAA"
	SyntheticTestDNSQueryTypeCNAME          = "CNAME"
	SyntheticTestDNSQueryTypeNS             = "NS"

	// Transport protocol constants
	SyntheticTestTransportTCP = "TCP"
	SyntheticTestTransportUDP = "UDP"

	// Operator constants
	SyntheticTestOperatorCONTAINS     = "CONTAINS"
	SyntheticTestOperatorEQUALS       = "EQUALS"
	SyntheticTestOperatorGREATER_THAN = "GREATER_THAN"
	SyntheticTestOperatorIS           = "IS"
	SyntheticTestOperatorLESS_THAN    = "LESS_THAN"
	SyntheticTestOperatorMATCHES      = "MATCHES"
	SyntheticTestOperatorNOT_MATCHES  = "NOT_MATCHES"

	// Synthetic type constants
	SyntheticTestTypeHTTPAction     = "HTTPAction"
	SyntheticTestTypeHTTPScript     = "HTTPScript"
	SyntheticTestTypeBrowserScript  = "BrowserScript"
	SyntheticTestTypeDNS            = "DNS"
	SyntheticTestTypeSSLCertificate = "SSLCertificate"
	SyntheticTestTypeWebpageAction  = "WebpageAction"
	SyntheticTestTypeWebpageScript  = "WebpageScript"

	// Resource description constants
	SyntheticTestDescResource         = "This resource manages Synthetic Tests in Instana."
	SyntheticTestDescID               = "The ID of the Synthetic Test."
	SyntheticTestDescLabel            = "Friendly name of the Synthetic test."
	SyntheticTestDescDescription      = "The description of the Synthetic test."
	SyntheticTestDescActive           = "Indicates if the Synthetic test is started or not."
	SyntheticTestDescApplicationID    = "Unique identifier of the Application Perspective."
	SyntheticTestDescApplications     = "Array of application IDs"
	SyntheticTestDescMobileApps       = "Array of mobile app IDs"
	SyntheticTestDescWebsites         = "Array of website IDs"
	SyntheticTestDescCustomProperties = "Name/value pairs to provide additional information of the Synthetic test."
	SyntheticTestDescLocations        = "Array of the PoP location IDs."
	SyntheticTestDescRbacTags         = "RBAC tags for access control"
	SyntheticTestDescPlaybackMode     = "Defines how the Synthetic test should be executed across multiple PoPs."
	SyntheticTestDescTestFrequency    = "How often the playback for a Synthetic test is scheduled."

	// HTTP Action block descriptions
	SyntheticTestDescHttpAction        = "The configuration of the synthetic alert of type http action."
	SyntheticTestDescMarkSyntheticCall = "Flag used to control if HTTP calls will be marked as synthetic calls."
	SyntheticTestDescRetries           = "Indicates how many attempts will be allowed to get a successful connection."
	SyntheticTestDescRetryInterval     = "The time interval between retries in seconds."
	SyntheticTestDescTimeout           = "The timeout to be used by the PoP playback engines running the test."
	SyntheticTestDescURL               = "The URL which is being tested."
	SyntheticTestDescOperation         = "The HTTP operation."
	SyntheticTestDescHeaders           = "An object with header/value pairs."
	SyntheticTestDescBody              = "The body content to send with the operation."
	SyntheticTestDescValidationString  = "An expression to be evaluated."
	SyntheticTestDescFollowRedirect    = "A boolean type, true by default; to allow redirect."
	SyntheticTestDescAllowInsecure     = "A boolean type, if set to true then allow insecure certificates."
	SyntheticTestDescExpectStatus      = "An integer type, by default, the Synthetic passes for any 2XX status code."
	SyntheticTestDescExpectMatch       = "An optional regular expression string to be used to check the test response."
	SyntheticTestDescExpectExists      = "JSON paths that must exist in the response"
	SyntheticTestDescExpectNotEmpty    = "JSON paths that must not be empty in the response"
	SyntheticTestDescExpectJson        = "Expected JSON structure in the response"

	// HTTP Script block descriptions
	SyntheticTestDescHttpScript = "The configuration of the synthetic alert of type http script."
	SyntheticTestDescScript     = "The Javascript content in plain text."
	SyntheticTestDescScriptType = "Script type (Basic or Jest)"
	SyntheticTestDescFileName   = "Script file name"
	SyntheticTestDescScripts    = "Multiple scripts configuration for Jest"
	SyntheticTestDescBundle     = "Bundle content"
	SyntheticTestDescScriptFile = "Script file content"

	// Browser configuration descriptions
	SyntheticTestDescBrowserScript = "Browser script configuration"
	SyntheticTestDescBrowser       = "Browser type (chrome or firefox)"
	SyntheticTestDescRecordVideo   = "Record video of the test execution"

	// DNS configuration descriptions
	SyntheticTestDescDNS              = "DNS test configuration"
	SyntheticTestDescLookup           = "Domain name to lookup"
	SyntheticTestDescServer           = "DNS server to query"
	SyntheticTestDescQueryType        = "DNS query type"
	SyntheticTestDescPort             = "DNS server port"
	SyntheticTestDescTransport        = "Transport protocol (TCP or UDP)"
	SyntheticTestDescAcceptCNAME      = "Accept CNAME records"
	SyntheticTestDescLookupServerName = "Lookup server name"
	SyntheticTestDescRecursiveLookups = "Enable recursive lookups"
	SyntheticTestDescServerRetries    = "Number of server retries"
	SyntheticTestDescQueryTime        = "Query time filter"
	SyntheticTestDescTargetValues     = "Target value filters"
	SyntheticTestDescFilterKey        = "Filter key"
	SyntheticTestDescFilterOperator   = "Filter operator"
	SyntheticTestDescFilterValue      = "Filter value"

	// SSL Certificate configuration descriptions
	SyntheticTestDescSSLCertificate       = "SSL certificate test configuration"
	SyntheticTestDescHostname             = "Hostname to check SSL certificate"
	SyntheticTestDescDaysRemainingCheck   = "Minimum days remaining before certificate expiration"
	SyntheticTestDescAcceptSelfSignedCert = "Accept self-signed certificates"
	SyntheticTestDescPortNumber           = "Port number"
	SyntheticTestDescValidationRules      = "SSL certificate validation rules"
	SyntheticTestDescValidationKey        = "Validation key"
	SyntheticTestDescValidationOperator   = "Validation operator"
	SyntheticTestDescValidationValue      = "Validation value"

	// Webpage configuration descriptions
	SyntheticTestDescWebpageAction = "Webpage action test configuration"
	SyntheticTestDescWebpageScript = "Webpage script test configuration"

	// RBAC tag descriptions
	SyntheticTestDescTagName  = "Tag name"
	SyntheticTestDescTagValue = "Tag value"

	// Error message constants
	SyntheticTestErrConfigRequired       = "Configuration required"
	SyntheticTestErrConfigRequiredMsg    = "Either http_action or http_script configuration must be provided"
	SyntheticTestErrInvalidConfig        = "Invalid configuration"
	SyntheticTestErrInvalidConfigMsg     = "Only one of http_action or http_script configuration can be provided"
	SyntheticTestErrInvalidHttpAction    = "Invalid HTTP Action configuration"
	SyntheticTestErrInvalidHttpActionMsg = "Exactly one HTTP Action configuration is required"
	SyntheticTestErrInvalidHttpScript    = "Invalid HTTP Script configuration"
	SyntheticTestErrInvalidHttpScriptMsg = "Exactly one HTTP Script configuration is required"
	SyntheticTestErrNoValidConfig        = "Invalid configuration"
	SyntheticTestErrNoValidConfigMsg     = "No valid configuration provided"

	// Validator description constants
	SyntheticTestValidatorURLRegex = "must be a valid URL with HTTP or HTTPS scheme"

	// Default values
	SyntheticTestDefaultPlaybackMode   = "Simultaneous"
	SyntheticTestDefaultTestFrequency  = int64(15)
	SyntheticTestDefaultRetries        = int64(0)
	SyntheticTestDefaultRetryInterval  = int64(1)
	SyntheticTestDefaultMarkSynthetic  = false
	SyntheticTestDefaultFollowRedirect = false
	SyntheticTestDefaultAllowInsecure  = false
	SyntheticTestDefaultRecordVideo    = false
	SyntheticTestDefaultActive         = true

	// Validation limits
	SyntheticTestMinLabelLength        = 0
	SyntheticTestMaxLabelLength        = 128
	SyntheticTestMinDescriptionLength  = 0
	SyntheticTestMaxDescriptionLength  = 512
	SyntheticTestMinTestFrequency      = int64(1)
	SyntheticTestMaxTestFrequency      = int64(1440)
	SyntheticTestMinRetries            = int64(0)
	SyntheticTestMaxRetries            = int64(2)
	SyntheticTestMinRetryInterval      = int64(1)
	SyntheticTestMaxRetryInterval      = int64(10)
	SyntheticTestMinHostnameLength     = 0
	SyntheticTestMaxHostnameLength     = 2047
	SyntheticTestMinDaysRemainingCheck = int64(1)
	SyntheticTestMaxDaysRemainingCheck = int64(365)
	SyntheticTestMinViolationsCount    = int64(1)
	SyntheticTestMaxViolationsCount    = int64(12)
)
