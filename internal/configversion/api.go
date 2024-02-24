package configversion

import (
	"net/http"

	"github.com/gorilla/mux"
	otfapi "github.com/tofutf/tofutf/internal/api"
	"github.com/tofutf/tofutf/internal/http/decode"
	"github.com/tofutf/tofutf/internal/tfeapi"
)

type api struct {
	*Service
	*tfeapi.Responder
}

func (a *api) addHandlers(r *mux.Router) {
	r = r.PathPrefix(otfapi.DefaultBasePath).Subrouter()
	r.HandleFunc("/configuration-versions/{id}/download", a.download).Methods("GET")
}

func (a *api) download(w http.ResponseWriter, r *http.Request) {
	id, err := decode.Param("id", r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	resp, err := a.DownloadConfig(r.Context(), id)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	w.Write(resp)
}
