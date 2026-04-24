# Production Environment Outputs

output "environment" {
  description = "Environment name"
  value       = "production"
}

output "application_id" {
  description = "Application configuration ID"
  value       = module.application.id
}

output "application_name" {
  description = "Application name"
  value       = module.application.label
}

output "email_channel_id" {
  description = "Email alert channel ID"
  value       = module.alerting.email_channel_id
}

output "slack_channel_id" {
  description = "Slack alert channel ID (if created)"
  value       = module.alerting.slack_channel_id
}

output "alert_ids" {
  description = "Map of all alert configuration IDs"
  value = {
    slowness_alert   = module.alerting.slowness_alert_id
    error_alert      = module.alerting.error_alert_id
    throughput_alert = module.alerting.throughput_alert_id
    log_alert        = module.alerting.log_alert_id
  }
}

output "rbac_role_ids" {
  description = "Map of all RBAC role IDs"
  value = {
    developer_role = module.rbac.developer_role_id
    viewer_role    = module.rbac.viewer_role_id
    admin_role     = module.rbac.admin_role_id
  }
}

output "deployment_summary" {
  description = "Summary of deployed resources"
  value = {
    environment      = "production"
    application_name = module.application.label
    alert_channels   = [module.alerting.email_channel_id, module.alerting.slack_channel_id]
    alerts_enabled   = var.enable_alerts
  }
}