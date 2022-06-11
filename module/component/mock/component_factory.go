// Code generated by mockery v2.12.1. DO NOT EDIT.

package component

import (
	mock "github.com/stretchr/testify/mock"

	component "github.com/onflow/flow-go/module/component"

	testing "testing"
)

// ComponentFactory is an autogenerated mock type for the ComponentFactory type
type ComponentFactory struct {
	mock.Mock
}

// Execute provides a mock function with given fields:
func (_m *ComponentFactory) Execute() (component.Component, error) {
	ret := _m.Called()

	var r0 component.Component
	if rf, ok := ret.Get(0).(func() component.Component); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(component.Component)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewComponentFactory creates a new instance of ComponentFactory. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewComponentFactory(t testing.TB) *ComponentFactory {
	mock := &ComponentFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
