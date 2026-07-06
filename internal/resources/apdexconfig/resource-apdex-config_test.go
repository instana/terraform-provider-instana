package apdexconfig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/instana/instana-go-client/api"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApdexConfigResourceHandle(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewApdexConfigResourceHandle()

		require.NotNil(t, handle)
		metadata := handle.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaApdexConfig, metadata.ResourceName)
		assert.Equal(t, int64(0), metadata.SchemaVersion)
	})
}

func TestMetaData(t *testing.T) {
	t.Run("should return metadata", func(t *testing.T) {
		resource := &apdexConfigResource{}
		metadata := resource.MetaData()
		require.NotNil(t, metadata)
	})
}

func TestSetComputedFields(t *testing.T) {
	t.Run("should not set any computed fields", func(t *testing.T) {
		resource := &apdexConfigResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaApdexConfig,
				Schema:        NewApdexConfigResourceHandle().MetaData().Schema,
				SchemaVersion: 0,
			},
		}
		ctx := context.Background()

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		diags := resource.SetComputedFields(ctx, plan)
		assert.False(t, diags.HasError())
	})
}

func TestMapStateToDataObject_ApplicationEntity(t *testing.T) {
	resource := &apdexConfigResource{}

	t.Run("should map application entity successfully with all fields", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			ApplicationEntityModel: &ApplicationApdexEntityModel{
				EntityID:         types.StringValue("app-123"),
				Threshold:        types.Int64Value(500),
				BoundaryScope:    types.StringValue("ALL"),
				IncludeInternal:  types.BoolValue(true),
				IncludeSynthetic: types.BoolValue(false),
				FilterExpression: types.StringValue("call.http.status EQUALS 200"),
			},
		}

		result, diags := resource.mapEntityFromState(entityModel)

		assert.False(t, diags.HasError())
		assert.Equal(t, "application", result.Type)
		assert.Equal(t, "app-123", result.EntityID)
		assert.Equal(t, 500, result.Threshold)
		require.NotNil(t, result.BoundaryScope)
		assert.Equal(t, "ALL", *result.BoundaryScope)
		require.NotNil(t, result.IncludeInternal)
		assert.True(t, *result.IncludeInternal)
		require.NotNil(t, result.IncludeSynthetic)
		assert.False(t, *result.IncludeSynthetic)
		require.NotNil(t, result.TagFilter)
	})

	t.Run("should map application entity with defaults for optional boolean fields", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			ApplicationEntityModel: &ApplicationApdexEntityModel{
				EntityID:      types.StringValue("app-456"),
				Threshold:     types.Int64Value(1000),
				BoundaryScope: types.StringValue("INBOUND"),
				// Optional fields not set
			},
		}

		result, diags := resource.mapEntityFromState(entityModel)

		assert.False(t, diags.HasError())
		assert.Equal(t, "application", result.Type)
		require.NotNil(t, result.IncludeInternal)
		assert.False(t, *result.IncludeInternal, "should default to false")
		require.NotNil(t, result.IncludeSynthetic)
		assert.False(t, *result.IncludeSynthetic, "should default to false")
	})

	t.Run("should return error when both plan and state are nil", func(t *testing.T) {
		ctx := context.Background()
		result, diags := resource.MapStateToDataObject(ctx, nil, nil)

		assert.Nil(t, result)
		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrMappingState)
	})

	t.Run("should return error when application entity missing entity_id", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			ApplicationEntityModel: &ApplicationApdexEntityModel{
				EntityID:      types.StringNull(),
				Threshold:     types.Int64Value(500),
				BoundaryScope: types.StringValue("ALL"),
			},
		}

		_, diags := resource.validateAndMapApplicationEntity(entityModel.ApplicationEntityModel)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrMissingEntity)
	})

	t.Run("should return error when application entity missing threshold", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			ApplicationEntityModel: &ApplicationApdexEntityModel{
				EntityID:      types.StringValue("app-123"),
				Threshold:     types.Int64Null(),
				BoundaryScope: types.StringValue("ALL"),
			},
		}

		_, diags := resource.validateAndMapApplicationEntity(entityModel.ApplicationEntityModel)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrMissingEntity)
	})

	t.Run("should return error when application entity missing boundary_scope", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			ApplicationEntityModel: &ApplicationApdexEntityModel{
				EntityID:      types.StringValue("app-123"),
				Threshold:     types.Int64Value(500),
				BoundaryScope: types.StringNull(),
			},
		}

		_, diags := resource.validateAndMapApplicationEntity(entityModel.ApplicationEntityModel)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrMissingEntity)
	})

	t.Run("should return error when filter expression is invalid", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			ApplicationEntityModel: &ApplicationApdexEntityModel{
				EntityID:         types.StringValue("app-123"),
				Threshold:        types.Int64Value(500),
				BoundaryScope:    types.StringValue("ALL"),
				FilterExpression: types.StringValue("invalid((expression"),
			},
		}

		_, diags := resource.validateAndMapApplicationEntity(entityModel.ApplicationEntityModel)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrParsingFilterExpression)
	})
}

