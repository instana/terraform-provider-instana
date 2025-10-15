package tf_framework

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &sloConfigResource{}
	_ resource.ResourceWithConfigure   = &sloConfigResource{}
	_ resource.ResourceWithImportState = &sloConfigResource{}
)

// NewSloConfigResource is a helper function to simplify the provider implementation
func NewSloConfigResource() resource.Resource {
	return &sloConfigResource{}
}

// sloConfigResource is the resource implementation
type sloConfigResource struct {
	client restapi.InstanaAPI
}

// Configure adds the provider configured client to the resource
func (r *sloConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(restapi.InstanaAPI)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected restapi.InstanaAPI, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata returns the resource type name
func (r *sloConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_slo_config"
}

// Schema defines the schema for the resource
func (r *sloConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to configure an SLO configuration in Instana.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the SLO configuration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the SLO configuration.",
				Required:    true,
			},
			"target": schema.Float64Attribute{
				Description: "The target value for the SLO.",
				Required:    true,
			},
			"tags": schema.ListAttribute{
				Description: "The tags for the SLO configuration.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"entity": schema.ListNestedBlock{
				Description: "The entity configuration for the SLO.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"application": schema.ListNestedBlock{
							Description: "Application entity configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"application_id": schema.StringAttribute{
										Description: "The ID of the application.",
										Required:    true,
									},
									"service_id": schema.StringAttribute{
										Description: "The ID of the service.",
										Optional:    true,
									},
									"endpoint_id": schema.StringAttribute{
										Description: "The ID of the endpoint.",
										Optional:    true,
									},
									"boundary_scope": schema.StringAttribute{
										Description: "The boundary scope for the entity (ALL, INBOUND).",
										Required:    true,
									},
									"include_synthetic": schema.BoolAttribute{
										Description: "Flag to indicate whether synthetic calls are included.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"include_internal": schema.BoolAttribute{
										Description: "Flag to indicate whether internal calls are included.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
									"filter_expression": schema.StringAttribute{
										Description: "The filter expression for the entity.",
										Optional:    true,
									},
								},
							},
						},
						"website": schema.ListNestedBlock{
							Description: "Website entity configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"website_id": schema.StringAttribute{
										Description: "The ID of the website.",
										Required:    true,
									},
									"beacon_type": schema.StringAttribute{
										Description: "The beacon type for the entity (pageLoad, resourceLoad, httpRequest, error, custom, pageChange).",
										Required:    true,
									},
									"filter_expression": schema.StringAttribute{
										Description: "The filter expression for the entity.",
										Optional:    true,
									},
								},
							},
						},
						"synthetic": schema.ListNestedBlock{
							Description: "Synthetic entity configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"synthetic_test_ids": schema.ListAttribute{
										Description: "The IDs of the synthetic tests.",
										Required:    true,
										ElementType: types.StringType,
									},
									"filter_expression": schema.StringAttribute{
										Description: "The filter expression for the entity.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
			"indicator": schema.ListNestedBlock{
				Description: "The indicator configuration for the SLO.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"time_based_latency": schema.ListNestedBlock{
							Description: "Time-based latency indicator configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"threshold": schema.Float64Attribute{
										Description: "The threshold for the indicator.",
										Required:    true,
									},
									"aggregation": schema.StringAttribute{
										Description: "The aggregation type for the indicator (SUM, MEAN, MAX, MIN, P25, P50, P75, P90, P95, P98, P99, P99_9, P99_99, DISTRIBUTION, DISTINCT_COUNT, SUM_POSITIVE, PER_SECOND).",
										Required:    true,
									},
								},
							},
						},
						"event_based_latency": schema.ListNestedBlock{
							Description: "Event-based latency indicator configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"threshold": schema.Float64Attribute{
										Description: "The threshold for the indicator.",
										Required:    true,
									},
								},
							},
						},
						"time_based_availability": schema.ListNestedBlock{
							Description: "Time-based availability indicator configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"threshold": schema.Float64Attribute{
										Description: "The threshold for the indicator.",
										Required:    true,
									},
									"aggregation": schema.StringAttribute{
										Description: "The aggregation type for the indicator (SUM, MEAN, MAX, MIN, P25, P50, P75, P90, P95, P98, P99, P99_9, P99_99, DISTRIBUTION, DISTINCT_COUNT, SUM_POSITIVE, PER_SECOND).",
										Required:    true,
									},
								},
							},
						},
						"event_based_availability": schema.ListNestedBlock{
							Description: "Event-based availability indicator configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{},
							},
						},
						"traffic": schema.ListNestedBlock{
							Description: "Traffic indicator configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"traffic_type": schema.StringAttribute{
										Description: "The traffic type for the indicator (all, erroneous).",
										Required:    true,
									},
									"threshold": schema.Float64Attribute{
										Description: "The threshold for the indicator.",
										Required:    true,
									},
									"aggregation": schema.StringAttribute{
										Description: "The aggregation type for the indicator (SUM, MEAN, MAX, MIN, P25, P50, P75, P90, P95, P98, P99, P99_9, P99_99, DISTRIBUTION, DISTINCT_COUNT, SUM_POSITIVE, PER_SECOND).",
										Required:    true,
									},
								},
							},
						},
						"custom": schema.ListNestedBlock{
							Description: "Custom indicator configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"good_event_filter_expression": schema.StringAttribute{
										Description: "The filter expression for good events.",
										Required:    true,
									},
									"bad_event_filter_expression": schema.StringAttribute{
										Description: "The filter expression for bad events.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
			"time_window": schema.ListNestedBlock{
				Description: "The time window configuration for the SLO.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"rolling": schema.ListNestedBlock{
							Description: "Rolling time window configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Description: "The duration for the time window.",
										Required:    true,
									},
									"duration_unit": schema.StringAttribute{
										Description: "The duration unit for the time window (day, week).",
										Required:    true,
									},
									"timezone": schema.StringAttribute{
										Description: "The timezone for the time window.",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("UTC"),
									},
								},
							},
						},
						"fixed": schema.ListNestedBlock{
							Description: "Fixed time window configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Description: "The duration for the time window.",
										Required:    true,
									},
									"duration_unit": schema.StringAttribute{
										Description: "The duration unit for the time window (day, week).",
										Required:    true,
									},
									"timezone": schema.StringAttribute{
										Description: "The timezone for the time window.",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("UTC"),
									},
									"start_timestamp": schema.Float64Attribute{
										Description: "The start timestamp for the time window.",
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state
func (r *sloConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan SloConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert from Terraform model to API object
	sloConfig, err := r.mapModelToAPIObject(ctx, plan, resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SLO configuration",
			fmt.Sprintf("Could not map SLO configuration to API object: %s", err),
		)
		return
	}

	// Create new SLO configuration
	sloConfigCreated, err := r.client.SloConfigs().Create(sloConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SLO configuration",
			fmt.Sprintf("Could not create SLO configuration: %s", err),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	err = r.mapAPIObjectToModel(ctx, sloConfigCreated, &plan, resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SLO configuration",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data
func (r *sloConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state SloConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed SLO configuration value from Instana
	sloConfig, err := r.client.SloConfigs().GetOne(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SLO configuration",
			fmt.Sprintf("Could not read SLO configuration ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Map response body to model
	err = r.mapAPIObjectToModel(ctx, sloConfig, &state, resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SLO configuration",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success
func (r *sloConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan SloConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state SloConfigModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve ID from state
	plan.ID = state.ID

	// Convert from Terraform model to API object
	sloConfig, err := r.mapModelToAPIObject(ctx, plan, resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating SLO configuration",
			fmt.Sprintf("Could not map SLO configuration to API object: %s", err),
		)
		return
	}

	// Update SLO configuration
	sloConfigUpdated, err := r.client.SloConfigs().Update(sloConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating SLO configuration",
			fmt.Sprintf("Could not update SLO configuration: %s", err),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	err = r.mapAPIObjectToModel(ctx, sloConfigUpdated, &plan, resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating SLO configuration",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success
func (r *sloConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state SloConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing SLO configuration
	err := r.client.SloConfigs().DeleteByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SLO configuration",
			fmt.Sprintf("Could not delete SLO configuration ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

// ImportState imports the resource into Terraform state
func (r *sloConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapModelToAPIObject maps the Terraform model to the API object
func (r *sloConfigResource) mapModelToAPIObject(ctx context.Context, model SloConfigModel, diags diag.Diagnostics) (*restapi.SloConfig, error) {
	sloConfig := &restapi.SloConfig{
		ID:     model.ID.ValueString(),
		Name:   model.Name.ValueString(),
		Target: model.Target.ValueFloat64(),
	}

	// Map tags
	if !model.Tags.IsNull() && !model.Tags.IsUnknown() {
		var tags []string
		diags.Append(model.Tags.ElementsAs(ctx, &tags, false)...)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to map tags: %v", diags)
		}
		sloConfig.Tags = tags
	}

	// Map entity
	entity, err := r.mapEntityModelToAPIObject(ctx, model, diags)
	if err != nil {
		return nil, fmt.Errorf("failed to map entity: %v", err)
	}
	sloConfig.Entity = entity

	// Map indicator
	indicator, err := r.mapIndicatorModelToAPIObject(ctx, model, diags)
	if err != nil {
		return nil, fmt.Errorf("failed to map indicator: %v", err)
	}
	sloConfig.Indicator = indicator

	// Map time window
	timeWindow, err := r.mapTimeWindowModelToAPIObject(ctx, model, diags)
	if err != nil {
		return nil, fmt.Errorf("failed to map time window: %v", err)
	}
	sloConfig.TimeWindow = timeWindow

	return sloConfig, nil
}

// mapEntityModelToAPIObject maps the entity model to the API object
func (r *sloConfigResource) mapEntityModelToAPIObject(ctx context.Context, model SloConfigModel, diags diag.Diagnostics) (interface{}, error) {
	var entityModels []SloEntityModel
	diags.Append(model.Entity.ElementsAs(ctx, &entityModels, false)...)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to map entity: %v", diags)
	}

	if len(entityModels) != 1 {
		return nil, fmt.Errorf("exactly one entity configuration is required")
	}

	entityModel := entityModels[0]

	// Check which entity type is set
	if entityModel.Application != nil {
		// Map application entity
		appEntity := restapi.SloApplicationEntity{
			Type:             "application",
			ApplicationID:    getStringPointer(entityModel.Application.ApplicationID),
			ServiceID:        getStringPointer(entityModel.Application.ServiceID),
			EndpointID:       getStringPointer(entityModel.Application.EndpointID),
			BoundaryScope:    getStringPointer(entityModel.Application.BoundaryScope),
			IncludeSynthetic: getBoolPointer(entityModel.Application.IncludeSynthetic),
			IncludeInternal:  getBoolPointer(entityModel.Application.IncludeInternal),
		}

		// Map filter expression if set
		if !entityModel.Application.FilterExpression.IsNull() && !entityModel.Application.FilterExpression.IsUnknown() {
			filterExpr, err := r.mapTagFilterStringToAPIModel(entityModel.Application.FilterExpression.ValueString())
			if err != nil {
				return nil, fmt.Errorf("failed to map application entity filter expression: %v", err)
			}
			appEntity.FilterExpression = filterExpr
		}

		return appEntity, nil
	} else if entityModel.Website != nil {
		// Map website entity
		websiteEntity := restapi.SloWebsiteEntity{
			Type:       "website",
			WebsiteId:  getStringPointer(entityModel.Website.WebsiteID),
			BeaconType: getStringPointer(entityModel.Website.BeaconType),
		}

		// Map filter expression if set
		if !entityModel.Website.FilterExpression.IsNull() && !entityModel.Website.FilterExpression.IsUnknown() {
			filterExpr, err := r.mapTagFilterStringToAPIModel(entityModel.Website.FilterExpression.ValueString())
			if err != nil {
				return nil, fmt.Errorf("failed to map website entity filter expression: %v", err)
			}
			websiteEntity.FilterExpression = filterExpr
		}

		return websiteEntity, nil
	} else if entityModel.Synthetic != nil {
		// Map synthetic entity
		syntheticEntity := restapi.SloSyntheticEntity{
			Type: "synthetic",
		}

		// Map synthetic test IDs
		if !entityModel.Synthetic.SyntheticTestIDs.IsNull() && !entityModel.Synthetic.SyntheticTestIDs.IsUnknown() {
			var testIDs []string
			diags.Append(entityModel.Synthetic.SyntheticTestIDs.ElementsAs(ctx, &testIDs, false)...)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to map synthetic test IDs: %v", diags)
			}

			// Convert string slice to interface slice
			testIDsInterface := make([]interface{}, len(testIDs))
			for i, id := range testIDs {
				testIDsInterface[i] = id
			}
			syntheticEntity.SyntheticTestIDs = testIDsInterface
		}

		// Map filter expression if set
		if !entityModel.Synthetic.FilterExpression.IsNull() && !entityModel.Synthetic.FilterExpression.IsUnknown() {
			filterExpr, err := r.mapTagFilterStringToAPIModel(entityModel.Synthetic.FilterExpression.ValueString())
			if err != nil {
				return nil, fmt.Errorf("failed to map synthetic entity filter expression: %v", err)
			}
			syntheticEntity.FilterExpression = filterExpr
		}

		return syntheticEntity, nil
	}

	return nil, fmt.Errorf("exactly one entity type configuration is required")
}

// mapIndicatorModelToAPIObject maps the indicator model to the API object
func (r *sloConfigResource) mapIndicatorModelToAPIObject(ctx context.Context, model SloConfigModel, diags diag.Diagnostics) (interface{}, error) {
	var indicatorModels []SloIndicatorModel
	diags.Append(model.Indicator.ElementsAs(ctx, &indicatorModels, false)...)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to map indicator: %v", diags)
	}

	if len(indicatorModels) != 1 {
		return nil, fmt.Errorf("exactly one indicator configuration is required")
	}

	indicatorModel := indicatorModels[0]

	// Check which indicator type is set
	if indicatorModel.TimeBasedLatency != nil {
		return restapi.SloTimeBasedLatencyIndicator{
			Blueprint:   "latency",
			Type:        "timeBased",
			Threshold:   indicatorModel.TimeBasedLatency.Threshold.ValueFloat64(),
			Aggregation: indicatorModel.TimeBasedLatency.Aggregation.ValueString(),
		}, nil
	} else if indicatorModel.EventBasedLatency != nil {
		return restapi.SloEventBasedLatencyIndicator{
			Blueprint: "latency",
			Type:      "eventBased",
			Threshold: indicatorModel.EventBasedLatency.Threshold.ValueFloat64(),
		}, nil
	} else if indicatorModel.TimeBasedAvailability != nil {
		return restapi.SloTimeBasedAvailabilityIndicator{
			Blueprint:   "availability",
			Type:        "timeBased",
			Threshold:   indicatorModel.TimeBasedAvailability.Threshold.ValueFloat64(),
			Aggregation: indicatorModel.TimeBasedAvailability.Aggregation.ValueString(),
		}, nil
	} else if indicatorModel.EventBasedAvailability != nil {
		return restapi.SloEventBasedAvailabilityIndicator{
			Blueprint: "availability",
			Type:      "eventBased",
		}, nil
	} else if indicatorModel.Traffic != nil {
		return restapi.SloTrafficIndicator{
			Blueprint:   "traffic",
			TrafficType: indicatorModel.Traffic.TrafficType.ValueString(),
			Threshold:   indicatorModel.Traffic.Threshold.ValueFloat64(),
			Aggregation: indicatorModel.Traffic.Aggregation.ValueString(),
		}, nil
	} else if indicatorModel.Custom != nil {
		customIndicator := restapi.SloCustomIndicator{
			Type:      "eventBased",
			Blueprint: "custom",
		}

		// Map good event filter expression
		if !indicatorModel.Custom.GoodEventFilterExpression.IsNull() && !indicatorModel.Custom.GoodEventFilterExpression.IsUnknown() {
			filterExpr, err := r.mapTagFilterStringToAPIModel(indicatorModel.Custom.GoodEventFilterExpression.ValueString())
			if err != nil {
				return nil, fmt.Errorf("failed to map good event filter expression: %v", err)
			}
			customIndicator.GoodEventFilterExpression = filterExpr
		}

		// Map bad event filter expression if set
		if !indicatorModel.Custom.BadEventFilterExpression.IsNull() && !indicatorModel.Custom.BadEventFilterExpression.IsUnknown() {
			filterExpr, err := r.mapTagFilterStringToAPIModel(indicatorModel.Custom.BadEventFilterExpression.ValueString())
			if err != nil {
				return nil, fmt.Errorf("failed to map bad event filter expression: %v", err)
			}
			customIndicator.BadEventFilterExpression = filterExpr
		}

		return customIndicator, nil
	}

	return nil, fmt.Errorf("exactly one indicator type configuration is required")
}

// mapTimeWindowModelToAPIObject maps the time window model to the API object
func (r *sloConfigResource) mapTimeWindowModelToAPIObject(ctx context.Context, model SloConfigModel, diags diag.Diagnostics) (interface{}, error) {
	var timeWindowModels []SloTimeWindowModel
	diags.Append(model.TimeWindow.ElementsAs(ctx, &timeWindowModels, false)...)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to map time window: %v", diags)
	}

	if len(timeWindowModels) != 1 {
		return nil, fmt.Errorf("exactly one time window configuration is required")
	}

	timeWindowModel := timeWindowModels[0]

	// Check which time window type is set
	if timeWindowModel.Rolling != nil {
		return restapi.SloRollingTimeWindow{
			Type:         "rolling",
			Duration:     int(timeWindowModel.Rolling.Duration.ValueInt64()),
			DurationUnit: timeWindowModel.Rolling.DurationUnit.ValueString(),
			Timezone:     timeWindowModel.Rolling.Timezone.ValueString(),
		}, nil
	} else if timeWindowModel.Fixed != nil {
		return restapi.SloFixedTimeWindow{
			Type:         "fixed",
			Duration:     int(timeWindowModel.Fixed.Duration.ValueInt64()),
			DurationUnit: timeWindowModel.Fixed.DurationUnit.ValueString(),
			Timezone:     timeWindowModel.Fixed.Timezone.ValueString(),
			StartTime:    timeWindowModel.Fixed.StartTime.ValueFloat64(),
		}, nil
	}

	return nil, fmt.Errorf("exactly one time window type configuration is required")
}

// mapAPIObjectToModel maps the API object to the Terraform model
func (r *sloConfigResource) mapAPIObjectToModel(ctx context.Context, sloConfig *restapi.SloConfig, model *SloConfigModel, diags diag.Diagnostics) error {
	model.ID = types.StringValue(sloConfig.ID)
	model.Name = types.StringValue(sloConfig.Name)
	model.Target = types.Float64Value(sloConfig.Target)

	// Map tags
	if tags, ok := sloConfig.Tags.([]interface{}); ok {
		tagValues := make([]string, 0, len(tags))
		for _, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				tagValues = append(tagValues, tagStr)
			}
		}
		model.Tags = basetypes.NewListValueMust(types.StringType, convertToStringValues(tagValues))
	} else if tags, ok := sloConfig.Tags.([]string); ok {
		model.Tags = basetypes.NewListValueMust(types.StringType, convertToStringValues(tags))
	} else {
		model.Tags = types.ListNull(types.StringType)
	}
	// Map entity
	entity, err := r.mapAPIEntityToModel(ctx, sloConfig, diags)
	if err != nil {
		return fmt.Errorf("failed to map entity from API: %v", err)
	}
	model.Entity = entity

	// Map indicator
	indicator, err := r.mapAPIIndicatorToModel(ctx, sloConfig, diags)
	if err != nil {
		return fmt.Errorf("failed to map indicator from API: %v", err)
	}
	model.Indicator = indicator

	// Map time window
	timeWindow, err := r.mapAPITimeWindowToModel(ctx, sloConfig, diags)
	if err != nil {
		return fmt.Errorf("failed to map time window from API: %v", err)
	}
	model.TimeWindow = timeWindow

	return nil
}

// mapAPIEntityToModel maps the API entity to the Terraform model
func (r *sloConfigResource) mapAPIEntityToModel(ctx context.Context, sloConfig *restapi.SloConfig, diags diag.Diagnostics) (types.List, error) {
	if sloConfig.Entity == nil {
		return types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{
			"application": types.ListType{ElemType: types.ObjectType{}},
			"website":     types.ListType{ElemType: types.ObjectType{}},
			"synthetic":   types.ListType{ElemType: types.ObjectType{}},
		}}), nil
	}

	entityAttrValues := make(map[string]attr.Value)

	// Check entity type and map accordingly
	if entity, ok := sloConfig.Entity.(map[string]interface{}); ok {
		entityType, ok := entity["type"].(string)
		if !ok {
			return types.ListNull(types.ObjectType{}), fmt.Errorf("entity type not found in API response")
		}

		switch entityType {
		case "application":
			// Map application entity
			appEntity := SloApplicationEntityModel{}

			if appID, ok := entity["applicationId"].(string); ok {
				appEntity.ApplicationID = types.StringValue(appID)
			}

			if serviceID, ok := entity["serviceId"].(string); ok && serviceID != "" {
				appEntity.ServiceID = types.StringValue(serviceID)
			}

			if endpointID, ok := entity["endpointId"].(string); ok && endpointID != "" {
				appEntity.EndpointID = types.StringValue(endpointID)
			}

			if boundaryScope, ok := entity["boundaryScope"].(string); ok {
				appEntity.BoundaryScope = types.StringValue(boundaryScope)
			}

			if includeSynthetic, ok := entity["includeSynthetic"].(bool); ok {
				appEntity.IncludeSynthetic = types.BoolValue(includeSynthetic)
			}

			if includeInternal, ok := entity["includeInternal"].(bool); ok {
				appEntity.IncludeInternal = types.BoolValue(includeInternal)
			}

			// Map filter expression
			filterExpr, err := r.mapFilterExpressionFromAPIModel(entity["tagFilterExpression"])
			if err == nil && filterExpr != nil {
				appEntity.FilterExpression = types.StringValue(filterExpr.Render())
			}

			appEntityList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
				"application_id":    types.StringType,
				"service_id":        types.StringType,
				"endpoint_id":       types.StringType,
				"boundary_scope":    types.StringType,
				"include_synthetic": types.BoolType,
				"include_internal":  types.BoolType,
				"filter_expression": types.StringType,
			}}, []SloApplicationEntityModel{appEntity})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create application entity list: %v", diags)
			}

			entityAttrValues["application"] = appEntityList
			entityAttrValues["website"] = types.ListNull(types.ObjectType{})
			entityAttrValues["synthetic"] = types.ListNull(types.ObjectType{})

		case "website":
			// Map website entity
			websiteEntity := SloWebsiteEntityModel{}

			if websiteID, ok := entity["websiteId"].(string); ok {
				websiteEntity.WebsiteID = types.StringValue(websiteID)
			}

			if beaconType, ok := entity["beaconType"].(string); ok {
				websiteEntity.BeaconType = types.StringValue(beaconType)
			}

			// Map filter expression
			filterExpr, err := r.mapFilterExpressionFromAPIModel(entity["tagFilterExpression"])
			if err == nil && filterExpr != nil {
				websiteEntity.FilterExpression = types.StringValue(filterExpr.Render())
			}

			websiteEntityList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
				"website_id":        types.StringType,
				"beacon_type":       types.StringType,
				"filter_expression": types.StringType,
			}}, []SloWebsiteEntityModel{websiteEntity})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create website entity list: %v", diags)
			}

			entityAttrValues["application"] = types.ListNull(types.ObjectType{})
			entityAttrValues["website"] = websiteEntityList
			entityAttrValues["synthetic"] = types.ListNull(types.ObjectType{})

		case "synthetic":
			// Map synthetic entity
			syntheticEntity := SloSyntheticEntityModel{}

			// Map synthetic test IDs
			if testIDs, ok := entity["syntheticTestIds"].([]interface{}); ok {
				stringIDs := make([]string, 0, len(testIDs))
				for _, id := range testIDs {
					if strID, ok := id.(string); ok {
						stringIDs = append(stringIDs, strID)
					}
				}
				syntheticEntity.SyntheticTestIDs = basetypes.NewListValueMust(types.StringType, convertToStringValues(stringIDs))
			}

			// Map filter expression
			filterExpr, err := r.mapFilterExpressionFromAPIModel(entity["tagFilterExpression"])
			if err == nil && filterExpr != nil {
				syntheticEntity.FilterExpression = types.StringValue(filterExpr.Render())
			}

			syntheticEntityList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
				"synthetic_test_ids": types.ListType{ElemType: types.StringType},
				"filter_expression":  types.StringType,
			}}, []SloSyntheticEntityModel{syntheticEntity})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create synthetic entity list: %v", diags)
			}

			entityAttrValues["application"] = types.ListNull(types.ObjectType{})
			entityAttrValues["website"] = types.ListNull(types.ObjectType{})
			entityAttrValues["synthetic"] = syntheticEntityList
		}
	}

	// Convert to Terraform object
	entityValue, diags := types.ObjectValue(
		map[string]attr.Type{
			"application": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"application_id":    types.StringType,
				"service_id":        types.StringType,
				"endpoint_id":       types.StringType,
				"boundary_scope":    types.StringType,
				"include_synthetic": types.BoolType,
				"include_internal":  types.BoolType,
				"filter_expression": types.StringType,
			}}},
			"website": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"website_id":        types.StringType,
				"beacon_type":       types.StringType,
				"filter_expression": types.StringType,
			}}},
			"synthetic": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"synthetic_test_ids": types.ListType{ElemType: types.StringType},
				"filter_expression":  types.StringType,
			}}},
		},
		entityAttrValues,
	)

	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create entity object value: %v", diags)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
		"application": types.ListType{ElemType: types.ObjectType{}},
		"website":     types.ListType{ElemType: types.ObjectType{}},
		"synthetic":   types.ListType{ElemType: types.ObjectType{}},
	}}, []attr.Value{entityValue}), nil
}

