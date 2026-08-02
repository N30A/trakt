package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/N30A/trakt/models"
	"github.com/N30A/trakt/repos"
)

func (s *APIServer) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /devices", s.getDevices)
	mux.HandleFunc("GET /devices/{id}", s.getDevice)
	mux.HandleFunc("POST /devices", s.createDevice)
	mux.HandleFunc("DELETE /devices/{id}", s.deleteDevice)

	mux.HandleFunc("GET /positions", s.getPositions)
}

func (s *APIServer) getDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := s.deviceRepo.GetDevices(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	response := make([]deviceResponse, len(devices))
	for i, device := range devices {
		response[i] = deviceResponse{
			ID:       device.ID,
			UniqueID: device.UniqueID,
			Name:     device.Name,
		}
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *APIServer) getDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "device id must be a positive integer", http.StatusBadRequest)
		return
	}

	device, err := s.deviceRepo.GetDeviceByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repos.ErrNotFound) {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, deviceResponse{
		ID:       device.ID,
		UniqueID: device.UniqueID,
		Name:     device.Name,
	})
}

func (s *APIServer) createDevice(w http.ResponseWriter, r *http.Request) {
	var request createDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "malformed request", http.StatusBadRequest)
		return
	}

	uniqueID := strings.TrimSpace(request.UniqueID)
	name := strings.TrimSpace(request.Name)

	if uniqueID == "" || name == "" {
		http.Error(w, "unique_id or name must not be empty", http.StatusBadRequest)
		return
	}

	device, err := s.deviceRepo.AddDevice(r.Context(), models.Device{UniqueID: uniqueID, Name: name})
	if err != nil {
		if errors.Is(err, repos.ErrConflict) {
			http.Error(w, "device already exists", http.StatusConflict)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, deviceResponse{
		ID:       device.ID,
		UniqueID: device.UniqueID,
		Name:     device.Name,
	})
}

func (s *APIServer) deleteDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "device id must be a positive integer", http.StatusBadRequest)
		return
	}

	err = s.deviceRepo.DeleteDeviceByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repos.ErrNotFound) {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
