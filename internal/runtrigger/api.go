package runtrigger

import (
	"github.com/gorilla/mux"
	otfapi "github.com/tofutf/tofutf/internal/api"
	"github.com/tofutf/tofutf/internal/tfeapi"
)

type (
	api struct {
		*Service
		*tfeapi.Responder
	}
)

func (a *api) addHandlers(r *mux.Router) {
	r = r.PathPrefix(otfapi.DefaultBasePath).Subrouter()

}
