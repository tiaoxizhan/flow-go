package execution_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/atomic"

	"github.com/onflow/flow-go/consensus/hotstuff/helper"
	"github.com/onflow/flow-go/engine"
	execTestutil "github.com/onflow/flow-go/engine/execution/testutil"
	"github.com/onflow/flow-go/engine/testutil"
	testmock "github.com/onflow/flow-go/engine/testutil/mock"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/order"
	"github.com/onflow/flow-go/model/messages"
	"github.com/onflow/flow-go/module/signature"
	"github.com/onflow/flow-go/network/mocknetwork"
	"github.com/onflow/flow-go/network/stub"
	"github.com/onflow/flow-go/state/cluster"
	"github.com/onflow/flow-go/utils/unittest"
)

func sendBlock(exeNode *testmock.ExecutionNode, from flow.Identifier, proposal *messages.BlockProposal) error {
	return exeNode.FollowerEngine.Process(engine.ReceiveBlocks, from, proposal)
}

// Test when the ingestion engine receives a block, it will
// request collections from collection node, and send ER to
// verification node and consensus node.
// create a block that has two collections: col1 and col2;
// col1 has tx1 and tx2, col2 has tx3 and tx4.
// create another child block which will trigger the parent
// block to be incorporated and be passed to the ingestion engine
func TestExecutionFlow(t *testing.T) {
	hub := stub.NewNetworkHub()

	chainID := flow.Testnet

	colID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleCollection),
		unittest.WithKeys,
	)
	conID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleConsensus),
		unittest.WithKeys,
	)
	exeID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleExecution),
		unittest.WithKeys,
	)
	verID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleVerification),
		unittest.WithKeys,
	)

	identities := unittest.CompleteIdentitySet(colID, conID, exeID, verID).Sort(order.Canonical)

	// create execution node
	exeNode := testutil.ExecutionNode(t, hub, exeID, identities, 21, chainID)
	exeNode.Ready()
	defer exeNode.Done()

	genesis, err := exeNode.State.AtHeight(0).Head()
	require.NoError(t, err)

	tx1 := flow.TransactionBody{
		Script: []byte("transaction { execute { log(1) } }"),
	}

	tx2 := flow.TransactionBody{
		Script: []byte("transaction { execute { log(2) } }"),
	}

	tx3 := flow.TransactionBody{
		Script: []byte("transaction { execute { log(3) } }"),
	}

	tx4 := flow.TransactionBody{
		Script: []byte("transaction { execute { log(4) } }"),
	}

	col1 := flow.Collection{Transactions: []*flow.TransactionBody{&tx1, &tx2}}
	col2 := flow.Collection{Transactions: []*flow.TransactionBody{&tx3, &tx4}}

	collections := map[flow.Identifier]*flow.Collection{
		col1.ID(): &col1,
		col2.ID(): &col2,
	}

	clusterChainID := cluster.CanonicalClusterID(1, flow.IdentityList{colID})

	// signed by the only collector
	block := unittest.BlockWithParentAndProposerFixture(t, genesis, conID.NodeID)
	signerIndices, err := signature.EncodeSignersToIndices(
		[]flow.Identifier{colID.NodeID}, []flow.Identifier{colID.NodeID})
	require.NoError(t, err)
	block.SetPayload(flow.Payload{
		Guarantees: []*flow.CollectionGuarantee{
			{
				CollectionID:     col1.ID(),
				SignerIndices:    signerIndices,
				ChainID:          clusterChainID,
				ReferenceBlockID: genesis.ID(),
			},
			{
				CollectionID:     col2.ID(),
				SignerIndices:    signerIndices,
				ChainID:          clusterChainID,
				ReferenceBlockID: genesis.ID(),
			},
		},
	})

	child := unittest.BlockWithParentAndProposerFixture(t, block.Header, conID.NodeID)

	log.Info().Msgf("child block ID: %v, indices: %x", child.Header.ID(), child.Header.ParentVoterIndices)

	collectionNode := testutil.GenericNodeFromParticipants(t, hub, colID, identities, chainID)
	defer collectionNode.Done()
	verificationNode := testutil.GenericNodeFromParticipants(t, hub, verID, identities, chainID)
	defer verificationNode.Done()
	consensusNode := testutil.GenericNodeFromParticipants(t, hub, conID, identities, chainID)
	defer consensusNode.Done()

	// create collection node that can respond collections to execution node
	// check collection node received the collection request from execution node
	providerEngine := new(mocknetwork.Engine)
	provConduit, _ := collectionNode.Net.Register(engine.ProvideCollections, providerEngine)
	providerEngine.On("Process", mock.AnythingOfType("network.Channel"), exeID.NodeID, mock.Anything).
		Run(func(args mock.Arguments) {
			originID := args.Get(1).(flow.Identifier)
			req := args.Get(2).(*messages.EntityRequest)

			var entities []flow.Entity
			for _, entityID := range req.EntityIDs {
				coll, exists := collections[entityID]
				require.True(t, exists)
				entities = append(entities, coll)
			}

			var blobs [][]byte
			for _, entity := range entities {
				blob, _ := msgpack.Marshal(entity)
				blobs = append(blobs, blob)
			}

			res := &messages.EntityResponse{
				Nonce:     req.Nonce,
				EntityIDs: req.EntityIDs,
				Blobs:     blobs,
			}

			err := provConduit.Publish(res, originID)
			assert.NoError(t, err)
		}).
		Once().
		Return(nil)

	var lock sync.Mutex
	var receipt *flow.ExecutionReceipt

	// create verification engine that can create approvals and send to consensus nodes
	// check the verification engine received the ER from execution node
	verificationEngine := new(mocknetwork.Engine)
	_, _ = verificationNode.Net.Register(engine.ReceiveReceipts, verificationEngine)
	verificationEngine.On("Process", mock.AnythingOfType("network.Channel"), exeID.NodeID, mock.Anything).
		Run(func(args mock.Arguments) {
			lock.Lock()
			defer lock.Unlock()
			receipt, _ = args[2].(*flow.ExecutionReceipt)

			assert.Equal(t, block.ID(), receipt.ExecutionResult.BlockID)
		}).
		Return(nil).
		Once()

	// create consensus engine that accepts the result
	// check the consensus engine has received the result from execution node
	consensusEngine := new(mocknetwork.Engine)
	_, _ = consensusNode.Net.Register(engine.ReceiveReceipts, consensusEngine)
	consensusEngine.On("Process", mock.AnythingOfType("network.Channel"), exeID.NodeID, mock.Anything).
		Run(func(args mock.Arguments) {
			lock.Lock()
			defer lock.Unlock()

			receipt, _ = args[2].(*flow.ExecutionReceipt)

			assert.Equal(t, block.ID(), receipt.ExecutionResult.BlockID)
			assert.Equal(t, len(collections), len(receipt.ExecutionResult.Chunks)-1) // don't count system chunk

			for i, chunk := range receipt.ExecutionResult.Chunks {
				assert.EqualValues(t, i, chunk.CollectionIndex)
			}
		}).
		Return(nil).
		Once()

	// submit block from consensus node
	err = sendBlock(&exeNode, conID.NodeID, unittest.ProposalFromBlock(&block))
	require.NoError(t, err)

	// submit the child block from consensus node, which trigger the parent block
	// to be passed to BlockProcessable
	err = sendBlock(&exeNode, conID.NodeID, unittest.ProposalFromBlock(&child))
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		// when sendBlock returned, ingestion engine might not have processed
		// the block yet, because the process is async. we have to wait
		hub.DeliverAll()

		lock.Lock()
		defer lock.Unlock()
		return receipt != nil
	}, time.Second*10, time.Millisecond*500)

	// check that the block has been executed.
	exeNode.AssertHighestExecutedBlock(t, block.Header)

	myReceipt, err := exeNode.MyExecutionReceipts.MyReceipt(block.ID())
	require.NoError(t, err)
	require.NotNil(t, myReceipt)
	require.Equal(t, exeNode.Me.NodeID(), myReceipt.ExecutorID)

	providerEngine.AssertExpectations(t)
	verificationEngine.AssertExpectations(t)
	consensusEngine.AssertExpectations(t)
}

