# Synthetic Test Resource

Manages synthetic test configurations in Instana. Synthetic tests monitor the availability and performance of your applications and services from various locations around the world.

API Documentation: <https://instana.github.io/openapi/#operation/getSyntheticTests>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

The main change is that all test configuration blocks now use attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
```hcl
resource "instana_synthetic_test" "example" {
  label = "HTTP Test"
  
  http_action {
    url = "https://example.com"
    operation = "GET"
  }
  # rest of the configuration
}
```

#### NEW (v6.x) Syntax:
```hcl
resource "instana_synthetic_test" "example" {
  active            = true
  http_action = {
    allow_insecure   = true
    expect_status    = 200
    follow_redirect  = true
    headers = {
      test = "test"
    }
    mark_synthetic_call = true
    operation           = "GET"
    retries             = 0
    retry_interval      = 1
    timeout             = "0m"
    url                 = "var.url"
    validation_string   = "test"
  }
  label           = "test label"
  locations       = ["var.location"]
  playback_mode   = "Simultaneous"
  test_frequency  = 15
}
```

### Key Syntax Changes

1. **Test Configuration Blocks**: `http_action { }` → `http_action = { }`
2. **Nested Objects**: `scripts { }` → `scripts = { }`
3. **RBAC Tags**: `rbac_tags { }` (multiple) → `rbac_tags = [{ }]` (set)
4. **DNS Filters**: `target_values { }` (multiple) → `target_values = [{ }]` (set)
5. **SSL Validation**: `validation_rules { }` (multiple) → `validation_rules = [{ }]` (set)

## Example Usage

### HTTP Action Test

Basic HTTP GET request test:

```hcl
resource "instana_synthetic_test" "http_action_basic" {
  label          = "Basic HTTP Test"
  description    = "Monitor website availability"
  active         = true
  application_id = "my-app-id"
  locations      = [data.instana_synthetic_location.loc1.id]
  test_frequency = 15
  playback_mode  = "Simultaneous"
  
  http_action = {
    mark_synthetic_call = true
    retries             = 2
    retry_interval      = 1
    timeout             = "30s"
    url                 = "https://example.com"
    operation           = "GET"
    follow_redirect     = true
    allow_insecure      = false
    expect_status       = 200
  }
  
  custom_properties = {
    "environment" = "production"
    "team"        = "platform"
  }
}
```

### HTTP Action with Headers and Validation

POST request with custom headers and response validation:

```hcl
resource "instana_synthetic_test" "http_action_advanced" {
  label       = "API Endpoint Test"
  description = "Test API with authentication"
  active      = true
  locations   = [data.instana_synthetic_location.us_east.id]
  
  http_action = {
    url       = "https://api.example.com/v1/users"
    operation = "POST"
    headers = {
      "Authorization" = "Bearer token123"
      "Content-Type"  = "application/json"
    }
    body              = jsonencode({
      name  = "test"
      email = "test@example.com"
    })
    expect_status     = 201
    validation_string = "success"
  }
}
```

### HTTP Action with JSON Validation

Validate JSON response structure:

```hcl
resource "instana_synthetic_test" "http_json_validation" {
  label     = "JSON API Test"
  active    = true
  locations = [data.instana_synthetic_location.loc1.id]
  
  http_action = {
    url           = "https://api.example.com/status"
    operation     = "GET"
    expect_status = 200
    expect_json = {
      "status" = "ok"
    }
  }
}
```
### DNS Test

DNS resolution test:

```hcl
resource "instana_synthetic_test" "dns_test" {
  label           = "dns_test"
  active            = true
  dns = {
    accept_cname        = false
    lookup              = "var.url"
    lookup_server_name  = false
    mark_synthetic_call = true
    query_time = {
      key      = "responseTime"
      operator = "LESS_THAN"
      value    = 120
    }
    query_type        = "A"
    recursive_lookups = true
    retries           = 0
    retry_interval    = 1
    server            = "8.8.8.8"
    server_retries    = 1
    timeout           = "0m"
    transport         = "UDP"
  }
  locations       = ["b8dsyQt4fDukWzR9RMXW"]
  playback_mode   = "Simultaneous"
  test_frequency  = 15

}
```

### SSL Certificate Test

SSL certificate validation test:

```hcl
resource "instana_synthetic_test" "ssl_cert_test" {
  label       = "SSL Certificate Test"
  description = "Monitor SSL certificate expiration"
  active      = true
  locations   = [data.instana_synthetic_location.loc1.id]
  ssl_certificate = {
    accept_self_signed_certificate = false
    days_remaining_check           = 11
    hostname                       = "var.url"
    mark_synthetic_call            = true
    port                           = null
    retries                        = 0
    retry_interval                 = 1
    timeout                        = "0m"
  }
}
```
## Argument Reference

