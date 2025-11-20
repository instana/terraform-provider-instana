# Synthetic Test Resource

Manages synthetic test configurations in Instana. Synthetic tests monitor the availability and performance of your applications and services from various locations around the world.

API Documentation: <https://instana.github.io/openapi/#operation/getSyntheticTests>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

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
}
```

#### NEW (v6.x) Syntax:
```hcl
resource "instana_synthetic_test" "example" {
  label = "HTTP Test"
  
  http_action = {
    url = "https://example.com"
    operation = "GET"
  }
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
    expect_match      = "\"id\":\\s*\"[a-zA-Z0-9]+\""
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
    expect_exists = [
      "$.status",
      "$.data.version",
      "$.data.uptime"
    ]
    expect_not_empty = [
      "$.data.services"
    ]
    expect_json = {
      "status" = "ok"
    }
  }
}
```

### HTTP Script Test (Basic)

Simple HTTP script test:

```hcl
resource "instana_synthetic_test" "http_script_basic" {
  label       = "HTTP Script Test"
  description = "Custom HTTP script"
  active      = true
  locations   = [data.instana_synthetic_location.loc1.id]
  
  http_script = {
    mark_synthetic_call = true
    retries             = 1
    retry_interval      = 2
    timeout             = "1m"
    script_type         = "Basic"
    script              = <<-EOF
      const assert = require('assert');
      
      $http.get('https://example.com', function(err, response, body) {
        if (err) {
          console.error(err);
          return;
        }
        assert.equal(response.statusCode, 200, 'Expected 200 OK');
        assert(body.includes('Example'), 'Expected body to contain Example');
      });
    EOF
  }
}
```

### HTTP Script Test (Jest)

Advanced Jest-based HTTP script:

```hcl
resource "instana_synthetic_test" "http_script_jest" {
  label     = "Jest HTTP Script"
  active    = true
  locations = [data.instana_synthetic_location.loc1.id]
  
  http_script = {
    script_type = "Jest"
    file_name   = "test.js"
    scripts = {
      bundle      = "// Bundle content here"
      script_file = <<-EOF
        const axios = require('axios');
        
        test('API returns valid response', async () => {
          const response = await axios.get('https://api.example.com/health');
          expect(response.status).toBe(200);
          expect(response.data).toHaveProperty('status', 'healthy');
        });
      EOF
    }
  }
}
```

### Browser Script Test

Browser automation test:

```hcl
resource "instana_synthetic_test" "browser_script" {
  label       = "Browser Script Test"
  description = "Test user login flow"
  active      = true
  locations   = [data.instana_synthetic_location.loc1.id]
  
  browser_script = {
    mark_synthetic_call = true
    retries             = 1
    retry_interval      = 2
    timeout             = "2m"
    browser             = "chrome"
    record_video        = true
    script_type         = "Basic"
    script              = <<-EOF
      const { Builder, By, until } = require('selenium-webdriver');
      
      (async function() {
        let driver = await new Builder().forBrowser('chrome').build();
        try {
          await driver.get('https://example.com/login');
          await driver.findElement(By.id('username')).sendKeys('testuser');
          await driver.findElement(By.id('password')).sendKeys('password123');
          await driver.findElement(By.id('submit')).click();
          await driver.wait(until.urlContains('/dashboard'), 5000);
        } finally {
          await driver.quit();
        }
      })();
    EOF
  }
}
```

### DNS Test

DNS resolution test:

```hcl
resource "instana_synthetic_test" "dns_test" {
  label       = "DNS Resolution Test"
  description = "Verify DNS records"
  active      = true
  locations   = [data.instana_synthetic_location.loc1.id]
  
  dns = {
    lookup              = "example.com"
    server              = "8.8.8.8"
    query_type          = "A"
    port                = 53
    transport           = "UDP"
    accept_cname        = true
    recursive_lookups   = true
    server_retries      = 2
    query_time = {
      key      = "queryTime"
      operator = "LESS_THAN"
      value    = 100
    }
    target_values = [
      {
        key      = "A"
        operator = "EQUALS"
        value    = "93.184.216.34"
      }
    ]
  }
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
    hostname                      = "example.com"
    port                          = 443
    days_remaining_check          = 30
    accept_self_signed_certificate = false
    validation_rules = [
      {
        key      = "issuer"
        operator = "CONTAINS"
        value    = "DigiCert"
      },
      {
        key      = "subject"
        operator = "EQUALS"
        value    = "CN=example.com"
      }
    ]
  }
}
```

### Webpage Action Test

Simple webpage load test:

```hcl
resource "instana_synthetic_test" "webpage_action" {
  label       = "Webpage Load Test"
  description = "Monitor page load time"
  active      = true
  locations   = [data.instana_synthetic_location.loc1.id]
  
  webpage_action = {
    url          = "https://example.com"
    browser      = "chrome"
    record_video = false
    timeout      = "30s"
  }
}
```

### Webpage Script Test

Custom webpage interaction script:

```hcl
resource "instana_synthetic_test" "webpage_script" {
  label       = "Webpage Script Test"
  description = "Test complex user interactions"
  active      = true
  locations   = [data.instana_synthetic_location.loc1.id]
  
  webpage_script = {
    browser      = "firefox"
    record_video = true
    timeout      = "1m"
    file_name    = "interaction.js"
    script       = <<-EOF
      const { Builder, By } = require('selenium-webdriver');
      
      (async function() {
        let driver = await new Builder().forBrowser('firefox').build();
        try {
          await driver.get('https://example.com');
          await driver.findElement(By.css('.search-box')).sendKeys('test query');
          await driver.findElement(By.css('.search-button')).click();
        } finally {
          await driver.quit();
        }
      })();
    EOF
  }
}
```

### Test with Multiple Locations

Run test from multiple geographic locations:

```hcl
resource "instana_synthetic_test" "multi_location" {
  label          = "Multi-Location Test"
  description    = "Test from multiple regions"
  active         = true
  locations      = [
    data.instana_synthetic_location.us_east.id,
    data.instana_synthetic_location.eu_west.id,
    data.instana_synthetic_location.ap_south.id
  ]
  test_frequency = 10
  playback_mode  = "Staggered"
  
  http_action = {
    url           = "https://example.com"
    operation     = "GET"
    expect_status = 200
  }
}
```

### Test with RBAC Tags

Add RBAC tags for access control:

```hcl
resource "instana_synthetic_test" "with_rbac" {
  label     = "RBAC Tagged Test"
  active    = true
  locations = [data.instana_synthetic_location.loc1.id]
  
  rbac_tags = [
    {
      name  = "team"
      value = "platform"
    },
    {
      name  = "environment"
      value = "production"
    }
  ]
  
  http_action = {
    url       = "https://example.com"
    operation = "GET"
  }
}
```

### Test with Multiple Applications

Associate test with multiple applications:

```hcl
resource "instana_synthetic_test" "multi_app" {
  label        = "Multi-Application Test"
  active       = true
  locations    = [data.instana_synthetic_location.loc1.id]
  applications = [
    instana_application_config.app1.id,
    instana_application_config.app2.id
  ]
  
  http_action = {
    url       = "https://example.com"
    operation = "GET"
  }
}
```

### Test with Mobile Apps

Associate test with mobile applications:

```hcl
resource "instana_synthetic_test" "mobile_app" {
  label       = "Mobile App Backend Test"
  active      = true
  locations   = [data.instana_synthetic_location.loc1.id]
  mobile_apps = [
    "mobile-app-id-1",
    "mobile-app-id-2"
  ]
  
  http_action = {
    url       = "https://api.example.com/mobile/v1/status"
    operation = "GET"
  }
}
```

### Test with Websites

Associate test with website monitoring:

```hcl
resource "instana_synthetic_test" "website_test" {
  label    = "Website Monitoring Test"
  active   = true
  locations = [data.instana_synthetic_location.loc1.id]
  websites = [
    instana_website_monitoring_config.site1.id
  ]
  
  http_action = {
    url       = "https://example.com"
    operation = "GET"
  }
}
```

### Complex DNS Test

Advanced DNS test with multiple validations:

```hcl
resource "instana_synthetic_test" "dns_advanced" {
  label     = "Advanced DNS Test"
  active    = true
  locations = [data.instana_synthetic_location.loc1.id]
  
  dns = {
    lookup            = "example.com"
    server            = "1.1.1.1"
    query_type        = "ALL"
    port              = 53
    transport         = "TCP"
    accept_cname      = true
    lookup_server_name = true
    recursive_lookups = true
    server_retries    = 3
    query_time = {
      key      = "queryTime"
      operator = "LESS_THAN"
      value    = 50
    }
    target_values = [
      {
        key      = "A"
        operator = "EQUALS"
        value    = "93.184.216.34"
      },
      {
        key      = "AAAA"
        operator = "EQUALS"
        value    = "2606:2800:220:1:248:1893:25c8:1946"
      },
      {
        key      = "NS"
        operator = "CONTAINS"
        value    = "example.com"
      }
    ]
  }
}
```

### SSL Certificate with Multiple Validations

Comprehensive SSL certificate test:

```hcl
resource "instana_synthetic_test" "ssl_comprehensive" {
  label     = "Comprehensive SSL Test"
  active    = true
  locations = [data.instana_synthetic_location.loc1.id]
  
  ssl_certificate = {
    hostname             = "secure.example.com"
    port                 = 443
    days_remaining_check = 60
    validation_rules = [
      {
        key      = "issuer"
        operator = "CONTAINS"
        value    = "Let's Encrypt"
      },
      {
        key      = "subject"
        operator = "EQUALS"
        value    = "CN=secure.example.com"
      },
      {
        key      = "keySize"
        operator = "GREATER_OR_EQUAL_THAN"
        value    = "2048"
      },
      {
        key      = "signatureAlgorithm"
        operator = "EQUALS"
        value    = "SHA256withRSA"
      }
    ]
  }
}
```

### Environment-Specific Tests

Create tests for different environments:

```hcl
locals {
  environments = {
    production = {
      url       = "https://prod.example.com"
      frequency = 5
      locations = [data.instana_synthetic_location.us_east.id]
    }
    staging = {
      url       = "https://staging.example.com"
      frequency = 15
      locations = [data.instana_synthetic_location.us_west.id]
    }
  }
}

resource "instana_synthetic_test" "env_tests" {
  for_each = local.environments

  label          = "${title(each.key)} Environment Test"
  description    = "Monitor ${each.key} environment"
  active         = true
  locations      = each.value.locations
  test_frequency = each.value.frequency
  
  http_action = {
    url           = each.value.url
    operation     = "GET"
    expect_status = 200
  }
  
  custom_properties = {
    "environment" = each.key
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
* Use `playback_mode = "Staggered"` to distribute test execution across time
* RBAC tags control who can view and modify the test
* Custom properties help organize and filter tests
* Video recording is available for browser-based tests
* SSL certificate tests alert before certificates expire
* DNS tests validate DNS resolution and records
* HTTP scripts support both Basic and Jest frameworks
* Browser scripts use Selenium WebDriver
