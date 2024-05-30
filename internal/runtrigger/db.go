package runtrigger

import (
	"context"
	"log/slog"

	"github.com/jackc/pgtype"
	"github.com/tofutf/tofutf/internal/sql"
	"github.com/tofutf/tofutf/internal/sql/pggen"
)

// dbresult represents the result of a database query for a user.
type dbresult struct {
	RunTriggerID pgtype.Text        `json:"run_trigger_id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	SourceableID pgtype.Text        `json:"sourceable_id"`
	WorkspaceID  pgtype.Text        `json:"workspace_id"`
}

func (result dbresult) toRunTrigger() *RunTrigger {
	rt := RunTrigger{
		RunTriggerID: result.RunTriggerID.String,
		CreatedAt:    result.CreatedAt.Time,
		SourceableID: result.SourceableID.String,
		WorkspaceID:  result.WorkspaceID.String,
	}
	return &rt
}

// pgdb stores user resources in a postgres database
type pgdb struct {
	*sql.DB // provides access to generated SQL queries
	Logger  *slog.Logger
}

// CreateRunTrigger persists a RunTrigger to the DB.
func (db *pgdb) CreateRunTrigger(ctx context.Context, rt *RunTrigger) error {
	return db.Tx(ctx, func(ctx context.Context, q pggen.Querier) error {
		_, err := q.InsertRunTrigger(ctx, pggen.InsertRunTriggerParams{
			RunTriggerID: sql.String(rt.RunTriggerID),
			CreatedAt:    sql.Timestamptz(rt.CreatedAt),
			SourceableID: sql.String(rt.SourceableID),
			WorkspaceID:  sql.String(rt.WorkspaceID),
		})
		if err != nil {
			return sql.Error(err)
		}
		return nil
	})
}

func (db *pgdb) GetRunTrigger(ctx context.Context, rtID string) (*RunTrigger, error) {
	result, err := db.Conn(ctx).FindRunTriggerByID(ctx, sql.String(rtID))
	if err != nil {
		return nil, err
	}
	return dbresult(result).toRunTrigger(), nil
}

// DeleteRunTrigger persists a RunTrigger to the DB.
func (db *pgdb) DeleteRunTrigger(ctx context.Context, rtID string) error {
	_, err := db.Conn(ctx).DeleteRunTriggerByID(ctx, sql.String(rtID))
	if err != nil {
		return sql.Error(err)
	}
	return nil
}

func (db *pgdb) ListRunTriggersByWorkspaceID(ctx context.Context, workspaceID string) ([]*RunTrigger, error) {
	result, err := db.Conn(ctx).FindRunTriggersByWorkspaceID(ctx, sql.String(workspaceID))
	if err != nil {
		return nil, err
	}
	rts := make([]*RunTrigger, len(result))
	for i, r := range result {
		rts[i] = dbresult(r).toRunTrigger()
	}
	return rts, nil
}

func (db *pgdb) ListRunTriggersBySourceableID(ctx context.Context, sourceableID string) ([]*RunTrigger, error) {
	result, err := db.Conn(ctx).FindRunTriggersBySourceableID(ctx, sql.String(sourceableID))
	if err != nil {
		return nil, err
	}
	rts := make([]*RunTrigger, len(result))
	for i, r := range result {
		rts[i] = dbresult(r).toRunTrigger()
	}
	return rts, nil
}
