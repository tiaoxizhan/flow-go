// Code generated by mockery v2.13.0. DO NOT EDIT.

package mockfetcher

import (
	fetcher "github.com/onflow/flow-go/engine/verification/fetcher"
	mock "github.com/stretchr/testify/mock"

	verification "github.com/onflow/flow-go/model/verification"
)

// ChunkDataPackRequester is an autogenerated mock type for the ChunkDataPackRequester type
type ChunkDataPackRequester struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *ChunkDataPackRequester) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *ChunkDataPackRequester) Ready() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Request provides a mock function with given fields: request
func (_m *ChunkDataPackRequester) Request(request *verification.ChunkDataPackRequest) {
	_m.Called(request)
}

// WithChunkDataPackHandler provides a mock function with given fields: handler
func (_m *ChunkDataPackRequester) WithChunkDataPackHandler(handler fetcher.ChunkDataPackHandler) {
	_m.Called(handler)
}

type NewChunkDataPackRequesterT interface {
	mock.TestingT
	Cleanup(func())
}

// NewChunkDataPackRequester creates a new instance of ChunkDataPackRequester. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewChunkDataPackRequester(t NewChunkDataPackRequesterT) *ChunkDataPackRequester {
	mock := &ChunkDataPackRequester{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
