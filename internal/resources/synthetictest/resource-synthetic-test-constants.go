package synthetictest

const (
	// Resource description constants
	SyntheticTestDescResource         = "This resource manages Synthetic Tests in Instana."
	SyntheticTestDescID               = "The ID of the Synthetic Test."
	SyntheticTestDescLabel            = "Friendly name of the Synthetic test."
	SyntheticTestDescDescription      = "The description of the Synthetic test."
	SyntheticTestDescActive           = "Indicates if the Synthetic test is started or not."
	SyntheticTestDescApplicationID    = "Unique identifier of the Application Perspective."
	SyntheticTestDescCustomProperties = "Name/value pairs to provide additional information of the Synthetic test."
	SyntheticTestDescLocations        = "Array of the PoP location IDs."
	SyntheticTestDescPlaybackMode     = "Defines how the Synthetic test should be executed across multiple PoPs."
	SyntheticTestDescTestFrequency    = "How often the playback for a Synthetic test is scheduled."

	// HTTP Action block descriptions
	SyntheticTestDescHttpAction        = "The configuration of the synthetic alert of type http action."
	SyntheticTestDescMarkSyntheticCall = "Flag used to control if HTTP calls will be marked as synthetic calls."
	SyntheticTestDescRetries           = "Indicates how many attempts will be allowed to get a successful connection."
	SyntheticTestDescRetryInterval     = "The time interval between retries in seconds."
	SyntheticTestDescTimeout           = "The timeout to be used by the PoP playback engines running the test."
	SyntheticTestDescURL               = "The URL which is being tested."
	SyntheticTestDescOperation         = "The HTTP operation."
	SyntheticTestDescHeaders           = "An object with header/value pairs."
	SyntheticTestDescBody              = "The body content to send with the operation."
	SyntheticTestDescValidationString  = "An expression to be evaluated."
	SyntheticTestDescFollowRedirect    = "A boolean type, true by default; to allow redirect."
	SyntheticTestDescAllowInsecure     = "A boolean type, if set to true then allow insecure certificates."
	SyntheticTestDescExpectStatus      = "An integer type, by default, the Synthetic passes for any 2XX status code."
	SyntheticTestDescExpectMatch       = "An optional regular expression string to be used to check the test response."

	// HTTP Script block descriptions
	SyntheticTestDescHttpScript = "The configuration of the synthetic alert of type http script."
	SyntheticTestDescScript     = "The Javascript content in plain text."

	// Error message constants
	SyntheticTestErrConfigRequired       = "Configuration required"
	SyntheticTestErrConfigRequiredMsg    = "Either http_action or http_script configuration must be provided"
	SyntheticTestErrInvalidConfig        = "Invalid configuration"
	SyntheticTestErrInvalidConfigMsg     = "Only one of http_action or http_script configuration can be provided"
	SyntheticTestErrInvalidHttpAction    = "Invalid HTTP Action configuration"
	SyntheticTestErrInvalidHttpActionMsg = "Exactly one HTTP Action configuration is required"
	SyntheticTestErrInvalidHttpScript    = "Invalid HTTP Script configuration"
	SyntheticTestErrInvalidHttpScriptMsg = "Exactly one HTTP Script configuration is required"
	SyntheticTestErrNoValidConfig        = "Invalid configuration"
	SyntheticTestErrNoValidConfigMsg     = "No valid configuration provided"

	// Validator description constants
	SyntheticTestValidatorURLRegex = "must be a valid URL with HTTP or HTTPS scheme"
)
