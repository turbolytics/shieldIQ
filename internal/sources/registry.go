package sources

import (
	"net/http"

	"github.com/turbolytics/sqlsec/internal/sources/github"
)

type Validator interface {
	Validate(r *http.Request, secret string) error
}

type Parser interface {
	Parse(r *http.Request) (map[string]any, error)
}

type Registry struct {
	sources    map[string]any
	validators map[string]Validator
	parsers    map[string]Parser
}

func New() *Registry {
	return &Registry{
		sources:    make(map[string]any),
		validators: make(map[string]Validator),
		parsers:    make(map[string]Parser),
	}
}

func (r *Registry) Init() {
	r.sources["github"] = struct{}{}
	r.validators["github"] = &github.GithubValidator{}
	r.parsers["github"] = &github.GithubParser{}
}

func (r *Registry) IsEnabled(source string) bool {
	_, ok := r.sources[source]
	return ok
}

func (r *Registry) GetValidator(source string) Validator {
	return r.validators[source]
}

func (r *Registry) GetParser(source string) Parser {
	return r.parsers[source]
}

var DefaultRegistry *Registry

func init() {
	DefaultRegistry = New()
	DefaultRegistry.Init()
}
