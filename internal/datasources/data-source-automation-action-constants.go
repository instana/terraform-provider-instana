package datasources

// Data source name constants
const (
	// DataSourceInstanaAutomationActionFramework the name of the terraform-provider-instana data source to read automation actions
	DataSourceInstanaAutomationActionFramework = "automation_action"
)

// Field name constants for automation action datasource (field constants are shared with automation-action.go)
const (
	// AutomationActionFieldID constant value for the schema field id
	AutomationActionFieldID = "id"
)

// Description constants for automation action fields
const (
	// AutomationActionDescDataSource description for the data source
	AutomationActionDescDataSource = "Data source for an Instana automation action. Automation actions are used to execute scripts or HTTP requests."
	// AutomationActionDescID description for the ID field
	AutomationActionDescID = "The ID of the automation action."
	// AutomationActionDescName description for the name field
	AutomationActionDescName = "The name of the automation action."
	// AutomationActionDescDescription description for the description field
	AutomationActionDescDescription = "The description of the automation action."
	// AutomationActionDescType description for the type field
	AutomationActionDescType = "The type of the automation action."
	// AutomationActionDescTags description for the tags field
	AutomationActionDescTags = "The tags of the automation action."
)

// Error message constants
const (
	// AutomationActionErrUnexpectedConfigureType error message for unexpected configure type
	AutomationActionErrUnexpectedConfigureType = "Unexpected Data Source Configure Type"
	// AutomationActionErrUnexpectedConfigureTypeDetail error message detail for unexpected configure type
	AutomationActionErrUnexpectedConfigureTypeDetail = "Expected *restapi.ProviderMeta, got: %T. Please report this issue to the provider developers."
	// AutomationActionErrReadingActions error message for reading automation actions
	AutomationActionErrReadingActions = "Error reading automation actions"
	// AutomationActionErrReadingActionsDetail error message detail for reading automation actions
	AutomationActionErrReadingActionsDetail = "Could not read automation actions: %s"
	// AutomationActionErrNotFound error message for automation action not found
	AutomationActionErrNotFound = "Automation action not found"
	// AutomationActionErrNotFoundDetail error message detail for automation action not found
	AutomationActionErrNotFoundDetail = "No automation action found with name '%s' and type '%s'"
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
