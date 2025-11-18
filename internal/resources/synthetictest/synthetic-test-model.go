package synthetictest

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SyntheticTestModel represents the Terraform model for a Synthetic Test
type SyntheticTestModel struct {
	ID               types.String               `tfsdk:"id"`
	Label            types.String               `tfsdk:"label"`
	Description      types.String               `tfsdk:"description"`
	Active           types.Bool                 `tfsdk:"active"`
	ApplicationID    types.String               `tfsdk:"application_id"`
	Applications     types.Set                  `tfsdk:"applications"`
	MobileApps       types.Set                  `tfsdk:"mobile_apps"`
	Websites         types.Set                  `tfsdk:"websites"`
	CustomProperties types.Map                  `tfsdk:"custom_properties"`
	Locations        types.Set                  `tfsdk:"locations"`
	PlaybackMode     types.String               `tfsdk:"playback_mode"`
	TestFrequency    types.Int64                `tfsdk:"test_frequency"`
	RbacTags         types.Set                  `tfsdk:"rbac_tags"`
	HttpAction       *HttpActionConfigModel     `tfsdk:"http_action"`
	HttpScript       *HttpScriptConfigModel     `tfsdk:"http_script"`
	BrowserScript    *BrowserScriptConfigModel  `tfsdk:"browser_script"`
	DNS              *DNSConfigModel            `tfsdk:"dns"`
	SSLCertificate   *SSLCertificateConfigModel `tfsdk:"ssl_certificate"`
	WebpageAction    *WebpageActionConfigModel  `tfsdk:"webpage_action"`
	WebpageScript    *WebpageScriptConfigModel  `tfsdk:"webpage_script"`
}

// RbacTagModel represents an RBAC tag
type RbacTagModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// MultipleScriptsModel represents multiple scripts configuration
type MultipleScriptsModel struct {
	Bundle     types.String `tfsdk:"bundle"`
	ScriptFile types.String `tfsdk:"script_file"`
}

// HttpActionConfigModel represents the Terraform model for HTTP Action configuration
type HttpActionConfigModel struct {
	MarkSyntheticCall types.Bool   `tfsdk:"mark_synthetic_call"`
	Retries           types.Int64  `tfsdk:"retries"`
	RetryInterval     types.Int64  `tfsdk:"retry_interval"`
	Timeout           types.String `tfsdk:"timeout"`
	URL               types.String `tfsdk:"url"`
	Operation         types.String `tfsdk:"operation"`
	Headers           types.Map    `tfsdk:"headers"`
	Body              types.String `tfsdk:"body"`
	ValidationString  types.String `tfsdk:"validation_string"`
	FollowRedirect    types.Bool   `tfsdk:"follow_redirect"`
	AllowInsecure     types.Bool   `tfsdk:"allow_insecure"`
	ExpectStatus      types.Int64  `tfsdk:"expect_status"`
	ExpectMatch       types.String `tfsdk:"expect_match"`
	ExpectExists      types.Set    `tfsdk:"expect_exists"`
	ExpectNotEmpty    types.Set    `tfsdk:"expect_not_empty"`
	ExpectJson        types.String `tfsdk:"expect_json"`
}

// HttpScriptConfigModel represents the Terraform model for HTTP Script configuration
type HttpScriptConfigModel struct {
	MarkSyntheticCall types.Bool            `tfsdk:"mark_synthetic_call"`
	Retries           types.Int64           `tfsdk:"retries"`
	RetryInterval     types.Int64           `tfsdk:"retry_interval"`
	Timeout           types.String          `tfsdk:"timeout"`
	Script            types.String          `tfsdk:"script"`
	ScriptType        types.String          `tfsdk:"script_type"`
	FileName          types.String          `tfsdk:"file_name"`
	Scripts           *MultipleScriptsModel `tfsdk:"scripts"`
}