func deployContractBlock(t *testing.T, conID *flow.Identity, colID *flow.Identity, chain flow.Chain, seq uint64, parent *flow.Header, ref *flow.Header) (
	*flow.TransactionBody, *flow.Collection, flow.Block, *messages.BlockProposal, uint64) {
	// make tx
	tx := execTestutil.DeployCounterContractTransaction(chain.ServiceAddress(), chain)
	err := execTestutil.SignTransactionAsServiceAccount(tx, seq, chain)
	require.NoError(t, err)

	// make collection
	col := &flow.Collection{Transactions: []*flow.TransactionBody{tx}}

	signerIndices, err := signature.EncodeSignersToIndices(
		[]flow.Identifier{colID.NodeID}, []flow.Identifier{colID.NodeID})
	require.NoError(t, err)

	clusterChainID := cluster.CanonicalClusterID(1, flow.IdentityList{colID})

	// make block
	block := unittest.BlockWithParentAndProposerFixture(t, parent, conID.NodeID)
	block.SetPayload(flow.Payload{
		Guarantees: []*flow.CollectionGuarantee{
			{
				CollectionID:     col.ID(),
				SignerIndices:    signerIndices,
				ChainID:          clusterChainID,
				ReferenceBlockID: ref.ID(),
			},
		},
	})

	// make proposal
	proposal := unittest.ProposalFromBlock(&block)

	return tx, col, block, proposal, seq + 1
}

