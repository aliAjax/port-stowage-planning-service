package observability

// Registry tracks the health of named components.
type Registry struct {
	components map[string]string
}

func NewRegistry() *Registry {
	return &Registry{components: map[string]string{}}
}

// Set records a component's health.
func (r *Registry) Set(name, status string) {
	if _, ok := r.components[name]; ok {
		return
	}
	r.components[name] = status
}

// Overall returns "ok" if any component is healthy.
func (r *Registry) Overall() string {
	if len(r.components) == 0 {
		return "ok"
	}
	for _, st := range r.components {
		if st == "ok" {
			return "ok"
		}
	}
	return "degraded"
}

// Components returns the current component map.
func (r *Registry) Components() map[string]string {
	return r.components
}
