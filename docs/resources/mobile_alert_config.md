---
page_title: "instana_mobile_alert_config Resource - terraform-provider-instana"
subcategory: ""
description: |-
  Manages Mobile App Alert Configurations in Instana.
---

# instana_mobile_alert_config (Resource)

This resource manages Mobile App Alert Configurations in Instana. Mobile alert configurations allow you to define alert rules for mobile applications based on various metrics and thresholds.

## Example Usage

```terraform
resource "instana_mobile_alert_config" "example" {
  name            = "My Mobile Alert"
  description     = "Alert for mobile app performance issues"
  mobile_app_id   = "mobile-app-123"
  triggering      = false
  granularity     = 600000
  tag_filter      = "entity.type:mobileApp"

  custom_payload_field {
    type = "staticString"
    key  = "team"
    value = "mobile-team"
  }

  rules {
    rule {
      alert_type  = "slowness"
      metric_name = "mobile.app.crash.rate"
      aggregation = "MEAN"
    }
    
    threshold_operator = ">"
    
    threshold {
      warning {
        static {
          operator = ">="
          value    = 50
        }
      }
      
      critical {
        static {
          operator = ">="
          value    = 100
        }
      }
    }
  }

  time_threshold {
    violations_in_sequence {
      time_window = 600000
    }
  }

  alert_channels = {
    "5"  = ["alert-channel-id-1"]
    "10" = ["alert-channel-id-2"]
  }
}
```

## Example with User Impact Time Threshold

```terraform
resource "instana_mobile_alert_config" "user_impact_example" {
  name            = "Mobile User Impact Alert"
  description     = "Alert based on user impact"
  mobile_app_id   = "mobile-app-456"
  triggering      = true
  granularity     = 300000
  tag_filter      = "entity.type:mobileApp"

  custom_payload_field {
    type = "dynamic"
    key  = "app_version"
    value = "entity.app.version"
  }

  rules {
    rule {
      alert_type  = "errors"
      metric_name = "mobile.app.error.rate"
      aggregation = "SUM"
    }
    
    threshold_operator = ">="
    
    threshold {
      critical {
        static {
          operator = ">="
          value    = 10
        }
      }
    }
  }

  time_threshold {
    user_impact_of_violations_in_sequence {
      time_window = 900000
      users       = 100
      percentage  = 5.0
    }
  }
}
```

## Example with Violations in Period

```terraform
resource "instana_mobile_alert_config" "violations_period_example" {
  name            = "Mobile Violations Period Alert"
  description     = "Alert based on violations in period"
  mobile_app_id   = "mobile-app-789"
  triggering      = false
  granularity     = 600000
  tag_filter      = "entity.type:mobileApp"
  grace_period    = 300000

  custom_payload_field {
    type = "staticString"
    key  = "environment"
    value = "production"
  }

  rules {
    rule {
      alert_type  = "throughput"
      metric_name = "mobile.app.requests"
      aggregation = "P95"
    }
    
    threshold_operator = "<"
    
    threshold {
      warning {
        static {
          operator = "<="
          value    = 1000
        }
      }
    }
  }

  time_threshold {
    violations_in_period {
      time_window = 1800000
      violations  = 3
    }
  }
}
```

## Schema

### Required

