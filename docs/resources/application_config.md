# Application Configuration Resource

Management of application configurations (definition of application perspectives).

API Documentation: <https://instana.github.io/openapi/#operation/putApplicationConfig>

---
## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- `match_specification` has been replaced with `access_rules` (list attribute with `= [{ }]`)
- Enhanced validation for scope and boundary_scope values
- Improved tag filter parsing and normalization
- Better state management with computed fields
- Default values are now explicit (scope defaults to `INCLUDE_NO_DOWNSTREAM`, boundary_scope to `DEFAULT`)

#### OLD (v5.x) Syntax:

```hcl
resource "instana_application_config" "example" {
  label          = "My Application"
  scope          = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "INBOUND"
  tag_filter     = "service.name@dest EQUALS 'DC11'"
  
  match_specification {
    # ... match specification blocks
  }
}
```

#### NEW (v6.x) Syntax:

```hcl
resource "instana_application_config" "example" {
  label          = "My Application"
  scope          = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "INBOUND"
  tag_filter     = "service.name@dest EQUALS 'DC11'"
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "GLOBAL"
    related_id    = "group-id"
  }]
}
```


### Basic Application Configuration

```hcl
resource "instana_application_config" "application_perspective_config" {
  access_rules = [ 
    {
      access_type   = "READ_WRITE"
      relation_type = "GLOBAL"
    },
  ]
  boundary_scope = "INBOUND"
  label          = "Label" # Replace with your own value
  scope          = "INCLUDE_NO_DOWNSTREAM" # Replace with your own value
  tag_filter     = "((call.type@na EQUALS 'HTTP' AND service.name@dest EQUALS 'cart' AND kubernetes.namespace@dest EQUALS 'robot-shop') OR service.name@dest EQUALS 'catalogue')" # Replace with your own value
}
```

### Application with Tag Filtering

```hcl
resource "instana_application_config" "tf_b_application_perspective_config_8" {
  access_rules = [
    {
      access_type   = "READ_WRITE"
      related_id    = null
      relation_type = "GLOBAL"
    },
  ]
  boundary_scope = "ALL"
  label          = "Label" # Replace with your own value
  scope          = "INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING" # Replace with your own value
  tag_filter     = "(kubernetes.deployment.namespace@dest EQUALS 'release-pink' AND container.image.name@dest EQUALS 'containers.instana.io/synthetic/synthetic-playback-browserscript:1.296.0')" # Replace with your own value
}
```

## Generating Configuration from Existing Resources

If you have already created a application configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your application configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_application_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual application configuration ID from Instana

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

* `label` - Required - The name/label of the application perspective
* `scope` - Optional - The scope of the application perspective. Default value: `INCLUDE_NO_DOWNSTREAM`. Allowed values:
  - `INCLUDE_ALL_DOWNSTREAM` - Include all downstream services
  - `INCLUDE_NO_DOWNSTREAM` - Include only the matched services
  - `INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING` - Include direct database and messaging dependencies
* `boundary_scope` - Optional - The boundary scope of the application perspective. Default value: `DEFAULT`. Allowed values:
  - `INBOUND` - Only inbound calls
  - `ALL` - All calls (inbound and outbound)
  - `DEFAULT` - Default boundary behavior
* `tag_filter` - Optional - Specifies which entities should be included in the application using a tag filter expression
* `access_rules` - Required - List of access rules defining who can access this application perspective. [Details](#access-rules-argument-reference)

### Access Rules Argument Reference

* `access_type` - Required - The type of access granted. Allowed values:
  - `READ_WRITE` - Full read and write access
  - `READ_ONLY` - Read-only access
* `relation_type` - Required - The type of relation. 
* `related_id` - Optional - The ID of the related entity (e.g., RBAC group ID). If not specified, the rule applies to all users

### Tag Filter

The **tag_filter** defines which entities should be included into the application. It supports:

* logical AND and/or logical OR conjunctions whereas AND has higher precedence then OR
* comparison operators EQUALS, NOT_EQUAL, CONTAINS | NOT_CONTAIN, STARTS_WITH, ENDS_WITH, NOT_STARTS_WITH, NOT_ENDS_WITH, GREATER_OR_EQUAL_THAN, LESS_OR_EQUAL_THAN, LESS_THAN, GREATER_THAN
* unary operators IS_EMPTY, NOT_EMPTY, IS_BLANK, NOT_BLANK

The **tag_filter** is defined by the following eBNF:

```plain
tag_filter                := logical_or
logical_or                := logical_and OR logical_or | logical_and
logical_and               := primary_expression AND logical_and | bracket_expression
bracket_expression        := ( logical_or ) | primary_expression
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

#### Tag Filter Examples

```hcl
tag_filter     = "(kubernetes.deployment.namespace@dest EQUALS 'release-pink' AND container.image.name@dest EQUALS 'containers.instana.io/synthetic/synthetic-playback-browserscript:1.296.0')"

```

```hcl
tag_filter     = "(service.name@dest EQUALS 'payment' AND service.name@dest EQUALS 'user' AND service.name@dest EQUALS 'cart')"

```


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the application configuration

## Import

Application Configs can be imported using the `id`, e.g.:

```bash
terraform import instana_application_config.my_app_config 60845e4e5e6b9cf8fc2868da
```

## Best Practices

### Scope Selection

Choose the appropriate scope for your use case:

```hcl
# For frontend applications - include all downstream
scope = "INCLUDE_ALL_DOWNSTREAM"

# For microservices - include only the service
scope = "INCLUDE_NO_DOWNSTREAM"

# For services with databases - include immediate downstream
scope = "INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING"
```

### Boundary Scope

Select boundary scope based on monitoring needs:

```hcl
# For API services - monitor inbound calls only
boundary_scope = "INBOUND"

# For comprehensive monitoring - monitor all calls
boundary_scope = "ALL"

# For default behavior
boundary_scope = "DEFAULT"
```

### Tag Filter Best Practices

1. **Use Specific Filters**: Be as specific as possible to avoid including unwanted services
2. **Test Filters**: Verify your tag filter matches the expected services in Instana UI
3. **Use join() for Readability**: Break complex filters into multiple lines
4. **Document Complex Logic**: Add comments explaining complex filter logic

```hcl
tag_filter = join(" AND ", [
  # Primary service identification
  "service.name EQUALS 'api-service'",
  
  # Environment filtering
  "entity.zone EQUALS 'production'",
  
  # Exclude test traffic
  "NOT call.tag:test@na EQUALS 'true'"
])
```

### Access Rules

Configure appropriate access levels:

```hcl
access_rules = [
  {
    access_type   = "READ_WRITE"
    relation_type = "GLOBAL"
    related_id    = instana_rbac_group.admins.id
  },
  {
    access_type   = "READ_ONLY"
    relation_type = "GLOBAL"
    related_id    = instana_rbac_group.viewers.id
  }
]
```

### Naming Conventions

Use clear, descriptive labels:

```hcl
# Good
label = "Production E-commerce API"

# Avoid
label = "app1"
```

## Notes

- The resource ID is auto-generated by Instana upon creation
- Tag filters are normalized and stored in a canonical format
- Changes to `tag_filter` will update the application perspective immediately
- The `scope` setting affects which downstream services are included in metrics
- Use `boundary_scope` to control whether internal service-to-service calls are included
- Access rules integrate with Instana's RBAC system for fine-grained access control
- Tag filters support complex boolean logic with proper operator precedence (AND before OR)
