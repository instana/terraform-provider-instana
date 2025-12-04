package websitemonitoringconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// NewWebsiteMonitoringConfigResourceHandle creates the resource handle for Website Monitoring Configurations
func NewWebsiteMonitoringConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.WebsiteMonitoringConfig] {
	return &websiteMonitoringConfigResource{
		metaData: createResourceMetaData(),
	}
}

// websiteMonitoringConfigResource implements the resource framework for website monitoring configurations
type websiteMonitoringConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

// MetaData returns the resource metadata
func (r *websiteMonitoringConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

// GetRestResource returns the REST resource for the API
func (r *websiteMonitoringConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.WebsiteMonitoringConfig] {
	return api.WebsiteMonitoringConfig()
}

// SetComputedFields sets computed fields in the plan (no computed fields to set for this resource)
func (r *websiteMonitoringConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return diag.Diagnostics{}
}

// MapStateToDataObject maps Terraform state/plan to API data object
func (r *websiteMonitoringConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.WebsiteMonitoringConfig, diag.Diagnostics) {
	if err := validateMapStateToDataObjectInputs(ctx, plan, state); err != nil {
		return nil, diag.Diagnostics{err}
	}

	model, diags := extractModelFromPlanOrState(ctx, plan, state)
	if diags.HasError() {
		return nil, diags
	}

	apiObject, err := mapModelToAPIObject(model)
	if err != nil {
		diags.AddError(WebsiteMonitoringConfigErrMappingToAPI, err.Error())
		return nil, diags
	}

	return apiObject, diags
}

// UpdateState updates Terraform state with data from API object
func (r *websiteMonitoringConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.WebsiteMonitoringConfig) diag.Diagnostics {
	if err := validateUpdateStateInputs(ctx, state, apiObject); err != nil {
		return diag.Diagnostics{err}
	}

	model := mapAPIObjectToModel(apiObject)

	return updateTerraformState(ctx, state, model)
}

// createResourceMetaData creates the resource metadata with schema definition
func createResourceMetaData() resourcehandle.ResourceMetaData {
	return resourcehandle.ResourceMetaData{
		ResourceName:  ResourceInstanaWebsiteMonitoringConfig,
		Schema:        createResourceSchema(),
		SchemaVersion: WebsiteMonitoringConfigSchemaVersion,
	}
}

// createResourceSchema creates the Terraform schema for the resource
func createResourceSchema() schema.Schema {
	return schema.Schema{
		Description: WebsiteMonitoringConfigDescResource,
		Attributes:  createSchemaAttributes(),
	}
}

// createSchemaAttributes creates the schema attributes map
func createSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		WebsiteMonitoringConfigFieldID:      createIDAttribute(),
		WebsiteMonitoringConfigFieldName:    createNameAttribute(),
		WebsiteMonitoringConfigFieldAppName: createAppNameAttribute(),
	}
}

// createIDAttribute creates the ID schema attribute
func createIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:      true,
		Description:   WebsiteMonitoringConfigDescID,
		PlanModifiers: createIDPlanModifiers(),
	}
}

// createNameAttribute creates the name schema attribute
func createNameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: WebsiteMonitoringConfigDescName,
	}
}

// createAppNameAttribute creates the app_name schema attribute
func createAppNameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:    true,
		Description: WebsiteMonitoringConfigDescAppName,
	}
}

// createIDPlanModifiers creates plan modifiers for the ID field
func createIDPlanModifiers() []planmodifier.String {
	return []planmodifier.String{
		stringplanmodifier.UseStateForUnknown(),
	}
}

// validateMapStateToDataObjectInputs validates inputs for MapStateToDataObject method
func validateMapStateToDataObjectInputs(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) diag.Diagnostic {
	if ctx == nil {
		return diag.NewErrorDiagnostic(WebsiteMonitoringConfigErrInvalidInput, WebsiteMonitoringConfigErrNilContext)
	}

	if plan == nil && state == nil {
		return diag.NewErrorDiagnostic(WebsiteMonitoringConfigErrInvalidInput, "Both plan and state cannot be nil")
	}

	return nil
}

