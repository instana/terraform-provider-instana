package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SloConfigModel represents the data model for the SLO configuration resource
type SloConfigModel struct {
	ID         types.String    `tfsdk:"id"`
	Name       types.String    `tfsdk:"name"`
	Target     types.Float64   `tfsdk:"target"`
	Tags       []types.String  `tfsdk:"tags"`
	RbacTags   []RbacTagModel  `tfsdk:"rbac_tags"`
	Entity     EntityModel     `tfsdk:"entity"`
	Indicator  IndicatorModel  `tfsdk:"indicator"`
	TimeWindow TimeWindowModel `tfsdk:"time_window"`
}

// RbacTagModel represents an RBAC tag in the Terraform model
type RbacTagModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	ID          types.String `tfsdk:"id"`
}
type TimeWindowModel struct {
	FixedTimeWindowModel   *FixedTimeWindowModel   `tfsdk:"fixed"`
	RollingTimeWindowModel *RollingTimeWindowModel `tfsdk:"rolling"`
}
type EntityModel struct {
	ApplicationEntityModel *ApplicationEntityModel `tfsdk:"application"`
	WebsiteEntityModel     *WebsiteEntityModel     `tfsdk:"website"`
	SyntheticEntityModel   *SyntheticEntityModel   `tfsdk:"synthetic"`
}
type IndicatorModel struct {
	TimeBasedLatencyIndicatorModel       *TimeBasedLatencyIndicatorModel       `tfsdk:"time_based_latency"`
	EventBasedLatencyIndicatorModel      *EventBasedLatencyIndicatorModel      `tfsdk:"event_based_latency"`
	TimeBasedAvailabilityIndicatorModel  *TimeBasedAvailabilityIndicatorModel  `tfsdk:"time_based_availability"`
	EventBasedAvailabilityIndicatorModel *EventBasedAvailabilityIndicatorModel `tfsdk:"event_based_availability"`
	TrafficIndicatorModel                *TrafficIndicatorModel                `tfsdk:"traffic"`
	CustomIndicatorModel                 *CustomIndicatorModel                 `tfsdk:"custom"`
}

