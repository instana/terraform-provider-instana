# Application Module Outputs

output "id" {
  description = "ID of the created application configuration"
  value       = instana_application_config.this.id
}

output "label" {
  description = "Name/label of the application"
  value       = instana_application_config.this.label
}

output "scope" {
  description = "Scope of the application perspective"
  value       = instana_application_config.this.scope
}

output "boundary_scope" {
  description = "Boundary scope of the application"
  value       = instana_application_config.this.boundary_scope
}

output "tag_filter" {
  description = "Tag filter expression used"
  value       = instana_application_config.this.tag_filter
}