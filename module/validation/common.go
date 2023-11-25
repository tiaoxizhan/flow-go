package validation

import (
	"fmt"

	"github.com/onflow/flow-go/engine"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/filter"
	"github.com/onflow/flow-go/state/protocol"
)

// identityForNode ensures that `nodeID` is an authorized member of the network
// at the given block and returns the corresponding node's full identity.
// Error returns:
//   - sentinel engine.InvalidInputError is nodeID is NOT an authorized member of the network
//   - generic error indicating a fatal internal problem
func identityForNode(state protocol.State, blockID flow.Identifier, nodeID flow.Identifier) (*flow.Identity, error) {
	// get the identity of the origin node
	identity, err := state.AtBlockID(blockID).Identity(nodeID)
	if err != nil {
		if protocol.IsIdentityNotFound(err) {
			return nil, engine.NewInvalidInputErrorf("unknown node identity: %w", err)
		}
		// unexpected exception
		return nil, fmt.Errorf("failed to retrieve node identity: %w", err)
	}

	return identity, nil
}

// ensureNodeHasWeightAndRole checks whether, at the given block, `nodeID`
//   - has _positive_ weight
//   - and has the expected role
//   - is an active participant of the current epoch and not ejected (i.e. has `EpochParticipationStatusActive`)
//
// Returns the following errors:
//   - sentinel engine.InvalidInputError if any of the above-listed conditions are violated.
//
// Note: the method receives the identity as proof of its existence.
// Therefore, we consider the case where the respective identity is unknown to the
// protocol state as a symptom of a fatal implementation bug.
func ensureNodeHasWeightAndRole(identity *flow.Identity, expectedRole flow.Role) error {
	// check that the role is expected
	if identity.Role != expectedRole {
		return engine.NewInvalidInputErrorf("expected node %x to have role %s but got %s", identity.NodeID, expectedRole, identity.Role)
	}

	// check if the identity is a valid epoch participant(is active in the current epoch + not ejected)
	if !filter.IsValidCurrentEpochParticipant(identity) {
		return engine.NewInvalidInputErrorf("node (%x) is not an active participant, instead has status: %s", identity.NodeID,
			identity.EpochParticipationStatus.String())
	}

	return nil
}
