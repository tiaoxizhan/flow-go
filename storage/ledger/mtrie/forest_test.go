package mtrie_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go/storage/ledger/mtrie/proof"

	"github.com/dapperlabs/flow-go/storage/ledger/mtrie/trie"

	"github.com/dapperlabs/flow-go/storage/ledger/mtrie"
	"github.com/dapperlabs/flow-go/storage/ledger/mtrie/common"
	cstrie "github.com/dapperlabs/flow-go/storage/ledger/trie"
	"github.com/dapperlabs/flow-go/storage/ledger/utils"
)

// TestTrieOperations tests adding removing and retrieving Trie from Forrest
func TestTrieOperations(t *testing.T) {
	trieHeight := 17

	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	// Add trie
	nt, err := trie.NewEmptyMTrie(trieHeight, 1, []byte{})
	err = fStore.AddTrie(nt)
	require.NoError(t, err)

	// Get trie
	retnt, err := fStore.GetTrie(nt.RootHash())
	require.NoError(t, err)
	require.True(t, bytes.Equal(retnt.RootHash(), nt.RootHash()))
	require.Equal(t, fStore.Size(), 2)

	// Remove trie
	fStore.RemoveTrie(nt.RootHash())
	require.Equal(t, fStore.Size(), 1)
}

// TestTrieUpdate updates the empty trie with some values and verifies that the
// written values can be retrieved from the updated trie.
func TestTrieUpdate(t *testing.T) {
	trieHeight := 17

	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)
	rootHash := fStore.GetEmptyRootHash()

	k1 := []byte([]uint8{uint8(53), uint8(74)})
	v1 := []byte{'A'}
	keys := [][]byte{k1}
	values := [][]byte{v1}
	updatedTrie, err := fStore.Update(rootHash, keys, values)
	require.NoError(t, err)

	retValues, err := fStore.Read(updatedTrie.RootHash(), keys)
	require.NoError(t, err)
	require.True(t, bytes.Equal(retValues[0], values[0]))
}

// TestLeftEmptyInsert tests inserting a new value into an empty sub-trie:
//   1. we first construct a baseTrie holding a couple of values on the right branch [~]
//   2. we update a previously non-existent register on the left branch (X)
// We verify that values for _all_ keys in the updated Trie have correct values
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestLeftEmptyInsert(t *testing.T) {
	//////////////////////
	//     insert X     //
	//       ()         //
	//      /  \        //
	//    (X)  [~]      //
	//////////////////////
	trieHeight := 17
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(129), uint8(1)}) // key: 1000...
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(193), uint8(1)}) // key: 1100...
	v2 := []byte{'B'}
	keys := [][]byte{k1, k2}
	values := [][]byte{v1, v2}

	baseTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)
	// resulting base trie:
	// 16: ([],)[]
	// 		15: ([],)[1]
	// 			14: ([129 1],41)[10]
	// 			14: ([193 1],42)[11]
	fmt.Println("BASE TRIE:")
	fmt.Println(baseTrie.String())

	k3 := []byte([]uint8{uint8(1), uint8(1)})
	v3 := []byte{'C'}
	keys = [][]byte{k3}
	values = [][]byte{v3}
	updatedTrie, err := fStore.Update(baseTrie.RootHash(), keys, values)
	require.NoError(t, err)
	// expected updated Trie:
	// TODO: update Trie representation
	// 16: ([],)[]
	// 		15: ([1 1],43)[0]
	// 		15: ([],)[1]
	// 			14: ([129 1],41)[10]
	// 			14: ([193 1],42)[11]
	fmt.Println("UPDATED TRIE:")
	fmt.Println(updatedTrie.String())

	keys = [][]byte{k1, k2, k3}
	values = [][]byte{v1, v2, v3}
	retValues, err := fStore.Read(updatedTrie.RootHash(), keys)
	require.NoError(t, err)
	for i := range keys {
		require.True(t, bytes.Equal(retValues[i], values[i]))
	}
}

