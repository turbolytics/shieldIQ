package sources

type Registry struct {
	sources map[string]any
}

func New() *Registry {
	r := &Registry{
		sources: map[string]any{
			"github": struct{}{},
		},
	}

	return r
}

func (r *Registry) IsEnabled(source string) bool {
	_, ok := r.sources[source]
	return ok
}

var DefaultRegistry *Registry

func init() {
	DefaultRegistry = New()
}
