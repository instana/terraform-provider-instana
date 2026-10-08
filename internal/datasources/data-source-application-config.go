package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/client"
	models "github.com/instana/instana-go-client/shared/types"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// ApplicationConfigDataSourceModel represents the data model for the application configuration data source
type ApplicationConfigDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Scope         types.String `tfsdk:"scope"`
	BoundaryScope types.String `tfsdk:"boundary_scope"`
	TagFilter     types.String `tfsdk:"tag_filter"`
	AccessRules   types.List   `tfsdk:"access_rules"`
}

// ApplicationConfigAccessRuleDataSourceModel represents an access rule model in the application configuration data source
type ApplicationConfigAccessRuleDataSourceModel struct {
	AccessType   types.String `tfsdk:"access_type"`
	RelatedID    types.String `tfsdk:"related_id"`
	RelationType types.String `tfsdk:"relation_type"`
}

// NewApplicationConfigDataSource creates a new data source for application config
func NewApplicationConfigDataSource() datasource.DataSource {
	return &applicationConfigDataSource{}
}

type applicationConfigDataSource struct {
	instanaAPI client.InstanaAPI
}

func (d *applicationConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaApplicationConfig
}

func (d *applicationConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: ApplicationConfigDataSourceDesc,
		Attributes: map[string]schema.Attribute{
			ApplicationConfigDataSourceFieldID: schema.StringAttribute{
				Description: ApplicationConfigDataSourceDescID,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot(ApplicationConfigDataSourceFieldID), path.MatchRoot(ApplicationConfigDataSourceFieldName)),
				},
			},
			ApplicationConfigDataSourceFieldName: schema.StringAttribute{
				Description: ApplicationConfigDataSourceDescName,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot(ApplicationConfigDataSourceFieldID), path.MatchRoot(ApplicationConfigDataSourceFieldName)),
				},
			},
			ApplicationConfigDataSourceFieldLabel: schema.StringAttribute{
				Description: ApplicationConfigDataSourceDescLabel,
				Computed:    true,
			},
			ApplicationConfigDataSourceFieldScope: schema.StringAttribute{
				Description: ApplicationConfigDataSourceDescScope,
				Computed:    true,
			},
			ApplicationConfigDataSourceFieldBoundaryScope: schema.StringAttribute{
				Description: ApplicationConfigDataSourceDescBoundaryScope,
				Computed:    true,
			},
			ApplicationConfigDataSourceFieldTagFilter: schema.StringAttribute{
				Description: ApplicationConfigDataSourceDescTagFilter,
				Computed:    true,
			},
			ApplicationConfigDataSourceFieldAccessRules: schema.ListNestedAttribute{
				Description: ApplicationConfigDataSourceDescAccessRules,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						ApplicationConfigDataSourceFieldAccessType: schema.StringAttribute{
							Description: ApplicationConfigDataSourceDescAccessType,
							Computed:    true,
						},
						ApplicationConfigDataSourceFieldRelatedID: schema.StringAttribute{
							Description: ApplicationConfigDataSourceDescRelatedID,
							Computed:    true,
						},
						ApplicationConfigDataSourceFieldRelationType: schema.StringAttribute{
							Description: ApplicationConfigDataSourceDescRelationType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *applicationConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(*shared.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			ApplicationConfigDataSourceErrUnexpectedConfigureType,
			fmt.Sprintf("Expected *shared.ProviderMeta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.instanaAPI = providerMeta.InstanaAPI
}

func (d *applicationConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApplicationConfigDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	name := data.Name.ValueString()
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Attribute", ApplicationConfigDataSourceErrMissingLookupAttribute)
		return
	}

	var appConfig *api.ApplicationConfig
	var err error

	if id != "" {
		appConfig, err = d.instanaAPI.ApplicationConfigs().GetOne(id)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading application configuration",
				fmt.Sprintf(ApplicationConfigDataSourceErrReadByID, id, err),
			)
			return
		}
		if appConfig == nil {
			resp.Diagnostics.AddError(
				"Not found",
				fmt.Sprintf(ApplicationConfigDataSourceErrNotFoundByID, id),
			)
			return
		}
	} else {
		appConfigs, err := d.instanaAPI.ApplicationConfigs().GetAll()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading application configurations",
				fmt.Sprintf(ApplicationConfigDataSourceErrReadAll, err),
			)
			return
		}

		for _, config := range *appConfigs {
			if config.Label == name {
				appConfig = config
				break
			}
		}

		if appConfig == nil {
			resp.Diagnostics.AddError(
				"Not found",
				fmt.Sprintf(ApplicationConfigDataSourceErrNotFoundByName, name),
			)
			return
		}
	}

	data.ID = types.StringValue(appConfig.ID)
	data.Name = types.StringValue(appConfig.Label)
	data.Label = types.StringValue(appConfig.Label)
	data.Scope = types.StringValue(string(appConfig.Scope))
	data.BoundaryScope = types.StringValue(string(appConfig.BoundaryScope))

	// Map tag filter expression to normalized string
	if appConfig.TagFilterExpression != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(appConfig.TagFilterExpression)
		if err != nil {
			resp.Diagnostics.AddError(
				ApplicationConfigDataSourceErrConvertingTagFilter,
				fmt.Sprintf(ApplicationConfigDataSourceErrFailedToConvert, err),
			)
			return
		}
		data.TagFilter = util.SetStringPointerToState(normalizedTagFilterString)
	} else {
		data.TagFilter = types.StringNull()
	}

	// Map access rules
	accessRulesList, accessRulesDiags := d.mapAccessRulesToState(ctx, appConfig.AccessRules)
	resp.Diagnostics.Append(accessRulesDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.AccessRules = accessRulesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *applicationConfigDataSource) mapAccessRulesToState(ctx context.Context, accessRules []models.AccessRule) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	accessRuleAttrTypes := map[string]attr.Type{
		ApplicationConfigDataSourceFieldAccessType:   types.StringType,
		ApplicationConfigDataSourceFieldRelatedID:    types.StringType,
		ApplicationConfigDataSourceFieldRelationType: types.StringType,
	}

	if len(accessRules) == 0 {
		emptyList, listDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: accessRuleAttrTypes,
		}, []ApplicationConfigAccessRuleDataSourceModel{})
		diags.Append(listDiags...)
		return emptyList, diags
	}

	ruleModels := make([]ApplicationConfigAccessRuleDataSourceModel, len(accessRules))
	for i, rule := range accessRules {
		ruleModels[i] = ApplicationConfigAccessRuleDataSourceModel{
			AccessType:   types.StringValue(string(rule.AccessType)),
			RelationType: types.StringValue(string(rule.RelationType)),
			RelatedID:    util.SetStringPointerToState(rule.RelatedID),
		}
	}

	listValue, listDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: accessRuleAttrTypes,
	}, ruleModels)
	diags.Append(listDiags...)

	return listValue, diags
}
