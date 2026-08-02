package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/N30A/trakt/models"
	"github.com/N30A/trakt/repos"
)

func (s *APIServer) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	deviceHandler := newDeviceHandler(s.deviceRepo)

	mux.HandleFunc("GET /devices", deviceHandler.getDevices)
	mux.HandleFunc("GET /devices/{id}", deviceHandler.getDevice)
	mux.HandleFunc("POST /devices", deviceHandler.createDevice)
	mux.HandleFunc("PUT /devices/{id}", deviceHandler.updateDevice)
	mux.HandleFunc("DELETE /devices/{id}", deviceHandler.deleteDevice)

	mux.HandleFunc("GET /positions", s.getPositions)
}

func (s *APIServer) getPositions(w http.ResponseWriter, r *http.Request) {
	deviceID, hasDeviceID, err := parseDeviceIDQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	from, to, err := parseTimeRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var positions []models.Position

	if hasDeviceID {
		positions, err = s.positionRepo.GetPositionsByDevice(r.Context(), deviceID, from, to)
	} else {
		positions, err = s.positionRepo.GetPositions(r.Context(), from, to)
	}

	if err != nil {
		if errors.Is(err, repos.ErrNotFound) {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		log.Printf("failed to retrieve positions: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	response := make([]positionResponse, len(positions))
	for i, position := range positions {
		response[i] = positionResponse{
			ID:         position.ID,
			DeviceID:   position.DeviceID,
			Latitude:   position.Latitude,
			Longitude:  position.Longitude,
			FixTime:    position.FixTime,
			ServerTime: position.ServerTime,
			Protocol:   string(position.Protocol),
			Altitude:   position.Altitude,
			Speed:      position.Speed,
			Course:     position.Course,
			Accuracy:   position.Accuracy,
		}
	}

	writeJSON(w, http.StatusOK, response)
}
