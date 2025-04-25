# Automation Action Data Source

Data source to get the automation action from Instana API. This allows you to retrieve the automation action
by name and type and reference it in other resources such as Automation Policy.

API Documentation: <https://instana.github.io/openapi/#operation/getActions>

## Example Usage

```hcl
data "instana_automation_action" "hello_world" {
  name = "Hello world"
  type = "script"
}
```

## Argument Reference

* `name` - Required - the name of the automation action
* `type` - Required - the type of the automation action
