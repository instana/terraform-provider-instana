# 🚀 Instana Terraform Provider - SDKv2 to Plugin Framework Migration Guide

## Overview

This guide helps you understand the migration from Terraform Plugin SDKv2 to the newer Terraform Plugin Framework for the Instana Terraform Provider. This migration brings improved type safety, better validation, and enhanced developer experience.

---

## 📋 What Changed?

### High-Level Changes

The Instana Terraform Provider has been migrated from **Terraform Plugin SDKv2** to **Terraform Plugin Framework**. This is an internal architectural change that improves:

- **Type Safety**: Stronger typing with framework-native types
- **Validation**: Built-in validators for better input validation
- **State Management**: More explicit and predictable state handling
- **Error Handling**: Improved diagnostic messages
- **Performance**: Better resource lifecycle management

---

## 🔍 Key Difference: Block Schema vs Attribute Schema

### The Critical Change

The most significant change is the **internal schema representation**:

**SDKv2 (Old):**
- Used **block schema** (`schema.TypeList` with `Elem: &schema.Resource{}`)
- Nested configurations defined as blocks

**Plugin Framework (New):**
- Uses **attribute schema** (`schema.SingleNestedAttribute`, `schema.ListNestedAttribute`)
- Nested configurations defined as attributes
- **Block schema is NOT supported**

### Impact

⚠️ **Resources must be recreated** because the internal schema representation changed from blocks to attributes.

> **Note**: Your HCL configuration syntax remains identical, but Terraform will detect the schema change and require resource recreation.

---

## 📝 Examples

### Example 1: Simple Resource

**Your Configuration (Unchanged):**
```hcl
resource "instana_api_token" "example" {
  name                           = "my-api-token"
  can_configure_service_mapping  = true
  can_configure_eum_applications = true
}
```

**What Happens:**
```bash
terraform plan
# Output:
# instana_api_token.example must be replaced
-/+ resource "instana_api_token" "example" {
    ~ id   = "abc123" -> (known after apply)
      name = "my-api-token"
      # (forces replacement due to schema change)
}
```

---

### Example 2: Nested Configuration - Syntax Change Required

**SDKv2 (Old - Block Syntax):**
```hcl
resource "instana_automation_action" "http_action" {
  name        = "my-http-action"
  description = "HTTP action example"
  
  # Block syntax (OLD)
  http {
    host   = "https://example.com"
    method = "POST"
    
    auth {
      basic_auth {
        username = "user"
        password = "pass"
      }
    }
  }
}
```

**Plugin Framework (New - Attribute Syntax):**
```hcl
resource "instana_automation_action" "http_action" {
  name        = "my-http-action"
  description = "HTTP action example"
  
  # Attribute syntax (NEW) - Note the = sign
  http = {
    host   = "https://example.com"
    method = "POST"
    
    auth = {
      basic_auth = {
        username = "user"
        password = "pass"
      }
    }
  }
}
```

**Key Change:**
- **Block syntax**: `http { ... }` (no equals sign)
- **Attribute syntax**: `http = { ... }` (with equals sign)
- You must update your configuration to use `=` for nested objects

---

## 🔧 Migration Steps

### Recommended Approach: Using Import Blocks (Zero Downtime)

This approach uses Terraform's import blocks to migrate resources without destroying and recreating them, ensuring **zero downtime**.

> **Important**: After successful migration verification, the `migration/` directory becomes your **main working directory**. Keep this folder for ongoing Terraform operations and treat it as your primary workspace going forward.

#### Step 1: Backup Your Existing State

```bash
# Pull and backup your current state file from the old provider
terraform state pull > backup-$(date +%Y%m%d-%H%M%S).tfstate
```

#### Step 2: Create Migration Directory

```bash
# Create a new directory for migration
mkdir migration
cd migration
```

#### Step 3: Copy State File to Migration Directory

```bash
# Copy the existing state file to the migration directory
cp ../terraform.tfstate ./terraform.tfstate
```

#### Step 4: Generate Import Blocks from State

Use the provided migration script located in the `migration/` folder to automatically generate import blocks from your state file:

```bash
# Download the migration script from the provider repository
# Or copy it from: https://github.com/instana/terraform-provider-instana/blob/main/migration/migration-script.go

# Run the script to generate import blocks
go run migration-script.go terraform.tfstate import.tf
```