// TestLeftEmptyInsert tests inserting a new value into an empty sub-trie:
//   1. we first construct a baseTrie holding a couple of values on the left branch [~]
//   2. we update a previously non-existent register on the right branch (X)
// We verify that values for _all_ keys in the updated Trie have correct values
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestRightEmptyInsert(t *testing.T) {
	///////////////////////
	//     insert X      //
	//       ()          //
	//      /  \         //
	//    [~]  (X)       //
	///////////////////////
	trieHeight := 17
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(1), uint8(1)}) // key: 0000...
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(64), uint8(1)}) // key: 0100....
	v2 := []byte{'B'}
	keys := [][]byte{k1, k2}
	values := [][]byte{v1, v2}

	baseTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)
	// resulting base trie:
	// TODO: update Trie representation
	// 16: ([],)[]
	// 		15: ([],)[0]
	// 			14: ([ 1 1],41)[00]
	// 			14: ([64 1],42)[01]
	fmt.Println("BASE TRIE:")
	fmt.Println(baseTrie.String())

	k3 := []byte([]uint8{uint8(129), uint8(1)})
	v3 := []byte{'C'}
	keys = [][]byte{k3}
	values = [][]byte{v3}
	updatedTrie, err := fStore.Update(baseTrie.RootHash(), keys, values)
	require.NoError(t, err)
	// expected updated Trie:
	// TODO: update Trie representation
	// 16: ([],)[]
	// 		15: ([],)[0]
	// 			14: ([ 1 1],41)[00]
	// 			14: ([64 1],42)[01]
	// 		15: ([129 1],43)[1]
	fmt.Println("UPDATED TRIE:")
	fmt.Println(updatedTrie.String())

	keys = [][]byte{k1, k2, k3}
	values = [][]byte{v1, v2, v3}
	retValues, err := fStore.Read(updatedTrie.RootHash(), keys)
	require.NoError(t, err)
	for i := range keys {
		require.True(t, bytes.Equal(retValues[i], values[i]))
	}
}

// TestExpansionInsert tests inserting a new value into a populated sub-trie, where a
// leaf (holding a single value) would be replaced by an expanded sub-trie holding multiple value
//   1. we first construct a baseTrie holding a couple of values on the right branch [~]
//   2. we update a previously non-existent register on the right branch turning [~] to [~']
// We verify that values for _all_ keys in the updated Trie have correct values
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestExpansionInsert(t *testing.T) {
	////////////////////////
	// modify [~] -> [~'] //
	//       ()           //
	//      /  \          //
	//         [~]        //
	////////////////////////

	trieHeight := 17
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(129), uint8(1)}) // key: 1000000...
	v1 := []byte{'A'}
	keys := [][]byte{k1}
	values := [][]byte{v1}
	baseTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)
	// resulting base trie:
	// 16: ([129 1],41)[]
	fmt.Println("BASE TRIE:")
	fmt.Println(baseTrie.String())

	k2 := []byte([]uint8{uint8(130), uint8(1)}) // key: 1000001...
	v2 := []byte{'B'}
	keys = [][]byte{k2}
	values = [][]byte{v2}
	updatedTrie, err := fStore.Update(baseTrie.RootHash(), keys, values)
	require.NoError(t, err)
	// expected updated Trie:
	// 16: ([],)[]
	// 		15: ([],)[1]
	// 			14: ([],)[10]
	// 				13: ([],)[100]
	// 					12: ([],)[1000]
	// 						11: ([],)[10000]
	// 							10: ([],)[100000]
	// 								9: ([129 1],41)[1000000]
	// 								9: ([130 1],42)[1000001]
	fmt.Println("UPDATED TRIE:")
	fmt.Println(baseTrie.String())

	keys = [][]byte{k1, k2}
	values = [][]byte{v1, v2}
	retValues, err := fStore.Read(updatedTrie.RootHash(), keys)
	require.NoError(t, err)
	for i := range keys {
		require.True(t, bytes.Equal(retValues[i], values[i]))
	}
}

