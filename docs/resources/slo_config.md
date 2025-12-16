# SLO Configuration Resource

An SLO (Service Level Objective) is a specific, measurable target that defines the expected level of performance, reliability, or availability of a service, agreed upon between a service provider and its users or customers. For instance, an SLO could state that a specific SLI (Service Level Indicator), such as availability, must reach 99.9% over a set period.

API Documentation: <https://instana.github.io/openapi/#operation/createSloConfig

 **⚠️ BREAKING CHANGES - Plugin Framework Migration**

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- All nested configurations (`entity`, `indicator`, `time_window`) are now **single nested attributes** instead of blocks
- Use `attribute = { ... }` syntax instead of `attribute { ... }` block syntax
- All nested attributes within entity, indicator, and time_window follow the same pattern
- The `id` attribute is now computed with a special prefix (`terraform-slo-config-`)
- Tag filter expressions are validated and normalized
- RBAC tags are now structured as list of objects with `display_name` and `id`

#### OLD (v5.x) Syntax:

```hcl
resource "instana_slo_config" "example" {
  name   = "my-slo"
  target = 0.99
  
  entity {
    application {
      application_id = "app-123"
      boundary_scope = "ALL"
    }
  }
  
  indicator {
    time_based_latency {
      threshold   = 500
      aggregation = "P95"
    }
  }
  
  time_window {
    rolling {
      duration      = 7
      duration_unit = "day"
    }
  }
}
```

#### NEW (v6.x) Syntax:

```hcl
resource "instana_slo_config" "example" {
  name   = "my-slo"
  entity = {
    application = {
      application_id    = "Fwx93plkTtKv09MHqT5d6Q" # replace with actual application Id
      boundary_scope    = "ALL"
      include_internal  = false
      include_synthetic = false
    }
  }
  indicator = {
    time_based_availability = {
      aggregation = "MEAN"
      threshold   = 1
    }
  }
  tags      = ["test", "timeBased", "availability", "rolling", "andre"]
  target    = 0.99
  time_window = {
    rolling = {
      duration      = 1
      duration_unit = "week"
      timezone      = "UTC"
    }
  }
}
```

Please update your Terraform configurations to use the new attribute-based syntax with equals signs.

## Example Usage

### Application Entity SLOs

#### Basic Application SLO with Time-Based Latency
```hcl
resource "instana_slo_config" "app_latency_basic" {
  name   = "app-latency-basic"
  entity = {
    application = {
      application_id    = "Fwx93plkTtKv09MHqT5d6Q" # replace with actual application Id
      boundary_scope    = "ALL"
      endpoint_id       = "XKvV8V1KuxP2F1YDxYs_eVmhz4g" # replace with actual end point Id
      include_internal  = false
      include_synthetic = false
    }
  }
  indicator = {
    time_based_latency = {
      aggregation = "MEAN"
      threshold   = 100
    }
  }
  target    = 0.99
  time_window = {
    fixed = {
      duration        = 1
      duration_unit   = "week"
      start_timestamp = 1757961000000
      timezone        = "UTC"
    }
  }
}
```
#### Application SLO with Tag Filter Expression
```hcl
resource "instana_slo_config" "app_with_filter" {
  name   = "app-filtered-slo"
  target = 0.98
  
  entity = {
    application = {
      application_id     = "Fwx93plkTtKv09MHqT5d6Q" # replace with actual application Id
      boundary_scope     = "ALL"
      filter_expression  = "call.http.status@na EQUALS 200 AND call.http.method@na EQUALS 'GET'"
      include_internal   = true
      include_synthetic  = true
    }
  }
  
  indicator = {
    time_based_latency = {
      aggregation = "MEAN"
      threshold   = 100
    }
  }
  
  time_window = {
    rolling = {
      duration      = 7
      duration_unit = "day"
      timezone      = "UTC"

    }
  }
}
```

#### Application SLO with Event-Based Latency
```hcl
resource "instana_slo_config" "app_event_latency" {
  name   = "app-event-based-latency"
  target = 0.95
  
  entity = {
    application = {
      application_id = "Fwx93plkTtKv09MHqT5d6Q" # replace with actual application Id
      boundary_scope = "ALL"
    }
  }
  
  indicator = {
    event_based_latency = {
      threshold = 500
    }
  }
  
  time_window = {
    rolling = {
      duration      = 7
      duration_unit = "day"
    }
  }
}
```

