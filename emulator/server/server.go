package server

import (
	"encoding/hex"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/emulator"
	"github.com/dapperlabs/flow-go-sdk/emulator/storage"
	"github.com/dapperlabs/flow-go-sdk/emulator/storage/badger"
	"github.com/dapperlabs/flow-go-sdk/emulator/storage/memstore"
	"github.com/dapperlabs/flow-go-sdk/utils/graceful"
)

// EmulatorServer is a local server that runs a Flow Emulator instance.
//
// The server wraps an emulated blockchain instance with the Observation gRPC interface.
type EmulatorServer struct {
	logger         *logrus.Logger
	config         *Config
	backend        *Backend
	grpcServer     *GRPCServer
	httpServer     *HTTPServer
	blocksTicker   *BlocksTicker
	livenessTicker *LivenessTicker
	onCleanup      func()
}

const (
	defaultGRPCPort               = 3569
	defaultHTTPPort               = 8080
	defaultLivenessCheckTolerance = time.Second
)

var (
	defaultHTTPHeaders = []HTTPHeader{
		{
			Key:   "Access-Control-Allow-Origin",
			Value: "*",
		},
		{
			Key:   "Access-Control-Allow-Methods",
			Value: "POST, GET, OPTIONS, PUT, DELETE",
		},
		{
			Key:   "Access-Control-Allow-Headers",
			Value: "*",
		},
	}
)

// Config is the configuration for an emulator server.
type Config struct {
	GRPCPort       int
	GRPCDebug      bool
	HTTPPort       int
	HTTPHeaders    []HTTPHeader
	BlockTime      time.Duration
	RootAccountKey *flow.AccountPrivateKey
	Persist        bool
	// DBPath is the path to the Badger database on disk
	DBPath string
	// LivenessCheckTolerance is the tolerance level of the liveness check
	// e.g. how long we can go without answering before being considered not alive
	LivenessCheckTolerance time.Duration
}

// NewEmulatorServer creates a new instance of a Flow Emulator server.
func NewEmulatorServer(logger *logrus.Logger, conf *Config) *EmulatorServer {

	conf = sanitizeConfig(conf)

	store, closeStore, err := configureStore(logger, conf)
	if err != nil {
		logger.WithError(err).Error("❗  Failed to configure storage")
		return nil
	}

	blockchain, err := configureBlockchain(conf, store)
	if err != nil {
		logger.WithError(err).Error("❗  Failed to configure emulated blockchain")
		return nil
	}

	backend := configureBackend(logger, conf, blockchain)

	livenessTicker := NewLivenessTicker(conf.LivenessCheckTolerance)
	grpcServer := NewGRPCServer(logger, backend, conf.GRPCPort, conf.GRPCDebug)
	httpServer := NewHTTPServer(grpcServer, livenessTicker, conf.HTTPPort, conf.HTTPHeaders)

	server := &EmulatorServer{
		logger:         logger,
		config:         conf,
		backend:        backend,
		grpcServer:     grpcServer,
		httpServer:     httpServer,
		livenessTicker: livenessTicker,
		onCleanup: func() {
			err := closeStore()
			if err != nil {
				logger.WithError(err).Infof("Failed to close storage")
			}
		},
	}

	// only create blocks ticker if block time > 0
	if conf.BlockTime > 0 {
		server.blocksTicker = NewBlocksTicker(backend, conf.BlockTime)

	}

	address := blockchain.RootAccountAddress()
	prKey := blockchain.RootKey()
	prKeyBytes, _ := prKey.PrivateKey.Encode()

	logger.WithFields(logrus.Fields{
		"address": address.Hex(),
		"prKey":   hex.EncodeToString(prKeyBytes),
	}).Infof("⚙️   Using root account 0x%s", address.Hex())

	return server
}

// Start starts the Flow Emulator server.
func (s *EmulatorServer) Start() {
	defer s.cleanup()

	g := graceful.NewGroup()

	s.logger.
		WithField("port", s.config.GRPCPort).
		Infof("🌱  Starting gRPC server on port %d...", s.config.GRPCPort)

	s.logger.
		WithField("port", s.config.HTTPPort).
		Infof("🌱  Starting HTTP server on port %d...", s.config.HTTPPort)

	g.Add(s.grpcServer)
	g.Add(s.httpServer)
	g.Add(s.livenessTicker)

	// only start blocks ticker if it exists
	if s.blocksTicker != nil {
		g.Add(s.blocksTicker)
	}

	err := g.Start()
	if err != nil {
		s.logger.WithError(err).Error("❗  Server error")
	}

	s.logger.Info("🛑  Server stopped")
}

// cleanup cleans up the server.
// This MUST be called before the server process terminates.
func (e *EmulatorServer) cleanup() {
	e.onCleanup()
}

func configureStore(logger *logrus.Logger, conf *Config) (store storage.Store, close func() error, err error) {
	if conf.Persist {
		badgerStore, err := badger.New(
			badger.WithPath(conf.DBPath),
			badger.WithLogger(logger),
			badger.WithTruncate(true),
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to initialize Badger store")
		}

		close = func() error {
			err := badgerStore.Close()
			if err != nil {
				return err
			}

			return nil
		}

		return badgerStore, close, nil
	}

	store = memstore.New()
	close = func() error { return nil }

	return store, close, nil
}

func configureBlockchain(conf *Config, store storage.Store) (*emulator.Blockchain, error) {
	options := []emulator.Option{
		emulator.WithStore(store),
	}

	if conf.RootAccountKey != nil {
		options = append(options, emulator.WithRootAccountKey(*conf.RootAccountKey))
	}

	blockchain, err := emulator.NewBlockchain(options...)
	if err != nil {
		return nil, err
	}

	return blockchain, nil
}

func configureBackend(logger *logrus.Logger, conf *Config, blockchain *emulator.Blockchain) *Backend {
	backend := NewBackend(logger, blockchain)

	if conf.BlockTime == 0 {
		backend.EnableAutoMine()
	}

	return backend
}

func sanitizeConfig(conf *Config) *Config {
	if conf.GRPCPort == 0 {
		conf.GRPCPort = defaultGRPCPort
	}

	if conf.HTTPPort == 0 {
		conf.HTTPPort = defaultHTTPPort
	}

	if conf.HTTPHeaders == nil {
		conf.HTTPHeaders = defaultHTTPHeaders
	}

	if conf.LivenessCheckTolerance == 0 {
		conf.LivenessCheckTolerance = defaultLivenessCheckTolerance
	}

	return conf
}
