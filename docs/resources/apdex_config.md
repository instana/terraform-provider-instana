# Apdex Configuration Resource

Apdex (Application Performance Index) is an open standard for measuring application performance. It converts response times into a score between 0 and 1, where 1 represents perfect performance. The Apdex score is calculated based on a threshold value that defines the boundary between satisfactory and tolerable response times.

API Documentation: <https://instana.github.io/openapi/#tag/Apdex-Configuration>

**⚠️ BREAKING CHANGES - Plugin Framework Migration**

**This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**. While the basic structure remains similar, there are important syntax changes for nested configurations.

## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- The `apdex_entity` configuration is now a **single nested attribute** instead of a block
- Use `attribute = { ... }` syntax instead of `attribute { ... }` block syntax
- All nested attributes within entity follow the same pattern
- The `id` attribute is computed and auto-generated
- Tag filter expressions are validated and normalized
- RBAC tags are now structured as list of objects with `display_name` and `id`

#### OLD (v5.x) Syntax:

```hcl
resource "instana_apdex_config" "example" {
  apdex_name = "my-apdex"
  
  apdex_entity {
    application {
      entity_id      = "app-123"
      threshold      = 500
      boundary_scope = "ALL"
    }
  }
}
```

#### NEW (v6.x) Syntax:

```hcl
resource "instana_apdex_config" "example" {
  apdex_name = "my-apdex"
  
  apdex_entity = {
    application = {
      entity_id         = "btg-B701Rx6o9QNXUS4TVw"
      threshold         = 500
      boundary_scope    = "ALL"
      include_internal  = false
      include_synthetic = false
    }
  }
  
  tags = ["production", "critical"]
}
```

Please update your Terraform configurations to use the new attribute-based syntax with equals signs.

## Example Usage

### Application Entity Apdex

#### Basic Application Apdex
```hcl
resource "instana_apdex_config" "app_basic" {
  apdex_name = "app-apdex-basic"
  
  apdex_entity = {
    application = {
      entity_id      = "btg-B701Rx6o9QNXUS4TVw"  # Application ID
      threshold      = 500                        # 500ms threshold
      boundary_scope = "ALL"
    }
  }
  
  tags = ["application", "production"]
}
```

#### Application Apdex with All Options
```hcl
resource "instana_apdex_config" "app_complete" {
  apdex_name = "app-apdex-complete"
  
  apdex_entity = {
    application = {
      entity_id         = "btg-B701Rx6o9QNXUS4TVw"
      threshold         = 1000
      boundary_scope    = "INBOUND"
      include_internal  = true
      include_synthetic = true
      filter_expression = "(call.http.status@na EQUALS 200 AND call.http.method@na EQUALS 'GET')"
    }
  }
  
  tags = ["application", "filtered"]
  
  rbac_tags = [
    {
      display_name = "Team Platform"
      id           = "team-platform-id"
    },
    {
      display_name = "Environment Production"
      id           = "env-production-id"
    }
  ]
}
```

#### Application Apdex with Inbound Scope
```hcl
resource "instana_apdex_config" "app_inbound" {
  apdex_name = "app-apdex-inbound"
  
  apdex_entity = {
    application = {
      entity_id      = "btg-B701Rx6o9QNXUS4TVw"
      threshold      = 750
      boundary_scope = "INBOUND"  # Only measure inbound calls
    }
  }
}
```

### Website Entity Apdex

#### Basic Website Apdex with Page Load
```hcl
resource "instana_apdex_config" "website_pageload" {
  apdex_name = "website-pageload-apdex"
  
  apdex_entity = {
    website = {
      entity_id   = "NSQQfMicRkyl5lDAprqNSA"  # Website ID
      threshold   = 2000                       # 2000ms threshold
      beacon_type = "pageLoad"
    }
  }
  
  tags = ["website", "pageload"]
}
```

#### Website Apdex with HTTP Request
```hcl
resource "instana_apdex_config" "website_http" {
  apdex_name = "website-http-apdex"
  
  apdex_entity = {
    website = {
      entity_id   = "NSQQfMicRkyl5lDAprqNSA"
      threshold   = 500
      beacon_type = "httpRequest"
    }
  }
}
```

#### Website Apdex with Custom Beacon
```hcl
resource "instana_apdex_config" "website_custom" {
  apdex_name = "website-custom-apdex"
  
  apdex_entity = {
    website = {
      entity_id   = "NSQQfMicRkyl5lDAprqNSA"
      threshold   = 1500
      beacon_type = "custom"
    }
  }
  
  tags = ["website", "custom"]
}
```