- `name` (String) The name of the Mobile Alert Configuration. Used as a template for the title of alert/event notifications triggered by this Smart Alert configuration. Maximum length: 256 characters.
- `description` (String) The description of the Mobile Alert Configuration. Used as a template for the description of alert/event notifications triggered by this Smart Alert configuration. Maximum length: 65536 characters.
- `mobile_app_id` (String) ID of the mobile app that this Smart Alert configuration is applied to. Maximum length: 64 characters.
- `tag_filter` (String) The tag filter expression for the Mobile Alert Configuration. Defines which mobile app entities this alert applies to.
- `time_threshold` (Block) The type of threshold to define the criteria when the event and alert triggers and resolves. (see [below for nested schema](#nestedblock--time_threshold))

### Optional

- `severity` (Number) The severity of the alert when triggered (5 for Warning, 10 for Critical). **Deprecated** - use `rules` with `threshold` instead. Valid values: 5, 10.
- `triggering` (Boolean) Flag to indicate whether an Incident is also triggered or not. Default: `false`.
- `complete_tag_filter` (String) The complete tag filter expression for the Mobile Alert Configuration.
- `alert_channels` (Map of List of String) Set of alert channel IDs associated with the severity. The map key is the severity level ("5" for Warning, "10" for Critical), and the value is a list of alert channel IDs.
- `granularity` (Number) The evaluation granularity used for detection of violations of the defined threshold. Defines the size of the tumbling window used. Default: `600000` (10 minutes). Valid values: 60000, 300000, 600000, 900000, 1200000, 1800000.
- `grace_period` (Number) The duration (in milliseconds) for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires.
- `custom_payload_field` (Block List) Custom payload fields to send additional information in the alert notifications. Maximum 20 fields. (see [below for nested schema](#nestedblock--custom_payload_field))
- `rules` (Block List) A list of rules where each rule is associated with multiple thresholds and their corresponding severity levels. This enables more complex alert configurations with validations to ensure consistent and logical threshold-severity combinations. Maximum 1 rule. (see [below for nested schema](#nestedblock--rules))

### Read-Only

- `id` (String) The ID of the Mobile Alert Configuration.

<a id="nestedblock--custom_payload_field"></a>
### Nested Schema for `custom_payload_field`

Required:

- `type` (String) The type of custom payload field. Valid values: `staticString`, `dynamic`.
- `key` (String) A user-specified unique identifier or name for a custom payload entry.

Optional:

- `value` (String) The value of the custom payload field. For `staticString` type, this is a static string value. For `dynamic` type, this is a reference to a dynamic field (e.g., `entity.app.version`).

<a id="nestedblock--rules"></a>
### Nested Schema for `rules`

Optional:

- `rule` (Block) The mobile app alert rule configuration. (see [below for nested schema](#nestedblock--rules--rule))
- `threshold_operator` (String) The operator to apply for threshold comparison. Valid values: `>`, `>=`, `<`, `<=`.
- `threshold` (Block) Threshold configuration for different severity levels. (see [below for nested schema](#nestedblock--rules--threshold))

<a id="nestedblock--rules--rule"></a>
### Nested Schema for `rules.rule`

Required:

- `alert_type` (String) The type of alert rule (e.g., `slowness`, `errors`, `throughput`).
- `metric_name` (String) The metric name of the mobile alert rule.

Optional:

- `aggregation` (String) The aggregation function of the mobile alert rule. Valid values: `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`.

<a id="nestedblock--rules--threshold"></a>
### Nested Schema for `rules.threshold`

Optional:

- `warning` (Block) Warning threshold configuration. (see [below for nested schema](#nestedblock--rules--threshold--warning))
- `critical` (Block) Critical threshold configuration. (see [below for nested schema](#nestedblock--rules--threshold--critical))

<a id="nestedblock--rules--threshold--warning"></a>
<a id="nestedblock--rules--threshold--critical"></a>
### Nested Schema for `rules.threshold.warning` and `rules.threshold.critical`

Optional (exactly one must be specified):

- `static` (Block) Static threshold definition.
  - `operator` (String) Comparison operator for the static threshold. Valid values: `>`, `>=`, `<`, `<=`, `==`.
  - `value` (Number) The numeric value for the static threshold.

- `adaptive_baseline` (Block) Adaptive baseline threshold definition.
  - `deviation_factor` (Number) The numeric value for the deviation factor.
  - `adaptability` (Number) The numeric value for the adaptability.
  - `seasonality` (String) Value for the seasonality.

- `historic_baseline` (Block) Historic baseline threshold definition.
  - `seasonality` (String) Value for the seasonality.
  - `baseline` (Number) The baseline value.

<a id="nestedblock--time_threshold"></a>
### Nested Schema for `time_threshold`

Optional (exactly one must be specified):

- `violations_in_sequence` (Block) Time threshold based on violations in sequence.
  - `time_window` (Number) The time window (in milliseconds) of the time threshold. Default: `600000`.

- `user_impact_of_violations_in_sequence` (Block) Time threshold based on user impact of violations in sequence.
  - `time_window` (Number) The time window (in milliseconds) of the time threshold.
  - `users` (Number) The number of impacted users.
  - `percentage` (Number) The percentage of impacted users.

- `violations_in_period` (Block) Time threshold based on violations in period.
  - `time_window` (Number) The time window (in milliseconds) of the time threshold.
  - `violations` (Number) The number of violations that must appear in the period.

## Import

Mobile Alert Configurations can be imported using the alert configuration ID:

```shell
terraform import instana_mobile_alert_config.example alert-config-id-here