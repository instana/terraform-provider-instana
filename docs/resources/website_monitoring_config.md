# Website Monitoring Config Resource

> **⚠️ BREAKING CHANGES - Plugin Framework Migration**
>
> This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework. The schema has transitioned from **block structure to attribute format**.
>
> **Major Changes:**
> - All attributes remain top-level (no nested blocks in this resource)
> - The `id` attribute is now computed with a plan modifier for state management
> - The `app_name` attribute is computed and returned by the API
> - Attribute syntax remains the same, but schema validation is enhanced
>
> **Migration Example:**
> ```hcl
> # OLD (SDK v2)
> resource "instana_website_monitoring_config" "example" {
>   name = "my-website"
> }
>
> # NEW (Plugin Framework - Same Syntax)
> resource "instana_website_monitoring_config" "example" {
>   name = "my-website"
> }
> ```
>
> The syntax is fully compatible, but the framework provides better validation and state management.

Resource to configure websites in Instana for Real User Monitoring (RUM).

API Documentation: <https://instana.github.io/openapi/#tag/Website-Configuration

 **⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)**

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

### Syntax Changes Overview

- All attributes remain top-level (no nested blocks in this resource)
- The `id` attribute is now computed with a plan modifier for state management
- The `app_name` attribute is computed and returned by the API
- Attribute syntax remains the same, but schema validation is enhanced

#### OLD (v5.x) Syntax:

### Basic Website Monitoring Configuration

#### Simple Website Configuration
```hcl
resource "instana_website_monitoring_config" "basic" {
  name = "my-website"
}
```

#### Production Website Configuration
```hcl
resource "instana_website_monitoring_config" "production" {
  name = "production-website"
}
```

#### Staging Website Configuration
```hcl
resource "instana_website_monitoring_config" "staging" {
  name = "staging-website"
}
```

### Multiple Website Configurations

#### Multi-Environment Setup
```hcl
# Production environment
resource "instana_website_monitoring_config" "prod" {
  name = "production-ecommerce-site"
}

# Staging environment
resource "instana_website_monitoring_config" "staging" {
  name = "staging-ecommerce-site"
}

# Development environment
resource "instana_website_monitoring_config" "dev" {
  name = "development-ecommerce-site"
}
```

#### Multi-Region Website Monitoring
```hcl
# US region
resource "instana_website_monitoring_config" "us" {
  name = "website-us-region"
}

# EU region
resource "instana_website_monitoring_config" "eu" {
  name = "website-eu-region"
}

# APAC region
resource "instana_website_monitoring_config" "apac" {
  name = "website-apac-region"
}
```

### Application-Specific Configurations

#### E-commerce Website
```hcl
resource "instana_website_monitoring_config" "ecommerce" {
  name = "ecommerce-platform"
}
```

#### Corporate Website
```hcl
resource "instana_website_monitoring_config" "corporate" {
  name = "corporate-website"
}
```

#### Customer Portal
```hcl
resource "instana_website_monitoring_config" "customer_portal" {
  name = "customer-portal"
}
```

#### Admin Dashboard
```hcl
resource "instana_website_monitoring_config" "admin_dashboard" {
  name = "admin-dashboard"
}
```

#### Mobile Web Application
```hcl
resource "instana_website_monitoring_config" "mobile_web" {
  name = "mobile-web-app"
}
```

### Using with Variables

#### Parameterized Configuration
```hcl
variable "website_name" {
  description = "Name of the website to monitor"
  type        = string
}

resource "instana_website_monitoring_config" "parameterized" {
  name = var.website_name
}
```

#### Environment-Based Configuration
```hcl
variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "application_name" {
  description = "Application name"
  type        = string
}

resource "instana_website_monitoring_config" "env_based" {
  name = "${var.application_name}-${var.environment}"
}
```

### Using with Outputs

#### Exporting Website Configuration Details
```hcl
resource "instana_website_monitoring_config" "main" {
  name = "main-website"
}

output "website_id" {
  description = "The ID of the website monitoring configuration"
  value       = instana_website_monitoring_config.main.id
}

output "website_name" {
  description = "The name of the website monitoring configuration"
  value       = instana_website_monitoring_config.main.name
}

output "website_app_name" {
  description = "The application name assigned by Instana"
  value       = instana_website_monitoring_config.main.app_name
}
```

### Integration with Other Resources

