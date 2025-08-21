# SLO Configuration

An SLO(Service Level Objective) is a specific, measurable target that defines the expected level of performance, 
reliability, or availability of a service, agreed upon between a service provider and its users or customers.
For instance, an SLO could state that a specific SLI (Service Level Indicator), such as availability, must reach 99.9% 
over a set period.

API Documentation: <https://instana.github.io/openapi/#operation/createSloConfig>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

## Example Usage
Creating an `application` SLO with a `timebased latency` indicator and a `rolling` time window 

```hcl
resource "instana_slo_config" "slo_1" {
  name = "app_timebased_latency_rolling"
  target = 0.95
  tags = ["terraform", "app", "timebased", "latency", "fixed"]
  entity {
    application {
      application_id = "instana_application_config_id"
      boundary_scope = "ALL"
      include_internal = false
      include_synthetic = false
      filter_expression = "AND"
    }
  }
  indicator {
     time_based_latency {
       threshold = 13.1
       aggregation = "MEAN"
     }
  }
  time_window {
    rolling {
      duration = 1
      duration_unit = "day"
      timezone = "UTC"
    }
  }
}
```

Creating a `website` SLO with a `timebased availability` indicator and a `fixed` time window 

```hcl
resource "instana_slo_config" "website_3" {
  name = "website_timebased_availability_fixed"
  target = 0.91
  tags = ["terraform", "web", "timebased", "availability", "fixed"]
  entity {
    website {
     website_id = "instana_website_monitoring_config_id"
     beacon_type = "httpRequest"
     filter_expression = "AND"
    }
  }
  indicator {
    time_based_availability {
      threshold = 14.7
      aggregation = "MEAN"
    }
  }
   time_window {
     fixed {
       duration = 1
       duration_unit = "day"
       timezone = "Europe/Dublin"
       start_timestamp = var.fixed_timewindow_start_timestamp
     }
   }
}
```

Creating a `synthetic` monitoring SLO with the `all traffic` indicator and a `rolling` time window 

```hcl
resource "instana_slo_config" "synthetic_r_6" {
  name = "synthetic_traffic_all_rolling"
  target = 0.91
  tags = ["terraform", "synthetic", "traffic", "all", "rolling-time-window"]
  entity {
     synthetic {
       synthetic_test_ids = ["DrMyeGl08w79poguQ3mhH", "sYDtb2slIIolfXhPBnodSz" ]
       filter_expression = "AND"
     } 
  }
  indicator {
    traffic {
      traffic_type = "all"
      threshold = 14
      aggregation = "SUM"
    }
  }
   time_window {
     rolling {
       duration = 1
       duration_unit = "day"
       timezone = ""
     }
   }
}
``` 

## Argument Reference

