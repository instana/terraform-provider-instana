# Alerting Module Variables


# Application Reference
variable "application_id" {
  description = "ID of the application to monitor"
  type        = string
}

# General Alert Configuration
variable "alerts_enabled" {
  description = "Enable or disable all alerts"
  type        = bool
  default     = true
}

variable "alerts_triggering" {
  description = "Whether alerts should trigger incidents"
  type        = bool
  default     = true
}

variable "boundary_scope" {
  description = "Boundary scope for alerts"
  type        = string
  default     = "INBOUND"

  validation {
    condition     = contains(["INBOUND", "ALL", "DEFAULT"], var.boundary_scope)
    error_message = "Boundary scope must be one of: INBOUND, ALL, DEFAULT"
  }
}

variable "evaluation_type" {
  description = "Evaluation type for alerts"
  type        = string
  default     = "PER_AP"

  validation {
    condition     = contains(["PER_AP", "PER_AP_SERVICE", "PER_AP_ENDPOINT"], var.evaluation_type)
    error_message = "Evaluation type must be one of: PER_AP, PER_AP_SERVICE, PER_AP_ENDPOINT"
  }
}

variable "granularity" {
  description = "Evaluation granularity in milliseconds"
  type        = number
  default     = 600000

  validation {
    condition     = contains([300000, 600000, 900000, 1200000, 1800000], var.granularity)
    error_message = "Granularity must be one of: 300000, 600000, 900000, 1200000, 1800000"
  }
}

variable "time_window" {
  description = "Time window for violations in milliseconds"
  type        = number
  default     = 600000
}

# Email Channel Configuration
variable "create_email_channel" {
  description = "Create email alerting channel"
  type        = bool
  default     = true
}

variable "email_channel_name" {
  description = "Name of the email alerting channel"
  type        = string
  default     = "Email Alerts"
}

variable "email_addresses" {
  description = "List of email addresses for alerts"
  type        = list(string)
  default     = []
}

# Slack Channel Configuration
variable "create_slack_channel" {
  description = "Create Slack alerting channel"
  type        = bool
  default     = true
}

variable "slack_channel_name" {
  description = "Name of the Slack alerting channel"
  type        = string
  default     = "Slack Alerts"
}

variable "slack_webhook_url" {
  description = "Slack webhook URL"
  type        = string
  default     = null
  sensitive   = true
}

variable "slack_channel" {
  description = "Slack channel name"
  type        = string
  default     = "#alerts"
}

# Additional Channels
variable "additional_channel_ids" {
  description = "Additional alert channel IDs to include"
  type        = list(string)
  default     = []
}

# Slowness Alert Configuration
variable "create_slowness_alert" {
  description = "Create slowness alert"
  type        = bool
  default     = true
}

variable "slowness_alert_name" {
  description = "Name of the slowness alert"
  type        = string
  default     = "Application Slowness Alert"
}

variable "slowness_alert_description" {
  description = "Description of the slowness alert"
  type        = string
  default     = "Alert when application response time exceeds thresholds"
}

variable "slowness_aggregation" {
  description = "Aggregation method for slowness metric"
  type        = string
  default     = "P90"

  validation {
    condition     = contains(["SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"], var.slowness_aggregation)
    error_message = "Aggregation must be one of: SUM, MEAN, MAX, MIN, P25, P50, P75, P90, P95, P98, P99"
  }
}

variable "slowness_threshold_warning" {
  description = "Warning threshold for slowness in milliseconds"
  type        = number
  default     = 1000
}

variable "slowness_threshold_critical" {
  description = "Critical threshold for slowness in milliseconds"
  type        = number
  default     = 3000
}

# Error Alert Configuration
variable "create_error_alert" {
  description = "Create error rate alert"
  type        = bool
  default     = true
}

variable "error_alert_name" {
  description = "Name of the error alert"
  type        = string
  default     = "Application Error Rate Alert"
}

variable "error_alert_description" {
  description = "Description of the error alert"
  type        = string
  default     = "Alert when application error rate exceeds thresholds"
}

variable "error_aggregation" {
  description = "Aggregation method for error metric"
  type        = string
  default     = "MEAN"

  validation {
    condition     = contains(["SUM", "MEAN", "MAX", "MIN"], var.error_aggregation)
    error_message = "Aggregation must be one of: SUM, MEAN, MAX, MIN"
  }
}

variable "error_threshold_warning" {
  description = "Warning threshold for error rate (0.0 to 1.0)"
  type        = number
  default     = 0.05
}

variable "error_threshold_critical" {
  description = "Critical threshold for error rate (0.0 to 1.0)"
  type        = number
  default     = null
}

# Throughput Alert Configuration
variable "create_throughput_alert" {
  description = "Create throughput alert"
  type        = bool
  default     = true
}

variable "throughput_alert_name" {
  description = "Name of the throughput alert"
  type        = string
  default     = "Application Throughput Alert"
}

variable "throughput_alert_description" {
  description = "Description of the throughput alert"
  type        = string
  default     = "Alert when application throughput drops below threshold"
}

variable "throughput_threshold_warning" {
  description = "Warning threshold for throughput (calls per time window)"
  type        = number
  default     = 100
}

# Log Alert Configuration
variable "create_log_alert" {
  description = "Create log alert"
  type        = bool
  default     = true
}

variable "log_alert_name" {
  description = "Name of the log alert"
  type        = string
  default     = "Application Log Alert"
}

variable "log_alert_description" {
  description = "Description of the log alert"
  type        = string
  default     = "Alert when log patterns match specified criteria"
}

variable "log_query" {
  description = "Log query to match (e.g., 'log.level@na EQUALS 'ERROR'')"
  type        = string
  default     = "log.level@na EQUALS 'ERROR'"
}

variable "log_threshold_warning" {
  description = "Warning threshold for log occurrences"
  type        = number
  default     = 10
}

variable "log_threshold_critical" {
  description = "Critical threshold for log occurrences"
  type        = number
  default     = 50
}