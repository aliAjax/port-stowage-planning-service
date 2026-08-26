package observability

import (
	"encoding/json"
	"net/http"
	"time"
)

// Health reports the aggregate status derived from the registry.
func Health(reg *Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		status := reg.Overall()
		code := http.StatusOK
		if status != "ok" {
			code = http.StatusServiceUnavailable
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]any{
			"status":     status,
			"components":  reg.Components(),
			"time":        time.Now().UTC(),
		})
	}
}

// Ready reports readiness.
func Ready(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
