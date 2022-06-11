package common

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/onflow/flow-go-sdk"

	"github.com/onflow/flow-go/integration/testnet"
	"github.com/onflow/flow-go/integration/tests/lib"
	"github.com/onflow/flow-go/model/flow"
)

// SealingAndVerificationHappyPathTest evaluates the health of the happy path of verification and sealing. It
// deploys a transaction into the testnet hence causing an execution result with more than
// one chunk, assigns all chunks to the same single verification node in this testnet, and then verifies whether verification node
// generates a result approval for all chunks of that execution result.
// It also enables sealing based on result approvals and verifies whether the block of that specific multi-chunk execution result is sealed
// affected by the emitted result approvals.
//
// It returns the execution receipts and result approvals that are involved in this test.
func SealingAndVerificationHappyPathTest(
	t *testing.T,
	blockState *lib.BlockState,
	receiptState *lib.ReceiptState,
	approvalState *lib.ResultApprovalState,
	accessClient *testnet.Client,
	exe1Id flow.Identifier,
	exe2Id flow.Identifier,
	verId flow.Identifier,
	rootBlockId flow.Identifier) ([]*flow.ExecutionReceipt, []*flow.ResultApproval) {

	// wait for next height finalized (potentially first height), called blockA, just to make sure consensus progresses.
	blockA := blockState.WaitForHighestFinalizedProgress(t)
	t.Logf("blockA generated, height: %v ID: %v\n", blockA.Header.Height, blockA.Header.ID())

	// sends a transaction
	err := accessClient.DeployContract(context.Background(), sdk.Identifier(rootBlockId), lib.CounterContract)
	require.NoError(t, err, "could not deploy counter")

	// waits until for a different state commitment for a finalized block, call that block blockB,
	// which has more than one chunk on its execution result.
	blockB, _ := lib.WaitUntilFinalizedStateCommitmentChanged(t, blockState, receiptState, lib.WithMinimumChunks(2))
	t.Logf("got blockB height %v ID %v\n", blockB.Header.Height, blockB.Header.ID())

	// waits for the execution receipt of blockB from both execution nodes, and makes sure that there is no execution fork.
	receiptB1 := receiptState.WaitForReceiptFrom(t, blockB.Header.ID(), exe1Id)
	t.Logf("receipt for blockB generated by execution node-1: %x result ID: %x\n", exe1Id, receiptB1.ExecutionResult.ID())
	receiptB2 := receiptState.WaitForReceiptFrom(t, blockB.Header.ID(), exe2Id)
	t.Logf("receipt for blockB generated by execution node-2: %x result ID: %x\n", exe2Id, receiptB2.ExecutionResult.ID())

	require.Equal(t, receiptB1.ExecutionResult.ID(), receiptB2.ExecutionResult.ID(), "execution fork happened at blockB")
	resultB := receiptB1.ExecutionResult
	resultBId := resultB.ID()
	// re-evaluates that resultB has more than one chunk.
	require.Greater(t, len(resultB.Chunks), 1)
	t.Logf("receipt for blockB generated: result ID: %x with %d chunks\n", resultBId, len(resultB.Chunks))

	// waits till result approval emits for all chunks of resultB
	approvals := make([]*flow.ResultApproval, 0)
	for i := 0; i < len(resultB.Chunks); i++ {
		approval := approvalState.WaitForResultApproval(t, verId, resultBId, uint64(i))
		approvals = append(approvals, approval)
	}

	// waits until blockB is sealed by consensus nodes after result approvals for all of its chunks emitted.
	blockState.WaitForSealed(t, blockB.Header.Height)

	return []*flow.ExecutionReceipt{receiptB1, receiptB2}, approvals
}
