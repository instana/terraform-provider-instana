# Custom Dashboard

Management of custom dashboards.

API Documentation: <https://instana.github.io/openapi/#tag/Custom-Dashboards>

---
## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.

## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- `access_rule` now uses **list syntax** with `= [{ }]` instead of block syntax
- Enhanced validation for access types and relation types
- Improved JSON normalization for widgets
- Better state management

#### OLD (v5.x) Syntax:

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

#### NEW (v6.x) Syntax:

```hcl
resource "instana_custom_dashboard" "example" {
  title = "Example Dashboard"
  
  access_rule = [{
    access_type   = "READ_WRITE"
    relation_type = "USER"
    related_id    = "5ee8a3e8cd70020001ecb007" # replace with actual user id
  }, {
    access_type   = "READ"
    relation_type = "GLOBAL"
  }]
  
  widgets = file("${path.module}/widgets.json")
}
```

---

**Permissions Required:**
- `canCreatePublicCustomDashboards` (Creation of public custom dashboards)
- `canEditAllAccessibleCustomDashboards` (Management of all accessible custom dashboards)

The ID of the resource which is also used as unique identifier in Instana is auto generated!

## Example Usage

### Custom Dashboard

```hcl
resource "instana_custom_dashboard" "custom_dashboard" {
  access_rule = [
    {
      access_type   = "READ_WRITE"
      related_id    = "5ee8a3e8cd70020001ecb007" # replace with actual user id
      relation_type = "USER"
    },
  ]
  title = "custom_dashboard"
  widgets = jsonencode([{
    config = {
      comparisonDecreaseColor = "greenish"
      comparisonIncreaseColor = "redish"
      formatter               = "number.detailed"
      metricConfiguration = {
        aggregation      = "SUM"
        includeInternal  = false
        includeSynthetic = false
        metric           = "calls"
        source           = "APPLICATION"
        tagFilterExpression = {
          elements        = []
          logicalOperator = "AND"
          type            = "EXPRESSION"
        }
        threshold = {
          critical         = ""
          operator         = ">="
          thresholdEnabled = false
          warning          = ""
        }
        timeShift = 0
      }
    }
    height = 1
    id     = "id"
    title  = "Widget"
    type   = "bigNumber"
    width  = 1
    x      = 0
    y      = 0
  }])
}
```



## Generating Configuration from Existing Resources

If you have already created a custom dashboard in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your custom dashboard in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_custom_dashboard.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual custom dashboard ID from Instana

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
