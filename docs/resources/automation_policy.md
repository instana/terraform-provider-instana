# Automation Policy Resource

Management of Automation Policies.

API Documentation: <https://instana.github.io/openapi/#tag/Policies>

## Example Usage

### Create a policy

```hcl
resource "instana_automation_policy" "hello_world_policy" {
  name            = "Hello world"
  description     = "Sample policy for test"
  tags            = ["test", "hello"]

  trigger {
    id          = "r6mpl4BPRhm_77Vn"
    type        = "customEvent"
  }

  type_configuration {
    name = "manual"

    action {
      action_id = "daa60cb8-41fc-3d4c-868f-3c2aca5831fa"
      agent_id = "00:00:0a:ff:fe:0b:21:cc"

      input_parameters = {
        test1 = "value1"
        test2 = "value2"
      }
    }

    condition {
      query = "entity.agent.capability:action-script"
    }
  }
}
```

## Argument Reference

* `name` - Required - The name of the automation policy.
* `description` - Required - The description of the automation policy.
* `tags` - Optional - A list of tags for the automation policy.
* `trigger` - Required - The trigger for the automation policy [Details](#trigger-argument-reference)
* `type_configuration` - Required - A list of configurations with the list of actions to run and the mode (automatic or manual) in which the policy is run. [Details](#type-configuration-argument-reference)

### Trigger Argument Reference

* `id` - Required - Trigger (Instana event or Smart Alert) identifier.
* `type` - Required - Instana event or Smart Alert type.

### Type Configuration Argument Reference

* `name` - Required - The policy type (manual or automatic).
* `condition` - Optional - The condition that selects the list of entities on which the policy is run. Only for automatic policy type. [Details](#condition-argument-reference)
* `action` - Required - The configuration for the automation action.

### Condition Argument Reference

* `query` - Required - Dynamic Focus Query string that selects a list of entities on which the policy is run.

### Action Argument Reference

* `action_id` - Required - The identifier for the automation action.
* `agent_id` - The identifier for the agent host. Optional if the type configuration is manual. For automatic type configuration, the argument is required.
* `input_parameters` - Optional - Map with input parameters name and value.