# SLO Correction Configuration Resource

Manages an SLO Correction Configuration in Instana.

SLO correction windows let you exclude specific time periods from SLO calculations, providing a more accurate measurement of your Service Level Objective (SLO) performance. Common scenarios include:

- **Planned maintenance periods**
- **Non-business hours** (such as weekends, holidays, or overnight)
- **Isolated incidents or events** that do not represent normal operations

Excluding these intervals helps prevent temporary or expected disruptions from distorting your SLO results. This leads to a more accurate view of your service reliability and enables better decision-making.

API Documentation: <https://instana.github.io/openapi/#tag/SLO-Correction-Configuration>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

### Syntax Changes Overview

The main change is in how the `scheduling` block is defined. In v6, it uses attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
```hcl
resource "instana_slo_correction_config" "example" {
  name = "Maintenance Window"
  
  scheduling {
    start_time = 1718000880000
    duration = 60
    duration_unit = "MINUTE"
    recurrent_rule = "FREQ=WEEKLY;BYDAY=SU"
  }
}
```

#### NEW (v6.x) Syntax:
```hcl
resource "instana_slo_correction_config" "example" {
  name = "Maintenance Window"
  
  scheduling = {
    start_time = 1718000880000
    duration = 60
    duration_unit = "minute"
    recurrent_rule = "FREQ=WEEKLY;BYDAY=SU"
  }
}
```

### Key Syntax Changes

1. **Scheduling**: `scheduling { }` → `scheduling = { }`
2. **Duration Unit**: Case-insensitive but lowercase recommended: `MINUTE` → `minute`
3. **All nested objects**: Use `= { }` syntax

## Example Usage

### One-Time Correction Window

Create a one-time correction window for planned maintenance:

```hcl
resource "instana_slo_correction_config" "one_time_maintenance" {
  name = "Planned Maintenance - June 2024"
  description = "Database upgrade maintenance window"
  active = true
  
  scheduling = {
    start_time = 1718000880000
    duration = 120
    duration_unit = "minute"
    recurrent_rule = ""
  }
  
  slo_ids = ["api-latency-slo", "api-availability-slo"]
  tags = ["maintenance", "database-upgrade"]
}
```

### Daily Recurring Correction

Exclude non-business hours every day:

```hcl
resource "instana_slo_correction_config" "nightly_maintenance" {
  name = "Nightly Maintenance Window"
  description = "Daily maintenance window from 2 AM to 4 AM"
  active = true
  
  scheduling = {
    start_time = 1720000000000
    duration = 2
    duration_unit = "hour"
    recurrent_rule = "FREQ=DAILY"
  }
  
  slo_ids = ["prod-slo-1", "prod-slo-2"]
  tags = ["nightly-maintenance", "automated"]
}
```

### Weekly Recurring Correction

Exclude specific days of the week:

```hcl
resource "instana_slo_correction_config" "weekend_exclusion" {
  name = "Weekend Exclusion"
  description = "Exclude weekends from SLO calculations"
  active = true
  
  scheduling = {
    start_time = 1719000000000
    duration = 24
    duration_unit = "hour"
    recurrent_rule = "FREQ=WEEKLY;BYDAY=SA,SU"
  }
  
  slo_ids = ["business-hours-slo"]
  tags = ["weekend", "non-business-hours"]
}
```

### Weekday Business Hours Only

Monitor only during business hours on weekdays:

```hcl
resource "instana_slo_correction_config" "business_hours" {
  name = "Business Hours Only"
  description = "Exclude non-business hours on weekdays"
  active = true
  
  scheduling = {
    start_time = 1720000000000
    duration = 16
    duration_unit = "hour"
    recurrent_rule = "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"
  }
  
  slo_ids = ["customer-facing-slo"]
  tags = ["business-hours", "weekday"]
}
```

### Monthly Recurring Correction

Exclude first day of each month for maintenance:

```hcl
resource "instana_slo_correction_config" "monthly_maintenance" {
  name = "Monthly Maintenance"
  description = "First day of month maintenance window"
  active = true
  
  scheduling = {
    start_time = 1721000000000
    duration = 4
    duration_unit = "hour"
    recurrent_rule = "FREQ=MONTHLY;BYMONTHDAY=1"
  }
  
  slo_ids = ["system-slo"]
  tags = ["monthly-maintenance", "scheduled"]
}
```

