package alertingconfig

import (
	"context"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewAlertingConfigResourceHandleFramework creates the resource handle for Alerting Configuration
func NewAlertingConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.AlertingConfiguration] {
	return &alertingConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaAlertingConfigFramework,
			Schema: schema.Schema{
				Description: AlertingConfigDescResource,
				Attributes: map[string]schema.Attribute{
					AlertingConfigFieldID: schema.StringAttribute{
						Computed:    true,
						Description: AlertingConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					AlertingConfigFieldAlertName: schema.StringAttribute{
						Required:    true,
						Description: AlertingConfigDescAlertName,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 256),
						},
					},
					AlertingConfigFieldIntegrationIds: schema.SetAttribute{
						Required:    true,
						Description: AlertingConfigDescIntegrationIds,
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeBetween(0, 1024),
						},
					},
					AlertingConfigFieldEventFilterQuery: schema.StringAttribute{
						Optional:    true,
						Description: AlertingConfigDescEventFilterQuery,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 2048),
						},
					},
					AlertingConfigFieldEventFilterEventTypes: schema.SetAttribute{
						Optional:    true,
						Description: AlertingConfigDescEventFilterEventTypes,
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(supportedEventTypes...),
							),
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.RequiresReplace(),
						},
					},
					AlertingConfigFieldEventFilterRuleIDs: schema.SetAttribute{
						Optional:    true,
						Description: AlertingConfigDescEventFilterRuleIDs,
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeBetween(0, 1024),
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.RequiresReplace(),
						},
					},
					shared.DefaultCustomPayloadFieldsName: shared.GetCustomPayloadFieldsSchema(),
				},
			},
			SchemaVersion: 2,
		},
	}
}

type alertingConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *alertingConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *alertingConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingConfiguration] {
	return api.AlertingConfigurations()
}

func (r *alertingConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *alertingConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, config *restapi.AlertingConfiguration) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the config
	model := AlertingConfigModel{
		ID:        types.StringValue(config.ID),
		AlertName: types.StringValue(config.AlertName),
	}

	// Set integration IDs
	model.IntegrationIDs = config.IntegrationIDs

	// Set event filter query
	model.EventFilterQuery = util.SetStringPointerToState(config.EventFilteringConfiguration.Query)

	// Set event filter event types
	eventTypes := r.convertEventTypesToHarmonizedStringRepresentation(config.EventFilteringConfiguration.EventTypes)
	model.EventFilterEventTypes = eventTypes

	// Set event filter rule IDs
	model.EventFilterRuleIDs = config.EventFilteringConfiguration.RuleIDs

	// Convert custom payload fields to the appropriate Terraform types
	// Using the utility function from tfutils package for better maintainability and reusability
	// This handles both static string and dynamic custom payload field types
	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, config.CustomerPayloadFields)
	if payloadDiags.HasError() {
		return payloadDiags
	}
	model.CustomPayloadFields = customPayloadFieldsList

	// Set the entire model to state
	diags = state.Set(ctx, model)
	if diags.HasError() {
		return diags
	}

	return nil
}

func (r *alertingConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.AlertingConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model AlertingConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map ID
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}

	// Map alert name
	alertName := model.AlertName.ValueString()

	// Map integration IDs
	integrationIDs := model.IntegrationIDs

	// Map event filter query
	var query *string
	if !model.EventFilterQuery.IsNull() {
		queryStr := model.EventFilterQuery.ValueString()
		query = &queryStr
	}

	// Map event filter event types
	eventTypes := r.readEventTypesFromStrings(model.EventFilterEventTypes)

	// Map event filter rule IDs
	ruleIDs := model.EventFilterRuleIDs

	// Map custom payload fields
	var customerPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayloadFields.IsNull() {
		var payloadDiags diag.Diagnostics
		customerPayloadFields, payloadDiags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
		if payloadDiags.HasError() {
			diags.Append(payloadDiags...)
			return nil, diags
		}
	}

	return &restapi.AlertingConfiguration{
		ID:             id,
		AlertName:      alertName,
		IntegrationIDs: integrationIDs,
		EventFilteringConfiguration: restapi.EventFilteringConfiguration{
			Query:      query,
			RuleIDs:    ruleIDs,
			EventTypes: eventTypes,
		},
		CustomerPayloadFields: customerPayloadFields,
	}, diags
}

func (r *alertingConfigResourceFramework) convertEventTypesToHarmonizedStringRepresentation(input []restapi.AlertEventType) []string {
	result := make([]string, len(input))
	for i, v := range input {
		value := strings.ToLower(string(v))
		result[i] = value
	}
	return result
}

func (r *alertingConfigResourceFramework) readEventTypesFromStrings(input []string) []restapi.AlertEventType {
	result := make([]restapi.AlertEventType, len(input))
	for i, v := range input {
		value := strings.ToLower(v)
		result[i] = restapi.AlertEventType(value)
	}
	return result
}

var supportedEventTypes = convertSupportedEventTypesToStringSlice()

func convertSupportedEventTypesToStringSlice() []string {
	result := make([]string, len(restapi.SupportedAlertEventTypes))
	for i, t := range restapi.SupportedAlertEventTypes {
		result[i] = string(t)
	}
	return result
}
