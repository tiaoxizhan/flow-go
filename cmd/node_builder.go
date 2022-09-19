package cmd

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dgraph-io/badger/v2"
	madns "github.com/multiformats/go-multiaddr-dns"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/onflow/flow-go/admin/commands"
	"github.com/onflow/flow-go/crypto"
	"github.com/onflow/flow-go/fvm"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/chainsync"
	"github.com/onflow/flow-go/module/compliance"
	"github.com/onflow/flow-go/module/component"
	"github.com/onflow/flow-go/module/id"
	"github.com/onflow/flow-go/network"
	"github.com/onflow/flow-go/network/codec/cbor"
	"github.com/onflow/flow-go/network/p2p"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/state/protocol/events"
	bstorage "github.com/onflow/flow-go/storage/badger"
)

const NotSet = "not set"

type BuilderFunc func(nodeConfig *NodeConfig) error
type ReadyDoneFactory func(node *NodeConfig) (module.ReadyDoneAware, error)

// NodeBuilder declares the initialization methods needed to bootstrap up a Flow node
type NodeBuilder interface {
	// BaseFlags reads the command line arguments common to all nodes
	BaseFlags()

	// ExtraFlags reads the node specific command line arguments and adds it to the FlagSet
	ExtraFlags(f func(*pflag.FlagSet)) NodeBuilder

	// ParseAndPrintFlags parses and validates all the command line arguments
	ParseAndPrintFlags() error

	// Initialize performs all the initialization needed at the very start of a node
	Initialize() error

	// PrintBuildVersionDetails prints the node software build version
	PrintBuildVersionDetails()

	// InitIDProviders initializes the ID providers needed by various components
	InitIDProviders()

	// EnqueueNetworkInit enqueues the default networking layer.
	EnqueueNetworkInit()

	// EnqueueMetricsServerInit enqueues the metrics component.
	EnqueueMetricsServerInit()

	// EnqueueTracer enqueues the Tracer component.
	EnqueueTracer()

	// Module enables setting up dependencies of the engine with the builder context
	Module(name string, f BuilderFunc) NodeBuilder

	// Component adds a new component to the node that conforms to the ReadyDoneAware
	// interface, and throws a Fatal() when an irrecoverable error is encountered.
	//
	// The ReadyDoneFactory may return either a `Component` or `ReadyDoneAware` instance.
	// In both cases, the object is started according to its interface when the node is run,
	// and the node will wait for the component to exit gracefully.
	Component(name string, f ReadyDoneFactory) NodeBuilder

	// RestartableComponent adds a new component to the node that conforms to the ReadyDoneAware
	// interface, and calls the provided error handler when an irrecoverable error is encountered.
	// Use RestartableComponent if the component is not critical to the node's safe operation and
	// can/should be independently restarted when an irrecoverable error is encountered.
	//
	// Any irrecoverable errors thrown by the component will be passed to the provided error handler.
	RestartableComponent(name string, f ReadyDoneFactory, errorHandler component.OnError) NodeBuilder

	// ShutdownFunc adds a callback function that is called after all components have exited.
	// All shutdown functions are called regardless of errors returned by previous callbacks. Any
	// errors returned are captured and passed to the caller.
	ShutdownFunc(fn func() error) NodeBuilder

	// AdminCommand registers a new admin command with the admin server
	AdminCommand(command string, f func(config *NodeConfig) commands.AdminCommand) NodeBuilder

	// MustNot asserts that the given error must not occur.
	// If the error is nil, returns a nil log event (which acts as a no-op).
	// If the error is not nil, returns a fatal log event containing the error.
	MustNot(err error) *zerolog.Event

	// Build finalizes the node configuration in preparation for start and returns a Node
	// object that can be run
	Build() (Node, error)

	// PreInit registers a new PreInit function.
	// PreInit functions run before the protocol state is initialized or any other modules or components are initialized
	PreInit(f BuilderFunc) NodeBuilder

	// PostInit registers a new PreInit function.
	// PostInit functions run after the protocol state has been initialized but before any other modules or components
	// are initialized
	PostInit(f BuilderFunc) NodeBuilder

	// RegisterBadgerMetrics registers all badger related metrics
	RegisterBadgerMetrics() error

	// ValidateFlags sets any custom validation rules for the command line flags,
	// for example where certain combinations aren't allowed
	ValidateFlags(func() error) NodeBuilder
}

