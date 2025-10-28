package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SloConfigModel represents the data model for the SLO configuration resource
type SloConfigModel struct {
	ID         types.String   `tfsdk:"id"`
	Name       types.String   `tfsdk:"name"`
	Target     types.Float64  `tfsdk:"target"`
	Tags       []types.String `tfsdk:"tags"`
	RbacTags   []RbacTagModel `tfsdk:"rbac_tags"`
	Entity     EntityModel    `tfsdk:"entity"`
	Indicator  IndicatorModel `tfsdk:"indicator"`
	TimeWindow types.Object   `tfsdk:"time_window"`
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
	InfraEntityModel       *InfraEntityModel       `tfsdk:"infra"`
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
type InfraEntityModel struct {
	Type      string `json:"type"`
	InfraType string `json:"infraType"`
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
							SloConfigApplicationEntity: schema.ListNestedBlock{
								Description: "Application entity of SLO",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										SloConfigFieldApplicationID: schema.StringAttribute{
											Required:    true,
											Description: "The application ID of the entity",
										},
										SloConfigFieldBoundaryScope: schema.StringAttribute{
											Required:    true,
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
							},
							SloConfigWebsiteEntity: schema.ListNestedBlock{
								Description: "Website entity of SLO",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										SloConfigFieldWebsiteID: schema.StringAttribute{
											Required:    true,
											Description: "The website ID of the entity",
										},
										SloConfigFieldFilterExpression: schema.StringAttribute{
											Optional:    true,
											Description: "Entity filter",
										},
										SloConfigFieldBeaconType: schema.StringAttribute{
											Required:    true,
											Description: "The beacon type for the entity configuration (pageLoad, resourceLoad, httpRequest, error, custom, pageChange)",
										},
									},
								},
							},
							SloConfigSyntheticEntity: schema.ListNestedBlock{
								Description: "Synthetic entity of SLO",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										SloConfigFieldSyntheticTestIDs: schema.ListAttribute{
											ElementType: types.StringType,
											Required:    true,
											Description: "The synthetics ID of the entity",
										},
										SloConfigFieldFilterExpression: schema.StringAttribute{
											Optional:    true,
											Description: "Entity filter",
										},
									},
								},
							},
							SloConfigInfraEntity: schema.ListNestedBlock{
								Description: "Synthetic entity of SLO",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"infra_type": schema.StringAttribute{
											Optional:    true,
											Description: "Entity filter",
										},
									},
								},
							},
						},
					},
					SloConfigFieldSloIndicator: schema.SingleNestedBlock{
						Description: "The indicator to use for the SLO configuration",
						Blocks: map[string]schema.Block{
							"time_based_latency": schema.ListNestedBlock{
								Description: "Time-based latency indicator",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"threshold": schema.Float64Attribute{
											Required:    true,
											Description: "The threshold for the metric configuration",
										},
										"aggregation": schema.StringAttribute{
											Required:    true,
											Description: "The aggregation type for the metric configuration",
										},
									},
								},
							},
							"event_based_latency": schema.ListNestedBlock{
								Description: "Event-based latency indicator",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"threshold": schema.Float64Attribute{
											Required:    true,
											Description: "The threshold for the metric configuration",
										},
									},
								},
							},
							"time_based_availability": schema.ListNestedBlock{
								Description: "Time-based availability indicator",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"threshold": schema.Float64Attribute{
											Required:    true,
											Description: "The threshold for the metric configuration",
										},
										"aggregation": schema.StringAttribute{
											Required:    true,
											Description: "The aggregation type for the metric configuration",
										},
									},
								},
							},
							"event_based_availability": schema.ListNestedBlock{
								Description:  "Event-based availability indicator",
								NestedObject: schema.NestedBlockObject{},
							},
							"traffic": schema.ListNestedBlock{
								Description: "Traffic indicator",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"traffic_type": schema.StringAttribute{
											Required:    true,
											Description: "The traffic type for the indicator",
										},
										"threshold": schema.Float64Attribute{
											Required:    true,
											Description: "The threshold for the metric configuration",
										},
										"aggregation": schema.StringAttribute{
											Required:    true,
											Description: "The aggregation type for the metric configuration",
										},
									},
								},
							},
							"custom": schema.ListNestedBlock{
								Description: "Custom indicator",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"good_event_filter_expression": schema.StringAttribute{
											Required:    true,
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
					},
					SloConfigFieldSloTimeWindow: schema.SingleNestedBlock{
						Description: "The time window to use for the SLO configuration",
						Blocks: map[string]schema.Block{
							"rolling": schema.ListNestedBlock{
								Description: "Rolling time window",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"duration": schema.Int64Attribute{
											Required:    true,
											Description: "The duration of the time window",
										},
										"duration_unit": schema.StringAttribute{
											Required:    true,
											Description: "The duration unit of the time window (day, week)",
										},
										"timezone": schema.StringAttribute{
											Optional:    true,
											Description: "The timezone for the SLO configuration",
										},
									},
								},
							},
							"fixed": schema.ListNestedBlock{
								Description: "Fixed time window",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"duration": schema.Int64Attribute{
											Required:    true,
											Description: "The duration of the time window",
										},
										"duration_unit": schema.StringAttribute{
											Required:    true,
											Description: "The duration unit of the time window (day, week)",
										},
										"timezone": schema.StringAttribute{
											Optional:    true,
											Description: "The timezone for the SLO configuration",
										},
										"start_timestamp": schema.Float64Attribute{
											Required:    true,
											Description: "Time window start time",
										},
									},
								},
							},
						},
					},
				},
			},
			SchemaVersion: 1,
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

