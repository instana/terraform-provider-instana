# Automation Action Resource


Management of Automation Actions.

API Documentation: <https://instana.github.io/openapi/#tag/Action-Catalog>

---
## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.
 
## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- Action type blocks (`script`, `http`, `manual`, etc.) now use **object syntax** with `= { }`
- `input_parameter` now uses **list syntax** with `= [{ }]`
- Enhanced validation for action types and parameters
- Better error messages for configuration issues
- Improved state management

#### OLD (v5.x) Syntax:

```hcl
resource "instana_automation_action" "example" {
  name        = "Hello world"
  description = "Script action"
  tags        = ["test", "hello"]
  
  script {
    interpreter = "bash"
    content     = "test"
    timeout     = "10"
  }
  
  input_parameter {
    name        = "test"
    label       = "test parameter"
    description = "parameter for test"
    type        = "static"
    required    = true
    hidden      = false
    value       = ""
  }
}
```

#### NEW (v6.x) Syntax:

```hcl
resource "instana_automation_action" "example" {
  name        = "Hello world"
  description = "Script action"
  tags        = ["test", "hello"]
  
  script = {
    interpreter = "bash"
    content     = "test"
    timeout     = "10"
  }
  
  input_parameter = [{
    name        = "test"
    label       = "test parameter"
    description = "parameter for test"
    type        = "static"
    required    = true
    hidden      = false
    value       = ""
  }]
}
```

---

## Example Usage

### Script Action - Basic

```hcl
resource "instana_automation_action" "hello_world" {
  name        = "Hello World Script" # Replace with your own value
  description = "Simple hello world script" # Replace with your own value
  tags        = ["test", "demo"] # Replace with your own value
  
  script = {
    interpreter = "bash"
    content     = "test-content" # Replace with your own value
    timeout     = "10"
  }
}
```

### Script Action with Parameters

```hcl
resource "instana_automation_action" "restart_service" {
  name        = "Restart Service" # Replace with your own value
  description = "Restart a specified service" # Replace with your own value
  tags        = ["operations", "service-management"] # Replace with your own value
  
  script = {
    interpreter = "bash"
    content     = "test-content" # Replace with your own value
    timeout     = "30"
  }
  
  input_parameter = [{
    name        = "service_name" # Replace with your own value
    label       = "Service Name" # Replace with your own value
    description = "Name of the service to restart" # Replace with your own value
    type        = "static"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "wait_time" # Replace with your own value
    label       = "Wait Time (seconds)" # Replace with your own value
    description = "Time to wait after restart" # Replace with your own value
    type        = "static"
    required    = false
    hidden      = false
    value       = "5" # Replace with your own value
  }]
}
```

### HTTP Action - Basic

```hcl
resource "instana_automation_action" "health_check" {
  name        = "Service Health Check"
  description = "Check service health endpoint"
  tags        = ["monitoring", "health"]
  
  http = {
    host                      = "https://api.example.com/health" # Replace with your own value
    method                    = "GET"
    ignore_certificate_errors = false
    timeout                   = "10"
  }
}
```

### HTTP Action with Basic Auth

```hcl
resource "instana_automation_action" "basic_auth_call" {
  name        = "API with Basic Auth" # Replace with your own value
  description = "Call API with basic authentication" # Replace with your own value
  tags        = ["api", "legacy"] # Replace with your own value
  
  http = {
    host   = "https://legacy-api.example.com/endpoint" # Replace with your own value
    method = "GET"
    
    auth = {
      basic_auth = {
        username = "@@username@@" # Replace with your own value
        password = "@@password@@" # Replace with your own value
      }
    }
    
    timeout = "15"
  }
  
  input_parameter = [{
    name        = "username" # Replace with your own value
    label       = "Username" # Replace with your own value
    description = "API username" # Replace with your own value
    type        = "static"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "password" # Replace with your own value
    label       = "Password" # Replace with your own value
    description = "API password" # Replace with your own value
    type        = "vault"
    required    = true
    hidden      = true
    value       = ""
  }]
}
```

### Manual Action

```hcl
resource "instana_automation_action" "automation_action_manual" {
  description     = " Detects a rapid increase in the values of the erroneous call rate metric" # Replace with your own value
  manual = {
    content = "content" # Replace with your own value
  }
  name   = "erroneous call rate" # Replace with your own value
  tags   = ["tag"]
}
```

### JIRA Integration

```hcl
resource "instana_automation_action" "Jira_task" {
  description     = "Description"               # Replace with your own value
  jira = {
    description = "This is a task in Jira"      # Replace with your own value
    issue_type  = "Bug"                         # Replace with your own value
    operation   = "open"
    project     = "Project"                     # Replace with your own value
    title       = "Title"                       # Replace with your own value
  }
  manual = null
  name   = "Jira-test" # Replace with your own value
}
```

### GitHub Integration

