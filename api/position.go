package api

import (
	"log"
	"net/http"

	"github.com/N30A/trakt/position"
)

type positionHandler struct {
	service *position.PositionService
}

func newPositionHandler(service *position.PositionService) *positionHandler {
	return &positionHandler{service: service}
}

func (h *positionHandler) positions(w http.ResponseWriter, r *http.Request) {
	deviceID, err := parseDeviceIDQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	from, to, err := parseTimeRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var positions []position.Position

	positions, err = h.service.Positions(r.Context(), deviceID, from, to)
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
