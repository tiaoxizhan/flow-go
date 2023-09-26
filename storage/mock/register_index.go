// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// RegisterIndex is an autogenerated mock type for the RegisterIndex type
type RegisterIndex struct {
	mock.Mock
}

// FirstHeight provides a mock function with given fields:
func (_m *RegisterIndex) FirstHeight() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// Get provides a mock function with given fields: ID, height
func (_m *RegisterIndex) Get(ID flow.RegisterID, height uint64) ([]byte, error) {
	ret := _m.Called(ID, height)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.RegisterID, uint64) ([]byte, error)); ok {
		return rf(ID, height)
	}
	if rf, ok := ret.Get(0).(func(flow.RegisterID, uint64) []byte); ok {
		r0 = rf(ID, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.RegisterID, uint64) error); ok {
		r1 = rf(ID, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LatestHeight provides a mock function with given fields:
func (_m *RegisterIndex) LatestHeight() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// Store provides a mock function with given fields: entries, height
func (_m *RegisterIndex) Store(entries flow.RegisterEntries, height uint64) error {
	ret := _m.Called(entries, height)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.RegisterEntries, uint64) error); ok {
		r0 = rf(entries, height)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRegisterIndex interface {
	mock.TestingT
	Cleanup(func())
}

// NewRegisterIndex creates a new instance of RegisterIndex. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRegisterIndex(t mockConstructorTestingTNewRegisterIndex) *RegisterIndex {
	mock := &RegisterIndex{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