func (r *sloConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
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

	// Get time window data
	var timeWindow interface{}
	// Check if time window is set (it's a types.Object, so we can use IsNull)
	if !model.TimeWindow.IsNull() && !model.TimeWindow.IsUnknown() {
		// Convert the time window object to a tfsdk.Config for processing
		var planData tfsdk.Config
		if plan != nil {
			diags.Append(plan.Get(ctx, &planData)...)
		} else if state != nil {
			diags.Append(state.Get(ctx, &planData)...)
		}

		timeWindowData, timeWindowDiags := r.mapTimeWindowFromPlan(ctx, planData)
		diags.Append(timeWindowDiags...)
		if !diags.HasError() {
			timeWindow = timeWindowData
		}
	}

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
		TimeWindow: timeWindow,
		RbacTags:   rbacTagsList,
	}

	// Generate ID if needed
	if sloConfig.ID == "" {
		sloConfig.ID = SloConfigFromTerraformIdPrefix + RandomID()
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
		if tagsList, ok := apiObject.Tags.([]interface{}); ok {
			var tags []types.String
			for _, tag := range tagsList {
				if tagStr, ok := tag.(string); ok {
					tags = append(tags, types.StringValue(tagStr))
				}
			}
			model.Tags = tags
		}
	}

	// Set RBAC tags if present
	if apiObject.RbacTags != nil {
		if rbacTagsList, ok := apiObject.RbacTags.([]interface{}); ok {
			var rbacTags []types.Object
			for _, tag := range rbacTagsList {
				if rbacTagMap, ok := tag.(map[string]interface{}); ok {
					displayName, _ := rbacTagMap["displayName"].(string)
					id, _ := rbacTagMap["id"].(string)

					rbacTagObj, err := types.ObjectValueFrom(
						ctx,
						map[string]attr.Type{
							"display_name": types.StringType,
							"id":           types.StringType,
						},
						map[string]attr.Value{
							"display_name": types.StringValue(displayName),
							"id":           types.StringValue(id),
						},
					)
					if err == nil {
						rbacTags = append(rbacTags, rbacTagObj)
					}
				}
			}
			model.RbacTags = rbacTags
		}
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
		applicationIdStr := applicationModel.ApplicationID.ValueString()
		serviceID := applicationModel.ServiceID.ValueString()
		endpointID := applicationModel.EndpointID.ValueString()
		boundaryScope := applicationModel.BoundaryScope.ValueString()
		includeInternal := applicationModel.IncludeInternal.ValueBool()
		includeSynthetic := applicationModel.IncludeSynthetic.ValueBool()

		appEntityObj := restapi.SloEntity{
			Type:             SloConfigApplicationEntity,
			ApplicationID:    &applicationIdStr,
			ServiceID:        &serviceID,
			EndpointID:       &endpointID,
			BoundaryScope:    &boundaryScope,
			IncludeInternal:  &includeInternal,
			IncludeSynthetic: &includeSynthetic,
		}

		// Convert filter expression to API model if set
		var tagFilter *restapi.TagFilter
		if !applicationModel.FilterExpression.IsNull() && !applicationModel.FilterExpression.IsUnknown() {
			parser := tagfilter.NewParser()
			expr, err := parser.Parse(applicationModel.FilterExpression.ValueString())
			if err != nil {
				diags.AddError(
					"Error parsing filter expression",
					fmt.Sprintf("Could not parse filter expression: %s", err),
				)
				return nil, diags
			}

			mapper := tagfilter.NewMapper()
			tagFilter = mapper.ToAPIModel(expr)
		}
		appEntityObj.FilterExpression = tagFilter
		return appEntityObj, diags
	}

	// Check for website entity
	if entityObj.WebsiteEntityModel != nil {
		websiteModel := entityObj.WebsiteEntityModel
		websiteIdStr := websiteModel.WebsiteID.ValueString()
		beaconTypeStr := websiteModel.BeaconType.ValueString()

		websiteEntityObj := restapi.SloEntity{
			Type:       SloConfigWebsiteEntity,
			WebsiteId:  &websiteIdStr,
			BeaconType: &beaconTypeStr,
		}

		// Convert filter expression to API model if set
		var tagFilter *restapi.TagFilter
		if !websiteModel.FilterExpression.IsNull() && !websiteModel.FilterExpression.IsUnknown() {
			parser := tagfilter.NewParser()
			expr, err := parser.Parse(websiteModel.FilterExpression.ValueString())
			if err != nil {
				diags.AddError(
					"Error parsing filter expression",
					fmt.Sprintf("Could not parse filter expression: %s", err),
				)
				return nil, diags
			}

			mapper := tagfilter.NewMapper()
			tagFilter = mapper.ToAPIModel(expr)
		}
		websiteEntityObj.FilterExpression = tagFilter
		return websiteEntityObj, diags
	}

	// Check for synthetic entity
	if entityObj.SyntheticEntityModel != nil {
		syntheticModel := entityObj.SyntheticEntityModel

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
		if !syntheticModel.FilterExpression.IsNull() && !syntheticModel.FilterExpression.IsUnknown() {
			parser := tagfilter.NewParser()
			expr, err := parser.Parse(syntheticModel.FilterExpression.ValueString())
			if err != nil {
				diags.AddError(
					"Error parsing filter expression",
					fmt.Sprintf("Could not parse filter expression: %s", err),
				)
				return nil, diags
			}

			mapper := tagfilter.NewMapper()
			tagFilter = mapper.ToAPIModel(expr)
		}
		syntheticEntityObj.FilterExpression = tagFilter
		return syntheticEntityObj, diags
	}

	// Check for infra entity
	if entityObj.InfraEntityModel != nil {
		infraModel := entityObj.InfraEntityModel

		infraEntityObj := restapi.SloEntity{
			Type:      SloConfigInfraEntity,
			InfraType: infraModel.InfraType,
		}

		return infraEntityObj, diags
	}

	diags.AddError(
		"Missing entity configuration",
		"Exactly one entity configuration is required",
	)
	return restapi.SloEntity{}, diags
}

