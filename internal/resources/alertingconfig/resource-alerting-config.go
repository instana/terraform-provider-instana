package alertingconfig

import (
	"context"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
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

// NewAlertingConfigResourceHandle creates the resource handle for Alerting Configuration
func NewAlertingConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.AlertingConfiguration] {
	return &alertingConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName: ResourceInstanaAlertingConfig,
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
					AlertingConfigFieldEventFilterApplicationAlertIDs: schema.SetAttribute{
						Optional:    true,
						Description: AlertingConfigDescEventFilterApplicationAlertIDs,
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeBetween(0, 1024),
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
					shared.DefaultCustomPayloadFieldsName: shared.GetCustomPayloadFieldsSchema(),
				},
			},
			SchemaVersion: 1,
		},
	}
}

// ============================================================================
// Resource Implementation
// ============================================================================

type alertingConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

// MetaData returns the resource metadata
func (r *alertingConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

// GetRestResource returns the REST resource for alerting configurations
func (r *alertingConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingConfiguration] {
	return api.AlertingConfigurations()
}

// SetComputedFields sets computed fields in the plan
func (r *alertingConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// ============================================================================
// State Management
// ============================================================================

// UpdateState updates the Terraform state with the alerting configuration data from the API
func (r *alertingConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, config *restapi.AlertingConfiguration) diag.Diagnostics {
	// Create base model with ID and alert name
	model := AlertingConfigModel{
		ID:        types.StringValue(config.ID),
		AlertName: types.StringValue(config.AlertName),
	}

	// Map integration IDs
	integrationDiags := r.mapIntegrationIDs(ctx, config, &model)
	if integrationDiags.HasError() {
		return integrationDiags
	}

	// Map event filtering configuration
	filterDiags := r.mapEventFilteringConfig(ctx, config, &model)
	if filterDiags.HasError() {
		return filterDiags
	}

	// Map custom payload fields
	payloadDiags := r.mapCustomPayloadFields(ctx, config, &model)
	if payloadDiags.HasError() {
		return payloadDiags
	}

	// Set the entire model to state
	return state.Set(ctx, model)
}

// mapIntegrationIDs maps integration IDs from API to model
func (r *alertingConfigResource) mapIntegrationIDs(ctx context.Context, config *restapi.AlertingConfiguration, model *AlertingConfigModel) diag.Diagnostics {
	integrationIDs, diags := types.SetValueFrom(ctx, types.StringType, config.IntegrationIDs)
	if !diags.HasError() {
		model.IntegrationIDs = integrationIDs
	}
	return diags
}

// mapEventFilteringConfig maps event filtering configuration from API to model
func (r *alertingConfigResource) mapEventFilteringConfig(ctx context.Context, config *restapi.AlertingConfiguration, model *AlertingConfigModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Set event filter query
	model.EventFilterQuery = util.SetStringPointerToState(config.EventFilteringConfiguration.Query)

	// Set event filter event types
	eventTypes := r.convertEventTypesToHarmonizedStringRepresentation(config.EventFilteringConfiguration.EventTypes)
	if len(eventTypes) > 0 {
		eventTypesSet, eventDiags := types.SetValueFrom(ctx, types.StringType, eventTypes)
		if eventDiags.HasError() {
			return eventDiags
		}
		model.EventFilterEventTypes = eventTypesSet
	} else {
		model.EventFilterEventTypes = types.SetNull(types.StringType)
	}

	// Set event filter rule IDs
	if len(config.EventFilteringConfiguration.RuleIDs) > 0 {
		ruleIDsSet, ruleDiags := types.SetValueFrom(ctx, types.StringType, config.EventFilteringConfiguration.RuleIDs)
		if ruleDiags.HasError() {
			return ruleDiags
		}
		model.EventFilterRuleIDs = ruleIDsSet
	} else {
		model.EventFilterRuleIDs = types.SetNull(types.StringType)
	}

	// Set event filter application alert config IDs
	if len(config.EventFilteringConfiguration.ApplicationAlertConfigIds) > 0 {
		appAlertIDsSet, appAlertDiags := types.SetValueFrom(ctx, types.StringType, config.EventFilteringConfiguration.ApplicationAlertConfigIds)
		if appAlertDiags.HasError() {
			return appAlertDiags
		}
		model.EventFilterApplicationAlertIDs = appAlertIDsSet
	} else {
		model.EventFilterApplicationAlertIDs = types.SetNull(types.StringType)
	}

	return diags
}

// mapCustomPayloadFields maps custom payload fields from API to model
func (r *alertingConfigResource) mapCustomPayloadFields(ctx context.Context, config *restapi.AlertingConfiguration, model *AlertingConfigModel) diag.Diagnostics {
	customPayloadFieldsList, diags := shared.CustomPayloadFieldsToTerraform(ctx, config.CustomerPayloadFields)
	if !diags.HasError() {
		model.CustomPayloadFields = customPayloadFieldsList
	}
	return diags
}

// MapStateToDataObject converts Terraform state to API object
// This method maps the Terraform configuration to the Instana API alerting configuration structure
func (r *alertingConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.AlertingConfiguration, diag.Diagnostics) {
	// Get model from plan or state
	model, diags := r.getConfigModelFromPlanOrState(ctx, plan, state)
	if diags.HasError() {
		return nil, diags
	}

	// Extract basic fields
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}
	alertName := model.AlertName.ValueString()

	// Map integration IDs
	integrationIDs, intDiags := r.extractIntegrationIDs(ctx, model)
	if intDiags.HasError() {
		return nil, intDiags
	}

	// Map event filtering configuration
	eventFilterConfig, filterDiags := r.extractEventFilteringConfig(ctx, model)
	if filterDiags.HasError() {
		return nil, filterDiags
	}

	// Map custom payload fields
	customerPayloadFields, payloadDiags := r.extractCustomPayloadFields(ctx, model)
	if payloadDiags.HasError() {
		return nil, payloadDiags
	}

	return &restapi.AlertingConfiguration{
		ID:                          id,
		AlertName:                   alertName,
		IntegrationIDs:              integrationIDs,
		EventFilteringConfiguration: eventFilterConfig,
		CustomerPayloadFields:       customerPayloadFields,
	}, diags
}

// getConfigModelFromPlanOrState retrieves the model from either plan or state
func (r *alertingConfigResource) getConfigModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (AlertingConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model AlertingConfigModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	return model, diags
}

// extractIntegrationIDs extracts integration IDs from the model
func (r *alertingConfigResource) extractIntegrationIDs(ctx context.Context, model AlertingConfigModel) ([]string, diag.Diagnostics) {
	var integrationIDs []string
	var diags diag.Diagnostics

	if !model.IntegrationIDs.IsNull() {
		diags.Append(model.IntegrationIDs.ElementsAs(ctx, &integrationIDs, false)...)
	}

	return integrationIDs, diags
}

// extractEventFilteringConfig extracts event filtering configuration from the model
func (r *alertingConfigResource) extractEventFilteringConfig(ctx context.Context, model AlertingConfigModel) (restapi.EventFilteringConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Map event filter query
	var query *string
	if !model.EventFilterQuery.IsNull() {
		queryStr := model.EventFilterQuery.ValueString()
		query = &queryStr
	}

	// Map event filter event types
	var eventTypeStrs []string
	if !model.EventFilterEventTypes.IsNull() {
		diags.Append(model.EventFilterEventTypes.ElementsAs(ctx, &eventTypeStrs, false)...)
		if diags.HasError() {
			return restapi.EventFilteringConfiguration{}, diags
		}
	}
	eventTypes := r.readEventTypesFromStrings(eventTypeStrs)

	// Map event filter rule IDs
	var ruleIDs []string
	if !model.EventFilterRuleIDs.IsNull() {
		diags.Append(model.EventFilterRuleIDs.ElementsAs(ctx, &ruleIDs, false)...)
		if diags.HasError() {
			return restapi.EventFilteringConfiguration{}, diags
		}
	}

	// Map event filter application alert config IDs
	var applicationAlertConfigIds []string
	if !model.EventFilterApplicationAlertIDs.IsNull() {
		diags.Append(model.EventFilterApplicationAlertIDs.ElementsAs(ctx, &applicationAlertConfigIds, false)...)
		if diags.HasError() {
			return restapi.EventFilteringConfiguration{}, diags
		}
	}

	return restapi.EventFilteringConfiguration{
		Query:                     query,
		RuleIDs:                   ruleIDs,
		EventTypes:                eventTypes,
		ApplicationAlertConfigIds: applicationAlertConfigIds,
	}, diags
}

// extractCustomPayloadFields extracts custom payload fields from the model
func (r *alertingConfigResource) extractCustomPayloadFields(ctx context.Context, model AlertingConfigModel) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
	var customerPayloadFields []restapi.CustomPayloadField[any]
	var diags diag.Diagnostics

	if !model.CustomPayloadFields.IsNull() {
		customerPayloadFields, diags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
	}

	return customerPayloadFields, diags
}

