// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertOrganizationSQL = `INSERT INTO organizations (
    organization_id,
    created_at,
    updated_at,
    name,
    email,
    collaborator_auth_policy,
    cost_estimation_enabled,
    session_remember,
    session_timeout,
    allow_force_delete_workspaces
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
);`

type InsertOrganizationParams struct {
	ID                         pgtype.Text        `json:"id"`
	CreatedAt                  pgtype.Timestamptz `json:"created_at"`
	UpdatedAt                  pgtype.Timestamptz `json:"updated_at"`
	Name                       pgtype.Text        `json:"name"`
	Email                      pgtype.Text        `json:"email"`
	CollaboratorAuthPolicy     pgtype.Text        `json:"collaborator_auth_policy"`
	CostEstimationEnabled      pgtype.Bool        `json:"cost_estimation_enabled"`
	SessionRemember            pgtype.Int4        `json:"session_remember"`
	SessionTimeout             pgtype.Int4        `json:"session_timeout"`
	AllowForceDeleteWorkspaces pgtype.Bool        `json:"allow_force_delete_workspaces"`
}

// InsertOrganization implements Querier.InsertOrganization.
func (q *DBQuerier) InsertOrganization(ctx context.Context, params InsertOrganizationParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertOrganization")
	cmdTag, err := q.conn.Exec(ctx, insertOrganizationSQL, params.ID, params.CreatedAt, params.UpdatedAt, params.Name, params.Email, params.CollaboratorAuthPolicy, params.CostEstimationEnabled, params.SessionRemember, params.SessionTimeout, params.AllowForceDeleteWorkspaces)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertOrganization: %w", err)
	}
	return cmdTag, err
}

// InsertOrganizationBatch implements Querier.InsertOrganizationBatch.
func (q *DBQuerier) InsertOrganizationBatch(batch genericBatch, params InsertOrganizationParams) {
	batch.Queue(insertOrganizationSQL, params.ID, params.CreatedAt, params.UpdatedAt, params.Name, params.Email, params.CollaboratorAuthPolicy, params.CostEstimationEnabled, params.SessionRemember, params.SessionTimeout, params.AllowForceDeleteWorkspaces)
}

// InsertOrganizationScan implements Querier.InsertOrganizationScan.
func (q *DBQuerier) InsertOrganizationScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertOrganizationBatch: %w", err)
	}
	return cmdTag, err
}

const findOrganizationNameByWorkspaceIDSQL = `SELECT organization_name
FROM workspaces
WHERE workspace_id = $1
;`

// FindOrganizationNameByWorkspaceID implements Querier.FindOrganizationNameByWorkspaceID.
func (q *DBQuerier) FindOrganizationNameByWorkspaceID(ctx context.Context, workspaceID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindOrganizationNameByWorkspaceID")
	row := q.conn.QueryRow(ctx, findOrganizationNameByWorkspaceIDSQL, workspaceID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query FindOrganizationNameByWorkspaceID: %w", err)
	}
	return item, nil
}

// FindOrganizationNameByWorkspaceIDBatch implements Querier.FindOrganizationNameByWorkspaceIDBatch.
func (q *DBQuerier) FindOrganizationNameByWorkspaceIDBatch(batch genericBatch, workspaceID pgtype.Text) {
	batch.Queue(findOrganizationNameByWorkspaceIDSQL, workspaceID)
}

// FindOrganizationNameByWorkspaceIDScan implements Querier.FindOrganizationNameByWorkspaceIDScan.
func (q *DBQuerier) FindOrganizationNameByWorkspaceIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan FindOrganizationNameByWorkspaceIDBatch row: %w", err)
	}
	return item, nil
}

const findOrganizationByNameSQL = `SELECT * FROM organizations WHERE name = $1;`

type FindOrganizationByNameRow struct {
	OrganizationID             pgtype.Text        `json:"organization_id"`
	CreatedAt                  pgtype.Timestamptz `json:"created_at"`
	UpdatedAt                  pgtype.Timestamptz `json:"updated_at"`
	Name                       pgtype.Text        `json:"name"`
	SessionRemember            pgtype.Int4        `json:"session_remember"`
	SessionTimeout             pgtype.Int4        `json:"session_timeout"`
	Email                      pgtype.Text        `json:"email"`
	CollaboratorAuthPolicy     pgtype.Text        `json:"collaborator_auth_policy"`
	AllowForceDeleteWorkspaces pgtype.Bool        `json:"allow_force_delete_workspaces"`
	CostEstimationEnabled      pgtype.Bool        `json:"cost_estimation_enabled"`
}

// FindOrganizationByName implements Querier.FindOrganizationByName.
func (q *DBQuerier) FindOrganizationByName(ctx context.Context, name pgtype.Text) (FindOrganizationByNameRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindOrganizationByName")
	row := q.conn.QueryRow(ctx, findOrganizationByNameSQL, name)
	var item FindOrganizationByNameRow
	if err := row.Scan(&item.OrganizationID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.SessionRemember, &item.SessionTimeout, &item.Email, &item.CollaboratorAuthPolicy, &item.AllowForceDeleteWorkspaces, &item.CostEstimationEnabled); err != nil {
		return item, fmt.Errorf("query FindOrganizationByName: %w", err)
	}
	return item, nil
}

// FindOrganizationByNameBatch implements Querier.FindOrganizationByNameBatch.
func (q *DBQuerier) FindOrganizationByNameBatch(batch genericBatch, name pgtype.Text) {
	batch.Queue(findOrganizationByNameSQL, name)
}

// FindOrganizationByNameScan implements Querier.FindOrganizationByNameScan.
func (q *DBQuerier) FindOrganizationByNameScan(results pgx.BatchResults) (FindOrganizationByNameRow, error) {
	row := results.QueryRow()
	var item FindOrganizationByNameRow
	if err := row.Scan(&item.OrganizationID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.SessionRemember, &item.SessionTimeout, &item.Email, &item.CollaboratorAuthPolicy, &item.AllowForceDeleteWorkspaces, &item.CostEstimationEnabled); err != nil {
		return item, fmt.Errorf("scan FindOrganizationByNameBatch row: %w", err)
	}
	return item, nil
}

const findOrganizationByIDSQL = `SELECT * FROM organizations WHERE organization_id = $1;`

type FindOrganizationByIDRow struct {
	OrganizationID             pgtype.Text        `json:"organization_id"`
	CreatedAt                  pgtype.Timestamptz `json:"created_at"`
	UpdatedAt                  pgtype.Timestamptz `json:"updated_at"`
	Name                       pgtype.Text        `json:"name"`
	SessionRemember            pgtype.Int4        `json:"session_remember"`
	SessionTimeout             pgtype.Int4        `json:"session_timeout"`
	Email                      pgtype.Text        `json:"email"`
	CollaboratorAuthPolicy     pgtype.Text        `json:"collaborator_auth_policy"`
	AllowForceDeleteWorkspaces pgtype.Bool        `json:"allow_force_delete_workspaces"`
	CostEstimationEnabled      pgtype.Bool        `json:"cost_estimation_enabled"`
}

// FindOrganizationByID implements Querier.FindOrganizationByID.
func (q *DBQuerier) FindOrganizationByID(ctx context.Context, organizationID pgtype.Text) (FindOrganizationByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindOrganizationByID")
	row := q.conn.QueryRow(ctx, findOrganizationByIDSQL, organizationID)
	var item FindOrganizationByIDRow
	if err := row.Scan(&item.OrganizationID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.SessionRemember, &item.SessionTimeout, &item.Email, &item.CollaboratorAuthPolicy, &item.AllowForceDeleteWorkspaces, &item.CostEstimationEnabled); err != nil {
		return item, fmt.Errorf("query FindOrganizationByID: %w", err)
	}
	return item, nil
}

// FindOrganizationByIDBatch implements Querier.FindOrganizationByIDBatch.
func (q *DBQuerier) FindOrganizationByIDBatch(batch genericBatch, organizationID pgtype.Text) {
	batch.Queue(findOrganizationByIDSQL, organizationID)
}

// FindOrganizationByIDScan implements Querier.FindOrganizationByIDScan.
func (q *DBQuerier) FindOrganizationByIDScan(results pgx.BatchResults) (FindOrganizationByIDRow, error) {
	row := results.QueryRow()
	var item FindOrganizationByIDRow
	if err := row.Scan(&item.OrganizationID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.SessionRemember, &item.SessionTimeout, &item.Email, &item.CollaboratorAuthPolicy, &item.AllowForceDeleteWorkspaces, &item.CostEstimationEnabled); err != nil {
		return item, fmt.Errorf("scan FindOrganizationByIDBatch row: %w", err)
	}
	return item, nil
}

const findOrganizationByNameForUpdateSQL = `SELECT *
FROM organizations
WHERE name = $1
FOR UPDATE
;`

type FindOrganizationByNameForUpdateRow struct {
	OrganizationID             pgtype.Text        `json:"organization_id"`
	CreatedAt                  pgtype.Timestamptz `json:"created_at"`
	UpdatedAt                  pgtype.Timestamptz `json:"updated_at"`
	Name                       pgtype.Text        `json:"name"`
	SessionRemember            pgtype.Int4        `json:"session_remember"`
	SessionTimeout             pgtype.Int4        `json:"session_timeout"`
	Email                      pgtype.Text        `json:"email"`
	CollaboratorAuthPolicy     pgtype.Text        `json:"collaborator_auth_policy"`
	AllowForceDeleteWorkspaces pgtype.Bool        `json:"allow_force_delete_workspaces"`
	CostEstimationEnabled      pgtype.Bool        `json:"cost_estimation_enabled"`
}

// FindOrganizationByNameForUpdate implements Querier.FindOrganizationByNameForUpdate.
func (q *DBQuerier) FindOrganizationByNameForUpdate(ctx context.Context, name pgtype.Text) (FindOrganizationByNameForUpdateRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindOrganizationByNameForUpdate")
	row := q.conn.QueryRow(ctx, findOrganizationByNameForUpdateSQL, name)
	var item FindOrganizationByNameForUpdateRow
	if err := row.Scan(&item.OrganizationID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.SessionRemember, &item.SessionTimeout, &item.Email, &item.CollaboratorAuthPolicy, &item.AllowForceDeleteWorkspaces, &item.CostEstimationEnabled); err != nil {
		return item, fmt.Errorf("query FindOrganizationByNameForUpdate: %w", err)
	}
	return item, nil
}

// FindOrganizationByNameForUpdateBatch implements Querier.FindOrganizationByNameForUpdateBatch.
func (q *DBQuerier) FindOrganizationByNameForUpdateBatch(batch genericBatch, name pgtype.Text) {
	batch.Queue(findOrganizationByNameForUpdateSQL, name)
}

// FindOrganizationByNameForUpdateScan implements Querier.FindOrganizationByNameForUpdateScan.
func (q *DBQuerier) FindOrganizationByNameForUpdateScan(results pgx.BatchResults) (FindOrganizationByNameForUpdateRow, error) {
	row := results.QueryRow()
	var item FindOrganizationByNameForUpdateRow
	if err := row.Scan(&item.OrganizationID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.SessionRemember, &item.SessionTimeout, &item.Email, &item.CollaboratorAuthPolicy, &item.AllowForceDeleteWorkspaces, &item.CostEstimationEnabled); err != nil {
		return item, fmt.Errorf("scan FindOrganizationByNameForUpdateBatch row: %w", err)
	}
	return item, nil
}

const findOrganizationsSQL = `SELECT *
FROM organizations
WHERE name LIKE ANY($1)
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3
;`

type FindOrganizationsParams struct {
	Names  []string    `json:"names"`
	Limit  pgtype.Int8 `json:"limit"`
	Offset pgtype.Int8 `json:"offset"`
}

type FindOrganizationsRow struct {
	OrganizationID             pgtype.Text        `json:"organization_id"`
	CreatedAt                  pgtype.Timestamptz `json:"created_at"`
	UpdatedAt                  pgtype.Timestamptz `json:"updated_at"`
	Name                       pgtype.Text        `json:"name"`
	SessionRemember            pgtype.Int4        `json:"session_remember"`
	SessionTimeout             pgtype.Int4        `json:"session_timeout"`
	Email                      pgtype.Text        `json:"email"`
	CollaboratorAuthPolicy     pgtype.Text        `json:"collaborator_auth_policy"`
	AllowForceDeleteWorkspaces pgtype.Bool        `json:"allow_force_delete_workspaces"`
	CostEstimationEnabled      pgtype.Bool        `json:"cost_estimation_enabled"`
}

// FindOrganizations implements Querier.FindOrganizations.
func (q *DBQuerier) FindOrganizations(ctx context.Context, params FindOrganizationsParams) ([]FindOrganizationsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindOrganizations")
	rows, err := q.conn.Query(ctx, findOrganizationsSQL, params.Names, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query FindOrganizations: %w", err)
	}
	defer rows.Close()
	items := []FindOrganizationsRow{}
	for rows.Next() {
		var item FindOrganizationsRow
		if err := rows.Scan(&item.OrganizationID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.SessionRemember, &item.SessionTimeout, &item.Email, &item.CollaboratorAuthPolicy, &item.AllowForceDeleteWorkspaces, &item.CostEstimationEnabled); err != nil {
			return nil, fmt.Errorf("scan FindOrganizations row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindOrganizations rows: %w", err)
	}
	return items, err
}

// FindOrganizationsBatch implements Querier.FindOrganizationsBatch.
func (q *DBQuerier) FindOrganizationsBatch(batch genericBatch, params FindOrganizationsParams) {
	batch.Queue(findOrganizationsSQL, params.Names, params.Limit, params.Offset)
}

// FindOrganizationsScan implements Querier.FindOrganizationsScan.
func (q *DBQuerier) FindOrganizationsScan(results pgx.BatchResults) ([]FindOrganizationsRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindOrganizationsBatch: %w", err)
	}
	defer rows.Close()
	items := []FindOrganizationsRow{}
	for rows.Next() {
		var item FindOrganizationsRow
		if err := rows.Scan(&item.OrganizationID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.SessionRemember, &item.SessionTimeout, &item.Email, &item.CollaboratorAuthPolicy, &item.AllowForceDeleteWorkspaces, &item.CostEstimationEnabled); err != nil {
			return nil, fmt.Errorf("scan FindOrganizationsBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindOrganizationsBatch rows: %w", err)
	}
	return items, err
}

const countOrganizationsSQL = `SELECT count(*)
FROM organizations
WHERE name LIKE ANY($1)
;`

// CountOrganizations implements Querier.CountOrganizations.
func (q *DBQuerier) CountOrganizations(ctx context.Context, names []string) (pgtype.Int8, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CountOrganizations")
	row := q.conn.QueryRow(ctx, countOrganizationsSQL, names)
	var item pgtype.Int8
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query CountOrganizations: %w", err)
	}
	return item, nil
}

// CountOrganizationsBatch implements Querier.CountOrganizationsBatch.
func (q *DBQuerier) CountOrganizationsBatch(batch genericBatch, names []string) {
	batch.Queue(countOrganizationsSQL, names)
}

// CountOrganizationsScan implements Querier.CountOrganizationsScan.
func (q *DBQuerier) CountOrganizationsScan(results pgx.BatchResults) (pgtype.Int8, error) {
	row := results.QueryRow()
	var item pgtype.Int8
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan CountOrganizationsBatch row: %w", err)
	}
	return item, nil
}

const updateOrganizationByNameSQL = `UPDATE organizations
SET
    name = $1,
    email = $2,
    collaborator_auth_policy = $3,
    cost_estimation_enabled = $4,
    session_remember = $5,
    session_timeout = $6,
    allow_force_delete_workspaces = $7,
    updated_at = $8
