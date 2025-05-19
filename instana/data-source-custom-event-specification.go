package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewCustomEventDataSource creates a new DataSource for custom Event specifications
func NewCustomEventSpecificationDataSource() DataSource {
	return &customEventSpecificationDataSource{}
}

const DataSourceCustomEventSpec = "instana_custom_event_spec"

type customEventSpecificationDataSource struct{}

// CreateResource creates the terraform Resource for the data source for Instana custom event specifications
func (ds *customEventSpecificationDataSource) CreateResource() *schema.Resource {
	return &schema.Resource{
		ReadContext: ds.read,
		Schema: map[string]*schema.Schema{
			CustomEventSpecificationFieldName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the custom event specification.",
			},
			CustomEventSpecificationFieldDescription: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the custom event specification.",
			},
			CustomEventSpecificationFieldEntityType: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The entity type for which the custom event specification is created.",
			},
			CustomEventSpecificationFieldTriggering: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if an incident is triggered the custom event or not.",
			},
			CustomEventSpecificationFieldEnabled: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the custom event is enabled or not.",
			},
			CustomEventSpecificationFieldQuery: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dynamic focus query for the custom event specification.",
			},
			CustomEventSpecificationFieldExpirationTime: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The expiration time (grace period) to wait before the issue is closed.",
			},
		},
	}
}

func (ds *customEventSpecificationDataSource) read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	name := d.Get(CustomEventSpecificationFieldName).(string)
	entityType := d.Get(CustomEventSpecificationFieldEntityType).(string)

	data, err := instanaAPI.CustomEventSpecifications().GetAll()
	if err != nil {
		return diag.FromErr(err)
	}

	customEvent, err := ds.findCustomEventByNameAndEntityType(name, entityType, data)
	if err != nil {
		return diag.FromErr(err)
	}

	err = ds.updateState(d, customEvent)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func (ds *customEventSpecificationDataSource) findCustomEventByNameAndEntityType(name string, entityType string, data *[]*restapi.CustomEventSpecification) (*restapi.CustomEventSpecification, error) {
	for _, customEvent := range *data {
		if customEvent.Name == name && customEvent.EntityType == entityType {
			return customEvent, nil
		}
	}
	return nil, fmt.Errorf("no custom event specification found for name '%s' and entity type '%s'", name, entityType)
}

func (ds *customEventSpecificationDataSource) updateState(d *schema.ResourceData, customEventSpec *restapi.CustomEventSpecification) error {
	d.SetId(customEventSpec.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		CustomEventSpecificationFieldDescription:    customEventSpec.Description,
		CustomEventSpecificationFieldTriggering:     customEventSpec.Triggering,
		CustomEventSpecificationFieldEnabled:        customEventSpec.Enabled,
		CustomEventSpecificationFieldQuery:          customEventSpec.Query,
		CustomEventSpecificationFieldExpirationTime: customEventSpec.ExpirationTime,
	})
}
