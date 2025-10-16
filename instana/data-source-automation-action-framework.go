package instana

import (
	"context"
	"fmt"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceInstanaAutomationActionFramework the name of the terraform-provider-instana data source to read automation actions
const DataSourceInstanaAutomationActionFramework = "automation_action"

// AutomationActionDataSourceModel represents the data model for the automation action data source
type AutomationActionDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Tags        types.List   `tfsdk:"tags"`
}

// NewAutomationActionDataSourceFramework creates a new data source for automation actions
func NewAutomationActionDataSourceFramework() datasource.DataSource {
	return &automationActionDataSourceFramework{}
}

type automationActionDataSourceFramework struct {
	instanaAPI restapi.InstanaAPI
}

func (d *automationActionDataSourceFramework) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaAutomationActionFramework
}

func (d *automationActionDataSourceFramework) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for an Instana automation action. Automation actions are used to execute scripts or HTTP requests.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the automation action.",
				Computed:    true,
			},
			AutomationActionFieldName: schema.StringAttribute{
				Description: "The name of the automation action.",
				Required:    true,
			},
			AutomationActionFieldDescription: schema.StringAttribute{
				Description: "The description of the automation action.",
				Computed:    true,
			},
			AutomationActionFieldType: schema.StringAttribute{
				Description: "The type of the automation action.",
				Required:    true,
			},
			AutomationActionFieldTags: schema.ListAttribute{
				Description: "The tags of the automation action.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *automationActionDataSourceFramework) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(*ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderMeta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.instanaAPI = providerMeta.InstanaAPI
}

func (d *automationActionDataSourceFramework) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
			"Error reading automation actions",
			fmt.Sprintf("Could not read automation actions: %s", err),
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
			"Automation action not found",
			fmt.Sprintf("No automation action found with name '%s' and type '%s'", name, actionType),
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

// Made with Bob
