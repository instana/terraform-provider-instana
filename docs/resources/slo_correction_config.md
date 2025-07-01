# SLO Correction Configuration
Manages an SLO Correction Configuration in Instana.

SLO correction windows let you exclude specific time periods from SLO calculations, providing a more accurate measurement of your Service Level Objective (SLO) performance. Common scenarios include:

- **Planned maintenance periods**
- **Non-business hours** (such as weekends, holidays, or overnight)
- **Isolated incidents or events** that do not represent normal operations

Excluding these intervals helps prevent temporary or expected disruptions from distorting your SLO results. This leads to a more accurate view of your service reliability and enables better decision-making.

## Example Usage

The following example shows how to create a basic SLO correction configuration with a one-time correction window:

```hcl
resource "instana_slo_correction_config" "correction_1" {
    name = "Example SLO Correction terraform correction test config"
    description = "This is a test SLO correction config."
    active = true
    scheduling {
        start_time = 1718000880000
        duration = 60
        duration_unit = "MINUTE"
        recurrent_rule = ""
    }
    slo_ids = ["instana-slo-id-1", "instana-slo-id-2"]
    tags = ["tag1", "tag2", "tag3"]
}
```
This example demonstrates a configuration with a recurring correction window using a recurrence rule:

```hcl
resource "instana_slo_correction_config" "correction_simple" {
    name        = "Simple SLO Correction"
    description = "A simple test SLO correction config."
    active      = false

    scheduling {
        start_time     = 1719000000000
        duration       = 120
        duration_unit  = "MINUTE"
        recurrent_rule = "FREQ=WEEKLY;BYDAY=MO,WE,FR"
    }

    slo_ids = [
        "instana-slo-id-1",
        "instana-slo-id-2"
    ]
    tags = [
        "env:staging",
        "team:backend"
    ]
}
```

This example shows a daily recurring correction window that runs for 1 hour every day:

```hcl
resource "instana_slo_correction_config" "correction_daily" {
    name        = "Daily SLO Correction"
    description = "Correction runs every day for 1 hour."
    active      = true

    scheduling {
        start_time     = 1720000000000
        duration       = 60
        duration_unit  = "MINUTE"
        recurrent_rule = "FREQ=DAILY"
    }

    slo_ids = [
        "instana-slo-id-3"
    ]
    tags = [
        "env:production"
    ]
}
```

This example demonstrates a monthly recurring correction window that runs on the first day of each month for 2 hours:

```hcl
resource "instana_slo_correction_config" "correction_monthly" {
    name        = "Monthly SLO Correction"
    description = "Correction runs on the first day of each month for 2 hours."
    active      = true

    scheduling {
        start_time     = 1721000000000
        duration       = 120
        duration_unit  = "MINUTE"
        recurrent_rule = "FREQ=MONTHLY;BYMONTHDAY=1"
    }

    slo_ids = [
        "instana-slo-id-4"
    ]
    tags = [
        "env:qa"
    ]
}
```

## Argument Reference
- `name` (String, **Required**) – Name of the SLO correction configuration.
- `description` (String, **Optional**) – Description of the correction configuration.
- `active` (Boolean, **Required**) – Whether the correction configuration is active.
- `scheduling` (Block, **Required**) – Scheduling configuration for the correction window:
    - `duration_unit` (String, **Required**) – Unit for duration. Supported values: `MINUTE`, `HOUR`, `DAY`.
    - `duration` (Number, **Required**) – Duration of the correction window.
    - `recurrent_rule` (String, **Optional**) – Recurrence rule in [iCalendar RFC 5545](https://icalendar.org/iCalendar-RFC-5545/3-8-5-3-recurrence-rule.html) format. Leave empty for non-recurring corrections.
    - `start_time` (Number, **Required**) – Start time of the correction window in milliseconds since epoch (UTC).
- `slo_ids` (List of String, **Required**) – List of SLO IDs to which this correction applies.
- `tags` (List of String, **Optional**) – List of tags to associate with the correction configuration.

---

### Recurrent Rule Reference

The `recurrent_rule` argument allows you to define how often the correction window should repeat, using the [iCalendar RFC 5545](https://icalendar.org/iCalendar-RFC-5545/3-8-5-3-recurrence-rule.html) standard. This enables flexible scheduling for recurring corrections.

**Supported rule parts:**
- `FREQ` – Frequency of recurrence (`DAILY`, `WEEKLY`, `MONTHLY`, etc.).
- `INTERVAL` – Interval between recurrences (e.g., every 2 weeks: `INTERVAL=2`).
- `COUNT` – Number of occurrences.
- `UNTIL` – End date/time for the recurrence (in UTC, e.g., `UNTIL=20240630T235959Z`).
- `BYMONTH` – Specific months (e.g., `BYMONTH=1,7` for January and July).
- `BYDAY` – Specific days of the week (e.g., `BYDAY=MO,WE,FR` for Monday, Wednesday, Friday).
- `BYMONTHDAY` – Specific days of the month (e.g., `BYMONTHDAY=1,15` for the 1st and 15th).

**Examples:**
- `FREQ=WEEKLY;BYDAY=MO,WE,FR` – Every Monday, Wednesday, and Friday.
- `FREQ=DAILY;INTERVAL=2` – Every other day.
- `FREQ=MONTHLY;BYMONTHDAY=1` – On the first day of each month.
- `FREQ=WEEKLY;COUNT=5` – Weekly, only 5 times.
- `FREQ=DAILY;UNTIL=20240630T235959Z` – Daily until June 30, 2024.

Leave `recurrent_rule` empty for a one-time (non-recurring) correction window.
