// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	crypto "github.com/onflow/crypto"
	flow "github.com/onflow/flow-go/model/flow"

	messages "github.com/onflow/flow-go/model/messages"

	mock "github.com/stretchr/testify/mock"
)

// DKGBroker is an autogenerated mock type for the DKGBroker type
type DKGBroker struct {
	mock.Mock
}

// Broadcast provides a mock function with given fields: data
func (_m *DKGBroker) Broadcast(data []byte) {
	_m.Called(data)
}

// Disqualify provides a mock function with given fields: participant, log
func (_m *DKGBroker) Disqualify(participant int, log string) {
	_m.Called(participant, log)
}

// FlagMisbehavior provides a mock function with given fields: participant, log
func (_m *DKGBroker) FlagMisbehavior(participant int, log string) {
	_m.Called(participant, log)
}

// GetBroadcastMsgCh provides a mock function with given fields:
func (_m *DKGBroker) GetBroadcastMsgCh() <-chan messages.BroadcastDKGMessage {
	ret := _m.Called()

	var r0 <-chan messages.BroadcastDKGMessage
	if rf, ok := ret.Get(0).(func() <-chan messages.BroadcastDKGMessage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan messages.BroadcastDKGMessage)
		}
	}

	return r0
}

// GetIndex provides a mock function with given fields:
func (_m *DKGBroker) GetIndex() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetPrivateMsgCh provides a mock function with given fields:
func (_m *DKGBroker) GetPrivateMsgCh() <-chan messages.PrivDKGMessageIn {
	ret := _m.Called()

	var r0 <-chan messages.PrivDKGMessageIn
	if rf, ok := ret.Get(0).(func() <-chan messages.PrivDKGMessageIn); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan messages.PrivDKGMessageIn)
		}
	}

	return r0
}

// Poll provides a mock function with given fields: referenceBlock
func (_m *DKGBroker) Poll(referenceBlock flow.Identifier) error {
	ret := _m.Called(referenceBlock)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) error); ok {
		r0 = rf(referenceBlock)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PrivateSend provides a mock function with given fields: dest, data
func (_m *DKGBroker) PrivateSend(dest int, data []byte) {
	_m.Called(dest, data)
}

// Shutdown provides a mock function with given fields:
func (_m *DKGBroker) Shutdown() {
	_m.Called()
}

// SubmitResult provides a mock function with given fields: _a0, _a1
func (_m *DKGBroker) SubmitResult(_a0 crypto.PublicKey, _a1 []crypto.PublicKey) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(crypto.PublicKey, []crypto.PublicKey) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewDKGBroker interface {
	mock.TestingT
	Cleanup(func())
}

// NewDKGBroker creates a new instance of DKGBroker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDKGBroker(t mockConstructorTestingTNewDKGBroker) *DKGBroker {
	mock := &DKGBroker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
