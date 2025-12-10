# Website Monitoring Config Resource

Resource to configure websites in Instana for Real User Monitoring (RUM).

API Documentation: <https://instana.github.io/openapi/#tag/Website-Configuration

 **⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)**

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- All attributes remain top-level (no nested blocks in this resource)
- The `id` attribute is now computed with a plan modifier for state management
- The `app_name` attribute is computed and returned by the API
- Attribute syntax remains the same, but schema validation is enhanced

#### OLD (v5.x) Syntax:

```hcl
resource "instana_website_monitoring_config" "example" {
  name = "my-website"
}
```

#### NEW (v6.x) Syntax:

```hcl
resource "instana_website_monitoring_config" "example" {
  name = "my-website"
}
```

 **⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)**

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

### Syntax Changes Overview

- All attributes remain top-level (no nested blocks in this resource)
- The `id` attribute is now computed with a plan modifier for state management
- The `app_name` attribute is computed and returned by the API
- Attribute syntax remains the same, but schema validation is enhanced


### Basic Website Monitoring Configuration

####  Website Configuration
```hcl
resource "instana_website_monitoring_config" "basic" {
  name = "my-website"
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
  entity = {
    website = {
      website_id  = instana_website_monitoring_config.website.id
      beacon_type = "httpRequest"
    }
  }
  # slo configuration details ....
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
  
  # Alert configuration details...
}
```

## Generating Configuration from Existing Resources

If you have already created a website monitoring configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your website monitoring configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_website_monitoring_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual website monitoring configuration ID from Instana

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
2. **Install the Snippet**: Add the snippet to your website's HTML (typically in the `<head` section)
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

### Integration Patterns

**With Terraform Modules:**
```hcl
module "website_monitoring" {
  source = "./modules/website-monitoring"
  
  website_name = "my-application"
  environment  = "production"
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

### Best Practices Summary

1. **Unique Names**: Ensure each website has a unique, identifiable name
2. **Environment Separation**: Create separate configurations for different environments
3. **Use Outputs**: Export IDs and app names for use in other resources
4. **Document Purpose**: Use comments to document the purpose of each website configuration
5. **Consistent Naming**: Follow a consistent naming convention across your organization
6. **Version Control**: Track website configurations in version control
7. **Regular Review**: Periodically review and clean up unused website configurations