// ApplicationEntityModel represents an application entity in the Terraform model
type ApplicationEntityModel struct {
	ApplicationID    types.String `tfsdk:"application_id"`
	ServiceID        types.String `tfsdk:"service_id"`
	EndpointID       types.String `tfsdk:"endpoint_id"`
	BoundaryScope    types.String `tfsdk:"boundary_scope"`
	IncludeSynthetic types.Bool   `tfsdk:"include_synthetic"`
	IncludeInternal  types.Bool   `tfsdk:"include_internal"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// WebsiteEntityModel represents a website entity in the Terraform model
type WebsiteEntityModel struct {
	WebsiteID        types.String `tfsdk:"website_id"`
	BeaconType       types.String `tfsdk:"beacon_type"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// SyntheticEntityModel represents a synthetic entity in the Terraform model
type SyntheticEntityModel struct {
	SyntheticTestIDs []types.String `tfsdk:"synthetic_test_ids"`
	FilterExpression types.String   `tfsdk:"filter_expression"`
}

// TimeBasedLatencyIndicatorModel represents a time-based latency indicator in the Terraform model
type TimeBasedLatencyIndicatorModel struct {
	Threshold   types.Float64 `tfsdk:"threshold"`
	Aggregation types.String  `tfsdk:"aggregation"`
}

// EventBasedLatencyIndicatorModel represents an event-based latency indicator in the Terraform model
type EventBasedLatencyIndicatorModel struct {
	Threshold types.Float64 `tfsdk:"threshold"`
}

// TimeBasedAvailabilityIndicatorModel represents a time-based availability indicator in the Terraform model
type TimeBasedAvailabilityIndicatorModel struct {
	Threshold   types.Float64 `tfsdk:"threshold"`
	Aggregation types.String  `tfsdk:"aggregation"`
}

// EventBasedAvailabilityIndicatorModel represents an event-based availability indicator in the Terraform model
type EventBasedAvailabilityIndicatorModel struct {
	// No fields needed for this indicator type
}

// TrafficIndicatorModel represents a traffic indicator in the Terraform model
type TrafficIndicatorModel struct {
	TrafficType types.String  `tfsdk:"traffic_type"`
	Threshold   types.Float64 `tfsdk:"threshold"`
	Aggregation types.String  `tfsdk:"aggregation"`
	Operator    types.String  `tfsdk:"operator"`
}

// CustomIndicatorModel represents a custom indicator in the Terraform model
type CustomIndicatorModel struct {
	GoodEventFilterExpression types.String `tfsdk:"good_event_filter_expression"`
	BadEventFilterExpression  types.String `tfsdk:"bad_event_filter_expression"`
}

// RollingTimeWindowModel represents a rolling time window in the Terraform model
type RollingTimeWindowModel struct {
	Duration     types.Int64  `tfsdk:"duration"`
	DurationUnit types.String `tfsdk:"duration_unit"`
	Timezone     types.String `tfsdk:"timezone"`
}

// FixedTimeWindowModel represents a fixed time window in the Terraform model
type FixedTimeWindowModel struct {
	Duration       types.Int64   `tfsdk:"duration"`
	DurationUnit   types.String  `tfsdk:"duration_unit"`
	Timezone       types.String  `tfsdk:"timezone"`
	StartTimestamp types.Float64 `tfsdk:"start_timestamp"`
}

// ResourceInstanaSloConfigFramework the name of the terraform-provider-instana resource to manage SLO configurations
const ResourceInstanaSloConfigFramework = "slo_config"

// SloConfigFieldRbacTags is the field name for RBAC tags
const SloConfigFieldRbacTags = "rbac_tags"

// NewSloConfigResourceHandleFramework creates the resource handle for SLO Config
func NewSloConfigResourceHandleFramework() ResourceHandleFramework[*restapi.SloConfig] {
	return &sloConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaSloConfigFramework,
			Schema: schema.Schema{
				Description: "This resource manages SLO Configurations in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the SLO configuration",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					SloConfigFieldName: schema.StringAttribute{
						Required:    true,
						Description: "The name of the SLO configuration",
					},
					SloConfigFieldTarget: schema.Float64Attribute{
						Required:    true,
						Description: "The target of the SLO configuration",
					},
					SloConfigFieldTags: schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "The tags of the SLO configuration",
					},
					SloConfigFieldRbacTags: schema.ListNestedAttribute{
						Optional:    true,
						Description: "RBAC tags for the SLO configuration",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"display_name": schema.StringAttribute{
									Required:    true,
									Description: "Display name of the RBAC tag",
								},
								"id": schema.StringAttribute{
									Required:    true,
									Description: "ID of the RBAC tag",
								},
							},
						},
					},
				},
				Blocks: map[string]schema.Block{
					SloConfigFieldSloEntity: schema.SingleNestedBlock{
						Description: "The entity to use for the SLO configuration",
						Blocks: map[string]schema.Block{
							SloConfigApplicationEntity: schema.SingleNestedBlock{
								Description: "Application entity of SLO",
								Attributes: map[string]schema.Attribute{
									SloConfigFieldApplicationID: schema.StringAttribute{
										Optional:    true,
										Description: "The application ID of the entity",
									},
									SloConfigFieldBoundaryScope: schema.StringAttribute{
										Optional:    true,
										Description: "The boundary scope for the entity configuration (ALL, INBOUND)",
									},
									SloConfigFieldFilterExpression: schema.StringAttribute{
										Optional:    true,
										Description: "Entity filter",
									},
									SloConfigFieldIncludeInternal: schema.BoolAttribute{
										Optional:    true,
										Description: "Optional flag to indicate whether also internal calls are included",
									},
									SloConfigFieldIncludeSynthetic: schema.BoolAttribute{
										Optional:    true,
										Description: "Optional flag to indicate whether also synthetic calls are included in the scope or not",
									},
									SloConfigFieldServiceID: schema.StringAttribute{
										Optional:    true,
										Description: "The service ID of the entity",
									},
									SloConfigFieldEndpointID: schema.StringAttribute{
										Optional:    true,
										Description: "The endpoint ID of the entity",
									},
								},
							},
							SloConfigWebsiteEntity: schema.SingleNestedBlock{
								Description: "Website entity of SLO",
								Attributes: map[string]schema.Attribute{
									SloConfigFieldWebsiteID: schema.StringAttribute{
										Optional:    true,
										Description: "The website ID of the entity",
									},
									SloConfigFieldFilterExpression: schema.StringAttribute{
										Optional:    true,
										Description: "Entity filter",
									},
									SloConfigFieldBeaconType: schema.StringAttribute{
										Optional:    true,
										Description: "The beacon type for the entity configuration (pageLoad, resourceLoad, httpRequest, error, custom, pageChange)",
									},
								},
							},
							SloConfigSyntheticEntity: schema.SingleNestedBlock{
								Description: "Synthetic entity of SLO",
								Attributes: map[string]schema.Attribute{
									SloConfigFieldSyntheticTestIDs: schema.ListAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Description: "The synthetics ID of the entity",
									},
									SloConfigFieldFilterExpression: schema.StringAttribute{
										Optional:    true,
										Description: "Entity filter",
									},
								},
							},
						},
					},
					SloConfigFieldSloIndicator: schema.SingleNestedBlock{
						Description: "The indicator to use for the SLO configuration",
						Blocks: map[string]schema.Block{
							"time_based_latency": schema.SingleNestedBlock{
								Description: "Time-based latency indicator",
								Attributes: map[string]schema.Attribute{
									"threshold": schema.Float64Attribute{
										Optional:    true,
										Description: "The threshold for the metric configuration",
									},
									"aggregation": schema.StringAttribute{
										Optional:    true,
										Description: "The aggregation type for the metric configuration",
									},
								},
							},
							"event_based_latency": schema.SingleNestedBlock{
								Description: "Event-based latency indicator",
								Attributes: map[string]schema.Attribute{
									"threshold": schema.Float64Attribute{
										Optional:    true,
										Description: "The threshold for the metric configuration",
									},
								},
							},
							"time_based_availability": schema.SingleNestedBlock{
								Description: "Time-based availability indicator",
								Attributes: map[string]schema.Attribute{
									"threshold": schema.Float64Attribute{
										Optional:    true,
										Description: "The threshold for the metric configuration",
									},
									"aggregation": schema.StringAttribute{
										Optional:    true,
										Description: "The aggregation type for the metric configuration",
									},
								},
							},
							"event_based_availability": schema.SingleNestedBlock{
								Description: "Event-based availability indicator",
							},
							"traffic": schema.SingleNestedBlock{
								Description: "Traffic indicator",
								Attributes: map[string]schema.Attribute{
									"traffic_type": schema.StringAttribute{
										Optional:    true,
										Description: "The traffic type for the indicator",
									},
									"threshold": schema.Float64Attribute{
										Optional:    true,
										Description: "The threshold for the metric configuration",
									},
									"operator": schema.StringAttribute{
										Optional:    true,
										Description: "The aggregation type for the metric configuration",
										Validators: []validator.String{
											stringvalidator.OneOf(">", ">=", "<", "<="),
										},
									},
								},
							},
							"custom": schema.SingleNestedBlock{
								Description: "Custom indicator",
								Attributes: map[string]schema.Attribute{
									"good_event_filter_expression": schema.StringAttribute{
										Optional:    true,
										Description: "Good event filter expression",
									},
									"bad_event_filter_expression": schema.StringAttribute{
										Optional:    true,
										Description: "Bad event filter expression",
									},
								},
							},
						},
					},
					SloConfigFieldSloTimeWindow: schema.SingleNestedBlock{
						Description: "The time window to use for the SLO configuration",
						Blocks: map[string]schema.Block{
							"rolling": schema.SingleNestedBlock{
								Description: "Rolling time window",
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Optional:    true,
										Description: "The duration of the time window",
									},
									"duration_unit": schema.StringAttribute{
										Optional:    true,
										Description: "The duration unit of the time window (day, week)",
									},
									"timezone": schema.StringAttribute{
										Optional:    true,
										Description: "The timezone for the SLO configuration",
									},
								},
							},
							"fixed": schema.SingleNestedBlock{
								Description: "Fixed time window",
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Optional:    true,
										Description: "The duration of the time window",
									},
									"duration_unit": schema.StringAttribute{
										Optional:    true,
										Description: "The duration unit of the time window (day, week)",
									},
									"timezone": schema.StringAttribute{
										Optional:    true,
										Description: "The timezone for the SLO configuration",
									},
									"start_timestamp": schema.Float64Attribute{
										Optional:    true,
										Description: "Time window start time",
									},
								},
							},
						},
					},
				},
			},
			SchemaVersion:    1,
			SkipIDGeneration: true,
		},
	}
}

type sloConfigResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *sloConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *sloConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloConfig] {
	return api.SloConfigs()
}

func (r *sloConfigResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	var diags diag.Diagnostics
	diags.Append(plan.SetAttribute(ctx, path.Root("id"), types.StringValue(SloConfigFromTerraformIdPrefix+RandomID()))...)
	return diags
}

func (r *sloConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SloConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SloConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	} else {
		diags.AddError(
			"Error mapping state to data object",
			"Both plan and state are nil",
		)
		return nil, diags
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map ID
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}

	// Map name and target
	name := model.Name.ValueString()
	target := model.Target.ValueFloat64()

	// Convert tags to []interface{}
	var tagsList []string
	for _, tag := range model.Tags {
		if !tag.IsNull() && !tag.IsUnknown() {
			tagsList = append(tagsList, tag.ValueString())
		}
	}

	// Convert RBAC tags to []interface{}
	rbacTagsList := mapRbacTags(ctx, model.RbacTags)

	// Get entity data
	entityData, entityDiags := r.mapEntityFromState(ctx, model.Entity)
	diags.Append(entityDiags...)

	// Get indicator data
	indicator, indicatorDiags := r.mapIndicatorFromState(ctx, model.Indicator)
	diags.Append(indicatorDiags...)

	timeWindowData, timeWindowDiags := r.mapTimeWindowFromState(ctx, model.TimeWindow)
	diags.Append(timeWindowDiags...)

	if diags.HasError() {
		return nil, diags
	}

	// Create SLO config object
	sloConfig := &restapi.SloConfig{
		ID:         id,
		Name:       name,
		Target:     target,
		Tags:       tagsList,
		Entity:     entityData,
		Indicator:  indicator,
		TimeWindow: timeWindowData,
		RbacTags:   rbacTagsList,
	}

	return sloConfig, diags
}

func mapRbacTags(ctx context.Context, rbacTags []RbacTagModel) []restapi.RbacTag {
	var rbacTagsList []restapi.RbacTag
	for _, t := range rbacTags {
		rbacTagsList = append(rbacTagsList, restapi.RbacTag{
			DisplayName: t.DisplayName.ValueString(),
			ID:          t.ID.ValueString(),
		})
	}

	return rbacTagsList
}

