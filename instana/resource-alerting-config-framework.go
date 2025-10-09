package instana

import (
	"context"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
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

// ResourceInstanaAlertingConfigFramework the name of the terraform-provider-instana resource to manage alerting configurations
const ResourceInstanaAlertingConfigFramework = "alerting_config"

// AlertingConfigModel represents the data model for the alerting configuration resource
type AlertingConfigModel struct {
	ID                    types.String `tfsdk:"id"`
	AlertName             types.String `tfsdk:"alert_name"`
	IntegrationIDs        types.Set    `tfsdk:"integration_ids"`
	EventFilterQuery      types.String `tfsdk:"event_filter_query"`
	EventFilterEventTypes types.Set    `tfsdk:"event_filter_event_types"`
	EventFilterRuleIDs    types.Set    `tfsdk:"event_filter_rule_ids"`
	CustomPayloadFields   types.List   `tfsdk:"custom_payload_field"`
}

// NewAlertingConfigResourceHandleFramework creates the resource handle for Alerting Configuration
func NewAlertingConfigResourceHandleFramework() ResourceHandleFramework[*restapi.AlertingConfiguration] {
	return &alertingConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaAlertingConfigFramework,
			Schema: schema.Schema{
				Description: "This resource manages alerting configurations in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the alerting configuration.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					AlertingConfigFieldAlertName: schema.StringAttribute{
						Required:    true,
						Description: "Configures the alert name of the alerting configuration",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 256),
						},
					},
					AlertingConfigFieldIntegrationIds: schema.SetAttribute{
						Required:    true,
						Description: "Configures the list of Integration IDs (Alerting Channels).",
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeBetween(0, 1024),
						},
					},
					AlertingConfigFieldEventFilterQuery: schema.StringAttribute{
						Optional:    true,
						Description: "Configures a filter query to to filter rules or event types for a limited set of entities",
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 2048),
						},
					},
					AlertingConfigFieldEventFilterEventTypes: schema.SetAttribute{
						Optional:    true,
						Description: "Configures the list of Event Types IDs which should trigger an alert.",
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
						Description: "Configures the list of Rule IDs which should trigger an alert.",
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeBetween(0, 1024),
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.RequiresReplace(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					DefaultCustomPayloadFieldsName: schema.ListNestedBlock{
						Description: "Custom payload fields for the alerting configuration.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								CustomPayloadFieldsFieldKey: schema.StringAttribute{
									Required:    true,
									Description: "The key of the custom payload field",
								},
								CustomPayloadFieldsFieldStaticStringValue: schema.StringAttribute{
									Optional:    true,
									Description: "The value of a static string custom payload field",
								},
								CustomPayloadFieldsFieldDynamicValue: schema.ListNestedAttribute{
									Optional:    true,
									Description: "The value of a dynamic custom payload field",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											CustomPayloadFieldsFieldDynamicKey: schema.StringAttribute{
												Optional:    true,
												Description: "The key of the dynamic custom payload field",
											},
											CustomPayloadFieldsFieldDynamicTagName: schema.StringAttribute{
												Required:    true,
												Description: "The name of the tag of the dynamic custom payload field",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			SchemaVersion: 2,
		},
	}
}

type alertingConfigResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *alertingConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *alertingConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingConfiguration] {
	return api.AlertingConfigurations()
}

func (r *alertingConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *alertingConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, config *restapi.AlertingConfiguration) diag.Diagnostics {

	var diags diag.Diagnostics
	var model AlertingConfigModel

	diags.Append(state.Get(ctx, &model)...)

	// Set integration IDs
	integrationIDs, diags := types.SetValueFrom(ctx, types.StringType, config.IntegrationIDs)
	if diags.HasError() {
		return diags
	}
	model.IntegrationIDs = integrationIDs

	// Set event filter query
	if config.EventFilteringConfiguration.Query != nil {
		model.EventFilterQuery = types.StringValue(*config.EventFilteringConfiguration.Query)
	} else {
		model.EventFilterQuery = types.StringNull()
	}

	// Set event filter event types
	eventTypes := r.convertEventTypesToHarmonizedStringRepresentation(config.EventFilteringConfiguration.EventTypes)
	eventTypesSet, diags := types.SetValueFrom(ctx, types.StringType, eventTypes)
	if diags.HasError() {
		return diags
	}
	model.EventFilterEventTypes = eventTypesSet

	// Set event filter rule IDs
	ruleIDsSet, diags := types.SetValueFrom(ctx, types.StringType, config.EventFilteringConfiguration.RuleIDs)
	if diags.HasError() {
		return diags
	}
	model.EventFilterRuleIDs = ruleIDsSet

	model.CustomPayloadFields = types.ListNull(types.ObjectType{})

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
	var integrationIDs []string
	if !model.IntegrationIDs.IsNull() {
		diags.Append(model.IntegrationIDs.ElementsAs(ctx, &integrationIDs, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

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
			return nil, diags
		}
	}
	eventTypes := r.readEventTypesFromStrings(eventTypeStrs)

	// Map event filter rule IDs
	var ruleIDs []string
	if !model.EventFilterRuleIDs.IsNull() {
		diags.Append(model.EventFilterRuleIDs.ElementsAs(ctx, &ruleIDs, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map custom payload fields
	var customerPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayloadFields.IsNull() {
		var payloadDiags diag.Diagnostics
		customerPayloadFields, payloadDiags = BuildCustomPayloadFieldsTyped(ctx, model.CustomPayloadFields)
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

// Made with Bob
