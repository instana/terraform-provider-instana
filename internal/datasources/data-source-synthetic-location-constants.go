package datasources

// Data source name constants
const (
	// DataSourceInstanaSyntheticLocationFramework the name of the terraform-provider-instana data source to read synthetic locations
	DataSourceInstanaSyntheticLocationFramework = "synthetic_location"
	// DataSourceSyntheticLocation the name of the terraform-provider-instana data source for synthetic location specifications
	DataSourceSyntheticLocation = "instana_synthetic_location"
)

// Field name constants for synthetic location
const (
	// SyntheticLocationFieldLabel constant value for the schema field label
	SyntheticLocationFieldLabel = "label"
	// SyntheticLocationFieldDescription constant value for the computed schema field description
	SyntheticLocationFieldDescription = "description"
	// SyntheticLocationFieldLocationType constant value for the schema field location_type
	SyntheticLocationFieldLocationType = "location_type"
)

// Field ID constant
const (
	// SyntheticLocationFieldID constant value for the schema field id
	SyntheticLocationFieldID = "id"
)

// Description constants for synthetic location fields
const (
	// SyntheticLocationDescDataSource description for the data source
	SyntheticLocationDescDataSource = "Data source for Instana synthetic locations. Synthetic locations are the locations from which synthetic tests are executed."
	// SyntheticLocationDescID description for the ID field
	SyntheticLocationDescID = "The ID of the synthetic location."
	// SyntheticLocationDescLabel description for the label field
	SyntheticLocationDescLabel = "Friendly name of the Synthetic Location"
	// SyntheticLocationDescDescription description for the description field
	SyntheticLocationDescDescription = "The description of the Synthetic location"
	// SyntheticLocationDescLocationType description for the location_type field
	SyntheticLocationDescLocationType = "Indicates if the location is public or private"
)

// Error message constants
const (
	// SyntheticLocationErrUnexpectedConfigureType error message for unexpected configure type
	SyntheticLocationErrUnexpectedConfigureType = "Unexpected Data Source Configure Type"
	// SyntheticLocationErrUnexpectedConfigureTypeDetail error message detail for unexpected configure type
	SyntheticLocationErrUnexpectedConfigureTypeDetail = "Expected *ProviderMeta, got: %T. Please report this issue to the provider developers."
	// SyntheticLocationErrReadingLocations error message for reading synthetic locations
	SyntheticLocationErrReadingLocations = "Error reading synthetic locations"
	// SyntheticLocationErrReadingLocationsDetail error message detail for reading synthetic locations
	SyntheticLocationErrReadingLocationsDetail = "Could not read synthetic locations: %s"
	// SyntheticLocationErrNotFound error message for synthetic location not found
	SyntheticLocationErrNotFound = "No matching synthetic location found"
	// SyntheticLocationErrNotFoundDetail error message detail for synthetic location not found
	SyntheticLocationErrNotFoundDetail = "No synthetic location found with label '%s' and location type '%s'"
)

// Made with Bob
