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

	positionHandler := newPositionHandler(s.positionService)

	// GET /positions?device_id=1&from=<date>&to=<date>
	// GET /positions?device_id=1&device_id=2&from=<date>&to=<date>
	mux.HandleFunc("GET /positions", positionHandler.devicePositions)
	// GET /positions/latest
	mux.HandleFunc("GET /positions/latest", positionHandler.latestPositions)

	authHandler := newAuthHandler(s.authService)

	mux.HandleFunc("POST /auth/init", authHandler.init)
	mux.HandleFunc("POST /auth/login", authHandler.login)
}
