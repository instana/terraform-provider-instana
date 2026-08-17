# Session Settings Resource

Manages the tenant unit session settings in Instana, including the authentication token lifetime and idle session timeout.

API Documentation: [Instana REST API - Session Settings](https://instana.github.io/openapi/#tag/Session-Settings)

---

> ⚠️ **Singleton Resource — One Instance Per Tenant**
>
> `instana_session_settings` is a **tenant-level singleton**. There is exactly one set of session settings per Instana tenant unit. You must declare **at most one** `instana_session_settings` block across your entire Terraform configuration.
>
> Declaring multiple blocks is not supported: every `apply` overwrites the same tenant setting, and only the configuration applied **last** will be in effect. Terraform will not raise an error if you declare two blocks, but the behaviour is undefined and your configuration will be inconsistent.

---

## Example Usage

### Minimal — use server defaults

Both fields are optional. If omitted, Instana's default values are used:
- `token_life_time_in_millis` → `604800000` (7 days)
- `idle_time_in_millis` → `28800000` (8 hours)

```hcl
resource "instana_session_settings" "main" {}
```

### Custom token and idle timeouts

```hcl
resource "instana_session_settings" "main" {
  token_life_time_in_millis = 86400000   # 1 day
  idle_time_in_millis       = 1800000    # 30 minutes
}
```

### Using variables

```hcl
variable "token_lifetime_ms" {
  description = "Authentication token lifetime in milliseconds"
  type        = number
  default     = 604800000 # 7 days
}

variable "idle_timeout_ms" {
  description = "Idle session timeout in milliseconds"
  type        = number
  default     = 28800000 # 8 hours
}

resource "instana_session_settings" "main" {
  token_life_time_in_millis = var.token_lifetime_ms
  idle_time_in_millis       = var.idle_timeout_ms
}
```

### Prevent accidental deletion

Deleting this resource resets the tenant to Instana's built-in defaults. Use `prevent_destroy` in production environments to guard against accidental resets:

```hcl
resource "instana_session_settings" "main" {
  token_life_time_in_millis = 86400000
  idle_time_in_millis       = 3600000

  lifecycle {
    prevent_destroy = true
  }
}
```

---

## Argument Reference

### Optional Attributes

* `token_life_time_in_millis` - (Optional, Computed) Maximum lifetime of an authentication token in milliseconds.
  Defaults to `604800000` (7 days) if not set.
  Valid range: `600000` (10 minutes) to `604800000` (7 days).

  **Type:** `number`

* `idle_time_in_millis` - (Optional, Computed) Idle timeout before a session expires in milliseconds.
  Defaults to `28800000` (8 hours) if not set.
  Valid range: `600000` (10 minutes) to `28800000` (8 hours).

  **Type:** `number`

### Quick reference — common values

| Duration     | Milliseconds  |
|--------------|--------------|
| 10 minutes   | `600000`     |
| 30 minutes   | `1800000`    |
| 1 hour       | `3600000`    |
| 4 hours      | `14400000`   |
| 8 hours      | `28800000`   |
| 1 day        | `86400000`   |
| 7 days       | `604800000`  |

---

## Import

Session settings can be imported by providing any non-empty placeholder string as the ID — the actual value is ignored because this resource has no real ID. The current settings are read directly from the API during import.

### Using an import block (Terraform ≥ 1.5, recommended)

```hcl
import {
  to = instana_session_settings.main
  id = "session_settings"
}
```

Then generate the configuration automatically:

```bash
terraform plan -generate-config-out=generated.tf
```

Review `generated.tf`, then apply:

```bash
terraform apply
```

### Using the CLI import command

```bash
terraform import instana_session_settings.main session_settings
```

---

## Notes

### Singleton behaviour

This resource manages a single tenant-level configuration object. The Instana API exposes it at the fixed path `/api/settings/session` with no per-resource ID. All CRUD operations target that same endpoint:

| Terraform operation | API call           |
|---------------------|--------------------|
| `create` / `update` | `PUT /api/settings/session` |
| `read`              | `GET /api/settings/session` |
| `delete`            | `DELETE /api/settings/session` (reverts to defaults) |

### What happens on delete

Destroying this resource calls `DELETE /api/settings/session`, which **reverts the tenant to Instana's built-in defaults** — it does not permanently remove anything. Use `lifecycle { prevent_destroy = true }` in production to prevent unintentional resets.

### Required API token permission

The API token used by the provider must have the **`CanConfigureSessionSettings`** permission to manage this resource.
