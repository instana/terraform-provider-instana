# Application Module
# Creates an Instana application configuration with monitoring perspective

terraform {
  required_providers {
    instana = {
      source  = "instana/instana"
      version = ">= 7.0.0"
    }
  }
}

# Application Configuration
resource "instana_application_config" "this" {
  label          = var.application_name
  scope          = var.scope
  boundary_scope = var.boundary_scope
  tag_filter     = var.tag_filter
  
  # Access rules define who can access this application perspective
  access_rules = var.access_rules
}

# Note: This is a template. Implement specific resources based on your requirements.
# Refer to the Instana provider documentation for resource-specific configurations.