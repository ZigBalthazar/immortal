package relay

import (
	"log"

	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/handler"
	"github.com/dezh-tech/immortal/management"
	"github.com/dezh-tech/immortal/management/routes"
	"github.com/dezh-tech/immortal/metrics"
	"github.com/dezh-tech/immortal/relay/redis"
	"github.com/dezh-tech/immortal/websocket"
)

// Relay keeps all concepts such as server, database and manages them.
type Relay struct {
	config          config.Config
	websocketServer *websocket.Server
	database        *database.Database
	redis           *redis.Redis
	management      *management.Server
}

// NewRelay creates a new relay.
func New(cfg *config.Config) (*Relay, error) {
	db, err := database.Connect(cfg.Database)
	if err != nil {
		return nil, err
	}

	err = cfg.LoadParameters(db)
	if err != nil {
		return nil, err
	}

	h := handler.New(db, cfg.Parameters.Handler)

	m := metrics.New()

	r, err := redis.New(cfg.RedisConf)
	if err != nil {
		return nil, err
	}

	ws, err := websocket.New(cfg.WebsocketServer, cfg.GetNIP11Documents(), h, m, r)
	if err != nil {
		return nil, err
	}

	managementServer, err := management.NewServer(cfg.Management)
	if err != nil {
		return nil, err
	}

	return &Relay{
		config:          *cfg,
		websocketServer: ws,
		database:        db,
		redis:           r,
		management: managementServer,
	}, nil
}

// Start runs the relay and its children.
func (r *Relay) Start() chan error {
	log.Println("relay started successfully...")
	errCh := make(chan error, 2)

	go func() {
		if err := r.websocketServer.Start(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		routes.ConfigureRoutes(r.management)
		if err := r.management.Start(); err != nil {
			errCh <- err
		}
	}()

	return errCh
}

// Stop shutdowns the relay and its children gracefully.
func (r *Relay) Stop() error {
	log.Println("stopping relay...")

	if err := r.websocketServer.Stop(); err != nil {
		return err
	}

	return r.database.Stop()
}
