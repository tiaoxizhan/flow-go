// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	fvm "github.com/onflow/flow-go/fvm"
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"

	programs "github.com/onflow/flow-go/fvm/programs"

	state "github.com/onflow/flow-go/fvm/state"
)

// VirtualMachine is an autogenerated mock type for the VirtualMachine type
type VirtualMachine struct {
	mock.Mock
}

// GetAccount provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *VirtualMachine) GetAccount(_a0 fvm.Context, _a1 flow.Address, _a2 state.View, _a3 *programs.Programs) (*flow.Account, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *flow.Account
	if rf, ok := ret.Get(0).(func(fvm.Context, flow.Address, state.View, *programs.Programs) *flow.Account); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(fvm.Context, flow.Address, state.View, *programs.Programs) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Run provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *VirtualMachine) Run(_a0 fvm.Context, _a1 fvm.Procedure, _a2 state.View, _a3 *programs.Programs) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(fvm.Context, fvm.Procedure, state.View, *programs.Programs) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunV2 provides a mock function with given fields: _a0, _a1, _a2
func (_m *VirtualMachine) RunV2(_a0 fvm.Context, _a1 fvm.Procedure, _a2 state.View) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(fvm.Context, fvm.Procedure, state.View) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewVirtualMachine interface {
	mock.TestingT
	Cleanup(func())
}

// NewVirtualMachine creates a new instance of VirtualMachine. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewVirtualMachine(t mockConstructorTestingTNewVirtualMachine) *VirtualMachine {
	mock := &VirtualMachine{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
