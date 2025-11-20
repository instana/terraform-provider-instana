# Automation Action Resource


Management of Automation Actions.

API Documentation: <https://instana.github.io/openapi/#tag/Action-Catalog>

---
## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.
 
## Migration Guide (v5 to v6)
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
    content     = filebase64("test.sh")
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
    content     = filebase64("test.sh")
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
  name        = "Hello World Script"
  description = "Simple hello world script"
  tags        = ["test", "demo"]
  
  script = {
    interpreter = "bash"
    content     = filebase64("${path.module}/scripts/hello.sh")
    timeout     = "10"
  }
}
```

### Script Action with Parameters

```hcl
resource "instana_automation_action" "restart_service" {
  name        = "Restart Service"
  description = "Restart a specified service"
  tags        = ["operations", "service-management"]
  
  script = {
    interpreter = "bash"
    content     = filebase64("${path.module}/scripts/restart_service.sh")
    timeout     = "30"
  }
  
  input_parameter = [{
    name        = "service_name"
    label       = "Service Name"
    description = "Name of the service to restart"
    type        = "static"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "wait_time"
    label       = "Wait Time (seconds)"
    description = "Time to wait after restart"
    type        = "static"
    required    = false
    hidden      = false
    value       = "5"
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
    host                      = "https://api.example.com/health"
    method                    = "GET"
    ignore_certificate_errors = false
    timeout                   = "10"
  }
}
```

### HTTP Action with Authentication

```hcl
resource "instana_automation_action" "api_call_with_auth" {
  name        = "API Call with Bearer Token"
  description = "Make authenticated API call"
  tags        = ["api", "integration"]
  
  http = {
    host   = "https://api.example.com/v1/resource"
    method = "POST"
    
    headers = {
      "Content-Type" = "application/json"
      "Accept"       = "application/json"
    }
    
    body = jsonencode({
      action = "process"
      data   = "@@input_data@@"
    })
    
    auth = {
      token = {
        bearer_token = "@@api_token@@"
      }
    }
    
    timeout = "30"
  }
  
  input_parameter = [{
    name        = "api_token"
    label       = "API Token"
    description = "Bearer token for authentication"
    type        = "vault"
    required    = true
    hidden      = true
    value       = ""
  }, {
    name        = "input_data"
    label       = "Input Data"
    description = "Data to process"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }]
}
```

### HTTP Action with Basic Auth

```hcl
resource "instana_automation_action" "basic_auth_call" {
  name        = "API with Basic Auth"
  description = "Call API with basic authentication"
  tags        = ["api", "legacy"]
  
  http = {
    host   = "https://legacy-api.example.com/endpoint"
    method = "GET"
    
    auth = {
      basic_auth = {
        username = "@@username@@"
        password = "@@password@@"
      }
    }
    
    timeout = "15"
  }
  
  input_parameter = [{
    name        = "username"
    label       = "Username"
    description = "API username"
    type        = "static"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "password"
    label       = "Password"
    description = "API password"
    type        = "vault"
    required    = true
    hidden      = true
    value       = ""
  }]
}
```

### HTTP Action with API Key

```hcl
resource "instana_automation_action" "api_key_call" {
  name        = "API with API Key"
  description = "Call API with API key authentication"
  tags        = ["api", "integration"]
  
  http = {
    host   = "https://api.example.com/v2/data"
    method = "POST"
    
    headers = {
      "Content-Type" = "application/json"
    }
    
    body = jsonencode({
      query = "@@query@@"
    })
    
    auth = {
      api_key = {
        key          = "X-API-Key"
        value        = "@@api_key@@"
        key_location = "header"
      }
    }
    
    timeout = "20"
  }
  
  input_parameter = [{
    name        = "api_key"
    label       = "API Key"
    description = "API key for authentication"
    type        = "vault"
    required    = true
    hidden      = true
    value       = ""
  }, {
    name        = "query"
    label       = "Query"
    description = "Query to execute"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }]
}
```

### Manual Action

```hcl
resource "instana_automation_action" "manual_intervention" {
  name        = "Manual Intervention Required"
  description = "Notify team for manual intervention"
  tags        = ["manual", "escalation"]
  
  manual = {
    content = <<-EOT
      # Manual Intervention Required
      
      ## Issue Details
      - Alert: @@alert_name@@
      - Severity: @@severity@@
      - Time: @@timestamp@@
      
      ## Required Actions
      1. Review the alert details
      2. Check system logs
      3. Verify service health
      4. Take corrective action
      5. Document resolution
      
      ## Escalation
      If issue persists, escalate to on-call engineer.
    EOT
  }
  
  input_parameter = [{
    name        = "alert_name"
    label       = "Alert Name"
    description = "Name of the triggered alert"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "severity"
    label       = "Severity"
    description = "Alert severity level"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }]
}
```

### JIRA Integration

```hcl
resource "instana_automation_action" "create_jira_ticket" {
  name        = "Create JIRA Ticket"
  description = "Create JIRA ticket for incident"
  tags        = ["jira", "ticketing"]
  
  jira = {
    project     = "OPS"
    operation   = "create"
    issue_type  = "Incident"
    title       = "Instana Alert: @@alert_name@@"
    description = "Alert triggered at @@timestamp@@\n\nDetails: @@alert_details@@"
    assignee    = "@@assignee@@"
    labels      = "instana,automated,@@severity@@"
  }
  
  input_parameter = [{
    name        = "alert_name"
    label       = "Alert Name"
    description = "Name of the alert"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "alert_details"
    label       = "Alert Details"
    description = "Detailed alert information"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "assignee"
    label       = "Assignee"
    description = "JIRA user to assign ticket to"
    type        = "static"
    required    = false
    hidden      = false
    value       = "unassigned"
  }, {
    name        = "severity"
    label       = "Severity"
    description = "Alert severity"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }]
}
```

### GitHub Integration

```hcl
resource "instana_automation_action" "create_github_issue" {
  name        = "Create GitHub Issue"
  description = "Create GitHub issue for tracking"
  tags        = ["github", "issue-tracking"]
  
  github = {
    owner      = "myorg"
    repo       = "infrastructure"
    operation  = "create"
    title      = "[Instana] @@alert_name@@"
    body       = "## Alert Details\n\n@@alert_details@@\n\n## Timestamp\n@@timestamp@@"
    assignees  = "@@assignees@@"
    labels     = "instana,automated,@@severity@@"
  }
  
  input_parameter = [{
    name        = "alert_name"
    label       = "Alert Name"
    description = "Name of the alert"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "alert_details"
    label       = "Alert Details"
    description = "Detailed alert information"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "assignees"
    label       = "Assignees"
    description = "GitHub users to assign (comma-separated)"
    type        = "static"
    required    = false
    hidden      = false
    value       = ""
  }, {
    name        = "severity"
    label       = "Severity"
    description = "Alert severity"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }]
}
```

### GitLab Integration

```hcl
resource "instana_automation_action" "create_gitlab_issue" {
  name        = "Create GitLab Issue"
  description = "Create GitLab issue for incident tracking"
  tags        = ["gitlab", "issue-tracking"]
  
  gitlab = {
    project_id  = "12345"
    operation   = "create"
    title       = "[Instana Alert] @@alert_name@@"
    description = "Alert triggered: @@alert_details@@"
    labels      = "instana,automated,@@severity@@"
    issue_type  = "incident"
  }
  
  input_parameter = [{
    name        = "alert_name"
    label       = "Alert Name"
    description = "Name of the alert"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "alert_details"
    label       = "Alert Details"
    description = "Detailed alert information"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "severity"
    label       = "Severity"
    description = "Alert severity"
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
  name        = "Runbook Reference"
  description = "Link to relevant runbook"
  tags        = ["documentation", "runbook"]
  
  doc_link = {
    url = "https://wiki.example.com/runbooks/@@runbook_id@@"
  }
  
  input_parameter = [{
    name        = "runbook_id"
    label       = "Runbook ID"
    description = "ID of the relevant runbook"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }]
}
```

### Ansible Integration

```hcl
resource "instana_automation_action" "ansible_playbook" {
  name        = "Run Ansible Playbook"
  description = "Execute Ansible playbook for remediation"
  tags        = ["ansible", "automation"]
  
  ansible = {
    workflow_id        = "@@workflow_id@@"
    playbook_id        = "@@playbook_id@@"
    playbook_file_name = "remediation.yml"
    url                = "https://ansible.example.com"
    host_id            = "@@target_host@@"
  }
  
  input_parameter = [{
    name        = "workflow_id"
    label       = "Workflow ID"
    description = "Ansible workflow identifier"
    type        = "static"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "playbook_id"
    label       = "Playbook ID"
    description = "Ansible playbook identifier"
    type        = "static"
    required    = true
    hidden      = false
    value       = ""
  }, {
    name        = "target_host"
    label       = "Target Host"
    description = "Host to run playbook on"
    type        = "dynamic"
    required    = true
    hidden      = false
    value       = ""
  }]
}
```

### Complex HTTP Action with Multiple Parameters

```hcl
resource "instana_automation_action" "complex_api_call" {
  name        = "Complex API Integration"
  description = "Advanced API call with multiple parameters"
  tags        = ["api", "integration", "advanced"]
  
  http = {
    host   = "https://api.example.com/v1/incidents"
    method = "POST"
    
    headers = {
      "Content-Type"  = "application/json"
      "Accept"        = "application/json"
      "X-Request-ID"  = "@@request_id@@"
      "X-Correlation" = "@@correlation_id@@"
    }
    
    body = jsonencode({
      incident = {
        title       = "@@incident_title@@"
        description = "@@incident_description@@"
        severity    = "@@severity@@"
        source      = "instana"
        metadata = {
          alert_id    = "@@alert_id@@"
          timestamp   = "@@timestamp@@"
          environment = "@@environment@@"
        }
      }
    })
    
    auth = {
      token = {
        bearer_token = "@@api_token@@"
      }
    }
    
    ignore_certificate_errors = false
    timeout                   = "30"
    language                  = "json"
    content_type              = "application/json"
  }
  
  input_parameter = [
    {
      name        = "api_token"
      label       = "API Token"
      description = "Bearer token for API authentication"
      type        = "vault"
      required    = true
      hidden      = true
      value       = ""
    },
    {
      name        = "incident_title"
      label       = "Incident Title"
      description = "Title of the incident"
      type        = "dynamic"
      required    = true
      hidden      = false
      value       = ""
    },
    {
      name        = "incident_description"
      label       = "Incident Description"
      description = "Detailed description of the incident"
      type        = "dynamic"
      required    = true
      hidden      = false
      value       = ""
    },
    {
      name        = "severity"
      label       = "Severity"
      description = "Incident severity level"
      type        = "dynamic"
      required    = true
      hidden      = false
      value       = ""
    },
    {
      name        = "alert_id"
      label       = "Alert ID"
      description = "Instana alert identifier"
      type        = "dynamic"
      required    = true
      hidden      = false
      value       = ""
    },
    {
      name        = "timestamp"
      label       = "Timestamp"
      description = "Alert timestamp"
      type        = "dynamic"
      required    = true
      hidden      = false
      value       = ""
    },
    {
      name        = "environment"
      label       = "Environment"
      description = "Environment name"
      type        = "static"
      required    = true
      hidden      = false
      value       = "production"
    },
    {
      name        = "request_id"
      label       = "Request ID"
      description = "Unique request identifier"
      type        = "dynamic"
      required    = false
      hidden      = false
      value       = ""
    },
    {
      name        = "correlation_id"
      label       = "Correlation ID"
      description = "Correlation identifier for tracking"
      type        = "dynamic"
      required    = false
      hidden      = false
      value       = ""
    }
  ]
}
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