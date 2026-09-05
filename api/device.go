package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/N30A/trakt/device"
)

type deviceHandler struct {
	repo *device.DeviceRepo
}

func newDeviceHandler(repo *device.DeviceRepo) *deviceHandler {
	return &deviceHandler{repo: repo}
}

func (h *deviceHandler) getDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.repo.GetDevices(r.Context())
	if err != nil {
		log.Printf("failed to retrieve devices: %v", err)
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

func (h *deviceHandler) getDevice(w http.ResponseWriter, r *http.Request) {
	id, err := parseDeviceIDPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dev, err := h.repo.GetDeviceByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, device.ErrNotFound) {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		log.Printf("failed to retrieve device %d: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, deviceResponse{
		ID:       dev.ID,
		UniqueID: dev.UniqueID,
		Name:     dev.Name,
	})
}

func (h *deviceHandler) createDevice(w http.ResponseWriter, r *http.Request) {
	var request createDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	uniqueID := strings.TrimSpace(request.UniqueID)
	name := strings.TrimSpace(request.Name)

	if uniqueID == "" || name == "" {
		http.Error(w, "unique_id and name must not be empty", http.StatusBadRequest)
		return
	}

	dev, err := h.repo.AddDevice(r.Context(), device.Device{UniqueID: uniqueID, Name: name})
	if err != nil {
		if errors.Is(err, device.ErrConflict) {
			http.Error(w, "device already exists", http.StatusConflict)
			return
		}

		log.Printf("failed to create device: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, deviceResponse{
		ID:       dev.ID,
		UniqueID: dev.UniqueID,
		Name:     dev.Name,
	})
}

func (h *deviceHandler) updateDevice(w http.ResponseWriter, r *http.Request) {
	id, err := parseDeviceIDPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request updateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	uniqueID := strings.TrimSpace(request.UniqueID)
	name := strings.TrimSpace(request.Name)

	if uniqueID == "" || name == "" {
		http.Error(w, "unique_id and name must not be empty", http.StatusBadRequest)
		return
	}

	dev, err := h.repo.UpdateDevice(r.Context(), device.Device{ID: id, UniqueID: uniqueID, Name: name})
	if err != nil {
		if errors.Is(err, device.ErrNotFound) {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		log.Printf("failed to update device %d: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, deviceResponse{
		ID:       dev.ID,
		UniqueID: dev.UniqueID,
		Name:     dev.Name,
	})
}

func (h *deviceHandler) deleteDevice(w http.ResponseWriter, r *http.Request) {
	id, err := parseDeviceIDPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteDeviceByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, device.ErrNotFound) {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		log.Printf("failed to delete device %d: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
