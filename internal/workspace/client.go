package workspace

import (
	"context"
	"fmt"
	"net/url"

	types "github.com/hashicorp/go-tfe"
	otfapi "github.com/tofutf/tofutf/internal/api"
	"github.com/tofutf/tofutf/internal/resource"
)

type Client struct {
	*otfapi.Client
}

func (c *Client) GetByName(ctx context.Context, organization, workspace string) (*types.Workspace, error) {
	path := fmt.Sprintf("organizations/%s/workspaces/%s", organization, workspace)
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var ws types.Workspace
	if err := c.Do(ctx, req, &ws); err != nil {
		return nil, err
	}
	return &ws, nil
}

func (c *Client) Get(ctx context.Context, workspaceID string) (*types.Workspace, error) {
	path := fmt.Sprintf("workspaces/%s", workspaceID)
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var ws types.Workspace
	if err := c.Do(ctx, req, &ws); err != nil {
		return nil, err
	}
	return &ws, nil
}

func (c *Client) List(ctx context.Context, opts ListOptions) (*resource.Page[*types.Workspace], error) {
	u := fmt.Sprintf("organizations/%s/workspaces", url.QueryEscape(*opts.Organization))
	req, err := c.NewRequest("GET", u, &opts)
	if err != nil {
		return nil, err
	}
	var page resource.Page[*types.Workspace]
	if err = c.Do(ctx, req, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

func (c *Client) Update(ctx context.Context, workspaceID string, opts types.WorkspaceUpdateOptions) (*types.Workspace, error) {
	// Pre-emptively validate options
	if _, err := Update(&types.Workspace{}, opts); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("workspaces/%s", workspaceID)
	req, err := c.NewRequest("PATCH", path, &opts)
	if err != nil {
		return nil, err
	}

	var ws types.Workspace
	if err := c.Do(ctx, req, &ws); err != nil {
		return nil, err
	}

	return &ws, nil
}

func (c *Client) Lock(ctx context.Context, workspaceID string, runID *string) (*types.Workspace, error) {
	path := fmt.Sprintf("workspaces/%s/actions/lock", workspaceID)
	req, err := c.NewRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	var ws types.Workspace
	if err := c.Do(ctx, req, &ws); err != nil {
		return nil, err
	}

	return &ws, nil
}

func (c *Client) Unlock(ctx context.Context, workspaceID string, runID *string, force bool) (*types.Workspace, error) {
	var u string
	if force {
		u = fmt.Sprintf("workspaces/%s/actions/unlock", workspaceID)
	} else {
		u = fmt.Sprintf("workspaces/%s/actions/force-unlock", workspaceID)
	}
	req, err := c.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	var ws types.Workspace
	if err := c.Do(ctx, req, &ws); err != nil {
		return nil, err
	}

	return &ws, nil
}
