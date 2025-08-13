# Infrastructure Alert Configuration Resource

Management of Infrastructure alert configurations (Infrastructure Smart Alerts).

API Documentation: <https://instana.github.io/openapi/#tag/Infrastructure-Alert-Configuration>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

## Example Usage

```hcl
resource "instana_infra_alert_config" "example" {
  name = "test-alert"
  description = "test-alert-description"
  alert_channels {
    warning = ["alert-channel-id-1"]
    critical = ["alert-channel-id-2"]
  }
  group_by = ["metricId"]
  granularity = 600000
  evaluation_type = "CUSTOM"
  tag_filter = "host.fqdn@na STARTS_WITH 'fooBar'"

  rules {
    generic_rule {
      metric_name = "cpu\\.(nice|user|sys|wait)"
      entity_type = "host"
      aggregation = "MEAN"
      cross_series_aggregation = "SUM"
      regex = true
      threshold_operator = ">="

      threshold {
        
        critical {
          static {
            value = 5.0
          }
        }
        
        warning {
          static {
            value = 3.0
          }
        }
        
      }
    }
  }

  time_threshold {
    violations_in_sequence {
      time_window = 600000
    }
  }

  custom_payload_field {
    key = "test1"
    value = "foo"
  }

  custom_payload_field {
    key = "test2"
    dynamic_value {
      key = "dynamic-value-key"
      tag_name = "dynamic-value-tag-name"
    }
  }
}
```

## Argument Reference

* `name` - Required - The name for the infrastructure alert configuration
* `description` - Required - The description text of the infrastructure alert config
* `alert_channels` - Optional - Set of alert channel IDs associated with the severity. [Details](#alert-channels-reference)
* `group_by` - Optional - The grouping tags used to group the metric results.
* `granularity` - Optional - default `600000` - The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used. Allowed values: `300000`, `600000`, `900000`, `1200000`, `800000`
* `evaluation_type` - Optional - default `CUSTOM` - The evaluation type of the infrastructure alert config. Allowed values: `CUSTOM`, `PER_ENTITY`. [Details](#evaluation-type-reference)
* `tag_filter` - Optional - The tag filter of the global application alert config. [Details](#tag-filter-argument-reference)
* `rules` - Required - A list of rules where each rule is associated with multiple thresholds and their corresponding severity levels. This enables more complex alert configurations with validations to ensure consistent and logical threshold-severity combinations. [Details](#rules-argument-reference)
* `time_threshold` - Required - Indicates the type of violation of the defined threshold.  [Details](#time-threshold-argument-reference)
* `custom_payload_filed` - Optional - An optional list of custom payload fields (static key/value pairs added to the event).  [Details](#custom-payload-field-argument-reference)

### Alert Channels Reference

* `warning` - Optional - List of alert channel IDs associated with the warning severity 
* `critical` - Optional - List of alert channel IDs associated with the critical severity

### Evaluation Type Reference

* `CUSTOM`: This is a default option. With this, you can combine all metrics in scope into a single metric per group.
* `PER_ENTITY`: To monitor each metric individually and trigger alerts for each individual entity.

### Tag Filter Argument Reference
The **tag_filter** defines which entities should be included into the application. It supports:

* logical AND and/or logical OR conjunctions whereas AND has higher precedence then OR
* comparison operators EQUALS, NOT_EQUAL, CONTAINS | NOT_CONTAIN, STARTS_WITH, ENDS_WITH, NOT_STARTS_WITH, NOT_ENDS_WITH, GREATER_OR_EQUAL_THAN, LESS_OR_EQUAL_THAN, LESS_THAN, GREATER_THAN
* unary operators IS_EMPTY, NOT_EMPTY, IS_BLANK, NOT_BLANK.

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

### Rules Argument Reference

* `generic_rule` - Required - A generic rule based on custom aggregated metric [Details](#generic-rule-argument-reference)

#### Generic Rule Argument Reference 

* `metric_name` - Required - The metric name of the infrastructure alert rule
* `entity_type` - Required - The entity type of the infrastructure alert rule
* `aggregation` - Required - The aggregation function of the infra alert rule. Supported values `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `SUM`, `PER_SECOND`
* `cross_series_aggregation` - Required - Cross-series aggregation function of the infra alert rule. Supported values `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `SUM`
* `regex` - Optional - Indicates if the given metric name follows regex pattern or not
* `threshold_operator` - Required - The operator which will be applied to evaluate the threshold
* `threshold` - Required - Indicates the type of threshold associated with given severity this alert rule is evaluated on.  [Details](#threshold-rule-argument-reference)
* `time_threshold` - Required - Indicates the type of violation of the defined threshold.  [Details](#time-threshold-argument-reference)

#### Threshold Rule Argument Reference

At least one of the elements below must be configured

* `warning` - Optional - Threshold associated with the warning severity [Details](#threshold-argument-reference)
* `critical` - Optional - Threshold associated with the critical severity [Details](#threshold-argument-reference)

##### Threshold Argument Reference

* `static` - Required - Static threshold definition. [Details](#static-threshold-argument-reference)

###### Static Threshold Argument Reference

* `operator` - Required - The operator which will be applied to evaluate the threshold. Supported values: `>`, `>=`, `<`, `<=`
* `value` - Optional - The value of the static threshold

### Time Threshold Argument Reference

* `violations_in_sequence` - Required - Time threshold base on violations in sequence. [Details](#violations-in-sequence-time-threshold-argument-reference)

#### Violations In Sequence Time Threshold Argument Reference

* `time_window` - Required - The time window if the time threshold

### Custom Payload Field Argument Reference

* `key` - Required - The key of the custom payload field
* `value` - Optional - The static string value of the custom payload field. Either `value` or `dynamic_value` must be defined.
* `dynamic_value` - Optional - The dynamic value of the custom payload field [Details](#dynamic-custom-payload-field-value). Either `value` or `dynamic_value` must be defined.

#### Dynamic Custom Payload Field Value
* `key` - Optional - The key of the tag which should be added to the payload
* `tag_name` - Required - The name of the tag which should be added to the payload
