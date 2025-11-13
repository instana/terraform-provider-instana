# SLI Configuration Resource

> **⚠️ BREAKING CHANGES in v3.0.0**
> 
> This resource has been migrated from Terraform SDK v2 to the Plugin Framework. While most configurations remain compatible, there are important syntax changes you need to be aware of.
>
> **Key Changes:**
> - `metric_configuration` block now uses attribute syntax `= { }` instead of block syntax
> - `sli_entity` block now uses attribute syntax `= { }` instead of block syntax
> - All nested entity type blocks use attribute syntax
> - **IMPORTANT**: SLI Configurations are immutable - updates will result in errors
> - See [Migration Guide](#migration-guide-v2-to-v3) below for detailed examples

Management of SLI configurations. A service level indicator (SLI) is the defined quantitative measure of one characteristic of the level of service that is provided to a customer. Common examples of such indicators are error rate or response latency of a service.

API Documentation: <https://instana.github.io/openapi/#operation/createSli>
**Note:** SLI Configurations cannot be changed. An update of the resource will result in an error. To update an SLI you need to create a new SLI and delete the old one.

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

**Note:** SLI Configurations cannot be changed. An update of the resource will result in an error. To update an SLI you need to create a new SLI and delete the old one.

## Migration Guide (v2 to v3)

### Syntax Changes Overview

The main changes are in how `metric_configuration` and `sli_entity` blocks are defined. In v3, all nested configurations use attribute syntax instead of block syntax.

#### OLD (v2.x) Syntax:
```hcl
resource "instana_sli_config" "example" {
  name = "API Latency SLI"
  
  metric_configuration {
    metric_name = "latency"
    aggregation = "P95"
    threshold = 500
  }
  
  sli_entity {
    application_time_based {
      application_id = "app-123"
      boundary_scope = "INBOUND"
    }
  }
}
```

#### NEW (v3.x) Syntax:
```hcl
resource "instana_sli_config" "example" {
  name = "API Latency SLI"
  
  metric_configuration = {
    metric_name = "latency"
    aggregation = "P95"
    threshold = 500
  }
  
  sli_entity = {
    application_time_based = {
      application_id = "app-123"
      boundary_scope = "INBOUND"
    }
  }
}
```

### Key Syntax Changes

1. **Metric Configuration**: `metric_configuration { }` → `metric_configuration = { }`
2. **SLI Entity**: `sli_entity { }` → `sli_entity = { }`
3. **Entity Types**: `application_time_based { }` → `application_time_based = { }`
4. **All nested objects**: Use `= { }` syntax

## Example Usage

### Application Time-Based SLI (Latency)

```hcl
resource "instana_sli_config" "app_latency" {
  name = "API Latency P95"
  initial_evaluation_timestamp = 0
  
  metric_configuration = {
    metric_name = "latency"
    aggregation = "P95"
    threshold = 500
  }
  
  sli_entity = {
    application_time_based = {
      application_id = "app-prod-api"
      boundary_scope = "INBOUND"
    }
  }
}
```

### Application Time-Based SLI with Service Scope

```hcl
resource "instana_sli_config" "service_latency" {
  name = "Payment Service Latency"
  
  metric_configuration = {
    metric_name = "latency"
    aggregation = "P99"
    threshold = 1000
  }
  
  sli_entity = {
    application_time_based = {
      application_id = "payment-app"
      service_id = "payment-service-id"
      boundary_scope = "ALL"
    }
  }
}
```

### Application Time-Based SLI with Endpoint Scope

```hcl
resource "instana_sli_config" "endpoint_latency" {
  name = "Checkout Endpoint Latency"
  
  metric_configuration = {
    metric_name = "latency"
    aggregation = "P95"
    threshold = 300
  }
  
  sli_entity = {
    application_time_based = {
      application_id = "ecommerce-app"
      service_id = "checkout-service"
      endpoint_id = "POST-/api/checkout"
      boundary_scope = "INBOUND"
    }
  }
}
```

### Application Event-Based SLI (Availability)

```hcl
resource "instana_sli_config" "app_availability" {
  name = "API Availability"
  
  sli_entity = {
    application_event_based = {
      application_id = "api-app"
      boundary_scope = "INBOUND"
      good_event_filter_expression = "call.http.status@na LESS_THAN 500"
      bad_event_filter_expression = "call.http.status@na GREATER_OR_EQUAL_THAN 500"
    }
  }
}
```

### Application Event-Based with Internal and Synthetic Calls

```hcl
resource "instana_sli_config" "comprehensive_availability" {
  name = "Comprehensive Availability SLI"
  
  sli_entity = {
    application_event_based = {
      application_id = "main-app"
      boundary_scope = "ALL"
      include_internal = true
      include_synthetic = true
      good_event_filter_expression = "call.http.status@na EQUALS 200"
      bad_event_filter_expression = "call.http.status@na NOT_EQUAL 200"
    }
  }
}
```

### Application Event-Based with Service Scope

```hcl
resource "instana_sli_config" "service_availability" {
  name = "User Service Availability"
  
  sli_entity = {
    application_event_based = {
      application_id = "user-management"
      service_id = "user-service-id"
      boundary_scope = "INBOUND"
      good_event_filter_expression = "call.http.status@na LESS_THAN 400"
      bad_event_filter_expression = "call.http.status@na GREATER_OR_EQUAL_THAN 400"
    }
  }
}
```

### Application Event-Based with Endpoint Scope

```hcl
resource "instana_sli_config" "endpoint_availability" {
  name = "Login Endpoint Availability"
  
  sli_entity = {
    application_event_based = {
      application_id = "auth-app"
      service_id = "auth-service"
      endpoint_id = "POST-/api/login"
      boundary_scope = "INBOUND"
      good_event_filter_expression = "call.http.status@na EQUALS 200"
      bad_event_filter_expression = "call.http.status@na NOT_EQUAL 200"
    }
  }
}
```

### Website Time-Based SLI (Page Load Time)

```hcl
resource "instana_sli_config" "website_page_load" {
  name = "Homepage Load Time"
  
  metric_configuration = {
    metric_name = "page.load.time"
    aggregation = "P75"
    threshold = 2000
  }
  
  sli_entity = {
    website_time_based = {
      website_id = "website-prod"
      beacon_type = "pageLoad"
    }
  }
}
```

### Website Time-Based with Filter

```hcl
resource "instana_sli_config" "filtered_page_load" {
  name = "Product Pages Load Time"
  
  metric_configuration = {
    metric_name = "page.load.time"
    aggregation = "P90"
    threshold = 3000
  }
  
  sli_entity = {
    website_time_based = {
      website_id = "ecommerce-website"
      beacon_type = "pageLoad"
      filter_expression = "page.url@na CONTAINS '/products/'"
    }
  }
}
```

### Website Event-Based SLI (Error Rate)

```hcl
resource "instana_sli_config" "website_errors" {
  name = "Website Error Rate"
  
  sli_entity = {
    website_event_based = {
      website_id = "website-prod"
      beacon_type = "error"
      good_event_filter_expression = "beacon.error.type@na IS_EMPTY"
      bad_event_filter_expression = "beacon.error.type@na NOT_EMPTY"
    }
  }
}
```

### Website HTTP Request SLI

```hcl
resource "instana_sli_config" "website_api_calls" {
  name = "Website API Call Success Rate"
  
  sli_entity = {
    website_event_based = {
      website_id = "spa-website"
      beacon_type = "httpRequest"
      good_event_filter_expression = "http.status@na LESS_THAN 400"
      bad_event_filter_expression = "http.status@na GREATER_OR_EQUAL_THAN 400"
    }
  }
}
```

### Website Resource Load SLI

```hcl
resource "instana_sli_config" "resource_load_time" {
  name = "Resource Load Performance"
  
  metric_configuration = {
    metric_name = "resource.load.time"
    aggregation = "P95"
    threshold = 1000
  }
  
  sli_entity = {
    website_time_based = {
      website_id = "website-prod"
      beacon_type = "resourceLoad"
      filter_expression = "resource.type@na EQUALS 'script'"
    }
  }
}
```

### Custom Beacon SLI

```hcl
resource "instana_sli_config" "custom_metric" {
  name = "Custom Business Metric"
  
  metric_configuration = {
    metric_name = "custom.transaction.time"
    aggregation = "MEAN"
    threshold = 500
  }
  
  sli_entity = {
    website_time_based = {
      website_id = "business-app"
      beacon_type = "custom"
    }
  }
}
```

### Page Change SLI (SPA)

```hcl
resource "instana_sli_config" "spa_navigation" {
  name = "SPA Navigation Performance"
  
  metric_configuration = {
    metric_name = "page.change.time"
    aggregation = "P90"
    threshold = 500
  }
  
  sli_entity = {
    website_time_based = {
      website_id = "spa-website"
      beacon_type = "pageChange"
    }
  }
}
```

### Multi-Environment SLI Setup

```hcl
locals {
  environments = {
    production = {
      app_id = "prod-app-id"
      threshold = 500
      aggregation = "P95"
    }
    staging = {
      app_id = "staging-app-id"
      threshold = 1000
      aggregation = "P90"
    }
  }
}

resource "instana_sli_config" "env_latency" {
  for_each = local.environments

  name = "${each.key} API Latency"
  
  metric_configuration = {
    metric_name = "latency"
    aggregation = each.value.aggregation
    threshold = each.value.threshold
  }
  
  sli_entity = {
    application_time_based = {
      application_id = each.value.app_id
      boundary_scope = "INBOUND"
    }
  }
}
```

### Complex Event Filter Example

```hcl
resource "instana_sli_config" "complex_availability" {
  name = "Complex Availability SLI"
  
  sli_entity = {
    application_event_based = {
      application_id = "complex-app"
      boundary_scope = "INBOUND"
      good_event_filter_expression = "(call.http.status@na EQUALS 200 OR call.http.status@na EQUALS 201) AND call.error.count@na EQUALS 0"
      bad_event_filter_expression = "call.http.status@na GREATER_OR_EQUAL_THAN 400 OR call.error.count@na GREATER_THAN 0"
    }
  }
}
```

### Different Aggregation Types

```hcl
# P99 Latency
resource "instana_sli_config" "p99_latency" {
  name = "P99 Latency SLI"
  
  metric_configuration = {
    metric_name = "latency"
    aggregation = "P99"
    threshold = 2000
  }
  
  sli_entity = {
    application_time_based = {
      application_id = "critical-app"
      boundary_scope = "INBOUND"
    }
  }
}

# Mean Latency
resource "instana_sli_config" "mean_latency" {
  name = "Mean Latency SLI"
  
  metric_configuration = {
    metric_name = "latency"
    aggregation = "MEAN"
    threshold = 300
  }
  
  sli_entity = {
    application_time_based = {
      application_id = "standard-app"
      boundary_scope = "ALL"
    }
  }
}

# Max Latency
resource "instana_sli_config" "max_latency" {
  name = "Max Latency SLI"
  
  metric_configuration = {
    metric_name = "latency"
    aggregation = "MAX"
    threshold = 5000
  }
  
  sli_entity = {
    application_time_based = {
      application_id = "batch-app"
      boundary_scope = "ALL"
    }
  }
}
```

### Error Rate SLI

```hcl
resource "instana_sli_config" "error_rate" {
  name = "Application Error Rate"
  
  metric_configuration = {
    metric_name = "errors"
    aggregation = "SUM"
    threshold = 10
  }
  
  sli_entity = {
    application_time_based = {
      application_id = "monitored-app"
      boundary_scope = "INBOUND"
    }
  }
}
```

### Throughput SLI

```hcl
resource "instana_sli_config" "throughput" {
  name = "API Throughput"
  
  metric_configuration = {
    metric_name = "calls"
    aggregation = "PER_SECOND"
    threshold = 100
  }
  
  sli_entity = {
    application_time_based = {
      application_id = "high-traffic-app"
      boundary_scope = "INBOUND"
    }
  }
}
```

## Argument Reference

* `name` - Required - The name of the SLI configuration (max 256 characters)
* `initial_evaluation_timestamp` - Optional - The initial evaluation timestamp for the SLI config (Unix timestamp in milliseconds)
* `metric_configuration` - Optional - Configuration block to describe the metric the SLI config is based on [Details](#metric-configuration-reference). Required for `application_time_based` and `website_time_based` SLI entities
* `sli_entity` - Required - Configuration block to describe the entity the SLI config is based on [Details](#sli-entity-reference)

### Metric Configuration Reference

* `metric_name` - Required - Name of the metric to monitor
* `aggregation` - Required - The aggregation type for the metric configuration. Allowed values: `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`, `PER_SECOND`
* `threshold` - Required - Threshold value for the metric configuration (must be at least 0.000001)

### SLI Entity Reference

Exactly one of the elements below must be configured:

* `application_event_based` - Optional - Event-based SLI entity configuration for applications [Details](#application-event-based-sli-entity-reference)
* `application_time_based` - Optional - Time-based SLI entity configuration for applications [Details](#application-time-based-sli-entity-reference)
* `website_event_based` - Optional - Event-based SLI entity configuration for websites [Details](#website-event-based-sli-entity-reference)
* `website_time_based` - Optional - Time-based SLI entity configuration for websites [Details](#website-time-based-sli-entity-reference)

#### Application Event-Based SLI Entity Reference

Used for availability and success rate SLIs based on call/event filtering.

* `application_id` - Required - The application ID of the entity
* `boundary_scope` - Required - The boundary scope of the entity. Allowed values: `ALL`, `INBOUND`
* `good_event_filter_expression` - Required - Tag filter expression to match good events/calls [Details](#tag-filter-expression-reference)
* `bad_event_filter_expression` - Required - Tag filter expression to match bad events/calls [Details](#tag-filter-expression-reference)
* `service_id` - Optional - The service ID to scope the SLI to a specific service
* `endpoint_id` - Optional - The endpoint ID to scope the SLI to a specific endpoint
* `include_internal` - Optional - Flag to indicate whether internal calls are included in the scope. Default: `false`
* `include_synthetic` - Optional - Flag to indicate whether synthetic calls are included in the scope. Default: `false`

#### Application Time-Based SLI Entity Reference

Used for latency and performance SLIs based on metric aggregation.

* `application_id` - Required - The application ID of the entity
* `boundary_scope` - Required - The boundary scope of the entity. Allowed values: `ALL`, `INBOUND`
* `service_id` - Optional - The service ID to scope the SLI to a specific service
* `endpoint_id` - Optional - The endpoint ID to scope the SLI to a specific endpoint

#### Website Event-Based SLI Entity Reference

Used for website availability and error rate SLIs based on beacon filtering.

* `website_id` - Required - The website ID of the entity
* `beacon_type` - Required - The beacon type of the entity. Allowed values: `pageLoad`, `resourceLoad`, `httpRequest`, `error`, `custom`, `pageChange`
* `good_event_filter_expression` - Required - Tag filter expression to match good events/beacons [Details](#tag-filter-expression-reference)
* `bad_event_filter_expression` - Required - Tag filter expression to match bad events/beacons [Details](#tag-filter-expression-reference)

#### Website Time-Based SLI Entity Reference

Used for website performance SLIs based on metric aggregation.

* `website_id` - Required - The website ID of the entity
* `beacon_type` - Required - The beacon type of the entity. Allowed values: `pageLoad`, `resourceLoad`, `httpRequest`, `error`, `custom`, `pageChange`
* `filter_expression` - Optional - Tag filter expression to match events/beacons [Details](#tag-filter-expression-reference)

#### Tag Filter Expression Reference

The **tag_filter** defines which calls/events/beacons should be included. It supports:

* Logical AND and/or logical OR conjunctions (AND has higher precedence than OR)
* Comparison operators: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `LESS_THAN`, `GREATER_THAN`
* Unary operators: `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `NOT_BLANK`

The **tag_filter** is defined by the following eBNF:

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

## Attributes Reference

* `id` - The ID of the SLI configuration

## Import

SLI Configs can be imported using the `id`, e.g.:

```bash
$ terraform import instana_sli_config.my_sli 60845e4e5e6b9cf8fc2868da
```

## Notes

* **IMPORTANT**: SLI Configurations are immutable and cannot be updated. Any change will result in an error
* To modify an SLI, you must create a new SLI configuration and delete the old one
* The ID is auto-generated by Instana
* Use `application_time_based` for latency and performance metrics
* Use `application_event_based` for availability and success rate metrics
* Use `website_time_based` for website performance metrics
* Use `website_event_based` for website availability and error rate metrics
* Tag filter expressions support complex logical conditions
* The `boundary_scope` determines whether to monitor all calls or only inbound calls
* SLIs are typically used in conjunction with SLO configurations
* The `initial_evaluation_timestamp` can be used to backfill SLI data