// mapAPIIndicatorToModel maps the API indicator to the Terraform model
func (r *sloConfigResource) mapAPIIndicatorToModel(ctx context.Context, sloConfig *restapi.SloConfig, diags diag.Diagnostics) (types.List, error) {
	if sloConfig.Indicator == nil {
		return types.ListNull(types.ObjectType{}), nil
	}

	indicatorAttrValues := make(map[string]attr.Value)
	indicatorAttrValues["time_based_latency"] = types.ListNull(types.ObjectType{})
	indicatorAttrValues["event_based_latency"] = types.ListNull(types.ObjectType{})
	indicatorAttrValues["time_based_availability"] = types.ListNull(types.ObjectType{})
	indicatorAttrValues["event_based_availability"] = types.ListNull(types.ObjectType{})
	indicatorAttrValues["traffic"] = types.ListNull(types.ObjectType{})
	indicatorAttrValues["custom"] = types.ListNull(types.ObjectType{})

	// Check indicator type and map accordingly
	if indicator, ok := sloConfig.Indicator.(map[string]interface{}); ok {
		blueprint, _ := indicator["blueprint"].(string)
		indicatorType, _ := indicator["type"].(string)

		switch {
		case blueprint == "latency" && indicatorType == "timeBased":
			// Map time-based latency indicator
			latencyIndicator := SloTimeBasedLatencyIndicatorModel{}

			if threshold, ok := indicator["threshold"].(float64); ok {
				latencyIndicator.Threshold = types.Float64Value(threshold)
			}

			if aggregation, ok := indicator["aggregation"].(string); ok {
				latencyIndicator.Aggregation = types.StringValue(aggregation)
			}

			latencyIndicatorList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
				"threshold":   types.Float64Type,
				"aggregation": types.StringType,
			}}, []SloTimeBasedLatencyIndicatorModel{latencyIndicator})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create time-based latency indicator list: %v", diags)
			}

			indicatorAttrValues["time_based_latency"] = latencyIndicatorList

		case blueprint == "latency" && indicatorType == "eventBased":
			// Map event-based latency indicator
			latencyIndicator := SloEventBasedLatencyIndicatorModel{}

			if threshold, ok := indicator["threshold"].(float64); ok {
				latencyIndicator.Threshold = types.Float64Value(threshold)
			}

			latencyIndicatorList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
				"threshold": types.Float64Type,
			}}, []SloEventBasedLatencyIndicatorModel{latencyIndicator})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create event-based latency indicator list: %v", diags)
			}

			indicatorAttrValues["event_based_latency"] = latencyIndicatorList

		case blueprint == "availability" && indicatorType == "timeBased":
			// Map time-based availability indicator
			availabilityIndicator := SloTimeBasedAvailabilityIndicatorModel{}

			if threshold, ok := indicator["threshold"].(float64); ok {
				availabilityIndicator.Threshold = types.Float64Value(threshold)
			}

			if aggregation, ok := indicator["aggregation"].(string); ok {
				availabilityIndicator.Aggregation = types.StringValue(aggregation)
			}

			availabilityIndicatorList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
				"threshold":   types.Float64Type,
				"aggregation": types.StringType,
			}}, []SloTimeBasedAvailabilityIndicatorModel{availabilityIndicator})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create time-based availability indicator list: %v", diags)
			}

			indicatorAttrValues["time_based_availability"] = availabilityIndicatorList

		case blueprint == "availability" && indicatorType == "eventBased":
			// Map event-based availability indicator
			availabilityIndicatorList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{}}, []SloEventBasedAvailabilityIndicatorModel{{}})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create event-based availability indicator list: %v", diags)
			}

			indicatorAttrValues["event_based_availability"] = availabilityIndicatorList

		case blueprint == "traffic":
			// Map traffic indicator
			trafficIndicator := SloTrafficIndicatorModel{}

			if trafficType, ok := indicator["trafficType"].(string); ok {
				trafficIndicator.TrafficType = types.StringValue(trafficType)
			}

			if threshold, ok := indicator["threshold"].(float64); ok {
				trafficIndicator.Threshold = types.Float64Value(threshold)
			}

			if aggregation, ok := indicator["aggregation"].(string); ok {
				trafficIndicator.Aggregation = types.StringValue(aggregation)
			}

			trafficIndicatorList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
				"traffic_type": types.StringType,
				"threshold":    types.Float64Type,
				"aggregation":  types.StringType,
			}}, []SloTrafficIndicatorModel{trafficIndicator})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create traffic indicator list: %v", diags)
			}

			indicatorAttrValues["traffic"] = trafficIndicatorList

		case blueprint == "custom":
			// Map custom indicator
			customIndicator := SloCustomIndicatorModel{}

			// Map good event filter expression
			goodFilterExpr, err := r.mapFilterExpressionFromAPIModel(indicator["goodEventsFilter"])
			if err == nil && goodFilterExpr != nil {
				customIndicator.GoodEventFilterExpression = types.StringValue(goodFilterExpr.Render())
			}

			// Map bad event filter expression
			badFilterExpr, err := r.mapFilterExpressionFromAPIModel(indicator["badEventsFilter"])
			if err == nil && badFilterExpr != nil {
				customIndicator.BadEventFilterExpression = types.StringValue(badFilterExpr.Render())
			}

			customIndicatorList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
				"good_event_filter_expression": types.StringType,
				"bad_event_filter_expression":  types.StringType,
			}}, []SloCustomIndicatorModel{customIndicator})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create custom indicator list: %v", diags)
			}

			indicatorAttrValues["custom"] = customIndicatorList
		}
	}

	// Convert to Terraform object
	indicatorValue, diags := types.ObjectValue(
		map[string]attr.Type{
			"time_based_latency": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"threshold":   types.Float64Type,
				"aggregation": types.StringType,
			}}},
			"event_based_latency": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"threshold": types.Float64Type,
			}}},
			"time_based_availability": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"threshold":   types.Float64Type,
				"aggregation": types.StringType,
			}}},
			"event_based_availability": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{}}},
			"traffic": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"traffic_type": types.StringType,
				"threshold":    types.Float64Type,
				"aggregation":  types.StringType,
			}}},
			"custom": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"good_event_filter_expression": types.StringType,
				"bad_event_filter_expression":  types.StringType,
			}}},
		},
		indicatorAttrValues,
	)

	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create indicator object value: %v", diags)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
		"time_based_latency":       types.ListType{ElemType: types.ObjectType{}},
		"event_based_latency":      types.ListType{ElemType: types.ObjectType{}},
		"time_based_availability":  types.ListType{ElemType: types.ObjectType{}},
		"event_based_availability": types.ListType{ElemType: types.ObjectType{}},
		"traffic":                  types.ListType{ElemType: types.ObjectType{}},
		"custom":                   types.ListType{ElemType: types.ObjectType{}},
	}}, []attr.Value{indicatorValue}), nil
}

