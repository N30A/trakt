package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

type deviceResponse struct {
	ID       int    `json:"id"`
	UniqueID string `json:"unique_id"`
	Name     string `json:"name"`
}

type positionResponse struct {
	ID         int       `json:"id"`
	DeviceID   int       `json:"device_id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	FixTime    time.Time `json:"fix_time"`
	ServerTime time.Time `json:"server_time"`
	Protocol   string    `json:"protocol"`
	Altitude   *float64  `json:"altitude"`
	Speed      *float64  `json:"speed"`
	Course     *float64  `json:"course"`
	Accuracy   *float64  `json:"accuracy"`
}
