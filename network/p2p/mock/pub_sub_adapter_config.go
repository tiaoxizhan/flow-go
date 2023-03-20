// Code generated by mockery v2.21.4. DO NOT EDIT.

package mockp2p

import (
	p2p "github.com/onflow/flow-go/network/p2p"
	mock "github.com/stretchr/testify/mock"

	routing "github.com/libp2p/go-libp2p/core/routing"
)

// PubSubAdapterConfig is an autogenerated mock type for the PubSubAdapterConfig type
type PubSubAdapterConfig struct {
	mock.Mock
}

// WithAppSpecificRpcInspector provides a mock function with given fields: inspector
func (_m *PubSubAdapterConfig) WithAppSpecificRpcInspector(inspector p2p.BasicGossipSubRPCInspector) {
	_m.Called(inspector)
}

// WithMessageIdFunction provides a mock function with given fields: f
func (_m *PubSubAdapterConfig) WithMessageIdFunction(f func([]byte) string) {
	_m.Called(f)
}

// WithRoutingDiscovery provides a mock function with given fields: _a0
func (_m *PubSubAdapterConfig) WithRoutingDiscovery(_a0 routing.ContentRouting) {
	_m.Called(_a0)
}

// WithScoreOption provides a mock function with given fields: _a0
func (_m *PubSubAdapterConfig) WithScoreOption(_a0 p2p.ScoreOptionBuilder) {
	_m.Called(_a0)
}

// WithScoreTracer provides a mock function with given fields: tracer
func (_m *PubSubAdapterConfig) WithScoreTracer(tracer p2p.PeerScoreTracer) {
	_m.Called(tracer)
}

// WithSubscriptionFilter provides a mock function with given fields: _a0
func (_m *PubSubAdapterConfig) WithSubscriptionFilter(_a0 p2p.SubscriptionFilter) {
	_m.Called(_a0)
}

// WithTracer provides a mock function with given fields: t
func (_m *PubSubAdapterConfig) WithTracer(t p2p.PubSubTracer) {
	_m.Called(t)
}

type mockConstructorTestingTNewPubSubAdapterConfig interface {
	mock.TestingT
	Cleanup(func())
}

// NewPubSubAdapterConfig creates a new instance of PubSubAdapterConfig. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPubSubAdapterConfig(t mockConstructorTestingTNewPubSubAdapterConfig) *PubSubAdapterConfig {
	mock := &PubSubAdapterConfig{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