#### Application SLO with Time-Based Availability
```hcl
resource "instana_slo_config" "app_availability" {
  name   = "app-availability-slo"
  target = 0.999
  
  entity = {
    application = {
      application_id = "Fwx93plkTtKv09MHqT5d6Q" # replace with actual application Id
      boundary_scope = "ALL"
    }
  }
  
  indicator = {
    time_based_availability = {
      threshold   = 1
      aggregation = "MEAN"
    }
  }
  
  time_window = {
    rolling = {
      duration      = 30
      duration_unit = "day"
    }
  }
}
```

### Website Entity SLOs

#### Basic Website SLO with Time-Based Latency
```hcl
resource "instana_slo_config" "website_latency" {
  name   = "website-latency-slo"
  target = 0.95
  
  entity = {
    website = {
      website_id  = "NSQQfMicRkyl5lDAprqNSA" # replace with actual website Id
      beacon_type = "httpRequest"
    }
  }
  
  indicator = {
    time_based_latency = {
      threshold   = 2000.0
      aggregation = "P95"
    }
  }
  
  time_window = {
    rolling = {
      duration      = 7
      duration_unit = "day"
    }
  }
}
```

#### Website SLO with Filter Expression
```hcl
resource "instana_slo_config" "website_filtered" {
  name   = "website-filtered-slo"
  target = 0.98
  tags   = ["website", "production"]
  
  entity = {
    website = {
      website_id        = "NSQQfMicRkyl5lDAprqNSA" # replace with actual website Id
      beacon_type       = "httpRequest"
      filter_expression = "beacon.page.name@na CONTAINS 'checkout'"
    }
  }
  
  indicator = {
    time_based_availability = {
      threshold   = 1
      aggregation = "MEAN"
    }
  }
  
  time_window = {
    fixed = {
      duration        = 30
      duration_unit   = "day"
      timezone        = "Europe/London"
      start_timestamp = 1704067200.0
    }
  }
}
```

#### Website SLO with Traffic Indicator
```hcl
resource "instana_slo_config" "website_traffic" {
  name   = "website-traffic-slo"
  target = 0.90
  
  entity = {
    website = {
      website_id  = "NSQQfMicRkyl5lDAprqNSA" # replace with actual website Id
      beacon_type = "httpRequest"
    }
  }
  
  indicator = {
    traffic = {
      threshold    = 185
      traffic_type = "erroneous"
    }
  }
  
  time_window = {
    rolling = {
      duration      = 1
      duration_unit = "day"
    }
  }
}
```

### Synthetic Entity SLOs

#### Basic Synthetic SLO
```hcl
resource "instana_slo_config" "synthetic_basic" {
  name   = "synthetic-basic-slo"
  target = 0.99
  
  entity = {
    synthetic = {
      synthetic_test_ids = [
        "Ob5Aa1F9DagsiNzBaCXX",  # replace with actual synthetic test Ids
        "rXFgndZUj5kDJ2cFbhu1"
      ]
    }
  }
  
  indicator = {
    time_based_availability = {
      threshold   = 0.0
      aggregation = "MEAN"
    }
  }
  
  time_window = {
    rolling = {
      duration      = 7
      duration_unit = "day"
    }
  }
}
```

### Custom Indicator SLOs

#### Custom SLO with Good Events Only
```hcl
resource "instana_slo_config" "custom_good_events" {
  name   = "custom-good-events-slo"
  target = 0.99
  
  entity = {
    application = {
      application_id = "Fwx93plkTtKv09MHqT5d6Q" # replace with actual application Id
      boundary_scope = "ALL"
    }
  }
  
  indicator = {
    custom = {
      good_event_filter_expression = "call.http.status@na EQUALS 200"
    }
  }
  
  time_window = {
    rolling = {
      duration      = 7
      duration_unit = "day"
    }
  }
}
```

