package api

import "net/http"

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

	positionHandler := newPositionHandler(s.positionRepo, s.deviceRepo)

	mux.HandleFunc("GET /positions", positionHandler.getPositions)
	mux.HandleFunc("GET /positions/latest", positionHandler.getLatestPosition)
}
