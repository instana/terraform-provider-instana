package instana

import (
	"log"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaAutomationAction the name of the terraform-provider-instana resource to manage automation actions
const ResourceInstanaAutomationAction = "instana_automation_action"

// NewAutomationActionResourceHandle creates the resource handle for Automation Actions
func NewAutomationActionResourceHandle() ResourceHandle[*restapi.AutomationAction] {
	return &AutomationActionResource {
		metaData: ResourceMetaData {
			ResourceName: ResourceInstanaAutomationAction,
			Schema: map[string]*schema.Schema {
				AutomationActionFieldName: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The name of the automation action",
				},
				AutomationActionFieldDescription: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The description of the automation action",
				},
				AutomationActionFieldTags: {
					Type:     schema.TypeList,
					Elem: 	  &schema.Schema {
						Type: schema.TypeString,
					},
					Optional: true,
					Description: "The tags of the automation action.",
				},
				AutomationActionFieldTimeout: {
					Type:         schema.TypeInt,
					Optional:     true,
					Description:  "The timeout of the automation action",
				},
				AutomationActionFieldType: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The type of the automation action",
					ValidateFunc: validation.StringInSlice([]string{"SCRIPT", "HTTP"}, false),
				},
				AutomationActionFieldInputParameter: {
					Type: schema.TypeSet,
					Set: func(i interface{}) int {
						return schema.HashString(i.(map[string]interface{})[AutomationActionParameterFieldName])
					},
					Optional: true,
					MinItems: 0,
					Elem: &schema.Resource {
						Schema: map[string]*schema.Schema {
							AutomationActionParameterFieldName: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The label of the input parameter.",
							},
							AutomationActionParameterFieldLabel: {
								Type:     schema.TypeString,
								Required: true,
								Description: "The label of the input parameter.",
							},
							AutomationActionParameterFieldDescription: {
								Type:     schema.TypeString,
								Optional: true,
								Description: "The description of the input parameter.",
							},
							AutomationActionParameterFieldType: {
								Type:     schema.TypeString,
								Required: true,
								Description: "The type of the input parameter.",
								ValidateFunc: validation.StringInSlice([]string{"static", "dynamic"}, false),
							},
							AutomationActionParameterFieldRequired: {
								Type:     schema.TypeBool,
								Required: true,
								Description: "Indicates if the input parameter is required.",
							},
							AutomationActionParameterFieldHidden: {
								Type:     schema.TypeBool,
								Optional: true,
								Description: "Indicates if the input parameter is hidden.",
							},
							AutomationActionParameterFieldSecured: {
								Type:     schema.TypeBool,
								Optional: true,
								Description: "Indicates if the input parameter is secured.",
							},
							AutomationActionParameterFieldValue: {
								Type:     schema.TypeString,
								Required: true,
								Description: "The value of the input parameter.",
							},
						},
					},
					Description: "A list of input parameters.",
				},
			},
			SchemaVersion: 0,
		},
	}
}

type AutomationActionResource struct {
	metaData ResourceMetaData
}

func (r *AutomationActionResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *AutomationActionResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *AutomationActionResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AutomationAction] {
	return api.AutomationActions()
}

func (r *AutomationActionResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *AutomationActionResource) UpdateState(d *schema.ResourceData, automationAction *restapi.AutomationAction) error {
	log.Printf("INFO: UpdateState \n")

	d.SetId(automationAction.ID)
	return tfutils.UpdateState(d, map[string]interface{} {
		AutomationActionFieldName:            		automationAction.Name,
		AutomationActionFieldDescription:           automationAction.Description,
		AutomationActionFieldTags:					automationAction.Tags,
		AutomationActionFieldTimeout:				automationAction.Timeout,
		AutomationActionFieldType:					automationAction.Type,
		AutomationActionFieldInputParameter:		r.mapInputParametersToSchema(automationAction),
	})
}

func (r *AutomationActionResource) mapInputParametersToSchema(action *restapi.AutomationAction) []interface{} {
	result := make([]interface{}, len(action.InputParameters))

	i := 0
	for _, v := range action.InputParameters {
		val := v

		item := make(map[string]interface{})
		item[AutomationActionParameterFieldName] = val.Name
		item[AutomationActionParameterFieldDescription] = val.Description
		item[AutomationActionParameterFieldLabel] = val.Label
		item[AutomationActionParameterFieldHidden] = val.Hidden
		item[AutomationActionParameterFieldSecured] = val.Secured
		item[AutomationActionParameterFieldType] = val.Type
		item[AutomationActionParameterFieldValue] = val.Value
		
		result[i] = item
		i++
	}
	return result
}

func (r *AutomationActionResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.AutomationAction, error) {
	return &restapi.AutomationAction {
		ID:            		d.Id(),
		Name:          		d.Get(AutomationActionFieldName).(string),
		Description:   		d.Get(AutomationActionFieldDescription).(string),
		Type:          		d.Get(AutomationActionFieldType).(string),
		Tags:          		d.Get(AutomationActionFieldTags),
		Timeout:	   		d.Get(AutomationActionFieldTimeout).(int),
		InputParameters:	r.mapInputParametersFromSchema(d),
	}, nil
}

func (r *AutomationActionResource) mapInputParametersFromSchema(d *schema.ResourceData) []restapi.Parameter {
	val, ok := d.GetOk(AutomationActionFieldInputParameter)

	if ok && val != nil {
		schemaParameters := val.(*schema.Set).List()
		result := make([]restapi.Parameter, len(schemaParameters))
		
		i := 0
		for _, v := range schemaParameters {
			param := v.(map[string]interface{})
			
			result[i] = restapi.Parameter {
				Name: 			param[AutomationActionParameterFieldName].(string),
				Description: 	param[AutomationActionParameterFieldDescription].(string),
				Label: 			param[AutomationActionParameterFieldLabel].(string),
				Required: 		param[AutomationActionParameterFieldRequired].(bool),
				Hidden:			param[AutomationActionParameterFieldHidden].(bool),
				Secured: 		param[AutomationActionParameterFieldSecured].(bool),
				Type: 			param[AutomationActionParameterFieldType].(string),
				Value: 			param[AutomationActionParameterFieldValue].(string),
			}
			i++
		}
		return result
	}

	return []restapi.Parameter{}
}
