package organization

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/rbac"
)

// Authorizer authorizes access to an organization
type Authorizer struct {
	logr.Logger
}

func (a *Authorizer) CanAccess(ctx context.Context, action rbac.Action, name string) (internal.Subject, error) {
	subj, err := internal.SubjectFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if internal.SkipAuthz(ctx) {
		return subj, nil
	}
	if subj.CanAccessOrganization(action, name) {
		return subj, nil
	}
	a.Error(nil, "unauthorized action", "organization", name, "action", action.String(), "subject", subj)
	return nil, internal.ErrAccessNotPermitted
}