#### Custom SLO with Good and Bad Events
```hcl
resource "instana_slo_config" "custom_good_bad_events" {
  name   = "custom-good-bad-events-slo"
  target = 0.95
  tags   = ["custom", "advanced"]
  
  entity = {
    application = {
      application_id = "Fwx93plkTtKv09MHqT5d6Q" # replace with actual application Id
      boundary_scope = "ALL"
      filter_expression = "call.http.status@na GREATER_OR_EQUAL_THAN 200 AND call.http.status@na LESS_THAN 400"
    }
  }
  
  indicator = {
    custom = {
      good_event_filter_expression = "call.http.status@na GREATER_OR_EQUAL_THAN 200 AND call.http.status@na LESS_THAN 400"
      bad_event_filter_expression  = "call.http.status@na GREATER_OR_EQUAL_THAN 500"
    }
  }
  
  time_window = {
    rolling = {
      duration      = 30
      duration_unit = "day"
    }
  }
}
```

### SLOs with RBAC Tags

#### SLO with RBAC Tags
```hcl
resource "instana_slo_config" "with_rbac_tags" {
  name   = "slo-with-rbac"
  target = 0.99
  tags   = ["production", "critical"]
  
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
  
  entity = {
    application = {
      application_id = "Fwx93plkTtKv09MHqT5d6Q" # replace with actual application Id
      boundary_scope = "ALL"
    }
  }
  
  indicator = {
    time_based_latency = {
      threshold   = 500.0
      aggregation = "P95"
    }
  }
  
  time_window = {
    rolling = {
      duration      = 7
      duration_unit = "day"
    }
  }
}
```

## Generating Configuration from Existing Resources

If you have already created a SLO configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your SLO configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_slo_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual SLO configuration ID from Instana

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

* `id` - (Computed) The unique identifier of the SLO configuration (auto-generated with prefix `terraform-slo-config-`)
* `name` - (Required) The name of the SLO configuration
* `target` - (Required) The target SLO value (e.g., 0.99 for 99% availability). Must be a float between 0 and 1
* `tags` - (Optional) A list of tags associated with the SLO configuration
* `rbac_tags` - (Optional) A list of RBAC tags for access control. Each tag contains:
  * `display_name` - (Required) The display name of the RBAC tag
  * `id` - (Required) The ID of the RBAC tag
