package unicastcache_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/module/metrics"
	p2ptest "github.com/onflow/flow-go/network/p2p/test"
	unicastcache "github.com/onflow/flow-go/network/p2p/unicast/cache"
	"github.com/onflow/flow-go/network/p2p/unicast/unicastmodel"
	"github.com/onflow/flow-go/utils/unittest"
)

// TestNewDialConfigCache tests the creation of a new DialConfigCache.
// It asserts that the cache is created and its size is 0.
func TestNewDialConfigCache(t *testing.T) {
	sizeLimit := uint32(100)
	logger := zerolog.Nop()
	collector := metrics.NewNoopCollector()
	cache := unicastcache.NewDialConfigCache(sizeLimit, logger, collector, dialConfigFixture)
	require.NotNil(t, cache)
	require.Equalf(t, uint(0), cache.Size(), "cache size must be 0")
}

// dialConfigFixture returns a dial config fixture.
// The dial config is initialized with the default values.
func dialConfigFixture() unicastmodel.DialConfig {
	return unicastmodel.DialConfig{
		DialAttemptBudget:           unicastmodel.MaxDialAttemptTimes,
		StreamCreationAttemptBudget: unicastmodel.MaxStreamCreationAttemptTimes,
	}
}

// TestDialConfigCache_Adjust tests the Adjust method of the DialConfigCache. It asserts that the dial config is initialized, adjusted,
// and stored in the cache.
func TestDialConfigCache_Adjust_Init(t *testing.T) {
	sizeLimit := uint32(100)
	logger := zerolog.Nop()
	collector := metrics.NewNoopCollector()

	dialFactoryCalled := 0
	dialConfigFactory := func() unicastmodel.DialConfig {
		require.Less(t, dialFactoryCalled, 2, "dial config factory must be called at most twice")
		dialFactoryCalled++
		return dialConfigFixture()
	}
	adjustFuncIncrement := func(cfg unicastmodel.DialConfig) (unicastmodel.DialConfig, error) {
		cfg.DialAttemptBudget++
		return cfg, nil
	}

	cache := unicastcache.NewDialConfigCache(sizeLimit, logger, collector, dialConfigFactory)
	require.NotNil(t, cache)
	require.Zerof(t, cache.Size(), "cache size must be 0")

	peerID1 := p2ptest.PeerIdFixture(t)
	peerID2 := p2ptest.PeerIdFixture(t)

	// Initializing the dial config for peerID1 through GetOrInit.
	// dial config for peerID1 does not exist in the cache, so it must be initialized when using GetOrInit.
	cfg, err := cache.GetOrInit(peerID1)
	require.NoError(t, err)
	require.NotNil(t, cfg, "dial config must not be nil")
	require.Equal(t, dialConfigFixture(), *cfg, "dial config must be initialized with the default values")
	require.Equal(t, uint(1), cache.Size(), "cache size must be 1")

	// Initializing and adjusting the dial config for peerID2 through Adjust.
	// dial config for peerID2 does not exist in the cache, so it must be initialized when using Adjust.
	cfg, err = cache.Adjust(peerID2, adjustFuncIncrement)
	require.NoError(t, err)
	// adjusting a non-existing dial config must not initialize the config.
	require.Equal(t, uint(2), cache.Size(), "cache size must be 2")
	require.Equal(t, cfg.LastSuccessfulDial, dialConfigFixture().LastSuccessfulDial, "last successful dial must be 0")
	require.Equal(t, cfg.DialAttemptBudget, dialConfigFixture().DialAttemptBudget+1, "dial backoff must be adjusted")
	require.Equal(t, cfg.StreamCreationAttemptBudget, dialConfigFixture().StreamCreationAttemptBudget, "stream backoff must be 1")

	// Retrieving the dial config of peerID2 through GetOrInit.
	// retrieve the dial config for peerID2 and assert than it is initialized with the default values; and the adjust function is applied.
	cfg, err = cache.GetOrInit(peerID2)
	require.NoError(t, err, "dial config must exist in the cache")
	require.NotNil(t, cfg, "dial config must not be nil")
	// retrieving an existing dial config must not change the cache size.
	require.Equal(t, uint(2), cache.Size(), "cache size must be 2")
	// config should be the same as the one returned by Adjust.
	require.Equal(t, cfg.LastSuccessfulDial, dialConfigFixture().LastSuccessfulDial, "last successful dial must be 0")
	require.Equal(t, cfg.DialAttemptBudget, dialConfigFixture().DialAttemptBudget+1, "dial backoff must be adjusted")
	require.Equal(t, cfg.StreamCreationAttemptBudget, dialConfigFixture().StreamCreationAttemptBudget, "stream backoff must be 1")

	// Adjusting the dial config of peerID1 through Adjust.
	// dial config for peerID1 already exists in the cache, so it must be adjusted when using Adjust.
	cfg, err = cache.Adjust(peerID1, adjustFuncIncrement)
	require.NoError(t, err)
	// adjusting an existing dial config must not change the cache size.
	require.Equal(t, uint(2), cache.Size(), "cache size must be 2")
	require.Equal(t, cfg.LastSuccessfulDial, dialConfigFixture().LastSuccessfulDial, "last successful dial must be 0")
	require.Equal(t, cfg.DialAttemptBudget, dialConfigFixture().DialAttemptBudget+1, "dial backoff must be adjusted")
	require.Equal(t, cfg.StreamCreationAttemptBudget, dialConfigFixture().StreamCreationAttemptBudget, "stream backoff must be 1")

	// Recurring adjustment of the dial config of peerID1 through Adjust.
	// dial config for peerID1 already exists in the cache, so it must be adjusted when using Adjust.
	cfg, err = cache.Adjust(peerID1, adjustFuncIncrement)
	require.NoError(t, err)
	// adjusting an existing dial config must not change the cache size.
	require.Equal(t, uint(2), cache.Size(), "cache size must be 2")
	require.Equal(t, cfg.LastSuccessfulDial, dialConfigFixture().LastSuccessfulDial, "last successful dial must be 0")
	require.Equal(t, cfg.DialAttemptBudget, dialConfigFixture().DialAttemptBudget+2, "dial backoff must be adjusted")
	require.Equal(t, cfg.StreamCreationAttemptBudget, dialConfigFixture().StreamCreationAttemptBudget, "stream backoff must be 1")
}