// TestFullHouseInsert tests inserting a new value into a populated sub-trie, where a
// leaf's value is overridden _and_ further values are added which all fall into a subtree that
// replaces the leaf:
//   1. we first construct a baseTrie holding a couple of values on the right branch [~]
//   2. we update a previously non-existent register on the right branch turning [~] to [~']
// We verify that values for _all_ keys in the updated Trie have correct values
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestFullHouseInsert(t *testing.T) {
	///////////////////////
	//   insert ~1<X<~2  //
	//       ()          //
	//      /  \         //
	//    [~1]  [~2]     //
	///////////////////////

	trieHeight := 17
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	// key-value pair (k0,v0) forms [~1]; (k1,v1) and (k2,v2) form [~2]
	k0 := []byte([]uint8{uint8(64), uint8(1)}) // key: 0100...
	v0 := []byte{'0'}
	k1 := []byte([]uint8{uint8(129), uint8(1)}) // key: 1000...
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(193), uint8(1)}) // key: 1100....
	v2 := []byte{'B'}
	keys := [][]byte{k0, k1, k2}
	values := [][]byte{v0, v1, v2}

	baseTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)
	// expected trie:
	// TODO: update Trie representation
	// 16: ([],)[]
	// 		15: ([],)[1]
	// 			14: ([129 1],41)[10]
	// 			14: ([193 1],42)[11]
	fmt.Println("BASE TRIE:")
	fmt.Println(baseTrie.String())

	// we update value for key k1 and in addition add key-value (k3,v3) pair that has the same prefix `10` as k1
	v1 = []byte{'X'}
	k3 := []byte([]uint8{uint8(160), uint8(1)}) // key: 1010...
	v3 := []byte{'C'}
	keys = [][]byte{k1, k3}
	values = [][]byte{v1, v3}
	updatedTrie, err := fStore.Update(baseTrie.RootHash(), keys, values)
	require.NoError(t, err)
	// expected trie:
	// TODO: update Trie representation
	// 16: ([],)[]
	// 		15: ([],)[1]
	// 			14: ([],)[10]
	// 				13: ([129 1],41)[100]
	// 				13: ([160 1],43)[101]
	// 			14: ([193 1],42)[11]
	fmt.Println("UPDATED TRIE:")
	fmt.Println(baseTrie.String())

	keys = [][]byte{k1, k2, k3}
	values = [][]byte{v1, v2, v3}
	retValues, err := fStore.Read(updatedTrie.RootHash(), keys)
	require.NoError(t, err)
	for i := range keys {
		require.True(t, bytes.Equal(retValues[i], values[i]))
	}
}

// TestLeafInsert inserts two keys, which only differ in their last bit.
// I.e. the trie needs to be expanded to its hull depth
// We verify that values for _all_ keys in the updated Trie have correct values
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestLeafInsert(t *testing.T) {
	///////////////////////
	//   insert 1, 2     //
	//       ()          //
	//      /  \         //
	//     ()   ...      //
	//          /  \     //
	//         ()  ()    //
	///////////////////////
	trieHeight := 17
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(1), uint8(0)})
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(1), uint8(1)})
	v2 := []byte{'B'}
	keys := [][]byte{k1, k2}
	values := [][]byte{v1, v2}

	testTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)
	// expected trie:
	// 16: ([],)[]
	// 		15: ([],)[0]
	// 			14: ([],)[00]
	//               ...
	// 					1: ([],)[000000010000000]
	// 						0: ([1 0],41)[0000000100000000]
	// 						0: ([1 1],42)[0000000100000001]

	retValues, err := fStore.Read(testTrie.RootHash(), keys)
	require.NoError(t, err)
	for i := range keys {
		require.True(t, bytes.Equal(retValues[i], values[i]))
	}
}

// TestOverrideValue overrides an existing value in the trie (without any expansion)
// We verify that values for _all_ keys in the updated Trie have correct values
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestOverrideValue(t *testing.T) {
	trieHeight := 17 // should be key size (in bits) + 1
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(53), uint8(74)})
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(116), uint8(129)})
	v2 := []byte{'B'}
	keys := [][]byte{k1, k2}
	values := [][]byte{v1, v2}

	baseTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)

	k3 := []byte([]uint8{uint8(53), uint8(74)})
	v3 := []byte{'C'}
	keys = [][]byte{k3}
	values = [][]byte{v3}
	updatedTrie, err := fStore.Update(baseTrie.RootHash(), keys, values)
	require.NoError(t, err)

	retValues, err := fStore.Read(updatedTrie.RootHash(), keys)
	require.NoError(t, err)
	require.True(t, bytes.Equal(retValues[0], values[0]))
}

// TestDuplicateOverride tests behaviour when the updates contain two different values for the
// same key. I.e. we update with (k0, v0) and (k0, v1)
// We expect that the _last_ written value is persisted in the Trie
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestDuplicateOverride(t *testing.T) {
	trieHeight := 17 // should be key size (in bits) + 1
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k0 := []byte([]uint8{uint8(53), uint8(74)})
	v0 := []byte{'A'}
	baseTrie, err := fStore.Update(fStore.GetEmptyRootHash(), [][]byte{k0}, [][]byte{v0})
	require.NoError(t, err)

	v1 := []byte{'B'}
	v2 := []byte{'C'}
	updatedTrie, err := fStore.Update(baseTrie.RootHash(), [][]byte{k0, k0}, [][]byte{v1, v2})
	require.NoError(t, err)

	retValues, err := fStore.Read(updatedTrie.RootHash(), [][]byte{k0})
	require.NoError(t, err)
	require.Equal(t, retValues, [][]byte{v2})
}