// BaseConfig is the general config for the NodeBuilder and the command line params
// For a node running as a standalone process, the config fields will be populated from the command line params,
// while for a node running as a library, the config fields are expected to be initialized by the caller.
type BaseConfig struct {
	nodeIDHex                   string
	AdminAddr                   string
	AdminCert                   string
	AdminKey                    string
	AdminClientCAs              string
	BindAddr                    string
	NodeRole                    string
	DynamicStartupANAddress     string
	DynamicStartupANPubkey      string
	DynamicStartupEpochPhase    string
	DynamicStartupEpoch         string
	DynamicStartupSleepInterval time.Duration
	datadir                     string
	secretsdir                  string
	secretsDBEnabled            bool
	InsecureSecretsDB           bool
	level                       string
	metricsPort                 uint
	BootstrapDir                string
	PeerUpdateInterval          time.Duration
	UnicastMessageTimeout       time.Duration
	DNSCacheTTL                 time.Duration
	profilerEnabled             bool
	uploaderEnabled             bool
	profilerDir                 string
	profilerInterval            time.Duration
	profilerDuration            time.Duration
	profilerMemProfileRate      int
	tracerEnabled               bool
	tracerSensitivity           uint
	MetricsEnabled              bool
	guaranteesCacheSize         uint
	receiptsCacheSize           uint
	db                          *badger.DB
	PreferredUnicastProtocols   []string
	// NetworkConnectionPruning determines whether connections to nodes
	// that are not part of protocol state should be trimmed
	// TODO: solely a fallback mechanism, can be removed upon reliable behavior in production.
	NetworkConnectionPruning        bool
	NetworkReceivedMessageCacheSize uint32
	HeroCacheMetricsEnable          bool
	SyncCoreConfig                  chainsync.Config
	CodecFactory                    func() network.Codec
	// ComplianceConfig configures either the compliance engine (consensus nodes)
	// or the follower engine (all other node roles)
	ComplianceConfig compliance.Config

	// UnicastRateLimitDryRun will disable connection disconnects and gating when unicast rate limiters are configured
	UnicastRateLimitDryRun bool
	// UnicastMessageRateLimit amount of unicast messages that can be sent by a peer per second.
	UnicastMessageRateLimit int
	// UnicastMessageBurstLimit amount of unicast messages that can be sent by a peer at once.
	UnicastMessageBurstLimit int
	// UnicastBandwidthRateLimit bandwidth size in bytes a peer is allowed to send via unicast streams per second.
	UnicastBandwidthRateLimit int
	// UnicastBandwidthBurstLimit bandwidth size in bytes a peer is allowed to send via unicast streams at once.
	UnicastBandwidthBurstLimit int
}

// NodeConfig contains all the derived parameters such the NodeID, private keys etc. and initialized instances of
// structs such as DB, Network etc. The NodeConfig is composed of the BaseConfig and is updated in the
// NodeBuilder functions as a node is bootstrapped.
type NodeConfig struct {
	Cancel context.CancelFunc // cancel function for the context that is passed to the networking layer
	BaseConfig
	Logger            zerolog.Logger
	NodeID            flow.Identifier
	Me                module.Local
	Tracer            module.Tracer
	MetricsRegisterer prometheus.Registerer
	Metrics           Metrics
	DB                *badger.DB
	SecretsDB         *badger.DB
	Storage           Storage
	ProtocolEvents    *events.Distributor
	State             protocol.State
	Resolver          madns.BasicResolver
	Middleware        network.Middleware
	Network           network.Network
	PingService       network.PingService
	MsgValidators     []network.MessageValidator
	FvmOptions        []fvm.Option
	StakingKey        crypto.PrivateKey
	NetworkKey        crypto.PrivateKey

	// ID providers
	IdentityProvider             id.IdentityProvider
	IDTranslator                 p2p.IDTranslator
	SyncEngineIdentifierProvider id.IdentifierProvider

	// root state information
	RootSnapshot protocol.Snapshot
	// cached properties of RootSnapshot for convenience
	RootBlock   *flow.Block
	RootQC      *flow.QuorumCertificate
	RootResult  *flow.ExecutionResult
	RootSeal    *flow.Seal
	RootChainID flow.ChainID
	SporkID     flow.Identifier

	// bootstrapping options
	SkipNwAddressBasedValidations bool
}

func DefaultBaseConfig() *BaseConfig {
	homedir, _ := os.UserHomeDir()
	datadir := filepath.Join(homedir, ".flow", "database")

	// NOTE: if the codec used in the network component is ever changed any code relying on
	// the message format specific to the codec must be updated. i.e: the AuthorizedSenderValidator.
	codecFactory := func() network.Codec { return cbor.NewCodec() }

	return &BaseConfig{
		nodeIDHex:                       NotSet,
		AdminAddr:                       NotSet,
		AdminCert:                       NotSet,
		AdminKey:                        NotSet,
		AdminClientCAs:                  NotSet,
		BindAddr:                        NotSet,
		BootstrapDir:                    "bootstrap",
		datadir:                         datadir,
		secretsdir:                      NotSet,
		secretsDBEnabled:                true,
		level:                           "info",
		PeerUpdateInterval:              p2p.DefaultPeerUpdateInterval,
		UnicastMessageTimeout:           p2p.DefaultUnicastTimeout,
		metricsPort:                     8080,
		profilerEnabled:                 false,
		uploaderEnabled:                 false,
		profilerDir:                     "profiler",
		profilerInterval:                15 * time.Minute,
		profilerDuration:                10 * time.Second,
		profilerMemProfileRate:          runtime.MemProfileRate,
		tracerEnabled:                   false,
		tracerSensitivity:               4,
		MetricsEnabled:                  true,
		receiptsCacheSize:               bstorage.DefaultCacheSize,
		guaranteesCacheSize:             bstorage.DefaultCacheSize,
		NetworkReceivedMessageCacheSize: p2p.DefaultReceiveCacheSize,

		// By default we let networking layer trim connections to all nodes that
		// are no longer part of protocol state.
		NetworkConnectionPruning: p2p.ConnectionPruningEnabled,

		HeroCacheMetricsEnable: false,
		SyncCoreConfig:         chainsync.DefaultConfig(),
		CodecFactory:           codecFactory,
		ComplianceConfig:       compliance.DefaultConfig(),
	}
}
