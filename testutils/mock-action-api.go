package testutils

import (
	"github.com/instana/instana-go-client/shared/rest"
)

// ActionCall records a single call to ActionRestResource.Execute
type ActionCall struct {
	Method rest.ActionMethod
	ID     string
	Action string
}

// FakeActionRestResource is a rest.ActionRestResource which records all calls
type FakeActionRestResource struct {
	Calls []ActionCall
	Err   error
}

// Execute records the call and returns the configured error
func (f *FakeActionRestResource) Execute(method rest.ActionMethod, id string, action string) error {
	f.Calls = append(f.Calls, ActionCall{Method: method, ID: id, Action: action})
	return f.Err
}

// FakeActionInstanaAPI is a MockInstanaAPI which returns a FakeActionRestResource from ResourceActions
// and records the requested base path
type FakeActionInstanaAPI struct {
	MockInstanaAPI
	Actions  *FakeActionRestResource
	BasePath string
}

// NewFakeActionInstanaAPI creates a FakeActionInstanaAPI with a FakeActionRestResource returning the given error
func NewFakeActionInstanaAPI(err error) *FakeActionInstanaAPI {
	return &FakeActionInstanaAPI{Actions: &FakeActionRestResource{Err: err}}
}

// ResourceActions returns the fake action resource and records the base path
func (f *FakeActionInstanaAPI) ResourceActions(basePath string) rest.ActionRestResource {
	f.BasePath = basePath
	return f.Actions
}