* `name` - Required - The name of the SLO configuration. Must be a string
* `target` - Required - The target SLO value (e.g., 0.99 for 99% availability). Must be a float.
* `tags` - Optional - A list of tags associated with the SLO configuration. 
* `entity` element - Required - A resource block describing the entity the SLO configuration is based on. [Details](#entity-reference)
* `indicator` element - Required - A resource block describing the indicator (metric) the SLO configuration is based on. [Details](#indicator-reference)
* `time_window` element - Required - A resource block describing the time window for the SLO configuration. [Details](#time-window-reference)

### Entity element Reference

One of the elements below must be configured:

* `application` - Application-based SLO entity configuration. [Details](#application-slo-entity-reference)
* `website` - Website-based SLO entity configuration. [Details](#website-slo-entity-reference)
* `synthetic` - Synthetic-based SLO entity configuration. [Details](#synthetic-slo-entity-reference)

#### Application SLO Entity Reference

* `application_id` - Required - The application ID of the entity.
* `boundary_scope` - Required - The boundary scope of the entity. Allowed values: `ALL`, `INBOUND`.
* `filter_expression` - Optional - Tag filter expression to match events/calls. [Details](#tag-filter-expression-reference)
* `include_internal` - Optional - Flag to indicate whether internal calls are included in the scope. Must be a boolean. Defaults to `false`.
* `include_synthetic` - Optional - Flag to indicate whether synthetic calls are included in the scope. Must be a boolean. Defaults to `false`.
* `service_id` - Optional - The service ID of the entity.
* `endpoint_id` - Optional - The endpoint ID of the entity.

#### Website SLO Entity Reference

* `website_id` - Required - The website ID of the entity.
* `filter_expression` - Optional - Tag filter expression to match events/calls. [Details](#tag-filter-expression-reference)
* `beacon_type` - Required - The beacon type of the entity. Allowed value:  `httpRequest`.

#### Synthetic SLO Entity Reference

* `synthetic_test_ids` - Required - A list of synthetic test IDs for the entity. Must contain at least one ID.
* `filter_expression` - Optional - Tag filter expression to match events/calls. [Details](#tag-filter-expression-reference)

### Indicator element Reference

Exactly one of the elements below must be configured:

* `time_based_latency` - Optional - Time-based latency indicator configuration. [Details](#time-based-latency-indicator-reference)
* `event_based_latency` - Optional - Event-based latency indicator configuration. [Details](#event-based-latency-indicator-reference)
* `time_based_availability` - Optional - Time-based availability indicator configuration. [Details](#time-based-availability-indicator-reference)
* `event_based_availability` - Optional - Event-based availability indicator configuration. [Details](#event-based-availability-indicator-reference)
* `traffic` - Optional - Traffic indicator configuration. [Details](#traffic-indicator-reference)
* `custom` - Optional - Custom indicator configuration. [Details](#custom-indicator-reference)

#### Time-based Latency Indicator Reference

* `threshold` - Required - The threshold for the metric. Must be a float greater than 0.
* `aggregation` - Required - The aggregation type for the metric. Allowed values: `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`.

#### Event-based Latency Indicator Reference

* `threshold` - Required - The threshold for the metric. Must be a float greater than 0.

#### Time-based Availability Indicator Reference

* `threshold` - Required - The threshold for the metric. Must be a float greater than 0.
* `aggregation` - Required - The aggregation type for the metric. Allowed value: `MEAN`

#### Event-based Availability Indicator Reference

No additional arguments are required for this indicator type.

#### Traffic Indicator Reference

* `traffic_type` - Required - The type of traffic to measure. Allowed values: `all`, `erroneous`.
* `threshold` - Required - The threshold for the metric. Must be a float greater than 0.

#### Custom Indicator Reference

* `good_event_filter_expression` - Required - Tag filter expression to match good events/calls. [Details](#tag-filter-expression-reference)
* `bad_event_filter_expression` - Optional - Tag filter expression to match bad events/calls. If this field is empty the opposit value of the good event is used as bad event filter. [Details](#tag-filter-expression-reference)

### Time Window element Reference

One of the elements below must be configured:

* `rolling` - Rolling time window configuration. [Details](#rolling-time-window-reference)
* `fixed` - Fixed time window configuration. [Details](#fixed-time-window-reference)

#### Rolling Time Window Reference

* `duration` - Required - The duration of the time window. Must be an integer.
* `duration_unit` - Required - The unit of the duration. Allowed values: `day`, `week`.
* `timezone` - Optional - The timezone associated with the SLO configuration. Defaults to `UTC` when not specified.

Fixed Time Window Reference

* `duration` - Required - The duration of the time window. Must be an integer.
* `duration_unit` - Required - The unit of the duration. Allowed values: `day`, `week`.
* `timezone` - Optional - The timezone associated with the SLO configuration. Defaults to `UTC` when not specified.
* `initial_evaluation_timestamp` - Required - the starting timestamp for the Fixed Time Window.

### Tag Filter Expression Reference

The `tag_filter` defines which calls/events should be included. It supports:

* Logical `AND` and/or logical `OR` conjunctions whereas `AND` has higher precedence than `OR`
* Comparison operators: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `LESS_THAN`, `GREATER_THAN`
* Unary operators: `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `NOT_BLANK`

The `tag_filter` is defined by the following eBNF:

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

## Import

SLO Configs can be imported using the `id`, e.g.:

```
$ terraform import instana_slo_config.example_slo 60845e4e5e6fbf76pf8fc2868da
```
