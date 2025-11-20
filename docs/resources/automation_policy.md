# Automation Policy Resource

Management of Automation Policies.

API Documentation: <https://instana.github.io/openapi/#tag/Policies>

---
## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block struture.

## Migration Guide (v5 to v6)

### Syntax Changes Overview

**OLD Syntax (SDK v2):**
```hcl
resource "instana_automation_policy" "example" {
  name        = "Hello world"
  description = "Sample policy"
  tags        = ["test", "hello"]
  
  trigger {
    id   = "r6mpl4BPRhm_77Vn"
    type = "customEvent"
  }
  
  type_configuration {
    name = "manual"
    
    action {
      action_id = "action-id"
      agent_id  = "agent-id"
      
      input_parameters = {
        test1 = "value1"
      }
    }
    
    condition {
      query = "entity.agent.capability:action-script"
    }
  }
}
```

**NEW Syntax (Plugin Framework):**
```hcl
resource "instana_automation_policy" "example" {
  name        = "Hello world"
  description = "Sample policy"
  tags        = ["test", "hello"]
  
  trigger = {
    id   = "r6mpl4BPRhm_77Vn"
    type = "customEvent"
  }
  
  type_configuration = [{
    name = "manual"
    
    action = [{
      action = {
        id          = "action-id"
        name        = "Action Name"
        description = "Action description"
        
        # Action type configuration (script, http, etc.)
        script = {
          content     = filebase64("script.sh")
          interpreter = "bash"
        }
      }
      agent_id = "agent-id"
    }]
    
    condition = {
      query = "entity.agent.capability:action-script"
    }
  }]
}
```

**Key Changes:**
- `trigger` now uses **object syntax** with `= { }`
- `type_configuration` now uses **list syntax** with `= [{ }]`
- `action` within type_configuration uses **list syntax** with `= [{ }]`
- `condition` now uses **object syntax** with `= { }`
- Actions now require full action definition (not just action_id)
- Enhanced validation and better error messages
- Improved state management

---


## Example Usage

### Manual Policy with Script Action

```hcl
resource "instana_automation_policy" "manual_restart" {
  name        = "Manual Service Restart"
  description = "Manual policy to restart services"
  tags        = ["manual", "restart", "operations"]
  
  trigger = {
    id   = instana_custom_event_specification.service_down.id
    type = "customEvent"
  }
  
  type_configuration = [{
    name = "manual"
    
    action = [{
      action = {
        id          = instana_automation_action.restart_service.id
        name        = "Restart Service"
        description = "Restart the affected service"
        
        script = {
          content     = filebase64("${path.module}/scripts/restart.sh")
          interpreter = "bash"
          timeout     = "30"
        }
        
        input_parameter = [{
          name        = "service_name"
          label       = "Service Name"
          description = "Name of service to restart"
          type        = "dynamic"
          required    = true
          hidden      = false
          value       = ""
        }]
      }
      agent_id = "00:00:0a:ff:fe:0b:21:cc"
    }]
  }]
}
```

### Automatic Policy with Condition

```hcl
resource "instana_automation_policy" "auto_remediation" {
  name        = "Automatic Service Remediation"
  description = "Automatically restart services on failure"
  tags        = ["automatic", "remediation", "production"]
  
  trigger = {
    id   = instana_custom_event_specification.service_error.id
    type = "customEvent"
  }
  
  type_configuration = [{
    name = "automatic"
    
    condition = {
      query = "entity.agent.capability:action-script AND entity.zone:production"
    }
    
    action = [{
      action = {
        id          = instana_automation_action.restart_service.id
        name        = "Restart Service"
        description = "Automatically restart failed service"
        
        script = {
          content     = filebase64("${path.module}/scripts/auto_restart.sh")
          interpreter = "bash"
          timeout     = "30"
        }
      }
      agent_id = "00:00:0a:ff:fe:0b:21:cc"
    }]
  }]
}
```

### Policy with HTTP Action