### Bi-Weekly Correction

Every other week maintenance:

```hcl
resource "instana_slo_correction_config" "biweekly_maintenance" {
  name = "Bi-Weekly Maintenance"
  description = "Maintenance every other Sunday"
  active = true
  
  scheduling = {
    start_time = 1720000000000
    duration = 6
    duration_unit = "hour"
    recurrent_rule = "FREQ=WEEKLY;INTERVAL=2;BYDAY=SU"
  }
  
  slo_ids = ["infrastructure-slo"]
  tags = ["biweekly", "maintenance"]
}
```

### Limited Occurrence Correction

Correction that runs only a specific number of times:

```hcl
resource "instana_slo_correction_config" "limited_correction" {
  name = "Limited Correction Window"
  description = "Runs only 5 times"
  active = true
  
  scheduling = {
    start_time = 1720000000000
    duration = 2
    duration_unit = "hour"
    recurrent_rule = "FREQ=WEEKLY;BYDAY=WE;COUNT=5"
  }
  
  slo_ids = ["migration-slo"]
  tags = ["temporary", "migration"]
}
```

### Time-Bounded Recurring Correction

Correction with an end date:

```hcl
resource "instana_slo_correction_config" "time_bounded" {
  name = "Q4 Maintenance Windows"
  description = "Weekly maintenance until end of Q4"
  active = true
  
  scheduling = {
    start_time = 1720000000000
    duration = 3
    duration_unit = "hour"
    recurrent_rule = "FREQ=WEEKLY;BYDAY=SA;UNTIL=20241231T235959Z"
  }
  
  slo_ids = ["quarterly-slo"]
  tags = ["q4-2024", "maintenance"]
}
```

### Holiday Exclusions

Exclude specific holidays:

```hcl
resource "instana_slo_correction_config" "holidays" {
  name = "Holiday Exclusions"
  description = "Exclude major holidays from SLO"
  active = true
  
  scheduling = {
    start_time = 1735689600000  # Christmas Day 2024
    duration = 24
    duration_unit = "hour"
    recurrent_rule = "FREQ=YEARLY;BYMONTH=12;BYMONTHDAY=25"
  }
  
  slo_ids = ["customer-service-slo"]
  tags = ["holiday", "christmas"]
}
```

### Multi-Day Correction Window

Correction spanning multiple days:

```hcl
resource "instana_slo_correction_config" "multi_day" {
  name = "Extended Maintenance"
  description = "3-day maintenance window"
  active = true
  
  scheduling = {
    start_time = 1720000000000
    duration = 3
    duration_unit = "day"
    recurrent_rule = ""
  }
  
  slo_ids = ["system-migration-slo"]
  tags = ["extended-maintenance", "migration"]
}
```

### Disabled Correction for Testing

Create a disabled correction for testing:

```hcl
resource "instana_slo_correction_config" "test_correction" {
  name = "Test Correction (Disabled)"
  description = "Test correction configuration - not active"
  active = false
  
  scheduling = {
    start_time = 1720000000000
    duration = 60
    duration_unit = "minute"
    recurrent_rule = ""
  }
  
  slo_ids = ["test-slo"]
  tags = ["test"]
}
```

### Multi-Environment Setup

Create corrections for different environments:

```hcl
locals {
  environments = {
    production = {
      slo_ids = ["prod-slo-1", "prod-slo-2"]
      duration = 2
      recurrence = "FREQ=WEEKLY;BYDAY=SU"
    }
    staging = {
      slo_ids = ["staging-slo-1"]
      duration = 4
      recurrence = "FREQ=DAILY"
    }
  }
}

resource "instana_slo_correction_config" "env_maintenance" {
  for_each = local.environments

  name = "${each.key} Maintenance Window"
  description = "Maintenance window for ${each.key}"
  active = true
  
  scheduling = {
    start_time = 1720000000000
    duration = each.value.duration
    duration_unit = "hour"
    recurrent_rule = each.value.recurrence
  }
  
  slo_ids = each.value.slo_ids
  tags = ["${each.key}", "maintenance"]
}
```

### Complex Recurrence Pattern

Advanced recurrence with multiple conditions:

