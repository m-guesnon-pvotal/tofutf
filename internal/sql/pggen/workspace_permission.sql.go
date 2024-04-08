// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const upsertWorkspacePermissionSQL = `INSERT INTO workspace_permissions (
    workspace_id,
    team_id,
    role
) VALUES (
    $1,
    $2,
    $3
) ON CONFLICT (workspace_id, team_id) DO UPDATE SET role = $3;`

type UpsertWorkspacePermissionParams struct {
	WorkspaceID pgtype.Text `json:"workspace_id"`
	TeamID      pgtype.Text `json:"team_id"`
	Role        pgtype.Text `json:"role"`
}

// UpsertWorkspacePermission implements Querier.UpsertWorkspacePermission.
func (q *DBQuerier) UpsertWorkspacePermission(ctx context.Context, params UpsertWorkspacePermissionParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpsertWorkspacePermission")
	cmdTag, err := q.conn.Exec(ctx, upsertWorkspacePermissionSQL, params.WorkspaceID, params.TeamID, params.Role)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query UpsertWorkspacePermission: %w", err)
	}
	return cmdTag, err
}

// UpsertWorkspacePermissionBatch implements Querier.UpsertWorkspacePermissionBatch.
func (q *DBQuerier) UpsertWorkspacePermissionBatch(batch genericBatch, params UpsertWorkspacePermissionParams) {
	batch.Queue(upsertWorkspacePermissionSQL, params.WorkspaceID, params.TeamID, params.Role)
}

// UpsertWorkspacePermissionScan implements Querier.UpsertWorkspacePermissionScan.
func (q *DBQuerier) UpsertWorkspacePermissionScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec UpsertWorkspacePermissionBatch: %w", err)
	}
	return cmdTag, err
}

const findWorkspacePermissionsByWorkspaceIDSQL = `SELECT *
FROM workspace_permissions
WHERE workspace_id = $1;`

type FindWorkspacePermissionsByWorkspaceIDRow struct {
	WorkspaceID pgtype.Text `json:"workspace_id"`
	TeamID      pgtype.Text `json:"team_id"`
	Role        pgtype.Text `json:"role"`
}

// FindWorkspacePermissionsByWorkspaceID implements Querier.FindWorkspacePermissionsByWorkspaceID.
func (q *DBQuerier) FindWorkspacePermissionsByWorkspaceID(ctx context.Context, workspaceID pgtype.Text) ([]FindWorkspacePermissionsByWorkspaceIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindWorkspacePermissionsByWorkspaceID")
	rows, err := q.conn.Query(ctx, findWorkspacePermissionsByWorkspaceIDSQL, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("query FindWorkspacePermissionsByWorkspaceID: %w", err)
	}
	defer rows.Close()
	items := []FindWorkspacePermissionsByWorkspaceIDRow{}
	for rows.Next() {
		var item FindWorkspacePermissionsByWorkspaceIDRow
		if err := rows.Scan(&item.WorkspaceID, &item.TeamID, &item.Role); err != nil {
			return nil, fmt.Errorf("scan FindWorkspacePermissionsByWorkspaceID row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindWorkspacePermissionsByWorkspaceID rows: %w", err)
	}
	return items, err
}

// FindWorkspacePermissionsByWorkspaceIDBatch implements Querier.FindWorkspacePermissionsByWorkspaceIDBatch.
func (q *DBQuerier) FindWorkspacePermissionsByWorkspaceIDBatch(batch genericBatch, workspaceID pgtype.Text) {
	batch.Queue(findWorkspacePermissionsByWorkspaceIDSQL, workspaceID)
}

// FindWorkspacePermissionsByWorkspaceIDScan implements Querier.FindWorkspacePermissionsByWorkspaceIDScan.
func (q *DBQuerier) FindWorkspacePermissionsByWorkspaceIDScan(results pgx.BatchResults) ([]FindWorkspacePermissionsByWorkspaceIDRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindWorkspacePermissionsByWorkspaceIDBatch: %w", err)
	}
	defer rows.Close()
	items := []FindWorkspacePermissionsByWorkspaceIDRow{}
	for rows.Next() {
		var item FindWorkspacePermissionsByWorkspaceIDRow
		if err := rows.Scan(&item.WorkspaceID, &item.TeamID, &item.Role); err != nil {
			return nil, fmt.Errorf("scan FindWorkspacePermissionsByWorkspaceIDBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindWorkspacePermissionsByWorkspaceIDBatch rows: %w", err)
	}
	return items, err
}

const deleteWorkspacePermissionByIDSQL = `DELETE
FROM workspace_permissions
WHERE workspace_id = $1
AND team_id = $2;`

// DeleteWorkspacePermissionByID implements Querier.DeleteWorkspacePermissionByID.
func (q *DBQuerier) DeleteWorkspacePermissionByID(ctx context.Context, workspaceID pgtype.Text, teamID pgtype.Text) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteWorkspacePermissionByID")
	cmdTag, err := q.conn.Exec(ctx, deleteWorkspacePermissionByIDSQL, workspaceID, teamID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query DeleteWorkspacePermissionByID: %w", err)
	}
	return cmdTag, err
}

// DeleteWorkspacePermissionByIDBatch implements Querier.DeleteWorkspacePermissionByIDBatch.
func (q *DBQuerier) DeleteWorkspacePermissionByIDBatch(batch genericBatch, workspaceID pgtype.Text, teamID pgtype.Text) {
	batch.Queue(deleteWorkspacePermissionByIDSQL, workspaceID, teamID)
}

// DeleteWorkspacePermissionByIDScan implements Querier.DeleteWorkspacePermissionByIDScan.
func (q *DBQuerier) DeleteWorkspacePermissionByIDScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec DeleteWorkspacePermissionByIDBatch: %w", err)
	}
	return cmdTag, err
}
