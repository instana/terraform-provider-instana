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

> **Note**: This approach is **not recommended if you are using Terraform modules**. For module-based configurations, use the [Alternative Approach: Using Import CLI Commands](#alternative-approach-using-import-cli-commands-for-module-based-configurations) instead.

This approach uses Terraform's import blocks to migrate resources without destroying and recreating them, ensuring **zero downtime**.

> **Important**: After successful migration verification, the `migration/` directory becomes your **main working directory**. Keep this folder for ongoing Terraform operations and treat it as your primary workspace going forward.

#### Step 1: Backup Your Existing State

```bash
# Pull and backup your current state file from the old provider
terraform state pull > backup-state.tfstate
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
cp ../backup-state.tfstate ./terraform.tfstate
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

> **Important Note:** After importing, review the values in the generated configuration file carefully. Some values may have been replaced during the import process, and you might need to update references in the generated file manually. For example, if a generated resource uses an ID value that needs to be referenced in subsequent imports or configurations, ensure you update those references with the correct values before proceeding.

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

---

### Alternative Approach: Using Import CLI Commands (For Module-Based Configurations)

If you are using **Terraform modules** in your configuration, the import blocks approach may not work correctly. In this case, use the CLI-based import command approach instead.

> **Note**: This approach is recommended when your Terraform configuration uses modules (e.g., `module.alerts`, `module.monitoring`) or when you prefer more control over the import process.

> **Important**: During migration, work with a **local state file** in the migration directory for verification. Only after successful verification should you replace the original folder and state file.

#### Step 1: Backup Your Existing State

```bash
# Pull and backup your current state file from the old provider
terraform state pull > backup-state.tfstate
```

#### Step 2: Copy Entire Folder to Migration Directory

Create a complete copy of your existing Terraform configuration folder:

```bash
# From your project root directory
# Copy the entire folder to a new migration directory
cp -r . migration/

# Navigate to the migration directory
cd migration/

# Rename the backup state file to terraform.tfstate for local work
mv backup-state.tfstate terraform.tfstate

# Comment out remote backend configuration if you're using one
# Edit your terraform configuration file (e.g., main.tf or backend.tf)
# Comment out the backend block to use local state during migration
```

**Example: Comment out remote backend**
```hcl
terraform {
  required_providers {
    instana = {
      source  = "instana/instana"
      version = "~> 5.0"  # Old version, will be updated in Step 5
    }
  }
  
  # Comment out remote backend to use local state during migration
  # backend "s3" {
  #   bucket = "my-terraform-state"
  #   key    = "prod/terraform.tfstate"
  #   region = "us-east-1"
  # }
}
```

**Important:**
- This preserves your entire configuration structure including modules, variables, and all files
- The backup state file from Step 1 is renamed to `terraform.tfstate` for local work
- The original backup (`backup-state.tfstate`) remains in the parent directory untouched
- Comment out remote backend to work with local state during migration
- This ensures we don't accidentally modify the original remote state file
- Keep the original folder completely untouched as a backup

#### Step 3: Generate Import Commands

Use the provided script to generate CLI import commands from your state file:

```bash
# Download the script from the provider repository
# Or copy it from: https://github.com/instana/terraform-provider-instana/blob/main/migration/generate-import-commands.go

# Run the script to generate import commands (from migration directory)
go run ../migration/generate-import-commands.go terraform.tfstate import-commands.sh
```

**Output:**
```
Reading state file: terraform.tfstate
Output file: import-commands.sh

Generated import command for: instana_alerting_channel.team_email (ID: pTFqA1Uw6ErD0Un5)
Generated import command for: module.alerts.instana_synthetic_alert_config.latency_alert (ID: alert-789)
Generated import command for: module.app.instana_application_config.backend (ID: app-456)

✓ Successfully generated 3 import command(s)
✓ Import commands written to: import-commands.sh
```

The generated script will look like:
```bash
#!/bin/bash
# Terraform Import Commands
# Generated from terraform.tfstate

set -e  # Exit on error

echo "Starting Terraform import process..."

# Root module resources
terraform import instana_alerting_channel.team_email pTFqA1Uw6ErD0Un5

# Module: module.alerts
terraform import module.alerts.instana_synthetic_alert_config.latency_alert alert-789

# Module: module.app
terraform import module.app.instana_application_config.backend app-456

echo ""
echo "✓ Import process completed successfully!"
echo "✓ Total resources imported: 3"
```

#### Step 4: Remove Instana Resources from State

Remove all Instana resources from the state file to prepare for re-import with the new provider:

```bash
# List all Instana resources in the state
terraform state list | grep "instana_"

# Remove all Instana resources at once
terraform state list | grep "instana_" | xargs -I {} terraform state rm {}
```

**Expected Output:**
```
Removed instana_alerting_channel.team_email
Removed module.alerts.instana_synthetic_alert_config.latency_alert
Removed module.app.instana_application_config.backend
```

**Important:**
- This preserves non-Instana resources (e.g., AWS, Azure resources) in your state
- Only Instana resources will be re-imported with the new provider
- The state file now contains only non-Instana resources

#### Step 5: Update Provider Configuration

Update your `main.tf` (or provider configuration file) to use the new Plugin Framework provider version:

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

# Keep all other provider configurations
# provider "aws" { ... }
# provider "azurerm" { ... }
```

**Important:**
- Only update the Instana provider version
- Keep all other provider configurations unchanged
- Keep all module definitions and resource configurations as-is for now

#### Step 6: Initialize with New Provider

```bash
# Remove old provider plugins
rm -rf .terraform/

# Initialize with the new Plugin Framework provider
terraform init
```

**Expected Output:**
```
Initializing provider plugins...
- Finding instana/instana versions matching ">= 6.0.0"...
- Installing instana/instana v6.0.0...
- Installed instana/instana v6.0.0

Terraform has been successfully initialized!
```

This will download the new provider version (>= 6.0.0).

#### Step 8: Run Import Commands

Execute the generated import script to re-import all Instana resources with the new provider:

```bash
# Make the script executable
chmod +x import-commands.sh

# Run all import commands
./import-commands.sh
```

**Expected Output:**
```
Starting Terraform import process...
instana_alerting_channel.team_email: Importing from ID "pTFqA1Uw6ErD0Un5"...
instana_alerting_channel.team_email: Import prepared!
  Prepared instana_alerting_channel for import
instana_alerting_channel.team_email: Import complete!

module.alerts.instana_synthetic_alert_config.latency_alert: Importing from ID "alert-789"...
module.alerts.instana_synthetic_alert_config.latency_alert: Import prepared!
  Prepared instana_synthetic_alert_config for import
module.alerts.instana_synthetic_alert_config.latency_alert: Import complete!

✓ Import process completed successfully!
✓ Total resources imported: 3
```

#### Step 9: Verify State File

```bash
# Check that all resources are properly imported
terraform state list

# Expected output should include all your resources:
# instana_alerting_channel.team_email
# module.alerts.instana_synthetic_alert_config.latency_alert
# module.app.instana_application_config.backend
# aws_instance.example  (if you have non-Instana resources)
```

#### Step 10: Update Resource Configurations Manually

After importing, you need to manually update your resource configurations to match the new provider syntax. Refer to the [Breaking Changes](#-breaking-changes) section and the provider documentation for each resource type.

**Example Updates:**

Old syntax (SDKv2):
```hcl
resource "instana_alerting_channel" "team_email" {
  name = "Team Email"
  email {
    emails = ["team@example.com"]
  }
}
```

New syntax (Plugin Framework):
```hcl
resource "instana_alerting_channel" "team_email" {
  name = "Team Email"
  email = {
    emails = ["team@example.com"]
  }
}
```

**Key Changes to Look For:**
- Nested blocks changed to object syntax (use `=` instead of block syntax)
- Check each resource type in the documentation for specific changes
- Update all Instana resources in your configuration files and modules

#### Step 11: Verify Configuration

```bash
# Run plan to check for any drift or required changes
terraform plan
```

**Expected Output:**
- Should show "No changes. Your infrastructure matches the configuration." if all configurations are correctly updated
- If changes are shown, review and update your configuration accordingly
- Ensure no resources show as "must be replaced"

**If you see changes:**
```bash
# Review the changes carefully
terraform plan

# Update your configuration files to match the expected state
# Re-run plan until you see "No changes"
```

#### Step 12: Final Verification and Backup

Once everything is verified:

```bash
# Create a verified backup of the migrated state
cp terraform.tfstate verified-migration-state.tfstate

# Verify one more time
terraform plan
# Should show: "No changes. Your infrastructure matches the configuration."
```

#### Step 13: Configure Remote Backend (If Previously Used)

If you were using a remote backend before migration, reconfigure it and push the verified state:

```bash
# Uncomment your backend configuration in your terraform files
# Edit main.tf or backend.tf and uncomment the backend block

# Example: Uncomment remote backend
# terraform {
#   backend "s3" {
#     bucket = "my-terraform-state"
#     key    = "prod/terraform.tfstate"
#     region = "us-east-1"
#   }
# }
```

After uncommenting the backend configuration:

```bash
# Re-initialize to configure the remote backend
terraform init

# Push the verified local state to the remote backend
terraform state push terraform.tfstate
```

**Expected Output:**
```
Initializing the backend...
Successfully configured the backend "s3"!

Terraform has been successfully initialized!
```

**Important:**
- Only push the state after thorough verification
- Ensure the remote backend is accessible and properly configured
- The local `terraform.tfstate` file will be uploaded to the remote backend
- Keep a backup of the local state file before pushing

#### Step 14: Replace Original Folder (After Successful Verification)

**Only after thorough verification and testing**, replace your original folder with the migration folder:

```bash
# From the parent directory (not from migration/)
cd ..

# Backup the original folder
mv original-folder original-folder-backup-$(date +%Y%m%d)

# Replace with the migration folder
mv migration/ original-folder/

# Update your CI/CD pipelines to use the updated configuration
```

**Important Notes:**
- **Do NOT replace your production folder until you have thoroughly tested and verified the migration**
- Test the migration in a non-production environment first
- Keep the original folder backup for rollback purposes
- If using remote backend, ensure Step 13 is completed successfully before replacing the folder
- Update all CI/CD pipelines and documentation to reference the updated configuration
- The migrated folder becomes your main working directory after successful verification

#### Troubleshooting

**Module Not Found Error:**
```bash
Error: Module not found: module.alerts
```
**Solution:** Ensure your module definitions exist in `main.tf` and module directories are copied correctly.

**Resource Already Exists:**
```bash
Error: Resource already managed by Terraform
```
**Solution:** The resource is already in state. Skip that import command or use `terraform state rm` first.

**ID Format Issues:**
Some resources require specific ID formats. Check the provider documentation for the correct format. The generated script uses IDs from your existing state, which should be correct.

---

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