// TestUpdateWithWrongKeySize verifies that attempting to update a trie with wrong key size errors
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestUpdateWithWrongKeySize(t *testing.T) {
	trieHeight := 17 // should be key size (in bits) + 1
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	// short key
	key1 := make([]byte, 1)
	utils.SetBit(key1, 5)
	value1 := []byte{'a'}
	keys := [][]byte{key1}
	values := [][]byte{value1}

	testTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.Error(t, err)

	// long key
	key2 := make([]byte, 33)
	utils.SetBit(key2, 5)
	value2 := []byte{'a'}
	keys = [][]byte{key2}
	values = [][]byte{value2}

	_, err = fStore.Update(testTrie.RootHash(), keys, values)
	require.Error(t, err)
}

// TestReadOrder tests that values from reading a trie are delivered in the order as specified by the keys
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestReadOrder(t *testing.T) {
	trieHeight := 17 // should be key size (in bits) + 1
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(116), uint8(74)})
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(53), uint8(129)})
	v2 := []byte{'B'}
	keys := [][]byte{k1, k2}
	values := [][]byte{v1, v2}
	testTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)

	retValues, err := fStore.Read(testTrie.RootHash(), [][]byte{k1, k2})
	require.NoError(t, err)
	require.True(t, bytes.Equal(retValues[0], values[0]))
	require.True(t, bytes.Equal(retValues[1], values[1]))

	retValues, err = fStore.Read(testTrie.RootHash(), [][]byte{k2, k1})
	require.NoError(t, err)
	require.True(t, bytes.Equal(retValues[1], values[0]))
	require.True(t, bytes.Equal(retValues[0], values[1]))
}

// TestMixRead tests reading a mixture of set and unset registers.
// We expect the default value (empty slice) to be returned for unset registers.
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestMixRead(t *testing.T) {
	trieHeight := 17 // should be key size (in bits) + 1
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(125), uint8(23)}) // key: 01111101...
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(178), uint8(152)}) // key: 10110010...
	v2 := []byte{'B'}
	keys := [][]byte{k1, k2}
	values := [][]byte{v1, v2}

	baseTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)

	k3 := []byte([]uint8{uint8(110), uint8(48)}) // key: 01101110...
	v3 := []byte{}
	k4 := []byte([]uint8{uint8(23), uint8(82)}) // key: 00010111...
	v4 := []byte{}
	readKeys := [][]byte{k1, k2, k3, k4}
	expectedValues := [][]byte{v1, v2, v3, v4}

	retValues, err := fStore.Read(baseTrie.RootHash(), readKeys)
	require.NoError(t, err)
	require.Equal(t, retValues, expectedValues)
	// TODO: cleanup commented out code
	//for i := range keys {
	//	require.True(t, bytes.Equal(retValues[i], expectedValues[i]))
	//}
}

// TestReadWithDuplicatedKeys reads a the values for two keys, where both keys have the same value.
// We expect that we receive the respective value twice in the return.
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestReadWithDuplicatedKeys(t *testing.T) {
	trieHeight := 17 // should be key size (in bits) + 1
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(53), uint8(74)})
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(116), uint8(129)})
	v2 := []byte{'B'}
	k3 := []byte([]uint8{uint8(53), uint8(74)})

	keys := [][]byte{k1, k2}
	values := [][]byte{v1, v2}
	testTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)

	keys = [][]byte{k1, k2, k3}
	expectedValues := [][]byte{v1, v2, v1}

	retValues, err := fStore.Read(testTrie.RootHash(), keys)
	require.NoError(t, err)
	require.Equal(t, retValues, expectedValues)
	// TODO: cleanup commented out code
	//require.True(t, bytes.Equal(retValues[0], values[0]))
	//require.True(t, bytes.Equal(retValues[1], values[1]))
	//require.True(t, bytes.Equal(retValues[2], values[2]))
}

