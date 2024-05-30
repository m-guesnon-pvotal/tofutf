// Package ghapphandler provides a handler for the github app webhook endpoint.
package ghapphandler

import (
	"errors"
	"github.com/tofutf/tofutf/internal/api"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/github"
	"github.com/tofutf/tofutf/internal/vcs"
	"github.com/tofutf/tofutf/internal/vcsprovider"
)

type Handler struct {
	Logger *slog.Logger
	vcs.Publisher

	VCSProviders *vcsprovider.Service
	GithubApps   *github.Service
}

func (h *Handler) AddHandlers(r *mux.Router) {
	r.Handle(github.AppEventsPath, h)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// permit handler to talk to services
	ctx := internal.AddSubjectToContext(r.Context(), &internal.Superuser{Username: "github-app-event-handler"})
	// retrieve github app config; if one hasn't been configured then return a
	// 400
	app, err := h.GithubApps.GetApp(ctx)
	if err != nil {
		api.HandleError(w, err, http.StatusBadRequest)
		return
	}
	h.Logger.Debug("received vcs event", "github_app", app)

	// use github-specific handler to unmarshal event
	payload, err := github.HandleEvent(r, app.WebhookSecret)
	// either ignore the event, return an error, or publish the event onwards
	var ignore vcs.ErrIgnoreEvent
	if errors.As(err, &ignore) {
		h.Logger.Info("ignoring event: "+err.Error(), "github_app", app)
		w.WriteHeader(http.StatusOK)
		return
	} else if err != nil {
		h.Logger.Error("handling vcs event", "github_app", app, "err", err)
		api.HandleError(w, err, http.StatusBadRequest)
		return
	}
	// relay a copy of the event for each vcs provider configured with the
	// github app install that triggered the event.
	providers, err := h.VCSProviders.ListVCSProvidersByGithubAppInstall(ctx, *payload.GithubAppInstallID)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	for _, prov := range providers {
		h.Publish(vcs.Event{
			EventHeader:  vcs.EventHeader{VCSProviderID: prov.ID},
			EventPayload: *payload,
		})
	}
	w.WriteHeader(http.StatusOK)
}