func (r *sloConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SloConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the API object
	model := SloConfigModel{
		ID:     types.StringValue(apiObject.ID),
		Name:   types.StringValue(apiObject.Name),
		Target: types.Float64Value(apiObject.Target),
	}

	// Set tags if present
	if apiObject.Tags != nil {
		var tags []types.String
		for _, tag := range apiObject.Tags {
			tags = append(tags, types.StringValue(tag))
		}
		model.Tags = tags
	}

	// Set RBAC tags if present
	if apiObject.RbacTags != nil {
		var rbacTags []RbacTagModel
		for _, tag := range apiObject.RbacTags {
			rbacTags = append(rbacTags, RbacTagModel{
				DisplayName: types.StringValue(tag.DisplayName),
				ID:          types.StringValue(tag.ID),
			})
		}
		model.RbacTags = rbacTags
	}

	// Map entity
	entityData, entityDiags := r.mapEntityToState(ctx, apiObject)
	diags.Append(entityDiags...)
	if !diags.HasError() {
		model.Entity = entityData
	}

	// Map indicator
	indicatorData, indicatorDiags := r.mapIndicatorToState(ctx, apiObject)
	diags.Append(indicatorDiags...)
	if !diags.HasError() {
		model.Indicator = indicatorData
	}

	// Map time window
	timeWindowData, timeWindowDiags := r.mapTimeWindowToState(ctx, apiObject)
	diags.Append(timeWindowDiags...)
	if !diags.HasError() {
		model.TimeWindow = timeWindowData
	}

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	if diags.HasError() {
		return diags
	}

	return diags
}

// Helper methods for mapping entity from plan
func (r *sloConfigResourceFramework) mapEntityFromState(ctx context.Context, entityObj EntityModel) (restapi.SloEntity, diag.Diagnostics) {
	var diags diag.Diagnostics
	// Check for application entity
	if entityObj.ApplicationEntityModel != nil {
		applicationModel := entityObj.ApplicationEntityModel
		if applicationModel.ApplicationID.IsUnknown() || applicationModel.ApplicationID.IsNull() ||
			applicationModel.BoundaryScope.IsNull() || applicationModel.BoundaryScope.IsUnknown() {
			diags.AddError(
				"application_id and boundary_scope are required for application entity",
				"application_id and boundary_scope are required for application entity",
			)
			return restapi.SloEntity{}, diags
		}
		applicationIdStr := applicationModel.ApplicationID.ValueString()
		serviceID := setStringPointer(applicationModel.ServiceID)
		endpointID := setStringPointer(applicationModel.EndpointID)
		boundaryScope := setStringPointer(applicationModel.BoundaryScope)
		includeInternal := applicationModel.IncludeInternal.ValueBool()
		includeSynthetic := applicationModel.IncludeSynthetic.ValueBool()

		appEntityObj := restapi.SloEntity{
			Type:             SloConfigApplicationEntity,
			ApplicationID:    &applicationIdStr,
			ServiceID:        serviceID,
			EndpointID:       endpointID,
			BoundaryScope:    boundaryScope,
			IncludeInternal:  &includeInternal,
			IncludeSynthetic: &includeSynthetic,
		}

		// Convert filter expression to API model if set
		var tagFilter *restapi.TagFilter
		tagFilter, tagDiags := mapTagFilterFromState(applicationModel.FilterExpression, tagFilter)
		if tagDiags.HasError() {
			return restapi.SloEntity{}, tagDiags
		}
		appEntityObj.FilterExpression = tagFilter
		return appEntityObj, diags
	}

	// Check for website entity
	if entityObj.WebsiteEntityModel != nil {
		websiteModel := entityObj.WebsiteEntityModel

		if websiteModel.WebsiteID.IsNull() || websiteModel.WebsiteID.IsUnknown() ||
			websiteModel.BeaconType.IsNull() || websiteModel.BeaconType.IsUnknown() {
			diags.AddError(
				"website_id and beacon_type are required for website entity",
				"website_id and beacon_type are required for website entity",
			)
			return restapi.SloEntity{}, diags
		}

		websiteIdStr := setStringPointer(websiteModel.WebsiteID)
		beaconTypeStr := setStringPointer(websiteModel.BeaconType)

		websiteEntityObj := restapi.SloEntity{
			Type:       SloConfigWebsiteEntity,
			WebsiteId:  websiteIdStr,
			BeaconType: beaconTypeStr,
		}

		// Convert filter expression to API model if set
		var tagFilter *restapi.TagFilter
		tagFilter, tagDiags := mapTagFilterFromState(websiteModel.FilterExpression, tagFilter)
		if tagDiags.HasError() {
			return restapi.SloEntity{}, tagDiags
		}
		websiteEntityObj.FilterExpression = tagFilter
		return websiteEntityObj, diags
	}

	// Check for synthetic entity
	if entityObj.SyntheticEntityModel != nil {
		syntheticModel := entityObj.SyntheticEntityModel

		if len(syntheticModel.SyntheticTestIDs) == 0 {
			diags.AddError(
				"synthetic_test_ids is required for synthetic entity",
				"synthetic_test_ids is required for synthetic entity",
			)
			return restapi.SloEntity{}, diags
		}

		// Convert synthetic test IDs to []interface{}
		var testIDs []interface{}
		for _, id := range syntheticModel.SyntheticTestIDs {
			if !id.IsNull() && !id.IsUnknown() {
				testIDs = append(testIDs, id.ValueString())
			}
		}

		syntheticEntityObj := restapi.SloEntity{
			Type:             SloConfigSyntheticEntity,
			SyntheticTestIDs: testIDs,
		}

		// Convert filter expression to API model if set
		var tagFilter *restapi.TagFilter
		tagFilter, tagDiags := mapTagFilterFromState(syntheticModel.FilterExpression, tagFilter)
		if tagDiags.HasError() {
			return restapi.SloEntity{}, tagDiags
		}
		syntheticEntityObj.FilterExpression = tagFilter
		return syntheticEntityObj, diags
	}

	diags.AddError(
		"Missing entity configuration",
		"Exactly one entity configuration is required",
	)
	return restapi.SloEntity{}, diags
}

