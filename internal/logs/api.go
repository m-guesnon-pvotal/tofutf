package logs

import (
	"bytes"
	"io"
	"net/http"

	otfapi "github.com/tofutf/tofutf/internal/api"

	"github.com/gorilla/mux"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/http/decode"
)

type api struct {
	internal.Verifier // for verifying upload url

	svc *Service
}

func (a *api) addHandlers(r *mux.Router) {
	// client is typically terraform-cli
	signed := r.PathPrefix("/signed/{signature.expiry}").Subrouter()
	signed.Use(internal.VerifySignedURL(a.Verifier))
	signed.HandleFunc("/runs/{run_id}/logs/{phase}", a.getLogs).Methods("GET")

	// client is typically otf-agent
	r = r.PathPrefix(otfapi.DefaultBasePath).Subrouter()
	r.HandleFunc("/runs/{run_id}/logs/{phase}", a.putLogs).Methods("PUT")
}

func (a *api) getLogs(w http.ResponseWriter, r *http.Request) {
	var opts internal.GetChunkOptions
	if err := decode.All(&opts, r); err != nil {
		a.svc.logger.Error("failed to decode get logs options", "err", err)
		otfapi.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}
	chunk, err := a.svc.GetChunk(r.Context(), opts)
	if err != nil {
		a.svc.logger.Error("failed to copy buffer", "err", err)
		otfapi.HandleError(w, err, http.StatusNotFound)
		return
	}
	if _, err := w.Write(chunk.Data); err != nil {
		a.svc.logger.Error("failed to put logs", "err", err)
		otfapi.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (a *api) putLogs(w http.ResponseWriter, r *http.Request) {
	var opts internal.PutChunkOptions
	if err := decode.All(&opts, r); err != nil {
		a.svc.logger.Error("failed to decode put logs options", "err", err)
		otfapi.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r.Body); err != nil {
		a.svc.logger.Error("failed to copy buffer", "err", err)
		otfapi.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}
	opts.Data = buf.Bytes()
	if err := a.svc.PutChunk(r.Context(), opts); err != nil {
		a.svc.logger.Error("failed to put logs", "err", err)
		otfapi.HandleError(w, err, http.StatusNotFound)
		return
	}
}