// ============================================================================
// Helper Methods
// ============================================================================

// convertEventTypesToHarmonizedStringRepresentation converts event types to lowercase string representation
// This ensures consistent representation of event types in Terraform state
func (r *alertingConfigResource) convertEventTypesToHarmonizedStringRepresentation(input []restapi.AlertEventType) []string {
	result := make([]string, len(input))
	for i, v := range input {
		value := strings.ToLower(string(v))
		result[i] = value
	}
	return result
}

// readEventTypesFromStrings converts string slice to AlertEventType slice
// This method normalizes the input by converting to lowercase before creating the event type
func (r *alertingConfigResource) readEventTypesFromStrings(input []string) []restapi.AlertEventType {
	result := make([]restapi.AlertEventType, len(input))
	for i, v := range input {
		value := strings.ToLower(v)
		result[i] = restapi.AlertEventType(value)
	}
	return result
}

// ============================================================================
// Package-level Variables and Functions
// ============================================================================

// supportedEventTypes contains all supported event types as strings for validation
var supportedEventTypes = convertSupportedEventTypesToStringSlice()

// convertSupportedEventTypesToStringSlice converts the supported event types to a string slice
// This is used for schema validation to ensure only valid event types are accepted
func convertSupportedEventTypesToStringSlice() []string {
	result := make([]string, len(restapi.SupportedAlertEventTypes))
	for i, t := range restapi.SupportedAlertEventTypes {
		result[i] = string(t)
	}
	return result
}
