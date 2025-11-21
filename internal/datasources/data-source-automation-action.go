package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Constants are now defined in data-source-automation-action-constants.go

// AutomationActionDataSourceModel represents the data model for the automation action data source
type AutomationActionDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Tags        types.List   `tfsdk:"tags"`
}

// NewAutomationActionDataSource creates a new data source for automation actions
func NewAutomationActionDataSource() datasource.DataSource {
	return &automationActionDataSource{}
}

type automationActionDataSource struct {
	instanaAPI restapi.InstanaAPI
}

func (d *automationActionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaAutomationAction
}

func (d *automationActionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: AutomationActionDescDataSource,
		Attributes: map[string]schema.Attribute{
			AutomationActionFieldID: schema.StringAttribute{
				Description: AutomationActionDescID,
				Computed:    true,
			},
			AutomationActionFieldName: schema.StringAttribute{
				Description: AutomationActionDescName,
				Required:    true,
			},
			AutomationActionFieldDescription: schema.StringAttribute{
				Description: AutomationActionDescDescription,
				Computed:    true,
			},
			AutomationActionFieldType: schema.StringAttribute{
				Description: AutomationActionDescType,
				Required:    true,
			},
			AutomationActionFieldTags: schema.ListAttribute{
				Description: AutomationActionDescTags,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *automationActionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(*restapi.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			AutomationActionErrUnexpectedConfigureType,
			fmt.Sprintf(AutomationActionErrUnexpectedConfigureTypeDetail, req.ProviderData),
		)
		return
	}

	d.instanaAPI = providerMeta.InstanaAPI
}

func (d *automationActionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AutomationActionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the name and type from the configuration
	name := data.Name.ValueString()
	actionType := data.Type.ValueString()

	// Get all automation actions
	actions, err := d.instanaAPI.AutomationActions().GetAll()
	if err != nil {
		resp.Diagnostics.AddError(
			AutomationActionErrReadingActions,
			fmt.Sprintf(AutomationActionErrReadingActionsDetail, err),
		)
		return
	}

	// Find the action with the matching name and type
	var matchingAction *restapi.AutomationAction
	for _, action := range *actions {
		if action.Name == name && strings.EqualFold(action.Type, actionType) {
			matchingAction = action
			break
		}
	}

	if matchingAction == nil {
		resp.Diagnostics.AddError(
			AutomationActionErrNotFound,
			fmt.Sprintf(AutomationActionErrNotFoundDetail, name, actionType),
		)
		return
	}

	// Update the data model with the action details
	data.ID = types.StringValue(matchingAction.ID)
	data.Description = types.StringValue(matchingAction.Description)

	// Handle tags based on their type
	if matchingAction.Tags != nil {
		switch v := matchingAction.Tags.(type) {
		case []interface{}:
			elements := make([]types.String, len(v))
			for i, tag := range v {
				if strTag, ok := tag.(string); ok {
					elements[i] = types.StringValue(strTag)
				} else {
					elements[i] = types.StringValue(fmt.Sprintf("%v", tag))
				}
			}
			tagsList, diags := types.ListValueFrom(ctx, types.StringType, elements)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			data.Tags = tagsList
		default:
			data.Tags = types.ListNull(types.StringType)
		}
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