// TestReadNonExistKey tests reading an unset registers.
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestReadNonExistKey(t *testing.T) {
	trieHeight := 17 // should be key size (in bits) + 1
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(53), uint8(74)})
	v1 := []byte{'A'}
	keys := [][]byte{k1}
	values := [][]byte{v1}
	testTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)

	k2 := []byte([]uint8{uint8(116), uint8(129)})
	retValues, err := fStore.Read(testTrie.RootHash(), [][]byte{k2})
	require.NoError(t, err)
	require.Equal(t, retValues, [][]byte{[]byte{}})
	// TODO: cleanup commented out code
	//require.Equal(t, len(retValues[0]), 0)
}

// TestReadWithWrongKeySize verifies that attempting to read a trie with wrong key size errors
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestReadWithWrongKeySize(t *testing.T) {
	trieHeight := 17 // should be key size (in bits) + 1
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	// setup
	key1 := make([]byte, 2)
	utils.SetBit(key1, 5)
	value1 := []byte{'a'}
	keys := [][]byte{key1}
	values := [][]byte{value1}
	testTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)

	// key too short
	key2 := make([]byte, 1)
	utils.SetBit(key2, 5)
	keys = [][]byte{key2}
	_, err = fStore.Read(testTrie.RootHash(), keys)
	require.Error(t, err)

	// key too long
	key3 := make([]byte, 33)
	utils.SetBit(key3, 5)
	keys = [][]byte{key3}
	_, err = fStore.Read(testTrie.RootHash(), keys)
	require.Error(t, err)
}

// TODO test read (multiple non exist in a branch)
// [AlexH] doesn't TestMixRead do this test?

// TestForkingUpdates updates a base trie in two different ways. We expect
// that for each update, a new trie is added to the forrest preserving the
// updated values independently of the other update.
func TestForkingUpdates(t *testing.T) {
	trieHeight := 17 // should be key size (in bits) + 1
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(53), uint8(74)})
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(116), uint8(129)})
	v2 := []byte{'B'}
	keys := [][]byte{k1, k2}
	values := [][]byte{v1, v2}
	baseTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)

	// update baseTrie -> updatedTrieA
	v1a := []byte{'C'}
	k3a := []byte([]uint8{uint8(116), uint8(22)})
	v3a := []byte{'D'}
	keysA := [][]byte{k1, k3a}
	valuesA := [][]byte{v1a, v3a}
	updatedTrieA, err := fStore.Update(baseTrie.RootHash(), keysA, valuesA)
	require.NoError(t, err)

	// update baseTrie -> updatedTrieB
	v1b := []byte{'C'}
	k3b := []byte([]uint8{uint8(116), uint8(22)})
	v3b := []byte{'D'}
	keysB := [][]byte{k1, k3b}
	valuesB := [][]byte{v1b, v3b}
	updatedTrieB, err := fStore.Update(baseTrie.RootHash(), keysB, valuesB)
	require.NoError(t, err)

	// Verify values preserved
	retValues, err := fStore.Read(baseTrie.RootHash(), keys) // reading from original Trie
	require.NoError(t, err)
	require.Equal(t, retValues, values)

	retValues, err = fStore.Read(updatedTrieA.RootHash(), keysA) // reading from updatedTrieA
	require.NoError(t, err)
	require.Equal(t, retValues, valuesA)

	retValues, err = fStore.Read(updatedTrieB.RootHash(), keysB) // reading from updatedTrieB
	require.NoError(t, err)
	require.Equal(t, retValues, valuesB)

	// TODO: cleanup commented out code
	//keys = [][]byte{k1, k2, k3}
	//retValues, err := fStore.Read(keys, rootHash21)
	//require.NoError(t, err)
	//require.True(t, bytes.Equal(retValues[0], v1))
	//require.True(t, bytes.Equal(retValues[1], v2))
	//require.True(t, bytes.Equal(retValues[2], []byte{}))
	//
	//retValues, err = fStore.Read(keys, rootHash22)
	//require.NoError(t, err)
	//require.True(t, bytes.Equal(retValues[0], v1p))
	//require.True(t, bytes.Equal(retValues[1], []byte{}))
	//require.True(t, bytes.Equal(retValues[2], v3))
}