```hcl
resource "instana_automation_policy" "webhook_notification" {
  name        = "Webhook Notification Policy"
  description = "Send webhook notification on alert"
  tags        = ["notification", "webhook"]
  
  trigger = {
    id   = instana_application_alert_config.critical_alert.id
    type = "applicationAlert"
  }
  
  type_configuration = [{
    name = "automatic"
    
    condition = {
      query = "entity.type:service"
    }
    
    action = [{
      action = {
        id          = instana_automation_action.webhook_call.id
        name        = "Send Webhook"
        description = "Send notification via webhook"
        
        http = {
          host   = "https://hooks.example.com/instana"
          method = "POST"
          
          headers = {
            "Content-Type" = "application/json"
          }
          
          body = jsonencode({
            alert_name = "@@alert_name@@"
            severity   = "@@severity@@"
            timestamp  = "@@timestamp@@"
          })
          
          auth = {
            token = {
              bearer_token = "@@webhook_token@@"
            }
          }
          
          timeout = "15"
        }
        
        input_parameter = [{
          name        = "webhook_token"
          label       = "Webhook Token"
          description = "Authentication token"
          type        = "vault"
          required    = true
          hidden      = true
          value       = ""
        }]
      }
    }]
  }]
}
```

### Policy with Multiple Actions

```hcl
resource "instana_automation_policy" "multi_action_policy" {
  name        = "Multi-Step Remediation"
  description = "Execute multiple remediation steps"
  tags        = ["multi-step", "remediation"]
  
  trigger = {
    id   = instana_custom_event_specification.critical_issue.id
    type = "customEvent"
  }
  
  type_configuration = [{
    name = "automatic"
    
    condition = {
      query = "entity.zone:production AND entity.type:service"
    }
    
    action = [
      {
        action = {
          id          = instana_automation_action.collect_logs.id
          name        = "Collect Logs"
          description = "Collect diagnostic logs"
          
          script = {
            content     = filebase64("${path.module}/scripts/collect_logs.sh")
            interpreter = "bash"
            timeout     = "60"
          }
        }
        agent_id = "00:00:0a:ff:fe:0b:21:cc"
      },
      {
        action = {
          id          = instana_automation_action.restart_service.id
          name        = "Restart Service"
          description = "Restart the affected service"
          
          script = {
            content     = filebase64("${path.module}/scripts/restart.sh")
            interpreter = "bash"
            timeout     = "30"
          }
        }
        agent_id = "00:00:0a:ff:fe:0b:21:cc"
      },
      {
        action = {
          id          = instana_automation_action.notify_team.id
          name        = "Notify Team"
          description = "Send notification to team"
          
          http = {
            host   = "https://api.example.com/notify"
            method = "POST"
            
            body = jsonencode({
              message = "Service restarted automatically"
            })
            
            timeout = "10"
          }
        }
      }
    ]
  }]
}
```

### Policy with JIRA Integration

```hcl
resource "instana_automation_policy" "create_jira_ticket" {
  name        = "Create JIRA Ticket on Alert"
  description = "Automatically create JIRA ticket for incidents"
  tags        = ["jira", "ticketing", "automatic"]
  
  trigger = {
    id   = instana_application_alert_config.production_alert.id
    type = "applicationAlert"
  }
  
  type_configuration = [{
    name = "automatic"
    
    condition = {
      query = "entity.zone:production"
    }
    
    action = [{
      action = {
        id          = instana_automation_action.jira_ticket.id
        name        = "Create JIRA Ticket"
        description = "Create incident ticket in JIRA"
        
        jira = {
          project     = "OPS"
          operation   = "create"
          issue_type  = "Incident"
          title       = "Instana Alert: @@alert_name@@"
          description = "Alert Details: @@alert_details@@"
          assignee    = "ops-team"
          labels      = "instana,automated,production"
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
          name        = "alert_details"
          label       = "Alert Details"
          description = "Detailed alert information"
          type        = "dynamic"
          required    = true
          hidden      = false
          value       = ""
        }]
      }
    }]
  }]
}
```

### Policy with Scheduled Trigger

```hcl
resource "instana_automation_policy" "scheduled_maintenance" {
  name        = "Scheduled Maintenance"
  description = "Run scheduled maintenance tasks"
  tags        = ["scheduled", "maintenance"]
  
  trigger = {
    id          = "scheduled-trigger-id"
    type        = "scheduled"
    name        = "Weekly Maintenance"
    description = "Weekly maintenance window"
    
    scheduling = {
      start_time     = 1640995200000  # Unix timestamp in milliseconds
      duration       = 2
      duration_unit  = "HOUR"
      recurrent      = true
      recurrent_rule = "FREQ=WEEKLY;BYDAY=SU;BYHOUR=2"
    }
  }
  
  type_configuration = [{
    name = "automatic"
    
    action = [{
      action = {
        id          = instana_automation_action.maintenance_script.id
        name        = "Run Maintenance"
        description = "Execute maintenance tasks"
        
        script = {
          content     = filebase64("${path.module}/scripts/maintenance.sh")
          interpreter = "bash"
          timeout     = "300"
        }
      }
      agent_id = "00:00:0a:ff:fe:0b:21:cc"
    }]
  }]
}
```

