# Mobile App Config Resource

Resource to configure mobile applications in Instana for Mobile App Monitoring.

API Documentation: [Instana REST API - Mobile App Config](https://developer.ibm.com/apis/catalog/instana--instana-rest-api/api/API--instana--instana-rest-api-documentation#postMobileAppConfig)

## Example Usage

### Mobile App Configuration

```hcl
resource "instana_mobile_app_config" "example" {
  name = "my-mobile-app"
}
```

### Using with Variables

#### Parameterized Configuration
```hcl
variable "app_name" {
  description = "Name of the mobile app to monitor"
  type        = string
}

resource "instana_mobile_app_config" "parameterized" {
  name = var.app_name
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

resource "instana_mobile_app_config" "env_based" {
  name = "${var.application_name}-${var.environment}"
}
```

### Integration with Other Resources

#### Mobile App with Alert Configuration
```hcl
# Create mobile app monitoring configuration
resource "instana_mobile_app_config" "mobile_app" {
  name = "production-mobile-app"
}

# Create mobile app alert configuration
resource "instana_mobile_alert_config" "mobile_alert" {
  name          = "mobile-app-performance-alert"
  mobile_app_id = instana_mobile_app_config.mobile_app.id
  
  # Alert configuration details...
}
```


## Generating Configuration from Existing Resources

If you have already created a mobile app configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your mobile app configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_mobile_app_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual mobile app configuration ID from Instana

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

* `name` - (Required) The name of the mobile app configuration. This name will be used to identify the mobile app in Instana.

**Type:** `string`

**Constraints:** Must be between 1 and 128 characters

### Computed Attributes

* `id` - (Computed) The unique identifier of the mobile app configuration

**Type:** `string`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the mobile app configuration

## Import

Mobile App Configs can be imported using the `id`, e.g.:

```bash
$ terraform import instana_mobile_app_config.my_app K3bP-bmCRkyimNai9vvq8o
```

## Notes

### Mobile App Monitoring Setup

After creating a mobile app configuration:

1. **Obtain the App Key**: Use the Instana UI to get the app key for your mobile application
2. **Install the SDK**: Add the Instana mobile SDK to your iOS or Android application
3. **Configure the SDK**: Initialize the SDK with your app key
4. **Verify Data Collection**: Check the Instana UI to confirm that data is being collected

### Mobile App ID Usage

The `id` of the mobile app configuration is used when:
- Creating alert configurations for the mobile app
- Referencing the mobile app in other Instana resources

### Naming Best Practices

1. **Be Descriptive**: Use clear, descriptive names that identify the app's purpose
2. **Include Platform**: Consider including the platform in the name (e.g., "myapp-ios", "myapp-android")
3. **Include Environment**: Consider including the environment in the name (e.g., "prod-myapp", "staging-myapp")
4. **Use Consistent Format**: Maintain a consistent naming convention across all mobile apps
5. **Avoid Special Characters**: Stick to alphanumeric characters, hyphens, and underscores
6. **Keep It Concise**: While being descriptive, avoid overly long names (max 128 characters)

### Integration Patterns

**With Terraform Modules:**
```hcl
module "mobile_app_monitoring" {
  source = "./modules/mobile-app-monitoring"
  
  app_name    = "my-application"
  platform    = "ios"
  environment = "production"
}
```

### Lifecycle Management

**Prevent Accidental Deletion:**
```hcl
resource "instana_mobile_app_config" "critical" {
  name = "critical-production-app"
  
  lifecycle {
    prevent_destroy = true
  }
}
```

**Ignore External Changes:**
```hcl
resource "instana_mobile_app_config" "managed_externally" {
  name = "partially-managed-app"
  
  lifecycle {
    ignore_changes = [name]
  }
}
```

### Best Practices Summary

1. **Unique Names**: Ensure each mobile app has a unique, identifiable name
2. **Platform Separation**: Create separate configurations for iOS and Android versions
3. **Environment Separation**: Create separate configurations for different environments
4. **Use Outputs**: Export IDs for use in other resources
5. **Document Purpose**: Use comments to document the purpose of each mobile app configuration
6. **Consistent Naming**: Follow a consistent naming convention across your organization
7. **Version Control**: Track mobile app configurations in version control
8. **Regular Review**: Periodically review and clean up unused mobile app configurations
