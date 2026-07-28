package groupmapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	"github.com/instana/terraform-provider-instana/internal/util"
)

// NewGroupMappingResourceHandle creates the resource handle for RBAC Group Mappings
func NewGroupMappingResourceHandle() resourcehandle.ResourceHandle[*api.GroupMapping] {
	return &groupMappingResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaGroupMapping,
			Schema:        buildGroupMappingSchema(),
			SchemaVersion: 0,
		},
	}
}

// buildGroupMappingSchema constructs the Terraform schema for the group mapping resource
func buildGroupMappingSchema() schema.Schema {
	return schema.Schema{
		Description: GroupMappingDescResource,
		Attributes: map[string]schema.Attribute{
			GroupMappingFieldID: schema.StringAttribute{
				Computed:    true,
				Description: GroupMappingDescID,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			GroupMappingFieldKey: schema.StringAttribute{
				Required:    true,
				Description: GroupMappingDescKey,
			},
			GroupMappingFieldValue: schema.StringAttribute{
				Required:    true,
				Description: GroupMappingDescValue,
			},
			GroupMappingFieldGroupID: schema.StringAttribute{
				Required:    true,
				Description: GroupMappingDescGroupID,
			},
			GroupMappingFieldTeamID: schema.StringAttribute{
				Optional:    true,
				Description: GroupMappingDescTeamID,
			},
		},
	}
}

type groupMappingResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *groupMappingResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *groupMappingResource) GetRestResource(instanaAPI client.InstanaAPI) rest.RestResource[*api.GroupMapping] {
	return instanaAPI.GroupMappings()
}

func (r *groupMappingResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// UpdateState updates the Terraform state with data from the API response
func (r *groupMappingResource) UpdateState(_ context.Context, state *tfsdk.State, _ *tfsdk.Plan, mapping *api.GroupMapping) diag.Diagnostics {
	model := GroupMappingModel{
		ID:      types.StringValue(mapping.ID),
		Key:     types.StringValue(mapping.Key),
		Value:   types.StringValue(mapping.Value),
		GroupID: types.StringValue(mapping.GroupID),
		TeamID:  util.SetStringPointerToState(mapping.TeamID),
	}
	return state.Set(context.Background(), model)
}

// MapStateToDataObject maps Terraform state/plan to an API GroupMapping object
func (r *groupMappingResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*api.GroupMapping, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model GroupMappingModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	id := ""
	if !model.ID.IsNull() && !model.ID.IsUnknown() {
		id = model.ID.ValueString()
	}

	var teamID *string
	if !model.TeamID.IsNull() && !model.TeamID.IsUnknown() {
		v := model.TeamID.ValueString()
		teamID = &v
	}

	return &api.GroupMapping{
		ID:      id,
		Key:     model.Key.ValueString(),
		Value:   model.Value.ValueString(),
		GroupID: model.GroupID.ValueString(),
		TeamID:  teamID,
	}, diags
}

// GetStateUpgraders returns the state upgraders for this resource (none needed for v0)
func (r *groupMappingResource) GetStateUpgraders(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
