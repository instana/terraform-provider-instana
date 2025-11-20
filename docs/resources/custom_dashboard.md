# Custom Dashboard

Management of custom dashboards.

API Documentation: <https://instana.github.io/openapi/#tag/Custom-Dashboards>

---
## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block struture.

## Migration Guide (v5 to v6)

### Syntax Changes Overview

**OLD Syntax (SDK v2):**
```hcl
resource "instana_custom_dashboard" "example" {
  title = "Example Dashboard"
  
  access_rule { 
    access_type   = "READ_WRITE"
    relation_type = "USER"
    related_id    = "user-id-1"
  }
  
  access_rule {
    access_type   = "READ"
    relation_type = "GLOBAL"
  }
  
  widgets = file("${path.module}/widgets.json")
}
```

**NEW Syntax (Plugin Framework):**
```hcl
resource "instana_custom_dashboard" "example" {
  title = "Example Dashboard"
  
  access_rule = [{
    access_type   = "READ_WRITE"
    relation_type = "USER"
    related_id    = "user-id-1"
  }, {
    access_type   = "READ"
    relation_type = "GLOBAL"
  }]
  
  widgets = file("${path.module}/widgets.json")
}
```

**Key Changes:**
- `access_rule` now uses **list syntax** with `= [{ }]` instead of block syntax
- Enhanced validation for access types and relation types
- Improved JSON normalization for widgets
- Better state management

---

**Permissions Required:**
- `canCreatePublicCustomDashboards` (Creation of public custom dashboards)
- `canEditAllAccessibleCustomDashboards` (Management of all accessible custom dashboards)

The ID of the resource which is also used as unique identifier in Instana is auto generated!

## Example Usage

### Basic Dashboard

```hcl
resource "instana_custom_dashboard" "example" {
  title = "Production Monitoring Dashboard"
  
  access_rule = [{
    access_type   = "READ_WRITE"
    relation_type = "GLOBAL"
  }]
  
  widgets = file("${path.module}/dashboards/production.json")
}
```

### Dashboard with User-Specific Access

```hcl
resource "instana_custom_dashboard" "team_dashboard" {
  title = "Team Dashboard"
  
  access_rule = [{
    access_type   = "READ_WRITE"
    relation_type = "USER"
    related_id    = "user-123"
  }, {
    access_type   = "READ_WRITE"
    relation_type = "USER"
    related_id    = "user-456"
  }, {
    access_type   = "READ"
    relation_type = "GLOBAL"
  }]
  
  widgets = file("${path.module}/dashboards/team.json")
}
```

### Dashboard with Role-Based Access

```hcl
resource "instana_custom_dashboard" "ops_dashboard" {
  title = "Operations Dashboard"
  
  access_rule = [{
    access_type   = "READ_WRITE"
    relation_type = "ROLE"
    related_id    = instana_rbac_role.ops_admin.id
  }, {
    access_type   = "READ"
    relation_type = "ROLE"
    related_id    = instana_rbac_role.ops_viewer.id
  }]
  
  widgets = file("${path.module}/dashboards/operations.json")
}
```

### Dashboard with Team Access

```hcl
resource "instana_custom_dashboard" "platform_dashboard" {
  title = "Platform Engineering Dashboard"
  
  access_rule = [{
    access_type   = "READ_WRITE"
    relation_type = "TEAM"
    related_id    = "team-platform-engineering"
  }]
  
  widgets = file("${path.module}/dashboards/platform.json")
}
```

### Dashboard with API Token Access

```hcl
resource "instana_custom_dashboard" "api_dashboard" {
  title = "API Metrics Dashboard"
  
  access_rule = [{
    access_type   = "READ"
    relation_type = "API_TOKEN"
    related_id    = instana_api_token.dashboard_viewer.id
  }]
  
  widgets = file("${path.module}/dashboards/api_metrics.json")
}
```

### Dashboard with Mixed Access Levels

