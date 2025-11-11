package automationaction

// ResourceInstanaAutomationActionFramework the name of the terraform-provider-instana resource to manage automation actions
const ResourceInstanaAutomationActionFramework = "automation_action"

// Resource description
const AutomationActionDescResource = "This resource manages automation actions in Instana."

// Note: AutomationActionDescID, AutomationActionDescName, AutomationActionDescDescription,
// and AutomationActionDescTags are already defined in data-source-automation-action-constants.go

// Field descriptions - Script configuration
const AutomationActionDescScript = "Script configuration for the automation action."
const AutomationActionDescScriptContent = "The script content."
const AutomationActionDescScriptInterpreter = "The script interpreter."
const AutomationActionDescScriptTimeout = "The timeout for script execution in seconds."
const AutomationActionDescScriptSource = "The source of the script."

// Field descriptions - HTTP configuration
const AutomationActionDescHttp = "HTTP configuration for the automation action."
const AutomationActionDescHttpHost = "The URL of the HTTP request."
const AutomationActionDescHttpMethod = "The HTTP method."
const AutomationActionDescHttpBody = "The body of the HTTP request."
const AutomationActionDescHttpHeaders = "The headers of the HTTP request."
const AutomationActionDescHttpIgnoreCertErrors = "Whether to ignore certificate errors for the request."
const AutomationActionDescHttpTimeout = "The timeout for HTTP request execution in seconds."
const AutomationActionDescHttpLanguage = "The language for the HTTP request."
const AutomationActionDescHttpContentType = "The content type for the HTTP request."

// Field descriptions - HTTP Auth configuration
const AutomationActionDescHttpAuth = "Authentication configuration for the HTTP request."
const AutomationActionDescHttpAuthBasic = "Basic authentication configuration."
const AutomationActionDescHttpAuthBasicUsername = "The username for basic authentication."
const AutomationActionDescHttpAuthBasicPassword = "The password for basic authentication."
const AutomationActionDescHttpAuthToken = "Bearer token authentication configuration."
const AutomationActionDescHttpAuthBearerToken = "The bearer token for authentication."
const AutomationActionDescHttpAuthApiKey = "API key authentication configuration."
const AutomationActionDescHttpAuthApiKeyKey = "The API key name."
const AutomationActionDescHttpAuthApiKeyValue = "The API key value."
const AutomationActionDescHttpAuthApiKeyLocation = "Where to add the API key (header or query)."

// Field descriptions - Manual action configuration
const AutomationActionDescManual = "Manual action configuration."
const AutomationActionDescManualContent = "The content for the manual action."

// Field descriptions - Jira action configuration
const AutomationActionDescJira = "Jira action configuration."
const AutomationActionDescJiraProject = "The Jira project."
const AutomationActionDescJiraOperation = "The Jira operation type."
const AutomationActionDescJiraIssueType = "The Jira issue type."
const AutomationActionDescJiraDescription = "The Jira issue description."
const AutomationActionDescJiraAssignee = "The Jira issue assignee."
const AutomationActionDescJiraTitle = "The Jira issue title."
const AutomationActionDescJiraLabels = "The Jira issue labels."
const AutomationActionDescJiraComment = "The Jira issue comment."

// Field descriptions - GitHub action configuration
const AutomationActionDescGitHub = "GitHub action configuration."
const AutomationActionDescGitHubOwner = "The GitHub repository owner."
const AutomationActionDescGitHubRepo = "The GitHub repository name."
const AutomationActionDescGitHubTitle = "The GitHub issue title."
const AutomationActionDescGitHubBody = "The GitHub issue body."
const AutomationActionDescGitHubOperation = "The GitHub operation type."
const AutomationActionDescGitHubAssignees = "The GitHub issue assignees."
const AutomationActionDescGitHubLabels = "The GitHub issue labels."
const AutomationActionDescGitHubComment = "The GitHub issue comment."

// Field descriptions - DocLink action configuration
const AutomationActionDescDocLink = "Documentation link action configuration."
const AutomationActionDescDocLinkUrl = "The URL to the documentation."