// Helper methods for mapping indicator from plan
func (r *sloConfigResourceFramework) mapIndicatorFromState(ctx context.Context, indicatorModel IndicatorModel) (restapi.SloIndicator, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check which indicator type is set
	var indicatorBlock types.Object
	diags.Append(indicatorModel.GetAttribute(ctx, path.Root(SloConfigFieldSloIndicator), &indicatorBlock)...)

	if diags.HasError() || indicatorBlock.IsNull() || indicatorBlock.IsUnknown() {
		return nil, diags
	}

	// Check for time-based latency indicator
	var timeBasedLatency types.Object
	timeBasedLatencyPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("time_based_latency")
	diags.Append(indicatorModel.GetAttribute(ctx, timeBasedLatencyPath, &timeBasedLatency)...)

	if !timeBasedLatency.IsNull() && !timeBasedLatency.IsUnknown() {
		// Get time-based latency fields
		var threshold types.Float64
		var aggregation types.String

		timeBasedLatencyPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("time_based_latency").AtListIndex(0)

		diags.Append(indicatorModel.GetAttribute(ctx, timeBasedLatencyPath.AtName("threshold"), &threshold)...)
		diags.Append(indicatorModel.GetAttribute(ctx, timeBasedLatencyPath.AtName("aggregation"), &aggregation)...)

		if diags.HasError() {
			return nil, diags
		}

		// Create time-based latency indicator
		return restapi.SloTimeBasedLatencyIndicator{
			Blueprint:   SloConfigAPIIndicatorBlueprintLatency,
			Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
			Threshold:   threshold.ValueFloat64(),
			Aggregation: aggregation.ValueString(),
		}, diags
	}

	// Check for event-based latency indicator
	var eventBasedLatency types.Object
	eventBasedLatencyPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("event_based_latency")
	diags.Append(indicatorModel.GetAttribute(ctx, eventBasedLatencyPath, &eventBasedLatency)...)

	if !eventBasedLatency.IsNull() && !eventBasedLatency.IsUnknown() {
		// Get event-based latency fields
		var threshold types.Float64

		eventBasedLatencyPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("event_based_latency").AtListIndex(0)

		diags.Append(indicatorModel.GetAttribute(ctx, eventBasedLatencyPath.AtName("threshold"), &threshold)...)

		if diags.HasError() {
			return nil, diags
		}

		// Create event-based latency indicator
		return restapi.SloEventBasedLatencyIndicator{
			Blueprint: SloConfigAPIIndicatorBlueprintLatency,
			Type:      SloConfigAPIIndicatorMeasurementTypeEventBased,
			Threshold: threshold.ValueFloat64(),
		}, diags
	}

	// Check for time-based availability indicator
	var timeBasedAvailability types.Object
	timeBasedAvailabilityPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("time_based_availability")
	diags.Append(indicatorModel.GetAttribute(ctx, timeBasedAvailabilityPath, &timeBasedAvailability)...)

	if !timeBasedAvailability.IsNull() && !timeBasedAvailability.IsUnknown() {
		// Get time-based availability fields
		var threshold types.Float64
		var aggregation types.String

		timeBasedAvailabilityPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("time_based_availability").AtListIndex(0)

		diags.Append(indicatorModel.GetAttribute(ctx, timeBasedAvailabilityPath.AtName("threshold"), &threshold)...)
		diags.Append(indicatorModel.GetAttribute(ctx, timeBasedAvailabilityPath.AtName("aggregation"), &aggregation)...)

		if diags.HasError() {
			return nil, diags
		}

		// Create time-based availability indicator
		return restapi.SloTimeBasedAvailabilityIndicator{
			Blueprint:   SloConfigAPIIndicatorBlueprintAvailability,
			Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
			Threshold:   threshold.ValueFloat64(),
			Aggregation: aggregation.ValueString(),
		}, diags
	}

	// Check for event-based availability indicator
	var eventBasedAvailability types.Object
	eventBasedAvailabilityPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("event_based_availability")
	diags.Append(indicatorModel.GetAttribute(ctx, eventBasedAvailabilityPath, &eventBasedAvailability)...)

	if !eventBasedAvailability.IsNull() && !eventBasedAvailability.IsUnknown() {
		// Create event-based availability indicator
		return restapi.SloEventBasedAvailabilityIndicator{
			Blueprint: SloConfigAPIIndicatorBlueprintAvailability,
			Type:      SloConfigAPIIndicatorMeasurementTypeEventBased,
		}, diags
	}

	// Check for traffic indicator
	var traffic types.Object
	trafficPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("traffic")
	diags.Append(indicatorModel.GetAttribute(ctx, trafficPath, &traffic)...)

	if !traffic.IsNull() && !traffic.IsUnknown() {
		// Get traffic fields
		var trafficType types.String
		var threshold types.Float64
		var aggregation types.String

		trafficPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("traffic").AtListIndex(0)

		diags.Append(indicatorModel.GetAttribute(ctx, trafficPath.AtName("traffic_type"), &trafficType)...)
		diags.Append(indicatorModel.GetAttribute(ctx, trafficPath.AtName("threshold"), &threshold)...)
		diags.Append(indicatorModel.GetAttribute(ctx, trafficPath.AtName("aggregation"), &aggregation)...)

		if diags.HasError() {
			return nil, diags
		}

		// Create traffic indicator
		return restapi.SloTrafficIndicator{
			Blueprint:   SloConfigAPIIndicatorBlueprintTraffic,
			TrafficType: trafficType.ValueString(),
			Threshold:   threshold.ValueFloat64(),
			Aggregation: aggregation.ValueString(),
		}, diags
	}

	// Check for custom indicator
	var custom types.Object
	customPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("custom")
	diags.Append(indicatorModel.GetAttribute(ctx, customPath, &custom)...)

	if !custom.IsNull() && !custom.IsUnknown() {
		// Get custom fields
		var goodEventFilterExpression, badEventFilterExpression types.String

		customPath := path.Root(SloConfigFieldSloIndicator).AtListIndex(0).AtName("custom").AtListIndex(0)

		diags.Append(indicatorModel.GetAttribute(ctx, customPath.AtName("good_event_filter_expression"), &goodEventFilterExpression)...)
		diags.Append(indicatorModel.GetAttribute(ctx, customPath.AtName("bad_event_filter_expression"), &badEventFilterExpression)...)

		if diags.HasError() {
			return nil, diags
		}

		// Convert filter expressions to API model
		var goodEventFilter, badEventFilter *restapi.TagFilter

		if !goodEventFilterExpression.IsNull() && !goodEventFilterExpression.IsUnknown() {
			parser := tagfilter.NewParser()
			expr, err := parser.Parse(goodEventFilterExpression.ValueString())
			if err != nil {
				diags.AddError(
					"Error parsing good event filter expression",
					fmt.Sprintf("Could not parse filter expression: %s", err),
				)
				return nil, diags
			}

			mapper := tagfilter.NewMapper()
			goodEventFilter = mapper.ToAPIModel(expr)
		}

		if !badEventFilterExpression.IsNull() && !badEventFilterExpression.IsUnknown() {
			parser := tagfilter.NewParser()
			expr, err := parser.Parse(badEventFilterExpression.ValueString())
			if err != nil {
				diags.AddError(
					"Error parsing bad event filter expression",
					fmt.Sprintf("Could not parse filter expression: %s", err),
				)
				return nil, diags
			}

			mapper := tagfilter.NewMapper()
			badEventFilter = mapper.ToAPIModel(expr)
		}

		// Create custom indicator
		return restapi.SloCustomIndicator{
			Type:                      SloConfigAPIIndicatorMeasurementTypeEventBased,
			Blueprint:                 SloConfigAPIIndicatorBlueprintCustom,
			GoodEventFilterExpression: goodEventFilter,
			BadEventFilterExpression:  badEventFilter,
		}, diags
	}

	diags.AddError(
		"Missing indicator configuration",
		"Exactly one indicator configuration is required",
	)
	return nil, diags
}

