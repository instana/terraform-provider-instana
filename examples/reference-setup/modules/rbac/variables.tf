# RBAC Module Variables

# Developer Role Configuration
variable "create_developer_role" {
  description = "Create developer role"
  type        = bool
  default     = true
}

variable "developer_role_name" {
  description = "Name of the developer role"
  type        = string
  default     = "Developer"
}

variable "developer_permissions" {
  description = "List of permissions for developers"
  type        = list(string)
  default = [
    "CAN_CONFIGURE_APPLICATIONS",
    "CAN_VIEW_TRACE_DETAILS",
    "CAN_VIEW_LOGS",
    "CAN_CONFIGURE_APPLICATION_SMART_ALERTS"
  ]
}

variable "developer_user_ids" {
  description = "List of user IDs to assign to developer role"
  type        = list(string)
  default     = []
}

# Viewer Role Configuration
variable "create_viewer_role" {
  description = "Create viewer role"
  type        = bool
  default     = true
}

variable "viewer_role_name" {
  description = "Name of the viewer role"
  type        = string
  default     = "Viewer"
}

variable "viewer_permissions" {
  description = "List of permissions for viewers"
  type        = list(string)
  default = [
    "CAN_VIEW_TRACE_DETAILS",
    "CAN_VIEW_LOGS",
    "CAN_VIEW_AUDIT_LOG"
  ]
}

variable "viewer_user_ids" {
  description = "List of user IDs to assign to viewer role"
  type        = list(string)
  default     = []
}

# Admin Role Configuration
variable "create_admin_role" {
  description = "Create admin role"
  type        = bool
  default     = true
}

variable "admin_role_name" {
  description = "Name of the admin role"
  type        = string
  default     = "Administrator"
}

variable "admin_permissions" {
  description = "List of permissions for administrators"
  type        = list(string)
  default = [
    "CAN_CONFIGURE_APPLICATIONS",
    "CAN_CONFIGURE_EUM_APPLICATIONS",
    "CAN_CONFIGURE_AGENTS",
    "CAN_VIEW_TRACE_DETAILS",
    "CAN_VIEW_LOGS",
    "CAN_CONFIGURE_SESSION_SETTINGS",
    "CAN_CONFIGURE_INTEGRATIONS",
    "CAN_CONFIGURE_GLOBAL_APPLICATION_SMART_ALERTS",
    "CAN_CONFIGURE_API_TOKENS",
    "CAN_CONFIGURE_SERVICE_LEVELS",
    "CAN_VIEW_AUDIT_LOG",
    "CAN_CONFIGURE_EVENTS_AND_ALERTS",
    "CAN_CONFIGURE_MAINTENANCE_WINDOWS",
    "CAN_CONFIGURE_APPLICATION_SMART_ALERTS",
    "CAN_CONFIGURE_USERS",
    "CAN_CONFIGURE_TEAMS"
  ]
}

variable "admin_user_ids" {
  description = "List of user IDs to assign to admin role"
  type        = list(string)
  default     = []
}