```hcl
resource "instana_custom_dashboard" "comprehensive_dashboard" {
  title = "Comprehensive Monitoring Dashboard"
  
  access_rule = [
    {
      access_type   = "READ_WRITE"
      relation_type = "USER"
      related_id    = "admin-user-id"
    },
    {
      access_type   = "READ_WRITE"
      relation_type = "TEAM"
      related_id    = "platform-team-id"
    },
    {
      access_type   = "READ"
      relation_type = "ROLE"
      related_id    = "viewer-role-id"
    },
    {
      access_type   = "READ"
      relation_type = "GLOBAL"
    }
  ]
  
  widgets = file("${path.module}/dashboards/comprehensive.json")
}
```

### Dashboard with Templated Widgets

```hcl
resource "instana_custom_dashboard" "dynamic_dashboard" {
  title = "Dynamic Environment Dashboard - ${var.environment}"
  
  access_rule = [{
    access_type   = "READ_WRITE"
    relation_type = "GLOBAL"
  }]
  
  widgets = templatefile("${path.module}/dashboards/template.json", {
    environment = var.environment
    region      = var.region
    app_name    = var.application_name
  })
}
```

### Infrastructure Monitoring Dashboard

```hcl
resource "instana_custom_dashboard" "infrastructure" {
  title = "Infrastructure Health Dashboard"
  
  access_rule = [{
    access_type   = "READ_WRITE"
    relation_type = "TEAM"
    related_id    = "infrastructure-team"
  }, {
    access_type   = "READ"
    relation_type = "GLOBAL"
  }]
  
  widgets = jsonencode([
    {
      type = "metric"
      title = "CPU Usage"
      metric = {
        plugin = "host"
        metric = "cpu.usage"
        aggregation = "mean"
      }
      timeframe = {
        to = 0
        windowSize = 3600000
      }
    },
    {
      type = "metric"
      title = "Memory Usage"
      metric = {
        plugin = "host"
        metric = "memory.used"
        aggregation = "mean"
      }
      timeframe = {
        to = 0
        windowSize = 3600000
      }
    }
  ])
}
```

### Application Performance Dashboard

```hcl
resource "instana_custom_dashboard" "application_performance" {
  title = "Application Performance Metrics"
  
  access_rule = [{
    access_type   = "READ_WRITE"
    relation_type = "TEAM"
    related_id    = "development-team"
  }]
  
  widgets = jsonencode([
    {
      type = "application"
      title = "Application Calls"
      applicationId = instana_application_config.my_app.id
      metric = "calls"
      timeframe = {
        to = 0
        windowSize = 7200000
      }
    },
    {
      type = "application"
      title = "Error Rate"
      applicationId = instana_application_config.my_app.id
      metric = "errors"
      timeframe = {
        to = 0
        windowSize = 7200000
      }
    },
    {
      type = "application"
      title = "Latency P95"
      applicationId = instana_application_config.my_app.id
      metric = "latency"
      aggregation = "p95"
      timeframe = {
        to = 0
        windowSize = 7200000
      }
    }
  ])
}
```

### Multi-Environment Dashboard

```hcl
locals {
  environments = ["dev", "staging", "production"]
}

resource "instana_custom_dashboard" "multi_env" {
  for_each = toset(local.environments)
  
  title = "${title(each.value)} Environment Dashboard"
  
  access_rule = [{
    access_type   = "READ_WRITE"
    relation_type = "TEAM"
    related_id    = "${each.value}-team"
  }, {
    access_type   = "READ"
    relation_type = "GLOBAL"
  }]
  
  widgets = templatefile("${path.module}/dashboards/environment.json", {
    environment = each.value
    zone        = "${each.value}-zone"
  })
}
```

## Argument Reference