func makePanicBlock(t *testing.T, conID *flow.Identity, colID *flow.Identity, chain flow.Chain, seq uint64, parent *flow.Header, ref *flow.Header) (
	*flow.TransactionBody, *flow.Collection, flow.Block, *messages.BlockProposal, uint64) {
	// make tx
	tx := execTestutil.CreateCounterPanicTransaction(chain.ServiceAddress(), chain.ServiceAddress())
	err := execTestutil.SignTransactionAsServiceAccount(tx, seq, chain)
	require.NoError(t, err)

	// make collection
	col := &flow.Collection{Transactions: []*flow.TransactionBody{tx}}

	clusterChainID := cluster.CanonicalClusterID(1, flow.IdentityList{colID})
	// make block
	block := unittest.BlockWithParentAndProposerFixture(t, parent, conID.NodeID)

	signerIndices, err := signature.EncodeSignersToIndices(
		[]flow.Identifier{colID.NodeID}, []flow.Identifier{colID.NodeID})
	require.NoError(t, err)

	block.SetPayload(flow.Payload{
		Guarantees: []*flow.CollectionGuarantee{
			{CollectionID: col.ID(), SignerIndices: signerIndices, ChainID: clusterChainID, ReferenceBlockID: ref.ID()},
		},
	})

	proposal := unittest.ProposalFromBlock(&block)

	return tx, col, block, proposal, seq + 1
}

