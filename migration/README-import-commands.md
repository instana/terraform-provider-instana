# Terraform Import Commands Generator

This script generates Terraform import CLI commands from a `terraform.tfstate` file, including support for module resources.

## Features

- ✅ Generates executable shell script with `terraform import` commands
- ✅ **Full module support** - handles nested modules (e.g., `module.alerts`, `module.app.module.db`)
- ✅ Supports indexed resources (count and for_each)
- ✅ Escapes special characters in resource IDs
- ✅ Groups commands by module for better organization
- ✅ Only processes Instana provider resources (instana_*)

## Usage

### Basic Usage

```bash
# Default: reads terraform.tfstate, outputs to import-commands.sh
go run migration/generate-import-commands.go

# Custom paths
go run migration/generate-import-commands.go path/to/state.tfstate output.sh
```

### Execute Generated Commands

```bash
# Make the script executable
chmod +x import-commands.sh

# Run all import commands
./import-commands.sh

# Or run individual commands as needed
terraform import instana_alerting_channel.example abc123
```

## Module Support

The script fully supports Terraform modules and generates the correct import commands with module paths.

### Example State File Structure

```json
{
  "resources": [
    {
      "mode": "managed",
      "type": "instana_alerting_channel",
      "name": "root_channel",
      "instances": [{"attributes": {"id": "channel-001"}}]
    },
    {
      "module": "module.alerts",
      "mode": "managed",
      "type": "instana_synthetic_alert_config",
      "name": "latency_alert",
      "instances": [{"attributes": {"id": "alert-789"}}]
    },
    {
      "module": "module.application.module.endpoints",
      "mode": "managed",
      "type": "instana_application_config",
      "name": "app_backend",
      "instances": [{"attributes": {"id": "app-456"}}]
    }
  ]
}
```

### Generated Import Commands

```bash
# Root module resources
terraform import instana_alerting_channel.root_channel channel-001

# Module: module.alerts
terraform import module.alerts.instana_synthetic_alert_config.latency_alert alert-789

# Module: module.application.module.endpoints
terraform import module.application.module.endpoints.instana_application_config.app_backend app-456
```

## Indexed Resources Support

The script handles both `count` and `for_each` indexed resources:

### Count Example
```bash
# State: resource with index_key: 0
terraform import instana_alerting_channel.example[0] channel-001

# State: resource with index_key: 1
terraform import instana_alerting_channel.example[1] channel-002
```

### For_each Example
```bash
# State: resource with index_key: "prod"
terraform import instana_alerting_channel.example["prod"] channel-prod

# State: resource with index_key: "staging"
terraform import instana_alerting_channel.example["staging"] channel-staging
```

### Module with Indexed Resources
```bash
# Module resource with count
terraform import module.alerts.instana_synthetic_alert_config.alert[0] alert-001

# Module resource with for_each
terraform import module.alerts.instana_synthetic_alert_config.alert["critical"] alert-002
```

## Output Format

The generated shell script includes:

1. **Header** - Usage instructions and setup
2. **Grouped Commands** - Organized by module (root, then each module)
3. **Comments** - Module names for easy navigation
4. **Footer** - Success message with count

Example output structure:
```bash
#!/bin/bash
# Terraform Import Commands
# Generated from terraform.tfstate

set -e  # Exit on error

echo "Starting Terraform import process..."

# Root module resources
terraform import instana_alerting_channel.team_email email-channel-123
terraform import instana_synthetic_test.homepage_test synthetic-test-001

# Module: module.alerts
terraform import module.alerts.instana_synthetic_alert_config.latency_alert alert-config-789

# Module: module.application.module.endpoints
terraform import module.application.module.endpoints.instana_application_config.app_backend app-backend-456

echo ""
echo "✓ Import process completed successfully!"
echo "✓ Total resources imported: 4"
```

## Important Notes

### Module Directory Structure

If your modules are in subdirectories, you may need to:

1. **Run from root directory** where your main Terraform configuration is located
2. **Use `-chdir` flag** if needed:
   ```bash
   terraform -chdir=path/to/root import module.alerts.instana_synthetic_alert_config.alert alert-123
   ```

### Before Running Imports

1. Ensure you have the correct Terraform configuration files
2. Initialize Terraform: `terraform init`
3. Verify module structure matches your configuration
4. Back up your state file before importing

### Selective Import

You can edit the generated script to:
- Comment out resources you don't want to import
- Run specific module imports separately
- Add custom logic between imports

## Troubleshooting

### Module Not Found Error

If you get "module not found" errors:
```bash
Error: Module not found: module.alerts
```

**Solution**: Ensure your Terraform configuration includes the module definition:
```hcl
module "alerts" {
  source = "./modules/alerts"
  # ... configuration
}
```

### Resource Already Exists

If a resource is already in state:
```bash
Error: Resource already managed by Terraform
```

**Solution**: Remove that import command or use `terraform state rm` first.

### ID Format Issues

Some resources require specific ID formats. Check the provider documentation for the correct format.

## Examples

### Simple Root Resource
```bash
terraform import instana_alerting_channel.email_alerts pTFqA1Uw6ErD0Un5
```

### Module Resource
```bash
terraform import module.monitoring.instana_synthetic_test.api_check test-abc123
```

### Nested Module Resource
```bash
terraform import module.app.module.database.instana_application_config.db_config app-xyz789
```

### Indexed Module Resource
```bash
terraform import module.alerts.instana_synthetic_alert_config.checks["production"] alert-prod-001
```

## Related Scripts

- `migration-script.go` - Generates Terraform import blocks (for Terraform 1.5+)
- `generate-import-commands.go` - Generates CLI import commands (this script)

Both scripts read from the same state file but generate different output formats.