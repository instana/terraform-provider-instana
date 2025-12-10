# SLI Configuration Resource

Management of SLI configurations. A service level indicator (SLI) is the defined quantitative measure of one characteristic of the level of service that is provided to a customer. Common examples of such indicators are error rate or response latency of a service.

API Documentation: <https://instana.github.io/openapi/#operation/createSli>
**Note:** SLI Configurations cannot be changed. An update of the resource will result in an error. To update an SLI you need to create a new SLI and delete the old one.

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

The main changes are in how `metric_configuration` and `sli_entity` blocks are defined. In v6, all nested configurations use attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
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

#### NEW (v6.x) Syntax:
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


## Generating Configuration from Existing Resources

If you have already created a SLI configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your SLI configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_sli_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual SLI configuration ID from Instana

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
