package runtrigger

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tofutf/tofutf/internal/run"
	"github.com/tofutf/tofutf/internal/workspace"

	"github.com/gorilla/mux"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/http/html"
	"github.com/tofutf/tofutf/internal/organization"
	"github.com/tofutf/tofutf/internal/rbac"
	"github.com/tofutf/tofutf/internal/sql"
	"github.com/tofutf/tofutf/internal/tfeapi"
)

type (
	Service struct {
		site         internal.Authorizer // authorizes site access
		organization internal.Authorizer // authorizes org access

		WorkspaceService *workspace.Service
		RunService       *run.Service
		Logger           *slog.Logger

		logger *slog.Logger
		db     *pgdb
		tfeapi *tfe
		api    *api
	}

	Options struct {
		SiteToken        string
		WorkspaceService *workspace.Service
		RunService       *run.Service
		Logger           *slog.Logger

		*sql.DB
		*tfeapi.Responder
		html.Renderer
	}
)

func NewService(opts Options) *Service {
	svc := Service{
		logger:           opts.Logger,
		organization:     &organization.Authorizer{Logger: opts.Logger},
		site:             &internal.SiteAuthorizer{Logger: opts.Logger},
		db:               &pgdb{opts.DB, opts.Logger},
		RunService:       opts.RunService,
		WorkspaceService: opts.WorkspaceService,
	}

	svc.tfeapi = &tfe{
		Service:   &svc,
		Responder: opts.Responder,
	}
	svc.api = &api{
		Service:   &svc,
		Responder: opts.Responder,
	}

	opts.RunService.AfterSuccessfulApply(func(ctx context.Context, r *run.Run) error {
		targets, err := svc.ListBySourceableID(ctx, r.WorkspaceID, WorkspaceSourceable)
		if err != nil {
			return err
		}

		for _, target := range targets {
			ctx := internal.AddSkipAuthz(ctx)
			msg := fmt.Sprintf("Triggered by %s", target.SourceableID)
			_, err := opts.RunService.Create(ctx, target.WorkspaceID, run.CreateOptions{
				Message: &msg,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return &svc
}

func (a *Service) AddHandlers(r *mux.Router) {
	a.tfeapi.addHandlers(r)
	a.api.addHandlers(r)
}

func (a *Service) Create(ctx context.Context, workspaceID string, opts CreateOptions) (*RunTrigger, error) {
	subject, err := a.site.CanAccess(ctx, rbac.UpdateWorkspaceAction, "")
	if err != nil {
		return nil, err
	}

	rt := NewRunTrigger(workspaceID, opts)

	if err := a.db.CreateRunTrigger(ctx, rt); err != nil {
		a.logger.Error("creating run trigger", "workspace", workspaceID, "subject", subject, "err", err)
		return nil, err
	}

	a.logger.Info("created run trigger", "workspace", workspaceID, "subject", subject)

	return rt, nil
}

func (a *Service) GetRunTrigger(ctx context.Context, rtID string) (*RunTrigger, error) {
	subject, err := a.site.CanAccess(ctx, rbac.GetUserAction, "")
	if err != nil {
		return nil, err
	}

	user, err := a.db.GetRunTrigger(ctx, rtID)
	if err != nil {
		a.logger.Debug("retrieving run trigger", "run trigger ID", rtID, "subject", subject, "err", err)
		return nil, err
	}

	a.logger.Debug("retrieved run trigger", "run trigger ID", rtID, "subject", subject)

	return user, nil
}

// ListByWorkspaceID lists all run triggers for a given workspace.
func (a *Service) ListByWorkspaceID(ctx context.Context, workspaceID string) ([]*RunTrigger, error) {
	_, err := a.site.CanAccess(ctx, rbac.UpdateWorkspaceAction, "")
	if err != nil {
		return nil, err
	}

	return a.db.ListRunTriggersByWorkspaceID(ctx, workspaceID)
}

// ListBySourceableID lists all run triggers for a given sourceable.
func (a *Service) ListBySourceableID(ctx context.Context, sourceableID string, sourceableType SourceableType) ([]*RunTrigger, error) {
	_, err := a.site.CanAccess(ctx, rbac.UpdateWorkspaceAction, "")
	if err != nil {
		return nil, err
	}

	return a.db.ListRunTriggersBySourceableID(ctx, sourceableID)
}

func (a *Service) Delete(ctx context.Context, rtID string) error {
	subject, err := a.site.CanAccess(ctx, rbac.UpdateWorkspaceAction, "")
	if err != nil {
		return err
	}

	err = a.db.DeleteRunTrigger(ctx, rtID)
	if err != nil {
		a.logger.Error("deleting run trigger", "run trigger ID", rtID, "subject", subject, "err", err)
		return err
	}

	a.logger.Info("deleted run trigger", "run trigger ID", rtID, "subject", subject)

	return nil
}
