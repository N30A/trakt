package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/N30A/trakt/internal/auth"
	"github.com/N30A/trakt/internal/device"
	"github.com/N30A/trakt/internal/position"
)

const (
	host = "0.0.0.0"
	port = "8080"
)

type APIServer struct {
	deviceRepo      *device.DeviceRepo
	positionService *position.PositionService
	authService     *auth.AuthService
	server          *http.Server
}

func NewAPIServer(deviceRepo *device.DeviceRepo, positionService *position.PositionService, authService *auth.AuthService) *APIServer {
	mux := http.NewServeMux()
	server := &APIServer{
		deviceRepo:      deviceRepo,
		positionService: positionService,
		authService:     authService,
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