WHERE name = $9
RETURNING organization_id;`

type UpdateOrganizationByNameParams struct {
	NewName                    pgtype.Text        `json:"new_name"`
	Email                      pgtype.Text        `json:"email"`
	CollaboratorAuthPolicy     pgtype.Text        `json:"collaborator_auth_policy"`
	CostEstimationEnabled      pgtype.Bool        `json:"cost_estimation_enabled"`
	SessionRemember            pgtype.Int4        `json:"session_remember"`
	SessionTimeout             pgtype.Int4        `json:"session_timeout"`
	AllowForceDeleteWorkspaces pgtype.Bool        `json:"allow_force_delete_workspaces"`
	UpdatedAt                  pgtype.Timestamptz `json:"updated_at"`
	Name                       pgtype.Text        `json:"name"`
}

// UpdateOrganizationByName implements Querier.UpdateOrganizationByName.
func (q *DBQuerier) UpdateOrganizationByName(ctx context.Context, params UpdateOrganizationByNameParams) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateOrganizationByName")
	row := q.conn.QueryRow(ctx, updateOrganizationByNameSQL, params.NewName, params.Email, params.CollaboratorAuthPolicy, params.CostEstimationEnabled, params.SessionRemember, params.SessionTimeout, params.AllowForceDeleteWorkspaces, params.UpdatedAt, params.Name)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateOrganizationByName: %w", err)
	}
	return item, nil
}

// UpdateOrganizationByNameBatch implements Querier.UpdateOrganizationByNameBatch.
func (q *DBQuerier) UpdateOrganizationByNameBatch(batch genericBatch, params UpdateOrganizationByNameParams) {
	batch.Queue(updateOrganizationByNameSQL, params.NewName, params.Email, params.CollaboratorAuthPolicy, params.CostEstimationEnabled, params.SessionRemember, params.SessionTimeout, params.AllowForceDeleteWorkspaces, params.UpdatedAt, params.Name)
}

// UpdateOrganizationByNameScan implements Querier.UpdateOrganizationByNameScan.
func (q *DBQuerier) UpdateOrganizationByNameScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateOrganizationByNameBatch row: %w", err)
	}
	return item, nil
}

const deleteOrganizationByNameSQL = `DELETE
FROM organizations
WHERE name = $1
RETURNING organization_id;`

// DeleteOrganizationByName implements Querier.DeleteOrganizationByName.
func (q *DBQuerier) DeleteOrganizationByName(ctx context.Context, name pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteOrganizationByName")
	row := q.conn.QueryRow(ctx, deleteOrganizationByNameSQL, name)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteOrganizationByName: %w", err)
	}
	return item, nil
}

// DeleteOrganizationByNameBatch implements Querier.DeleteOrganizationByNameBatch.
func (q *DBQuerier) DeleteOrganizationByNameBatch(batch genericBatch, name pgtype.Text) {
	batch.Queue(deleteOrganizationByNameSQL, name)
}

// DeleteOrganizationByNameScan implements Querier.DeleteOrganizationByNameScan.
func (q *DBQuerier) DeleteOrganizationByNameScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DeleteOrganizationByNameBatch row: %w", err)
	}
	return item, nil
}
