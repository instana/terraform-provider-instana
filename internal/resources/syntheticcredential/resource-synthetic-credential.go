package syntheticcredential

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
)

// NewSyntheticCredentialResourceHandle creates the resource handle for Synthetic Credentials
func NewSyntheticCredentialResourceHandle() resourcehandle.ResourceHandle[*api.SyntheticCredential] {
	return &syntheticCredentialResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:     ResourceInstanaSyntheticCredential,
			Schema:           buildSyntheticCredentialSchema(),
			SchemaVersion:    1,
			SkipIDGeneration: true,
			ResourceIDField:  strPtr(SyntheticCredentialFieldCredentialName),
		},
	}
}

func strPtr(s string) *string {
	return &s
}

func buildSyntheticCredentialSchema() schema.Schema {
	return schema.Schema{
		Description: SyntheticCredentialDescResource,
		Attributes: map[string]schema.Attribute{
			SyntheticCredentialFieldCredentialName: schema.StringAttribute{
				Required:    true,
				Description: SyntheticCredentialDescCredentialName,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			SyntheticCredentialFieldCredentialValue: schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: SyntheticCredentialDescCredentialValue,
				// The API never returns this value. UseStateForUnknown keeps the existing
				// value in state across reads so no spurious diff is produced.
				// Optional (not Required) so that `terraform import` succeeds — import
				// only populates credential_name; the user must add credential_value to
				// their config and run `terraform apply` afterwards.
				// Enforcement that the value is present happens in MapStateToDataObject
				// (called only during apply, not plan) so that post-import plans succeed.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			SyntheticCredentialFieldApplications: schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Description: SyntheticCredentialDescApplications,
				ElementType: types.StringType,
			},
			SyntheticCredentialFieldMobileApps: schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Description: SyntheticCredentialDescMobileApps,
				ElementType: types.StringType,
			},
			SyntheticCredentialFieldWebsites: schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Description: SyntheticCredentialDescWebsites,
				ElementType: types.StringType,
			},
			SyntheticCredentialFieldRbacTags: schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: SyntheticCredentialDescRbacTags,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						SyntheticCredentialFieldRbacTagID: schema.StringAttribute{
							Required:    true,
							Description: SyntheticCredentialDescRbacTagID,
						},
						SyntheticCredentialFieldRbacTagName: schema.StringAttribute{
							Required:    true,
							Description: SyntheticCredentialDescRbacTagName,
						},
					},
				},
			},
		},
	}
}

type syntheticCredentialResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *syntheticCredentialResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *syntheticCredentialResource) GetRestResource(instanaAPI client.InstanaAPI) rest.RestResource[*api.SyntheticCredential] {
	return instanaAPI.SyntheticCredentials()
}

func (r *syntheticCredentialResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *syntheticCredentialResource) GetStateUpgraders(_ context.Context) map[int64]resource.StateUpgrader {
	return nil
}

// MapStateToDataObject maps the Terraform plan/state to the API model
func (r *syntheticCredentialResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*api.SyntheticCredential, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SyntheticCredentialModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}
	if diags.HasError() {
		return nil, diags
	}

	cred := &api.SyntheticCredential{
		CredentialName: model.CredentialName.ValueString(),
	}
	// credential_value is write-only: the API never returns it. We only include it
	// in the payload when the user has actually provided it in their config.
	// If it is null here during a create/update (plan is non-nil), the user omitted
	// it from their config — return an error rather than silently sending an empty
	// value to the API.
	// During delete (state is non-nil, plan is nil) the value may be null if the
	// resource was imported without it; that is fine because delete only needs the name.
	if !model.CredentialValue.IsNull() && !model.CredentialValue.IsUnknown() {
		cred.CredentialValue = model.CredentialValue.ValueString()
	} else if plan != nil {
		// plan != nil means this is a create or update — value is required.
		diags.AddAttributeError(
			path.Root(SyntheticCredentialFieldCredentialValue),
			"Missing required attribute",
			fmt.Sprintf(
				"%s must be set before applying. Add it to your configuration and run terraform apply again.",
				SyntheticCredentialFieldCredentialValue,
			),
		)
		return nil, diags
	}

	// Applications
	if !model.Applications.IsNull() && !model.Applications.IsUnknown() {
		var apps []string
		diags.Append(model.Applications.ElementsAs(ctx, &apps, false)...)
		cred.Applications = apps
	}

	// MobileApps
	if !model.MobileApps.IsNull() && !model.MobileApps.IsUnknown() {
		var mobileApps []string
		diags.Append(model.MobileApps.ElementsAs(ctx, &mobileApps, false)...)
		cred.MobileApps = mobileApps
	}

	// Websites
	if !model.Websites.IsNull() && !model.Websites.IsUnknown() {
		var websites []string
		diags.Append(model.Websites.ElementsAs(ctx, &websites, false)...)
		cred.Websites = websites
	}

	// RbacTags
	if !model.RbacTags.IsNull() && !model.RbacTags.IsUnknown() {
		var tagModels []SyntheticCredentialRbacTagModel
		diags.Append(model.RbacTags.ElementsAs(ctx, &tagModels, false)...)
		if !diags.HasError() {
			for _, t := range tagModels {
				cred.RbacTags = append(cred.RbacTags, api.RbacTag{
					ID:          t.ID.ValueString(),
					DisplayName: t.DisplayName.ValueString(),
				})
			}
		}
	}

	return cred, diags
}

