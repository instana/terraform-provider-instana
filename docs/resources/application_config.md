# Application Configuration Resource

Management of application configurations (definition of application perspectives).

API Documentation: <https://instana.github.io/openapi/#operation/putApplicationConfig>

---
## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block struture.


## Migration Guide (v5 to v6)

### Syntax Changes Overview

**OLD Syntax (SDK v2):**
```hcl
resource "instana_application_config" "example" {
  label          = "My Application"
  scope          = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "INBOUND"
  tag_filter     = "agent.tag:stage EQUALS 'test'"
  
  match_specification {
    # ... match specification blocks
  }
}
```

**NEW Syntax (Plugin Framework):**
```hcl
resource "instana_application_config" "example" {
  label          = "My Application"
  scope          = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "INBOUND"
  tag_filter     = "agent.tag:stage EQUALS 'test'"
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
    related_id    = "group-id"
  }]
}
```

**Key Changes:**
- `match_specification` has been replaced with `access_rules` (list attribute with `= [{ }]`)
- Enhanced validation for scope and boundary_scope values
- Improved tag filter parsing and normalization
- Better state management with computed fields
- Default values are now explicit (scope defaults to `INCLUDE_NO_DOWNSTREAM`, boundary_scope to `DEFAULT`)

---


## Example Usage

### Basic Application Configuration