```hcl
resource "instana_automation_action" "create_github_issue" {
  name        = "Create GitHub Issue" # Replace with your own value
  description = "Create GitHub issue for tracking" # Replace with your own value
  tags        = ["github", "issue-tracking"]
  github = {
    owner      = "myorg"                        # Replace with your own value
    repo       = "infrastructure"               # Replace with your own value
    operation  = "create"
    title      = "[Instana] @@alert_name@@"
    body       = "## Alert Details\n\n@@alert_details@@\n\n## Timestamp\n@@timestamp@@"
    assignees  = "@@assignees@@"
    labels     = "instana,automated,@@severity@@" # Replace with your own value
  }
  input_parameter = [{
    name        = "alert_name" # Replace with your own value
    label       = "Alert Name" # Replace with your own value
    description = "Name of the alert" # Replace with your own value
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = "" # Replace with your own value
  }]
}
```

### GitLab Integration

```hcl
resource "instana_automation_action" "create_gitlab_issue" {
  name        = "Create GitLab Issue" # Replace with your own value
  description = "Create GitLab issue for incident tracking" # Replace with your own value
  tags        = ["gitlab", "issue-tracking"]
  
  gitlab = {
    project_id  = "12345"                       # Replace with your own value
    operation   = "create"
    title       = "[Instana Alert]"             # Replace with your own value
    description = "Alert triggered"             # Replace with your own value
    labels      = "instana,automated"
    issue_type  = "incident"
  }
  
  input_parameter = [{
    name        = "alert_name" # Replace with your own value
    label       = "Alert Name" # Replace with your own value
    description = "Name of the alert" # Replace with your own value
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }]
}
```

### Documentation Link

```hcl
resource "instana_automation_action" "runbook_link" {
  name        = "Runbook Reference" # Replace with your own value
  description = "Link to relevant runbook" # Replace with your own value
  tags        = ["documentation", "runbook"]
  
  doc_link = {
    url = "https://wiki.example.com/runbooks/@@runbook_id@@" # Replace with your own value
  }
  
  input_parameter = [{
    name        = "runbook_id" # Replace with your own value
    label       = "Runbook ID" # Replace with your own value
    description = "ID of the relevant runbook" # Replace with your own value
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }]
}
```

### Ansible Integration

```hcl
resource "instana_automation_action" "tf_b_ansible_1" {
  ansible = {
    host_id            = "<host_id>"            # Replace with your own value
    playbook_file_name = "<file_name>"          # Replace with your own value
    playbook_id        = "<id>"                 # Replace with your own value
    url                = "<url>"                # Replace with your own value
  }
  description = "Test"
  input_parameter = [
    {
      description = "test" # Replace with your own value
      hidden      = false
      label       = "Label"
      name        = "Name" # Replace with your own value
      required    = false
      type        = "static"
      value       = ""
    },
    {
      description = "" # Replace with your own value
      hidden      = false
      label       = "Label" # Replace with your own value
      name        = "Name" # Replace with your own value
      required    = false
      type        = "static"
      value       = ""
    },
  ]
  name   = "Ansible Test"
}

```

## Generating Configuration from Existing Resources

If you have already created a automation action in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your automation action in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_automation_action.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual automation action ID from Instana

3. **Generate the Configuration**: Run the following Terraform command:

```bash
terraform plan -generate-config-out=generated.tf
```

This will:
- Import the existing resource state
- Generate the complete Terraform configuration in `generated.tf`
- Show you what will be imported

4. **Review and Apply**: Review the generated configuration in `generated.tf` and make any necessary adjustments.

   - **To import the existing resource**: Keep the import block and run `terraform apply`. This will import the resource into your Terraform state and link it to the existing Instana resource.
   
   - **To create a new resource**: If you only need the configuration structure as a template, remove the import block from your configuration. Modify the generated configuration as needed, and when you run `terraform apply`, it will create a new resource in Instana instead of importing the existing one.

```bash
terraform apply
```

## Argument Reference