#### Website Apdex with Filter Expression
```hcl
resource "instana_apdex_config" "website_filtered" {
  apdex_name = "website-filtered-apdex"
  
  apdex_entity = {
    website = {
      entity_id         = "NSQQfMicRkyl5lDAprqNSA"
      threshold         = 1000
      beacon_type       = "httpRequest"
      filter_expression = "(beacon.page.name@na CONTAINS 'checkout' AND beacon.http.status@na LESS_THAN 400)"
    }
  }
  
  tags = ["website", "checkout", "filtered"]
}
```

### Apdex with RBAC Tags

#### Apdex with Access Control
```hcl
resource "instana_apdex_config" "with_rbac" {
  apdex_name = "apdex-with-rbac"
  
  apdex_entity = {
    application = {
      entity_id      = "btg-B701Rx6o9QNXUS4TVw"
      threshold      = 500
      boundary_scope = "ALL"
    }
  }
  
  tags = ["production", "critical"]
  
  rbac_tags = [
    {
      display_name = "Team Platform"
      id           = "vHvRLxteSBm2hluwj5xcaQ"
    },
    {
      display_name = "Environment Production"
      id           = "prod-env-id"
    }
  ]
}
```

## Generating Configuration from Existing Resources

If you have already created an Apdex configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your Apdex configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_apdex_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual Apdex configuration ID from Instana

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

### Top-Level Attributes

* `id` - (Computed) The unique identifier of the Apdex configuration (auto-generated)
* `apdex_name` - (Required) The name of the Apdex configuration
* `tags` - (Optional) A set of tags associated with the Apdex configuration. Must contain at least one tag if specified
* `rbac_tags` - (Optional) A list of RBAC tags for access control. Must contain at least one tag if specified. Each tag contains:
  * `display_name` - (Required) The display name of the RBAC tag
  * `id` - (Required) The ID of the RBAC tag
