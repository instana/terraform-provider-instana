# Staging Environment Variables

# Instana Provider Configuration
variable "instana_api_token" {
  description = "Instana API token for authentication. Can be set via INSTANA_API_TOKEN environment variable"
  type        = string
  sensitive   = true
  default     = null
}

variable "instana_endpoint" {
  description = "Instana API endpoint (e.g., your-tenant.instana.io). Can be set via INSTANA_ENDPOINT environment variable"
  type        = string
  default     = null
}

# Application Configuration
variable "application_name" {
  description = "Name of the application to monitor"
  type        = string
  default     = "my-application"
}

variable "application_scope" {
  description = "Scope of the application perspective"
  type        = string
  default     = "INCLUDE_NO_DOWNSTREAM"
}

variable "application_boundary_scope" {
  description = "Boundary scope of the application perspective"
  type        = string
  default     = "INBOUND"
}

variable "application_tag_filter" {
  description = "Tag filter expression to define which entities are included in the application"
  type        = string
  default     = "service.name@dest EQUALS 'my-service'"
}

# Alerting Configuration
variable "enable_alerts" {
  description = "Enable or disable alert configurations"
  type        = bool
  default     = true
}

variable "alerts_triggering" {
  description = "Whether alerts should trigger incidents"
  type        = bool
  default     = true  # Production-like for staging
}

variable "alert_evaluation_type" {
  description = "Evaluation type for alerts"
  type        = string
  default     = "PER_AP"
}

variable "alert_granularity" {
  description = "Evaluation granularity in milliseconds"
  type        = number
  default     = 600000
}

# Email Channel Configuration
variable "create_email_channel" {
  description = "Create email alerting channel"
  type        = bool
  default     = true
}

variable "alert_email_addresses" {
  description = "List of email addresses to receive alerts"
  type        = list(string)
  default     = ["staging-team@example.com"]
}

# Slack Channel Configuration
variable "create_slack_channel" {
  description = "Create Slack alerting channel"
  type        = bool
  default     = true
}

variable "slack_webhook_url" {
  description = "Slack webhook URL for notifications"
  type        = string
  sensitive   = true
  default     = null
}

variable "slack_channel" {
  description = "Slack channel name for notifications"
  type        = string
  default     = "#staging-alerts"
}

# Slowness Alert Configuration (Production-like thresholds)
variable "create_slowness_alert" {
  description = "Create slowness alert"
  type        = bool
  default     = true
}

variable "latency_threshold_warning" {
  description = "Warning threshold for application latency in milliseconds"
  type        = number
  default     = 1500  # Production-like
}

variable "latency_threshold_critical" {
  description = "Critical threshold for application latency in milliseconds"
  type        = number
  default     = 3500  # Production-like
}

# Error Alert Configuration
variable "create_error_alert" {
  description = "Create error rate alert"
  type        = bool
  default     = true
}

variable "error_rate_threshold" {
  description = "Threshold for application error rate (0.0 to 1.0)"
  type        = number
  default     = 0.07  # Production-like
}

# Throughput Alert Configuration
variable "create_throughput_alert" {
  description = "Create throughput alert"
  type        = bool
  default     = true
}

variable "throughput_threshold" {
  description = "Threshold for application throughput (calls per time window)"
  type        = number
  default     = 100
}

# Log Alert Configuration
variable "create_log_alert" {
  description = "Create log alert"
  type        = bool
  default     = true
}

variable "log_query" {
  description = "Log query filter expression"
  type        = string
  default     = "log.level@na EQUALS 'ERROR'"
}

variable "log_threshold_warning" {
  description = "Warning threshold for log occurrences"
  type        = number
  default     = 15  # Production-like
}

variable "log_threshold_critical" {
  description = "Critical threshold for log occurrences"
  type        = number
  default     = 75  # Production-like
}

# RBAC Configuration
variable "create_developer_role" {
  description = "Create developer RBAC role"
  type        = bool
  default     = true
}

variable "developer_user_ids" {
  description = "List of user IDs to assign to developer role"
  type        = list(string)
  default     = []
}

variable "create_viewer_role" {
  description = "Create viewer RBAC role"
  type        = bool
  default     = true
}

variable "viewer_user_ids" {
  description = "List of user IDs to assign to viewer role"
  type        = list(string)
  default     = []
}

variable "create_admin_role" {
  description = "Create admin RBAC role"
  type        = bool
  default     = true
}

variable "admin_user_ids" {
  description = "List of user IDs to assign to admin role"
  type        = list(string)
  default     = []
}