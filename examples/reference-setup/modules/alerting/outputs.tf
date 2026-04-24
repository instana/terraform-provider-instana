# Alerting Module Outputs

# Alert Channel Outputs
output "email_channel_id" {
  description = "ID of the email alerting channel"
  value       = var.create_email_channel ? instana_alerting_channel.email[0].id : null
}

output "slack_channel_id" {
  description = "ID of the Slack alerting channel"
  value       = var.create_slack_channel ? instana_alerting_channel.slack[0].id : null
}

output "all_channel_ids" {
  description = "List of all created alert channel IDs"
  value       = local.all_channel_ids
}

# Alert Configuration Outputs
output "slowness_alert_id" {
  description = "ID of the slowness alert configuration"
  value       = var.create_slowness_alert ? instana_application_alert_config.slowness[0].id : null
}

output "error_alert_id" {
  description = "ID of the error rate alert configuration"
  value       = var.create_error_alert ? instana_application_alert_config.errors[0].id : null
}

output "throughput_alert_id" {
  description = "ID of the throughput alert configuration"
  value       = var.create_throughput_alert ? instana_application_alert_config.throughput[0].id : null
}

output "log_alert_id" {
  description = "ID of the log alert configuration"
  value       = var.create_log_alert ? instana_log_alert_config.logs[0].id : null
}

# Summary Output
output "alert_summary" {
  description = "Summary of created alerting resources"
  value = {
    channels_created = {
      email = var.create_email_channel
      slack = var.create_slack_channel
    }
    alerts_created = {
      slowness   = var.create_slowness_alert
      errors     = var.create_error_alert
      throughput = var.create_throughput_alert
      logs       = var.create_log_alert
    }
    total_channels = length(local.all_channel_ids)
  }
}