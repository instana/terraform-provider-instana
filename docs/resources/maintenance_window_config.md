# Maintenance Window Configuration Resource

Manages maintenance window configurations in Instana. Maintenance windows allow you to suppress alert notifications during planned downtime or maintenance periods, preventing false alarms and reducing alert noise.

API Documentation: <https://developer.ibm.com/apis/catalog/instana--instana-rest-api/api/API--instana--instana-rest-api-documentation#getMaintenanceConfigV2>

## Example Usage

### One-Time Maintenance Window (Hours)

Schedule a single maintenance window for a planned deployment:

```hcl
resource "instana_maintenance_window_config" "deployment_window" {
  name  = "Planned Deployment - Production"
  query = "entity.zone:production"

  scheduling = {
    start = 2088055029000  # Unix timestamp in milliseconds
    type  = "ONE_TIME"

    duration = {
      amount = 2
      unit   = "HOURS"
    }

    rrule       = null
    timezone_id = null
  }

  tag_filter_expression         = null
  tag_filter_expression_enabled = false
}
```

### One-Time Maintenance Window (Minutes)

Schedule a short maintenance window for a quick patch:

```hcl
resource "instana_maintenance_window_config" "quick_patch" {
  name  = "Quick Patch Window"
  query = "entity.zone:production"

  scheduling = {
    start = 2088055029000  # Unix timestamp in milliseconds
    type  = "ONE_TIME"

    duration = {
      amount = 14
      unit   = "MINUTES"
    }

    rrule       = null
    timezone_id = null
  }

  tag_filter_expression         = null
  tag_filter_expression_enabled = false
}
```

### One-Time Maintenance Window (Days)

Schedule a full-day maintenance window for a major upgrade:

```hcl
resource "instana_maintenance_window_config" "major_upgrade" {
  name  = "Major Upgrade - Full Day"
  query = "entity.application.id:\"vcPQQcz-RP6VQsg9o2sDIw\"" # replace with valid application Id

  scheduling = {
    start = 2088055029000  # Unix timestamp in milliseconds
    type  = "ONE_TIME"

    duration = {
      amount = 1
      unit   = "DAYS"
    }

    rrule       = null
    timezone_id = null
  }

  tag_filter_expression         = null
  tag_filter_expression_enabled = false
}
```

### One-Time Maintenance Window with Tag Filter (Synthetic Tests)

Suppress alerts only for specific synthetic tests during a maintenance window:

```hcl
resource "instana_maintenance_window_config" "synthetic_test_maintenance" {
  name  = "One-time MW with filtering on synthetic tests"
  query = ""

  scheduling = {
    start = 2088055029000  # Unix timestamp in milliseconds
    type  = "ONE_TIME"

    duration = {
      amount = 1
      unit   = "HOURS"
    }

    rrule       = null
    timezone_id = null
  }

  tag_filter_expression         = "synthetic.testName@na EQUALS 'My Synthetic Test'"
  tag_filter_expression_enabled = true
}
```

### Recurrent Maintenance Window (Monthly - 4th Friday)

Schedule a recurring maintenance window on the 4th Friday of every month:

```hcl
resource "instana_maintenance_window_config" "monthly_friday_maintenance" {
  name  = "Monthly Maintenance - 4th Friday"
  query = "entity.zone:production"

  scheduling = {
    start       = 2088055029000  # Unix timestamp in milliseconds (first occurrence)
    type        = "RECURRENT"
    rrule       = "FREQ=MONTHLY;INTERVAL=1;UNTIL=20371114T000000Z;BYDAY=+4FR"
    timezone_id = "UTC"

    duration = {
      amount = 2
      unit   = "HOURS"
    }
  }

  tag_filter_expression         = null
  tag_filter_expression_enabled = false
}
```

### Recurrent Maintenance Window (Monthly - 4th Thursday, Limited Occurrences)

Schedule a recurring maintenance window on the 4th Thursday of every month, limited to 5 occurrences:

```hcl
resource "instana_maintenance_window_config" "monthly_thursday_limited" {
  name  = "Monthly Maintenance - 4th Thursday (5 occurrences)"
  query = ""

  scheduling = {
    start       = 2088055029000  # Unix timestamp in milliseconds (first occurrence)
    type        = "RECURRENT"
    rrule       = "FREQ=MONTHLY;INTERVAL=1;COUNT=5;BYDAY=+4TH"
    timezone_id = "UTC"

    duration = {
      amount = 3
      unit   = "HOURS"
    }
  }

  tag_filter_expression         = null
  tag_filter_expression_enabled = false
}
```

### Recurrent Maintenance Window (Weekly)

Schedule a recurring maintenance window every Sunday:

```hcl
resource "instana_maintenance_window_config" "weekly_maintenance" {
  name  = "Weekly Maintenance Window"
  query = "entity.zone:production"

  scheduling = {
    start       = 2088055029000  # Unix timestamp in milliseconds (first occurrence)
    type        = "RECURRENT"
    rrule       = "FREQ=WEEKLY;BYDAY=SU;INTERVAL=1"
    timezone_id = "America/New_York"

    duration = {
      amount = 4
      unit   = "HOURS"
    }
  }

  tag_filter_expression         = null
  tag_filter_expression_enabled = false
}
```

