package observability

import (
	"encoding/json"
	"net/http"
	"time"
)

// Health reports the overall status, always ok.
func Health(reg *Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "time": time.Now().UTC()})
	}
}

// Ready reports readiness.
func Ready(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