func TestMapStateToDataObject_WebsiteEntity(t *testing.T) {
	resource := &apdexConfigResource{}

	t.Run("should map website entity successfully with all fields", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			WebsiteEntityModel: &WebsiteApdexEntityModel{
				EntityID:         types.StringValue("website-123"),
				Threshold:        types.Int64Value(1000),
				BeaconType:       types.StringValue("pageLoad"),
				FilterExpression: types.StringValue("beacon.page.name CONTAINS 'checkout'"),
			},
		}

		result, diags := resource.mapEntityFromState(entityModel)

		assert.False(t, diags.HasError())
		assert.Equal(t, "website", result.Type)
		assert.Equal(t, "website-123", result.EntityID)
		assert.Equal(t, 1000, result.Threshold)
		require.NotNil(t, result.BeaconType)
		assert.Equal(t, "pageLoad", *result.BeaconType)
		require.NotNil(t, result.TagFilter)
	})

	t.Run("should map website entity without filter expression", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			WebsiteEntityModel: &WebsiteApdexEntityModel{
				EntityID:   types.StringValue("website-456"),
				Threshold:  types.Int64Value(2000),
				BeaconType: types.StringValue("httpRequest"),
			},
		}

		result, diags := resource.mapEntityFromState(entityModel)

		assert.False(t, diags.HasError())
		assert.Equal(t, "website", result.Type)
		assert.NotNil(t, result.TagFilter)
		assert.Equal(t, TagFilterTypeExpression, string(result.TagFilter.Type))
	})

	t.Run("should return error when website entity missing entity_id", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			WebsiteEntityModel: &WebsiteApdexEntityModel{
				EntityID:   types.StringNull(),
				Threshold:  types.Int64Value(1000),
				BeaconType: types.StringValue("pageLoad"),
			},
		}

		_, diags := resource.validateAndMapWebsiteEntity(entityModel.WebsiteEntityModel)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrMissingEntity)
	})

	t.Run("should return error when website entity missing threshold", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			WebsiteEntityModel: &WebsiteApdexEntityModel{
				EntityID:   types.StringValue("website-123"),
				Threshold:  types.Int64Null(),
				BeaconType: types.StringValue("pageLoad"),
			},
		}

		_, diags := resource.validateAndMapWebsiteEntity(entityModel.WebsiteEntityModel)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrMissingEntity)
	})

	t.Run("should return error when website entity missing beacon_type", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			WebsiteEntityModel: &WebsiteApdexEntityModel{
				EntityID:   types.StringValue("website-123"),
				Threshold:  types.Int64Value(1000),
				BeaconType: types.StringNull(),
			},
		}

		_, diags := resource.validateAndMapWebsiteEntity(entityModel.WebsiteEntityModel)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrMissingEntity)
	})

	t.Run("should return error when filter expression is invalid", func(t *testing.T) {
		entityModel := &ApdexEntityModel{
			WebsiteEntityModel: &WebsiteApdexEntityModel{
				EntityID:         types.StringValue("website-123"),
				Threshold:        types.Int64Value(1000),
				BeaconType:       types.StringValue("pageLoad"),
				FilterExpression: types.StringValue("invalid((expression"),
			},
		}

		_, diags := resource.validateAndMapWebsiteEntity(entityModel.WebsiteEntityModel)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrParsingFilterExpression)
	})
}

