package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/N30A/trakt/position"
)

type positionHandler struct {
	service *position.PositionService
}

func newPositionHandler(service *position.PositionService) *positionHandler {
	return &positionHandler{service: service}
}

func parseDeviceIDs(r *http.Request) ([]int, error) {
	query := r.URL.Query()
	key := "device_id"

	if !query.Has(key) {
		return nil, errors.New("at least one device_id is required")
	}

	deviceIDs := make([]int, len(query[key]))
	for i, id := range query[key] {
		deviceID, err := strconv.Atoi(strings.TrimSpace(id))
		if err != nil || deviceID <= 0 {
			return nil, errors.New("device_id must be a positive integer")
		}
		deviceIDs[i] = deviceID
	}

	return deviceIDs, nil
}

func parseDateRange(r *http.Request) (time.Time, time.Time, error) {
	query := r.URL.Query()
	fromStr := strings.TrimSpace(query.Get("from"))
	toStr := strings.TrimSpace(query.Get("to"))

	if fromStr == "" || toStr == "" {
		return time.Time{}, time.Time{}, errors.New("from and to must be specified together")
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("from must be a valid RFC3339 timestamp")
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("to must be a valid RFC3339 timestamp")
	}

	return from, to, nil
}

func (h *positionHandler) devicePositions(w http.ResponseWriter, r *http.Request) {
	deviceIDs, err := parseDeviceIDs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	from, to, err := parseDateRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var positions []position.Position

	positions, err = h.service.DevicePositions(r.Context(), deviceIDs, from, to)
	if err != nil {
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
			Altitude:   position.Altitude,
			Speed:      position.Speed,
			Course:     position.Course,
			Accuracy:   position.Accuracy,
		}
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *positionHandler) latestPositions(w http.ResponseWriter, r *http.Request) {
	var positions []position.Position

	positions, err := h.service.LatestPositions(r.Context())
	if err != nil {
		log.Printf("failed to retrieve latest positions: %v", err)
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
			Altitude:   position.Altitude,
			Speed:      position.Speed,
			Course:     position.Course,
			Accuracy:   position.Accuracy,
		}
	}

	writeJSON(w, http.StatusOK, response)
}
