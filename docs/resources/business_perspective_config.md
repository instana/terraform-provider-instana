# Business Perspective Config Resource

Resource to configure business perspectives in Instana for Business Monitoring.

API Documentation: [Instana REST API - Business Perspectives](https://developer.ibm.com/apis/catalog/instana--instana-rest-api/api/API--instana--instana-rest-api-documentation#createBusinessPerspective)

## Example Usage

### Basic Business Perspective Configuration

```hcl
resource "instana_business_perspective_config" "example" {
  name        = "Order Processing"
  description = "Monitors checkout and order processing workflows"
  tag_filter  = "service.name EQUALS 'order-service' OR service.name EQUALS 'payment-service'"
}
```

### Minimal Business Perspective Configuration

```hcl
resource "instana_business_perspective_config" "minimal" {
  name = "Payment Services"
}
```

### Business Perspective with Tag Filter Expression

```hcl
resource "instana_business_perspective_config" "complex" {
  name        = "Critical Business Flow"
  description = "Monitors frontend and checkout backend services in production"
  tag_filter  = "entity.type EQUALS 'service' AND (service.name EQUALS 'cart-service' OR service.name EQUALS 'checkout-service')"
}
```

### Business Perspective with RBAC Tags

```hcl
resource "instana_business_perspective_config" "with_rbac_tags" {
  name        = "Payment Flow Perspective"
  description = "Business perspective for payment operations"
  tag_filter  = "service.name EQUALS 'payment-service'"

  rbac_tags = [
    {
      id           = "0xHKaxfaS161Al6Qc23g4w"
      display_name = "team 1"
    }
  ]
}
```

### Integration with RBAC Team

```hcl
resource "instana_business_perspective_config" "order_flow" {
  name        = "Order Flow Perspective"
  description = "Business perspective for order operations"
  tag_filter  = "service.name EQUALS 'order-service'"
}

resource "instana_rbac_team" "devops" {
  name                 = "DevOps Team"
  business_perspectives = [instana_business_perspective_config.order_flow.id]
}
```

## Generating Configuration from Existing Resources

If you have already created a business perspective configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your business perspective in Instana (from the Instana UI URL or via the API).

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_business_perspective_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform resource identifier
- `resource_id` with your actual business perspective ID from Instana

3. **Generate the Configuration**: Run the following Terraform command:

```bash
terraform plan -generate-config-out=generated.tf
```

4. **Review and Apply**: Review the generated configuration in `generated.tf` and apply:

```bash
terraform apply
```

## Argument Reference

### Required Attributes

* `name` - (Required) The name of the business perspective configuration.

**Type:** `string`

**Constraints:** Max length: 100 characters

### Optional Attributes

* `description` - (Optional, Computed) The description of the business perspective configuration.

**Type:** `string`

**Constraints:** Max length: 300 characters

* `tag_filter` - (Optional) The tag filter expression used to filter services, endpoints, or calls included in this business perspective.

**Type:** `string`

* `rbac_tags` - (Optional, Computed) List of RBAC tags (teams) assigned to this business perspective for access control. Each tag object contains:
  * `id` - (Required) The unique identifier of the RBAC team.
  * `display_name` - (Required) The display name of the RBAC team.

**Type:** `list(object)`

### Computed Attributes

* `id` - (Computed) The unique identifier of the business perspective configuration.

**Type:** `string`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the business perspective configuration.

## Import

Business Perspective Configs can be imported using the `id`, e.g.:

```bash
$ terraform import instana_business_perspective_config.example gBwNFIbjS6Ozwt0a12regg
```