func TestMapEntityFromState_Validation(t *testing.T) {
	resource := &apdexConfigResource{}

	t.Run("should return error when no entity is configured", func(t *testing.T) {
		entityModel := &ApdexEntityModel{}

		_, diags := resource.mapEntityFromState(entityModel)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrMissingEntity)
	})

	t.Run("should return error when entity model is nil", func(t *testing.T) {
		_, diags := resource.mapEntityFromState(nil)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrMissingEntity)
	})

	t.Run("should process first entity when both are configured", func(t *testing.T) {
		// Note: Terraform schema should prevent this, but implementation handles it gracefully
		entityModel := &ApdexEntityModel{
			ApplicationEntityModel: &ApplicationApdexEntityModel{
				EntityID:      types.StringValue("app-123"),
				Threshold:     types.Int64Value(500),
				BoundaryScope: types.StringValue("ALL"),
			},
			WebsiteEntityModel: &WebsiteApdexEntityModel{
				EntityID:   types.StringValue("website-123"),
				Threshold:  types.Int64Value(1000),
				BeaconType: types.StringValue("pageLoad"),
			},
		}

		result, diags := resource.mapEntityFromState(entityModel)

		// First entity (application) takes precedence
		assert.False(t, diags.HasError())
		assert.Equal(t, "application", result.Type)
		assert.Equal(t, "app-123", result.EntityID)
	})
}

func TestMapTagsFromState(t *testing.T) {
	resource := &apdexConfigResource{}

	t.Run("should map tags successfully", func(t *testing.T) {
		tags, _ := types.SetValueFrom(context.Background(), types.StringType, []string{"tag1", "tag2", "tag3"})

		result := resource.mapTagsFromState(tags)

		assert.Len(t, result, 3)
		assert.Contains(t, result, "tag1")
		assert.Contains(t, result, "tag2")
		assert.Contains(t, result, "tag3")
	})

	t.Run("should return empty slice for null tags", func(t *testing.T) {
		tags := types.SetNull(types.StringType)

		result := resource.mapTagsFromState(tags)

		assert.Empty(t, result)
	})

	t.Run("should return empty slice for unknown tags", func(t *testing.T) {
		tags := types.SetUnknown(types.StringType)

		result := resource.mapTagsFromState(tags)

		assert.Empty(t, result)
	})
}

func TestMapRbacTagsToState(t *testing.T) {
	resource := &apdexConfigResource{}

	t.Run("should map RBAC tags to state successfully", func(t *testing.T) {
		rbacTags := []api.RbacTag{
			{
				DisplayName: "Team A",
				ID:          "team-a-id",
			},
			{
				DisplayName: "Team B",
				ID:          "team-b-id",
			},
		}

		result := resource.mapRbacTagsToState(rbacTags)

		assert.NotNil(t, result)
		assert.False(t, result.IsNull())
		assert.False(t, result.IsUnknown())

		var tagModels []RbacTagModel
		result.ElementsAs(context.Background(), &tagModels, false)
		assert.Equal(t, 2, len(tagModels))
		assert.Equal(t, "Team A", tagModels[0].DisplayName.ValueString())
		assert.Equal(t, "team-a-id", tagModels[0].ID.ValueString())
		assert.Equal(t, "Team B", tagModels[1].DisplayName.ValueString())
		assert.Equal(t, "team-b-id", tagModels[1].ID.ValueString())
	})

	t.Run("should return empty list for empty RBAC tags", func(t *testing.T) {
		rbacTags := []api.RbacTag{}

		result := resource.mapRbacTagsToState(rbacTags)

		assert.NotNil(t, result)
		assert.False(t, result.IsNull())
		var tagModels []RbacTagModel
		result.ElementsAs(context.Background(), &tagModels, false)
		assert.Equal(t, 0, len(tagModels))
	})

	t.Run("should return empty list for nil RBAC tags", func(t *testing.T) {
		result := resource.mapRbacTagsToState(nil)

		assert.NotNil(t, result)
		assert.False(t, result.IsNull())
		var tagModels []RbacTagModel
		result.ElementsAs(context.Background(), &tagModels, false)
		assert.Equal(t, 0, len(tagModels))
	})
}

