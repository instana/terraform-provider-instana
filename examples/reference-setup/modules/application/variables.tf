# Application Module Variables

variable "application_name" {
  description = "Name of the application to monitor"
  type        = string
}

variable "scope" {
  description = "Scope of the application perspective"
  type        = string
  default     = "INCLUDE_NO_DOWNSTREAM"
  
  validation {
    condition     = contains(["INCLUDE_ALL_DOWNSTREAM", "INCLUDE_NO_DOWNSTREAM", "INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING"], var.scope)
    error_message = "Scope must be one of: INCLUDE_ALL_DOWNSTREAM, INCLUDE_NO_DOWNSTREAM, INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING"
  }
}

variable "boundary_scope" {
  description = "Boundary scope of the application perspective"
  type        = string
  default     = "INBOUND"
  
  validation {
    condition     = contains(["INBOUND", "ALL", "DEFAULT"], var.boundary_scope)
    error_message = "Boundary scope must be one of: INBOUND, ALL, DEFAULT"
  }
}

variable "tag_filter" {
  description = "Tag filter expression to define which entities are included in the application"
  type        = string
}

variable "access_rules" {
  description = "List of access rules defining who can access this application perspective"
  type = list(object({
    access_type   = string
    relation_type = string
    related_id    = optional(string)
  }))
  default = [{
    access_type   = "READ_WRITE"
    relation_type = "GLOBAL"
  }]
}