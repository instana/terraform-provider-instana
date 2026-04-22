# Basic Instana Terraform Provider Example

This is a modular example of using the Instana Terraform provider (version >= 7.0.0) with S3 backend/local for state storage. Resources are organized into reusable modules for better maintainability and scalability.

## Overview

This example creates:
- **Application Module**: Application perspective for monitoring
- **Alerting Module**: Alert channels (Email, Slack) and alert configurations (Slowness, Error Rate, Throughput, Log Alerts)
- **Automation Module**: Automation actions and policies for incident response
- **Dashboard Module**: Custom dashboards for visualization
- **Infrastructure Module**: Infrastructure monitoring and alert configurations
- **Mobile Module**: Mobile application monitoring and alerts
- **RBAC Module**: Role-based access control (Developer, Viewer, Admin roles)
- **SLO Module**: Service Level Objectives and SLI configurations
- **Synthetic Module**: Synthetic monitoring tests and alerts
- **Website Module**: Website monitoring and alert configurations

## Architecture

```
examples/reference-architecture/
├── environments/              # Environment-specific configurations
│   ├── dev/                  # Development environment
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── backend.hcl
│   │   └── terraform.tfvars.example
│   ├── staging/              # Staging environment
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── backend.hcl
│   │   └── terraform.tfvars.example
│   └── production/           # Production environment
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       ├── backend.hcl
│       └── terraform.tfvars.example
├── modules/                  # Shared reusable modules
│   ├── application/          # Application configuration module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   ├── alerting/             # Alerting channels and alerts module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   ├── automation/           # Automation actions and policies module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   ├── dashboard/            # Custom dashboards module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   ├── infrastructure/       # Infrastructure monitoring module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   ├── mobile/               # Mobile app monitoring module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   ├── rbac/                 # RBAC roles module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   ├── slo/                  # SLO configuration module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   ├── synthetic/            # Synthetic monitoring module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   ├── website/              # Website monitoring module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   └── [PLACEHOLDER]/               # [PLACEHOLDER] Add your additional module here
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       └── README.md
└── README.md                 # This file
```

## Prerequisites

1. **Terraform**: Version 1.0.0 or higher (recommended 1.5.0+)
2. **Instana Account**: Active Instana tenant with API access
3. **AWS Account**: For S3 backend state storage (optional)
4. **Instana API Token**: Generate from Instana UI (Settings → API Tokens)

## Quick Start

This example uses a multi-environment structure. Choose your target environment (dev, staging, or production) and follow these steps:

### 1. Choose Your Environment

```bash
# Navigate to your target environment
cd environments/dev        # For development
# OR
cd environments/staging    # For staging
# OR
cd environments/production # For production
```

### 2. Set Environment Variables

```bash
# Required: Instana authentication
export INSTANA_API_TOKEN="your-api-token-here"
export INSTANA_ENDPOINT="your-tenant.instana.io"

# Optional: AWS credentials for S3 backend
export AWS_ACCESS_KEY_ID="your-aws-access-key"
export AWS_SECRET_ACCESS_KEY="your-aws-secret-key"
export AWS_REGION="us-east-1"
```

### 3. Configure Backend (S3)

Edit `backend.hcl` in your chosen environment directory with your S3 bucket details:

```hcl
# For dev environment
bucket         = "your-terraform-state-bucket"
key            = "instana/dev/terraform.tfstate"
region         = "us-east-1"
encrypt        = true
dynamodb_table = "terraform-state-lock"
```

**Note:** Each environment uses a different state file key to maintain isolation.

### 4. Initialize Terraform

```bash
# From within your environment directory (e.g., environments/dev/)
terraform init -backend-config=backend.hcl
```

### 5. Configure Variables

Create `terraform.tfvars` from the example:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your environment-specific values. Each environment has different default thresholds:

**Development:**
```hcl
application_name             = "my-web-app"
application_tag_filter       = "service.name@dest EQUALS 'web-service'"
alert_email_addresses        = ["dev-team@example.com"]
latency_threshold_warning    = 2000  # More lenient
latency_threshold_critical   = 5000
error_rate_threshold         = 0.10
```

**Staging:**
```hcl
application_name             = "my-web-app"
alert_email_addresses        = ["staging-team@example.com"]
create_slack_channel         = true
slack_webhook_url            = "https://hooks.slack.com/..."
latency_threshold_warning    = 1500  # Production-like
latency_threshold_critical   = 3500
error_rate_threshold         = 0.07
```

