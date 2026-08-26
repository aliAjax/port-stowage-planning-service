package observability

import "sync"

// Registry tracks the health of named components.
type Registry struct {
	mu         sync.RWMutex
	components map[string]string
}

func NewRegistry() *Registry {
	return &Registry{components: map[string]string{}}
}

// Set records a component's health, overwriting any prior status.
func (r *Registry) Set(name, status string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.components[name] = status
}

// Overall reports the aggregate status: "ok" when every component is
// healthy, "degraded" when at least one is not, and "unknown" when no
// component has been registered.
func (r *Registry) Overall() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.components) == 0 {
		return "unknown"
	}
	for _, st := range r.components {
		if st != "ok" {
			return "degraded"
		}
	}
	return "ok"
}

// Components returns a copy of the current component map.
func (r *Registry) Components() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]string, len(r.components))
	for k, v := range r.components {
		out[k] = v
	}
	return out
}
