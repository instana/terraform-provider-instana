package mobileappconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"

	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
)

// NewMobileAppConfigResourceHandle creates the resource handle for Mobile App Configurations
func NewMobileAppConfigResourceHandle() resourcehandle.ResourceHandle[*api.MobileAppConfig] {
	return &mobileAppConfigResource{
		metaData: createResourceMetaData(),
	}
}

// mobileAppConfigResource implements the resource framework for mobile app configurations
type mobileAppConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

// MetaData returns the resource metadata
func (r *mobileAppConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

// GetRestResource returns the REST resource for the API
func (r *mobileAppConfigResource) GetRestResource(api client.InstanaAPI) rest.RestResource[*api.MobileAppConfig] {
	return api.MobileAppConfig()
}

// SetComputedFields sets computed fields in the plan (no computed fields to set for this resource)
func (r *mobileAppConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return diag.Diagnostics{}
}

// MapStateToDataObject maps Terraform state/plan to API data object
func (r *mobileAppConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*api.MobileAppConfig, diag.Diagnostics) {
	if err := validateMapStateToDataObjectInputs(ctx, plan, state); err != nil {
		return nil, diag.Diagnostics{err}
	}

	model, diags := extractModelFromPlanOrState(ctx, plan, state)
	if diags.HasError() {
		return nil, diags
	}

	apiObject, err := mapModelToAPIObject(model)
	if err != nil {
		diags.AddError(MobileAppConfigErrMappingToAPI, err.Error())
		return nil, diags
	}

	return apiObject, diags
}

// UpdateState updates Terraform state with data from API object
func (r *mobileAppConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *api.MobileAppConfig) diag.Diagnostics {
	if err := validateUpdateStateInputs(ctx, state, apiObject); err != nil {
		return diag.Diagnostics{err}
	}

	model := mapAPIObjectToModel(apiObject)

	return updateTerraformState(ctx, state, model)
}

// createResourceMetaData creates the resource metadata with schema definition
func createResourceMetaData() resourcehandle.ResourceMetaData {
	return resourcehandle.ResourceMetaData{
		ResourceName:  ResourceInstanaMobileAppConfig,
		Schema:        createResourceSchema(),
		SchemaVersion: MobileAppConfigSchemaVersion,
	}
}

// createResourceSchema creates the Terraform schema for the resource
func createResourceSchema() schema.Schema {
	return schema.Schema{
		Description: MobileAppConfigDescResource,
		Attributes:  createSchemaAttributes(),
	}
}

// createSchemaAttributes creates the schema attributes map
func createSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		MobileAppConfigFieldID:   createIDAttribute(),
		MobileAppConfigFieldName: createNameAttribute(),
	}
}

// createIDAttribute creates the ID schema attribute
func createIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:      true,
		Description:   MobileAppConfigDescID,
		PlanModifiers: createIDPlanModifiers(),
	}
}

// createNameAttribute creates the name schema attribute
func createNameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: MobileAppConfigDescName,
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
		return diag.NewErrorDiagnostic(MobileAppConfigErrInvalidInput, MobileAppConfigErrNilContext)
	}

	if plan == nil && state == nil {
		return diag.NewErrorDiagnostic(MobileAppConfigErrInvalidInput, "Both plan and state cannot be nil")
	}

	return nil
}

// validateUpdateStateInputs validates inputs for UpdateState method
func validateUpdateStateInputs(ctx context.Context, state *tfsdk.State, apiObject *api.MobileAppConfig) diag.Diagnostic {
	if ctx == nil {
		return diag.NewErrorDiagnostic(MobileAppConfigErrInvalidInput, MobileAppConfigErrNilContext)
	}

	if state == nil {
		return diag.NewErrorDiagnostic(MobileAppConfigErrInvalidInput, MobileAppConfigErrNilState)
	}

	if apiObject == nil {
		return diag.NewErrorDiagnostic(MobileAppConfigErrInvalidInput, MobileAppConfigErrNilAPIObject)
	}

	return nil
}

// extractModelFromPlanOrState extracts the model from plan or state
func extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*MobileAppConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model MobileAppConfigModel

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
func extractModelFromPlan(ctx context.Context, plan *tfsdk.Plan, model *MobileAppConfigModel) diag.Diagnostics {
	diags := plan.Get(ctx, model)
	if diags.HasError() {
		diags.AddError(MobileAppConfigErrRetrievingPlan, "Failed to extract model from plan")
	}
	return diags
}

// extractModelFromState extracts model from Terraform state
func extractModelFromState(ctx context.Context, state *tfsdk.State, model *MobileAppConfigModel) diag.Diagnostics {
	diags := state.Get(ctx, model)
	if diags.HasError() {
		diags.AddError(MobileAppConfigErrRetrievingState, "Failed to extract model from state")
	}
	return diags
}

// mapModelToAPIObject converts Terraform model to API object
func mapModelToAPIObject(model *MobileAppConfigModel) (*api.MobileAppConfig, error) {
	if model == nil {
		return nil, fmt.Errorf("model cannot be nil")
	}

	if err := validateModelFields(model); err != nil {
		return nil, err
	}

	return createAPIObjectFromModel(model), nil
}

// validateModelFields validates required fields in the model
func validateModelFields(model *MobileAppConfigModel) error {
	if model.Name.IsNull() || model.Name.IsUnknown() {
		return fmt.Errorf("name field is required and cannot be null or unknown")
	}

	nameValue := model.Name.ValueString()
	if len(nameValue) < MobileAppConfigMinNameLength || len(nameValue) > MobileAppConfigMaxNameLength {
		return fmt.Errorf("name length must be between %d and %d characters",
			MobileAppConfigMinNameLength, MobileAppConfigMaxNameLength)
	}

	return nil
}

// createAPIObjectFromModel creates API object from validated model
func createAPIObjectFromModel(model *MobileAppConfigModel) *api.MobileAppConfig {
	return &api.MobileAppConfig{
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
func mapAPIObjectToModel(apiObject *api.MobileAppConfig) *MobileAppConfigModel {
	return &MobileAppConfigModel{
		ID:   createStringValue(apiObject.ID),
		Name: createStringValue(apiObject.Name),
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
func updateTerraformState(ctx context.Context, state *tfsdk.State, model *MobileAppConfigModel) diag.Diagnostics {
	diags := state.Set(ctx, model)
	if diags.HasError() {
		diags.AddError(MobileAppConfigErrUpdatingState, "Failed to update Terraform state")
	}
	return diags
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *mobileAppConfigResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: resourcehandle.CreateStateUpgraderForVersion(1),
	}
}
