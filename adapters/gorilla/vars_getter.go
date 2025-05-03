package gorilla

import (
	"net/http"

	"github.com/gorilla/mux"
)

type MuxVarsGetter struct{}

func (m *MuxVarsGetter) GetVar(r *http.Request, key string) (string, bool) {
	vars := mux.Vars(r)
	v, ok := vars[key]
	return v, ok
}

func New() *MuxVarsGetter {
	return &MuxVarsGetter{}
}
