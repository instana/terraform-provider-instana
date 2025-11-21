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

### Step 1: Backup Your State

```bash
terraform state pull > backup.tfstate
```

### Step 2: Update Provider Version

```hcl
terraform {
  required_providers {
    instana = {
      source  = "instana/instana"
      version = ">= 6.0.0"  # Plugin Framework version
    }
  }
}
```

### Step 3: Initialize and Plan

```bash
terraform init -upgrade
terraform plan
```

**Expected Output:**
- All existing resources will show as "must be replaced"
- This is normal and expected due to schema changes

### Step 4: Minimize Downtime (Optional)

Use `create_before_destroy` lifecycle rule:

```hcl
resource "instana_api_token" "example" {
  name = "my-api-token"
  
  lifecycle {
    create_before_destroy = true
  }
}
```

### Step 5: Apply Changes

```bash
terraform apply
```

⚠️ **Warning**: Resources will be destroyed and recreated. Plan for downtime or maintenance window.

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

- ⚠️ **Configuration syntax change required** - Update from block syntax to attribute syntax (add `=` sign)
- ⚠️ **Resource recreation IS required** - Block schema is not supported in Plugin Framework for Instana resources
- ⚠️ **Plan for downtime** - Resources will be destroyed and recreated
- ✅ **Better validation** - Errors caught at plan time instead of apply time
- ✅ **Improved type safety** - More reliable state management

### Migration Checklist:

- [ ] Backup state files
- [ ] Update configuration syntax (block `{}` to attribute `= {}`)
- [ ] Test in non-production environment
- [ ] Schedule maintenance window
- [ ] Update provider version
- [ ] Run `terraform init -upgrade`
- [ ] Review `terraform plan` output
- [ ] Apply changes during maintenance window
- [ ] Verify resources are working correctly

---

## 📚 Additional Resources

- [Terraform Plugin Framework Documentation](https://developer.hashicorp.com/terraform/plugin/framework)
- [Instana Provider Documentation](https://registry.terraform.io/providers/instana/instana/latest/docs)
- [GitHub Issues](https://github.com/instana/terraform-provider-instana/issues)

---
