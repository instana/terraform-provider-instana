# IP Filtering Resource

Manages tenant unit IP filtering settings in Instana, including allow/block traffic rules based on IP address or range, fallback deny-all setting, and support access.

API Documentation: [Instana REST API - IP Filtering](https://instana.github.io/openapi/#tag/IP-Filtering)

---

> ⚠️ **Singleton Resource — One Instance Per Tenant**
>
> `instana_ip_filtering` is a **tenant-level singleton**. There is exactly one IP filtering configuration per Instana tenant unit. You must declare **at most one** `instana_ip_filtering` block across your entire Terraform configuration.
>
> Declaring multiple blocks is not supported: every `apply` overwrites the same tenant setting, and only the configuration applied **last** will be in effect.

---

## Example Usage

### Basic Example — Allow specific IP and subnet, deny all others

```hcl
resource "instana_ip_filtering" "main" {
  enabled              = true
  deny_all             = true
  support_access_enabled = true

  rules = [
    {
      target = "192.168.1.0/24"
      block  = false
    },
    {
      target = "10.0.0.1"
      block  = false
    }
  ]
}
```

### Block specific IPs while allowing others

```hcl
resource "instana_ip_filtering" "main" {
  enabled              = true
  deny_all             = false
  support_access_enabled = false

  rules = [
    {
      target = "198.51.100.23"
      block  = true
    },
    {
      target = "203.0.113.0/24"
      block  = true
    }
  ]
}
```

### Prevent accidental deletion

Deleting this resource calls `DELETE /api/settings/ip-filtering`, removing the IP filtering configuration. Use `prevent_destroy` in production environments to guard against accidental deletion:

```hcl
resource "instana_ip_filtering" "main" {
  enabled  = true
  deny_all = true

  rules = [
    {
      target = "10.0.0.0/8"
      block  = false
    }
  ]

  lifecycle {
    prevent_destroy = true
  }
}
```

---

## Argument Reference

### Required Attributes

* `enabled` - (Required) If the IP filtering configuration should be enforced. See `active` status for whether it is currently enforced (temporary verification period vs permanent).
  **Type:** `bool`

* `deny_all` - (Required) Fallback configuration indicating if unmatched requests should be denied (`true`) or allowed (`false`).
  **Type:** `bool`

* `rules` - (Required) Traffic rules enforcing allow or block based on incoming IP address or range. Rules are evaluated in the given order. Between 1 and 75 rules are required.
  **Type:** `list(object)`

  Each `rules` item supports the following:
  * `target` - (Required) Target of this rule. Can be an individual IPv4/IPv6 address (e.g. `127.0.0.1`, `::1`) or an IP address range in CIDR notation (e.g. `10.0.0.0/8`, `192.168.1.0/24`).
    **Type:** `string`
  * `block` - (Required) Indicates if matching traffic should be blocked/denied (`true`) or allowed (`false`).
    **Type:** `bool`

### Optional Attributes

* `support_access_enabled` - (Optional, Computed) Flag for support access functionality. Defaults to `false`.
  **Type:** `bool`

### Read-Only / Computed Attributes

* `active` - (Computed) Indicates if the configuration is currently enforced. Can be permanently active (if verified), or temporary after creation (during verification window).
  **Type:** `bool`

* `last_changed_at` - (Computed) Unix timestamp in milliseconds when the configuration was last changed.
  **Type:** `number`

* `last_verified_at` - (Computed) Unix timestamp in milliseconds when the configuration was last verified.
  **Type:** `number`

---

## Import

IP filtering can be imported using any non-empty placeholder string as the ID — the actual value is ignored because this resource is a singleton with no per-resource ID:

### Using an import block (Terraform ≥ 1.5, recommended)

```hcl
import {
  to = instana_ip_filtering.main
  id = "ip_filtering"
}
```

### Using the CLI import command

```bash
terraform import instana_ip_filtering.main ip_filtering
```

---

## Notes

### Required API token permission

The API token used by the provider must have the **`CanConfigureIPFiltering`** permission to manage this resource.