// mapAPITimeWindowToModel maps the API time window to the Terraform model
func (r *sloConfigResource) mapAPITimeWindowToModel(ctx context.Context, sloConfig *restapi.SloConfig, diags diag.Diagnostics) (types.List, error) {
	if sloConfig.TimeWindow == nil {
		return types.ListNull(types.ObjectType{}), nil
	}

	timeWindowAttrValues := make(map[string]attr.Value)
	timeWindowAttrValues["rolling"] = types.ListNull(types.ObjectType{})
	timeWindowAttrValues["fixed"] = types.ListNull(types.ObjectType{})

	// Check time window type and map accordingly
	if timeWindow, ok := sloConfig.TimeWindow.(map[string]interface{}); ok {
		timeWindowType, _ := timeWindow["type"].(string)

		switch timeWindowType {
		case "rolling":
			// Map rolling time window
			rollingTimeWindow := SloRollingTimeWindowModel{}

			if duration, ok := timeWindow["duration"].(float64); ok {
				rollingTimeWindow.Duration = types.Int64Value(int64(duration))
			}

			if durationUnit, ok := timeWindow["durationUnit"].(string); ok {
				rollingTimeWindow.DurationUnit = types.StringValue(durationUnit)
			}

			if timezone, ok := timeWindow["timezone"].(string); ok {
				rollingTimeWindow.Timezone = types.StringValue(timezone)
			} else {
				rollingTimeWindow.Timezone = types.StringValue("UTC")
			}

			rollingTimeWindowList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
				"duration":      types.Int64Type,
				"duration_unit": types.StringType,
				"timezone":      types.StringType,
			}}, []SloRollingTimeWindowModel{rollingTimeWindow})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create rolling time window list: %v", diags)
			}

			timeWindowAttrValues["rolling"] = rollingTimeWindowList

		case "fixed":
			// Map fixed time window
			fixedTimeWindow := SloFixedTimeWindowModel{}

			if duration, ok := timeWindow["duration"].(float64); ok {
				fixedTimeWindow.Duration = types.Int64Value(int64(duration))
			}

			if durationUnit, ok := timeWindow["durationUnit"].(string); ok {
				fixedTimeWindow.DurationUnit = types.StringValue(durationUnit)
			}

			if timezone, ok := timeWindow["timezone"].(string); ok {
				fixedTimeWindow.Timezone = types.StringValue(timezone)
			} else {
				fixedTimeWindow.Timezone = types.StringValue("UTC")
			}

			if startTime, ok := timeWindow["startTimestamp"].(float64); ok {
				fixedTimeWindow.StartTime = types.Float64Value(startTime)
			}

			fixedTimeWindowList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
				"duration":        types.Int64Type,
				"duration_unit":   types.StringType,
				"timezone":        types.StringType,
				"start_timestamp": types.Float64Type,
			}}, []SloFixedTimeWindowModel{fixedTimeWindow})

			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create fixed time window list: %v", diags)
			}

			timeWindowAttrValues["fixed"] = fixedTimeWindowList
		}
	}

	// Convert to Terraform object
	timeWindowValue, diags := types.ObjectValue(
		map[string]attr.Type{
			"rolling": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"duration":      types.Int64Type,
				"duration_unit": types.StringType,
				"timezone":      types.StringType,
			}}},
			"fixed": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"duration":        types.Int64Type,
				"duration_unit":   types.StringType,
				"timezone":        types.StringType,
				"start_timestamp": types.Float64Type,
			}}},
		},
		timeWindowAttrValues,
	)

	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), fmt.Errorf("failed to create time window object value: %v", diags)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
		"rolling": types.ListType{ElemType: types.ObjectType{}},
		"fixed":   types.ListType{ElemType: types.ObjectType{}},
	}}, []attr.Value{timeWindowValue}), nil
}

// mapFilterExpressionFromAPIModel is a helper method to map filter expression from API model
func (r *sloConfigResource) mapFilterExpressionFromAPIModel(filter interface{}) (*tagfilter.FilterExpression, error) {
	if filter == nil {
		return nil, nil
	}

	if tagFilterString, ok := filter.(*restapi.TagFilter); ok {
		mapper := tagfilter.NewMapper()
		t, err := mapper.FromAPIModel(tagFilterString)
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	return nil, nil
}

// mapTagFilterStringToAPIModel maps a tag filter string to the API model
func (r *sloConfigResource) mapTagFilterStringToAPIModel(input string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

// Helper functions
func getStringPointer(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	v := value.ValueString()
	return &v
}

func getBoolPointer(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	v := value.ValueBool()
	return &v
}

func convertToStringValues(values []string) []attr.Value {
	result := make([]attr.Value, len(values))
	for i, v := range values {
		result[i] = types.StringValue(v)
	}
	return result
}

// Made with Bob
