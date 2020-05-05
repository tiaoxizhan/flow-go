// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import flow "github.com/dapperlabs/flow-go/model/flow"
import mock "github.com/stretchr/testify/mock"

import opentracing "github.com/opentracing/opentracing-go"
import time "time"

// Metrics is an autogenerated mock type for the Metrics type
type Metrics struct {
	mock.Mock
}

// BadgerLSMSize provides a mock function with given fields: sizeBytes
func (_m *Metrics) BadgerLSMSize(sizeBytes int64) {
	_m.Called(sizeBytes)
}

// BadgerNumBlockedPuts provides a mock function with given fields: n
func (_m *Metrics) BadgerNumBlockedPuts(n int64) {
	_m.Called(n)
}

// BadgerNumBytesRead provides a mock function with given fields: n
func (_m *Metrics) BadgerNumBytesRead(n int64) {
	_m.Called(n)
}

// BadgerNumBytesWritten provides a mock function with given fields: n
func (_m *Metrics) BadgerNumBytesWritten(n int64) {
	_m.Called(n)
}

// BadgerNumGets provides a mock function with given fields: n
func (_m *Metrics) BadgerNumGets(n int64) {
	_m.Called(n)
}

// BadgerNumMemtableGets provides a mock function with given fields: n
func (_m *Metrics) BadgerNumMemtableGets(n int64) {
	_m.Called(n)
}

// BadgerNumPuts provides a mock function with given fields: n
func (_m *Metrics) BadgerNumPuts(n int64) {
	_m.Called(n)
}

// BadgerNumReads provides a mock function with given fields: n
func (_m *Metrics) BadgerNumReads(n int64) {
	_m.Called(n)
}

// BadgerNumWrites provides a mock function with given fields: n
func (_m *Metrics) BadgerNumWrites(n int64) {
	_m.Called(n)
}

// BadgerVLogSize provides a mock function with given fields: sizeBytes
func (_m *Metrics) BadgerVLogSize(sizeBytes int64) {
	_m.Called(sizeBytes)
}

// CollectionGuaranteed provides a mock function with given fields: collection
func (_m *Metrics) CollectionGuaranteed(collection flow.LightCollection) {
	_m.Called(collection)
}

// CollectionProposed provides a mock function with given fields: collection
func (_m *Metrics) CollectionProposed(collection flow.LightCollection) {
	_m.Called(collection)
}

// CollectionsInFinalizedBlock provides a mock function with given fields: count
func (_m *Metrics) CollectionsInFinalizedBlock(count int) {
	_m.Called(count)
}

// CollectionsPerBlock provides a mock function with given fields: count
func (_m *Metrics) CollectionsPerBlock(count int) {
	_m.Called(count)
}

// ExecutionGasUsedPerBlock provides a mock function with given fields: gas
func (_m *Metrics) ExecutionGasUsedPerBlock(gas uint64) {
	_m.Called(gas)
}

// ExecutionStateReadsPerBlock provides a mock function with given fields: reads
func (_m *Metrics) ExecutionStateReadsPerBlock(reads uint64) {
	_m.Called(reads)
}

// ExecutionStateStorageDiskTotal provides a mock function with given fields: bytes
func (_m *Metrics) ExecutionStateStorageDiskTotal(bytes int64) {
	_m.Called(bytes)
}

// ExecutionStorageStateCommitment provides a mock function with given fields: bytes
func (_m *Metrics) ExecutionStorageStateCommitment(bytes int64) {
	_m.Called(bytes)
}

// FinalizedBlocks provides a mock function with given fields: count
func (_m *Metrics) FinalizedBlocks(count int) {
	_m.Called(count)
}

// FinishBlockReceivedToExecuted provides a mock function with given fields: blockID
func (_m *Metrics) FinishBlockReceivedToExecuted(blockID flow.Identifier) {
	_m.Called(blockID)
}

// FinishBlockToSeal provides a mock function with given fields: blockID
func (_m *Metrics) FinishBlockToSeal(blockID flow.Identifier) {
	_m.Called(blockID)
}

// FinishCollectionToFinalized provides a mock function with given fields: collectionID
func (_m *Metrics) FinishCollectionToFinalized(collectionID flow.Identifier) {
	_m.Called(collectionID)
}

// HotStuffBusyDuration provides a mock function with given fields: duration, event
func (_m *Metrics) HotStuffBusyDuration(duration time.Duration, event string) {
	_m.Called(duration, event)
}

