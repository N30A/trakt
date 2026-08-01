package api

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/N30A/trakt/repos"
)

const (
	host = "0.0.0.0"
	port = "8080"
)

type APIServer struct {
	deviceRepo   *repos.DeviceRepo
	positionRepo *repos.PositionRepo
	server       *http.Server
}

func New(deviceRepo *repos.DeviceRepo, positionRepo *repos.PositionRepo) *APIServer {
	mux := http.NewServeMux()
	server := &APIServer{
		deviceRepo:   deviceRepo,
		positionRepo: positionRepo,
		server: &http.Server{
			Addr:    net.JoinHostPort(host, port),
			Handler: mux,
		},
	}

	server.registerRoutes(mux)
	return server
}

func (s *APIServer) Name() string {
	return "api"
}

func (s *APIServer) Start() error {
	log.Printf("%s listening on %s", s.Name(), s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *APIServer) Stop(ctx context.Context) error {
	log.Printf("stopping %s", s.Name())
	return s.server.Shutdown(ctx)
}
