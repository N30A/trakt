package osmand

// https://www.traccar.org/osmand

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/N30A/trakt/models"
	"github.com/N30A/trakt/repos"
)

const (
	host = "0.0.0.0"
	port = "5055"
)

type OsmandServer struct {
	deviceRepo   *repos.DeviceRepo
	positionRepo *repos.PositionRepo
	server       *http.Server
}

func New(deviceRepo *repos.DeviceRepo, positionRepo *repos.PositionRepo) *OsmandServer {
	mux := http.NewServeMux()
	server := &OsmandServer{
		deviceRepo:   deviceRepo,
		positionRepo: positionRepo,
		server: &http.Server{
			Addr:    net.JoinHostPort(host, port),
			Handler: mux,
		},
	}

	mux.HandleFunc("POST /", server.handler)
	return server
}

func (s *OsmandServer) Name() string {
	return "osmand"
}

func (s *OsmandServer) Start() error {
	log.Printf("%s listening on %s", s.Name(), s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *OsmandServer) Stop(ctx context.Context) error {
	log.Printf("stopping %s\n", s.Name())
	return s.server.Shutdown(ctx)
}

func (s *OsmandServer) handler(w http.ResponseWriter, r *http.Request) {
	// expects to recive: application/x-www-form-urlencoded
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

	device, err := s.deviceRepo.GetDeviceByUniqueID(r.Context(), parsed.DeviceUniqueID)
	if errors.Is(err, repos.ErrNotFound) {
		log.Printf("unknown osmand device: %s", parsed.DeviceUniqueID)
		http.Error(w, "unknown device", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("failed to get device %s: %v", parsed.DeviceUniqueID, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	position := models.Position{
		DeviceID:   device.ID,
		Latitude:   parsed.Latitude,
		Longitude:  parsed.Longitude,
		FixTime:    parsed.Timestamp,
		ServerTime: serverTime,
		Protocol:   models.ProtocolOsmand,
		Altitude:   parsed.Altitude,
		Speed:      parsed.Speed,
		Course:     parsed.Course,
		Accuracy:   parsed.Accuracy,
	}

	if err := s.positionRepo.AddPosition(r.Context(), position); err != nil {
		log.Printf("failed to add position for device %s: %v", device.UniqueID, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