// validateUpdateStateInputs validates inputs for UpdateState method
func validateUpdateStateInputs(ctx context.Context, state *tfsdk.State, apiObject *restapi.WebsiteMonitoringConfig) diag.Diagnostic {
	if ctx == nil {
		return diag.NewErrorDiagnostic(WebsiteMonitoringConfigErrInvalidInput, WebsiteMonitoringConfigErrNilContext)
	}

	if state == nil {
		return diag.NewErrorDiagnostic(WebsiteMonitoringConfigErrInvalidInput, WebsiteMonitoringConfigErrNilState)
	}

	if apiObject == nil {
		return diag.NewErrorDiagnostic(WebsiteMonitoringConfigErrInvalidInput, WebsiteMonitoringConfigErrNilAPIObject)
	}

	return nil
}

// extractModelFromPlanOrState extracts the model from plan or state
func extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*WebsiteMonitoringConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model WebsiteMonitoringConfigModel

	if plan != nil {
		diags.Append(extractModelFromPlan(ctx, plan, &model)...)
	} else {
		diags.Append(extractModelFromState(ctx, state, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &model, diags
}

// extractModelFromPlan extracts model from Terraform plan
func extractModelFromPlan(ctx context.Context, plan *tfsdk.Plan, model *WebsiteMonitoringConfigModel) diag.Diagnostics {
	diags := plan.Get(ctx, model)
	if diags.HasError() {
		diags.AddError(WebsiteMonitoringConfigErrRetrievingPlan, "Failed to extract model from plan")
	}
	return diags
}

// extractModelFromState extracts model from Terraform state
func extractModelFromState(ctx context.Context, state *tfsdk.State, model *WebsiteMonitoringConfigModel) diag.Diagnostics {
	diags := state.Get(ctx, model)
	if diags.HasError() {
		diags.AddError(WebsiteMonitoringConfigErrRetrievingState, "Failed to extract model from state")
	}
	return diags
}

// mapModelToAPIObject converts Terraform model to API object
func mapModelToAPIObject(model *WebsiteMonitoringConfigModel) (*restapi.WebsiteMonitoringConfig, error) {
	if model == nil {
		return nil, fmt.Errorf("model cannot be nil")
	}

	if err := validateModelFields(model); err != nil {
		return nil, err
	}

	return createAPIObjectFromModel(model), nil
}

// validateModelFields validates required fields in the model
func validateModelFields(model *WebsiteMonitoringConfigModel) error {
	if model.Name.IsNull() || model.Name.IsUnknown() {
		return fmt.Errorf("name field is required and cannot be null or unknown")
	}

	nameValue := model.Name.ValueString()
	if len(nameValue) < WebsiteMonitoringConfigMinNameLength || len(nameValue) > WebsiteMonitoringConfigMaxNameLength {
		return fmt.Errorf("name length must be between %d and %d characters",
			WebsiteMonitoringConfigMinNameLength, WebsiteMonitoringConfigMaxNameLength)
	}

	return nil
}

// createAPIObjectFromModel creates API object from validated model
func createAPIObjectFromModel(model *WebsiteMonitoringConfigModel) *restapi.WebsiteMonitoringConfig {
	return &restapi.WebsiteMonitoringConfig{
		ID:   extractStringValue(model.ID),
		Name: extractStringValue(model.Name),
	}
}

// extractStringValue safely extracts string value from types.String
func extractStringValue(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}

// mapAPIObjectToModel converts API object to Terraform model
func mapAPIObjectToModel(apiObject *restapi.WebsiteMonitoringConfig) *WebsiteMonitoringConfigModel {
	return &WebsiteMonitoringConfigModel{
		ID:      createStringValue(apiObject.ID),
		Name:    createStringValue(apiObject.Name),
		AppName: createStringValue(apiObject.AppName),
	}
}

// createStringValue creates types.String from string value
func createStringValue(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

// updateTerraformState updates the Terraform state with the model
func updateTerraformState(ctx context.Context, state *tfsdk.State, model *WebsiteMonitoringConfigModel) diag.Diagnostics {
	diags := state.Set(ctx, model)
	if diags.HasError() {
		diags.AddError(WebsiteMonitoringConfigErrUpdatingState, "Failed to update Terraform state")
	}
	return diags
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *websiteMonitoringConfigResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: resourcehandle.CreateStateUpgraderForVersion(1),
	}
}