func makeSuccessBlock(t *testing.T, conID *flow.Identity, colID *flow.Identity, chain flow.Chain, seq uint64, parent *flow.Header, ref *flow.Header) (
	*flow.TransactionBody, *flow.Collection, flow.Block, *messages.BlockProposal, uint64) {
	tx := execTestutil.AddToCounterTransaction(chain.ServiceAddress(), chain.ServiceAddress())
	err := execTestutil.SignTransactionAsServiceAccount(tx, seq, chain)
	require.NoError(t, err)

	signerIndices, err := signature.EncodeSignersToIndices(
		[]flow.Identifier{colID.NodeID}, []flow.Identifier{colID.NodeID})
	require.NoError(t, err)
	clusterChainID := cluster.CanonicalClusterID(1, flow.IdentityList{colID})

	col := &flow.Collection{Transactions: []*flow.TransactionBody{tx}}
	block := unittest.BlockWithParentAndProposerFixture(t, parent, conID.NodeID)
	block.SetPayload(flow.Payload{
		Guarantees: []*flow.CollectionGuarantee{
			{CollectionID: col.ID(), SignerIndices: signerIndices, ChainID: clusterChainID, ReferenceBlockID: ref.ID()},
		},
	})

	proposal := unittest.ProposalFromBlock(&block)

	return tx, col, block, proposal, seq + 1
}

// Test a successful tx should change the statecommitment,
// but a failed Tx should not change the statecommitment.
func TestFailedTxWillNotChangeStateCommitment(t *testing.T) {
	hub := stub.NewNetworkHub()

	chainID := flow.Emulator

	colID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleCollection),
		unittest.WithKeys,
	)
	conID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleConsensus),
		unittest.WithKeys,
	)
	exe1ID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleExecution),
		unittest.WithKeys,
	)

	identities := unittest.CompleteIdentitySet(colID, conID, exe1ID)
	key := unittest.NetworkingPrivKeyFixture()
	identities[3].NetworkPubKey = key.PublicKey()

	collectionNode := testutil.GenericNodeFromParticipants(t, hub, colID, identities, chainID)
	defer collectionNode.Done()
	consensusNode := testutil.GenericNodeFromParticipants(t, hub, conID, identities, chainID)
	defer consensusNode.Done()
	exe1Node := testutil.ExecutionNode(t, hub, exe1ID, identities, 27, chainID)
	exe1Node.Ready()
	defer exe1Node.Done()

	genesis, err := exe1Node.State.AtHeight(0).Head()
	require.NoError(t, err)

	seq := uint64(0)

	chain := exe1Node.ChainID.Chain()

	// transaction that will change state and succeed, used to test that state commitment changes
	// genesis <- block1 [tx1] <- block2 [tx2] <- block3 [tx3] <- child
	_, col1, block1, proposal1, seq := deployContractBlock(t, conID, colID, chain, seq, genesis, genesis)

	// we don't set the proper sequence number of this one
	_, col2, block2, proposal2, _ := makePanicBlock(t, conID, colID, chain, uint64(0), block1.Header, genesis)

	_, col3, block3, proposal3, seq := makeSuccessBlock(t, conID, colID, chain, seq, block2.Header, genesis)

	_, _, _, proposal4, _ := makeSuccessBlock(t, conID, colID, chain, seq, block3.Header, genesis)
	// seq++

	// setup mocks and assertions
	collectionEngine := mockCollectionEngineToReturnCollections(
		t,
		&collectionNode,
		[]*flow.Collection{col1, col2, col3},
	)

	receiptsReceived := atomic.Uint64{}

	consensusEngine := new(mocknetwork.Engine)
	_, _ = consensusNode.Net.Register(engine.ReceiveReceipts, consensusEngine)
	consensusEngine.On("Process", mock.AnythingOfType("network.Channel"), mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			receiptsReceived.Inc()
			originID := args[1].(flow.Identifier)
			receipt := args[2].(*flow.ExecutionReceipt)
			finalState, _ := receipt.ExecutionResult.FinalStateCommitment()
			consensusNode.Log.Debug().
				Hex("origin", originID[:]).
				Hex("block", receipt.ExecutionResult.BlockID[:]).
				Hex("final_state_commit", finalState[:]).
				Msg("execution receipt delivered")
		}).Return(nil)

	// submit block2 from consensus node to execution node 1
	err = sendBlock(&exe1Node, conID.NodeID, proposal1)
	require.NoError(t, err)

	err = sendBlock(&exe1Node, conID.NodeID, proposal2)
	assert.NoError(t, err)

	// ensure block 1 has been executed
	hub.DeliverAllEventually(t, func() bool {
		return receiptsReceived.Load() == 1
	})
	exe1Node.AssertHighestExecutedBlock(t, block1.Header)

	scExe1Genesis, err := exe1Node.ExecutionState.StateCommitmentByBlockID(context.Background(), genesis.ID())
	assert.NoError(t, err)

	scExe1Block1, err := exe1Node.ExecutionState.StateCommitmentByBlockID(context.Background(), block1.ID())
	assert.NoError(t, err)
	assert.NotEqual(t, scExe1Genesis, scExe1Block1)

	// submit block 3 and block 4 from consensus node to execution node 1 (who have block1),
	err = sendBlock(&exe1Node, conID.NodeID, proposal3)
	assert.NoError(t, err)

	err = sendBlock(&exe1Node, conID.NodeID, proposal4)
	assert.NoError(t, err)

	// ensure block 1, 2 and 3 have been executed
	hub.DeliverAllEventually(t, func() bool {
		return receiptsReceived.Load() == 3
	})

	// ensure state has been synced across both nodes
	exe1Node.AssertHighestExecutedBlock(t, block3.Header)
	// exe2Node.AssertHighestExecutedBlock(t, block3.Header)

	// verify state commitment of block 2 is the same as block 1, since tx failed on seq number verification
	scExe1Block2, err := exe1Node.ExecutionState.StateCommitmentByBlockID(context.Background(), block2.ID())
	assert.NoError(t, err)
	// TODO this is no longer valid because the system chunk can change the state
	//assert.Equal(t, scExe1Block1, scExe1Block2)
	_ = scExe1Block2

	collectionEngine.AssertExpectations(t)
	consensusEngine.AssertExpectations(t)
}

