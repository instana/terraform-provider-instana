package restapi

// InfraAlertEvaluationType custom type representing the infrastructure alert evaluation type from the Instana API
type InfraAlertEvaluationType string

// InfraAlertEvaluationTypes custom type representing a slice of InfraAlertEvaluationType
type InfraAlertEvaluationTypes []InfraAlertEvaluationType

// ToStringSlice returns the corresponding string representations
func (types InfraAlertEvaluationTypes) ToStringSlice() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = string(v)
	}
	return result
}

const (
	// EvaluationTypePerEntity constant value for InfraAlertEvaluationType PER_ENTITY
	EvaluationTypePerEntity = InfraAlertEvaluationType("PER_ENTITY")
	// EvaluationTypeCustom constant value for InfraAlertEvaluationType CUSTOM
	EvaluationTypeCustom = InfraAlertEvaluationType("CUSTOM")
)

// SupportedInfraAlertEvaluationTypes list of all supported InfraAlertEvaluationTypes
var SupportedInfraAlertEvaluationTypes = InfraAlertEvaluationTypes{
	EvaluationTypePerEntity,
	EvaluationTypeCustom,
}