// TestDialConfigCache_Adjust tests the Adjust method of the DialConfigCache. It asserts that the dial config is adjusted,
// and stored in the cache as expected under concurrent adjustments.
func TestDialConfigCache_Concurrent_Adjust(t *testing.T) {
	sizeLimit := uint32(100)
	logger := zerolog.Nop()
	collector := metrics.NewNoopCollector()

	cache := unicastcache.NewDialConfigCache(sizeLimit, logger, collector, func() unicastmodel.DialConfig {
		return unicastmodel.DialConfig{} // empty dial config
	})
	require.NotNil(t, cache)
	require.Zerof(t, cache.Size(), "cache size must be 0")

	peerIds := make([]peer.ID, sizeLimit)
	for i := 0; i < int(sizeLimit); i++ {
		peerId := p2ptest.PeerIdFixture(t)
		require.NotContainsf(t, peerIds, peerId, "peer id must be unique")
		peerIds[i] = peerId
	}

	wg := sync.WaitGroup{}
	for i := 0; i < int(sizeLimit); i++ {
		// adjusts the ith dial config for peerID i times, concurrently.
		for j := 0; j < i+1; j++ {
			wg.Add(1)
			go func(peerId peer.ID) {
				defer wg.Done()
				_, err := cache.Adjust(peerId, func(cfg unicastmodel.DialConfig) (unicastmodel.DialConfig, error) {
					cfg.DialAttemptBudget++
					return cfg, nil
				})
				require.NoError(t, err)
			}(peerIds[i])
		}
	}

	unittest.RequireReturnsBefore(t, wg.Wait, time.Second*1, "adjustments must be done on time")

	// assert that the cache size is equal to the size limit.
	require.Equal(t, uint(sizeLimit), cache.Size(), "cache size must be equal to the size limit")

	// assert that the dial config for each peer is adjusted i times, concurrently.
	for i := 0; i < int(sizeLimit); i++ {
		wg.Add(1)
		go func(j int) {
			peerID := peerIds[j]
			cfg, err := cache.GetOrInit(peerID)
			require.NoError(t, err)
			require.Equal(t, uint64(j+1), cfg.DialAttemptBudget, fmt.Sprintf("peerId %s dial backoff must be adjusted %d times got: %d", peerID, j+1, cfg.DialAttemptBudget))
		}(i)
	}

	unittest.RequireReturnsBefore(t, wg.Wait, time.Second*1, "retrievals must be done on time")
}