// BrowserScriptConfigModel represents the Terraform model for Browser Script configuration
type BrowserScriptConfigModel struct {
	MarkSyntheticCall types.Bool            `tfsdk:"mark_synthetic_call"`
	Retries           types.Int64           `tfsdk:"retries"`
	RetryInterval     types.Int64           `tfsdk:"retry_interval"`
	Timeout           types.String          `tfsdk:"timeout"`
	Script            types.String          `tfsdk:"script"`
	ScriptType        types.String          `tfsdk:"script_type"`
	FileName          types.String          `tfsdk:"file_name"`
	Scripts           *MultipleScriptsModel `tfsdk:"scripts"`
	Browser           types.String          `tfsdk:"browser"`
	RecordVideo       types.Bool            `tfsdk:"record_video"`
}

// DNSFilterQueryTimeModel represents DNS query time filter
type DNSFilterQueryTimeModel struct {
	Key      types.String `tfsdk:"key"`
	Operator types.String `tfsdk:"operator"`
	Value    types.Int64  `tfsdk:"value"`
}

// DNSFilterTargetValueModel represents DNS target value filter
type DNSFilterTargetValueModel struct {
	Key      types.String `tfsdk:"key"`
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
}

// DNSConfigModel represents the Terraform model for DNS configuration
type DNSConfigModel struct {
	MarkSyntheticCall types.Bool               `tfsdk:"mark_synthetic_call"`
	Retries           types.Int64              `tfsdk:"retries"`
	RetryInterval     types.Int64              `tfsdk:"retry_interval"`
	Timeout           types.String             `tfsdk:"timeout"`
	Lookup            types.String             `tfsdk:"lookup"`
	Server            types.String             `tfsdk:"server"`
	QueryType         types.String             `tfsdk:"query_type"`
	Port              types.Int64              `tfsdk:"port"`
	Transport         types.String             `tfsdk:"transport"`
	AcceptCNAME       types.Bool               `tfsdk:"accept_cname"`
	LookupServerName  types.Bool               `tfsdk:"lookup_server_name"`
	RecursiveLookups  types.Bool               `tfsdk:"recursive_lookups"`
	ServerRetries     types.Int64              `tfsdk:"server_retries"`
	QueryTime         *DNSFilterQueryTimeModel `tfsdk:"query_time"`
	TargetValues      types.Set                `tfsdk:"target_values"`
}

// SSLCertificateValidationModel represents SSL certificate validation rule
type SSLCertificateValidationModel struct {
	Key      types.String `tfsdk:"key"`
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
}

// SSLCertificateConfigModel represents the Terraform model for SSL Certificate configuration
type SSLCertificateConfigModel struct {
	MarkSyntheticCall    types.Bool   `tfsdk:"mark_synthetic_call"`
	Retries              types.Int64  `tfsdk:"retries"`
	RetryInterval        types.Int64  `tfsdk:"retry_interval"`
	Timeout              types.String `tfsdk:"timeout"`
	Hostname             types.String `tfsdk:"hostname"`
	DaysRemainingCheck   types.Int64  `tfsdk:"days_remaining_check"`
	AcceptSelfSignedCert types.Bool   `tfsdk:"accept_self_signed_certificate"`
	Port                 types.Int64  `tfsdk:"port"`
	ValidationRules      types.Set    `tfsdk:"validation_rules"`
}

// WebpageActionConfigModel represents the Terraform model for Webpage Action configuration
type WebpageActionConfigModel struct {
	MarkSyntheticCall types.Bool   `tfsdk:"mark_synthetic_call"`
	Retries           types.Int64  `tfsdk:"retries"`
	RetryInterval     types.Int64  `tfsdk:"retry_interval"`
	Timeout           types.String `tfsdk:"timeout"`
	URL               types.String `tfsdk:"url"`
	Browser           types.String `tfsdk:"browser"`
	RecordVideo       types.Bool   `tfsdk:"record_video"`
}

// WebpageScriptConfigModel represents the Terraform model for Webpage Script configuration
type WebpageScriptConfigModel struct {
	MarkSyntheticCall types.Bool   `tfsdk:"mark_synthetic_call"`
	Retries           types.Int64  `tfsdk:"retries"`
	RetryInterval     types.Int64  `tfsdk:"retry_interval"`
	Timeout           types.String `tfsdk:"timeout"`
	Script            types.String `tfsdk:"script"`
	FileName          types.String `tfsdk:"file_name"`
	Browser           types.String `tfsdk:"browser"`
	RecordVideo       types.Bool   `tfsdk:"record_video"`
}
