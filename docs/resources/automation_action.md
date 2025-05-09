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

  script {
    interpreter     = "bash"
    content         = filebase64("test.sh")
    timeout         = "10"
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


### Create a HTTP action
```hcl
resource "instana_automation_action" "http_sample" {
  name            = "Instana health"
  description     = "Instana health status check"
  tags            = ["test"]

  http { 
    host                      = "@@instana_api_host@@/api/instana/health"
    method                    = "POST"
    ignore_certificate_errors = true
    timeout                   = "10"

    headers = {
      "Authentication"  = "Bearer <token>"
      "Accept-Language" = "application/json"
      "Content-Type"    = "application/json"
    } 
    body    = "{}"

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

* `name` - Required - The name of the automation action.
* `description` - Required - The description of the automation action.
* `tags` - Optional - A list of tags for the automation action.
* `input_parameter` - Optional - A list of input parameters [Details](#input-parameter-argument-reference)

Exactly on of the following blocks must be provided:
* `script` - Optional - Http Action Configuration block [Details](#script-argument-reference)
* `http` - Optional - HTTP Script Configuration block [Details](#http-argument-reference)

### Input Parameter Argument Reference

* `name` - Required - The name of the input parameter.
* `label` - Optional - The label of the input parameter.
* `description` - Optional - The description of the input parameter.
* `type` - Required - The type of the input parameter. It can be static or dynamic.
* `required` - Required - Indicates if the input parameter is required.
* `hidden` - Optional - Indicates if the input parameter is hidden. By default it is false.
* `secured` - Optional - Indicates if the input parameter is secured. By default it is false.
* `value` - Required - The value of the input parameter.

### Script Argument Reference

* `content` - Required - Base64 encoded script content.
* `interpreter` - Optional - The interpreter for script execution.
* `timeout` - Optional - The timeout of the automation action.

### Http Argument Reference

* `host` - Required - The host for the http request.
* `method` - Required - The method for http request.
* `ignore_certificate_errors` - Optional - Indicates if the http request ignores the certificate errors.
* `headers` - Optional - The headers of the http request.
* `body` - Optional - The body content for the http request.
* `timeout` - Optional - The timeout of the automation action.