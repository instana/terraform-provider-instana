# RBAC Module Outputs

output "developer_role_id" {
  description = "ID of the developer role"
  value       = var.create_developer_role ? instana_rbac_role.developer[0].id : null
}

output "viewer_role_id" {
  description = "ID of the viewer role"
  value       = var.create_viewer_role ? instana_rbac_role.viewer[0].id : null
}

output "admin_role_id" {
  description = "ID of the admin role"
  value       = var.create_admin_role ? instana_rbac_role.admin[0].id : null
}

output "roles_summary" {
  description = "Summary of created RBAC roles"
  value = {
    developer_created = var.create_developer_role
    viewer_created    = var.create_viewer_role
    admin_created     = var.create_admin_role
    total_roles       = (var.create_developer_role ? 1 : 0) + (var.create_viewer_role ? 1 : 0) + (var.create_admin_role ? 1 : 0)
  }
}