func mapTagFilterFromState(filterExpression types.String, tagFilter *restapi.TagFilter) (*restapi.TagFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	if !filterExpression.IsNull() && !filterExpression.IsUnknown() {
		parser := tagfilter.NewParser()
		expr, err := parser.Parse(filterExpression.ValueString())
		if err != nil {
			diags.AddError(
				"Error parsing filter expression",
				fmt.Sprintf("Could not parse filter expression: %s", err),
			)
			return nil, diags
		}
		mapper := tagfilter.NewMapper()
		tagFilter = mapper.ToAPIModel(expr)
	} else {
		operator := restapi.LogicalOperatorType("AND")
		tagFilter = &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		}
	}
	return tagFilter, diags
}

// Helper methods for mapping indicator from plan
func (r *sloConfigResourceFramework) mapIndicatorFromState(ctx context.Context, indicatorModel IndicatorModel) (restapi.SloIndicator, diag.Diagnostics) {
	var diags diag.Diagnostics

	defaultAggregation := "MEAN"
	// Check for time-based latency indicator
	if indicatorModel.TimeBasedLatencyIndicatorModel != nil {
		model := indicatorModel.TimeBasedLatencyIndicatorModel

		if model.Threshold.IsNull() || model.Threshold.IsUnknown() ||
			model.Aggregation.IsNull() || model.Aggregation.IsUnknown() {
			diags.AddError(
				"threshold and  aggregation are required for time_based_latency indicator",
				"threshold and  aggregation are required for time_based_latency indicator",
			)
			return restapi.SloIndicator{}, diags
		}

		threshold := model.Threshold.ValueFloat64()
		aggregation := model.Aggregation.ValueString()

		// Create time-based latency indicator
		return restapi.SloIndicator{
			Blueprint:   SloConfigAPIIndicatorBlueprintLatency,
			Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
			Threshold:   threshold,
			Aggregation: &aggregation,
		}, diags
	}

	// Check for event-based latency indicator
	if indicatorModel.EventBasedLatencyIndicatorModel != nil {
		model := indicatorModel.EventBasedLatencyIndicatorModel
		if model.Threshold.IsNull() || model.Threshold.IsUnknown() {
			diags.AddError(
				"threshold is required for event_based_latency indicator",
				"threshold is required for event_based_latency indicator",
			)
			return restapi.SloIndicator{}, diags
		}
		threshold := model.Threshold.ValueFloat64()

		// Create event-based latency indicator
		return restapi.SloIndicator{
			Blueprint:   SloConfigAPIIndicatorBlueprintLatency,
			Type:        SloConfigAPIIndicatorMeasurementTypeEventBased,
			Threshold:   threshold,
			Aggregation: &defaultAggregation,
		}, diags
	}

	// Check for time-based availability indicator
	if indicatorModel.TimeBasedAvailabilityIndicatorModel != nil {
		model := indicatorModel.TimeBasedAvailabilityIndicatorModel
		if model.Threshold.IsNull() || model.Threshold.IsUnknown() ||
			model.Aggregation.IsNull() || model.Aggregation.IsUnknown() {
			diags.AddError(
				"threshold and  aggregation are required for time_based_availability indicator",
				"threshold and  aggregation are required for time_based_availability indicator",
			)
			return restapi.SloIndicator{}, diags
		}
		threshold := model.Threshold.ValueFloat64()
		aggregation := model.Aggregation.ValueString()

		// Create time-based availability indicator
		return restapi.SloIndicator{
			Blueprint:   SloConfigAPIIndicatorBlueprintAvailability,
			Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
			Threshold:   threshold,
			Aggregation: &aggregation,
		}, diags
	}

	// Check for event-based availability indicator
	if indicatorModel.EventBasedAvailabilityIndicatorModel != nil {
		// Create event-based availability indicator
		return restapi.SloIndicator{
			Blueprint:   SloConfigAPIIndicatorBlueprintAvailability,
			Type:        SloConfigAPIIndicatorMeasurementTypeEventBased,
			Aggregation: &defaultAggregation,
		}, diags
	}

	// Check for traffic indicator
	if indicatorModel.TrafficIndicatorModel != nil {
		model := indicatorModel.TrafficIndicatorModel
		if model.Threshold.IsNull() || model.Threshold.IsUnknown() ||
			model.Operator.IsNull() || model.Operator.IsUnknown() {
			diags.AddError(
				"threshold and  operator are required for time_based_latency indicator",
				"threshold and  operator are required for time_based_latency indicator",
			)
			return restapi.SloIndicator{}, diags
		}
		trafficType := model.TrafficType.ValueString()
		threshold := model.Threshold.ValueFloat64()
		operator := model.Operator.ValueString()

		// Create traffic indicator
		return restapi.SloIndicator{
			Blueprint:   SloConfigAPIIndicatorBlueprintTraffic,
			Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
			TrafficType: &trafficType,
			Threshold:   threshold,
			Operator:    &operator,
			Aggregation: &defaultAggregation,
		}, diags
	}

	// Check for custom indicator
	if indicatorModel.CustomIndicatorModel != nil {
		model := indicatorModel.CustomIndicatorModel
		if model.GoodEventFilterExpression.IsNull() || model.GoodEventFilterExpression.IsUnknown() {
			diags.AddError(
				"good_event_filter_expression is required for custom indicator",
				"good_event_filter_expression is required for custom indicator",
			)
			return restapi.SloIndicator{}, diags
		}
		// Convert good event filter expressions to API model
		var goodEventFilter, badEventFilter *restapi.TagFilter

		goodEventFilter, tagDiags := mapTagFilterFromState(model.GoodEventFilterExpression, goodEventFilter)
		if tagDiags.HasError() {
			return restapi.SloIndicator{}, tagDiags
		}

		// Convert bad event filter expression to API model if set
		badEventFilter, tagBadDiags := mapTagFilterFromState(model.BadEventFilterExpression, badEventFilter)
		if tagBadDiags.HasError() {
			return restapi.SloIndicator{}, tagDiags
		}

		// Create custom indicator
		return restapi.SloIndicator{
			Type:                      SloConfigAPIIndicatorMeasurementTypeEventBased,
			Blueprint:                 SloConfigAPIIndicatorBlueprintCustom,
			GoodEventFilterExpression: goodEventFilter,
			BadEventFilterExpression:  badEventFilter,
			Aggregation:               &defaultAggregation,
		}, diags
	}

	diags.AddError(
		"Missing indicator configuration",
		"Exactly one indicator configuration is required",
	)
	return restapi.SloIndicator{}, diags
}

