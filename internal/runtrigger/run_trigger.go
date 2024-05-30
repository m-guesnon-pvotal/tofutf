package runtrigger

import (
	"github.com/tofutf/tofutf/internal"
	"time"
)

const (
	WorkspaceSourceable SourceableType = "workspaces"
)

type (
	SourceableType string
	// RunTrigger spawns new runs in response to other runs completing
	RunTrigger struct {
		RunTriggerID string    `json:"run_trigger_id"`
		CreatedAt    time.Time `json:"created_at"`
		SourceableID string    `json:"sourceable_id"`
		WorkspaceID  string    `json:"workspace_id"`
	}

	CreateOptions struct {
		SourceableID   string
		SourceableType SourceableType
	}
)

func NewRunTrigger(workspaceID string, opts CreateOptions) *RunTrigger {

	// Unused SourceableType. This is for later support
	rt := &RunTrigger{
		RunTriggerID: internal.NewID("rt"),
		CreatedAt:    internal.CurrentTimestamp(nil),
		WorkspaceID:  workspaceID,
		SourceableID: opts.SourceableID,
	}
	return rt
}
