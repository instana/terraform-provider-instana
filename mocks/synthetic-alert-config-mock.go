package mocks

import (
	"reflect"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"go.uber.org/mock/gomock"
)

// SyntheticAlertConfigs mocks base method.
func (m *MockInstanaAPI) SyntheticAlertConfigs() restapi.RestResource[*restapi.SyntheticAlertConfig] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyntheticAlertConfigs")
	ret0, _ := ret[0].(restapi.RestResource[*restapi.SyntheticAlertConfig])
	return ret0
}

// SyntheticAlertConfigs indicates an expected call of SyntheticAlertConfigs.
func (mr *MockInstanaAPIMockRecorder) SyntheticAlertConfigs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyntheticAlertConfigs", reflect.TypeOf((*MockInstanaAPI)(nil).SyntheticAlertConfigs))
}
