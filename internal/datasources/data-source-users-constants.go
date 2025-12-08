package datasources

// Data source specific constants for users (plural)
const (
	// UsersDataSourceFieldID constant for the ID field
	UsersDataSourceFieldID = "id"

	// Description constants for data source
	UsersDescDataSource   = "Data source for multiple Instana users. Allows querying multiple users by their email addresses."
	UsersDescID           = "The unique identifier for this datasource query."
	UsersDescEmails       = "List of email addresses to query. Returns users matching these emails."
	UsersDescUsers        = "List of users matching the provided email addresses."
	UsersDescUserID       = "The unique identifier of the user."
	UsersDescUserEmail    = "The email address of the user."
	UsersDescUserFullName = "The full name of the user."

	// Error message constants
	UsersErrUnexpectedConfigureType = "Unexpected Data Source Configure Type"
	UsersErrReadingUsers            = "Error reading users"

	// UsersFieldEmails constant value for the schema field emails
	UsersFieldEmails = "emails"
	// UsersFieldUsers constant value for the schema field users
	UsersFieldUsers = "users"
	// UsersFieldUserID constant value for the nested user id field
	UsersFieldUserID = "id"
	// UsersFieldUserEmail constant value for the nested user email field
	UsersFieldUserEmail = "email"
	// UsersFieldUserFullName constant value for the nested user full_name field
	UsersFieldUserFullName = "full_name"
)

// Made with Bob
