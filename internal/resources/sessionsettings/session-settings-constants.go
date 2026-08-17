package sessionsettings

// ResourceInstanaSessionSettings is the name of the terraform resource for session settings.
const ResourceInstanaSessionSettings = "session_settings"

// Schema field name constants
const (
	// SessionSettingsFieldTokenLifeTimeInMillis constant for the token lifetime field
	SessionSettingsFieldTokenLifeTimeInMillis = "token_life_time_in_millis"
	// SessionSettingsFieldIdleTimeInMillis constant for the idle time field
	SessionSettingsFieldIdleTimeInMillis = "idle_time_in_millis"
)

// Resource description constants
const (
	// SessionSettingsDescResource describes the resource purpose
	SessionSettingsDescResource = "Manages the tenant unit session settings in Instana, " +
		"including token lifetime and idle timeout configuration. " +
		"This is a singleton resource — only one instance exists per tenant unit."
	// SessionSettingsDescTokenLifeTime describes the token lifetime field
	SessionSettingsDescTokenLifeTime = "Maximum lifetime of an authentication token in milliseconds. " +
		"Valid range: 600000 (10 min) to 604800000 (7 days)."
	// SessionSettingsDescIdleTime describes the idle time field
	SessionSettingsDescIdleTime = "Idle timeout before a session expires in milliseconds. " +
		"Valid range: 600000 (10 min) to 28800000 (8 hours)."
)

// Validation constants enforced by the Instana API
const (
	// SessionSettingsMinTokenLifeTimeInMillis minimum token lifetime (10 minutes)
	SessionSettingsMinTokenLifeTimeInMillis = 600000
	// SessionSettingsMaxTokenLifeTimeInMillis maximum token lifetime (7 days)
	SessionSettingsMaxTokenLifeTimeInMillis = 604800000
	// SessionSettingsDefaultTokenLifeTimeInMillis default token lifetime (7 days)
	SessionSettingsDefaultTokenLifeTimeInMillis = 604800000
	// SessionSettingsMinIdleTimeInMillis minimum idle time (10 minutes)
	SessionSettingsMinIdleTimeInMillis = 600000
	// SessionSettingsMaxIdleTimeInMillis maximum idle time (8 hours)
	SessionSettingsMaxIdleTimeInMillis = 28800000
	// SessionSettingsDefaultIdleTimeInMillis default idle time (8 hours)
	SessionSettingsDefaultIdleTimeInMillis = 28800000
)

