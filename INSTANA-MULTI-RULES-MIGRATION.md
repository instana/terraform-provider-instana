# 🚀 Instana Terraform Provider Multi-Rule Migration Guide

## Overview

Starting with provider version **v7.10.0**, three resources support **multiple rules** in a single configuration. To enable this, the `rules` attribute (or the `threshold` rule inside `rules`) changed from a **single nested object** to a **list of objects** (`ListNestedAttribute`). This is a **breaking change for your `.tf` files**. This guide walks you through how to **safely upgrade** your existing configurations.

| Resource | What changed | Rule limit | State schema version |
|----------|--------------|------------|----------------------|
| `instana_infra_alert_config` | `rules = { ... }` → `rules = [ { ... } ]` | 1–5 rules | 2 → 3 |
| `instana_log_alert_config` | `rules = { ... }` → `rules = [ { ... } ]` | At least 1 (the API currently accepts exactly 1) | 1 → 2 |
| `instana_custom_event_specification` | `rules = { threshold = { ... } }` → `rules = { threshold = [ { ... } ] }` | 1–5 threshold rules | 1 → 2 |

> 📝 Multiple rules are combined with the top-level `rule_logical_operator` (`AND` or `OR`, default `AND`).

---

## 🔗 Provider Links

- GitHub: [instana/terraform-provider-instana](https://github.com/instana/terraform-provider-instana)
- Terraform Registry: [instana/instana](https://registry.terraform.io/providers/instana/instana/latest)
- Resource documentation:
    - [`instana_infra_alert_config`](docs/resources/infra_alert_config.md)
    - [`instana_log_alert_config`](docs/resources/log_alert_config.md)
    - [`instana_custom_event_specification`](docs/resources/custom_event_specification.md)

---

## 🧭 Migration Steps

### Step 1: Identify Affected Resources 🔎

Find every resource of the three types above in your Terraform code:

```bash
grep -rnE 'resource "instana_(infra_alert_config|log_alert_config|custom_event_specification)"' .
```

Only the following need to change:

- `instana_infra_alert_config` and `instana_log_alert_config`: every resource.
- `instana_custom_event_specification`: only resources that use a `threshold` rule. The other rule types (`entity_count`, `entity_count_verification`, `entity_verification`, `host_availability`, `system`) are **unchanged**.

### Step 2: Update Your Configuration ✏️

Wrap the affected definition in square brackets `[ ... ]`. Your values do **not** change.

#### `instana_infra_alert_config`

##### Before (v6.0.0 – v7.9.4):
```hcl
rules = {
  generic_rule = {
    aggregation              = "MEAN"
    cross_series_aggregation = "MEAN"
    entity_type              = "kubernetesHorizontalPodAutoscaler"
    metric_name              = "maxReplicas"
    threshold_operator       = ">="
    threshold = {
      critical = {
        static = {
          value = 90
        }
      }
    }
  }
}
```

##### After (v7.10.0+):
```hcl
rules = [
  {
    generic_rule = {
      aggregation              = "MEAN"
      cross_series_aggregation = "MEAN"
      entity_type              = "kubernetesHorizontalPodAutoscaler"
      metric_name              = "maxReplicas"
      threshold_operator       = ">="
      threshold = {
        critical = {
          static = {
            value = 90
          }
        }
      }
    }
  }
]
```

#### `instana_log_alert_config`

##### Before (v6.0.0 – v7.9.4):
```hcl
rules = {
  metric_name        = "logCount"
  alert_type         = "log.count"
  aggregation        = "SUM"
  threshold_operator = ">="
  threshold = {
    critical = {
      static = {
        value = 2
      }
    }
  }
}
```

##### After (v7.10.0+):
```hcl
rules = [
  {
    metric_name        = "logCount"
    alert_type         = "log.count"
    aggregation        = "SUM"
    threshold_operator = ">="
    threshold = {
      critical = {
        static = {
          value = 2
        }
      }
    }
  }
]
```

> ⚠️ The Instana API currently accepts **exactly one** rule per log alert configuration. Configuring more than one rule is rejected by the API at apply time until the backend lifts this limit.

#### `instana_custom_event_specification`

Only the `threshold` rule becomes a list. The `rules` attribute itself stays an object.

##### Before (v6.0.0 – v7.9.4):
```hcl
rules = {
  threshold = {
    severity           = "critical"
    metric_name        = "cpu.usage"
    rollup             = 0
    window             = 60000
    aggregation        = "avg"
    condition_operator = ">"
    condition_value    = 1.0
  }
}
```

##### After (v7.10.0+):
```hcl
rules = {
  threshold = [
    {
      severity           = "critical"
      metric_name        = "cpu.usage"
      rollup             = 0
      window             = 60000
      aggregation        = "avg"
      condition_operator = ">"
      condition_value    = 1.0
    }
  ]
}
```

### Step 3: Upgrade the Provider ⬆️

Update the provider version constraint and reinitialize:

```hcl
terraform {
  required_providers {
    instana = {
      source  = "instana/instana"
      version = ">= 7.10.0"
    }
  }
}
```

```bash
terraform init -upgrade
```

> 🛡️ It is highly recommended testing this in a non-production workspace/environment first.

### Step 4: Verify with `terraform plan` ✅

```bash
terraform plan
```

- Because you only wrapped existing values in `[ ... ]`, the plan should show **no changes** for the migrated resources.
- If the plan reports a type error on `rules` or `threshold`, you missed a spot. Check that the brackets are in place for that resource.

### Step 5: Apply 🚢

```bash
terraform apply
```

---

## 💾 State Migration

You **do not** need to edit or migrate your state manually. The provider's built-in state upgrader runs automatically the first time Terraform reads your state with v7.10.0+:

| Resource | Schema version upgrade | What the upgrader does |
|----------|------------------------|------------------------|
| `instana_infra_alert_config` | 2 → 3 | Wraps the single `rules` object into a one-element list |
| `instana_log_alert_config` | 1 → 2 | Wraps the single `rules` object into a one-element list |
| `instana_custom_event_specification` | 1 → 2 | Wraps the single `threshold` rule into a one-element list |

> 📝 Only your `.tf` configuration files need to be updated.

---

## ➕ Using Multiple Rules

After migrating, you can add more rules to the list.

### `instana_infra_alert_config` (1–5 rules)

```hcl
resource "instana_infra_alert_config" "multi_rule_alert" {
  name            = "Multi-Rule Infrastructure Alert"
  description     = "Alert combining CPU and memory rules with OR logic"
  evaluation_type = "PER_ENTITY"
  granularity     = 600000

  rule_logical_operator = "OR"

  rules = [
    {
      generic_rule = {
        aggregation              = "MEAN"
        cross_series_aggregation = "MEAN"
        entity_type              = "kubernetesPod"
        metric_name              = "cpuUsageToLimitRatio"
        threshold_operator       = ">="
        threshold = {
          critical = {
            static = {
              value = 90
            }
          }
        }
      }
    },
    {
      generic_rule = {
        aggregation              = "MEAN"
        cross_series_aggregation = "MEAN"
        entity_type              = "kubernetesPod"
        metric_name              = "memoryUsageToLimitRatio"
        threshold_operator       = ">="
        threshold = {
          critical = {
            static = {
              value = 85
            }
          }
        }
      }
    }
  ]

  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
}
```

### `instana_custom_event_specification` (1–5 threshold rules)

See [Multiple Threshold Rules with Logical Operator](docs/resources/custom_event_specification.md) in the resource documentation for a full example. When configuring multiple threshold rules:

- All threshold rules must share the same `severity`, `window` and `rollup`.
- Threshold rules cannot be combined with other rule types. Only one rule in total is allowed unless all rules are of type `threshold`.
- Entity types `application`, `service` and `endpoint` do not support multiple rules.
- Combine them with `rule_logical_operator` (`AND` by default, or `OR`).

---

## ✅ Final Tips

- Use `git` to track the change and allow for safe rollback.
- **Migrate one resource type at a time** and run `terraform plan` between each step.
- Run `terraform plan` before `terraform apply` and confirm there are **no unexpected changes**.
- If you generate configuration from existing resources (`terraform plan -generate-config-out`), check that `rules` / `threshold` is emitted as a list.

---

### 🧠 Example Workflow

```bash
# Step 1: Upgrade the provider
terraform init -upgrade

# Step 2: After wrapping rules/threshold in [ ... ], preview changes
terraform plan

# Step 3: Apply if everything looks correct (no unexpected changes)
terraform apply

# Step 4: Verify state
terraform state list
```
