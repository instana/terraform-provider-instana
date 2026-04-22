# RBAC Module
# Creates RBAC roles for access control

terraform {
  required_providers {
    instana = {
      source  = "instana/instana"
      version = ">= 7.0.0"
    }
  }
}

# Developer Role
resource "instana_rbac_role" "developer" {
  count = var.create_developer_role ? 1 : 0
  
  name = var.developer_role_name
  
  permissions = var.developer_permissions
  
  member = [
    for user_id in var.developer_user_ids : {
      user_id = user_id
    }
  ]
}

# Viewer Role
resource "instana_rbac_role" "viewer" {
  count = var.create_viewer_role ? 1 : 0
  
  name = var.viewer_role_name
  
  permissions = var.viewer_permissions
  
  member = [
    for user_id in var.viewer_user_ids : {
      user_id = user_id
    }
  ]
}

# Admin Role
resource "instana_rbac_role" "admin" {
  count = var.create_admin_role ? 1 : 0
  
  name = var.admin_role_name
  
  permissions = var.admin_permissions
  
  member = [
    for user_id in var.admin_user_ids : {
      user_id = user_id
    }
  ]
}

# Note: This is a template. Implement specific resources based on your requirements.
# Refer to the Instana provider documentation for resource-specific configurations.