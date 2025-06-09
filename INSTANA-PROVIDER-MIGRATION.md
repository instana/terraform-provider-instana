# 🚀 Instana Terraform Provider Migration Guide

## Overview

Instana has transitioned its official Terraform provider from a community-maintained version to a new, officially supported one. This guide walks you through how to **safely migrate** from the deprecated provider (`gessnerfl/instana`) to the new one (`instana/instana`).

---

## 🔗 Provider Repositories & Registry Links

### Deprecated Provider
- GitHub: [gessnerfl/terraform-provider-instana](https://github.com/gessnerfl/terraform-provider-instana)
- Terraform Registry: [gessnerfl/instana](https://registry.terraform.io/providers/gessnerfl/instana/latest)

### New Official Provider
- GitHub: [instana/terraform-provider-instana](https://github.com/instana/terraform-provider-instana)
- Terraform Registry: [instana/instana](https://registry.terraform.io/providers/instana/instana/latest)

---

## 🧭 Migration Steps

### Step 1: Keep Both Providers Temporarily 🔄

To ensure a smooth migration, configure **both** providers in your Terraform code:

```hcl
terraform {
  required_providers {
    instana = {
      source  = "gessnerfl/instana"
      version = "1.6.1"
    }

    ibm-instana = {
      source  = "instana/instana"
      version = "3.0.0"
    }
  }
}
```

> 📝 We're using `ibm-instana` as an alias for the new provider to avoid naming conflicts with the old one.

### Step 2: Specify Provider Meta-Arguments 🏷️

When using both the old and new providers in your Terraform configuration, Terraform cannot automatically determine which provider to use for a given resource. To avoid ambiguity, you must **explicitly define the `provider` meta-argument** for each resource.

Update your resource blocks like this:
```hcl
resource "instana_custom_event_spec_threshold_rule" "rule_abc" {
provider = instana
...
}

resource "instana_custom_event_specification" "rule_xyz" {
provider = ibm-instana
...
}
```

**Use:**

- `instana` for the old provider (`gessnerfl/instana`)
- `ibm-instana` for the new provider (`instana/instana`)

### Step 3: Migrate Resources Gradually 🐢

Migrate resources **one at a time** from the old provider to the new one. This approach minimizes risk and helps isolate any issues during the transition.

For each resource you wish to migrate:

1. Update the `provider` meta-argument from the old to the new provider.
2. Run the following Terraform commands to reinitialize and apply the changes:

   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

3. Verify that the resource is now managed by the new provider (`instana/instana`).

> 🛡️ It is highly recommended testing this in a non-production workspace/environment first.

### Step 4: Handle Orphaned Resources 👻

If you encounter an error like this:
```scss
instana_custom_event_spec_threshold_rule.some_rule (orphan) its original provider configuration at provider["registry.terraform.io/gessnerfl/instana"] is required, but it has been removed.
```
This means that resources in your Terraform state still reference the old provider.

To fix this:

- **Do NOT remove the old provider yet.**
- Make sure all resources still associated with `gessnerfl/instana` are either:
    - **Migrated** to the new provider using the `provider` meta-argument
    - **OR Destroyed** if no longer needed

### Step 5: Clean Up Old Provider 🧹

Once all resources are successfully migrated and no references remain to the old provider:

1. Remove all old provider references in your Terraform code.
2. Delete the old provider from your `required_providers` block:

```hcl
terraform {
  required_providers {
    ibm-instana = {
      source  = "instana/instana"
      version = "3.0.0"
    }
  }
}
```
---

## ✅ Final Tips

- Use `git` to track each step of your migration and allow for safe rollback.
- **Apply in stages** and verify behavior using `terraform plan` and `terraform state list`.

### 🔄 Apply in Stages

This means you **don’t try to migrate everything at once**. Instead:

1. **Pick a small set of resources** to migrate first.
2. Make the necessary Terraform changes (e.g., update provider blocks, meta-arguments).
3. Test just those changes using Terraform.
4. Confirm the outcome is as expected.
5. Repeat the process for the next group of resources.

> ⚠️ This reduces risk — if something goes wrong, it's easier to debug and roll back.

---

### 🔍 Verify Behavior Using `terraform plan` and `terraform state list`

#### ✅ `terraform plan`
- Shows what Terraform **intends to do** based on your code and current state.

#### 📋 `terraform state list`
- Lists all resources Terraform currently tracks in its state file.
- Useful for verifying:
    - Whether a resource is still associated with the old provider.
    - If a resource is correctly tracked under the new provider.

---

### 🧠 Example Workflow

```bash
# Step 1: Preview changes
terraform plan

# Step 2: Apply if everything looks correct
terraform apply

# Step 3: Verify state
terraform state list
```


