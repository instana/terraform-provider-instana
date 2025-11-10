package instana

// import (
// 	"context"
// 	"fmt"
// 	"strings"
// 	"github.com/gessnerfl/terraform-provider-instana/tfutils"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

// 	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// const DataSourceAutomationAction = "instana_automation_action"

// // NewAutomationActionDataSource creates a new DataSource for Automation Action
// func NewAutomationActionDataSource() DataSource {
// 	return &automationActionDataSource{}
// }

// type automationActionDataSource struct{}

// // CreateResource creates the terraform Resource for the data source for Instana automation actions
// func (ds *automationActionDataSource) CreateResource() *schema.Resource {
// 	return &schema.Resource {
// 		ReadContext: ds.read,
// 		Schema: map[string]*schema.Schema {
// 			AutomationActionFieldName: {
// 				Type:        schema.TypeString,
// 				Required:    true,
// 				Description: "The name of the automation action",
// 			},
// 			AutomationActionFieldDescription: {
// 				Type:        schema.TypeString,
// 				Computed:    true,
// 				Description: "The description of the automation action.",
// 			},
// 			AutomationActionFieldType: {
// 				Type:        schema.TypeString,
// 				Required:    true,
// 				Description: "The type of the automation action.",
// 			},
// 			AutomationActionFieldTags: {
// 				Type:     schema.TypeList,
// 				Elem: 	  &schema.Schema {
// 					Type: schema.TypeString,
// 				},
// 				Computed:    true,
// 				Description: "The tags of the automation action.",
// 			},
// 		},
// 	}
// }

// func (ds *automationActionDataSource) read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	providerMeta := meta.(*ProviderMeta)
// 	instanaAPI := providerMeta.InstanaAPI

// 	name := d.Get(AutomationActionFieldName).(string)
// 	actionType := d.Get(AutomationActionFieldType).(string)

// 	data, err := instanaAPI.AutomationActions().GetAll()
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	action, err := ds.findActionByNameAndType(name, actionType, data)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	err = ds.updateState(d, action)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return nil
// }

// func (ds *automationActionDataSource) findActionByNameAndType(name string, actionType string, data *[]*restapi.AutomationAction) (*restapi.AutomationAction, error) {
// 	for _, action := range *data {
// 		if action.Name == name && strings.EqualFold(action.Type, actionType) {
// 			return action, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("no automation action found for name '%s' and type '%s'", name, actionType)
// }

// func (ds *automationActionDataSource) updateState(d *schema.ResourceData, automationAction *restapi.AutomationAction) error {
// 	d.SetId(automationAction.ID)
// 	return tfutils.UpdateState(d, map[string]interface{} {
// 		AutomationActionFieldName:        	automationAction.Name,
// 	    AutomationActionFieldDescription:  	automationAction.Description,
// 		AutomationActionFieldType: 			automationAction.Type,
// 		AutomationActionFieldTags: 			automationAction.Tags,
// 	})
// }