**Output:**
```
Reading state file: terraform.tfstate
Output file: import.tf

Generated import block for: instana_infra_alert_config.my_alert (ID: abc123...)
Generated import block for: instana_log_alert_config.my_log_alert (ID: xyz789...)

✓ Successfully generated 2 import block(s)
✓ Import blocks written to: import.tf
```

#### Step 5: Remove Instana Resources from State

To preserve non-Instana resources in your state, remove only the Instana resources:

```bash
# List all Instana resources in the state
terraform state list | grep "instana_"

# Remove each Instana resource from the state
# Replace with your actual resource names from the list above
terraform state rm instana_infra_alert_config.my_alert
terraform state rm instana_log_alert_config.my_log_alert
# ... repeat for all Instana resources
```

**Alternative - Automated removal:**
```bash
# Remove all Instana resources at once
terraform state list | grep "instana_" | xargs -I {} terraform state rm {}
```

**Important:**
- This preserves non-Instana resources (e.g., AWS, Azure resources) in your state
- Only Instana resources will be re-imported with the new provider
- If you have ONLY Instana resources, you can delete the entire state file instead: `rm terraform.tfstate`

#### Step 6: Set Up Terraform Configuration

Create a `main.tf` file with the updated provider version and your resource definitions:

```hcl
terraform {
  required_providers {
    instana = {
      source  = "instana/instana"
      version = ">= 6.0.0"  # Plugin Framework version
    }
    # Keep other provider configurations (AWS, Azure, etc.)
  }
}

provider "instana" {
  api_token = var.instana_api_token
  endpoint  = var.instana_endpoint
}

# Keep other provider configurations
# provider "aws" { ... }
# provider "azurerm" { ... }

# Copy Non Instana resource definitions from the old configuration to main.tf if exists

# Non-Instana resources (keep as-is)
# resource "aws_instance" "example" { ... }
# resource "azurerm_resource_group" "example" { ... }

```

**Important:** Copy ALL resource definitions (non-Instana) to maintain consistency. Only Instana resources need syntax updates.

#### Step 7: Initialize with New Provider

```bash
# Initialize with the new Plugin Framework provider
terraform init
```

This will download the new provider version (>= 6.0.0) and create a fresh state.

#### Step 8: Plan the Import

```bash
terraform plan -generate-config-out=generated.tf
```

**Expected Output:**
- Terraform will show that resources will be imported (not created)
- Import blocks will bring existing resources into the new state
- Configuration changes will be shown if syntax updates are needed
- No resources should show as "must be replaced"

Example output:
```
Terraform will perform the following actions:

  # instana_log_alert_config.my_log_alert will be imported
    resource "instana_log_alert_config" "my_log_alert" {
        id          = "YzrwxoGDS0mY77DBk_8EqQ"
        name        = "My Log Alert"
        # ... other attributes
    }

Plan: 2 to import, 0 to add, 0 to change, 0 to destroy.
```

#### Step 9: Apply the Migration

```bash
terraform apply
```

This will import all resources into the **new state file** with the Plugin Framework provider. The resources themselves are not recreated - only the state representation changes.

#### Step 10: Verify Migration

```bash
# Check that all resources are properly managed in the new state
terraform state list

# Verify no drift
terraform plan
# Should show: "No changes. Your infrastructure matches the configuration."
```

#### Step 11: Clean Up and Finalize

```bash
# Remove import blocks from main.tf (they're only needed once)
# Edit main.tf and delete all import { } blocks
# Keep only the resource definitions

# Keep the old state backup for rollback purposes (in parent directory)
# The new state is in migration/terraform.tfstate
```

**Important Notes:**
- The **new state file** (`migration/terraform.tfstate`) is created by the import process
- The **old state file** (backup) should be kept for rollback purposes in the parent directory
- Do NOT copy the old state file to the migration directory - let Terraform create a new one
- **After successful migration verification, the `migration/` directory becomes your main working directory**
- Continue all future Terraform operations from the `migration/` folder
- Update your CI/CD pipelines and documentation to reference the new directory path


## 🎯 Best Practices

### 1. Test in Non-Production First

Always test the migration in a development or staging environment:

```bash
terraform plan
terraform apply
```

### 2. Schedule Maintenance Window

Coordinate with your team as resources will be recreated:
- API tokens will be regenerated
- Alert configurations will be briefly unavailable
- Monitoring may have gaps during recreation