func mockCollectionEngineToReturnCollections(t *testing.T, collectionNode *testmock.GenericNode, cols []*flow.Collection) *mocknetwork.Engine {
	collectionEngine := new(mocknetwork.Engine)
	colConduit, _ := collectionNode.Net.Register(engine.RequestCollections, collectionEngine)

	// make lookup
	colMap := make(map[flow.Identifier][]byte)
	for _, col := range cols {
		blob, _ := msgpack.Marshal(col)
		colMap[col.ID()] = blob
	}
	collectionEngine.On("Process", mock.AnythingOfType("network.Channel"), mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			originID := args[1].(flow.Identifier)
			req := args[2].(*messages.EntityRequest)
			blob, ok := colMap[req.EntityIDs[0]]
			if !ok {
				assert.FailNow(t, "requesting unexpected collection", req.EntityIDs[0])
			}
			res := &messages.EntityResponse{Blobs: [][]byte{blob}, EntityIDs: req.EntityIDs[:1]}
			err := colConduit.Publish(res, originID)
			assert.NoError(t, err)
		}).
		Return(nil).
		Times(len(cols))
	return collectionEngine
}

// Test the receipt will be sent to multiple verification nodes
func TestBroadcastToMultipleVerificationNodes(t *testing.T) {
	hub := stub.NewNetworkHub()

	chainID := flow.Emulator

	colID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleCollection),
		unittest.WithKeys,
	)
	conID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleConsensus),
		unittest.WithKeys,
	)
	exeID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleExecution),
		unittest.WithKeys,
	)
	ver1ID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleVerification),
		unittest.WithKeys,
	)
	ver2ID := unittest.IdentityFixture(
		unittest.WithRole(flow.RoleVerification),
		unittest.WithKeys,
	)

	identities := unittest.CompleteIdentitySet(colID, conID, exeID, ver1ID, ver2ID)

	exeNode := testutil.ExecutionNode(t, hub, exeID, identities, 21, chainID)
	exeNode.Ready()
	defer exeNode.Done()

	verification1Node := testutil.GenericNodeFromParticipants(t, hub, ver1ID, identities, chainID)
	defer verification1Node.Done()
	verification2Node := testutil.GenericNodeFromParticipants(t, hub, ver2ID, identities, chainID)
	defer verification2Node.Done()

	genesis, err := exeNode.State.AtHeight(0).Head()
	require.NoError(t, err)

	block := unittest.BlockWithParentAndProposerFixture(t, genesis, conID.NodeID)
	block.Header.View = 42
	block.Header.LastViewTC = helper.MakeTC(
		helper.WithTCView(block.Header.View-1),
		helper.WithTCSigners(block.Header.ParentVoterIndices),
		helper.WithTCNewestQC(
			helper.MakeQC(
				helper.WithQCView(genesis.View),
				helper.WithQCSigners(block.Header.ParentVoterIndices))))
	block.SetPayload(flow.Payload{})
	proposal := unittest.ProposalFromBlock(&block)

	child := unittest.BlockWithParentAndProposerFixture(t, block.Header, conID.NodeID)

	actualCalls := atomic.Uint64{}

	verificationEngine := new(mocknetwork.Engine)
	_, _ = verification1Node.Net.Register(engine.ReceiveReceipts, verificationEngine)
	_, _ = verification2Node.Net.Register(engine.ReceiveReceipts, verificationEngine)
	verificationEngine.On("Process", mock.AnythingOfType("network.Channel"), exeID.NodeID, mock.Anything).
		Run(func(args mock.Arguments) {
			actualCalls.Inc()

			var receipt *flow.ExecutionReceipt
			receipt, _ = args[2].(*flow.ExecutionReceipt)

			assert.Equal(t, block.ID(), receipt.ExecutionResult.BlockID)
		}).
		Return(nil)

	err = sendBlock(&exeNode, exeID.NodeID, proposal)
	require.NoError(t, err)

	err = sendBlock(&exeNode, conID.NodeID, unittest.ProposalFromBlock(&child))
	require.NoError(t, err)

	hub.DeliverAllEventually(t, func() bool {
		return actualCalls.Load() == 2
	})

	verificationEngine.AssertExpectations(t)
}

