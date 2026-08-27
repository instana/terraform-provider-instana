# Release Resource

Manages release markers in Instana. Release markers allow you to annotate your monitoring timeline with deployment and release events, making it easy to correlate changes in service behaviour with specific releases or deployments.

API Documentation: <https://developer.ibm.com/apis/catalog/instana--instana-rest-api/api/API--instana--instana-rest-api-documentation#postRelease>

## Example Usage

### Minimal Release

Create a simple release marker with just a name and start timestamp:

```hcl
resource "instana_release" "deployment" {
  name  = "frontend/release-2000"
  start = 1742349976000
}
```

### Release Scoped to Application Perspectives

Associate a release with specific application perspectives so it is visible in those views:

```hcl
resource "instana_release" "backend_deployment" {
  name  = "backend/v1.5.0"
  start = 1742349976000

  applications = [
    { name = "checkout-app" },
    { name = "payment-app" },
  ]
}
```

### Release Scoped to Services

Associate a release with specific services:

```hcl
resource "instana_release" "service_release" {
  name  = "payment-service/v2.0.0"
  start = 1742349976000

  services = [
    { name = "payment" },
    { name = "order-processor" },
  ]
}
```

### Release Scoped to a Service in a Specific Application

Restrict a service scope to a particular application perspective using `scoped_to`:

```hcl
resource "instana_release" "scoped_service_release" {
  name  = "auth-service/v3.1.0"
  start = 1742349976000

  services = [
    {
      name = "auth-service"
      scoped_to = {
        application_name = "checkout-app"
      }
    },
  ]
}
```

### Release Scoped to a Service in a Specific Environment

Restrict a service scope to a particular environment:

```hcl
resource "instana_release" "production_release" {
  name  = "inventory/v4.0.0"
  start = 1742349976000

  services = [
    {
      name = "inventory-service"
      scoped_to = {
        environment_name = "production"
      }
    },
  ]
}
```

### Full Release with Applications and Services

Combine application perspectives, services, and full service scoping in a single release marker:

```hcl
resource "instana_release" "full_release" {
  name  = "platform/v5.0.0"
  start = 1742349976000

  applications = [
    { name = "platform-app" },
  ]

  services = [
    {
      name = "api-gateway"
      scoped_to = {
        application_name = "platform-app"
        environment_name = "production"
      }
    },
    { name = "cache-service" },
  ]
}
```

### Dynamic Release Timestamp

Use Terraform's `timestamp()` function to mark a release at the time of the Terraform apply:

> **Note:** Using `timestamp()` causes the resource to be recreated on every apply since the value changes each run. Use a fixed epoch millisecond value for stable releases.

```hcl
locals {
  # Convert RFC3339 timestamp to milliseconds since epoch
  release_time_ms = parseint(formatdate("YYYYMMDDhhmmss", timestamp()), 10)
}

resource "instana_release" "ci_deployment" {
  name  = "ci/build-${var.build_number}"
  start = 1742349976000  # Use a fixed value in practice
}
```

## Generating Configuration from Existing Resources

If you have already created a release in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

### Steps to Generate Configuration:

1. **Get the Resource ID**: Locate the ID of your release in Instana. You can find this in the Instana UI under Releases or via the API (`GET /api/releases`).

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_release.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired Terraform resource name
- `resource_id` with your actual release ID from Instana (e.g., `Tiu16hLCTniHDtHb_uDV1w`)

3. **Generate the Configuration**: Run the following Terraform command:

```bash
terraform plan -generate-config-out=generated.tf
```

4. **Review and Apply**:

   - **To import the existing resource**: Keep the import block and run `terraform apply`. This will import the release into your Terraform state.

   - **To create a new resource**: Remove the import block, modify the generated configuration as needed, and run `terraform apply` to create a new release.

```bash
terraform apply
```

## Argument Reference

* `name` - Required - The name of the release. For example: `frontend/release-2000`. Maximum 256 characters.
* `start` - Required - The timestamp (in milliseconds since epoch) for when the release occurs. Must be at least `1`. For example: `1742349976000` is Wednesday, 19 March 2025 02:06:16 UTC.
* `applications` - Optional - A list of application perspectives where this release marker will be visible. Maximum 10 entries. Each entry contains: [Details](#application-scope-reference)
* `services` - Optional - A list of services where this release marker will be visible. Maximum 10 entries. Each entry contains: [Details](#service-scope-reference)

### Application Scope Reference

* `name` - Required - The name of the Application Perspective. For example: `app1`. Maximum 256 characters.

### Service Scope Reference

* `name` - Required - The name of the Service. For example: `payment`. Maximum 256 characters.
* `scoped_to` - Optional - An optional scope restriction to limit the service to a specific application or environment. [Details](#scoped-to-reference)

### Scoped To Reference

* `application_name` - Optional - The name of the Application Perspective to scope the service to.
* `environment_name` - Optional - The name of the environment to scope the service to.

## Attributes Reference

* `id` - The unique ID of the release (auto-generated by Instana, e.g., `Tiu16hLCTniHDtHb_uDV1w`).
* `last_updated` - The timestamp (in milliseconds since epoch) of the last update to this release (set by Instana).

## Import

Releases can be imported using the release `id`, e.g.:

```bash
$ terraform import instana_release.example Tiu16hLCTniHDtHb_uDV1w
```

## Notes

* The `id` and `last_updated` fields are computed by Instana and cannot be set in configuration.
* The `start` timestamp is in **milliseconds** since the Unix epoch, not seconds.
* `applications` and `services` each support a maximum of **10** entries per the Instana API.
* Release markers are visible in Instana's Application and Service dashboards when `applications` or `services` scopes are configured.
* Deleting this resource will remove the release marker from Instana.
* Required API token permission: `CanConfigureReleases`.