### Policy with Manual Action

```hcl
resource "instana_automation_policy" "manual_intervention" {
  name        = "Manual Intervention Required"
  description = "Provide runbook for manual intervention"
  tags        = ["manual", "runbook"]
  
  trigger = {
    id   = instana_custom_event_specification.complex_issue.id
    type = "customEvent"
  }
  
  type_configuration = [{
    name = "manual"
    
    action = [{
      action = {
        id          = instana_automation_action.runbook.id
        name        = "Intervention Runbook"
        description = "Manual intervention instructions"
        
        manual = {
          content = <<-EOT
            # Manual Intervention Required
            
            ## Issue Details
            - Alert: @@alert_name@@
            - Severity: @@severity@@
            - Timestamp: @@timestamp@@
            
            ## Investigation Steps
            1. Check service logs: `kubectl logs -n production service-name`
            2. Verify database connectivity
            3. Check external dependencies
            4. Review recent deployments
            
            ## Remediation Steps
            1. If database issue: Restart database connection pool
            2. If external dependency: Check API status
            3. If deployment issue: Consider rollback
            
            ## Escalation
            If issue persists after 30 minutes, escalate to on-call engineer.
          EOT
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
          name        = "severity"
          label       = "Severity"
          description = "Alert severity"
          type        = "dynamic"
          required    = true
          hidden      = false
          value       = ""
        }, {
          name        = "timestamp"
          label       = "Timestamp"
          description = "Alert timestamp"
          type        = "dynamic"
          required    = true
          hidden      = false
          value       = ""
        }]
      }
      agent_id = "00:00:0a:ff:fe:0b:21:cc"
    }]
  }]
}
```

### Comprehensive Policy Example

```hcl
resource "instana_automation_policy" "comprehensive_policy" {
  name        = "Production Incident Response"
  description = "Comprehensive incident response automation"
  tags        = ["production", "incident-response", "critical"]
  
  trigger = {
    id   = instana_application_alert_config.critical_production_alert.id
    type = "applicationAlert"
  }
  
  type_configuration = [{
    name = "automatic"
    
    condition = {
      query = join(" AND ", [
        "entity.zone:production",
        "entity.type:service",
        "entity.agent.capability:action-script"
      ])
    }
    
    action = [
      {
        action = {
          id          = instana_automation_action.collect_diagnostics.id
          name        = "Collect Diagnostics"
          description = "Collect system diagnostics"
          
          script = {
            content     = filebase64("${path.module}/scripts/diagnostics.sh")
            interpreter = "bash"
            timeout     = "60"
          }
          
          input_parameter = [{
            name        = "service_name"
            label       = "Service Name"
            description = "Name of the affected service"
            type        = "dynamic"
            required    = true
            hidden      = false
            value       = ""
          }]
        }
        agent_id = "00:00:0a:ff:fe:0b:21:cc"
      },
      {
        action = {
          id          = instana_automation_action.create_incident.id
          name        = "Create Incident"
          description = "Create incident in ticketing system"
          
          http = {
            host   = "https://api.example.com/incidents"
            method = "POST"
            
            headers = {
              "Content-Type" = "application/json"
            }
            
            body = jsonencode({
              title       = "@@alert_name@@"
              description = "@@alert_details@@"
              severity    = "@@severity@@"
              source      = "instana"
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
            description = "API authentication token"
            type        = "vault"
            required    = true
            hidden      = true
            value       = ""
          }, {
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
            description = "Alert severity level"
            type        = "dynamic"
            required    = true
            hidden      = false
            value       = ""
          }]
        }
      },
      {
        action = {
          id          = instana_automation_action.notify_oncall.id
          name        = "Notify On-Call"
          description = "Notify on-call engineer"
          
          http = {
            host   = "https://api.pagerduty.com/incidents"
            method = "POST"
            
            headers = {
              "Content-Type"  = "application/json"
              "Authorization" = "Token token=@@pagerduty_token@@"
            }
            
            body = jsonencode({
              incident = {
                type    = "incident"
                title   = "@@alert_name@@"
                service = {
                  id   = "@@service_id@@"
                  type = "service_reference"
                }
                urgency = "high"
              }
            })
            
            timeout = "15"
          }
          
          input_parameter = [{
            name        = "pagerduty_token"
            label       = "PagerDuty Token"
            description = "PagerDuty API token"
            type        = "vault"
            required    = true
            hidden      = true
            value       = ""
          }, {
            name        = "service_id"
            label       = "Service ID"
            description = "PagerDuty service ID"
            type        = "static"
            required    = true
            hidden      = false
            value       = ""
          }]
        }
      }
    ]
  }]
}
```