// Helper methods for mapping time window from plan
func (r *sloConfigResourceFramework) mapTimeWindowFromPlan(ctx context.Context, planData tfsdk.Config) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check which time window type is set
	var timeWindowBlock types.Object
	diags.Append(planData.GetAttribute(ctx, path.Root(SloConfigFieldSloTimeWindow), &timeWindowBlock)...)

	if diags.HasError() || timeWindowBlock.IsNull() || timeWindowBlock.IsUnknown() {
		return nil, diags
	}

	// Check for rolling time window
	var rolling types.Object
	rollingPath := path.Root(SloConfigFieldSloTimeWindow).AtListIndex(0).AtName("rolling")
	diags.Append(planData.GetAttribute(ctx, rollingPath, &rolling)...)

	if !rolling.IsNull() && !rolling.IsUnknown() {
		// Get rolling time window fields
		var duration types.Int64
		var durationUnit, timezone types.String

		rollingPath := path.Root(SloConfigFieldSloTimeWindow).AtListIndex(0).AtName("rolling").AtListIndex(0)

		diags.Append(planData.GetAttribute(ctx, rollingPath.AtName("duration"), &duration)...)
		diags.Append(planData.GetAttribute(ctx, rollingPath.AtName("duration_unit"), &durationUnit)...)
		diags.Append(planData.GetAttribute(ctx, rollingPath.AtName("timezone"), &timezone)...)

		if diags.HasError() {
			return nil, diags
		}

		// Create rolling time window
		timeWindow := restapi.SloRollingTimeWindow{
			Type:         SloConfigRollingTimeWindow,
			Duration:     int(duration.ValueInt64()),
			DurationUnit: durationUnit.ValueString(),
		}

		if !timezone.IsNull() && !timezone.IsUnknown() {
			timeWindow.Timezone = timezone.ValueString()
		}

		return timeWindow, diags
	}

	// Check for fixed time window
	var fixed types.Object
	fixedPath := path.Root(SloConfigFieldSloTimeWindow).AtListIndex(0).AtName("fixed")
	diags.Append(planData.GetAttribute(ctx, fixedPath, &fixed)...)

	if !fixed.IsNull() && !fixed.IsUnknown() {
		// Get fixed time window fields
		var duration types.Int64
		var durationUnit, timezone types.String
		var startTimestamp types.Float64

		fixedPath := path.Root(SloConfigFieldSloTimeWindow).AtListIndex(0).AtName("fixed").AtListIndex(0)

		diags.Append(planData.GetAttribute(ctx, fixedPath.AtName("duration"), &duration)...)
		diags.Append(planData.GetAttribute(ctx, fixedPath.AtName("duration_unit"), &durationUnit)...)
		diags.Append(planData.GetAttribute(ctx, fixedPath.AtName("timezone"), &timezone)...)
		diags.Append(planData.GetAttribute(ctx, fixedPath.AtName("start_timestamp"), &startTimestamp)...)

		if diags.HasError() {
			return nil, diags
		}

		// Create fixed time window
		timeWindow := restapi.SloFixedTimeWindow{
			Type:         SloConfigFixedTimeWindow,
			Duration:     int(duration.ValueInt64()),
			DurationUnit: durationUnit.ValueString(),
			StartTime:    startTimestamp.ValueFloat64(),
		}

		if !timezone.IsNull() && !timezone.IsUnknown() {
			timeWindow.Timezone = timezone.ValueString()
		}

		return timeWindow, diags
	}

	diags.AddError(
		"Missing time window configuration",
		"Exactly one time window configuration is required",
	)
	return nil, diags
}

