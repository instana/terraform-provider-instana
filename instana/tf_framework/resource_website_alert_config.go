package tf_framework

import (
	"context"
	"fmt"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &websiteAlertConfigResource{}
	_ resource.ResourceWithConfigure   = &websiteAlertConfigResource{}
	_ resource.ResourceWithImportState = &websiteAlertConfigResource{}
)

// NewWebsiteAlertConfigResource is a helper function to simplify the provider implementation
func NewWebsiteAlertConfigResource() resource.Resource {
	return &websiteAlertConfigResource{}
}

// websiteAlertConfigResource is the resource implementation
type websiteAlertConfigResource struct {
	client restapi.InstanaAPI
}

// Configure adds the provider configured client to the resource
func (r *websiteAlertConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *websiteAlertConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_website_alert_config"
}

// Schema defines the schema for the resource
func (r *websiteAlertConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to configure a Website Alert Configuration in Instana.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the Website Alert Configuration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Website Alert Configuration.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 256),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the Website Alert Configuration.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 65536),
				},
			},
			"severity": schema.StringAttribute{
				Description: "The severity of the alert when triggered.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("warning", "critical"),
				},
			},
			"triggering": schema.BoolAttribute{
				Description: "Flag to indicate whether also an Incident is triggered or not.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"website_id": schema.StringAttribute{
				Description: "Unique ID of the website.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 64),
				},
			},
			"tag_filter": schema.StringAttribute{
				Description: "The tag filter expression for the Website Alert Configuration.",
				Optional:    true,
			},
			"alert_channel_ids": schema.SetAttribute{
				Description: "List of IDs of alert channels defined in Instana.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"granularity": schema.Int64Attribute{
				Description: "The evaluation granularity used for detection of violations of the defined threshold.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(600000),
			},
			"custom_payload_fields": schema.SetNestedAttribute{
				Description: "A list of custom payload fields.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "The key of the custom payload field.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value of a static string custom payload field.",
							Optional:    true,
						},
					},
				},
			},
			"threshold": schema.SingleNestedAttribute{
				Description: "The threshold configuration for the Website Alert Configuration.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"operator": schema.StringAttribute{
						Description: "The operator of the threshold.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf(">", ">=", "<", "<=", "=="),
						},
					},
					"value": schema.Float64Attribute{
						Description: "The value of the threshold.",
						Required:    true,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				Description: "Indicates the type of rule this alert configuration is about.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"slowness": schema.ListNestedBlock{
							Description: "Rule based on the slowness of the configured alert configuration target.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Description: "The metric name of the website alert rule.",
										Required:    true,
									},
									"aggregation": schema.StringAttribute{
										Description: "The aggregation function of the website alert rule.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
										},
									},
								},
							},
						},
						"specific_js_error": schema.ListNestedBlock{
							Description: "Rule based on a specific javascript error of the configured alert configuration target.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Description: "The metric name of the website alert rule.",
										Required:    true,
									},
									"aggregation": schema.StringAttribute{
										Description: "The aggregation function of the website alert rule.",
										Optional:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
										},
									},
									"operator": schema.StringAttribute{
										Description: "The operator which will be applied to evaluate this rule.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
										},
									},
									"value": schema.StringAttribute{
										Description: "The value identify the specific javascript error.",
										Optional:    true,
									},
								},
							},
						},
						"status_code": schema.ListNestedBlock{
							Description: "Rule based on the HTTP status code of the configured alert configuration target.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Description: "The metric name of the website alert rule.",
										Required:    true,
									},
									"aggregation": schema.StringAttribute{
										Description: "The aggregation function of the website alert rule.",
										Optional:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
										},
									},
									"operator": schema.StringAttribute{
										Description: "The operator which will be applied to evaluate this rule.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
										},
									},
									"value": schema.StringAttribute{
										Description: "The value identify the specific http status code.",
										Required:    true,
									},
								},
							},
						},
						"throughput": schema.ListNestedBlock{
							Description: "Rule based on the throughput of the configured alert configuration target.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Description: "The metric name of the website alert rule.",
										Required:    true,
									},
									"aggregation": schema.StringAttribute{
										Description: "The aggregation function of the website alert rule.",
										Optional:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
										},
									},
								},
							},
						},
					},
				},
			},
			"time_threshold": schema.ListNestedBlock{
				Description: "Indicates the type of violation of the defined threshold.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"user_impact_of_violations_in_sequence": schema.ListNestedBlock{
							Description: "Time threshold base on user impact of violations in sequence.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Description: "The time window if the time threshold.",
										Optional:    true,
									},
									"impact_measurement_method": schema.StringAttribute{
										Description: "The impact method of the time threshold based on user impact of violations in sequence.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("AGGREGATED", "PER_WINDOW"),
										},
									},
									"user_percentage": schema.Float64Attribute{
										Description: "The percentage of impacted users of the time threshold based on user impact of violations in sequence.",
										Optional:    true,
									},
									"users": schema.Int64Attribute{
										Description: "The number of impacted users of the time threshold based on user impact of violations in sequence.",
										Optional:    true,
									},
								},
							},
						},
						"violations_in_period": schema.ListNestedBlock{
							Description: "Time threshold base on violations in period.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Description: "The time window if the time threshold.",
										Optional:    true,
									},
									"violations": schema.Int64Attribute{
										Description: "The violations appeared in the period.",
										Optional:    true,
									},
								},
							},
						},
						"violations_in_sequence": schema.ListNestedBlock{
							Description: "Time threshold base on violations in sequence.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Description: "The time window if the time threshold.",
										Optional:    true,
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
func (r *websiteAlertConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan WebsiteAlertConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert from Terraform model to API object
	websiteAlertConfig, err := r.mapModelToAPIObject(ctx, plan, resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Website Alert Configuration",
			fmt.Sprintf("Could not map Website Alert Configuration to API object: %s", err),
		)
		return
	}

	// Create new Website Alert Configuration
	websiteAlertConfigCreated, err := r.client.WebsiteAlertConfig().Create(websiteAlertConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Website Alert Configuration",
			fmt.Sprintf("Could not create Website Alert Configuration: %s", err),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	err = r.mapAPIObjectToModel(ctx, websiteAlertConfigCreated, &plan, resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Website Alert Configuration",
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
func (r *websiteAlertConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state WebsiteAlertConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Website Alert Configuration value from Instana
	websiteAlertConfig, err := r.client.WebsiteAlertConfig().GetOne(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Website Alert Configuration",
			fmt.Sprintf("Could not read Website Alert Configuration ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Map response body to model
	err = r.mapAPIObjectToModel(ctx, websiteAlertConfig, &state, resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Website Alert Configuration",
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
func (r *websiteAlertConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan WebsiteAlertConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state WebsiteAlertConfigModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve ID from state
	plan.ID = state.ID

	// Convert from Terraform model to API object
	websiteAlertConfig, err := r.mapModelToAPIObject(ctx, plan, resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Website Alert Configuration",
			fmt.Sprintf("Could not map Website Alert Configuration to API object: %s", err),
		)
		return
	}

	// Update Website Alert Configuration
	websiteAlertConfigUpdated, err := r.client.WebsiteAlertConfig().Update(websiteAlertConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Website Alert Configuration",
			fmt.Sprintf("Could not update Website Alert Configuration: %s", err),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	err = r.mapAPIObjectToModel(ctx, websiteAlertConfigUpdated, &plan, resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Website Alert Configuration",
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
func (r *websiteAlertConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state WebsiteAlertConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Website Alert Configuration
	err := r.client.WebsiteAlertConfig().DeleteByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Website Alert Configuration",
			fmt.Sprintf("Could not delete Website Alert Configuration ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

// ImportState imports the resource into Terraform state
func (r *websiteAlertConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapModelToAPIObject maps the Terraform model to the API object
func (r *websiteAlertConfigResource) mapModelToAPIObject(ctx context.Context, model WebsiteAlertConfigModel, diags diag.Diagnostics) (*restapi.WebsiteAlertConfig, error) {
	// Convert severity from Terraform representation to API representation
	severity, err := convertSeverityFromTerraformToInstanaAPIRepresentation(model.Severity.ValueString())
	if err != nil {
		return nil, err
	}

	// Map tag filter expression
	var tagFilter *restapi.TagFilter
	if !model.TagFilter.IsNull() && !model.TagFilter.IsUnknown() {
		tagFilter, err = r.mapTagFilterExpressionFromSchema(model.TagFilter.ValueString())
		if err != nil {
			return nil, err
		}
	}

	// Map alert channel IDs
	var alertChannelIDs []string
	if !model.AlertChannelIDs.IsNull() && !model.AlertChannelIDs.IsUnknown() {
		diags.Append(model.AlertChannelIDs.ElementsAs(ctx, &alertChannelIDs, false)...)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to map alert channel IDs: %v", diags)
		}
	}

	// Map custom payload fields
	var customPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayloadFields.IsNull() && !model.CustomPayloadFields.IsUnknown() {
		var customPayloadFieldModels []CustomPayloadFieldModel
		diags.Append(model.CustomPayloadFields.ElementsAs(ctx, &customPayloadFieldModels, false)...)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to map custom payload fields: %v", diags)
		}

		for _, field := range customPayloadFieldModels {
			customPayloadFields = append(customPayloadFields, restapi.CustomPayloadField[any]{
				Key:   field.Key.ValueString(),
				Value: field.Value.ValueString(),
			})
		}
	}

	// Map rule
	rule, err := r.mapRuleFromModel(ctx, model, diags)
	if err != nil {
		return nil, err
	}

	// Map threshold
	threshold, err := r.mapThresholdFromModel(ctx, model, diags)
	if err != nil {
		return nil, err
	}

	// Map time threshold
	timeThreshold, err := r.mapTimeThresholdFromModel(ctx, model, diags)
	if err != nil {
		return nil, err
	}

	// Create API object
	websiteAlertConfig := &restapi.WebsiteAlertConfig{
		ID:                    model.ID.ValueString(),
		Name:                  model.Name.ValueString(),
		Description:           model.Description.ValueString(),
		Severity:              severity,
		Triggering:            model.Triggering.ValueBool(),
		WebsiteID:             model.WebsiteID.ValueString(),
		TagFilterExpression:   tagFilter,
		AlertChannelIDs:       alertChannelIDs,
		Granularity:           restapi.Granularity(model.Granularity.ValueInt64()),
		CustomerPayloadFields: customPayloadFields,
		Rule:                  *rule,
		Threshold:             *threshold,
		TimeThreshold:         *timeThreshold,
	}

	return websiteAlertConfig, nil
}

// mapRuleFromModel maps the rule from the Terraform model to the API object
func (r *websiteAlertConfigResource) mapRuleFromModel(ctx context.Context, model WebsiteAlertConfigModel, diags diag.Diagnostics) (*restapi.WebsiteAlertRule, error) {
	var ruleModels []WebsiteAlertRuleModel
	diags.Append(model.Rule.ElementsAs(ctx, &ruleModels, false)...)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to map rule: %v", diags)
	}

	if len(ruleModels) != 1 {
		return nil, fmt.Errorf("exactly one rule configuration is required")
	}

	ruleModel := ruleModels[0]

	// Check which rule type is set
	if !ruleModel.Slowness.IsNull() && !ruleModel.Slowness.IsUnknown() {
		var slownessModels []WebsiteAlertRuleConfigModel
		diags.Append(ruleModel.Slowness.ElementsAs(ctx, &slownessModels, false)...)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to map slowness rule: %v", diags)
		}

		if len(slownessModels) != 1 {
			return nil, fmt.Errorf("exactly one slowness rule configuration is required")
		}

		slownessModel := slownessModels[0]
		aggregation := restapi.Aggregation(slownessModel.Aggregation.ValueString())

		return &restapi.WebsiteAlertRule{
			AlertType:   "slowness",
			MetricName:  slownessModel.MetricName.ValueString(),
			Aggregation: &aggregation,
		}, nil
	} else if !ruleModel.SpecificJsError.IsNull() && !ruleModel.SpecificJsError.IsUnknown() {
		var specificJsErrorModels []WebsiteAlertRuleConfigModel
		diags.Append(ruleModel.SpecificJsError.ElementsAs(ctx, &specificJsErrorModels, false)...)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to map specific JS error rule: %v", diags)
		}

		if len(specificJsErrorModels) != 1 {
			return nil, fmt.Errorf("exactly one specific JS error rule configuration is required")
		}

		specificJsErrorModel := specificJsErrorModels[0]
		var aggregationPtr *restapi.Aggregation
		if !specificJsErrorModel.Aggregation.IsNull() && !specificJsErrorModel.Aggregation.IsUnknown() {
			aggregation := restapi.Aggregation(specificJsErrorModel.Aggregation.ValueString())
			aggregationPtr = &aggregation
		}

		operator := restapi.ExpressionOperator(specificJsErrorModel.Operator.ValueString())
		var valuePtr *string
		if !specificJsErrorModel.Value.IsNull() && !specificJsErrorModel.Value.IsUnknown() {
			value := specificJsErrorModel.Value.ValueString()
			valuePtr = &value
		}

		return &restapi.WebsiteAlertRule{
			AlertType:   "specificJsError",
			MetricName:  specificJsErrorModel.MetricName.ValueString(),
			Aggregation: aggregationPtr,
			Operator:    &operator,
			Value:       valuePtr,
		}, nil
	} else if !ruleModel.StatusCode.IsNull() && !ruleModel.StatusCode.IsUnknown() {
		var statusCodeModels []WebsiteAlertRuleConfigModel
		diags.Append(ruleModel.StatusCode.ElementsAs(ctx, &statusCodeModels, false)...)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to map status code rule: %v", diags)
		}

		if len(statusCodeModels) != 1 {
			return nil, fmt.Errorf("exactly one status code rule configuration is required")
		}

		statusCodeModel := statusCodeModels[0]
		var aggregationPtr *restapi.Aggregation
		if !statusCodeModel.Aggregation.IsNull() && !statusCodeModel.Aggregation.IsUnknown() {
			aggregation := restapi.Aggregation(statusCodeModel.Aggregation.ValueString())
			aggregationPtr = &aggregation
		}

		operator := restapi.ExpressionOperator(statusCodeModel.Operator.ValueString())
		value := statusCodeModel.Value.ValueString()

		return &restapi.WebsiteAlertRule{
			AlertType:   "statusCode",
			MetricName:  statusCodeModel.MetricName.ValueString(),
			Aggregation: aggregationPtr,
			Operator:    &operator,
			Value:       &value,
		}, nil
	} else if !ruleModel.Throughput.IsNull() && !ruleModel.Throughput.IsUnknown() {
		var throughputModels []WebsiteAlertRuleConfigModel
		diags.Append(ruleModel.Throughput.ElementsAs(ctx, &throughputModels, false)...)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to map throughput rule: %v", diags)
		}

		if len(throughputModels) != 1 {
			return nil, fmt.Errorf("exactly one throughput rule configuration is required")
		}

		throughputModel := throughputModels[0]
		var aggregationPtr *restapi.Aggregation
		if !throughputModel.Aggregation.IsNull() && !throughputModel.Aggregation.IsUnknown() {
			aggregation := restapi.Aggregation(throughputModel.Aggregation.ValueString())
			aggregationPtr = &aggregation
		}

		return &restapi.WebsiteAlertRule{
			AlertType:   "throughput",
			MetricName:  throughputModel.MetricName.ValueString(),
			Aggregation: aggregationPtr,
		}, nil
	}

	return nil, fmt.Errorf("exactly one rule type configuration is required")
}

// mapThresholdFromModel maps the threshold from the Terraform model to the API object
func (r *websiteAlertConfigResource) mapThresholdFromModel(ctx context.Context, model WebsiteAlertConfigModel, diags diag.Diagnostics) (*restapi.Threshold, error) {
	var thresholdModel ThresholdModel
	diags.Append(model.Threshold.As(ctx, &thresholdModel, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to map threshold: %v", diags)
	}

	return &restapi.Threshold{
		Operator: restapi.ThresholdOperator(thresholdModel.Operator.ValueString()),
		Value:    thresholdModel.Value.ValueFloat64(),
	}, nil
}

// mapTimeThresholdFromModel maps the time threshold from the Terraform model to the API object
func (r *websiteAlertConfigResource) mapTimeThresholdFromModel(ctx context.Context, model WebsiteAlertConfigModel, diags diag.Diagnostics) (*restapi.WebsiteTimeThreshold, error) {
	var timeThresholdModels []WebsiteTimeThresholdModel
	diags.Append(model.TimeThreshold.ElementsAs(ctx, &timeThresholdModels, false)...)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to map time threshold: %v", diags)
	}

	if len(timeThresholdModels) != 1 {
		return nil, fmt.Errorf("exactly one time threshold configuration is required")
	}

	timeThresholdModel := timeThresholdModels[0]

	// Check which time threshold type is set
	if !timeThresholdModel.UserImpactOfViolationsInSequence.IsNull() && !timeThresholdModel.UserImpactOfViolationsInSequence.IsUnknown() {
		var userImpactModels []WebsiteUserImpactOfViolationsInSequenceModel
		diags.Append(timeThresholdModel.UserImpactOfViolationsInSequence.ElementsAs(ctx, &userImpactMo

// Made with Bob