// Helper methods for mapping time window from plan
func (r *sloConfigResourceFramework) mapTimeWindowFromState(ctx context.Context, timeWindowModel TimeWindowModel) (restapi.SloTimeWindow, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check for rolling time window
	if timeWindowModel.RollingTimeWindowModel != nil {
		rollingModel := timeWindowModel.RollingTimeWindowModel

		if rollingModel.Duration.IsNull() || rollingModel.Duration.IsUnknown() ||
			rollingModel.DurationUnit.IsNull() || rollingModel.DurationUnit.IsUnknown() {

			diags.AddError(
				"duration and duration_unit are required for rolling time window",
				"duration and duration_unit are required for rolling time window",
			)
			return restapi.SloTimeWindow{}, diags
		}

		duration := int(rollingModel.Duration.ValueInt64())
		durationUnit := rollingModel.DurationUnit.ValueString()

		// Create rolling time window
		timeWindow := restapi.SloTimeWindow{
			Type:         SloConfigRollingTimeWindow,
			Duration:     duration,
			DurationUnit: durationUnit,
		}

		// Set timezone if present
		if !rollingModel.Timezone.IsNull() && !rollingModel.Timezone.IsUnknown() {
			timeWindow.Timezone = rollingModel.Timezone.ValueString()
		}

		return timeWindow, diags
	}

	// Check for fixed time window
	if timeWindowModel.FixedTimeWindowModel != nil {
		fixedModel := timeWindowModel.FixedTimeWindowModel

		if fixedModel.Duration.IsNull() || fixedModel.Duration.IsUnknown() ||
			fixedModel.DurationUnit.IsNull() || fixedModel.DurationUnit.IsUnknown() ||
			fixedModel.StartTimestamp.IsNull() || fixedModel.StartTimestamp.IsUnknown() {

			diags.AddError(
				"duration,duration_unit,start_timestamp are required for fixed time window",
				"duration,duration_unit,start_timestamp are required for fixed time window",
			)
			return restapi.SloTimeWindow{}, diags
		}
		duration := int(fixedModel.Duration.ValueInt64())
		durationUnit := fixedModel.DurationUnit.ValueString()
		startTime := fixedModel.StartTimestamp.ValueFloat64()

		// Create fixed time window
		timeWindow := restapi.SloTimeWindow{
			Type:         SloConfigFixedTimeWindow,
			Duration:     duration,
			DurationUnit: durationUnit,
			StartTime:    startTime,
		}

		// Set timezone if present
		if !fixedModel.Timezone.IsNull() && !fixedModel.Timezone.IsUnknown() {
			timeWindow.Timezone = fixedModel.Timezone.ValueString()
		}

		return timeWindow, diags
	}

	diags.AddError(
		"Missing time window configuration",
		"Exactly one time window configuration is required",
	)

	// Return an empty SloTimeWindow with error diagnostics
	return restapi.SloTimeWindow{}, diags
}

