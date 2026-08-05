package syntheticcredential

import "github.com/hashicorp/terraform-plugin-framework/types"

// SyntheticCredentialModel is the Terraform state model for a Synthetic Credential
type SyntheticCredentialModel struct {
	CredentialName  types.String `tfsdk:"credential_name"`
	CredentialValue types.String `tfsdk:"credential_value"`
	Applications    types.Set    `tfsdk:"applications"`
	MobileApps      types.Set    `tfsdk:"mobile_apps"`
	Websites        types.Set    `tfsdk:"websites"`
	RbacTags        types.Set    `tfsdk:"rbac_tags"`
}

// SyntheticCredentialRbacTagModel represents an RBAC tag within the credential resource
type SyntheticCredentialRbacTagModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
}