**Production:**
```hcl
application_name             = "my-web-app"
alert_email_addresses        = ["ops@example.com", "oncall@example.com"]
create_slack_channel         = true
slack_webhook_url            = "https://hooks.slack.com/..."
latency_threshold_warning    = 1000  # Strict
latency_threshold_critical   = 3000
error_rate_threshold         = 0.05
```

### 6. Deploy

```bash
# Review the execution plan
terraform plan

# Apply the configuration
terraform apply

# View outputs
terraform output
```

### 7. Deploy to Other Environments

Repeat steps 1-6 for each environment you want to deploy:

```bash
# Deploy to staging
cd ../staging
terraform init -backend-config=backend.hcl
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars
terraform apply

# Deploy to production
cd ../production
terraform init -backend-config=backend.hcl
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars
terraform apply
```

## Outputs

After deployment, view the created resources from within your environment directory:

```bash
# Make sure you're in the environment directory
cd environments/dev  # or staging, or production

# All outputs
terraform output

# Specific output
terraform output application_id

# JSON format
terraform output -json
```

**Available Outputs:** (Add outputs.tf in each environment as needed)
- `application_id` - Application configuration ID
- `application_label` - Application name
- `email_channel_id` - Email alert channel ID
- `slack_channel_id` - Slack alert channel ID
- `alert_ids` - Map of all alert configuration IDs
- `rbac_role_ids` - Map of all RBAC role IDs



## Troubleshooting

### Module Not Found

```bash
# Re-initialize to download modules
terraform init -upgrade
```

### Authentication Issues

```bash
# Verify environment variables
echo $INSTANA_API_TOKEN
echo $INSTANA_ENDPOINT

# Test API connectivity
curl -H "Authorization: apiToken $INSTANA_API_TOKEN" \
  https://$INSTANA_ENDPOINT/api/application-monitoring/settings/application
```

### State Lock Issues

```bash
# Force unlock (use with caution)
terraform force-unlock <lock-id>
```

### Tag Filter Not Matching

1. Test filter in Instana UI first
2. Check entity origin (`@src`, `@dest`, `@na`)
3. Verify tag names and values
4. Use `terraform plan` to preview changes

## Cleanup

To destroy resources in a specific environment:

```bash
# Navigate to the environment
cd environments/dev  # or staging, or production

# Review what will be destroyed
terraform plan -destroy

# Destroy all resources in this environment
terraform destroy

# Destroy specific module in this environment
terraform destroy -target=module.rbac
```

**Important:** Each environment is isolated, so destroying one environment does not affect others.

## Best Practices

1. **Use Modules**: Keep resources organized and reusable
2. **Version Control**: Commit `.tf` files, exclude `.tfvars` and state files
3. **Remote State**: Always use S3 backend for team environments
4. **State Locking**: Enable DynamoDB locking to prevent conflicts
5. **Sensitive Data**: Use environment variables for secrets
6. **Tag Filters**: Test in Instana UI before applying
7. **Incremental Changes**: Enable features gradually
8. **Documentation**: Document custom configurations

## Multi-Environment Architecture

This example uses a **directory-based approach** for complete environment isolation. Each environment (dev, staging, production) has its own:

- Configuration files (`main.tf`, `variables.tf`)
- State file (via separate S3 keys)
- Variable values (`terraform.tfvars`)
- Backend configuration (`backend.hcl`)


### Shared Modules

All environments use the same shared modules from `modules/`:
- `application/` - Application monitoring configuration
- `alerting/` - Alert channels and configurations
- `rbac/` - Role-based access control
- Plus 7 additional modules (automation, dashboard, infrastructure, mobile, slo, synthetic, website)

This ensures consistency while allowing environment-specific customization through variables.

### State Management

Each environment maintains separate state files:

```
S3 Bucket: your-terraform-state-bucket
├── instana/dev/terraform.tfstate
├── instana/staging/terraform.tfstate
└── instana/production/terraform.tfstate
```

**Benefits:**
- Complete isolation between environments
- Independent deployment cycles
- No risk of cross-environment changes
- Separate access controls per environment

### Deployment Workflow

Typical workflow for promoting changes across environments:

```bash
# 1. Test in development
cd environments/dev
terraform apply

# 2. Validate in staging
cd ../staging
terraform apply

# 3. Deploy to production (after approval)
cd ../production
terraform apply
```
