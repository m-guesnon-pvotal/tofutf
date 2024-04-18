// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ genericConn = (*pgx.Conn)(nil)

const insertRepoConnectionSQL = `INSERT INTO repo_connections (
    vcs_provider_id,
    repo_path,
    workspace_id,
    module_id
) VALUES (
    $1,
    $2,
    $3,
    $4
);`

type InsertRepoConnectionParams struct {
	VCSProviderID pgtype.Text `json:"vcs_provider_id"`
	RepoPath      pgtype.Text `json:"repo_path"`
	WorkspaceID   pgtype.Text `json:"workspace_id"`
	ModuleID      pgtype.Text `json:"module_id"`
}

// InsertRepoConnection implements Querier.InsertRepoConnection.
func (q *DBQuerier) InsertRepoConnection(ctx context.Context, params InsertRepoConnectionParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertRepoConnection")
	cmdTag, err := q.conn.Exec(ctx, insertRepoConnectionSQL, params.VCSProviderID, params.RepoPath, params.WorkspaceID, params.ModuleID)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertRepoConnection: %w", err)
	}
	return cmdTag, err
}

const deleteWorkspaceConnectionByIDSQL = `DELETE
FROM repo_connections
WHERE workspace_id = $1
RETURNING *;`

type DeleteWorkspaceConnectionByIDRow struct {
	ModuleID      pgtype.Text `json:"module_id"`
	WorkspaceID   pgtype.Text `json:"workspace_id"`
	RepoPath      pgtype.Text `json:"repo_path"`
	VCSProviderID pgtype.Text `json:"vcs_provider_id"`
}

// DeleteWorkspaceConnectionByID implements Querier.DeleteWorkspaceConnectionByID.
func (q *DBQuerier) DeleteWorkspaceConnectionByID(ctx context.Context, workspaceID pgtype.Text) (DeleteWorkspaceConnectionByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteWorkspaceConnectionByID")
	rows, err := q.conn.Query(ctx, deleteWorkspaceConnectionByIDSQL, workspaceID)
	if err != nil {
		return DeleteWorkspaceConnectionByIDRow{}, fmt.Errorf("query DeleteWorkspaceConnectionByID: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (DeleteWorkspaceConnectionByIDRow, error) {
		var item DeleteWorkspaceConnectionByIDRow
		if err := row.Scan(&item.ModuleID, // 'module_id', 'ModuleID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.WorkspaceID,   // 'workspace_id', 'WorkspaceID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.RepoPath,      // 'repo_path', 'RepoPath', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.VCSProviderID, // 'vcs_provider_id', 'VCSProviderID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const deleteModuleConnectionByIDSQL = `DELETE
FROM repo_connections
WHERE module_id = $1
RETURNING *;`

type DeleteModuleConnectionByIDRow struct {
	ModuleID      pgtype.Text `json:"module_id"`
	WorkspaceID   pgtype.Text `json:"workspace_id"`
	RepoPath      pgtype.Text `json:"repo_path"`
	VCSProviderID pgtype.Text `json:"vcs_provider_id"`
}

// DeleteModuleConnectionByID implements Querier.DeleteModuleConnectionByID.
func (q *DBQuerier) DeleteModuleConnectionByID(ctx context.Context, moduleID pgtype.Text) (DeleteModuleConnectionByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteModuleConnectionByID")
	rows, err := q.conn.Query(ctx, deleteModuleConnectionByIDSQL, moduleID)
	if err != nil {
		return DeleteModuleConnectionByIDRow{}, fmt.Errorf("query DeleteModuleConnectionByID: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (DeleteModuleConnectionByIDRow, error) {
		var item DeleteModuleConnectionByIDRow
		if err := row.Scan(&item.ModuleID, // 'module_id', 'ModuleID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.WorkspaceID,   // 'workspace_id', 'WorkspaceID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.RepoPath,      // 'repo_path', 'RepoPath', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.VCSProviderID, // 'vcs_provider_id', 'VCSProviderID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