// TestIdenticalUpdateAppliedTwice updates a base trie in the same way twice.
// Hence, the forest should de-duplicate the resulting two version of the identical trie
// without an error.
func TestIdenticalUpdateAppliedTwice(t *testing.T) {
	trieHeight := 17 // should be key size (in bits) + 1
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(53), uint8(74)})
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(116), uint8(129)})
	v2 := []byte{'B'}
	keys := [][]byte{k1, k2}
	values := [][]byte{v1, v2}
	baseTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)

	k3 := []byte([]uint8{uint8(116), uint8(22)})
	v3 := []byte{'D'}
	updatedTrieA, err := fStore.Update(baseTrie.RootHash(), [][]byte{k3}, [][]byte{v3})
	require.NoError(t, err)
	updatedTrieB, err := fStore.Update(baseTrie.RootHash(), [][]byte{k3}, [][]byte{v3})
	require.NoError(t, err)
	require.Equal(t, updatedTrieA.RootHash(), updatedTrieB.RootHash())

	retValuesA, err := fStore.Read(updatedTrieA.RootHash(), [][]byte{k1, k2, k3})
	require.NoError(t, err)
	require.Equal(t, retValuesA, [][]byte{v1, v2, v3})

	retValuesB, err := fStore.Read(updatedTrieA.RootHash(), [][]byte{k1, k2, k3})
	require.NoError(t, err)
	require.Equal(t, retValuesB, [][]byte{v1, v2, v3})
}

// TestRandomUpdateReadProof tests a read proof against the ChainSafe Trie implementation
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestRandomUpdateReadProof(t *testing.T) {
	keyByteSize := 2
	trieHeight := keyByteSize*8 + 1
	rep := 10
	maxNumKeysPerStep := 10
	maxValueSize := 6 // bytes
	rand.Seed(time.Now().UnixNano())
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir) // clean up

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)
	emptyTrieHash := fStore.GetEmptyRootHash()
	latestValueByKey := make(map[string][]byte) // map store

	for e := 0; e < rep; e++ {
		keys := common.GetRandomKeysRandN(maxNumKeysPerStep, keyByteSize)
		values := common.GetRandomValues(len(keys), maxValueSize)

		// update map store with key values
		// we use this at the end of each step to check all existing keys
		for i, k := range keys {
			latestValueByKey[string(k)] = values[i]
		}

		// test reading for non-existing keys
		nonExistingKeys := make([][]byte, 0)
		for _, k := range keys {
			if _, ok := latestValueByKey[string(k)]; !ok {
				nonExistingKeys = append(nonExistingKeys, k)
			}
		}
		retValues, err := fStore.Read(emptyTrieHash, nonExistingKeys)
		require.NoError(t, err, "error reading - non existing keys")
		for i := range retValues {
			require.True(t, len(retValues[i]) == 0)
		}

		// test update
		updatedTrie, err := fStore.Update(emptyTrieHash, keys, values)
		updatedTrieHash := updatedTrie.RootHash()
		require.NoError(t, err, "error updating")

		// test read
		retValues, err = fStore.Read(updatedTrieHash, keys)
		require.NoError(t, err, "error reading")
		for i := range values {
			require.True(t, bytes.Equal(values[i], retValues[i]))
		}

		// test proof (mix of existing and non existing keys)
		proofKeys := make([][]byte, 0)
		proofValues := make([][]byte, 0)
		for i, k := range keys {
			proofKeys = append(proofKeys, k)
			proofValues = append(proofValues, values[i])
		}

		for _, k := range nonExistingKeys {
			proofKeys = append(proofKeys, k)
			proofValues = append(proofValues, []byte{})
		}

		batchProof, err := fStore.Proofs(updatedTrieHash, proofKeys)
		require.NoError(t, err, "error generating proofs")
		require.True(t, batchProof.Verify(proofKeys, proofValues, updatedTrieHash, trieHeight))

		psmt, err := cstrie.NewPSMT(updatedTrieHash, trieHeight, proofKeys, proofValues, proof.EncodeBatchProof(batchProof))
		require.NoError(t, err, "error building partial trie")
		require.True(t, bytes.Equal(psmt.GetRootHash(), updatedTrieHash))

		// check values for all existing keys
		allKeys := make([][]byte, 0, len(latestValueByKey))
		allValues := make([][]byte, 0, len(latestValueByKey))
		for k, v := range latestValueByKey {
			allKeys = append(allKeys, []byte(k))
			allValues = append(allValues, v)
		}
		retValues, err = fStore.Read(updatedTrieHash, allKeys)
		require.NoError(t, err)
		for i, v := range allValues {
			require.True(t, bytes.Equal(v, retValues[i]))
		}
	}
}

