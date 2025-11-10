package instana

// Data source name constants
const (
	// DataSourceInstanaCustomEventSpecificationFramework the name of the terraform-provider-instana data source to read custom event specifications
	DataSourceInstanaCustomEventSpecificationFramework = "custom_event_spec"
)

// Field name constants for custom event specification datasource (field constants are shared with resource-custom-event-specification-framework.go)
const (
	// CustomEventSpecificationFieldID constant value for the schema field id
	CustomEventSpecificationFieldID = "id"
)

// Description constants for custom event specification fields
const (
	// CustomEventSpecificationDescDataSource description for the data source
	CustomEventSpecificationDescDataSource = "Data source for an Instana custom event specification. Custom events are user-defined events in Instana."
	// CustomEventSpecificationDescID description for the ID field
	CustomEventSpecificationDescID = "The ID of the custom event specification."
	// CustomEventSpecificationDescName description for the name field
	CustomEventSpecificationDescName = "The name of the custom event specification."
	// CustomEventSpecificationDescDescription description for the description field
	CustomEventSpecificationDescDescription = "The description of the custom event specification."
	// CustomEventSpecificationDescEntityType description for the entity_type field
	CustomEventSpecificationDescEntityType = "The entity type for which the custom event specification is created."
	// CustomEventSpecificationDescTriggering description for the triggering field
	CustomEventSpecificationDescTriggering = "Indicates if an incident is triggered the custom event or not."
	// CustomEventSpecificationDescEnabled description for the enabled field
	CustomEventSpecificationDescEnabled = "Indicates if the custom event is enabled or not."
	// CustomEventSpecificationDescQuery description for the query field
	CustomEventSpecificationDescQuery = "Dynamic focus query for the custom event specification."
	// CustomEventSpecificationDescExpirationTime description for the expiration_time field
	CustomEventSpecificationDescExpirationTime = "The expiration time (grace period) to wait before the issue is closed."
)

// Error message constants
const (
	// CustomEventSpecificationErrUnexpectedConfigureType error message for unexpected configure type
	CustomEventSpecificationErrUnexpectedConfigureType = "Unexpected Data Source Configure Type"
	// CustomEventSpecificationErrUnexpectedConfigureTypeDetail error message detail for unexpected configure type
	CustomEventSpecificationErrUnexpectedConfigureTypeDetail = "Expected *ProviderMeta, got: %T. Please report this issue to the provider developers."
	// CustomEventSpecificationErrReadingSpecs error message for reading custom event specifications
	CustomEventSpecificationErrReadingSpecs = "Error reading custom event specifications"
	// CustomEventSpecificationErrReadingSpecsDetail error message detail for reading custom event specifications
	CustomEventSpecificationErrReadingSpecsDetail = "Could not read custom event specifications: %s"
	// CustomEventSpecificationErrNotFound error message for custom event specification not found
	CustomEventSpecificationErrNotFound = "Custom event specification not found"
	// CustomEventSpecificationErrNotFoundDetail error message detail for custom event specification not found
	CustomEventSpecificationErrNotFoundDetail = "No custom event specification found for name '%s' and entity type '%s'"
)

// Made with Bob