* `entity` - (Required) The entity configuration for the SLO. Must contain exactly one entity type - [Details](#entity-attribute)
* `indicator` - (Required) The indicator (metric) configuration for the SLO. Must contain exactly one indicator type - [Details](#indicator-attribute)
* `time_window` - (Required) The time window configuration for the SLO. Must contain exactly one time window type - [Details](#time-window-attribute)

### Entity Attribute

The `entity` attribute must contain exactly one of the following:

* `application` - Application-based SLO entity - [Details](#application-entity-attributes)
* `website` - Website-based SLO entity - [Details](#website-entity-attributes)
* `synthetic` - Synthetic-based SLO entity - [Details](#synthetic-entity-attributes)

#### Application Entity Attributes

* `application_id` - (Required) The application ID
* `boundary_scope` - (Required) The boundary scope. Valid values: `ALL`, `INBOUND`
* `filter_expression` - (Optional) Tag filter expression to match events/calls - [Details](#tag-filter-expression-syntax)
* `include_internal` - (Optional) Include internal calls in the scope. Defaults to `false`
* `include_synthetic` - (Optional) Include synthetic calls in the scope. Defaults to `false`
* `service_id` - (Optional) The service ID to filter by
* `endpoint_id` - (Optional) The endpoint ID to filter by

#### Website Entity Attributes

* `website_id` - (Required) The website monitoring configuration ID
* `beacon_type` - (Required) The beacon type. Valid value: `httpRequest`
* `filter_expression` - (Optional) Tag filter expression to match events/calls - [Details](#tag-filter-expression-syntax)

#### Synthetic Entity Attributes

* `synthetic_test_ids` - (Required) A set of synthetic test IDs. Must contain at least one ID.
* `filter_expression` - (Optional) Tag filter expression to match events/calls - [Details](#tag-filter-expression-syntax)

### Indicator Attribute

The `indicator` attribute must contain exactly one of the following:

* `time_based_latency` - Time-based latency indicator - [Details](#time-based-latency-indicator-attributes)
* `event_based_latency` - Event-based latency indicator - [Details](#event-based-latency-indicator-attributes)
* `time_based_availability` - Time-based availability indicator - [Details](#time-based-availability-indicator-attributes)
* `event_based_availability` - Event-based availability indicator - [Details](#event-based-availability-indicator-attributes)
* `traffic` - Traffic indicator - [Details](#traffic-indicator-attributes)
* `custom` - Custom indicator - [Details](#custom-indicator-attributes)

#### Time-Based Latency Indicator Attributes

* `threshold` - (Required) The latency threshold in milliseconds. Must be greater than 0
* `aggregation` - (Required) The aggregation type. Valid values: `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`

#### Event-Based Latency Indicator Attributes

* `threshold` - (Required) The latency threshold in milliseconds. Must be greater than 0

#### Time-Based Availability Indicator Attributes

* `threshold` - (Required) The error rate threshold. Typically `0.0` for availability
* `aggregation` - (Required) The aggregation type. Valid value: `MEAN`

#### Event-Based Availability Indicator Attributes

No additional attributes required. This is an empty object: `event_based_availability = {}`

#### Traffic Indicator Attributes

* `traffic_type` - (Required) The type of traffic to measure. Valid values: `all`, `erroneous`
* `threshold` - (Required) The traffic threshold. Must be greater than 0
* `operator` - (Required) The comparison operator. Valid values: ``, `=`, `<`, `<=`

#### Custom Indicator Attributes

* `good_event_filter_expression` - (Required) Tag filter expression to match good events/calls - [Details](#tag-filter-expression-syntax)
* `bad_event_filter_expression` - (Optional) Tag filter expression to match bad events/calls. If omitted, the opposite of good events is used - [Details](#tag-filter-expression-syntax)

### Time Window Attribute

The `time_window` attribute must contain exactly one of the following:

* `rolling` - Rolling time window - [Details](#rolling-time-window-attributes)
* `fixed` - Fixed time window - [Details](#fixed-time-window-attributes)

#### Rolling Time Window Attributes

* `duration` - (Required) The duration of the time window. Must be an integer
* `duration_unit` - (Required) The unit of the duration. Valid values: `day`, `week`
* `timezone` - (Optional) The timezone for the time window. Defaults to `UTC`

#### Fixed Time Window Attributes

* `duration` - (Required) The duration of the time window. Must be an integer
* `duration_unit` - (Required) The unit of the duration. Valid values: `day`, `week`
* `timezone` - (Optional) The timezone for the time window. Defaults to `UTC`
* `start_timestamp` - (Required) The starting timestamp for the fixed time window (Unix timestamp as float)

### Tag Filter Expression Syntax

The `filter_expression` defines which calls/events should be included. It supports:

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
string_value              := "'" <string "'"
number_value              := (+-)?[0-9]+
boolean_value             := TRUE | FALSE
identifier                := [a-zA-Z_][\.a-zA-Z0-9_\-/]*
```

**Examples:**

```hcl
# Simple equality
filter_expression = "call.http.status EQUALS 200"

# Multiple conditions with AND
filter_expression = "call.http.status EQUALS 200 AND call.http.method EQUALS 'GET'"

# Range check
filter_expression = "call.http.status GREATER_OR_EQUAL_THAN 200 AND call.http.status LESS_THAN 400"

# OR condition
filter_expression = "call.http.status EQUALS 200 OR call.http.status EQUALS 201"

# Complex expression
filter_expression = "(call.http.status EQUALS 200 OR call.http.status EQUALS 201) AND call.duration LESS_THAN 1000"

# String operations
filter_expression = "call.http.path STARTS_WITH '/api/' AND call.http.path NOT_CONTAINS '/internal/'"

# Unary operators
filter_expression = "call.error NOT_EMPTY"
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the SLO configuration (prefixed with `terraform-slo-config-`)

## Import

SLO Configs can be imported using the `id`, e.g.:

```bash
$ terraform import instana_slo_config.example terraform-slo-config-60845e4e5e6fbf76pf8fc2868da
```

## Notes

### Time Window Selection

Choose between rolling and fixed time windows based on your needs:

**Rolling Time Windows:**
- Continuously evaluate SLO over the last N days/weeks
- Best for ongoing monitoring
- Automatically adjusts as time progresses

**Fixed Time Windows:**
- Evaluate SLO for a specific time period
- Best for quarterly/monthly reviews
- Requires explicit start timestamp

### Aggregation Types

Different aggregation types serve different purposes:

- **MEAN**: Average latency across all requests
- **P50** (Median): 50% of requests are faster than this
- **P95**: 95% of requests are faster than this (common choice)
- **P99**: 99% of requests are faster than this (strict SLO)
- **MAX**: Maximum latency observed
- **MIN**: Minimum latency observed