// HotStuffIdleDuration provides a mock function with given fields: duration
func (_m *Metrics) HotStuffIdleDuration(duration time.Duration) {
	_m.Called(duration)
}

// HotStuffWaitDuration provides a mock function with given fields: duration, event
func (_m *Metrics) HotStuffWaitDuration(duration time.Duration, event string) {
	_m.Called(duration, event)
}

// MadeBlockProposal provides a mock function with given fields:
func (_m *Metrics) MadeBlockProposal() {
	_m.Called()
}

// MempoolApprovalsSize provides a mock function with given fields: size
func (_m *Metrics) MempoolApprovalsSize(size uint) {
	_m.Called(size)
}

// MempoolGuaranteesSize provides a mock function with given fields: size
func (_m *Metrics) MempoolGuaranteesSize(size uint) {
	_m.Called(size)
}

// MempoolReceiptsSize provides a mock function with given fields: size
func (_m *Metrics) MempoolReceiptsSize(size uint) {
	_m.Called(size)
}

// MempoolSealsSize provides a mock function with given fields: size
func (_m *Metrics) MempoolSealsSize(size uint) {
	_m.Called(size)
}

// NetworkMessageReceived provides a mock function with given fields: sizeBytes, topic
func (_m *Metrics) NetworkMessageReceived(sizeBytes int, topic string) {
	_m.Called(sizeBytes, topic)
}

// NetworkMessageSent provides a mock function with given fields: sizeBytes, topic
func (_m *Metrics) NetworkMessageSent(sizeBytes int, topic string) {
	_m.Called(sizeBytes, topic)
}

// NewestKnownQC provides a mock function with given fields: view
func (_m *Metrics) NewestKnownQC(view uint64) {
	_m.Called(view)
}

// OnChunkDataAdded provides a mock function with given fields: chunkID, size
func (_m *Metrics) OnChunkDataAdded(chunkID flow.Identifier, size float64) {
	_m.Called(chunkID, size)
}

// OnChunkDataRemoved provides a mock function with given fields: chunkID, size
func (_m *Metrics) OnChunkDataRemoved(chunkID flow.Identifier, size float64) {
	_m.Called(chunkID, size)
}

// OnChunkVerificationFinished provides a mock function with given fields: chunkID
func (_m *Metrics) OnChunkVerificationFinished(chunkID flow.Identifier) {
	_m.Called(chunkID)
}

// OnChunkVerificationStarted provides a mock function with given fields: chunkID
func (_m *Metrics) OnChunkVerificationStarted(chunkID flow.Identifier) {
	_m.Called(chunkID)
}

// OnResultApproval provides a mock function with given fields:
func (_m *Metrics) OnResultApproval() {
	_m.Called()
}

// PendingBlocks provides a mock function with given fields: n
func (_m *Metrics) PendingBlocks(n uint) {
	_m.Called(n)
}

// PendingClusterBlocks provides a mock function with given fields: n
func (_m *Metrics) PendingClusterBlocks(n uint) {
	_m.Called(n)
}

// SealsInFinalizedBlock provides a mock function with given fields: count
func (_m *Metrics) SealsInFinalizedBlock(count int) {
	_m.Called(count)
}

// StartBlockReceivedToExecuted provides a mock function with given fields: blockID
func (_m *Metrics) StartBlockReceivedToExecuted(blockID flow.Identifier) {
	_m.Called(blockID)
}

// StartBlockToSeal provides a mock function with given fields: blockID
func (_m *Metrics) StartBlockToSeal(blockID flow.Identifier) {
	_m.Called(blockID)
}

// StartCollectionToFinalized provides a mock function with given fields: collectionID
func (_m *Metrics) StartCollectionToFinalized(collectionID flow.Identifier) {
	_m.Called(collectionID)
}

// StartNewView provides a mock function with given fields: view
func (_m *Metrics) StartNewView(view uint64) {
	_m.Called(view)
}

// Tracer provides a mock function with given fields:
func (_m *Metrics) Tracer() opentracing.Tracer {
	ret := _m.Called()

	var r0 opentracing.Tracer
	if rf, ok := ret.Get(0).(func() opentracing.Tracer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(opentracing.Tracer)
		}
	}

	return r0
}

// TransactionReceived provides a mock function with given fields: txID
func (_m *Metrics) TransactionReceived(txID flow.Identifier) {
	_m.Called(txID)
}