// UpdateState maps the API response back to the Terraform state.
//
// The associations GET endpoint returns only credentialName, applications, applicationLabels,
// createdAt and modifiedAt — it does NOT return credentialValue, rbacTags, mobileApps or
// websites. Those fields are preserved verbatim from the plan (after create/update) or the
// previous state (after read) so that Terraform does not see a spurious diff.
func (r *syntheticCredentialResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *api.SyntheticCredential) diag.Diagnostics {
	var diags diag.Diagnostics

	// src is the plan when called after create/update, or the prior state when called after read.
	type stateGetter interface {
		GetAttribute(context.Context, path.Path, interface{}) diag.Diagnostics
	}
	var src stateGetter
	if plan != nil {
		src = plan
	} else if state != nil {
		src = state
	}

	// Preserve fields that the API never echoes back.
	var credentialValue types.String
	var rbacTags types.Set
	var mobileApps types.Set
	var websites types.Set

	if src != nil {
		diags.Append(src.GetAttribute(ctx, path.Root(SyntheticCredentialFieldCredentialValue), &credentialValue)...)
		diags.Append(src.GetAttribute(ctx, path.Root(SyntheticCredentialFieldRbacTags), &rbacTags)...)
		diags.Append(src.GetAttribute(ctx, path.Root(SyntheticCredentialFieldMobileApps), &mobileApps)...)
		diags.Append(src.GetAttribute(ctx, path.Root(SyntheticCredentialFieldWebsites), &websites)...)
	}
	if diags.HasError() {
		return diags
	}

	// Ensure preserved sets are never null/unknown to avoid framework type errors.
	if rbacTags.IsNull() || rbacTags.IsUnknown() {
		rbacTags = emptyRbacTagSet()
	}
	if mobileApps.IsNull() || mobileApps.IsUnknown() {
		mobileApps = types.SetValueMust(types.StringType, nil)
	}
	if websites.IsNull() || websites.IsUnknown() {
		websites = types.SetValueMust(types.StringType, nil)
	}

	model := SyntheticCredentialModel{
		CredentialName:  types.StringValue(apiObject.CredentialName),
		CredentialValue: credentialValue,
		// applications is returned by the associations endpoint — use the API value.
		Applications: stringSliceToSet(apiObject.Applications),
		// The remaining fields are not returned by the API — preserve from plan/state.
		MobileApps: mobileApps,
		Websites:   websites,
		RbacTags:   rbacTags,
	}

	diags.Append(state.Set(ctx, &model)...)
	return diags
}

// emptyRbacTagSet returns an empty set typed for the rbac_tag nested object.
func emptyRbacTagSet() types.Set {
	return types.SetValueMust(types.ObjectType{AttrTypes: rbacTagAttrTypes()}, nil)
}

// stringSliceToSet converts a string slice to a types.Set
func stringSliceToSet(items []string) types.Set {
	if len(items) == 0 {
		return types.SetValueMust(types.StringType, []attr.Value{})
	}
	vals := make([]attr.Value, len(items))
	for i, s := range items {
		vals[i] = types.StringValue(s)
	}
	set, _ := types.SetValue(types.StringType, vals)
	return set
}

// rbacTagAttrTypes returns the attribute types map for the rbac_tag nested object
func rbacTagAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		SyntheticCredentialFieldRbacTagID:   types.StringType,
		SyntheticCredentialFieldRbacTagName: types.StringType,
	}
}

