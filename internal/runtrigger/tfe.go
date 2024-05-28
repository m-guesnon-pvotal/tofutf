package runtrigger

import (
	"github.com/hashicorp/jsonapi"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/http/decode"
	"github.com/tofutf/tofutf/internal/tfeapi"

	types "github.com/hashicorp/go-tfe"
)

type (
	tfe struct {
		*Service
		*tfeapi.Responder
	}
)

func (a *tfe) addHandlers(r *mux.Router) {
	r = r.PathPrefix(tfeapi.APIPrefixV2).Subrouter()

	r.HandleFunc("/run-triggers/{run_trigger_id}", a.getRunTrigger).Methods("GET")
	r.HandleFunc("/run-triggers/{run_trigger_id}", a.deleteRunTrigger).Methods("DELETE")

	r.HandleFunc("/workspaces/{workspace_id}/run-triggers", a.listRunTriggers).Methods("GET")
	r.HandleFunc("/workspaces/{workspace_id}/run-triggers", a.addRunTrigger).Methods("POST")
}

func (a *tfe) toRunTrigger(from *RunTrigger) *types.RunTrigger {
	to := &types.RunTrigger{
		ID:        from.RunTriggerID,
		CreatedAt: from.CreatedAt,
		Workspace: &types.Workspace{
			ID: from.WorkspaceID,
		},
		SourceableChoice: &types.SourceableChoice{
			Workspace: &types.Workspace{
				ID: from.SourceableID,
			},
		},
	}
	return to
}

// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run-triggers#show-a-run-trigger
func (a *tfe) getRunTrigger(w http.ResponseWriter, r *http.Request) {
	var params struct {
		RunTriggerID string `schema:"run_trigger_id,required"`
	}
	if err := decode.Route(&params, r); err != nil {
		tfeapi.Error(w, err)
		return
	}
	rt, err := a.GetRunTrigger(r.Context(), params.RunTriggerID)
	if err != nil {
		tfeapi.Error(w, &internal.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	err = jsonapi.MarshalPayload(w, a.toRunTrigger(rt))
	if err != nil {
		tfeapi.Error(w, &internal.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run-triggers#list-run-triggers
func (a *tfe) listRunTriggers(w http.ResponseWriter, r *http.Request) {
	var params struct {
		WorkspaceID string `schema:"workspace_id,required"`
	}
	if err := decode.Route(&params, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	rts, err := a.ListByWorkspaceID(r.Context(), params.WorkspaceID)
	if err != nil {
		tfeapi.Error(w, &internal.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}
	resp := types.RunTriggerList{
		Items:      make([]*types.RunTrigger, 0),
		Pagination: &types.Pagination{},
	}

	for _, rt := range rts {
		resp.Items = append(resp.Items, a.toRunTrigger(rt))
	}

	err = jsonapi.MarshalPayload(w, resp)
	if err != nil {
		tfeapi.Error(w, &internal.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run-triggers#create-a-run-trigger
func (a *tfe) addRunTrigger(w http.ResponseWriter, r *http.Request) {
	var params struct {
		WorkspaceID string `schema:"workspace_id,required"`
	}
	if err := decode.Route(&params, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	var opts types.RunTriggerCreateOptions
	if err := jsonapi.UnmarshalPayload(r.Body, &opts); err != nil {
		tfeapi.Error(w, err)
		return
	}

	createOpts := CreateOptions{}

	if opts.Sourceable != nil {
		createOpts.SourceableID = opts.Sourceable.ID
		createOpts.SourceableType = WorkspaceSourceable
		if rt, err := a.Create(r.Context(), params.WorkspaceID, createOpts); err != nil {
			tfeapi.Error(w, err)
			return
		} else {
			resp := a.toRunTrigger(rt)
			err = jsonapi.MarshalPayload(w, resp)
			if err != nil {

				tfeapi.Error(w, &internal.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
				return
			}

			w.WriteHeader(http.StatusCreated)
			return
		}

	}
	tfeapi.Error(w, &internal.HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Missing sourceable",
	})
}

// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run-triggers#delete-a-run-trigger
func (a *tfe) deleteRunTrigger(w http.ResponseWriter, r *http.Request) {
	var params struct {
		RunTriggerID string `schema:"run_trigger_id,required"`
	}
	if err := decode.Route(&params, r); err != nil {
		tfeapi.Error(w, err)
		return
	}
	if err := a.Delete(r.Context(), params.RunTriggerID); err != nil {
		tfeapi.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
