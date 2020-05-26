// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import context "context"
import delta "github.com/dapperlabs/flow-go/engine/execution/state/delta"
import entity "github.com/dapperlabs/flow-go/module/mempool/entity"
import execution "github.com/dapperlabs/flow-go/engine/execution"
import flow "github.com/dapperlabs/flow-go/model/flow"
import mock "github.com/stretchr/testify/mock"

// ComputationManager is an autogenerated mock type for the ComputationManager type
type ComputationManager struct {
	mock.Mock
}

// ComputeBlock provides a mock function with given fields: ctx, block, view
func (_m *ComputationManager) ComputeBlock(ctx context.Context, block *entity.ExecutableBlock, view *delta.View) (*execution.ComputationResult, error) {
	ret := _m.Called(ctx, block, view)

	var r0 *execution.ComputationResult
	if rf, ok := ret.Get(0).(func(context.Context, *entity.ExecutableBlock, *delta.View) *execution.ComputationResult); ok {
		r0 = rf(ctx, block, view)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.ComputationResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.ExecutableBlock, *delta.View) error); ok {
		r1 = rf(ctx, block, view)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExecuteScript provides a mock function with given fields: _a0, _a1, _a2
func (_m *ComputationManager) ExecuteScript(_a0 []byte, _a1 *flow.Header, _a2 *delta.View) ([]byte, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte, *flow.Header, *delta.View) []byte); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, *flow.Header, *delta.View) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccount provides a mock function with given fields: addr, header, view
func (_m *ComputationManager) GetAccount(addr flow.Address, header *flow.Header, view *delta.View) (*flow.Account, error) {
	ret := _m.Called(addr, header, view)

	var r0 *flow.Account
	if rf, ok := ret.Get(0).(func(flow.Address, *flow.Header, *delta.View) *flow.Account); ok {
		r0 = rf(addr, header, view)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Address, *flow.Header, *delta.View) error); ok {
		r1 = rf(addr, header, view)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
