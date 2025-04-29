# Automation Action Resource

Management of Automation Actions.

API Documentation: <https://instana.github.io/openapi/#tag/Action-Catalog>

## Example Usage


### Create a script action
```hcl
resource "instana_automation_action" "hello_world" {
  name            = "Hello world"
  description     = "Script action for test"
  tags            = ["test", "hello"]
  timeout         = "10"
  type            = "script"

  interpreter     = "bash"
  content         = <<EOF
    echo "Hello world!"
  EOF

  input_parameter {
      name        = "test"
      label       = "test parameter"
      description = "parameter for test"
      type        = "static"
      required    = true
      hidden      = false
      secured     = false
      value       = ""
  }
}
```


### Create a HTTP action
```hcl
resource "instana_automation_action" "http_sample" {
  name            = "Instana health"
  description     = "Instana health status check"
  tags            = ["test"]
  timeout         = "10"

  http { 
    host                = "@@instana_api_host@@/api/instana/health"
    method              = "GET"
    ignoreCertErrors    = true
  }

  input_parameter {
      name        = "test"
      label       = "test parameter"
      description = "parameter for test"
      type        = "static"
      required    = true
      hidden      = false
      secured     = false
      value       = ""
  }
}
```


## Argument Reference

* `label` - Required - The name of the synthetic monitor
* `description` - Optional - The name of the synthetic monitor
* `active` - Optional - Enables/disables the synthetic monitor (defaults to true)
* `application_id` - Optional - Unique identifier of the Application Perspective.
* `custom_properties` - Optional - A map of key/values which are used as tags
* `locations` - Required - A list of strings with location IDs 
* `playback_mode` - Optional - Defines how the Synthetic test should be executed across multiple PoPs (defaults to Simultaneous)
* `test_frequency` - Optional - how often the playback for a synthetic monitor is scheduled (defaults to 15 seconds)

Exactly on of the following configuration blocks must be provided:
* `http_action` - Optional - Http Action Configuration block [Details](#http-action-configuration)
* `http_script` - Optional - HTTP Script Configuration block [Details](#http-script-configuration)

## Import

Automation actions can be imported using the `id`, e.g.:

```
$ terraform import instana_automation_action.script_action cl1g4qrmo26x930s17i2
```