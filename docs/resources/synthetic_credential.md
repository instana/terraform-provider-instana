# Synthetic Credential Resource

Manages synthetic credentials in Instana. Synthetic credentials store encrypted secret values (such as passwords, API tokens, or other sensitive strings) that can be referenced inside synthetic API script tests without exposing them in plain text.

API Documentation: [Synthetic credentials](https://developer.ibm.com/apis/catalog/instana--instana-rest-api/api/API--instana--instana-rest-api-documentation#createSyntheticCredential)

## Example Usage

### Basic Credential

Create a simple credential with a name and a secret value:

```hcl
resource "instana_synthetic_credential" "api_token" {
  credential_name  = "api_token"
  credential_value = "my-secret-token" # use a secret store or variable in practice
}
```

### Credential Scoped to Applications

Restrict which Application Perspectives can use this credential:

```hcl
resource "instana_synthetic_credential" "db_password" {
  credential_name  = "db_password"
  credential_value = var.db_password

  applications = [
    "NQIAsj1tSJm2zxMQx78MSA", # replace with actual application IDs
    "ZpQrTs2uVJn3axNRy89LTA",
  ]
}
```

### Credential with RBAC Tags

Control which teams can view and manage the credential:

```hcl
resource "instana_synthetic_credential" "service_key" {
  credential_name  = "service_key"
  credential_value = var.service_key

  rbac_tags = [
    {
      id           = "team-id-1234"  # replace with actual team ID
      display_name = "Platform Team" # replace with actual team display name
    }
  ]
}
```

## Using Credentials in Synthetic Tests

Synthetic credentials can be referenced inside **script-based** synthetic tests — specifically `http_script` (API script) and `browser_script` tests. 

Inside a script, reference a credential by name using the `$secure.<credential_name>` syntax. The credential value is resolved by the Instana runtime before script execution. Do not build the reference dynamically using variables, string concatenation, or `eval()` — only the literal `$secure.credentialName` format is supported.

> **Note:** Only the `$secure.credentialName` format is supported. The `$secure[credentialName]` bracket notation is not supported.

### HTTP Script (API Script) with Credentials

The following example creates a credential and then references it in an API script test that authenticates against a protected endpoint:

```hcl
resource "instana_synthetic_credential" "api_key" {
  credential_name  = "api_key"
  credential_value = var.api_key
}

resource "instana_synthetic_test" "api_script_with_credentials" {
  label          = "Authenticated API Script Test"
  description    = "Test a protected API endpoint using a stored credential"
  active         = true
  locations      = ["b8dsyQt4fDukWzR9RMXW"] # replace with actual location IDs
  test_frequency = 15
  playback_mode  = "Simultaneous"

  http_script = {
    mark_synthetic_call = true
    retries             = 1
    retry_interval      = 1
    timeout             = "30s"
    script              = <<-EOT
      const assert = require('assert');

      // Reference the stored credential using $secure.<credential_name>
      const response = await $http.get({
        url: 'https://api.example.com/v1/protected-resource',
        headers: {
          'Authorization': 'Bearer ' + $secure.api_key,
          'Content-Type':  'application/json'
        }
      });

      assert.equal(response.statusCode, 200, 'Expected HTTP 200 OK');
    EOT
  }

  depends_on = [instana_synthetic_credential.api_key]
}
```

### Browser Script with Credentials

The following example creates username and password credentials and then references them in a browser script that logs into a web application:

```hcl
resource "instana_synthetic_credential" "login_username" {
  credential_name  = "login_username"
  credential_value = var.login_username
}

resource "instana_synthetic_credential" "login_password" {
  credential_name  = "login_password"
  credential_value = var.login_password
}

resource "instana_synthetic_test" "browser_script_with_credentials" {
  label          = "Login Flow Browser Script Test"
  description    = "Test the login flow using stored credentials"
  active         = true
  locations      = ["b8dsyQt4fDukWzR9RMXW"] # replace with actual location IDs
  test_frequency = 30
  playback_mode  = "Simultaneous"

  browser_script = {
    mark_synthetic_call = true
    retries             = 1
    retry_interval      = 1
    timeout             = "60s"
    browser             = "chrome"
    record_video        = false
    script              = <<-EOT
      // Fill in the login form using stored credentials
      await page.fill('#username', $secure.login_username);
      await page.fill('#password', $secure.login_password);
      await page.click('#submit-btn');

      // Verify successful login
      await page.waitForSelector('h1');
    EOT
  }

  depends_on = [
    instana_synthetic_credential.login_username,
    instana_synthetic_credential.login_password,
  ]
}
```

> **Tip:** The `depends_on` argument is recommended when the credential and the test are defined in the same Terraform configuration. This ensures the credential is created before the synthetic test.

## Generating Configuration from Existing Resources

If you have already created a synthetic credential in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Credential Name**: Locate the name of your synthetic credential in Instana. You can find this in the Instana UI under **Synthetic Settings → Credentials** or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_synthetic_credential.example
  id = "my_credential_name"
}
```

Replace:
- `example` with your desired Terraform block name
- `my_credential_name` with the actual credential name from Instana

3. **Generate the Configuration**: Run the following Terraform command:

```bash
terraform plan -generate-config-out=generated.tf
```

This will:
- Import the existing resource state
- Generate the complete Terraform configuration in `generated.tf`
- Show you what will be imported

4. **Add the Credential Value**: Because the credential value is **write-only** (the API never returns it after creation), the generated configuration will not contain `credential_value`. You must add it manually before running `terraform apply`:

```hcl
# generated.tf — add credential_value before applying
resource "instana_synthetic_credential" "example" {
  credential_name  = "my_credential_name"
  credential_value = var.my_credential_value # add this line
}
```

5. **Review and Apply**: Review the configuration and run:

```bash
terraform apply
```

## Argument Reference

* `credential_name` - **Required** - The unique name that identifies the credential. Used as the resource ID. Must start with a letter and can only contain letters, numbers and underscores. Maximum length is 64 characters. **Changing this value forces the resource to be destroyed and re-created.**
* `credential_value` - **Optional (Required on create/update)** - The secret value of the credential. This field is **write-only**: Instana stores it encrypted and never returns it via the API. The field must be provided when creating or updating the resource. After an `import`, add this attribute to your configuration and run `terraform apply` to complete the import. See [Write-Only Credential Value](#write-only-credential-value) for details.
* `applications` - Optional - Set of Application Perspective IDs that are allowed to use this credential. When empty, the credential is not scoped to any application. Computed: the API returns the current list on read.
* `mobile_apps` - Optional - Set of mobile app IDs that are allowed to use this credential. Computed: preserved from state on read as the API does not return this field.
* `websites` - Optional - Set of website IDs that are allowed to use this credential. Computed: preserved from state on read as the API does not return this field.
* `rbac_tags` - Optional - Set of RBAC tags for access control. Computed: preserved from state on read as the API does not return this field. [Details](#rbac-tags-reference)

### RBAC Tags Reference

* `id` - Required - The ID of the RBAC tag (team).
* `display_name` - Required - The display name of the RBAC tag (team).

## Attributes Reference

* `credential_name` - The name of the synthetic credential (also serves as the resource ID).

## Write-Only Credential Value

The `credential_value` attribute is **write-only**. Instana stores the value encrypted and never exposes it through the API. As a result:

* Terraform cannot detect drift on `credential_value` — if the value is changed in Instana outside of Terraform, a plan will show no diff.
* After `terraform import`, the `credential_value` field will be absent from state. Add it to your configuration and run `terraform apply` to populate it.
* Store sensitive values in a secrets manager or use Terraform input variables marked `sensitive = true` rather than hardcoding them in your configuration.

## Import

Synthetic credentials can be imported using the credential name as the `id`, e.g.:

```bash
$ terraform import instana_synthetic_credential.example my_credential_name
```

After importing, add `credential_value` to your configuration and run `terraform apply` to complete the import.

## Notes

* The credential name serves as both the resource identifier and the Terraform resource ID — changing it forces a destroy and re-create.
* Credential names must start with a letter and may only contain letters, numbers, and underscores (max 64 characters).
* The credential value is encrypted at rest by Instana and is never returned by the API.
* The `applications`, `mobile_apps`, `websites`, and `rbac_tags` fields control the scope and visibility of the credential.
* The `applications` field is returned and refreshed on every read; `mobile_apps`, `websites`, and `rbac_tags` are preserved from state because the API does not return them.
