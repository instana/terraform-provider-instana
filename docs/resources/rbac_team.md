# RBAC Team Resource

Management of Teams for role-based access control (RBAC). RBAC teams allow you to define scope access and assign roles to team members for specific resources.

API Documentation: <https://instana.github.io/openapi/#tag/Teams>


## Example Usage

### Basic Team

```hcl
resource "instana_rbac_team" "example" {
  tag = "DevOps Team"
  
  info = {
    description = "Team for DevOps engineers"
  }
}
```

### Team with Members

```hcl
resource "instana_rbac_team" "dev_team" {
  tag = "Development Team Tf"
  
  info = {
    description = "Development team with application access"
  }
  
  member = [
    {
      user_id = "64c17435ad9cd800014d0973" # replace with valid user_id
      roles = [
        {
          role_id = "vCrV7DEdSHS9McI6cU3quA" # replace with valid role_id
        }
      ]
    },
    {
      user_id = "68e8daa9155aab0001bd6144"
      roles = [
        {
          role_id = "-3"
        }
      ]
    }
  ]

    scope = {
    applications = ["app-id-1", "app-id-2", "app-id-3"]
    access_permissions = [
      "LIMITED_APPLICATIONS_SCOPE"
    ]
  }
}
```

### Team with Application Scope

```hcl
resource "instana_rbac_team" "app_team" {
  tag = "Application Team"
  
  info = {
    description = "Team with access to specific applications"
  }
  
  scope = {
    applications = ["app-id-1", "app-id-2", "app-id-3"]
    access_permissions = [
      "LIMITED_APPLICATIONS_SCOPE"
    ]
  }
}
```

### Team with Kubernetes Scope

```hcl
resource "instana_rbac_team" "k8s_team" {
  tag = "Kubernetes Operations"
  
  info = {
    description = "Team managing Kubernetes clusters"
  }
  
  scope = {
    kubernetes_clusters = ["cluster-uuid-1", "cluster-uuid-2"] # provide valid cluster names
    kubernetes_namespaces = ["namespace-uuid-1", "namespace-uuid-2"] # provide valid namespace names
    access_permissions = [
      "LIMITED_KUBERNETES_SCOPE"
    ]
  }
}
```

### Team with Infrastructure Filter

```hcl
resource "instana_rbac_team" "infra_team" {
  tag = "Infrastructure Team"
  
  info = {
    description = "Team with infrastructure access"
  }
  
  scope = {
    infra_dfq_filter = "entity.zone:us-east-1"
    access_permissions = [
      "LIMITED_INFRASTRUCTURE_SCOPE"
    ]
  }
}
```

### Team with Multiple Scopes

```hcl
resource "instana_rbac_team" "platform_team" {
  tag = "Platform Engineering"
  
  info = {
    description = "Platform team with comprehensive access"
  }
  
  scope = {
    applications = ["app-1", "app-2"] # replace with valid application Ids
    kubernetes_clusters = ["k8s-cluster-1"] # replace with valid cluster Ids
    websites = ["website-1"]
    mobile_apps = ["mobile-app-1"]
    synthetic_tests = ["test-1", "test-2"]
    
    access_permissions = [
      "LIMITED_APPLICATIONS_SCOPE",
      "LIMITED_KUBERNETES_SCOPE",
      "LIMITED_WEBSITES_SCOPE",
      "LIMITED_MOBILE_APPS_SCOPE",
      "LIMITED_SYNTHETICS_SCOPE"
    ]
  }
}
```

### Team with Log and Action Filters

```hcl
resource "instana_rbac_team" "monitoring_team" {
  tag = "Monitoring Team"
  
  info = {
    description = "Team with log and action access"
  }
  
  scope = {
    log_filter = "service.name:my-service"
    action_filter = "action.type:deployment"
    access_permissions = [
      "LIMITED_LOGS_SCOPE",
      "LIMITED_AUTOMATION_SCOPE"
    ]
  }
}
```

### Team with SLO Access

```hcl
resource "instana_rbac_team" "slo_team" {
  tag = "SLO Management Team"
  
  info = {
    description = "Team managing service level objectives"
  }
  
  scope = {
    slo_ids = ["slo-id-1", "slo-id-2", "slo-id-3"]    # replace with valid slo Ids
    access_permissions = [
      "LIMITED_SERVICE_LEVEL_SCOPE"
    ]
  }
}
```

### Team with Synthetic Monitoring

```hcl
resource "instana_rbac_team" "synthetic_team" {
  tag = "Synthetic Monitoring Team"
  
  info = {
    description = "Team managing synthetic tests"
  }
  
  scope = {
    synthetic_tests = ["test-1", "test-2"]  # replace with valid synthetic Ids
    synthetic_credentials = ["cred-1", "cred-2"]
    access_permissions = [
      "LIMITED_SYNTHETICS_SCOPE"
    ]
  }
}
```

### Team with Business Perspectives

```hcl
resource "instana_rbac_team" "business_team" {
  tag = "Business Operations"
  
  info = {
    description = "Team with business perspective access"
  }
  
  scope = {
    business_perspectives = ["bp-id-1", "bp-id-2"]
    access_permissions = [
      "LIMITED_BIZOPS_SCOPE"
    ]
  }
}
```

### Team with Restricted Application Filter

