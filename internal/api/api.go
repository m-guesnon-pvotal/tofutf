// Package api provides commmon functionality for the OTF API
package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

const (
	DefaultBasePath = "/otfapi"
	PingEndpoint    = "ping"
	DefaultAddress  = "localhost:8080"
)

type Handlers struct{}

func (h *Handlers) AddHandlers(r *mux.Router) {
	// basic no-op ping handler
	r.HandleFunc(path.Join(DefaultBasePath, PingEndpoint), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleError(w http.ResponseWriter, err error, code int) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf("failed to handle request"), pcs[0])
	r.Add("err", err)
	_ = slog.Default().Handler().Handle(context.Background(), r)
	http.Error(w, err.Error(), code)
}
