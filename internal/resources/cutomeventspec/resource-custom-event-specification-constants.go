package instana

// Resource description constants
const (
	CustomEventSpecificationResourceDescResource                 = "This resource represents a custom event specification in Instana"
	CustomEventSpecificationResourceDescID                       = "The ID of the custom event specification"
	CustomEventSpecificationResourceDescName                     = "The name of the custom event specification"
	CustomEventSpecificationResourceDescEntityType               = "The entity type of the custom event specification"
	CustomEventSpecificationResourceDescQuery                    = "The query of the custom event specification"
	CustomEventSpecificationResourceDescTriggering               = "Indicates if the custom event specification is triggering"
	CustomEventSpecificationResourceDescDescription              = "The description of the custom event specification"
	CustomEventSpecificationResourceDescExpirationTime           = "The expiration time of the custom event specification"
	CustomEventSpecificationResourceDescEnabled                  = "Indicates if the custom event specification is enabled"
	CustomEventSpecificationResourceDescRuleLogicalOperator      = "The logical operator for the rules (AND, OR)"
	CustomEventSpecificationResourceDescRules                    = "The rules of the custom event specification"
	CustomEventSpecificationResourceDescEntityCountRules         = "Entity count rules"
	CustomEventSpecificationResourceDescEntityCountVerification  = "Entity count verification rules"
	CustomEventSpecificationResourceDescEntityVerification       = "Entity verification rules"
	CustomEventSpecificationResourceDescHostAvailability         = "Host availability rules"
	CustomEventSpecificationResourceDescSystemRules              = "System rules"
	CustomEventSpecificationResourceDescThresholdRules           = "Threshold rules"
	CustomEventSpecificationResourceDescSeverity                 = "The severity of the rule (warning, critical)"
	CustomEventSpecificationResourceDescConditionOperator        = "The condition operator of the rule"
	CustomEventSpecificationResourceDescConditionValue           = "The condition value of the rule"
	CustomEventSpecificationResourceDescMatchingEntityType       = "The matching entity type of the rule"
	CustomEventSpecificationResourceDescMatchingOperator         = "The matching operator of the rule"
	CustomEventSpecificationResourceDescMatchingEntityLabel      = "The matching entity label of the rule"
	CustomEventSpecificationResourceDescOfflineDuration          = "The offline duration of the rule"
	CustomEventSpecificationResourceDescSystemRuleID             = "The system rule ID"
	CustomEventSpecificationResourceDescMetricName               = "The metric name of the rule"
	CustomEventSpecificationResourceDescRollup                   = "The rollup of the rule"
	CustomEventSpecificationResourceDescWindow                   = "The window of the rule"
	CustomEventSpecificationResourceDescAggregation              = "The aggregation of the rule"
	CustomEventSpecificationResourceDescCloseAfter               = "The close after duration of the rule"
	CustomEventSpecificationResourceDescTagFilter                = "The tag filter of the rule"
	CustomEventSpecificationResourceDescMetricPattern            = "The metric pattern of the rule"
	CustomEventSpecificationResourceDescMetricPatternPrefix      = "The prefix of the metric pattern"
	CustomEventSpecificationResourceDescMetricPatternPostfix     = "The postfix of the metric pattern"
	CustomEventSpecificationResourceDescMetricPatternPlaceholder = "The placeholder of the metric pattern"
	CustomEventSpecificationResourceDescMetricPatternOperator    = "The operator of the metric pattern"
)

// Error message constants
const (
	CustomEventSpecificationResourceErrParseTagFilter    = "Error parsing tag filter"
	CustomEventSpecificationResourceErrParseTagFilterMsg = "Failed to parse tag filter: %s"
)

// Made with Bob