```hcl
resource "instana_rbac_team" "restricted_team" {
  tag = "Restricted Access Team"
  
  info = {
    description = "Team with restricted application access"
  }
  
  scope = {
    restricted_application_filter = {
      label = "Production Services"
      scope = "INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING"
      tag_filter_expression = "service.name@dest EQUALS 'butler'"
    }
    access_permissions = [
      "LIMITED_APPLICATIONS_SCOPE"
    ]
  }
}
```

## Generating Configuration from Existing Resources

If you have already created a RBAC team in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your RBAC team in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_rbac_team.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual RBAC team ID from Instana

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

* `tag` - Required - The name/tag of the RBAC team
* `info` - Optional - Additional information about the team [Details](#info-reference)
* `member` - Optional - List of team members [Details](#member-reference)
* `scope` - Optional - Scope configuration for the team [Details](#scope-reference)

### Info Reference

* `description` - Optional - Description of the team

### Member Reference

* `user_id` - Required - The user ID of the team member
* `roles` - Optional - List of roles assigned to the member [Details](#role-reference)

### Role Reference

* `role_id` - Required - The ID of the role

### Scope Reference

* `access_permissions` - Optional - List of access permissions. Allowed values:
  * `LIMITED_APPLICATIONS_SCOPE` - Limited access to applications
  * `LIMITED_WEBSITES_SCOPE` - Limited access to websites
  * `LIMITED_KUBERNETES_SCOPE` - Limited access to Kubernetes
  * `LIMITED_MOBILE_APPS_SCOPE` - Limited access to mobile apps
  * `LIMITED_INFRASTRUCTURE_SCOPE` - Limited access to infrastructure
  * `LIMITED_SYNTHETICS_SCOPE` - Limited access to synthetic monitoring
  * `LIMITED_BIZOPS_SCOPE` - Limited access to business operations
  * `LIMITED_GEN_AI_SCOPE` - Limited access to Gen AI
  * `LIMITED_AUTOMATION_SCOPE` - Limited access to automation
  * `LIMITED_LOGS_SCOPE` - Limited access to logs
  * `LIMITED_ALERT_CHANNELS_SCOPE` - Limited access to alert channels
  * `LIMITED_VSPHERE_SCOPE` - Limited access to vSphere
  * `LIMITED_PHMC_SCOPE` - Limited access to PHMC
  * `LIMITED_POWERVC_SCOPE` - Limited access to PowerVC
  * `LIMITED_ZHMC_SCOPE` - Limited access to zHMC
  * `LIMITED_PCF_SCOPE` - Limited access to PCF
  * `LIMITED_OPENSTACK_SCOPE` - Limited access to OpenStack
  * `LIMITED_SAP_SCOPE` - Limited access to SAP
  * `LIMITED_NUTANIX_SCOPE` - Limited access to Nutanix
  * `LIMITED_XENSERVER_SCOPE` - Limited access to XenServer
  * `LIMITED_WINDOWS_HYPERVISOR_SCOPE` - Limited access to Windows Hypervisor
  * `LIMITED_LINUX_KVM_HYPERVISOR_SCOPE` - Limited access to Linux KVM Hypervisor
  * `LIMITED_AI_GATEWAY_SCOPE` - Limited access to AI Gateway
  * `LIMITED_SERVICE_LEVEL_SCOPE` - Limited access to service levels
* `applications` - Optional - List of application IDs accessible to the team
* `kubernetes_clusters` - Optional - List of Kubernetes cluster IDs accessible to the team
* `kubernetes_namespaces` - Optional - List of Kubernetes namespace IDs accessible to the team
* `mobile_apps` - Optional - List of mobile app IDs accessible to the team
* `websites` - Optional - List of website IDs accessible to the team
* `infra_dfq_filter` - Optional - Infrastructure dynamic focus query filter
* `action_filter` - Optional - Action filter for automation
* `log_filter` - Optional - Log filter query
* `business_perspectives` - Optional - List of business perspective IDs accessible to the team
* `slo_ids` - Optional - List of SLO IDs accessible to the team
* `synthetic_tests` - Optional - List of synthetic test IDs accessible to the team
* `synthetic_credentials` - Optional - List of synthetic credential IDs accessible to the team
* `tag_ids` - Optional - List of tag IDs accessible to the team
* `restricted_application_filter` - Optional - Restricted application filter configuration [Details](#restricted-application-filter-reference)

### Restricted Application Filter Reference

* `label` - Optional - Label for the restricted application filter
* `scope` - Optional - The scope of the filter. Allowed values:
  * `INCLUDE_NO_DOWNSTREAM` - Include no downstream services
  * `INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING` - Include immediate downstream database and messaging
  * `INCLUDE_ALL_DOWNSTREAM` - Include all downstream services
* `tag_filter_expression` - Optional - Tag filter expression for the restricted application filter

## Attributes Reference

* `id` - The ID of the RBAC team

## Import

RBAC Teams can be imported using the `id` of the team, e.g.:

```bash
$ terraform import instana_rbac_team.my_team 60845e4e5e6b9cf8fc2868da
```

## Notes

* The ID is auto-generated by Instana
* Teams provide a way to organize users and control their access to specific resources
* The `tag` field is required and serves as the team's name
* Members can have multiple roles assigned
* Scope configuration allows fine-grained access control to various Instana resources
* The `restricted_application_filter` provides advanced filtering capabilities for application access
* Tag filter expressions use Instana's tag filter syntax