// Helper methods for mapping entity, indicator, and time window to state
func (r *sloConfigResourceFramework) mapEntityToState(ctx context.Context, apiObject *restapi.SloConfig) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check entity type
	entityMap, ok := apiObject.Entity.(map[string]interface{})
	if !ok {
		diags.AddError(
			"Error mapping entity to state",
			fmt.Sprintf("Expected map[string]interface{}, got: %T", apiObject.Entity),
		)
		return types.ObjectNull(nil), diags
	}

	entityType, ok := entityMap["type"].(string)
	if !ok {
		diags.AddError(
			"Error mapping entity to state",
			"Entity type is missing or not a string",
		)
		return types.ObjectNull(nil), diags
	}

	// Create entity object based on type
	switch entityType {
	case SloConfigApplicationEntity:
		return r.mapApplicationEntityToState(ctx, entityMap)
	case SloConfigWebsiteEntity:
		return r.mapWebsiteEntityToState(ctx, entityMap)
	case SloConfigSyntheticEntity:
		return r.mapSyntheticEntityToState(ctx, entityMap)
	case SloConfigInfraEntity:
		return r.mapInfraEntityToState(ctx, entityMap)
	default:
		diags.AddError(
			"Error mapping entity to state",
			fmt.Sprintf("Unsupported entity type: %s", entityType),
		)
		return types.ObjectNull(nil), diags
	}
}

