package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func parseDeviceIDQuery(r *http.Request) (int, bool, error) {
	value := r.URL.Query().Get("deviceid")
	if value == "" {
		return 0, false, nil
	}
	id, err := strconv.Atoi(value)
	if err != nil || id <= 0 {
		return 0, true, errors.New("device id must be a positive integer")
	}
	return id, true, nil
}

func parseTimeRange(r *http.Request) (time.Time, time.Time, error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		return time.Time{}, time.Time{}, errors.New("from and to are required")
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
