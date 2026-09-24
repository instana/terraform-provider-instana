package ipfiltering

// ResourceInstanaIPFiltering is the name of the terraform resource for IP filtering.
const ResourceInstanaIPFiltering = "ip_filtering"

// Schema field name constants
const (
	// IPFilteringFieldEnabled constant for the enabled field
	IPFilteringFieldEnabled = "enabled"
	// IPFilteringFieldDenyAll constant for the deny_all field
	IPFilteringFieldDenyAll = "deny_all"
	// IPFilteringFieldSupportAccessEnabled constant for the support_access_enabled field
	IPFilteringFieldSupportAccessEnabled = "support_access_enabled"
	// IPFilteringFieldRules constant for the rules field
	IPFilteringFieldRules = "rules"
	// IPFilteringFieldRuleTarget constant for the target field in a rule
	IPFilteringFieldRuleTarget = "target"
	// IPFilteringFieldRuleBlock constant for the block field in a rule
	IPFilteringFieldRuleBlock = "block"
	// IPFilteringFieldActive constant for the active field
	IPFilteringFieldActive = "active"
	// IPFilteringFieldLastChangedAt constant for the last_changed_at field
	IPFilteringFieldLastChangedAt = "last_changed_at"
	// IPFilteringFieldLastVerifiedAt constant for the last_verified_at field
	IPFilteringFieldLastVerifiedAt = "last_verified_at"
)

// Resource description constants
const (
	// IPFilteringDescResource describes the resource purpose
	IPFilteringDescResource = "Manages tenant unit IP filtering settings in Instana, " +
		"including traffic rules enforcing allow/block by IP address or range, deny-all fallback, and support access. " +
		"This is a singleton resource — only one instance exists per tenant unit."
	// IPFilteringDescEnabled describes the enabled field
	IPFilteringDescEnabled = "If the IP filtering configuration should be enforced."
	// IPFilteringDescDenyAll describes the deny_all field
	IPFilteringDescDenyAll = "Fallback configuration if unmatched requests should be denied."
	// IPFilteringDescSupportAccessEnabled describes the support_access_enabled field
	IPFilteringDescSupportAccessEnabled = "Optional flag for future support access functionality. Defaults to false."
	// IPFilteringDescRules describes the rules field
	IPFilteringDescRules = "Traffic rules enforcing allow or block based on incoming IP address or range. Evaluated in given order (1 to 75 rules)."
	// IPFilteringDescRuleTarget describes the rule target field
	IPFilteringDescRuleTarget = "Target of this rule, can be an IP address or IP address range (CIDR)."
	// IPFilteringDescRuleBlock describes the rule block field
	IPFilteringDescRuleBlock = "If the traffic should be blocked/denied (true) or allowed (false)."
	// IPFilteringDescActive describes the active field
	IPFilteringDescActive = "Indicates if the configuration is currently enforced (permanent if verified, or temporary after creation)."
	// IPFilteringDescLastChangedAt describes the last_changed_at field
	IPFilteringDescLastChangedAt = "Timestamp when the configuration was last changed."
	// IPFilteringDescLastVerifiedAt describes the last_verified_at field
	IPFilteringDescLastVerifiedAt = "Timestamp when the configuration was last verified."
)

// Validation constants enforced by the Instana API
const (
	// IPFilteringMinRules minimum number of rules
	IPFilteringMinRules = 1
	// IPFilteringMaxRules maximum number of rules
	IPFilteringMaxRules = 75
)