func (r *sloConfigResourceFramework) mapApplicationEntityToState(ctx context.Context, entityMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract application entity fields
	applicationID, _ := entityMap["applicationId"].(string)
	serviceID, _ := entityMap["serviceId"].(string)
	endpointID, _ := entityMap["endpointId"].(string)
	boundaryScope, _ := entityMap["boundaryScope"].(string)

	var includeInternal, includeSynthetic bool
	if val, ok := entityMap["includeInternal"]; ok && val != nil {
		includeInternal, _ = val.(bool)
	}
	if val, ok := entityMap["includeSynthetic"]; ok && val != nil {
		includeSynthetic, _ = val.(bool)
	}

	// Handle filter expression
	var filterExpression string
	if tagFilter, ok := entityMap["tagFilterExpression"].(*restapi.TagFilter); ok && tagFilter != nil {
		mapper := tagfilter.NewMapper()
		expr, err := mapper.FromAPIModel(tagFilter)
		if err == nil && expr != nil {
			// Convert the expression to string format
			filterExpression = fmt.Sprintf("%v", expr)
		}
	}

	// Create application entity object
	appEntityObj := map[string]interface{}{
		SloConfigFieldApplicationID:    types.StringValue(applicationID),
		SloConfigFieldBoundaryScope:    types.StringValue(boundaryScope),
		SloConfigFieldIncludeInternal:  types.BoolValue(includeInternal),
		SloConfigFieldIncludeSynthetic: types.BoolValue(includeSynthetic),
	}

	if serviceID != "" {
		appEntityObj[SloConfigFieldServiceID] = types.StringValue(serviceID)
	} else {
		appEntityObj[SloConfigFieldServiceID] = types.StringNull()
	}

	if endpointID != "" {
		appEntityObj[SloConfigFieldEndpointID] = types.StringValue(endpointID)
	} else {
		appEntityObj[SloConfigFieldEndpointID] = types.StringNull()
	}

	if filterExpression != "" {
		appEntityObj[SloConfigFieldFilterExpression] = types.StringValue(filterExpression)
	} else {
		appEntityObj[SloConfigFieldFilterExpression] = types.StringNull()
	}

	// Create entity object
	entityObj := map[string]interface{}{
		SloConfigApplicationEntity: []interface{}{appEntityObj},
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		SloConfigApplicationEntity: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					SloConfigFieldApplicationID:    types.StringType,
					SloConfigFieldBoundaryScope:    types.StringType,
					SloConfigFieldIncludeInternal:  types.BoolType,
					SloConfigFieldIncludeSynthetic: types.BoolType,
					SloConfigFieldServiceID:        types.StringType,
					SloConfigFieldEndpointID:       types.StringType,
					SloConfigFieldFilterExpression: types.StringType,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, entityObj)
	if err != nil {
		diags.AddError(
			"Error creating application entity object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapWebsiteEntityToState(ctx context.Context, entityMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract website entity fields
	websiteID, _ := entityMap["websiteId"].(string)
	beaconType, _ := entityMap["beaconType"].(string)

	// Handle filter expression
	var filterExpression string
	if tagFilter, ok := entityMap["tagFilterExpression"].(*restapi.TagFilter); ok && tagFilter != nil {
		mapper := tagfilter.NewMapper()
		expr, err := mapper.FromAPIModel(tagFilter)
		if err == nil && expr != nil {
			// Convert the expression to string format
			filterExpression = fmt.Sprintf("%v", expr)
		}
	}

	// Create website entity object
	websiteEntityObj := map[string]interface{}{
		SloConfigFieldWebsiteID:  types.StringValue(websiteID),
		SloConfigFieldBeaconType: types.StringValue(beaconType),
	}

	if filterExpression != "" {
		websiteEntityObj[SloConfigFieldFilterExpression] = types.StringValue(filterExpression)
	} else {
		websiteEntityObj[SloConfigFieldFilterExpression] = types.StringNull()
	}

	// Create entity object
	entityObj := map[string]interface{}{
		SloConfigWebsiteEntity: []interface{}{websiteEntityObj},
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		SloConfigWebsiteEntity: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					SloConfigFieldWebsiteID:        types.StringType,
					SloConfigFieldBeaconType:       types.StringType,
					SloConfigFieldFilterExpression: types.StringType,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, entityObj)
	if err != nil {
		diags.AddError(
			"Error creating website entity object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapSyntheticEntityToState(ctx context.Context, entityMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract synthetic entity fields
	syntheticTestIDs, _ := entityMap["syntheticTestIds"].([]interface{})

	// Handle filter expression
	var filterExpression string
	if tagFilter, ok := entityMap["tagFilterExpression"].(*restapi.TagFilter); ok && tagFilter != nil {
		mapper := tagfilter.NewMapper()
		expr, err := mapper.FromAPIModel(tagFilter)
		if err == nil && expr != nil {
			// Convert the expression to string format
			filterExpression = fmt.Sprintf("%v", expr)
		}
	}

	// Convert synthetic test IDs to types.String
	var testIDs []types.String
	for _, id := range syntheticTestIDs {
		if idStr, ok := id.(string); ok {
			testIDs = append(testIDs, types.StringValue(idStr))
		}
	}

	// Create synthetic entity object
	syntheticEntityObj := map[string]interface{}{
		SloConfigFieldSyntheticTestIDs: testIDs,
	}

	if filterExpression != "" {
		syntheticEntityObj[SloConfigFieldFilterExpression] = types.StringValue(filterExpression)
	} else {
		syntheticEntityObj[SloConfigFieldFilterExpression] = types.StringNull()
	}

	// Create entity object
	entityObj := map[string]interface{}{
		SloConfigSyntheticEntity: []interface{}{syntheticEntityObj},
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		SloConfigSyntheticEntity: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					SloConfigFieldSyntheticTestIDs: types.ListType{
						ElemType: types.StringType,
					},
					SloConfigFieldFilterExpression: types.StringType,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, entityObj)
	if err != nil {
		diags.AddError(
			"Error creating synthetic entity object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapInfraEntityToState(ctx context.Context, entityMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract infra entity fields
	infraType, _ := entityMap["infraType"].(string)

	// Create infra entity object
	infraEntityObj := map[string]interface{}{
		"infra_type": types.StringValue(infraType),
	}

	// Create entity object
	entityObj := map[string]interface{}{
		SloConfigInfraEntity: []interface{}{infraEntityObj},
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		SloConfigInfraEntity: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"infra_type": types.StringType,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, entityObj)
	if err != nil {
		diags.AddError(
			"Error creating infra entity object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapIndicatorToState(ctx context.Context, apiObject *restapi.SloConfig) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check indicator type
	indicatorMap, ok := apiObject.Indicator.(map[string]interface{})
	if !ok {
		diags.AddError(
			"Error mapping indicator to state",
			fmt.Sprintf("Expected map[string]interface{}, got: %T", apiObject.Indicator),
		)
		return types.ObjectNull(nil), diags
	}

	blueprint, _ := indicatorMap["blueprint"].(string)
	indicatorType, _ := indicatorMap["type"].(string)

	// Create indicator object based on type and blueprint
	switch {
	case indicatorType == SloConfigAPIIndicatorMeasurementTypeTimeBased && blueprint == SloConfigAPIIndicatorBlueprintLatency:
		return r.mapTimeBasedLatencyIndicatorToState(ctx, indicatorMap)
	case indicatorType == SloConfigAPIIndicatorMeasurementTypeEventBased && blueprint == SloConfigAPIIndicatorBlueprintLatency:
		return r.mapEventBasedLatencyIndicatorToState(ctx, indicatorMap)
	case indicatorType == SloConfigAPIIndicatorMeasurementTypeTimeBased && blueprint == SloConfigAPIIndicatorBlueprintAvailability:
		return r.mapTimeBasedAvailabilityIndicatorToState(ctx, indicatorMap)
	case indicatorType == SloConfigAPIIndicatorMeasurementTypeEventBased && blueprint == SloConfigAPIIndicatorBlueprintAvailability:
		return r.mapEventBasedAvailabilityIndicatorToState(ctx, indicatorMap)
	case blueprint == SloConfigAPIIndicatorBlueprintTraffic:
		return r.mapTrafficIndicatorToState(ctx, indicatorMap)
	case indicatorType == SloConfigAPIIndicatorMeasurementTypeEventBased && blueprint == SloConfigAPIIndicatorBlueprintCustom:
		return r.mapCustomIndicatorToState(ctx, indicatorMap)
	default:
		diags.AddError(
			"Error mapping indicator to state",
			fmt.Sprintf("Unsupported indicator type: %s, blueprint: %s", indicatorType, blueprint),
		)
		return types.ObjectNull(nil), diags
	}
}

func (r *sloConfigResourceFramework) mapTimeBasedLatencyIndicatorToState(ctx context.Context, indicatorMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract time-based latency fields
	threshold, _ := indicatorMap["threshold"].(float64)
	aggregation, _ := indicatorMap["aggregation"].(string)

	// Create time-based latency indicator object
	indicatorObj := map[string]interface{}{
		"time_based_latency": []interface{}{
			map[string]interface{}{
				"threshold":   types.Float64Value(threshold),
				"aggregation": types.StringValue(aggregation),
			},
		},
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		"time_based_latency": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"threshold":   types.Float64Type,
					"aggregation": types.StringType,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, indicatorObj)
	if err != nil {
		diags.AddError(
			"Error creating time-based latency indicator object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapEventBasedLatencyIndicatorToState(ctx context.Context, indicatorMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract event-based latency fields
	threshold, _ := indicatorMap["threshold"].(float64)

	// Create event-based latency indicator object
	indicatorObj := map[string]interface{}{
		"event_based_latency": []interface{}{
			map[string]interface{}{
				"threshold": types.Float64Value(threshold),
			},
		},
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		"event_based_latency": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"threshold": types.Float64Type,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, indicatorObj)
	if err != nil {
		diags.AddError(
			"Error creating event-based latency indicator object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapTimeBasedAvailabilityIndicatorToState(ctx context.Context, indicatorMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract time-based availability fields
	threshold, _ := indicatorMap["threshold"].(float64)
	aggregation, _ := indicatorMap["aggregation"].(string)

	// Create time-based availability indicator object
	indicatorObj := map[string]interface{}{
		"time_based_availability": []interface{}{
			map[string]interface{}{
				"threshold":   types.Float64Value(threshold),
				"aggregation": types.StringValue(aggregation),
			},
		},
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		"time_based_availability": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"threshold":   types.Float64Type,
					"aggregation": types.StringType,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, indicatorObj)
	if err != nil {
		diags.AddError(
			"Error creating time-based availability indicator object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapEventBasedAvailabilityIndicatorToState(ctx context.Context, indicatorMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create event-based availability indicator object
	indicatorObj := map[string]interface{}{
		"event_based_availability": []interface{}{
			map[string]interface{}{},
		},
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		"event_based_availability": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, indicatorObj)
	if err != nil {
		diags.AddError(
			"Error creating event-based availability indicator object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapTrafficIndicatorToState(ctx context.Context, indicatorMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract traffic fields
	trafficType, _ := indicatorMap["trafficType"].(string)
	threshold, _ := indicatorMap["threshold"].(float64)
	aggregation, _ := indicatorMap["aggregation"].(string)

	// Create traffic indicator object
	indicatorObj := map[string]interface{}{
		"traffic": []interface{}{
			map[string]interface{}{
				"traffic_type": types.StringValue(trafficType),
				"threshold":    types.Float64Value(threshold),
				"aggregation":  types.StringValue(aggregation),
			},
		},
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		"traffic": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"traffic_type": types.StringType,
					"threshold":    types.Float64Type,
					"aggregation":  types.StringType,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, indicatorObj)
	if err != nil {
		diags.AddError(
			"Error creating traffic indicator object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapCustomIndicatorToState(ctx context.Context, indicatorMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Handle filter expressions
	var goodEventFilterExpression, badEventFilterExpression string

	if goodEventFilter, ok := indicatorMap["goodEventsFilter"].(*restapi.TagFilter); ok && goodEventFilter != nil {
		mapper := tagfilter.NewMapper()
		expr, err := mapper.FromAPIModel(goodEventFilter)
		if err == nil && expr != nil {
			// Convert the expression to string format
			goodEventFilterExpression = fmt.Sprintf("%v", expr)
		}
	}

	if badEventFilter, ok := indicatorMap["badEventsFilter"].(*restapi.TagFilter); ok && badEventFilter != nil {
		mapper := tagfilter.NewMapper()
		expr, err := mapper.FromAPIModel(badEventFilter)
		if err == nil && expr != nil {
			// Convert the expression to string format
			badEventFilterExpression = fmt.Sprintf("%v", expr)
		}
	}

	// Create custom indicator object
	customObj := map[string]interface{}{
		"good_event_filter_expression": types.StringValue(goodEventFilterExpression),
	}

	if badEventFilterExpression != "" {
		customObj["bad_event_filter_expression"] = types.StringValue(badEventFilterExpression)
	} else {
		customObj["bad_event_filter_expression"] = types.StringNull()
	}

	// Create indicator object
	indicatorObj := map[string]interface{}{
		"custom": []interface{}{customObj},
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		"custom": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"good_event_filter_expression": types.StringType,
					"bad_event_filter_expression":  types.StringType,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, indicatorObj)
	if err != nil {
		diags.AddError(
			"Error creating custom indicator object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapTimeWindowToState(ctx context.Context, apiObject *restapi.SloConfig) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check time window type
	timeWindowMap, ok := apiObject.TimeWindow.(map[string]interface{})
	if !ok {
		diags.AddError(
			"Error mapping time window to state",
			fmt.Sprintf("Expected map[string]interface{}, got: %T", apiObject.TimeWindow),
		)
		return types.ObjectNull(nil), diags
	}

	timeWindowType, _ := timeWindowMap["type"].(string)

	// Create time window object based on type
	switch timeWindowType {
	case SloConfigRollingTimeWindow:
		return r.mapRollingTimeWindowToState(ctx, timeWindowMap)
	case SloConfigFixedTimeWindow:
		return r.mapFixedTimeWindowToState(ctx, timeWindowMap)
	default:
		diags.AddError(
			"Error mapping time window to state",
			fmt.Sprintf("Unsupported time window type: %s", timeWindowType),
		)
		return types.ObjectNull(nil), diags
	}
}

func (r *sloConfigResourceFramework) mapRollingTimeWindowToState(ctx context.Context, timeWindowMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract rolling time window fields
	duration, _ := timeWindowMap["duration"].(float64)
	durationUnit, _ := timeWindowMap["durationUnit"].(string)
	timezone, _ := timeWindowMap["timezone"].(string)

	// Create rolling time window object
	timeWindowObj := map[string]interface{}{
		"rolling": []interface{}{
			map[string]interface{}{
				"duration":      types.Int64Value(int64(duration)),
				"duration_unit": types.StringValue(durationUnit),
			},
		},
	}

	if timezone != "" {
		timeWindowObj["rolling"].([]interface{})[0].(map[string]interface{})["timezone"] = types.StringValue(timezone)
	} else {
		timeWindowObj["rolling"].([]interface{})[0].(map[string]interface{})["timezone"] = types.StringNull()
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		"rolling": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"duration":      types.Int64Type,
					"duration_unit": types.StringType,
					"timezone":      types.StringType,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, timeWindowObj)
	if err != nil {
		diags.AddError(
			"Error creating rolling time window object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

func (r *sloConfigResourceFramework) mapFixedTimeWindowToState(ctx context.Context, timeWindowMap map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract fixed time window fields
	duration, _ := timeWindowMap["duration"].(float64)
	durationUnit, _ := timeWindowMap["durationUnit"].(string)
	timezone, _ := timeWindowMap["timezone"].(string)
	startTimestamp, _ := timeWindowMap["startTimestamp"].(float64)

	// Create fixed time window object
	timeWindowObj := map[string]interface{}{
		"fixed": []interface{}{
			map[string]interface{}{
				"duration":        types.Int64Value(int64(duration)),
				"duration_unit":   types.StringValue(durationUnit),
				"start_timestamp": types.Float64Value(startTimestamp),
			},
		},
	}

	if timezone != "" {
		timeWindowObj["fixed"].([]interface{})[0].(map[string]interface{})["timezone"] = types.StringValue(timezone)
	} else {
		timeWindowObj["fixed"].([]interface{})[0].(map[string]interface{})["timezone"] = types.StringNull()
	}

	// Convert to types.Object
	objType := map[string]attr.Type{
		"fixed": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"duration":        types.Int64Type,
					"duration_unit":   types.StringType,
					"timezone":        types.StringType,
					"start_timestamp": types.Float64Type,
				},
			},
		},
	}

	obj, err := types.ObjectValueFrom(ctx, objType, timeWindowObj)
	if err != nil {
		diags.AddError(
			"Error creating fixed time window object",
			fmt.Sprintf("Could not create object: %s", err),
		)
		return types.ObjectNull(nil), diags
	}

	return obj, diags
}

// Made with Bob
