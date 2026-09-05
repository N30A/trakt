package osmand

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/N30A/trakt/position"
)

const (
	host = "0.0.0.0"
	port = "5055"
)

type OsmAndServer struct {
	positionService *position.PositionService
	server          *http.Server
}

func New(positionService *position.PositionService) *OsmAndServer {
	mux := http.NewServeMux()
	server := &OsmAndServer{
		positionService: positionService,
		server: &http.Server{
			Addr:    net.JoinHostPort(host, port),
			Handler: mux,
		},
	}

	mux.HandleFunc("POST /", server.handler)
	return server
}

func (s *OsmAndServer) Name() string {
	return "osmand"
}

func (s *OsmAndServer) Start() error {
	log.Printf("%s listening on %s", s.Name(), s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *OsmAndServer) Stop(ctx context.Context) error {
	log.Printf("stopping %s\n", s.Name())
	return s.server.Shutdown(ctx)
}

func (s *OsmAndServer) handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	serverTime := time.Now().UTC()

	parsed, err := parsePosition(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := toPositionInput(parsed, serverTime)

	if err := s.positionService.SavePosition(r.Context(), input); err != nil {
		if errors.Is(err, position.ErrNotFound) {
			slog.Warn("unknown device", "protocol", input.Protocol, "device_unique_id", input.DeviceUniqueID)
			http.Error(w, "unknown device", http.StatusUnauthorized)
			return
		}

		slog.Error("failed to save position", "protocol", input.Protocol, "device_unique_id", input.DeviceUniqueID, "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
