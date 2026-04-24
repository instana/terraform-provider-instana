# Alerting Module
# Creates alerting channels and various types of alert configurations

terraform {
  required_providers {
    instana = {
      source  = "instana/instana"
      version = ">= 7.0.0"
    }
  }
}

# Email Alerting Channel
resource "instana_alerting_channel" "email" {
  count = var.create_email_channel ? 1 : 0

  name = var.email_channel_name

  email = {
    emails = var.email_addresses
  }
}

# Slack Alerting Channel
resource "instana_alerting_channel" "slack" {
  count = var.create_slack_channel ? 1 : 0

  name = var.slack_channel_name

  slack = {
    webhook_url = var.slack_webhook_url
    channel     = var.slack_channel
  }
}

# Application Alert - Slowness
resource "instana_application_alert_config" "slowness" {
  count = var.create_slowness_alert ? 1 : 0

  name            = var.slowness_alert_name
  description     = var.slowness_alert_description
  enabled         = var.alerts_enabled
  triggering      = var.alerts_triggering
  boundary_scope  = var.boundary_scope
  evaluation_type = var.evaluation_type
  granularity     = var.granularity

  application = [{
    application_id = var.application_id
    inclusive      = true
  }]

  rules = [{
    rule = {
      slowness = {
        metric_name = "latency"
        aggregation = var.slowness_aggregation
      }
    }
    threshold = {
      warning = var.slowness_threshold_warning != null ? {
        static = {
          value = var.slowness_threshold_warning
        }
      } : null
      critical = var.slowness_threshold_critical != null ? {
        static = {
          value = var.slowness_threshold_critical
        }
      } : null
    }
    threshold_operator = ">="
  }]

  time_threshold = {
    violations_in_sequence = {
      time_window = var.time_window
    }
  }

  alert_channels = local.alert_channels
}

# Application Alert - Error Rate
resource "instana_application_alert_config" "errors" {
  count = var.create_error_alert ? 1 : 0

  name            = var.error_alert_name
  description     = var.error_alert_description
  enabled         = var.alerts_enabled
  triggering      = var.alerts_triggering
  boundary_scope  = var.boundary_scope
  evaluation_type = var.evaluation_type
  granularity     = var.granularity

  application = [{
    application_id = var.application_id
    inclusive      = true
  }]

  rules = [{
    rule = {
      errors = {
        metric_name = "errors"
        aggregation = var.error_aggregation
      }
    }
    threshold = {
      warning = var.error_threshold_warning != null ? {
        static = {
          value = var.error_threshold_warning
        }
      } : null
      critical = var.error_threshold_critical != null ? {
        static = {
          value = var.error_threshold_critical
        }
      } : null
    }
    threshold_operator = ">="
  }]

  time_threshold = {
    violations_in_sequence = {
      time_window = var.time_window
    }
  }

  alert_channels = local.alert_channels
}

# Application Alert - Throughput
resource "instana_application_alert_config" "throughput" {
  count = var.create_throughput_alert ? 1 : 0

  name            = var.throughput_alert_name
  description     = var.throughput_alert_description
  enabled         = var.alerts_enabled
  triggering      = var.alerts_triggering
  boundary_scope  = var.boundary_scope
  evaluation_type = var.evaluation_type
  granularity     = var.granularity

  application = [{
    application_id = var.application_id
    inclusive      = true
  }]

  rules = [{
    rule = {
      throughput = {
        metric_name = "calls"
        aggregation = "SUM"
      }
    }
    threshold = {
      warning = var.throughput_threshold_warning != null ? {
        static = {
          value = var.throughput_threshold_warning
        }
      } : null
    }
    threshold_operator = "<="
  }]

  time_threshold = {
    violations_in_sequence = {
      time_window = var.time_window
    }
  }

  alert_channels = local.alert_channels
}

# Log Alert Configuration
resource "instana_log_alert_config" "logs" {
  count = var.create_log_alert ? 1 : 0

  name        = var.log_alert_name
  description = var.log_alert_description
  granularity = var.granularity

  tag_filter = var.log_query

  rules = {
    metric_name = "logCount"
    alert_type  = "log.count"
    aggregation = "SUM"
    threshold = {
      warning = var.log_threshold_warning != null ? {
        static = {
          value = var.log_threshold_warning
        }
      } : null
      critical = var.log_threshold_critical != null ? {
        static = {
          value = var.log_threshold_critical
        }
      } : null
    }
    threshold_operator = ">="
  }

  time_threshold = {
    violations_in_sequence = {
      time_window = var.time_window
    }
  }

  alert_channels = local.alert_channels
}

# Local values for alert channel routing
locals {
  email_channel_ids = var.create_email_channel ? [instana_alerting_channel.email[0].id] : []
  slack_channel_ids = var.create_slack_channel ? [instana_alerting_channel.slack[0].id] : []

  all_channel_ids = concat(
    local.email_channel_ids,
    local.slack_channel_ids,
    var.additional_channel_ids
  )

  alert_channels = length(local.all_channel_ids) > 0 ? {
    WARNING  = local.all_channel_ids
    CRITICAL = local.all_channel_ids
  } : {}
}

# Note: This is a template. Implement specific resources based on your requirements.
# Refer to the Instana provider documentation for resource-specific configurations.