#### Website with SLO Configuration
```hcl
# Create website monitoring configuration
resource "instana_website_monitoring_config" "website" {
  name = "production-website"
}

# Create SLO for the website
resource "instana_slo_config" "website_slo" {
  name   = "website-availability-slo"
  target = 0.99
  
  entity = {
    website = {
      website_id  = instana_website_monitoring_config.website.id
      beacon_type = "httpRequest"
    }
  }
  
  indicator = {
    time_based_availability = {
      threshold   = 0.0
      aggregation = "MEAN"
    }
  }
  
  time_window = {
    rolling = {
      duration      = 7
      duration_unit = "day"
    }
  }
}
```

#### Website with Alert Configuration
```hcl
# Create website monitoring configuration
resource "instana_website_monitoring_config" "website" {
  name = "critical-website"
}

# Create website alert configuration
resource "instana_website_alert_config" "website_alert" {
  name                 = "website-performance-alert"
  website_id           = instana_website_monitoring_config.website.id
  severity             = "warning"
  triggering_enabled   = true
  
  # Alert configuration details...
}
```

### Naming Conventions

#### Descriptive Naming
```hcl
# Good: Clear, descriptive names
resource "instana_website_monitoring_config" "customer_facing_portal" {
  name = "customer-facing-portal"
}

resource "instana_website_monitoring_config" "internal_admin_panel" {
  name = "internal-admin-panel"
}

resource "instana_website_monitoring_config" "public_marketing_site" {
  name = "public-marketing-site"
}
```

#### Hierarchical Naming
```hcl
# Organization-based naming
resource "instana_website_monitoring_config" "acme_prod_main" {
  name = "acme-prod-main-website"
}

resource "instana_website_monitoring_config" "acme_prod_api_docs" {
  name = "acme-prod-api-docs"
}

resource "instana_website_monitoring_config" "acme_staging_main" {
  name = "acme-staging-main-website"
}
```

### Dynamic Configuration with For-Each

#### Multiple Websites from List
```hcl
variable "websites" {
  description = "List of websites to monitor"
  type        = list(string)
  default     = ["main-site", "blog", "docs", "api-portal"]
}

resource "instana_website_monitoring_config" "websites" {
  for_each = toset(var.websites)
  
  name = each.value
}

output "website_ids" {
  description = "Map of website names to their IDs"
  value = {
    for name, config in instana_website_monitoring_config.websites :
    name => config.id
  }
}
```

#### Multiple Environments with For-Each
```hcl
variable "environments" {
  description = "Map of environments to website configurations"
  type = map(object({
    website_name = string
  }))
  default = {
    production = {
      website_name = "prod-website"
    }
    staging = {
      website_name = "staging-website"
    }
    development = {
      website_name = "dev-website"
    }
  }
}

resource "instana_website_monitoring_config" "environments" {
  for_each = var.environments
  
  name = each.value.website_name
}
```

### Complete Production Example

#### Full Production Setup
```hcl
# Variables
variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "application_name" {
  description = "Application name"
  type        = string
  default     = "ecommerce-platform"
}

# Website monitoring configuration
resource "instana_website_monitoring_config" "production" {
  name = "${var.application_name}-${var.environment}"
}

# SLO for website availability
resource "instana_slo_config" "website_availability" {
  name   = "${var.application_name}-availability-slo"
  target = 0.995
  tags   = [var.environment, "website", "availability"]
  
  entity = {
    website = {
      website_id  = instana_website_monitoring_config.production.id
      beacon_type = "httpRequest"
    }
  }
  
  indicator = {
    time_based_availability = {
      threshold   = 0.0
      aggregation = "MEAN"
    }
  }
  
  time_window = {
    rolling = {
      duration      = 30
      duration_unit = "day"
      timezone      = "UTC"
    }
  }
}

# SLO for website latency
resource "instana_slo_config" "website_latency" {
  name   = "${var.application_name}-latency-slo"
  target = 0.95
  tags   = [var.environment, "website", "latency"]
  
  entity = {
    website = {
      website_id  = instana_website_monitoring_config.production.id
      beacon_type = "httpRequest"
    }
  }
  
  indicator = {
    time_based_latency = {
      threshold   = 2000.0
      aggregation = "P95"
    }
  }
  
  time_window = {
    rolling = {
      duration      = 7
      duration_unit = "day"
      timezone      = "UTC"
    }
  }
}

# Outputs
output "website_id" {
  description = "Website monitoring configuration ID"
  value       = instana_website_monitoring_config.production.id
}

output "website_app_name" {
  description = "Website application name in Instana"
  value       = instana_website_monitoring_config.production.app_name
}

output "availability_slo_id" {
  description = "Availability SLO ID"
  value       = instana_slo_config.website_availability.id
}

output "latency_slo_id" {
  description = "Latency SLO ID"
  value       = instana_slo_config.website_latency.id
}
```