```hcl
resource "instana_application_config" "example" {
  label          = "Production API"
  scope          = "INCLUDE_NO_DOWNSTREAM"
  boundary_scope = "DEFAULT"
  tag_filter     = "service.name EQUALS 'api-service' AND entity.zone EQUALS 'production'"
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with Downstream Services

```hcl
resource "instana_application_config" "with_downstream" {
  label          = "E-commerce Platform"
  scope          = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "INBOUND"
  tag_filter     = "service.name EQUALS 'frontend' OR service.name EQUALS 'api-gateway'"
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with Tag-Based Filtering

```hcl
resource "instana_application_config" "tag_based" {
  label          = "Staging Environment"
  scope          = "INCLUDE_NO_DOWNSTREAM"
  boundary_scope = "ALL"
  tag_filter     = join(" AND ", [
    "agent.tag:stage EQUALS 'staging'",
    "aws.ec2.tag:environment EQUALS 'staging'",
    "call.tag:version@na STARTS_WITH 'v2'"
  ])
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Multi-Service Application

```hcl
resource "instana_application_config" "microservices" {
  label          = "Payment Platform"
  scope          = "INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING"
  boundary_scope = "INBOUND"
  tag_filter     = join(" OR ", [
    "service.name EQUALS 'payment-service'",
    "service.name EQUALS 'billing-service'",
    "service.name EQUALS 'invoice-service'"
  ])
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with AWS Tags

```hcl
resource "instana_application_config" "aws_app" {
  label          = "AWS Production Services"
  scope          = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "INBOUND"
  tag_filter     = join(" AND ", [
    "aws.ec2.tag:Environment EQUALS 'production'",
    "aws.ec2.tag:Team EQUALS 'platform'",
    "entity.type EQUALS 'service'"
  ])
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with Kubernetes Labels

```hcl
resource "instana_application_config" "k8s_app" {
  label          = "Kubernetes Production App"
  scope          = "INCLUDE_NO_DOWNSTREAM"
  boundary_scope = "ALL"
  tag_filter     = join(" AND ", [
    "kubernetes.pod.label:app EQUALS 'my-app'",
    "kubernetes.namespace EQUALS 'production'",
    "kubernetes.cluster.name EQUALS 'prod-cluster'"
  ])
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with Call Tag Filtering

```hcl
resource "instana_application_config" "call_filtered" {
  label          = "API v2 Endpoints"
  scope          = "INCLUDE_NO_DOWNSTREAM"
  boundary_scope = "INBOUND"
  tag_filter     = join(" AND ", [
    "call.tag:api-version@na EQUALS 'v2'",
    "call.http.status@na GREATER_THAN 0",
    "service.name STARTS_WITH 'api-'"
  ])
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with Source-Based Filtering

```hcl
resource "instana_application_config" "source_filtered" {
  label          = "External API Calls"
  scope          = "INCLUDE_NO_DOWNSTREAM"
  boundary_scope = "INBOUND"
  tag_filter     = join(" AND ", [
    "service.name@src EQUALS 'external-gateway'",
    "agent.tag:zone@src EQUALS 'dmz'"
  ])
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with Complex Tag Logic

```hcl
resource "instana_application_config" "complex" {
  label          = "Production Services (Excluding Test)"
  scope          = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "ALL"
  tag_filter     = join(" AND ", [
    "entity.zone EQUALS 'production'",
    "(service.name STARTS_WITH 'prod-' OR service.name STARTS_WITH 'api-')",
    "NOT service.name CONTAINS 'test'",
    "entity.type EQUALS 'service'"
  ])
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with Database Scope

```hcl
resource "instana_application_config" "with_databases" {
  label          = "Application with Data Layer"
  scope          = "INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING"
  boundary_scope = "INBOUND"
  tag_filter     = "service.name EQUALS 'api-service'"
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with RBAC Access Rules

```hcl
resource "instana_application_config" "with_rbac" {
  label          = "Restricted Application"
  scope          = "INCLUDE_NO_DOWNSTREAM"
  boundary_scope = "DEFAULT"
  tag_filter     = "service.name EQUALS 'sensitive-service'"
  
  access_rules = [
    {
      access_type   = "READ_WRITE"
      relation_type = "SOURCE"
      related_id    = instana_rbac_group.admins.id
    },
    {
      access_type   = "READ_ONLY"
      relation_type = "SOURCE"
      related_id    = instana_rbac_group.viewers.id
    }
  ]
}
```

### Multi-Region Application

```hcl
resource "instana_application_config" "multi_region" {
  label          = "Global Application"
  scope          = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "ALL"
  tag_filter     = join(" OR ", [
    "(entity.zone EQUALS 'us-east-1' AND service.name EQUALS 'api-service')",
    "(entity.zone EQUALS 'eu-west-1' AND service.name EQUALS 'api-service')",
    "(entity.zone EQUALS 'ap-southeast-1' AND service.name EQUALS 'api-service')"
  ])
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with HTTP Status Filtering

```hcl
resource "instana_application_config" "error_tracking" {
  label          = "Error Tracking Application"
  scope          = "INCLUDE_NO_DOWNSTREAM"
  boundary_scope = "INBOUND"
  tag_filter     = join(" AND ", [
    "service.name EQUALS 'api-service'",
    "call.http.status@na GREATER_OR_EQUAL_THAN 400"
  ])
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Application with Unary Operators

```hcl
resource "instana_application_config" "with_unary" {
  label          = "Services with Tags"
  scope          = "INCLUDE_NO_DOWNSTREAM"
  boundary_scope = "DEFAULT"
  tag_filter     = join(" AND ", [
    "service.name STARTS_WITH 'prod-'",
    "agent.tag:team NOT_EMPTY",
    "agent.tag:version NOT_BLANK"
  ])
  
  access_rules = [{
    access_type   = "READ_WRITE"
    relation_type = "SOURCE"
  }]
}
```

### Comprehensive Application Configuration

```hcl
resource "instana_application_config" "comprehensive" {
  label          = "Production E-commerce Platform"
  scope          = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "INBOUND"
  
  tag_filter = join(" AND ", [
    # Service identification
    "(service.name EQUALS 'frontend' OR service.name EQUALS 'api-gateway' OR service.name EQUALS 'checkout-service')",
    
    # Environment filtering
    "entity.zone EQUALS 'production'",
    "agent.tag:environment EQUALS 'prod'",
    
    # Exclude test traffic
    "NOT call.tag:test-traffic@na EQUALS 'true'",
    
    # Include only successful or client error calls
    "(call.http.status@na LESS_THAN 500 OR call.http.status@na IS_EMPTY)",
    
    # Kubernetes filtering
    "kubernetes.namespace EQUALS 'production'",
    
    # AWS filtering
    "aws.ec2.tag:CostCenter EQUALS 'ecommerce'"
  ])
  
  access_rules = [
    {
      access_type   = "READ_WRITE"
      relation_type = "SOURCE"
      related_id    = instana_rbac_group.platform_team.id
    },
    {
      access_type   = "READ_ONLY"
      relation_type = "SOURCE"
      related_id    = instana_rbac_group.developers.id
    }
  ]
}
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
* `relation_type` - Required - The type of relation. Allowed values:
  - `SOURCE` - Source-based relation
  - `DESTINATION` - Destination-based relation
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

**Basic Service Filtering:**
```hcl
tag_filter = "service.name EQUALS 'my-service' AND agent.tag:stage EQUALS 'PROD' AND call.http.status@na EQUALS 404"
```

**Calls Filtered on Source:**
```hcl
tag_filter = "service.name@src EQUALS 'my-service' AND agent.tag:stage@src EQUALS 'PROD'"
```

**Multiple Services with OR:**
```hcl
tag_filter = "service.name EQUALS 'service-a' OR service.name EQUALS 'service-b' OR service.name EQUALS 'service-c'"
```

**Complex Logic with Parentheses:**
```hcl
tag_filter = "(service.name STARTS_WITH 'api-' OR service.name STARTS_WITH 'web-') AND entity.zone EQUALS 'production'"
```

**Using Unary Operators:**
```hcl
tag_filter = "service.name NOT_EMPTY AND agent.tag:version NOT_BLANK"
```

**Numeric Comparisons:**
```hcl
tag_filter = "call.http.status@na GREATER_OR_EQUAL_THAN 400 AND call.http.status@na LESS_THAN 500"
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
    relation_type = "SOURCE"
    related_id    = instana_rbac_group.admins.id
  },
  {
    access_type   = "READ_ONLY"
    relation_type = "SOURCE"
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

## Common Tag Filter Patterns

### Environment-Based

```hcl
tag_filter = "agent.tag:environment EQUALS 'production'"
```

### Service Name Patterns

```hcl
tag_filter = "service.name STARTS_WITH 'prod-' AND service.name NOT_CONTAINS 'test'"
```

### Cloud Provider Tags

```hcl
# AWS
tag_filter = "aws.ec2.tag:Environment EQUALS 'production' AND aws.ec2.tag:Team EQUALS 'platform'"

# Kubernetes
tag_filter = "kubernetes.pod.label:app EQUALS 'my-app' AND kubernetes.namespace EQUALS 'production'"
```

### Call-Based Filtering

```hcl
tag_filter = "call.type@na EQUALS 'HTTP' AND call.http.status@na GREATER_OR_EQUAL_THAN 200 AND call.http.status@na LESS_THAN 300"
```

## Notes

- The resource ID is auto-generated by Instana upon creation
- Tag filters are normalized and stored in a canonical format
- Changes to `tag_filter` will update the application perspective immediately
- The `scope` setting affects which downstream services are included in metrics
- Use `boundary_scope` to control whether internal service-to-service calls are included
- Access rules integrate with Instana's RBAC system for fine-grained access control
- Tag filters support complex boolean logic with proper operator precedence (AND before OR)