```hcl
resource "instana_slo_correction_config" "complex_recurrence" {
  name = "Complex Maintenance Schedule"
  description = "First and third Monday of each month"
  active = true
  
  scheduling = {
    start_time = 1720000000000
    duration = 4
    duration_unit = "hour"
    recurrent_rule = "FREQ=MONTHLY;BYDAY=1MO,3MO"
  }
  
  slo_ids = ["monthly-maintenance-slo"]
  tags = ["complex-schedule", "monthly"]
}
```

## Argument Reference

* `name` - Required - Name of the SLO correction configuration (max 256 characters)
* `description` - Required - Description of the correction configuration
* `active` - Required - Boolean flag indicating whether the correction configuration is active
* `scheduling` - Required - Scheduling configuration for the correction window [Details](#scheduling-reference)
* `slo_ids` - Required - Set of SLO IDs to which this correction applies
* `tags` - Optional - Set of tags to associate with the correction configuration

### Scheduling Reference

* `start_time` - Required - Start time of the correction window in milliseconds since epoch (UTC)
* `duration` - Required - Duration of the correction window (integer)
* `duration_unit` - Required - Unit for duration. Allowed values: `millisecond`, `second`, `minute`, `hour`, `day`, `week`, `month`
* `recurrent_rule` - Optional - Recurrence rule in [iCalendar RFC 5545](https://icalendar.org/iCalendar-RFC-5545/3-8-5-3-recurrence-rule.html) format. Leave empty for non-recurring corrections
* `recurrent` - Optional - Computed boolean flag indicating if the correction is recurrent (automatically set based on `recurrent_rule`)

### Recurrent Rule Reference

The `recurrent_rule` argument allows you to define how often the correction window should repeat, using the [iCalendar RFC 5545](https://icalendar.org/iCalendar-RFC-5545/3-8-5-3-recurrence-rule.html) standard. This enables flexible scheduling for recurring corrections.

**Supported rule parts:**
- `FREQ` - Frequency of recurrence. Values: `DAILY`, `WEEKLY`, `MONTHLY`, `YEARLY`
- `INTERVAL` - Interval between recurrences (e.g., every 2 weeks: `INTERVAL=2`)
- `COUNT` - Number of occurrences (e.g., `COUNT=5` for 5 occurrences)
- `UNTIL` - End date/time for the recurrence in UTC (e.g., `UNTIL=20240630T235959Z`)
- `BYMONTH` - Specific months (e.g., `BYMONTH=1,7` for January and July)
- `BYDAY` - Specific days of the week (e.g., `BYDAY=MO,WE,FR` for Monday, Wednesday, Friday)
- `BYMONTHDAY` - Specific days of the month (e.g., `BYMONTHDAY=1,15` for the 1st and 15th)

**Common Examples:**
- `FREQ=WEEKLY;BYDAY=MO,WE,FR` - Every Monday, Wednesday, and Friday
- `FREQ=DAILY;INTERVAL=2` - Every other day
- `FREQ=MONTHLY;BYMONTHDAY=1` - On the first day of each month
- `FREQ=WEEKLY;COUNT=5` - Weekly, only 5 times
- `FREQ=DAILY;UNTIL=20240630T235959Z` - Daily until June 30, 2024
- `FREQ=MONTHLY;BYDAY=1MO` - First Monday of each month
- `FREQ=YEARLY;BYMONTH=12;BYMONTHDAY=25` - Every Christmas Day

Leave `recurrent_rule` empty for a one-time (non-recurring) correction window.

## Attributes Reference

* `id` - The ID of the SLO correction configuration

## Import

SLO correction configurations can be imported using the `id`, e.g.:

```bash
$ terraform import instana_slo_correction_config.example 60845e4e5e6b9cf8fc2868da
```

## Notes

* The ID is auto-generated by Instana
* Correction windows exclude time periods from SLO calculations
* Use `active = false` to temporarily disable a correction without deleting it
* The `start_time` must be in milliseconds since Unix epoch (UTC)
* Recurrent rules follow the iCalendar RFC 5545 standard
* Multiple SLOs can share the same correction window
* Tags help organize and categorize correction windows
* Correction windows are useful for planned maintenance, non-business hours, and known incidents
* The `recurrent` field is automatically computed based on whether `recurrent_rule` is set
