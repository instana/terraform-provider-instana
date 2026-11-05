package businessperspectiveconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"
	tag "github.com/instana/instana-go-client/shared/tagfilter"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// NewBusinessPerspectiveConfigResourceHandle creates the resource handle for Business Perspective Configuration
func NewBusinessPerspectiveConfigResourceHandle() resourcehandle.ResourceHandle[*api.BusinessPerspectiveConfig] {
	return &businessPerspectiveConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaBusinessPerspectiveConfig,
			Schema:        createResourceSchema(),
			SchemaVersion: BusinessPerspectiveConfigSchemaVersion,
		},
	}
}

func createResourceSchema() schema.Schema {
	return schema.Schema{
		Description: BusinessPerspectiveConfigDescResource,
		Attributes: map[string]schema.Attribute{
			BusinessPerspectiveConfigFieldID: schema.StringAttribute{
				Computed:    true,
				Description: BusinessPerspectiveConfigDescID,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			BusinessPerspectiveConfigFieldName: schema.StringAttribute{
				Required:    true,
				Description: BusinessPerspectiveConfigDescName,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 100),
				},
			},
			BusinessPerspectiveConfigFieldDescription: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: BusinessPerspectiveConfigDescDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 300),
				},
			},
			BusinessPerspectiveConfigFieldTagFilter: schema.StringAttribute{
				Optional:    true,
				Description: BusinessPerspectiveConfigDescTagFilter,
			},
			BusinessPerspectiveConfigFieldRbacTags: schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: BusinessPerspectiveConfigDescRbacTags,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						BusinessPerspectiveConfigFieldRbacTagID: schema.StringAttribute{
							Required:    true,
							Description: BusinessPerspectiveConfigDescRbacTagID,
						},
						BusinessPerspectiveConfigFieldRbacTagDisplayName: schema.StringAttribute{
							Required:    true,
							Description: BusinessPerspectiveConfigDescRbacTagDisplayName,
						},
					},
				},
			},
		},
	}
}

type businessPerspectiveConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *businessPerspectiveConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *businessPerspectiveConfigResource) GetRestResource(apiClient client.InstanaAPI) rest.RestResource[*api.BusinessPerspectiveConfig] {
	return apiClient.BusinessPerspectives()
}

func (r *businessPerspectiveConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *businessPerspectiveConfigResource) GetStateUpgraders(_ context.Context) map[int64]resource.StateUpgrader {
	return nil
}

func (r *businessPerspectiveConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, config *api.BusinessPerspectiveConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	if ctx == nil {
		diags.AddError("Invalid Context", "Context cannot be nil")
		return diags
	}
	if state == nil {
		diags.AddError("Invalid State", "State cannot be nil")
		return diags
	}
	if config == nil {
		diags.AddError("Invalid API Object", "API object cannot be nil")
		return diags
	}

	var model BusinessPerspectiveConfigModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return diags
	}

	model.ID = types.StringValue(config.ID)
	model.Name = types.StringValue(config.Name)

	if config.Description != "" {
		model.Description = types.StringValue(config.Description)
	} else if !model.Description.IsNull() && !model.Description.IsUnknown() && model.Description.ValueString() == "" {
		model.Description = types.StringValue("")
	} else {
		model.Description = types.StringNull()
	}

	if model.TagFilter.IsNull() || model.TagFilter.IsUnknown() {
		tagFilter, tagFilterDiags := r.mapTagFilterToModel(config.TagFilterExpression)
		diags.Append(tagFilterDiags...)
		model.TagFilter = tagFilter
	}

	model.RbacTags = r.mapRbacTagsToState(config.RbacTags)

	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *businessPerspectiveConfigResource) mapTagFilterToModel(tagFilterExpression *tag.TagFilter) (types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tagFilterExpression == nil {
		return types.StringNull(), diags
	}

	normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(tagFilterExpression)
	if err != nil {
		diags.AddError(
			BusinessPerspectiveConfigErrConvertingTagFilter,
			fmt.Sprintf(BusinessPerspectiveConfigErrFailedToConvert, err),
		)
		return types.StringNull(), diags
	}

	return util.SetStringPointerToState(normalizedTagFilterString), diags
}

func (r *businessPerspectiveConfigResource) mapRbacTagsToState(rbacTags []api.RbacTag) types.List {
	tagAttrTypes := map[string]attr.Type{
		BusinessPerspectiveConfigFieldRbacTagID:          types.StringType,
		BusinessPerspectiveConfigFieldRbacTagDisplayName: types.StringType,
	}

	if len(rbacTags) == 0 {
		emptyList, _ := types.ListValue(
			types.ObjectType{AttrTypes: tagAttrTypes},
			[]attr.Value{},
		)
		return emptyList
	}

	tagValues := make([]attr.Value, len(rbacTags))
	for i, tag := range rbacTags {
		tagObj, _ := types.ObjectValue(
			tagAttrTypes,
			map[string]attr.Value{
				BusinessPerspectiveConfigFieldRbacTagID:          types.StringValue(tag.ID),
				BusinessPerspectiveConfigFieldRbacTagDisplayName: types.StringValue(tag.DisplayName),
			},
		)
		tagValues[i] = tagObj
	}

	list, _ := types.ListValue(types.ObjectType{AttrTypes: tagAttrTypes}, tagValues)
	return list
}

func (r *businessPerspectiveConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*api.BusinessPerspectiveConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	if ctx == nil {
		diags.AddError("Invalid Context", "Context cannot be nil")
		return nil, diags
	}
	if plan == nil && state == nil {
		diags.AddError("Invalid Configuration", "Either plan or state must be provided")
		return nil, diags
	}

	var model BusinessPerspectiveConfigModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	id := ""
	if !model.ID.IsNull() && !model.ID.IsUnknown() {
		id = model.ID.ValueString()
	}

	var description string
	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		description = model.Description.ValueString()
	}

	var tagFilter *tag.TagFilter
	if !model.TagFilter.IsNull() && !model.TagFilter.IsUnknown() {
		parser := tagfilter.NewParser()
		expr, err := parser.Parse(model.TagFilter.ValueString())
		if err != nil {
			diags.AddError(
				BusinessPerspectiveConfigErrParsingTagFilter,
				fmt.Sprintf(BusinessPerspectiveConfigErrFailedToParse, err),
			)
			return nil, diags
		}
		mapper := tagfilter.NewMapper()
		tagFilter = mapper.ToAPIModel(expr)
	}

	rbacTags := r.mapRbacTagsFromState(ctx, model.RbacTags)

	return &api.BusinessPerspectiveConfig{
		ID:                  id,
		Name:                model.Name.ValueString(),
		Description:         description,
		TagFilterExpression: tagFilter,
		RbacTags:            rbacTags,
	}, diags
}

func (r *businessPerspectiveConfigResource) mapRbacTagsFromState(ctx context.Context, rbacTagsList types.List) []api.RbacTag {
	if rbacTagsList.IsNull() || rbacTagsList.IsUnknown() {
		return []api.RbacTag{}
	}

	var tagModels []RbacTagModel
	rbacTagsList.ElementsAs(ctx, &tagModels, false)

	if len(tagModels) == 0 {
		return []api.RbacTag{}
	}

	rbacTags := make([]api.RbacTag, len(tagModels))
	for i, tagModel := range tagModels {
		rbacTags[i] = api.RbacTag{
			ID:          tagModel.ID.ValueString(),
			DisplayName: tagModel.DisplayName.ValueString(),
		}
	}

	return rbacTags
}