// Helper methods for mapping entity, indicator, and time window to state
func (r *sloConfigResourceFramework) mapEntityToState(ctx context.Context, apiObject *restapi.SloConfig) (EntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	entityModel := EntityModel{
		ApplicationEntityModel: nil,
		WebsiteEntityModel:     nil,
		SyntheticEntityModel:   nil,
	}

	// Create entity object based on type
	switch apiObject.Entity.Type {
	case SloConfigApplicationEntity:
		appModel, appDiags := r.mapApplicationEntityToState(ctx, apiObject.Entity)
		diags.Append(appDiags...)
		if !diags.HasError() {
			entityModel.ApplicationEntityModel = &appModel
		}
	case SloConfigWebsiteEntity:
		websiteModel, websiteDiags := r.mapWebsiteEntityToState(ctx, apiObject.Entity)
		diags.Append(websiteDiags...)
		if !diags.HasError() {
			entityModel.WebsiteEntityModel = &websiteModel
		}
	case SloConfigSyntheticEntity:
		syntheticModel, syntheticDiags := r.mapSyntheticEntityToState(ctx, apiObject.Entity)
		diags.Append(syntheticDiags...)
		if !diags.HasError() {
			entityModel.SyntheticEntityModel = &syntheticModel
		}
	default:
		diags.AddError(
			"Error mapping entity to state",
			fmt.Sprintf("Unsupported entity type: %s", apiObject.Entity.Type),
		)
	}

	return entityModel, diags
}

func (r *sloConfigResourceFramework) mapApplicationEntityToState(ctx context.Context, entity restapi.SloEntity) (ApplicationEntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create application entity object
	appEntityObj := ApplicationEntityModel{
		ApplicationID:    handleStringValue(entity.ApplicationID),
		BoundaryScope:    handleStringValue(entity.BoundaryScope),
		IncludeInternal:  handleBooleanValue(entity.IncludeInternal),
		IncludeSynthetic: handleBooleanValue(entity.IncludeSynthetic),
		ServiceID:        handleStringValue(entity.ServiceID),
		EndpointID:       handleStringValue(entity.EndpointID),
	}

	// Handle filter expression
	if entity.FilterExpression != nil {
		filterExpression, err := tagfilter.MapTagFilterToNormalizedString(entity.FilterExpression)
		if err != nil {
			diags.AddError(
				"Error normalizing filter expression",
				"Could not normalize filter expression: "+err.Error(),
			)
			return ApplicationEntityModel{}, diags
		}
		appEntityObj.FilterExpression = handleStringValue(filterExpression)

	}

	return appEntityObj, diags
}

func (r *sloConfigResourceFramework) mapWebsiteEntityToState(ctx context.Context, entity restapi.SloEntity) (WebsiteEntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create website entity object
	websiteEntityObj := WebsiteEntityModel{
		WebsiteID:  handleStringValue(entity.WebsiteId),
		BeaconType: handleStringValue(entity.BeaconType),
	}
	// Handle filter expression
	if entity.FilterExpression != nil {
		filterExpression, err := tagfilter.MapTagFilterToNormalizedString(entity.FilterExpression)
		if err != nil {
			diags.AddError(
				"Error normalizing filter expression",
				"Could not normalize filter expression: "+err.Error(),
			)
			return WebsiteEntityModel{}, diags
		}
		websiteEntityObj.FilterExpression = handleStringValue(filterExpression)

	}

	return websiteEntityObj, diags
}

func (r *sloConfigResourceFramework) mapSyntheticEntityToState(ctx context.Context, entity restapi.SloEntity) (SyntheticEntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert synthetic test IDs to types.String
	var testIDs []types.String
	for _, id := range entity.SyntheticTestIDs {
		if idStr, ok := id.(string); ok {
			testIDs = append(testIDs, types.StringValue(idStr))
		}
	}

	// Create synthetic entity object
	syntheticEntityObj := SyntheticEntityModel{
		SyntheticTestIDs: testIDs,
	}
	// Handle filter expression
	if entity.FilterExpression != nil {
		filterExpression, err := tagfilter.MapTagFilterToNormalizedString(entity.FilterExpression)
		if err != nil {
			diags.AddError(
				"Error normalizing filter expression",
				"Could not normalize filter expression: "+err.Error(),
			)
			return SyntheticEntityModel{}, diags
		}
		syntheticEntityObj.FilterExpression = handleStringValue(filterExpression)

	}

	return syntheticEntityObj, diags
}