// TestProofGenerationInclusion tests that inclusion proofs generated by a Trie pass verification
// TODO: move to Trie test (as it directly tests trie update as opposed to forrest functions)
func TestProofGenerationInclusion(t *testing.T) {
	trieHeight := 17
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)
	emptyTrieHash := fStore.GetEmptyRootHash()

	k1 := []byte([]uint8{uint8(1), uint8(74)})
	v1 := []byte{'A'}

	k2 := []byte([]uint8{uint8(2), uint8(74)})
	v2 := []byte{'B'}

	k3 := []byte([]uint8{uint8(130), uint8(74)})
	v3 := []byte{'C'}

	k4 := []byte([]uint8{uint8(131), uint8(74)})
	v4 := []byte{'D'}

	keys := [][]byte{k1, k2, k3, k4}
	values := [][]byte{v1, v2, v3, v4}
	updatedTrie, err := fStore.Update(emptyTrieHash, keys, values)
	require.NoError(t, err)
	proof, err := fStore.Proofs(updatedTrie.RootHash(), keys)
	require.NoError(t, err)
	require.True(t, proof.Verify(keys, values, updatedTrie.RootHash(), trieHeight))
}

// TestPurgeAndLoad this test updates the Tries repeatedly until the forest should
// purge some Tries from memory. Thereafter, we read the trie that was purged, hence
// making the forest load the trie again.
// TODO: implement this functionality for MForest
func TestPurgeAndLoad(t *testing.T) {
	t.Skip("we don't have this functionality right now")
	keyByteSize := 2
	trieHeight := keyByteSize*8 + 1

	dir, err := ioutil.TempDir("", "test-mtrie")
	require.NoError(t, err)
	defer os.RemoveAll(dir) // clean up

	fStore, err := mtrie.NewMForest(trieHeight, dir, 2, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(1), uint8(74)})
	v1 := []byte{'A'}

	keys := [][]byte{k1}
	values := [][]byte{v1}
	updatedTrie1, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)

	k2 := []byte([]uint8{uint8(2), uint8(74)})
	v2 := []byte{'B'}
	keys = [][]byte{k2}
	values = [][]byte{v2}
	updatedTrie2, err := fStore.Update(updatedTrie1.RootHash(), keys, values)
	require.NoError(t, err)

	k3 := []byte([]uint8{uint8(130), uint8(74)})
	v3 := []byte{'C'}
	keys = [][]byte{k3}
	values = [][]byte{v3}
	updatedTrie3, err := fStore.Update(updatedTrie2.RootHash(), keys, values)
	require.NoError(t, err)

	k4 := []byte([]uint8{uint8(131), uint8(74)})
	v4 := []byte{'D'}
	keys = [][]byte{k4}
	values = [][]byte{v4}
	updatedTrie4, err := fStore.Update(updatedTrie3.RootHash(), keys, values)
	require.NoError(t, err)

	keys = [][]byte{k1}
	values = [][]byte{v1}
	retValues, err := fStore.Read(updatedTrie1.RootHash(), keys)
	require.NoError(t, err)
	require.True(t, bytes.Equal(values[0], retValues[0]))

	retValues, err = fStore.Read(updatedTrie2.RootHash(), keys)
	require.NoError(t, err)
	require.True(t, bytes.Equal(values[0], retValues[0]))

	retValues, err = fStore.Read(updatedTrie3.RootHash(), keys)
	require.NoError(t, err)
	require.True(t, bytes.Equal(values[0], retValues[0]))

	retValues, err = fStore.Read(updatedTrie4.RootHash(), keys)
	require.NoError(t, err)
	require.True(t, bytes.Equal(values[0], retValues[0]))
}

