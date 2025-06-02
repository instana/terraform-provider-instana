package instana

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	AutomationActionFieldName           = "name"
	AutomationActionFieldDescription    = "description"
	AutomationActionFieldTags           = "tags"
	AutomationActionFieldTimeout        = "timeout"
	AutomationActionFieldType           = "type"
	AutomationActionFieldInputParameter = "input_parameter"

	// script constants
	AutomationActionFieldScript      = "script"
	AutomationActionFieldContent     = "content"
	AutomationActionFieldInterpreter = "interpreter"

	// http constants
	AutomationActionFieldHttp             = "http"
	AutomationActionFieldMethod           = "method"
	AutomationActionFieldHost             = "host"
	AutomationActionFieldHeaders          = "headers"
	AutomationActionFieldBody             = "body"
	AutomationActionFieldIgnoreCertErrors = "ignore_certificate_errors"

	// input parameter constants
	AutomationActionParameterFieldName        = "name"
	AutomationActionParameterFieldLabel       = "label"
	AutomationActionParameterFieldDescription = "description"
	AutomationActionParameterFieldType        = "type"
	AutomationActionParameterFieldValue       = "value"
	AutomationActionParameterFieldRequired    = "required"
	AutomationActionParameterFieldHidden      = "hidden"
)

var supportedActionTypes = []string{
	"script",
	"http",
}

var (
	automationActionTimeoutSchema = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The timeout of the automation action.",
	}
	automationActionScriptSchema = &schema.Schema{
		Type:         schema.TypeList,
		MinItems:     0,
		MaxItems:     1,
		Optional:     true,
		Description:  "The configuration of the script action.",
		ExactlyOneOf: supportedActionTypes,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				AutomationActionFieldContent: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Base64 encoded script content.",
				},
				AutomationActionFieldInterpreter: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The interpreter for script execution.",
				},
				AutomationActionFieldTimeout: automationActionTimeoutSchema,
			},
		},
	}
	automationActionHttpSchema = &schema.Schema{
		Type:         schema.TypeList,
		MinItems:     0,
		MaxItems:     1,
		Optional:     true,
		Description:  "The configuration of the http action.",
		ExactlyOneOf: supportedActionTypes,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				AutomationActionFieldMethod: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The method for http request.",
					ValidateFunc: validation.StringInSlice([]string{"GET", "POST", "PUT", "DELETE"}, false),
				},
				AutomationActionFieldHost: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The host for the http request.",
				},
				AutomationActionFieldIgnoreCertErrors: {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Indicates if the http request ignores the certificate errors.",
				},
				AutomationActionFieldHeaders: {
					Type:        schema.TypeMap,
					Optional:    true,
					Description: "The headers for the http request.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				AutomationActionFieldBody: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The body content for http request.",
				},
				AutomationActionFieldTimeout: automationActionTimeoutSchema,
			},
		},
	}
	automationActionInputParameterSchema = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				AutomationActionParameterFieldName: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the input parameter.",
				},
				AutomationActionParameterFieldLabel: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The label of the input parameter.",
				},
				AutomationActionParameterFieldDescription: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The description of the input parameter.",
				},
				AutomationActionParameterFieldType: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The type of the input parameter.",
					ValidateFunc: validation.StringInSlice([]string{"static", "dynamic", "vault"}, false),
				},
				AutomationActionParameterFieldRequired: {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Indicates if the input parameter is required.",
				},
				AutomationActionParameterFieldHidden: {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Indicates if the input parameter is hidden.",
				},
				AutomationActionParameterFieldValue: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The value of the input parameter.",
				},
			},
		},
		Description: "A list of input parameters.",
	}
)