### Maintenance Window Scoped to Application Perspective

Target a specific application using its application ID:

```hcl
resource "instana_maintenance_window_config" "app_maintenance" {
  name  = "Application Maintenance Window"
  query = "entity.application.id:\"EtUaKj0sSu2tkQGY1D_MDg\"" # replace with valid applicaton Id

  scheduling = {
    start = 2088055029000  # Unix timestamp in milliseconds
    type  = "ONE_TIME"

    duration = {
      amount = 1
      unit   = "DAYS"
    }

    rrule       = null
    timezone_id = null
  }

  tag_filter_expression         = null
  tag_filter_expression_enabled = false
}
```

### Maintenance Window Scoped by Entity Tag

Target entities using a specific tag (e.g., deployment):

```hcl
resource "instana_maintenance_window_config" "deployment_maintenance" {
  name  = "Deployment Maintenance"
  query = "entity.tag:deployment=service-instance_f3c86baf-7d08-46d5-a01c-e67523d40820" # replace with valid tag query

  scheduling = {
    start = 2088055029000  # Unix timestamp in milliseconds
    type  = "ONE_TIME"

    duration = {
      amount = 14
      unit   = "MINUTES"
    }

    rrule       = null
    timezone_id = null
  }

  tag_filter_expression         = null
  tag_filter_expression_enabled = false
}
```

## Generating Configuration from Existing Resources

If you have already created a maintenance window configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your maintenance window configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_maintenance_window_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual maintenance window configuration ID from Instana

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

* `name` - Required - The name of the maintenance window configuration (1–256 characters)
* `query` - Required - Dynamic Focus Query (DFQ) that determines the scope of entities affected by the maintenance window (0–2048 characters). Use Instana's query syntax to target specific hosts, services, zones, or applications. Use an empty string `""` to apply to all entities
* `scheduling` - Required - Time scheduling configuration for the maintenance window [Details](#scheduling-reference)
* `tag_filter_expression_enabled` - Optional - Boolean flag to enable tag filter expression-based scoping. When `true`, the `tag_filter_expression` is used to further filter which alert notifications are muted
* `tag_filter_expression` - Optional - Tag filter expression used to filter alert notifications that will be muted during the maintenance window. Set to `null` when not used [Details](#tag-filter-expression-reference)

### Scheduling Reference

* `start` - Required - Start time of the maintenance window as a Unix timestamp in milliseconds (must be at least 1). For `RECURRENT` windows, this is the start time of the first occurrence
* `type` - Required - Type of maintenance window scheduling. Allowed values: `ONE_TIME`, `RECURRENT`
* `duration` - Required - Duration of each maintenance window occurrence [Details](#duration-reference)
* `rrule` - Optional - For `RECURRENT` maintenance windows, the recurrence rule following the [RRULE standard from the iCalendar specification (RFC 5545)](https://tools.ietf.org/html/rfc5545). Set to `null` for `ONE_TIME` windows. Required when `type` is `RECURRENT`. **Note:** Only the following RRULE tokens are supported by the Instana API: `BYDAY`, `BYMONTH`, `BYMONTHDAY`, `COUNT`, `FREQ`, `INTERVAL`, `UNTIL`
* `timezone_id` - Optional - Timezone ID for recurrent maintenance windows (e.g., `UTC`, `America/New_York`, `Europe/Berlin`). Set to `null` for `ONE_TIME` windows

### Duration Reference

* `amount` - Required - The numeric amount of time for the duration (must be at least 1)
* `unit` - Required - The unit of time for the duration. Allowed values: `MINUTES`, `HOURS`, `DAYS`

### Tag Filter Expression Reference

The **tag_filter_expression** defines which alert notifications should be muted during the maintenance window. It supports:

* Logical AND and/or logical OR conjunctions (AND has higher precedence than OR)
* Comparison operators: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `LESS_THAN`, `GREATER_THAN`
* Unary operators: `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `NOT_BLANK`

The **tag_filter_expression** is defined by the following eBNF:

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

**Common Tag Filter Examples:**
* `synthetic.testName@na EQUALS 'My Synthetic Test'` - Specific synthetic test by name
* `service.name@na EQUALS 'payment-service'` - Specific service
* `entity.zone@na EQUALS 'production'` - Specific zone
* `host.name@na STARTS_WITH 'prod-db'` - Hosts matching a prefix
* `kubernetes.namespace.name@na EQUALS 'default'` - Specific Kubernetes namespace

## Attributes Reference

* `id` - The ID of the maintenance window configuration (auto-generated by Instana)

## Import

Maintenance window configurations can be imported using the `id`, e.g.:

```bash
$ terraform import instana_maintenance_window_config.example 60845e4e5e6b9cf8fc2868da
```