func (r *sloConfigResourceFramework) mapIndicatorToState(ctx context.Context, apiObject *restapi.SloConfig) (IndicatorModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	indicator := apiObject.Indicator

	// Create indicator model
	indicatorModel := IndicatorModel{
		TimeBasedLatencyIndicatorModel:       nil,
		EventBasedLatencyIndicatorModel:      nil,
		TimeBasedAvailabilityIndicatorModel:  nil,
		EventBasedAvailabilityIndicatorModel: nil,
		TrafficIndicatorModel:                nil,
		CustomIndicatorModel:                 nil,
	}

	// Create indicator object based on type and blueprint
	switch {
	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeTimeBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintLatency:
		model := &TimeBasedLatencyIndicatorModel{
			Threshold:   types.Float64Value(indicator.Threshold),
			Aggregation: handleStringValue(indicator.Aggregation),
		}
		indicatorModel.TimeBasedLatencyIndicatorModel = model

	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeEventBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintLatency:
		model := &EventBasedLatencyIndicatorModel{
			Threshold: types.Float64Value(indicator.Threshold),
		}
		indicatorModel.EventBasedLatencyIndicatorModel = model

	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeTimeBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintAvailability:
		model := &TimeBasedAvailabilityIndicatorModel{
			Threshold:   types.Float64Value(indicator.Threshold),
			Aggregation: handleStringValue(indicator.Aggregation),
		}
		indicatorModel.TimeBasedAvailabilityIndicatorModel = model

	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeEventBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintAvailability:
		model := &EventBasedAvailabilityIndicatorModel{}
		indicatorModel.EventBasedAvailabilityIndicatorModel = model

	case indicator.Blueprint == SloConfigAPIIndicatorBlueprintTraffic:
		model := &TrafficIndicatorModel{
			TrafficType: handleStringValue(indicator.TrafficType),
			Threshold:   types.Float64Value(indicator.Threshold),
			Aggregation: handleStringValue(indicator.Aggregation),
		}
		indicatorModel.TrafficIndicatorModel = model

	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeEventBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintCustom:
		// Handle filter expressions
		model := &CustomIndicatorModel{}
		if indicator.GoodEventFilterExpression != nil {
			goodEventFilterExpression, err := tagfilter.MapTagFilterToNormalizedString(indicator.GoodEventFilterExpression)
			if err != nil {
				diags.AddError(
					"Error normalizing goodEventFilterExpression",
					"Could not normalize goodEventFilterExpression: "+err.Error(),
				)
				return IndicatorModel{}, diags
			}
			model.GoodEventFilterExpression = handleStringValue(goodEventFilterExpression)

		}

		if indicator.BadEventFilterExpression != nil {
			badEventFilterExpression, err := tagfilter.MapTagFilterToNormalizedString(indicator.BadEventFilterExpression)
			if err != nil {
				diags.AddError(
					"Error normalizing badEventFilterExpression",
					"Could not normalize badEventFilterExpression: "+err.Error(),
				)
				return IndicatorModel{}, diags
			}
			model.BadEventFilterExpression = handleStringValue(badEventFilterExpression)

		}

		indicatorModel.CustomIndicatorModel = model

	default:
		diags.AddError(
			"Error mapping indicator to state",
			fmt.Sprintf("Unsupported indicator type: %s, blueprint: %s", indicator.Type, indicator.Blueprint),
		)
	}

	return indicatorModel, diags
}

func (r *sloConfigResourceFramework) mapTimeWindowToState(ctx context.Context, apiObject *restapi.SloConfig) (TimeWindowModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	timeWindow := apiObject.TimeWindow

	// Create time window model
	timeWindowModel := TimeWindowModel{
		FixedTimeWindowModel:   nil,
		RollingTimeWindowModel: nil,
	}

	// Create time window object based on type
	switch timeWindow.Type {
	case SloConfigRollingTimeWindow:
		model := &RollingTimeWindowModel{
			Duration:     types.Int64Value(int64(timeWindow.Duration)),
			DurationUnit: types.StringValue(timeWindow.DurationUnit),
		}

		if timeWindow.Timezone != "" {
			model.Timezone = types.StringValue(timeWindow.Timezone)
		} else {
			model.Timezone = types.StringNull()
		}

		timeWindowModel.RollingTimeWindowModel = model

	case SloConfigFixedTimeWindow:
		model := &FixedTimeWindowModel{
			Duration:       types.Int64Value(int64(timeWindow.Duration)),
			DurationUnit:   types.StringValue(timeWindow.DurationUnit),
			StartTimestamp: types.Float64Value(timeWindow.StartTime),
		}

		if timeWindow.Timezone != "" {
			model.Timezone = types.StringValue(timeWindow.Timezone)
		} else {
			model.Timezone = types.StringNull()
		}

		timeWindowModel.FixedTimeWindowModel = model

	default:
		diags.AddError(
			"Error mapping time window to state",
			fmt.Sprintf("Unsupported time window type: %s", timeWindow.Type),
		)
	}

	return timeWindowModel, diags
}

// Made with Bob