func TestMapFilterExpressionToEntity(t *testing.T) {
	resource := &apdexConfigResource{}

	t.Run("should return default filter for empty expression", func(t *testing.T) {
		result, diags := resource.mapFilterExpressionToEntity(types.StringValue(""))

		assert.False(t, diags.HasError())
		require.NotNil(t, result)
	})

	t.Run("should parse valid filter expression", func(t *testing.T) {
		result, diags := resource.mapFilterExpressionToEntity(types.StringValue("call.http.status EQUALS 200"))

		assert.False(t, diags.HasError())
		require.NotNil(t, result)
	})

	t.Run("should parse complex filter expression with AND", func(t *testing.T) {
		result, diags := resource.mapFilterExpressionToEntity(types.StringValue("call.http.status EQUALS 200 AND call.http.method EQUALS 'GET'"))

		assert.False(t, diags.HasError())
		require.NotNil(t, result)
	})

	t.Run("should return error for invalid filter expression", func(t *testing.T) {
		_, diags := resource.mapFilterExpressionToEntity(types.StringValue("invalid((expression"))

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrParsingFilterExpression)
	})

	t.Run("should return error for malformed filter expression", func(t *testing.T) {
		_, diags := resource.mapFilterExpressionToEntity(types.StringValue("call.http.status EQUALS"))

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrParsingFilterExpression)
	})
}

func TestMapEntityToState(t *testing.T) {
	resource := &apdexConfigResource{}

	t.Run("should map application entity to state", func(t *testing.T) {
		boundaryScope := "ALL"
		includeInternal := true
		includeSynthetic := false

		entity := api.ApdexEntity{
			Type:             "application",
			EntityID:         "app-123",
			Threshold:        500,
			BoundaryScope:    &boundaryScope,
			IncludeInternal:  &includeInternal,
			IncludeSynthetic: &includeSynthetic,
		}

		result, diags := resource.mapEntityToState(entity)

		assert.False(t, diags.HasError())
		require.NotNil(t, result.ApplicationEntityModel)
		assert.Nil(t, result.WebsiteEntityModel)
		assert.Equal(t, "app-123", result.ApplicationEntityModel.EntityID.ValueString())
		assert.Equal(t, int64(500), result.ApplicationEntityModel.Threshold.ValueInt64())
		assert.Equal(t, "ALL", result.ApplicationEntityModel.BoundaryScope.ValueString())
		assert.True(t, result.ApplicationEntityModel.IncludeInternal.ValueBool())
		assert.False(t, result.ApplicationEntityModel.IncludeSynthetic.ValueBool())
	})

	t.Run("should map website entity to state", func(t *testing.T) {
		beaconType := "pageLoad"

		entity := api.ApdexEntity{
			Type:       "website",
			EntityID:   "website-456",
			Threshold:  1000,
			BeaconType: &beaconType,
		}

		result, diags := resource.mapEntityToState(entity)

		assert.False(t, diags.HasError())
		require.NotNil(t, result.WebsiteEntityModel)
		assert.Nil(t, result.ApplicationEntityModel)
		assert.Equal(t, "website-456", result.WebsiteEntityModel.EntityID.ValueString())
		assert.Equal(t, int64(1000), result.WebsiteEntityModel.Threshold.ValueInt64())
		assert.Equal(t, "pageLoad", result.WebsiteEntityModel.BeaconType.ValueString())
	})

	t.Run("should return error for unsupported entity type", func(t *testing.T) {
		entity := api.ApdexEntity{
			Type:      "unsupported",
			EntityID:  "test-123",
			Threshold: 500,
		}

		_, diags := resource.mapEntityToState(entity)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), ApdexConfigErrMappingEntityToState)
	})
}
