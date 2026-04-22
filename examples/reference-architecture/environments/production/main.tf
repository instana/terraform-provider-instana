# Production Environment Configuration
# This configuration uses shared modules from ../../modules/

terraform {
  required_providers {
    instana = {
      source  = "instana/instana"
      version = ">= 7.0.0"
    }
  }

  # S3 backend for state storage
  # Configure via: terraform init -backend-config=backend.hcl
  backend "s3" {}
}

# Provider configuration
provider "instana" {
  api_token = var.instana_api_token
  endpoint  = var.instana_endpoint
}

# Application Module
module "application" {
  source = "../../modules/application"
  
  application_name    = var.application_name
  scope               = var.application_scope
  boundary_scope      = var.application_boundary_scope
  tag_filter          = var.application_tag_filter
}

# Alerting Module
module "alerting" {
  source = "../../modules/alerting"
  
  # Link to application
  application_id = module.application.id
  
  # Email Channel Configuration
  create_email_channel = var.create_email_channel
  email_addresses      = var.alert_email_addresses
  
  # Slack Channel Configuration
  create_slack_channel = var.create_slack_channel
  slack_webhook_url    = var.slack_webhook_url
  slack_channel        = var.slack_channel
  
  # Alert Configuration
  alerts_enabled    = var.enable_alerts
  alerts_triggering = var.alerts_triggering
  boundary_scope    = var.application_boundary_scope
  evaluation_type   = var.alert_evaluation_type
  granularity       = var.alert_granularity
  
  # Slowness Alert (strict production thresholds)
  create_slowness_alert       = var.create_slowness_alert
  slowness_threshold_warning  = var.latency_threshold_warning
  slowness_threshold_critical = var.latency_threshold_critical
  
  # Error Alert
  create_error_alert      = var.create_error_alert
  error_threshold_warning = var.error_rate_threshold
  
  # Throughput Alert
  create_throughput_alert      = var.create_throughput_alert
  throughput_threshold_warning = var.throughput_threshold
  
  # Log Alert
  create_log_alert         = var.create_log_alert
  log_query                = var.log_query
  log_threshold_warning    = var.log_threshold_warning
  log_threshold_critical   = var.log_threshold_critical
}

# RBAC Module
module "rbac" {
  source = "../../modules/rbac"
  
  # Developer Role
  create_developer_role = var.create_developer_role
  developer_user_ids    = var.developer_user_ids
  
  # Viewer Role
  create_viewer_role = var.create_viewer_role
  viewer_user_ids    = var.viewer_user_ids
  
  # Admin Role
  create_admin_role = var.create_admin_role
  admin_user_ids    = var.admin_user_ids
}