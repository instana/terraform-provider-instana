package tf_framework

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SyntheticTestModel represents the Terraform model for a Synthetic Test
type SyntheticTestModel struct {
	ID               types.String `tfsdk:"id"`
	Label            types.String `tfsdk:"label"`
	Description      types.String `tfsdk:"description"`
	Active           types.Bool   `tfsdk:"active"`
	ApplicationID    types.String `tfsdk:"application_id"`
	CustomProperties types.Map    `tfsdk:"custom_properties"`
	Locations        types.Set    `tfsdk:"locations"`
	PlaybackMode     types.String `tfsdk:"playback_mode"`
	TestFrequency    types.Int64  `tfsdk:"test_frequency"`
	HttpAction       types.List   `tfsdk:"http_action"`
	HttpScript       types.List   `tfsdk:"http_script"`
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
}

// HttpScriptConfigModel represents the Terraform model for HTTP Script configuration
type HttpScriptConfigModel struct {
	MarkSyntheticCall types.Bool   `tfsdk:"mark_synthetic_call"`
	Retries           types.Int64  `tfsdk:"retries"`
	RetryInterval     types.Int64  `tfsdk:"retry_interval"`
	Timeout           types.String `tfsdk:"timeout"`
	Script            types.String `tfsdk:"script"`
}

// Made with Bob
