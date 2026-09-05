package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/N30A/trakt/device"
	"github.com/N30A/trakt/position"
)

type positionHandler struct {
	repo       *position.PositionRepo
	deviceRepo *device.DeviceRepo
}

func newPositionHandler(repo *position.PositionRepo, deviceRepo *device.DeviceRepo) *positionHandler {
	return &positionHandler{repo: repo, deviceRepo: deviceRepo}
}

func (h *positionHandler) getPositions(w http.ResponseWriter, r *http.Request) {
	deviceID, hasDeviceID, err := parseDeviceIDQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !hasDeviceID {
		http.Error(w, "device_id is required", http.StatusBadRequest)
		return
	}

	from, to, err := parseTimeRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dev, err := h.deviceRepo.GetDeviceByID(r.Context(), deviceID)
	if err != nil {
		switch {
		case errors.Is(err, device.ErrNotFound):
			http.Error(w, "device not found", http.StatusNotFound)
			return
		default:
			log.Printf("failed to retrieve device %d: %v", deviceID, err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	var positions []position.Position

	positions, err = h.repo.GetPositionsByDevice(r.Context(), dev.ID, from, to)
	if err != nil {
		log.Printf("failed to retrieve positions for device %d: %v", dev.ID, err)
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

func (h *positionHandler) getLatestPosition(w http.ResponseWriter, r *http.Request) {
	deviceID, hasDeviceID, err := parseDeviceIDQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var positions []position.Position

	if hasDeviceID {
		dev, err := h.deviceRepo.GetDeviceByID(r.Context(), deviceID)
		if err != nil {
			switch {
			case errors.Is(err, device.ErrNotFound):
				http.Error(w, "device not found", http.StatusNotFound)
				return
			default:
				log.Printf("failed to retrieve device %d: %v", deviceID, err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		pos, err := h.repo.GetLatestPositionByDevice(r.Context(), dev.ID)
		if err != nil {
			switch {
			case errors.Is(err, position.ErrNotFound):
				positions = []position.Position{}
			default:
				log.Printf("failed to retrieve latest position for device %d: %v", dev.ID, err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		} else {
			positions = []position.Position{pos}
		}
	} else {
		positions, err = h.repo.GetLatestPositions(r.Context())
		if err != nil {
			log.Printf("failed to retrieve latest positions: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
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