// TestTrieStoreAndLoad tests storing a trie to file and loading it again.
// We verify that that the values in the loaded trie are correct and that it has
// the expected root hash.
// TODO: implement this functionality for MForest
func TestTrieStoreAndLoad(t *testing.T) {
	trieHeight := 17
	dir, err := ioutil.TempDir("", "test-mtrie-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fStore, err := mtrie.NewMForest(trieHeight, dir, 5, nil)
	require.NoError(t, err)

	k1 := []byte([]uint8{uint8(1), uint8(74)})
	v1 := []byte{'A'}
	k2 := []byte([]uint8{uint8(2), uint8(74)})
	v2 := []byte{'B'}
	k3 := []byte([]uint8{uint8(130), uint8(74)})
	v3 := []byte{'C'}
	k4 := []byte([]uint8{uint8(131), uint8(74)})
	v4 := []byte{'D'}
	k5 := []byte([]uint8{uint8(132), uint8(74)})
	v5 := []byte{'E'}

	keys := [][]byte{k1, k2, k3, k4, k5}
	values := [][]byte{v1, v2, v3, v4, v5}
	testTrie, err := fStore.Update(fStore.GetEmptyRootHash(), keys, values)
	require.NoError(t, err)

	file, err := ioutil.TempFile("", "flow-mtrie-load")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	err = fStore.StoreTrie(testTrie.RootHash(), file.Name())
	require.NoError(t, err)

	// create new store
	newStore, err := mtrie.NewMForest(trieHeight, "", 5, nil)
	require.NoError(t, err)
	loadedTrie, err := newStore.LoadTrie(file.Name())
	require.NoError(t, err)
	require.Equal(t, loadedTrie.RootHash(), testTrie.RootHash())

	retValues, err := newStore.Read(loadedTrie.RootHash(), keys)
	require.NoError(t, err)
	require.Equal(t, retValues, values)
	//for i := range keys {
	//	require.True(t, bytes.Equal(values[i], retValues[i]))
	//}
}

// TODO comment out this for now
// func TestMForestAccuracy(t *testing.T) {
// 	trieHeight := 17 // should be key size (in bits) + 1
// 	experimentRep := 10

// 	dbDir := unittest.TempDir(t)
// 	smt, err := cstrie.NewSMT(dbDir, trieHeight, 10, 100, experimentRep)
// 	require.NoError(t, err)
// 	defer func() {
// 		smt.SafeClose()
// 		os.RemoveAll(dbDir)
// 	}()

// 	fStore, err := mtrie.NewMForest(trieHeight, dbDir, 5, nil)
// 	require.NoError(t, err)
// 	rootHash := fStore.GetEmptyRootHash()

// 	emptyTree := cstrie.GetDefaultHashForHeight(trieHeight - 1)
// 	require.NoError(t, err)
// 	rootHashForSMT := emptyTree
// 	for e := 0; e < experimentRep; e++ {
// 		// insert some values to an empty trie
// 		keys := make([][]byte, 0)
// 		values := make([][]byte, 0)
// 		rand.Seed(time.Now().UnixNano())

// 		// rejection sampling
// 		numberOfKeys := rand.Intn(20) + 1
// 		keyValueMap := make(map[string][]byte)
// 		i := 0
// 		for i < numberOfKeys {
// 			key := make([]byte, 2)
// 			rand.Read(key)
// 			// deduplicate
// 			if _, found := keyValueMap[string(key)]; !found {
// 				keys = append(keys, key)
// 				value := make([]byte, 4)
// 				rand.Read(value)
// 				keyValueMap[string(key)] = value
// 				values = append(values, value)
// 				i++
// 			}
// 		}

// 		newRootHash, err := fStore.Update(keys, values, rootHash)
// 		require.NoError(t, err, "error commiting changes")
// 		rootHash = newRootHash

// 		// check values
// 		retValues, err := fStore.Read(keys, rootHash)
// 		require.NoError(t, err)
// 		for i, k := range keys {
// 			require.True(t, bytes.Equal(keyValueMap[string(k)], retValues[i]))
// 		}

// 		// Test eqaulity to SMT
// 		newRootHashForSMT, err := smt.Update(keys, values, rootHashForSMT)
// 		require.NoError(t, err)
// 		rootHashForSMT = newRootHashForSMT
// 		require.True(t, bytes.Equal(newRootHashForSMT, newRootHash))

// 		// TODO test proofs for non-existing keys
// 		batchProof, err := fStore.Proofs(keys, rootHash)
// 		require.NoError(t, err, "error generating proofs")

// 		batchProofSMT, err := smt.GetBatchProof(keys, rootHashForSMT)
// 		require.NoError(t, err, "error generating proofs (SMT)")

// 		encodedProof := mtrie.EncodeBatchProof(batchProof)
// 		encodedProofSMT := cstrie.EncodeProof(batchProofSMT)

// 		for i := range encodedProof {
// 			require.True(t, bytes.Equal(encodedProof[i], encodedProofSMT[i]))
// 		}

// 		psmt, err := cstrie.NewPSMT(rootHash, trieHeight, keys, values, encodedProof)
// 		require.True(t, bytes.Equal(psmt.GetRootHash(), rootHash))
// 		require.NoError(t, err, "error building partial trie")

// 	}
// }