// Field descriptions - GitLab action configuration
const AutomationActionDescGitLab = "GitLab action configuration."
const AutomationActionDescGitLabProjectId = "The GitLab project ID."
const AutomationActionDescGitLabTitle = "The GitLab issue title."
const AutomationActionDescGitLabDescription = "The GitLab issue description."
const AutomationActionDescGitLabOperation = "The GitLab operation type."
const AutomationActionDescGitLabLabels = "The GitLab issue labels."
const AutomationActionDescGitLabIssueType = "The GitLab issue type."
const AutomationActionDescGitLabComment = "The GitLab issue comment."

// Field descriptions - Ansible action configuration
const AutomationActionDescAnsible = "Ansible action configuration."
const AutomationActionDescAnsibleWorkflowId = "The Ansible workflow ID."
const AutomationActionDescAnsiblePlaybookId = "The Ansible playbook ID."
const AutomationActionDescAnsiblePlaybookFileName = "The Ansible playbook file name."
const AutomationActionDescAnsibleUrl = "The Ansible URL."
const AutomationActionDescAnsibleHostId = "The host ID from which this action is created."

// Field descriptions - Input parameters
const AutomationActionDescInputParameter = "Input parameters for the automation action."
const AutomationActionDescParameterName = "The name of the parameter."
const AutomationActionDescParameterDescription = "The description of the parameter."
const AutomationActionDescParameterLabel = "The label of the parameter."
const AutomationActionDescParameterRequired = "Whether the parameter is required."
const AutomationActionDescParameterHidden = "Whether the parameter is hidden."
const AutomationActionDescParameterType = "The type of the parameter."
const AutomationActionDescParameterValue = "The value of the parameter."

// Error messages
const AutomationActionErrMappingTags = "Error mapping tags"
const AutomationActionErrTagNotString = "Tag at index %d is not a string"
const AutomationActionErrTagsFormat = "Tags are not in the expected format"
const AutomationActionErrUnmarshalHeaders = "Error unmarshaling HTTP headers"
const AutomationActionErrUnmarshalHeadersFailed = "Failed to unmarshal HTTP headers: %s"
const AutomationActionErrMarshalBasicAuth = "Error marshaling basic auth"
const AutomationActionErrMarshalBasicAuthFailed = "Failed to marshal basic auth: %s"
const AutomationActionErrMarshalBearerToken = "Error marshaling bearer token"
const AutomationActionErrMarshalBearerTokenFailed = "Failed to marshal bearer token: %s"
const AutomationActionErrMarshalApiKey = "Error marshaling API key"
const AutomationActionErrMarshalApiKeyFailed = "Failed to marshal API key: %s"
const AutomationActionErrMarshalNoAuth = "Error marshaling no auth"
const AutomationActionErrMarshalNoAuthFailed = "Failed to marshal no auth: %s"
const AutomationActionErrMarshalHeaders = "Error marshaling HTTP headers"
const AutomationActionErrMarshalHeadersFailed = "Failed to marshal HTTP headers: %s"
const AutomationActionErrInvalidConfig = "Invalid action configuration"
const AutomationActionErrInvalidConfigMsg = "One of script, http, manual, jira, github, doclink, gitlab, or ansible configuration must be provided"

// Field names for internal use
const AutomationActionFieldSource = "source"
const AutomationActionFieldLanguage = "language"
const AutomationActionFieldContentType = "content_type"
const AutomationActionFieldAuth = "auth"
const AutomationActionFieldBasicAuth = "basic_auth"
const AutomationActionFieldUsername = "username"
const AutomationActionFieldPassword = "password"
const AutomationActionFieldToken = "token"
const AutomationActionFieldBearerToken = "bearer_token"
const AutomationActionFieldApiKey = "api_key"
const AutomationActionFieldKey = "key"
const AutomationActionFieldValue = "value"
const AutomationActionFieldKeyLocation = "key_location"
const AutomationActionFieldManual = "manual"
const AutomationActionFieldJira = "jira"
const AutomationActionFieldProject = "project"
const AutomationActionFieldOperation = "operation"
const AutomationActionFieldIssueType = "issue_type"
const AutomationActionFieldAssignee = "assignee"
const AutomationActionFieldTitle = "title"
const AutomationActionFieldLabels = "labels"
const AutomationActionFieldComment = "comment"
const AutomationActionFieldGitHub = "github"
const AutomationActionFieldOwner = "owner"
const AutomationActionFieldRepo = "repo"
const AutomationActionFieldAssignees = "assignees"
const AutomationActionFieldDocLink = "doc_link"
const AutomationActionFieldUrl = "url"
const AutomationActionFieldGitLab = "gitlab"
const AutomationActionFieldProjectId = "project_id"
const AutomationActionFieldAnsible = "ansible"
const AutomationActionFieldWorkflowId = "workflow_id"
const AutomationActionFieldPlaybookId = "playbook_id"
const AutomationActionFieldPlaybookFileName = "playbook_file_name"
const AutomationActionFieldHostId = "host_id"