## Argument Reference

### Required Attributes

* `name` - (Required) The name of the website monitoring configuration. This name will be used to identify the website in Instana.

**Type:** `string`

### Computed Attributes

* `id` - (Computed) The unique identifier of the website monitoring configuration
* `app_name` - (Computed) The application name assigned by Instana for this website

**Type:** `string` for both attributes

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the website monitoring configuration
* `app_name` - The application name assigned by Instana

## Import

Website Monitoring Configs can be imported using the `id`, e.g.:

```bash
$ terraform import instana_website_monitoring_config.my_website 60845e4e5e6b9cf8fc2868da
```

## Notes

### Website Monitoring Setup

After creating a website monitoring configuration:

1. **Obtain the JavaScript Snippet**: Use the Instana UI to get the JavaScript snippet for your website
2. **Install the Snippet**: Add the snippet to your website's HTML (typically in the `<head>` section)
3. **Verify Data Collection**: Check the Instana UI to confirm that data is being collected

### Application Name

The `app_name` attribute is automatically generated by Instana and may differ from the `name` you provide. Use the `app_name` when referencing the website in other Instana configurations or APIs.

### Website ID Usage

The `id` of the website monitoring configuration is used when:
- Creating SLO configurations for the website
- Creating alert configurations for the website
- Referencing the website in other Instana resources

### Naming Best Practices

1. **Be Descriptive**: Use clear, descriptive names that identify the website's purpose
2. **Include Environment**: Consider including the environment in the name (e.g., "prod-website", "staging-website")
3. **Use Consistent Format**: Maintain a consistent naming convention across all websites
4. **Avoid Special Characters**: Stick to alphanumeric characters and hyphens
5. **Keep It Concise**: While being descriptive, avoid overly long names

### Common Use Cases

**Single Page Applications (SPAs):**
```hcl
resource "instana_website_monitoring_config" "spa" {
  name = "react-spa-application"
}
```

**Multi-Page Applications:**
```hcl
resource "instana_website_monitoring_config" "mpa" {
  name = "traditional-web-application"
}
```

**Progressive Web Apps (PWAs):**
```hcl
resource "instana_website_monitoring_config" "pwa" {
  name = "progressive-web-app"
}
```

**Mobile Web Applications:**
```hcl
resource "instana_website_monitoring_config" "mobile" {
  name = "mobile-web-application"
}
```

### Integration Patterns

**With Terraform Modules:**
```hcl
module "website_monitoring" {
  source = "./modules/website-monitoring"
  
  website_name = "my-application"
  environment  = "production"
}
```

**With Data Sources:**
```hcl
# Reference existing website configuration
data "instana_website_monitoring_config" "existing" {
  name = "existing-website"
}

# Use in SLO configuration
resource "instana_slo_config" "website_slo" {
  name   = "website-slo"
  target = 0.99
  
  entity = {
    website = {
      website_id  = data.instana_website_monitoring_config.existing.id
      beacon_type = "httpRequest"
    }
  }
  
  # ... rest of configuration
}
```

### Lifecycle Management

**Prevent Accidental Deletion:**
```hcl
resource "instana_website_monitoring_config" "critical" {
  name = "critical-production-website"
  
  lifecycle {
    prevent_destroy = true
  }
}
```

**Ignore External Changes:**
```hcl
resource "instana_website_monitoring_config" "managed_externally" {
  name = "partially-managed-website"
  
  lifecycle {
    ignore_changes = [name]
  }
}
```

### Monitoring Multiple Domains

When monitoring multiple domains or subdomains of the same application:

```hcl
# Main domain
resource "instana_website_monitoring_config" "main_domain" {
  name = "example-com"
}

# Subdomain - API
resource "instana_website_monitoring_config" "api_subdomain" {
  name = "api-example-com"
}

# Subdomain - Blog
resource "instana_website_monitoring_config" "blog_subdomain" {
  name = "blog-example-com"
}
```

### Best Practices Summary

1. **Unique Names**: Ensure each website has a unique, identifiable name
2. **Environment Separation**: Create separate configurations for different environments
3. **Use Outputs**: Export IDs and app names for use in other resources
4. **Document Purpose**: Use comments to document the purpose of each website configuration
5. **Consistent Naming**: Follow a consistent naming convention across your organization
6. **Version Control**: Track website configurations in version control
7. **Regular Review**: Periodically review and clean up unused website configurations
