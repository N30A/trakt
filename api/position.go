package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/N30A/trakt/models"
	"github.com/N30A/trakt/repos"
)

type positionHandler struct {
	repo       *repos.PositionRepo
	deviceRepo *repos.DeviceRepo
}

func newPositionHandler(repo *repos.PositionRepo, deviceRepo *repos.DeviceRepo) *positionHandler {
	return &positionHandler{repo: repo, deviceRepo: deviceRepo}
}

func (h *positionHandler) getPositions(w http.ResponseWriter, r *http.Request) {
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
		positions, err = h.repo.GetPositionsByDevice(r.Context(), deviceID, from, to)
	} else {
		positions, err = h.repo.GetPositions(r.Context(), from, to)
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

func (h *positionHandler) getLatestPosition(w http.ResponseWriter, r *http.Request) {
	deviceID, hasDeviceID, err := parseDeviceIDQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var positions []models.Position

	if hasDeviceID {
		device, err := h.deviceRepo.GetDeviceByID(r.Context(), deviceID)
		if err != nil {
			switch {
			case errors.Is(err, repos.ErrNotFound):
				http.Error(w, "device not found", http.StatusNotFound)
				return
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		position, err := h.repo.GetLatestPositionByDevice(r.Context(), device.ID)
		if err != nil {
			switch {
			case errors.Is(err, repos.ErrNotFound):
				positions = []models.Position{}
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		} else {
			positions = []models.Position{position}
		}
	} else {
		positions, err = h.repo.GetLatestPositions(r.Context())
		if err != nil {
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
			Protocol:   string(position.Protocol),
			Altitude:   position.Altitude,
			Speed:      position.Speed,
			Course:     position.Course,
			Accuracy:   position.Accuracy,
		}
	}

	writeJSON(w, http.StatusOK, response)
}