### 3. Use Version Constraints

Pin your provider version to avoid unexpected upgrades:

```hcl
terraform {
  required_providers {
    instana = {
      source  = "instana/instana"
      version = "~> 6.0.0"  # Allow patch updates only
    }
  }
}
```

---

## 🐛 Common Issues

### Issue: All Resources Show as "Must Be Replaced"

**Symptom:**
```
# All resources show:
-/+ resource "instana_xxx" "example" {
    # (forces replacement due to schema change)
}
```

**Solution:**
This is **expected behavior**. The schema change from blocks to attributes requires resource recreation. This cannot be avoided.

---

### Issue: Validation Errors During Plan

**Symptom:**
```
Error: Invalid Attribute Value
│ 
│ Attribute method must be one of: GET, POST, PUT, DELETE
```

**Solution:**
The Plugin Framework provides better validation. Fix the attribute value to meet the requirements. This error would have occurred at API level in SDKv2, now it's caught earlier.

---

## 🔄 Rollback Procedure

If you need to rollback:

### Step 1: Restore State Backup

```bash
terraform state push backup.tfstate
```

### Step 2: Downgrade Provider

```hcl
terraform {
  required_providers {
    instana = {
      source  = "instana/instana"
      version = "5.x.x"  # Previous SDKv2 version
    }
  }
}
```

### Step 3: Reinitialize

```bash
terraform init -upgrade
terraform plan
```

---

## ✅ Summary

### Key Points:

- ✅ **Zero-downtime migration possible** - Use import blocks to migrate without recreating resources
- ⚠️ **Configuration syntax change required** - Update from block syntax to attribute syntax (add `=` sign)
- ✅ **Automated import generation** - Use `migration-script.go` to generate import blocks automatically
- ✅ **New state file created** - Import process creates a fresh state with the new provider
- ✅ **Migration folder becomes main directory** - After verification, use `migration/` as your primary workspace
- ✅ **Better validation** - Errors caught at plan time instead of apply time
- ✅ **Improved type safety** - More reliable state management
- ⚠️ **Alternative approach available** - Can recreate resources if import is not suitable (requires downtime)

### Migration Checklist (Import-Based Approach - Recommended):

- [ ] Backup old state file (`terraform state pull > backup-$(date +%Y%m%d-%H%M%S).tfstate`)
- [ ] Create migration directory (`mkdir migration && cd migration`)
- [ ] Copy existing state file to migration directory (`cp ../terraform.tfstate ./terraform.tfstate`)
- [ ] Download `migration-script.go` from the provider repository
- [ ] Run migration script to generate import blocks (`go run migration-script.go terraform.tfstate import.tf`)
- [ ] Remove only Instana resources from state (`terraform state list | grep "instana_" | xargs -I {} terraform state rm {}`)
  - This preserves non-Instana resources (AWS, Azure, etc.) in the state
  - If you have ONLY Instana resources, you can delete the entire state file instead
- [ ] Create `main.tf` with updated Instana provider version (>= 6.0.0)
- [ ] Copy ALL provider configurations (Instana, AWS, Azure, etc.)
- [ ] Copy ALL resource definitions from old configuration to `main.tf`
  - Include both Instana and non-Instana resources
  - Non-Instana resources remain unchanged
- [ ] Update Instana resource syntax (block `{}` to attribute `= {}`) where needed
- [ ] Test in non-production environment first
- [ ] Run `terraform init` (downloads new Instana provider)
- [ ] Review `terraform plan` output:
  - Should show imports for Instana resources only
  - Should show no changes for non-Instana resources
- [ ] Apply migration (`terraform apply` - imports Instana resources into state)
- [ ] Verify all resources are properly managed (`terraform state list`)
- [ ] Verify no drift (`terraform plan` should show no changes)
- [ ] Remove import blocks from `main.tf` (only needed once)
- [ ] Keep old state backup for rollback purposes
- [ ] **After verification, keep the `migration/` folder as your main working directory**

## 📚 Additional Resources

- [Terraform Plugin Framework Documentation](https://developer.hashicorp.com/terraform/plugin/framework)
- [Instana Provider Documentation](https://registry.terraform.io/providers/instana/instana/latest/docs)
- [GitHub Issues](https://github.com/instana/terraform-provider-instana/issues)

---
