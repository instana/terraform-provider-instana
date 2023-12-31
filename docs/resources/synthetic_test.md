# Synthetic Test Resource

Synthetic test configuration used to manage synthetic tests in Instana API. Right now, only `HTTPActionConfiguration` 
and `HTTPScriptConfiguration` are supported.

API Documentation: <https://instana.github.io/openapi/#operation/getSyntheticTests>

## Example Usage


### Create a HTTPAction test
```hcl
resource "instana_synthetic_test" "http_action" {
  label          = "test"
  active         = true
  application_id = "my-app-id"
  locations      = [data.instana_synthetic_location.loc1.id]
  test_frequency = 10
  playback_mode  = "Staggered"
  
  http_action { 
    mark_synthetic_call = true
    retries             = 0
    retry_interval      = 1
    timeout             = "3m"
    url                 = "https://example.com"
    operation           = "GET"
  }
  
  custom_properties = {
    "foo" = "bar"
  }
}
```

### Create a HTTPScript test
```hcl
resource "instana_synthetic_test" "http_script" {
  label          = "test"
  active         = true
  application_id = "my-app-id"
  locations      = [data.instana_synthetic_location.loc1.id]
  
  http_script {
    synthetic_type = "HTTPScript"
    script         = <<EOF
      const assert = require('assert');

      $http.get('https://terraform.io',
      function(err, response, body) {
          if (err) {
              console.error(err);
              return;
          }
          assert.equal(response.statusCode, 200, 'Expected a 200 OK response');
      });
    EOF
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

### HTTP Action configuration

* `mark_synthetic_call` - Optional - flag used to control if HTTP calls will be marked as synthetic calls
* `retries` - Optional - Indicates how many attempts will be allowed to get a successful connection (defaults to 0)
* `retry_interval` - Optional - The time interval between retries in seconds (defaults to 1)
* `synthetic_type` - Required - The type of the Synthetic test (currently supports HTTPAction or HTTPScript)
* `timeout` - Optional - The timeout to be used by the PoP playback engines running the test
* `url` - Required when synthetic_type is set to HTTPAction - The URL which is being tested
* `operation` - Optional - The HTTP operation
* `headers` - Optional - An object with header/value pairs
* `body` - Optional - The body content to send with the operation
* `validation_string` - Optional - An expression to be evaluated
* `follow_redirect` - Optional - A boolean type, true by default; to allow redirect
* `allow_insecure` - Optional - A boolean type, if set to true then allow insecure certificates
* `expect_status` - Optional - An integer type, by default, the Synthetic passes for any 2XX status code
* `expect_match` - Optional - An optional regular expression string to be used to check the test response

### HTTP Script configuration

* `mark_synthetic_call` - Optional - flag used to control if HTTP calls will be marked as synthetic calls
* `retries` - Optional - Indicates how many attempts will be allowed to get a successful connection (defaults to 0)
* `retry_interval` - Optional - The time interval between retries in seconds (defaults to 1)
* `synthetic_type` - Required - The type of the Synthetic test (currently supports HTTPAction or HTTPScript)
* `timeout` - Optional - The timeout to be used by the PoP playback engines running the test
* `script` - Required  when synthetic_type is set to HTTPScript - The Javascript content in plain text

## Import

Synthetic monitors can be imported using the `id`, e.g.:

```
$ terraform import instana_synthetic_test.http_action cl1g4qrmo26x930s17i2
```