## Argument Reference

* `name` - Required - The name of the automation policy
* `description` - Required - The description of the automation policy
* `tags` - Optional - A list of tags for the automation policy (list of strings)
* `trigger` - Required - The trigger for the automation policy (object). [Details](#trigger-argument-reference)
* `type_configuration` - Required - A list of configurations with the list of actions to run and the mode (automatic or manual) in which the policy is run (list). [Details](#type-configuration-argument-reference)

### Trigger Argument Reference

* `id` - Required - Trigger (Instana event or Smart Alert) identifier
* `type` - Required - Instana event or Smart Alert type. Allowed values: `customEvent`, `applicationAlert`, `scheduled`, etc.
* `name` - Optional - Name of the trigger
* `description` - Optional - Description of the trigger
* `scheduling` - Optional - Scheduling configuration for scheduled triggers (object). [Details](#scheduling-argument-reference)

#### Scheduling Argument Reference

* `start_time` - Optional - Start time in Unix timestamp (milliseconds)
* `duration` - Optional - Duration of the scheduled window
* `duration_unit` - Optional - Unit of duration. Allowed values: `MINUTE`, `HOUR`, `DAY`
* `recurrent` - Optional - Whether the schedule is recurrent (boolean)
* `recurrent_rule` - Optional - Recurrence rule in iCalendar RRULE format (e.g., `FREQ=WEEKLY;BYDAY=SU`)

### Type Configuration Argument Reference

* `name` - Required - The policy type. Allowed values: `manual`, `automatic`
* `condition` - Optional - The condition that selects the list of entities on which the policy is run. Only for automatic policy type (object). [Details](#condition-argument-reference)
* `action` - Required - The configuration for the automation actions (list). [Details](#action-argument-reference)

#### Condition Argument Reference

* `query` - Required - Dynamic Focus Query string that selects a list of entities on which the policy is run

#### Action Argument Reference

* `action` - Required - The full automation action configuration (object). This includes all fields from `instana_automation_action` resource
* `agent_id` - Optional - The identifier for the agent host. Optional if the type configuration is manual. For automatic type configuration, the argument is required

The `action` object should contain the same fields as the `instana_automation_action` resource, including:
- `id` - The action identifier
- `name` - Action name
- `description` - Action description
- `tags` - Optional tags
- `input_parameter` - Optional input parameters
- One of the action type configurations (`script`, `http`, `manual`, `jira`, `github`, `gitlab`, `doc_link`, `ansible`)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the automation policy

## Import

Automation Policies can be imported using the `id`, e.g.:

```bash
terraform import instana_automation_policy.example 60845e4e5e6b9cf8fc2868da
```

## Best Practices

### Policy Types

Choose the appropriate policy type:

- **Manual**: Requires human approval before execution
  - Use for critical operations
  - Use when human judgment is needed
  - Use for operations with potential side effects

- **Automatic**: Executes automatically when triggered
  - Use for well-tested remediation
  - Use for non-critical operations
  - Use with appropriate conditions to limit scope

### Conditions

Use specific conditions to limit automatic policy execution:

```hcl
condition = {
  query = join(" AND ", [
    "entity.zone:production",
    "entity.type:service",
    "entity.agent.capability:action-script",
    "NOT entity.tag:no-auto-remediation"
  ])
}
```

### Multiple Actions

Order actions logically:

1. Diagnostic/data collection
2. Remediation
3. Notification

### Triggers

- Use appropriate trigger types for your use case
- For scheduled policies, use clear recurrence rules
- Test triggers before enabling automatic policies

### Tags

Use tags for organization:

```hcl
tags = ["environment:production", "team:platform", "criticality:high"]
```

## Notes

- The resource ID is auto-generated by Instana upon creation
- Manual policies require human approval before execution
- Automatic policies execute immediately when conditions are met
- Multiple actions execute in the order defined
- Conditions use Instana's Dynamic Focus query language
- Scheduled triggers use iCalendar RRULE format for recurrence
- Actions in policies must include full action configuration, not just references
- Agent ID is required for automatic policies but optional for manual policies