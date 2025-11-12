package restapi

import (
	"encoding/json"
	"fmt"
)

// NewHostAgentJSONUnmarshaller creates a new instance of a generic JSONUnmarshaller.
func NewHostAgentJSONUnmarshaller[T InstanaDataObject](objectType T) JSONUnmarshaller[T] {
	arrayType := make(map[string][]T)
	arrayType["items"] = []T{}

	return &hostAgentJSONUnmarshaller[T]{
		objectType: objectType,
		arrayType:  &arrayType,
	}
}

type hostAgentJSONUnmarshaller[T any] struct {
	objectType T
	arrayType  *map[string][]T
}

// UnmarshalJSON unmarshals JSON data into the target object.
func (u *hostAgentJSONUnmarshaller[T]) Unmarshal(data []byte) (T, error) {
	target := u.objectType
	if err := json.Unmarshal(data, &target); err != nil {
		return target, fmt.Errorf("failed to parse json: %w", err)
	}
	return target, nil
}

// UnmarshalJSONArray unmarshals JSON array data into a slice of target objects.
func (u *hostAgentJSONUnmarshaller[T]) UnmarshalArray(data []byte) (*[]T, error) {
	target := u.arrayType
	if err := json.Unmarshal(data, &target); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}
	hostAgents := (*target)["items"]
	return &hostAgents, nil
}
