# Automation Policy Resource

Management of Automation Policies.

API Documentation: <https://instana.github.io/openapi/#tag/Policies>

---
## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.

## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- `trigger` now uses **object syntax** with `= { }`
- `type_configuration` now uses **list syntax** with `= [{ }]`
- `action` within type_configuration uses **list syntax** with `= [{ }]`
- `condition` now uses **object syntax** with `= { }`
- Actions now require full action definition (not just action_id)
- Enhanced validation and better error messages
- Improved state management

#### OLD (v5.x) Syntax:

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

#### NEW (v6.x) Syntax:

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

---


## Example Usage

### Manual Policy with Script Action

```hcl
resource "instana_automation_policy" "automation_policy" {
  description = "Stop the KLO Log File Agent on Linux"
  name        = "Policy_Stop"
  tags        = []
  trigger = {
    description = "Stop ping-springboot-app.sh to trigger event."
    id          = "trigger_id"
    name        = "Legacy Event - Spring Boot App - Requests <= 0"
    type        = "customEvent"
  }
  type_configuration = [
    {
      action = [
        {
          action = {
            description = "Stop the KLO Log File Agent on Linux"
            id          = "action_id"
            input_parameter = [
              {
                description = "Description"
                hidden      = false
                label       = "Instance Name"
                name        = "instance"
                required    = true
                type        = "static"
                value       = ""
              },
            ]
            name   = "Name"
            script = {
              content     = "content"
            }
            tags = ["tag1", "tag2"]
          }
        },
      ]
      name      = "manual"
    },
  ]
}

```

### Automatic Policy with Condition

```hcl
resource "instana_automation_policy" "automation_policy_2" {
  description = "Update ad-service configMap"
  name        = "Update adService config map flagt"
  tags        = []
  trigger = {
    description = "The erroneous call rate is higher or equal to 40%."
    id          = "trigger-id"
    name        = "Erroneous call rate is higher than normal"
    type        = "applicationSmartAlert"
  }
  type_configuration = [
    {
      action = [
        {
          action = {
            ansible = {
              host_id            = "hostid"
              playbook_file_name = "filename"
              playbook_id        = "id"
              url                = "url"
            }
            description = "Update ad-service configMap"
            id          = "action-id"
            input_parameter = [
              {
                description = ""
                hidden      = false
                label       = "label"
                name        = "name"
                required    = false
                type        = "static"
                value       = "value"
              },
              {
                description = ""
                hidden      = false
                label       = "namespace"
                name        = "namespace"
                required    = false
                type        = "static"
                value       = "robot-shop"
              },
              {
                description = ""
                hidden      = true
                label       = "k8s_api_server"
                name        = "k8s_api_server"
                required    = true
                type        = "static"
                value       = "value"
              },
              {
                description = ""
                hidden      = true
                label       = "k8s_user"
                name        = "k8s_user"
                required    = true
                type        = "static"
                value       = "value"
              },
              {
                description = ""
                hidden      = true
                label       = "k8s_password"
                name        = "k8s_password"
                required    = true
                type        = "static"
                value       = "value"
              },
            ]
            name   = "Update ad-service configMap"
          }
          agent_id = "agent-id"
        },
      ]
      condition = {
        query = "entity.aws.elb.type:application"
      }
      name = "manual"
    },
    {
      action = [
        {
          action = {
            ansible = {
              host_id            = "host-id"
              playbook_file_name = "file-name"
              playbook_id        = "id"
              url                = "url"
              workflow_id        = null
            }
            description = "Update ad-service configMap"
            id          = "action-id"
            input_parameter = [
              {
                description = ""
                hidden      = false
                label       = "label"
                name        = "name"
                required    = false
                type        = "static"
                value       = "ad-service"
              },
              {
                description = ""
                hidden      = false
                label       = "namespace"
                name        = "namespace"
                required    = false
                type        = "static"
                value       = "robot-shop"
              },
              {
                description = ""
                hidden      = true
                label       = "k8s_api_server"
                name        = "k8s_api_server"
                required    = true
                type        = "static"
                value       = "value"
              },
              {
                description = ""
                hidden      = true
                label       = "k8s_user"
                name        = "k8s_user"
                required    = true
                type        = "static"
                value       = "value"
              },
              {
                description = ""
                hidden      = true
                label       = "k8s_password"
                name        = "k8s_password"
                required    = true
                type        = "static"
                value       = "vallue"
              },
            ]
            name   = "Update ad-service configMap"
          }
          agent_id = "agent-id"
        },
      ]
      condition = {
        query = "entity.aws.elb.type:application"
      }
      name = "automatic"
    },
  ]
}

```

## Generating Configuration from Existing Resources

If you have already created a automation policy in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your automation policy in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_automation_policy.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual automation policy ID from Instana

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