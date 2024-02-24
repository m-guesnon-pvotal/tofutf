package workspace

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/rbac"
)

// authorizer authorizes access to a workspace
type authorizer struct {
	logr.Logger

	db *pgdb
}

func (a *authorizer) CanAccess(ctx context.Context, action rbac.Action, workspaceID string) (internal.Subject, error) {
	subj, err := internal.SubjectFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if internal.SkipAuthz(ctx) {
		return subj, nil
	}
	policy, err := a.db.GetWorkspacePolicy(ctx, workspaceID)
	if err != nil {
		return nil, internal.ErrResourceNotFound
	}
	if subj.CanAccessWorkspace(action, policy) {
		return subj, nil
	}
	a.Error(nil, "unauthorized action", "workspace_id", workspaceID, "organization", policy.Organization, "action", action.String(), "subject", subj)
	return nil, internal.ErrAccessNotPermitted
}
