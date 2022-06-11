// Code generated by mockery v2.12.1. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	flow "github.com/onflow/flow-go/model/flow"

	testing "testing"
)

// SealingCore is an autogenerated mock type for the SealingCore type
type SealingCore struct {
	mock.Mock
}

// ProcessApproval provides a mock function with given fields: approval
func (_m *SealingCore) ProcessApproval(approval *flow.ResultApproval) error {
	ret := _m.Called(approval)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.ResultApproval) error); ok {
		r0 = rf(approval)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProcessFinalizedBlock provides a mock function with given fields: finalizedBlockID
func (_m *SealingCore) ProcessFinalizedBlock(finalizedBlockID flow.Identifier) error {
	ret := _m.Called(finalizedBlockID)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) error); ok {
		r0 = rf(finalizedBlockID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProcessIncorporatedResult provides a mock function with given fields: result
func (_m *SealingCore) ProcessIncorporatedResult(result *flow.IncorporatedResult) error {
	ret := _m.Called(result)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.IncorporatedResult) error); ok {
		r0 = rf(result)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewSealingCore creates a new instance of SealingCore. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewSealingCore(t testing.TB) *SealingCore {
	mock := &SealingCore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