// TestConcurrent_Adjust_And_Get_Is_Safe tests that concurrent adjustments and retrievals are safe, and do not cause error even if they cause eviction. The test stress tests the cache
// with 2 * SizeLimit concurrent operations (SizeLimit times concurrent adjustments and SizeLimit times concurrent retrievals).
// It asserts that the cache size is equal to the size limit, and the dial config for each peer is adjusted and retrieved correctly.
func TestConcurrent_Adjust_And_Get_Is_Safe(t *testing.T) {
	sizeLimit := uint32(100)
	logger := zerolog.Nop()
	collector := metrics.NewNoopCollector()

	cache := unicastcache.NewDialConfigCache(sizeLimit, logger, collector, dialConfigFixture)
	require.NotNil(t, cache)
	require.Zerof(t, cache.Size(), "cache size must be 0")

	wg := sync.WaitGroup{}
	for i := 0; i < int(sizeLimit); i++ {
		// concurrently adjusts the dial configs.
		wg.Add(1)
		go func() {
			defer wg.Done()
			peerId := p2ptest.PeerIdFixture(t)
			dialTime := time.Now()
			updatedConfig, err := cache.Adjust(peerId, func(cfg unicastmodel.DialConfig) (unicastmodel.DialConfig, error) {
				cfg.DialAttemptBudget = 1 // some random adjustment
				cfg.LastSuccessfulDial = dialTime
				cfg.StreamCreationAttemptBudget = 2 // some random adjustment
				cfg.ConsecutiveSuccessfulStream = 3 // some random adjustment
				return cfg, nil
			})
			require.NoError(t, err)                                      // concurrent adjustment must not fail.
			require.Equal(t, uint64(1), updatedConfig.DialAttemptBudget) // adjustment must be successful
			require.Equal(t, uint64(2), updatedConfig.StreamCreationAttemptBudget)
			require.Equal(t, uint64(3), updatedConfig.ConsecutiveSuccessfulStream)
			require.Equal(t, dialTime, updatedConfig.LastSuccessfulDial)
		}()
	}

	// assert that the dial config for each peer is adjusted i times, concurrently.
	for i := 0; i < int(sizeLimit); i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			peerId := p2ptest.PeerIdFixture(t)
			cfg, err := cache.GetOrInit(peerId)
			require.NoError(t, err)                                                        // concurrent retrieval must not fail.
			require.Equal(t, dialConfigFixture().DialAttemptBudget, cfg.DialAttemptBudget) // dial config must be initialized with the default values.
			require.Equal(t, dialConfigFixture().StreamCreationAttemptBudget, cfg.StreamCreationAttemptBudget)
			require.Equal(t, uint64(0), cfg.ConsecutiveSuccessfulStream)
			require.True(t, cfg.LastSuccessfulDial.IsZero())
		}()
	}

	unittest.RequireReturnsBefore(t, wg.Wait, time.Second*1, "all operations must be done on time")

	// cache was stress-tested with 2 * SizeLimit concurrent operations. Nevertheless, the cache size must be equal to the size limit due to LRU eviction.
	require.Equal(t, uint(sizeLimit), cache.Size(), "cache size must be equal to the size limit")
}

// TODO: test LRU eviction
