package healthcheck

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Warn("Invalid HTTP method for /healthz")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.ContentLength > 0 {
		log.Warn("Request to /healthz contains unexpected content")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := Check()
	if err != nil {
		log.WithError(err).Error("Health check failed")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	log.Info("Health check successful")
}