* `label` - Required - Name of the synthetic test (max 128 characters)
* `description` - Optional - Description of the synthetic test (max 512 characters)
* `active` - Optional - Boolean flag to enable/disable the test. Default: `true`
* `application_id` - Optional - Unique identifier of the Application Perspective (deprecated, use `applications` instead)
* `applications` - Optional - Set of application IDs to associate with this test
* `mobile_apps` - Optional - Set of mobile app IDs to associate with this test
* `websites` - Optional - Set of website IDs to associate with this test
* `custom_properties` - Optional - Map of key/value pairs used as tags
* `locations` - Required - Set of location IDs where the test should run
* `rbac_tags` - Optional - Set of RBAC tags for access control [Details](#rbac-tags-reference)
* `playback_mode` - Optional - How the test executes across multiple locations. Values: `Simultaneous`, `Staggered`. Default: `Simultaneous`
* `test_frequency` - Optional - How often the test runs in minutes (1-120). Default: `15`

**Exactly one of the following test configuration blocks must be provided:**
* `http_action` - Optional - HTTP action test configuration [Details](#http-action-reference)
* `http_script` - Optional - HTTP script test configuration [Details](#http-script-reference)
* `browser_script` - Optional - Browser script test configuration [Details](#browser-script-reference)
* `dns` - Optional - DNS test configuration [Details](#dns-reference)
* `ssl_certificate` - Optional - SSL certificate test configuration [Details](#ssl-certificate-reference)
* `webpage_action` - Optional - Webpage action test configuration [Details](#webpage-action-reference)
* `webpage_script` - Optional - Webpage script test configuration [Details](#webpage-script-reference)

### RBAC Tags Reference

* `name` - Required - Tag name
* `value` - Required - Tag value

### HTTP Action Reference

* `mark_synthetic_call` - Optional - Mark HTTP calls as synthetic. Default: `false`
* `retries` - Optional - Number of retry attempts (0-2). Default: `0`
* `retry_interval` - Optional - Time between retries in seconds (1-10). Default: `1`
* `timeout` - Optional - Timeout duration (e.g., "30s", "1m")
* `url` - Required - URL to test (must start with http:// or https://)
* `operation` - Optional - HTTP method. Values: `GET`, `HEAD`, `OPTIONS`, `PATCH`, `POST`, `PUT`, `DELETE`
* `headers` - Optional - Map of HTTP headers
* `body` - Optional - Request body content
* `validation_string` - Optional - String that must be present in response
* `follow_redirect` - Optional - Follow HTTP redirects. Default: `false`
* `allow_insecure` - Optional - Allow insecure SSL certificates. Default: `false`
* `expect_status` - Optional - Expected HTTP status code
* `expect_match` - Optional - Regular expression to match in response
* `expect_exists` - Optional - Set of JSON paths that must exist in response
* `expect_not_empty` - Optional - Set of JSON paths that must not be empty in response
* `expect_json` - Optional - Map of expected JSON key/value pairs

### HTTP Script Reference

* `mark_synthetic_call` - Optional - Mark HTTP calls as synthetic. Default: `false`
* `retries` - Optional - Number of retry attempts (0-2). Default: `0`
* `retry_interval` - Optional - Time between retries in seconds (1-10). Default: `1`
* `timeout` - Optional - Timeout duration
* `script` - Optional - JavaScript content (for single script)
* `script_type` - Optional - Script type. Values: `Basic`, `Jest`
* `file_name` - Optional - Script file name
* `scripts` - Optional - Multiple scripts configuration for Jest [Details](#scripts-reference)

#### Scripts Reference

* `bundle` - Optional - Bundle content
* `script_file` - Optional - Script file content

### Browser Script Reference

* `mark_synthetic_call` - Optional - Mark HTTP calls as synthetic. Default: `false`
* `retries` - Optional - Number of retry attempts (0-2). Default: `0`
* `retry_interval` - Optional - Time between retries in seconds (1-10). Default: `1`
* `timeout` - Optional - Timeout duration
* `script` - Optional - JavaScript content
* `script_type` - Optional - Script type. Values: `Basic`, `Jest`
* `file_name` - Optional - Script file name
* `scripts` - Optional - Multiple scripts configuration [Details](#scripts-reference)
* `browser` - Optional - Browser type. Values: `chrome`, `firefox`
* `record_video` - Optional - Record video of test execution. Default: `false`

### DNS Reference

* `mark_synthetic_call` - Optional - Mark calls as synthetic. Default: `false`
* `retries` - Optional - Number of retry attempts (0-2). Default: `0`
* `retry_interval` - Optional - Time between retries in seconds (1-10). Default: `1`
* `timeout` - Optional - Timeout duration
* `lookup` - Required - Domain name to lookup
* `server` - Required - DNS server to query
* `query_type` - Optional - DNS query type. Values: `ALL`, `ALL_CONDITIONS`, `ANY`, `A`, `AAAA`, `CNAME`, `NS`
* `port` - Optional - DNS server port
* `transport` - Optional - Transport protocol. Values: `TCP`, `UDP`
* `accept_cname` - Optional - Accept CNAME records
* `lookup_server_name` - Optional - Lookup server name
* `recursive_lookups` - Optional - Enable recursive lookups
* `server_retries` - Optional - Number of server retries
* `query_time` - Optional - Query time filter [Details](#query-time-reference)
* `target_values` - Optional - Set of target value filters [Details](#target-values-reference)

#### Query Time Reference

* `key` - Required - Filter key
* `operator` - Required - Filter operator. Values: `CONTAINS`, `EQUALS`, `GREATER_THAN`, `IS`, `LESS_THAN`, `MATCHES`, `NOT_MATCHES`
* `value` - Required - Filter value (integer)

#### Target Values Reference

* `key` - Required - Filter key. Values: `ALL`, `ALL_CONDITIONS`, `ANY`, `A`, `AAAA`, `CNAME`, `NS`
* `operator` - Required - Filter operator. Values: `CONTAINS`, `EQUALS`, `GREATER_THAN`, `IS`, `LESS_THAN`, `MATCHES`, `NOT_MATCHES`
* `value` - Required - Filter value (string)

### SSL Certificate Reference

* `mark_synthetic_call` - Optional - Mark calls as synthetic. Default: `false`
* `retries` - Optional - Number of retry attempts (0-2). Default: `0`
* `retry_interval` - Optional - Time between retries in seconds (1-10). Default: `1`
* `timeout` - Optional - Timeout duration
* `hostname` - Required - Hostname to check SSL certificate (max 2047 characters)
* `days_remaining_check` - Required - Minimum days remaining before certificate expiration (1-365)
* `accept_self_signed_certificate` - Optional - Accept self-signed certificates
* `port` - Optional - Port number
* `validation_rules` - Optional - Set of SSL certificate validation rules [Details](#validation-rules-reference)

#### Validation Rules Reference

* `key` - Required - Validation key (e.g., "issuer", "subject", "keySize", "signatureAlgorithm")
* `operator` - Required - Validation operator. Values: `CONTAINS`, `EQUALS`, `GREATER_THAN`, `IS`, `LESS_THAN`, `MATCHES`, `NOT_MATCHES`
* `value` - Required - Validation value (string)

### Webpage Action Reference

* `mark_synthetic_call` - Optional - Mark calls as synthetic. Default: `false`
* `retries` - Optional - Number of retry attempts (0-2). Default: `0`
* `retry_interval` - Optional - Time between retries in seconds (1-10). Default: `1`
* `timeout` - Optional - Timeout duration
* `url` - Required - URL to test (must start with http:// or https://)
* `browser` - Optional - Browser type. Values: `chrome`, `firefox`
* `record_video` - Optional - Record video of test execution. Default: `false`

### Webpage Script Reference

* `mark_synthetic_call` - Optional - Mark calls as synthetic. Default: `false`
* `retries` - Optional - Number of retry attempts (0-2). Default: `0`
* `retry_interval` - Optional - Time between retries in seconds (1-10). Default: `1`
* `timeout` - Optional - Timeout duration
* `script` - Required - JavaScript content
* `file_name` - Optional - Script file name
* `browser` - Optional - Browser type. Values: `chrome`, `firefox`
* `record_video` - Optional - Record video of test execution. Default: `false`

## Attributes Reference

* `id` - The ID of the synthetic test

## Import

Synthetic tests can be imported using the `id`, e.g.:

```bash
$ terraform import instana_synthetic_test.example cl1g4qrmo26x930s17i2
```

## Notes

* The ID is auto-generated by Instana
* Only one test configuration type can be specified per resource
* Tests run from specified locations at the configured frequency
* RBAC tags control who can view and modify the test
* Custom properties help organize and filter tests
* SSL certificate tests alert before certificates expire
* DNS tests validate DNS resolution and records
* HTTP scripts support both Basic and Jest frameworks
