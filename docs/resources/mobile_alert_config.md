# Mobile Alert Configuration Resource

This resource manages mobile app alert configurations in Instana. Mobile alerts monitor mobile application metrics and trigger notifications based on defined thresholds and conditions.

API Documentation: <https://developer.ibm.com/apis/catalog/instana--instana-rest-api/api/API--instana--instana-rest-api-documentation#createMobileAppAlertConfig>

## Example Usage

### Mobile App slowness Alert

```hcl
resource "instana_mobile_alert_config" "slowness" {
  name          = "tf High Crash Rate - $${severity}" # Use double $$ to define placeholders
  description   = "Alert when mobile app crash rate exceeds thresholds"
  mobile_app_id = "mobile-app-id" # replace with valid mobile App Id
  triggering    = true
  granularity   = 600000
  tag_filter    = "mobileBeacon.http.status@na EQUALS '500'"
  enabled = false
  rules = [{
    rule = {
      alert_type  = "slowness"
      metric_name = "httpLatency"
      aggregation = "P90"
    }
    threshold_operator = ">="
    threshold = {
      warning = {
        static = {
          value = 50
        }
      }
      critical = {
        static = {
          value = 100
        }
      }
    }
  }]

  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }

  custom_payload_field = [
    {
      key   = "team"
      value = "mobile-team"
    }
  ]
}
```

### Mobile App Crash Alert 

```hcl
resource "instana_mobile_alert_config" "crash_rate" {
  name          = "Mobile Error Rate Alert tf"
  description   = "Alert based on error rate and user impact"
  mobile_app_id = "mobile-app-id" # replace with valid mobile App Id
  triggering    = true
  granularity   = 300000

   alert_channels = {
    warning  = ["channel-id-1"] # replace with actual channel Ids
    critical = ["channel-id-2"]
  }

  rules = [{
    rule = {
      alert_type  = "crash"
      metric_name = "crashAffectedSessionRate"
      aggregation = "MEAN"
    }
    threshold_operator = ">="
    threshold = {
      critical = {
        static = {
          value = 10
        }
      }
    }
  }]

  time_threshold = {
    user_impact_of_violations_in_sequence = {
      time_window = 900000
      users       = 100
      percentage  = 1.0
    }
  }

  custom_payload_field = [
    {
      key = "applicationName"
      dynamic_value = {
        tag_name = "mobileBeacon.mobileApp.name"
      }
    }
  ]
}
```

### Mobile App Throughput with Violations in Period

```hcl
resource "instana_mobile_alert_config" "throughput" {
  name          = "Low Mobile App Throughput tf"
  description   = "Alert when throughput drops below threshold"
  mobile_app_id = "mobile-app-id" # replace with valid mobile App Id
  triggering    = false
  granularity   = 600000

  alert_channels = {
    warning  = ["channel-id-1"] # replace with actual channel Ids
    critical = ["channel-id-2"]
  }

  rules =[ {
    rule = {
      alert_type  = "throughput"
      metric_name = "views"
      aggregation = "SUM"
    }
    threshold_operator = "<"
    threshold = {
      warning = {
        static = {
          value = 1000
        }
      }
    }
  }]

  time_threshold = {
    violations_in_period = {
      time_window = 1800000
      violations  = 3
    }
  }

  custom_payload_field = [
    {
      key   = "environment"
      value = "production"
    }
  ]
}
```

## Generating Configuration from Existing Resources

If you have already created a mobile alert configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your mobile alert configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_mobile_alert_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual mobile alert configuration ID from Instana

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

