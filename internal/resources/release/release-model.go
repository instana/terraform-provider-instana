package release

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
)

// ReleaseModel represents the Terraform state model for a release
type ReleaseModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Start        types.Int64  `tfsdk:"start"`
	LastUpdated  types.Int64  `tfsdk:"last_updated"`
	Applications types.List   `tfsdk:"applications"`
	Services     types.List   `tfsdk:"services"`
}

// ApplicationModel represents a single application scope in a release or in scoped_to
type ApplicationModel struct {
	Name types.String `tfsdk:"name"`
}

// ServiceModel represents a single service scope in the release
type ServiceModel struct {
	Name     types.String   `tfsdk:"name"`
	ScopedTo *ScopedToModel `tfsdk:"scoped_to"`
}

// ScopedToModel restricts a service to specific application perspectives.
// The API requires at least one application when scoped_to is present.
type ScopedToModel struct {
	Applications types.List `tfsdk:"applications"`
}

type releaseResource struct {
	metaData resourcehandle.ResourceMetaData
}