* `name` - Required - The name of the automation action
* `description` - Required - The description of the automation action
* `tags` - Optional - A list of tags for the automation action (list of strings)
* `input_parameter` - Optional - A list of input parameters (list). [Details](#input-parameter-argument-reference)

**Exactly one of the following action type blocks must be provided:**

* `script` - Optional - Script Action Configuration (object). [Details](#script-argument-reference)
* `http` - Optional - HTTP Action Configuration (object). [Details](#http-argument-reference)
* `manual` - Optional - Manual Action Configuration (object). [Details](#manual-argument-reference)
* `jira` - Optional - JIRA Integration Configuration (object). [Details](#jira-argument-reference)
* `github` - Optional - GitHub Integration Configuration (object). [Details](#github-argument-reference)
* `gitlab` - Optional - GitLab Integration Configuration (object). [Details](#gitlab-argument-reference)
* `doc_link` - Optional - Documentation Link Configuration (object). [Details](#doc-link-argument-reference)
* `ansible` - Optional - Ansible Integration Configuration (object). [Details](#ansible-argument-reference)

### Input Parameter Argument Reference

* `name` - Required - The name of the input parameter
* `label` - Required - The label of the input parameter
* `description` - Required - The description of the input parameter
* `type` - Required - The type of the input parameter. Allowed values: `static`, `dynamic`, `vault`
* `required` - Required - Indicates if the input parameter is required (boolean)
* `hidden` - Required - Indicates if the input parameter is hidden (boolean)
* `value` - Required - The default value of the input parameter (can be empty string)

### Script Argument Reference

* `content` - Required - Base64 encoded script content (use `filebase64()` function)
* `interpreter` - Optional - The interpreter for script execution (e.g., `bash`, `python`, `powershell`)
* `timeout` - Optional - The timeout of the automation action in seconds (string)
* `source` - Optional - The source of the script

### HTTP Argument Reference

* `host` - Required - The host URL for the HTTP request
* `method` - Required - The HTTP method. Allowed values: `GET`, `POST`, `PUT`, `DELETE`
* `body` - Optional - The body content for the HTTP request (string)
* `headers` - Optional - The headers of the HTTP request (map of strings)
* `auth` - Optional - Authentication configuration (object). [Details](#auth-argument-reference)
* `ignore_certificate_errors` - Optional - Indicates if the HTTP request ignores certificate errors (boolean)
* `timeout` - Optional - The timeout of the automation action in seconds (string)
* `language` - Optional - The language/format of the request body (e.g., `json`, `xml`)
* `content_type` - Optional - The content type of the request

#### Auth Argument Reference

Exactly one of the following must be configured:

* `basic_auth` - Optional - Basic authentication configuration (object)
  * `username` - Required - Username for basic authentication
  * `password` - Required - Password for basic authentication
* `token` - Optional - Bearer token authentication configuration (object)
  * `bearer_token` - Required - Bearer token for authentication
* `api_key` - Optional - API key authentication configuration (object)
  * `key` - Required - The API key header/parameter name
  * `value` - Required - The API key value
  * `key_location` - Required - Where to place the API key. Allowed values: `header`, `query`

### Manual Argument Reference

* `content` - Required - The content/instructions for manual action (string, supports markdown)

### JIRA Argument Reference

* `project` - Optional - JIRA project key
* `operation` - Optional - Operation to perform (e.g., `create`, `update`)
* `issue_type` - Optional - JIRA issue type (e.g., `Bug`, `Incident`, `Task`)
* `title` - Optional - Issue title/summary
* `description` - Optional - Issue description
* `assignee` - Optional - JIRA username to assign the issue to
* `labels` - Optional - Comma-separated list of labels
* `comment` - Optional - Comment to add to the issue

### GitHub Argument Reference

* `owner` - Optional - GitHub repository owner/organization
* `repo` - Optional - GitHub repository name
* `operation` - Optional - Operation to perform (e.g., `create`, `update`)
* `title` - Optional - Issue title
* `body` - Optional - Issue body/description
* `assignees` - Optional - Comma-separated list of GitHub usernames to assign
* `labels` - Optional - Comma-separated list of labels
* `comment` - Optional - Comment to add to the issue

### GitLab Argument Reference

* `project_id` - Optional - GitLab project ID
* `operation` - Optional - Operation to perform (e.g., `create`, `update`)
* `title` - Optional - Issue title
* `description` - Optional - Issue description
* `labels` - Optional - Comma-separated list of labels
* `issue_type` - Optional - Issue type (e.g., `issue`, `incident`)
* `comment` - Optional - Comment to add to the issue

### Doc Link Argument Reference

* `url` - Required - URL to the documentation/runbook

### Ansible Argument Reference

* `workflow_id` - Optional - Ansible workflow identifier
* `playbook_id` - Optional - Ansible playbook identifier
* `playbook_file_name` - Optional - Name of the playbook file
* `url` - Optional - Ansible Tower/AWX URL
* `host_id` - Optional - Target host identifier

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the automation action

## Import

Automation Actions can be imported using the `id`, e.g.:

```bash
terraform import instana_automation_action.example daa60cb8-41fc-3d4c-868f-3c2aca5831fa
```

## Best Practices

### Script Actions

- Use `filebase64()` to encode script content
- Set appropriate timeouts based on script complexity
- Test scripts independently before automation
- Use input parameters for flexibility

### HTTP Actions

- Always use HTTPS for sensitive data
- Store credentials in vault-type parameters
- Set reasonable timeouts
- Handle authentication properly
- Use appropriate HTTP methods

### Input Parameters

- Use `vault` type for sensitive data (passwords, tokens)
- Use `dynamic` type for runtime values from alerts
- Use `static` type for configuration values
- Mark sensitive parameters as `hidden`
- Provide clear descriptions

### Tags

Use tags for organization and filtering:

```hcl
tags = ["environment:production", "team:platform", "type:remediation"]
```

## Notes

- The resource ID is auto-generated by Instana upon creation
- Only one action type can be configured per resource
- Input parameters support variable substitution using `@@parameter_name@@` syntax
- Script content must be base64 encoded
- HTTP actions support various authentication methods
- Tags help organize and filter actions in the Instana UI
- Timeout values are in seconds and should be strings