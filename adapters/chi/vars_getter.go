package chi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ChiVarsGetter struct{}

func (c *ChiVarsGetter) GetVar(r *http.Request, key string) (string, bool) {
	v := chi.URLParam(r, key)
	return v, v != ""
}

func New() *ChiVarsGetter {
	return &ChiVarsGetter{}
}