* `title` - Required - The name/title of the custom dashboard
* `access_rule` - Required - Configuration of access rules (sharing/permissions) of the custom dashboard (list). [Details](#access-rule-argument-reference)
* `widgets` - Required - JSON array of widget configurations. It is recommended to get this configuration via the `Edit as Json` feature of custom dashboards in Instana UI and to adopt the configuration afterwards

### Access Rule Argument Reference

* `access_type` - Required - Type of granted access. Allowed values: `READ`, `READ_WRITE`
* `relation_type` - Required - Type of the entity for which the access is granted. Allowed values: `USER`, `API_TOKEN`, `ROLE`, `TEAM`, `GLOBAL`
* `related_id` - Optional - The ID of the related entity for which access is granted. Required for all `relation_type` except `GLOBAL`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the custom dashboard

## Widgets Configuration

The `widgets` field accepts a JSON array of widget configurations. Each widget can have different types and configurations.

### Common Widget Properties

All widgets support these common properties:

* `type` - The type of widget (e.g., `metric`, `application`, `website`, `log`, `event`)
* `title` - The title displayed on the widget
* `timeframe` - Time range configuration
  * `to` - End time (0 for current time)
  * `windowSize` - Time window in milliseconds

### Widget Types

#### Metric Widget

```json
{
  "type": "metric",
  "title": "CPU Usage",
  "metric": {
    "plugin": "host",
    "metric": "cpu.usage",
    "aggregation": "mean"
  },
  "timeframe": {
    "to": 0,
    "windowSize": 3600000
  }
}
```

#### Application Widget

```json
{
  "type": "application",
  "title": "Application Calls",
  "applicationId": "app-id",
  "metric": "calls",
  "timeframe": {
    "to": 0,
    "windowSize": 3600000
  }
}
```

#### Website Widget

```json
{
  "type": "website",
  "title": "Page Load Time",
  "websiteId": "website-id",
  "metric": "pageLoadTime",
  "timeframe": {
    "to": 0,
    "windowSize": 3600000
  }
}
```

### Using Terraform Functions

#### file() Function

Load widgets from a JSON file:

```hcl
widgets = file("${path.module}/dashboards/my_dashboard.json")
```

#### templatefile() Function

Use templates with variables:

```hcl
widgets = templatefile("${path.module}/dashboards/template.json", {
  environment = var.environment
  app_id      = instana_application_config.my_app.id
})
```

#### jsonencode() Function

Define widgets inline:

```hcl
widgets = jsonencode([
  {
    type = "metric"
    title = "CPU Usage"
    metric = {
      plugin = "host"
      metric = "cpu.usage"
    }
  }
])
```

## Import

Custom Dashboards can be imported using the `id`, e.g.:

```bash
terraform import instana_custom_dashboard.example 60845e4e5e6b9cf8fc2868da
```

## Best Practices

### Access Control

Use appropriate access levels:

```hcl
# Admin access for team
access_rule = [{
  access_type   = "READ_WRITE"
  relation_type = "TEAM"
  related_id    = "admin-team"
}, {
  # Read-only for everyone else
  access_type   = "READ"
  relation_type = "GLOBAL"
}]
```

### Widget Organization

- Store widget configurations in separate JSON files
- Use descriptive file names (e.g., `production_monitoring.json`)
- Version control your widget configurations
- Use templates for similar dashboards across environments

### JSON Management

```hcl
# Good - External file
widgets = file("${path.module}/dashboards/production.json")

# Good - Template with variables
widgets = templatefile("${path.module}/dashboards/template.json", {
  env = var.environment
})

# Avoid - Large inline JSON
widgets = jsonencode([...]) # Only for simple cases
```

### Naming Conventions

Use clear, descriptive titles:

```hcl
# Good
title = "Production API Performance Dashboard"

# Avoid
title = "Dashboard 1"
```

### Access Rule Patterns

```hcl
# Pattern 1: Team-based access
access_rule = [{
  access_type   = "READ_WRITE"
  relation_type = "TEAM"
  related_id    = "platform-team"
}]

# Pattern 2: Mixed access levels
access_rule = [
  {
    access_type   = "READ_WRITE"
    relation_type = "TEAM"
    related_id    = "admin-team"
  },
  {
    access_type   = "READ"
    relation_type = "GLOBAL"
  }
]

# Pattern 3: User-specific
access_rule = [{
  access_type   = "READ_WRITE"
  relation_type = "USER"
  related_id    = "user-id"
}]
```

## Notes

- The resource ID is auto-generated by Instana upon creation
- Widget configurations are stored as JSON and normalized automatically
- Access rules determine who can view and edit the dashboard
- `GLOBAL` relation type makes the dashboard accessible to all users
- Use the Instana UI's "Edit as Json" feature to get widget configurations
- Widget JSON is normalized to ensure consistent state management
- Multiple access rules can be defined to grant different levels of access to different entities
- Changes to widgets will update the dashboard immediately