* `name` - Required - The name of the mobile alert configuration (max 256 characters). Used as a template for the title of alert/event notifications
* `description` - Required - The description of the mobile alert configuration (max 65536 characters). Used as a template for the description of alert/event notifications
* `mobile_app_id` - Required - ID of the mobile app that this alert configuration is applied to (max 64 characters)
* `tag_filter` - Required - The tag filter expression for the mobile alert configuration. Defines which mobile app entities this alert applies to [Details](#tag-filter-argument-reference)
* `time_threshold` - Required - The type of threshold to define the criteria when the event and alert triggers and resolves [Details](#time-threshold-argument-reference)
* `triggering` - Optional - Flag to indicate whether an Incident is also triggered. Default: `false`
* `alert_channels` - Optional - Set of alert channel IDs associated with the severity [Details](#alert-channels-reference)
* `granularity` - Optional - The evaluation granularity in milliseconds. Default: `600000` (10 minutes). Allowed values: `60000`, `300000`, `600000`, `900000`, `1200000`, `1800000`
* `grace_period` - Optional - The duration in milliseconds for which an alert remains open after conditions are no longer violated. The alert auto-closes once the grace period expires
* `custom_payload_field` - Optional - A list of custom payload fields (static key/value pairs or dynamic tag values added to the event). Maximum 20 fields [Details](#custom-payload-field-argument-reference)
* `rules` - Optional - A rule configuration with thresholds and their corresponding severity levels. Maximum 1 rule [Details](#rules-argument-reference)

### Alert Channels Reference

* `warning` - Optional - List of alert channel IDs associated with the warning severity
* `critical` - Optional - List of alert channel IDs associated with the critical severity

### Tag Filter Argument Reference

The **tag_filter** defines which mobile app entities should be included in the alert scope. It supports:

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

### Rules Argument Reference

* `rule` - Required - The mobile app alert rule configuration [Details](#rule-argument-reference)
* `threshold_operator` - Required - The operator to apply for threshold comparison. Supported values: `>`, `>=`, `<`, `<=`
* `threshold` - Required - Threshold configuration for different severity levels [Details](#threshold-rule-argument-reference)

#### Rule Argument Reference

* `alert_type` - Required - The type of alert rule. Supported values: `slowness`, `errors`, `throughput`
* `metric_name` - Required - The metric name of the mobile alert rule
* `aggregation` - Optional - The aggregation function of the mobile alert rule. Supported values: `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`

#### Threshold Rule Argument Reference

At least one of the elements below must be configured:

* `warning` - Optional - Threshold associated with the warning severity [Details](#threshold-argument-reference)
* `critical` - Optional - Threshold associated with the critical severity [Details](#threshold-argument-reference)

##### Threshold Argument Reference

Exactly one of the following must be specified:

* `static` - Static threshold definition [Details](#static-threshold-argument-reference)
* `adaptive_baseline` - Adaptive baseline threshold definition [Details](#adaptive-baseline-threshold-argument-reference)
* `historic_baseline` - Historic baseline threshold definition [Details](#historic-baseline-threshold-argument-reference)

###### Static Threshold Argument Reference

* `value` - Required - The numeric value for the static threshold

###### Adaptive Baseline Threshold Argument Reference

* `deviation_factor` - Required - The numeric value for the deviation factor
* `adaptability` - Required - The numeric value for the adaptability
* `seasonality` - Required - Value for the seasonality

###### Historic Baseline Threshold Argument Reference

* `seasonality` - Required - Value for the seasonality
* `baseline` - Required - The baseline value

### Time Threshold Argument Reference

Exactly one of the following must be specified:

* `violations_in_sequence` - Time threshold based on violations in sequence [Details](#violations-in-sequence-time-threshold-argument-reference)
* `user_impact_of_violations_in_sequence` - Time threshold based on user impact of violations in sequence [Details](#user-impact-violations-in-sequence-time-threshold-argument-reference)
* `violations_in_period` - Time threshold based on violations in period [Details](#violations-in-period-time-threshold-argument-reference)

#### Violations In Sequence Time Threshold Argument Reference

* `time_window` - Required - The time window in milliseconds. Default: `600000`

#### User Impact Violations In Sequence Time Threshold Argument Reference

* `time_window` - Required - The time window in milliseconds
* `users` - Optional - The number of impacted users
* `percentage` - Optional - The percentage of impacted users
* `impact_measurement_method` - Optional - The method to measure user impact. Supported values: `AGGREGATED`, `PER_WINDOW`

#### Violations In Period Time Threshold Argument Reference

* `time_window` - Required - The time window in milliseconds
* `violations` - Required - The number of violations that must appear in the period

### Custom Payload Field Argument Reference

* `key` - Required - The key of the custom payload field
* `value` - Optional - The static string value of the custom payload field. Either `value` or `dynamic_value` must be defined
* `dynamic_value` - Optional - The dynamic value of the custom payload field [Details](#dynamic-custom-payload-field-value). Either `value` or `dynamic_value` must be defined

#### Dynamic Custom Payload Field Value

* `key` - Optional - The key of the tag which should be added to the payload
* `tag_name` - Required - The name of the tag which should be added to the payload

## Attributes Reference

* `id` - The ID of the mobile alert configuration