* `apdex_entity` - (Required) The entity configuration for the Apdex. Must contain exactly one entity type - [Details](#apdex-entity-attribute)

### Apdex Entity Attribute

The `apdex_entity` attribute must contain exactly one of the following:

* `application` - Application-based Apdex entity - [Details](#application-entity-attributes)
* `website` - Website-based Apdex entity - [Details](#website-entity-attributes)

#### Application Entity Attributes

* `entity_id` - (Required) The application ID (boundary ID)
* `threshold` - (Required) The Apdex threshold value in milliseconds. Defines the boundary between satisfactory and tolerable response times. Must be greater than 0
* `boundary_scope` - (Required) The boundary scope. Valid values:
  * `ALL` - Measure all calls within the application
  * `INBOUND` - Measure only inbound calls to the application
* `filter_expression` - (Optional) Tag filter expression to match specific calls - [Details](#tag-filter-expression-syntax)
* `include_internal` - (Optional) Include internal calls in the Apdex calculation. Defaults to `false`
* `include_synthetic` - (Optional) Include synthetic calls in the Apdex calculation. Defaults to `false`

#### Website Entity Attributes

* `entity_id` - (Required) The website monitoring configuration ID
* `threshold` - (Required) The Apdex threshold value in milliseconds. Must be greater than 0
* `beacon_type` - (Required) The beacon type to measure. Valid values:
  * `pageLoad` - Measure page load performance
  * `httpRequest` - Measure HTTP request performance
  * `custom` - Measure custom beacon performance
* `filter_expression` - (Optional) Tag filter expression to match specific beacons - [Details](#tag-filter-expression-syntax)

### Tag Filter Expression Syntax

The `filter_expression` defines which calls/beacons should be included in the Apdex calculation. It supports:

* Logical `AND` and/or logical `OR` conjunctions (AND has higher precedence than OR)
* Comparison operators: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `LESS_THAN`, `GREATER_THAN`
* Unary operators: `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `NOT_BLANK`

**eBNF Grammar:**

```plain
tag_filter                := logical_or
logical_or                := logical_and OR logical_or | logical_and
logical_and               := primary_expression AND logical_and | primary_expression
primary_expression        := comparison | unary_operator_expression
comparison                := identifier comparison_operator value | identifier@entity_origin comparison_operator value | identifier:tag_key comparison_operator value | identifier:tag_key@entity_origin comparison_operator value
comparison_operator       := EQUALS | NOT_EQUAL | CONTAINS | NOT_CONTAIN | STARTS_WITH | ENDS_WITH | NOT_STARTS_WITH | NOT_ENDS_WITH | GREATER_OR_EQUAL_THAN | LESS_OR_EQUAL_THAN | LESS_THAN | GREATER_THAN
unary_operator_expression := identifier unary_operator | identifier@entity_origin unary_operator
unary_operator            := IS_EMPTY | NOT_EMPTY | IS_BLANK | NOT_BLANK
tag_key                   := identifier | string_value
entity_origin             := src | dest | na
value                     := string_value | number_value | boolean_value
string_value              := "'" <string> "'"
number_value              := (+-)?[0-9]+
boolean_value             := TRUE | FALSE
identifier                := [a-zA-Z_][\.a-zA-Z0-9_\-/]*
```

**Examples:**

```hcl
# Simple equality (no parentheses needed for single condition)
filter_expression = "call.http.status@na EQUALS 200"

# Multiple conditions with AND (parentheses required)
filter_expression = "(call.http.status@na EQUALS 200 AND call.http.method@na EQUALS 'GET')"

# Range check for successful responses (parentheses required)
filter_expression = "(call.http.status@na GREATER_OR_EQUAL_THAN 200 AND call.http.status@na LESS_THAN 400)"

# OR condition (parentheses required)
filter_expression = "(call.http.status@na EQUALS 200 OR call.http.status@na EQUALS 201)"

# Complex expression (parentheses required)
filter_expression = "((call.http.status@na EQUALS 200 OR call.http.status@na EQUALS 201) AND call.duration@na LESS_THAN 1000)"

# String operations (parentheses required for multiple conditions)
filter_expression = "(call.http.path@na STARTS_WITH '/api/' AND call.http.path@na NOT_CONTAINS '/internal/')"

# Unary operators (no parentheses needed for single condition)
filter_expression = "call.error@na NOT_EMPTY"

# Website beacon filtering (no parentheses needed for single condition)
filter_expression = "beacon.page.name@na CONTAINS 'checkout'"

# Website status filtering (no parentheses needed for single condition)
filter_expression = "beacon.http.status@na LESS_THAN 400"
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the Apdex configuration

## Import

Apdex Configs can be imported using the `id`, e.g.:

```bash
$ terraform import instana_apdex_config.example apdex-config-id-123
```

## Notes

### Understanding Apdex Scores

The Apdex score is calculated based on three categories of response times:

- **Satisfied**: Response time ≤ threshold (T)
- **Tolerating**: Response time > T and ≤ 4T
- **Frustrated**: Response time > 4T or error occurred

**Formula**: `Apdex = (Satisfied + (Tolerating / 2)) / Total Requests`

The score ranges from 0 to 1:
- **0.94 - 1.00**: Excellent
- **0.85 - 0.93**: Good
- **0.70 - 0.84**: Fair
- **0.50 - 0.69**: Poor
- **0.00 - 0.49**: Unacceptable

### Choosing the Right Threshold

The threshold value is critical for meaningful Apdex scores:

- **Too low**: Most requests will be "Tolerating" or "Frustrated", making the score less useful
- **Too high**: Most requests will be "Satisfied", hiding performance issues

**Best Practices:**
1. Start with your target response time (e.g., 500ms for web applications)
2. Analyze your actual response time distribution
3. Set the threshold at the 50th percentile (median) of your response times
4. Adjust based on user expectations and business requirements

### Boundary Scope Selection

For application entities, choose the appropriate boundary scope:

**ALL:**
- Measures all calls within the application boundary
- Includes both inbound and internal service-to-service calls
- Best for overall application health monitoring

**INBOUND:**
- Measures only calls entering the application from external sources
- Excludes internal service-to-service communication
- Best for user-facing performance monitoring

### Beacon Types for Websites

Different beacon types measure different aspects of website performance:

**pageLoad:**
- Measures complete page load time
- Includes HTML parsing, resource loading, and rendering
- Best for traditional multi-page applications

**httpRequest:**
- Measures individual HTTP request performance
- Best for API calls and AJAX requests
- Useful for single-page applications (SPAs)

**custom:**
- Measures custom-defined performance markers
- Requires custom instrumentation in your application
- Best for specific user interactions or business transactions

### Filter Expressions

Use filter expressions to:
- Focus on specific endpoints or pages
- Exclude health checks or monitoring traffic
- Measure performance for specific user segments
- Track critical business transactions

**Example Use Cases:**
```hcl
# Only measure checkout flow (single condition, no parentheses needed)
filter_expression = "call.http.path@na STARTS_WITH '/checkout/'"

# Exclude health checks (single condition, no parentheses needed)
filter_expression = "call.http.path@na NOT_EQUALS '/health'"

# Only successful requests (multiple conditions, parentheses required)
filter_expression = "(call.http.status@na GREATER_OR_EQUAL_THAN 200 AND call.http.status@na LESS_THAN 400)"

# Specific HTTP methods (multiple conditions with OR, parentheses required)
filter_expression = "(call.http.method@na EQUALS 'POST' OR call.http.method@na EQUALS 'PUT')"
```

### RBAC Tags

Use RBAC tags to:
- Control access to Apdex configurations
- Organize configurations by team or environment
- Implement multi-tenancy in shared Instana instances
- Enforce security policies

Each RBAC tag requires both a `display_name` (human-readable) and an `id` (system identifier).