// Field descriptions for internal fields
const AutomationActionDescFieldSource = "The source of the script"
const AutomationActionDescFieldLanguage = "The language for the HTTP request"
const AutomationActionDescFieldContentType = "The content type for the HTTP request"
const AutomationActionDescFieldAuthen = "Authentication for the HTTPS request"
const AutomationActionDescFieldUrl = "URL to remediation documentation"

// Internal field names used in API
const AutomationActionAPIFieldAuthen = "authen"
const AutomationActionAPIFieldTicketActionType = "ticketActionType"
const AutomationActionAPIFieldSummary = "summary"
const AutomationActionAPIFieldTicketType = "ticketType"
const AutomationActionAPIFieldProjectId = "projectId"
const AutomationActionAPIFieldWorkflowId = "workflowId"
const AutomationActionAPIFieldAnsibleUrl = "ansibleUrl"
const AutomationActionAPIFieldHostId = "hostId"
const AutomationActionAPIFieldPlaybookId = "playbookId"
const AutomationActionAPIFieldPlaybookFileName = "playbookFileName"

// Action type constants (already defined in lines 25-26, but adding descriptions)
const AutomationActionTypeManual = "MANUAL"
const AutomationActionTypeJira = "JIRA"
const AutomationActionTypeGitHub = "GITHUB"
const AutomationActionTypeDocLink = "DOC_LINK"
const AutomationActionTypeGitLab = "GITLAB"
const AutomationActionTypeAnsible = "ANSIBLE"

// Auth type constants
const AutomationActionAuthTypeBasicAuth = "basicAuth"
const AutomationActionAuthTypeBearerToken = "bearerToken"
const AutomationActionAuthTypeApiKey = "apiKey"
const AutomationActionAuthTypeNoAuth = "noAuth"

// Field descriptions for API fields
const AutomationActionDescAPIFieldContent = "Content for manual action"
const AutomationActionDescAPIFieldProject = "jira project"
const AutomationActionDescAPIFieldTicketActionType = "jira ticket type"
const AutomationActionDescAPIFieldIssueType = "jira issue type"
const AutomationActionDescAPIFieldBody = "jira issue description"
const AutomationActionDescAPIFieldAssignee = "jira issue assignee"
const AutomationActionDescAPIFieldSummary = "jira issue summary"
const AutomationActionDescAPIFieldLabels = "jira issue labels"
const AutomationActionDescAPIFieldComment = "jira issue comment"
const AutomationActionDescAPIFieldOwner = "github issue owner/repo"
const AutomationActionDescAPIFieldRepo = "github issue repo"
const AutomationActionDescAPIFieldTitle = "github issue title"
const AutomationActionDescAPIFieldGitHubBody = "github issue body"
const AutomationActionDescAPIFieldTicketType = "github issue type"
const AutomationActionDescAPIFieldAssignees = "github issue assignees"
const AutomationActionDescAPIFieldGitHubLabels = "github issue labels"
const AutomationActionDescAPIFieldGitHubComment = "github issue comment"
const AutomationActionDescAPIFieldProjectId = "gitlab projectId"
const AutomationActionDescAPIFieldGitLabTitle = "gitlab issue title"
const AutomationActionDescAPIFieldGitLabBody = "gitlab issue description"
const AutomationActionDescAPIFieldGitLabTicketActionType = "gitlab ticket type"
const AutomationActionDescAPIFieldGitLabLabels = "github issue labels"
const AutomationActionDescAPIFieldGitLabIssueType = "gitlab issue type"
const AutomationActionDescAPIFieldGitLabComment = "gitlab issue comment"
const AutomationActionDescAPIFieldWorkflowId = "The workflow ID"
const AutomationActionDescAPIFieldAnsibleUrl = "The ansible url"
const AutomationActionDescAPIFieldHostId = "The host ID from which this action is created"
const AutomationActionDescAPIFieldPlaybookId = "The playbook ID"
const AutomationActionDescAPIFieldPlaybookFileName = "The playbook filename"

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
