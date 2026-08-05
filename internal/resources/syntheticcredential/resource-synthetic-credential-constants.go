package syntheticcredential

// ResourceInstanaSyntheticCredential the name of the terraform-provider-instana resource to manage synthetic credentials
const ResourceInstanaSyntheticCredential = "synthetic_credential"

const (
	// Field name constants
	SyntheticCredentialFieldCredentialName  = "credential_name"
	SyntheticCredentialFieldCredentialValue = "credential_value"
	SyntheticCredentialFieldApplications    = "applications"
	SyntheticCredentialFieldMobileApps      = "mobile_apps"
	SyntheticCredentialFieldWebsites        = "websites"
	SyntheticCredentialFieldRbacTags        = "rbac_tags"
	SyntheticCredentialFieldRbacTagID       = "id"
	SyntheticCredentialFieldRbacTagName     = "display_name"

	// Description constants
	SyntheticCredentialDescResource       = "This resource manages Synthetic Credentials in Instana."
	SyntheticCredentialDescCredentialName = "The unique name of the credential. This serves as the identifier for the resource."
	SyntheticCredentialDescCredentialValue = "The secret value of the credential. This field is write-only and will not be read back from the API."
	SyntheticCredentialDescApplications   = "List of application IDs the credential is scoped to."
	SyntheticCredentialDescMobileApps     = "List of mobile app IDs the credential is scoped to."
	SyntheticCredentialDescWebsites       = "List of website IDs the credential is scoped to."
	SyntheticCredentialDescRbacTags       = "RBAC tags for access control."
	SyntheticCredentialDescRbacTagID      = "The ID of the RBAC tag."
	SyntheticCredentialDescRbacTagName    = "The display name of the RBAC tag."
)
