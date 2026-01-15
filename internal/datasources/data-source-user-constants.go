package datasources

// Data source specific constants for user
const (
	// UserDataSourceFieldID constant for the ID field
	UserDataSourceFieldID = "id"

	// Description constants for data source
	UserDescDataSource = "Data source for an Instana user. Users can be referenced in roles and teams."
	UserDescID         = "The unique identifier of the user."
	UserDescEmail      = "The email address of the user."
	UserDescFullName   = "The full name of the user."

	// Error message constants
	UserErrUnexpectedConfigureType = "Unexpected Data Source Configure Type"
	UserErrReadingUsers            = "Error reading users"
	UserErrUserNotFound            = "User not found"

	// UserFieldEmail constant value for the schema field email
	UserFieldEmail = "email"
	// UserFieldFullName constant value for the schema field full_name
	UserFieldFullName = "full_name"
)

// Made with Bob