// Test that when received the same state delta for the second time,
// the delta will be saved again without causing any error.
// func TestReceiveTheSameDeltaMultipleTimes(t *testing.T) {
// 	hub := stub.NewNetworkHub()
//
// 	chainID := flow.Mainnet
//
// 	colID := unittest.IdentityFixture(unittest.WithRole(flow.RoleCollection))
// 	exeID := unittest.IdentityFixture(unittest.WithRole(flow.RoleExecution))
// 	ver1ID := unittest.IdentityFixture(unittest.WithRole(flow.RoleVerification))
// 	ver2ID := unittest.IdentityFixture(unittest.WithRole(flow.RoleVerification))
//
// 	identities := unittest.CompleteIdentitySet(colID, exeID, ver1ID, ver2ID)
//
// 	exeNode := testutil.ExecutionNode(t, hub, exeID, identities, 21, chainID)
// 	defer exeNode.Done()
//
// 	genesis, err := exeNode.State.AtHeight(0).Head()
// 	require.NoError(t, err)
//
// 	delta := unittest.StateDeltaWithParentFixture(genesis)
// 	delta.ExecutableBlock.StartState = unittest.GenesisStateCommitment
// 	delta.EndState = unittest.GenesisStateCommitment
//
// 	fmt.Printf("block id: %v, delta for block (%v)'s parent id: %v\n", genesis.ID(), delta.Block.ID(), delta.ParentID())
// 	exeNode.IngestionEngine.SubmitLocal(delta)
// 	time.Sleep(time.Second)
//
// 	exeNode.IngestionEngine.SubmitLocal(delta)
// 	// handling the same delta again to verify the DB calls in saveExecutionResults
// 	// are idempotent, if they weren't, it will hit log.Fatal and crash before
// 	// sleep is done
// 	time.Sleep(time.Second)
//
// }
