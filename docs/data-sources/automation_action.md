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

* `name` - Required - Name of the automation action.
* `type` - Required - Type of the automation action.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - The automation action identifier.
* `description` - The automation action description.
* `tags` - List of automation actions tags.
