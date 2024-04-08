// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertModuleSQL = `INSERT INTO modules (
    module_id,
    created_at,
    updated_at,
    name,
    provider,
    status,
    organization_name
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
);`

type InsertModuleParams struct {
	ID               pgtype.Text        `json:"id"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
	Name             pgtype.Text        `json:"name"`
	Provider         pgtype.Text        `json:"provider"`
	Status           pgtype.Text        `json:"status"`
	OrganizationName pgtype.Text        `json:"organization_name"`
}

// InsertModule implements Querier.InsertModule.
func (q *DBQuerier) InsertModule(ctx context.Context, params InsertModuleParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertModule")
	cmdTag, err := q.conn.Exec(ctx, insertModuleSQL, params.ID, params.CreatedAt, params.UpdatedAt, params.Name, params.Provider, params.Status, params.OrganizationName)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertModule: %w", err)
	}
	return cmdTag, err
}

// InsertModuleBatch implements Querier.InsertModuleBatch.
func (q *DBQuerier) InsertModuleBatch(batch genericBatch, params InsertModuleParams) {
	batch.Queue(insertModuleSQL, params.ID, params.CreatedAt, params.UpdatedAt, params.Name, params.Provider, params.Status, params.OrganizationName)
}

// InsertModuleScan implements Querier.InsertModuleScan.
func (q *DBQuerier) InsertModuleScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertModuleBatch: %w", err)
	}
	return cmdTag, err
}

const insertModuleVersionSQL = `INSERT INTO module_versions (
    module_version_id,
    version,
    created_at,
    updated_at,
    module_id,
    status
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;`

type InsertModuleVersionParams struct {
	ModuleVersionID pgtype.Text        `json:"module_version_id"`
	Version         pgtype.Text        `json:"version"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	ModuleID        pgtype.Text        `json:"module_id"`
	Status          pgtype.Text        `json:"status"`
}

type InsertModuleVersionRow struct {
	ModuleVersionID pgtype.Text        `json:"module_version_id"`
	Version         pgtype.Text        `json:"version"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	Status          pgtype.Text        `json:"status"`
	StatusError     pgtype.Text        `json:"status_error"`
	ModuleID        pgtype.Text        `json:"module_id"`
}

// InsertModuleVersion implements Querier.InsertModuleVersion.
func (q *DBQuerier) InsertModuleVersion(ctx context.Context, params InsertModuleVersionParams) (InsertModuleVersionRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertModuleVersion")
	row := q.conn.QueryRow(ctx, insertModuleVersionSQL, params.ModuleVersionID, params.Version, params.CreatedAt, params.UpdatedAt, params.ModuleID, params.Status)
	var item InsertModuleVersionRow
	if err := row.Scan(&item.ModuleVersionID, &item.Version, &item.CreatedAt, &item.UpdatedAt, &item.Status, &item.StatusError, &item.ModuleID); err != nil {
		return item, fmt.Errorf("query InsertModuleVersion: %w", err)
	}
	return item, nil
}

// InsertModuleVersionBatch implements Querier.InsertModuleVersionBatch.
func (q *DBQuerier) InsertModuleVersionBatch(batch genericBatch, params InsertModuleVersionParams) {
	batch.Queue(insertModuleVersionSQL, params.ModuleVersionID, params.Version, params.CreatedAt, params.UpdatedAt, params.ModuleID, params.Status)
}

// InsertModuleVersionScan implements Querier.InsertModuleVersionScan.
func (q *DBQuerier) InsertModuleVersionScan(results pgx.BatchResults) (InsertModuleVersionRow, error) {
	row := results.QueryRow()
	var item InsertModuleVersionRow
	if err := row.Scan(&item.ModuleVersionID, &item.Version, &item.CreatedAt, &item.UpdatedAt, &item.Status, &item.StatusError, &item.ModuleID); err != nil {
		return item, fmt.Errorf("scan InsertModuleVersionBatch row: %w", err)
	}
	return item, nil
}

const listModulesByOrganizationSQL = `SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    m.organization_name,
    (r.*)::"repo_connections" AS module_connection,
    (
        SELECT array_agg(v.*) AS versions
        FROM module_versions v
        WHERE v.module_id = m.module_id
    ) AS versions
FROM modules m
LEFT JOIN repo_connections r USING (module_id)
WHERE m.organization_name = $1
;`

type ListModulesByOrganizationRow struct {
	ModuleID         pgtype.Text        `json:"module_id"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
	Name             pgtype.Text        `json:"name"`
	Provider         pgtype.Text        `json:"provider"`
	Status           pgtype.Text        `json:"status"`
	OrganizationName pgtype.Text        `json:"organization_name"`
	ModuleConnection *RepoConnections   `json:"module_connection"`
	Versions         []ModuleVersions   `json:"versions"`
}

// ListModulesByOrganization implements Querier.ListModulesByOrganization.
func (q *DBQuerier) ListModulesByOrganization(ctx context.Context, organizationName pgtype.Text) ([]ListModulesByOrganizationRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "ListModulesByOrganization")
	rows, err := q.conn.Query(ctx, listModulesByOrganizationSQL, organizationName)
	if err != nil {
		return nil, fmt.Errorf("query ListModulesByOrganization: %w", err)
	}
	defer rows.Close()
	items := []ListModulesByOrganizationRow{}
	moduleConnectionRow := q.types.newRepoConnections()
	versionsArray := q.types.newModuleVersionsArray()
	for rows.Next() {
		var item ListModulesByOrganizationRow
		if err := rows.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, &item.OrganizationName, moduleConnectionRow, versionsArray); err != nil {
			return nil, fmt.Errorf("scan ListModulesByOrganization row: %w", err)
		}
		if err := moduleConnectionRow.AssignTo(&item.ModuleConnection); err != nil {
			return nil, fmt.Errorf("assign ListModulesByOrganization row: %w", err)
		}
		if err := versionsArray.AssignTo(&item.Versions); err != nil {
			return nil, fmt.Errorf("assign ListModulesByOrganization row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close ListModulesByOrganization rows: %w", err)
	}
	return items, err
}

// ListModulesByOrganizationBatch implements Querier.ListModulesByOrganizationBatch.
func (q *DBQuerier) ListModulesByOrganizationBatch(batch genericBatch, organizationName pgtype.Text) {
	batch.Queue(listModulesByOrganizationSQL, organizationName)
}

// ListModulesByOrganizationScan implements Querier.ListModulesByOrganizationScan.
func (q *DBQuerier) ListModulesByOrganizationScan(results pgx.BatchResults) ([]ListModulesByOrganizationRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query ListModulesByOrganizationBatch: %w", err)
	}
	defer rows.Close()
	items := []ListModulesByOrganizationRow{}
	moduleConnectionRow := q.types.newRepoConnections()
	versionsArray := q.types.newModuleVersionsArray()
	for rows.Next() {
		var item ListModulesByOrganizationRow
		if err := rows.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, &item.OrganizationName, moduleConnectionRow, versionsArray); err != nil {
			return nil, fmt.Errorf("scan ListModulesByOrganizationBatch row: %w", err)
		}
		if err := moduleConnectionRow.AssignTo(&item.ModuleConnection); err != nil {
			return nil, fmt.Errorf("assign ListModulesByOrganization row: %w", err)
		}
		if err := versionsArray.AssignTo(&item.Versions); err != nil {
			return nil, fmt.Errorf("assign ListModulesByOrganization row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close ListModulesByOrganizationBatch rows: %w", err)
	}
	return items, err
}

const findModuleByNameSQL = `SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    m.organization_name,
    (r.*)::"repo_connections" AS module_connection,
    (
        SELECT array_agg(v.*) AS versions
        FROM module_versions v
        WHERE v.module_id = m.module_id
    ) AS versions
FROM modules m
LEFT JOIN repo_connections r USING (module_id)
WHERE m.organization_name = $1
AND   m.name = $2
AND   m.provider = $3
;`

type FindModuleByNameParams struct {
	OrganizationName pgtype.Text `json:"organization_name"`
	Name             pgtype.Text `json:"name"`
	Provider         pgtype.Text `json:"provider"`
}

type FindModuleByNameRow struct {
	ModuleID         pgtype.Text        `json:"module_id"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
	Name             pgtype.Text        `json:"name"`
	Provider         pgtype.Text        `json:"provider"`
	Status           pgtype.Text        `json:"status"`
	OrganizationName pgtype.Text        `json:"organization_name"`
	ModuleConnection *RepoConnections   `json:"module_connection"`
	Versions         []ModuleVersions   `json:"versions"`
}

// FindModuleByName implements Querier.FindModuleByName.
func (q *DBQuerier) FindModuleByName(ctx context.Context, params FindModuleByNameParams) (FindModuleByNameRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindModuleByName")
	row := q.conn.QueryRow(ctx, findModuleByNameSQL, params.OrganizationName, params.Name, params.Provider)
	var item FindModuleByNameRow
	moduleConnectionRow := q.types.newRepoConnections()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, &item.OrganizationName, moduleConnectionRow, versionsArray); err != nil {
		return item, fmt.Errorf("query FindModuleByName: %w", err)
	}
	if err := moduleConnectionRow.AssignTo(&item.ModuleConnection); err != nil {
		return item, fmt.Errorf("assign FindModuleByName row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByName row: %w", err)
	}
	return item, nil
}

// FindModuleByNameBatch implements Querier.FindModuleByNameBatch.
func (q *DBQuerier) FindModuleByNameBatch(batch genericBatch, params FindModuleByNameParams) {
	batch.Queue(findModuleByNameSQL, params.OrganizationName, params.Name, params.Provider)
}

// FindModuleByNameScan implements Querier.FindModuleByNameScan.
func (q *DBQuerier) FindModuleByNameScan(results pgx.BatchResults) (FindModuleByNameRow, error) {
	row := results.QueryRow()
	var item FindModuleByNameRow
	moduleConnectionRow := q.types.newRepoConnections()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, &item.OrganizationName, moduleConnectionRow, versionsArray); err != nil {
		return item, fmt.Errorf("scan FindModuleByNameBatch row: %w", err)
	}
	if err := moduleConnectionRow.AssignTo(&item.ModuleConnection); err != nil {
		return item, fmt.Errorf("assign FindModuleByName row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByName row: %w", err)
	}
	return item, nil
}

const findModuleByIDSQL = `SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    m.organization_name,
    (r.*)::"repo_connections" AS module_connection,
    (
        SELECT array_agg(v.*) AS versions
        FROM module_versions v
        WHERE v.module_id = m.module_id
    ) AS versions
FROM modules m
LEFT JOIN repo_connections r USING (module_id)
WHERE m.module_id = $1
;`

type FindModuleByIDRow struct {
	ModuleID         pgtype.Text        `json:"module_id"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
	Name             pgtype.Text        `json:"name"`
	Provider         pgtype.Text        `json:"provider"`
	Status           pgtype.Text        `json:"status"`
	OrganizationName pgtype.Text        `json:"organization_name"`
	ModuleConnection *RepoConnections   `json:"module_connection"`
	Versions         []ModuleVersions   `json:"versions"`
}

// FindModuleByID implements Querier.FindModuleByID.
func (q *DBQuerier) FindModuleByID(ctx context.Context, id pgtype.Text) (FindModuleByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindModuleByID")
	row := q.conn.QueryRow(ctx, findModuleByIDSQL, id)
	var item FindModuleByIDRow
	moduleConnectionRow := q.types.newRepoConnections()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, &item.OrganizationName, moduleConnectionRow, versionsArray); err != nil {
		return item, fmt.Errorf("query FindModuleByID: %w", err)
	}
	if err := moduleConnectionRow.AssignTo(&item.ModuleConnection); err != nil {
		return item, fmt.Errorf("assign FindModuleByID row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByID row: %w", err)
	}
	return item, nil
}

// FindModuleByIDBatch implements Querier.FindModuleByIDBatch.
func (q *DBQuerier) FindModuleByIDBatch(batch genericBatch, id pgtype.Text) {
	batch.Queue(findModuleByIDSQL, id)
}

// FindModuleByIDScan implements Querier.FindModuleByIDScan.
func (q *DBQuerier) FindModuleByIDScan(results pgx.BatchResults) (FindModuleByIDRow, error) {
	row := results.QueryRow()
	var item FindModuleByIDRow
	moduleConnectionRow := q.types.newRepoConnections()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, &item.OrganizationName, moduleConnectionRow, versionsArray); err != nil {
		return item, fmt.Errorf("scan FindModuleByIDBatch row: %w", err)
	}
	if err := moduleConnectionRow.AssignTo(&item.ModuleConnection); err != nil {
		return item, fmt.Errorf("assign FindModuleByID row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByID row: %w", err)
	}
	return item, nil
}

const findModuleByConnectionSQL = `SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    m.organization_name,
    (r.*)::"repo_connections" AS module_connection,
    (
        SELECT array_agg(v.*) AS versions
        FROM module_versions v
        WHERE v.module_id = m.module_id
    ) AS versions
FROM modules m
JOIN repo_connections r USING (module_id)
WHERE r.vcs_provider_id = $1
AND   r.repo_path = $2
;`

type FindModuleByConnectionRow struct {
	ModuleID         pgtype.Text        `json:"module_id"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
	Name             pgtype.Text        `json:"name"`
	Provider         pgtype.Text        `json:"provider"`
	Status           pgtype.Text        `json:"status"`
	OrganizationName pgtype.Text        `json:"organization_name"`
	ModuleConnection *RepoConnections   `json:"module_connection"`
	Versions         []ModuleVersions   `json:"versions"`
}

// FindModuleByConnection implements Querier.FindModuleByConnection.
func (q *DBQuerier) FindModuleByConnection(ctx context.Context, vcsProviderID pgtype.Text, repoPath pgtype.Text) (FindModuleByConnectionRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindModuleByConnection")
	row := q.conn.QueryRow(ctx, findModuleByConnectionSQL, vcsProviderID, repoPath)
	var item FindModuleByConnectionRow
	moduleConnectionRow := q.types.newRepoConnections()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, &item.OrganizationName, moduleConnectionRow, versionsArray); err != nil {
		return item, fmt.Errorf("query FindModuleByConnection: %w", err)
	}
	if err := moduleConnectionRow.AssignTo(&item.ModuleConnection); err != nil {
		return item, fmt.Errorf("assign FindModuleByConnection row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByConnection row: %w", err)
	}
	return item, nil
}

// FindModuleByConnectionBatch implements Querier.FindModuleByConnectionBatch.
func (q *DBQuerier) FindModuleByConnectionBatch(batch genericBatch, vcsProviderID pgtype.Text, repoPath pgtype.Text) {
	batch.Queue(findModuleByConnectionSQL, vcsProviderID, repoPath)
}

// FindModuleByConnectionScan implements Querier.FindModuleByConnectionScan.
func (q *DBQuerier) FindModuleByConnectionScan(results pgx.BatchResults) (FindModuleByConnectionRow, error) {
	row := results.QueryRow()
	var item FindModuleByConnectionRow
	moduleConnectionRow := q.types.newRepoConnections()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, &item.OrganizationName, moduleConnectionRow, versionsArray); err != nil {
		return item, fmt.Errorf("scan FindModuleByConnectionBatch row: %w", err)
	}
	if err := moduleConnectionRow.AssignTo(&item.ModuleConnection); err != nil {
		return item, fmt.Errorf("assign FindModuleByConnection row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByConnection row: %w", err)
	}
	return item, nil
}

const findModuleByModuleVersionIDSQL = `SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    m.organization_name,
    (r.*)::"repo_connections" AS module_connection,
    (
        SELECT array_agg(v.*) AS versions
        FROM module_versions v
        WHERE v.module_id = m.module_id
    ) AS versions
FROM modules m
JOIN module_versions mv USING (module_id)
LEFT JOIN repo_connections r USING (module_id)
WHERE mv.module_version_id = $1
;`

type FindModuleByModuleVersionIDRow struct {
	ModuleID         pgtype.Text        `json:"module_id"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
	Name             pgtype.Text        `json:"name"`
	Provider         pgtype.Text        `json:"provider"`
	Status           pgtype.Text        `json:"status"`
	OrganizationName pgtype.Text        `json:"organization_name"`
	ModuleConnection *RepoConnections   `json:"module_connection"`
	Versions         []ModuleVersions   `json:"versions"`
}

// FindModuleByModuleVersionID implements Querier.FindModuleByModuleVersionID.
func (q *DBQuerier) FindModuleByModuleVersionID(ctx context.Context, moduleVersionID pgtype.Text) (FindModuleByModuleVersionIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindModuleByModuleVersionID")
	row := q.conn.QueryRow(ctx, findModuleByModuleVersionIDSQL, moduleVersionID)
	var item FindModuleByModuleVersionIDRow
	moduleConnectionRow := q.types.newRepoConnections()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, &item.OrganizationName, moduleConnectionRow, versionsArray); err != nil {
		return item, fmt.Errorf("query FindModuleByModuleVersionID: %w", err)
	}
	if err := moduleConnectionRow.AssignTo(&item.ModuleConnection); err != nil {
		return item, fmt.Errorf("assign FindModuleByModuleVersionID row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByModuleVersionID row: %w", err)
	}
	return item, nil
}

// FindModuleByModuleVersionIDBatch implements Querier.FindModuleByModuleVersionIDBatch.
func (q *DBQuerier) FindModuleByModuleVersionIDBatch(batch genericBatch, moduleVersionID pgtype.Text) {
	batch.Queue(findModuleByModuleVersionIDSQL, moduleVersionID)
}

// FindModuleByModuleVersionIDScan implements Querier.FindModuleByModuleVersionIDScan.
func (q *DBQuerier) FindModuleByModuleVersionIDScan(results pgx.BatchResults) (FindModuleByModuleVersionIDRow, error) {
	row := results.QueryRow()
	var item FindModuleByModuleVersionIDRow
	moduleConnectionRow := q.types.newRepoConnections()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, &item.OrganizationName, moduleConnectionRow, versionsArray); err != nil {
		return item, fmt.Errorf("scan FindModuleByModuleVersionIDBatch row: %w", err)
	}
	if err := moduleConnectionRow.AssignTo(&item.ModuleConnection); err != nil {
		return item, fmt.Errorf("assign FindModuleByModuleVersionID row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByModuleVersionID row: %w", err)
	}
	return item, nil
}

const updateModuleStatusByIDSQL = `UPDATE modules
SET status = $1
WHERE module_id = $2
RETURNING module_id
;`

// UpdateModuleStatusByID implements Querier.UpdateModuleStatusByID.
func (q *DBQuerier) UpdateModuleStatusByID(ctx context.Context, status pgtype.Text, moduleID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateModuleStatusByID")
	row := q.conn.QueryRow(ctx, updateModuleStatusByIDSQL, status, moduleID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateModuleStatusByID: %w", err)
	}
	return item, nil
}

// UpdateModuleStatusByIDBatch implements Querier.UpdateModuleStatusByIDBatch.
func (q *DBQuerier) UpdateModuleStatusByIDBatch(batch genericBatch, status pgtype.Text, moduleID pgtype.Text) {
	batch.Queue(updateModuleStatusByIDSQL, status, moduleID)
}

// UpdateModuleStatusByIDScan implements Querier.UpdateModuleStatusByIDScan.
func (q *DBQuerier) UpdateModuleStatusByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateModuleStatusByIDBatch row: %w", err)
	}
	return item, nil
}

const insertModuleTarballSQL = `INSERT INTO module_tarballs (
    tarball,
    module_version_id
) VALUES (
    $1,
    $2
)
RETURNING module_version_id;`

// InsertModuleTarball implements Querier.InsertModuleTarball.
func (q *DBQuerier) InsertModuleTarball(ctx context.Context, tarball []byte, moduleVersionID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertModuleTarball")
	row := q.conn.QueryRow(ctx, insertModuleTarballSQL, tarball, moduleVersionID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query InsertModuleTarball: %w", err)
	}
	return item, nil
}

// InsertModuleTarballBatch implements Querier.InsertModuleTarballBatch.
func (q *DBQuerier) InsertModuleTarballBatch(batch genericBatch, tarball []byte, moduleVersionID pgtype.Text) {
	batch.Queue(insertModuleTarballSQL, tarball, moduleVersionID)
}

// InsertModuleTarballScan implements Querier.InsertModuleTarballScan.
func (q *DBQuerier) InsertModuleTarballScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan InsertModuleTarballBatch row: %w", err)
	}
	return item, nil
}

const findModuleTarballSQL = `SELECT tarball
FROM module_tarballs
WHERE module_version_id = $1
;`

// FindModuleTarball implements Querier.FindModuleTarball.
func (q *DBQuerier) FindModuleTarball(ctx context.Context, moduleVersionID pgtype.Text) ([]byte, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindModuleTarball")
	row := q.conn.QueryRow(ctx, findModuleTarballSQL, moduleVersionID)
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query FindModuleTarball: %w", err)
	}
	return item, nil
}

// FindModuleTarballBatch implements Querier.FindModuleTarballBatch.
func (q *DBQuerier) FindModuleTarballBatch(batch genericBatch, moduleVersionID pgtype.Text) {
	batch.Queue(findModuleTarballSQL, moduleVersionID)
}

// FindModuleTarballScan implements Querier.FindModuleTarballScan.
func (q *DBQuerier) FindModuleTarballScan(results pgx.BatchResults) ([]byte, error) {
	row := results.QueryRow()
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan FindModuleTarballBatch row: %w", err)
	}
	return item, nil
}

const updateModuleVersionStatusByIDSQL = `UPDATE module_versions
SET
    status = $1,
    status_error = $2
WHERE module_version_id = $3
RETURNING *
;`

type UpdateModuleVersionStatusByIDParams struct {
	Status          pgtype.Text `json:"status"`
	StatusError     pgtype.Text `json:"status_error"`
	ModuleVersionID pgtype.Text `json:"module_version_id"`
}

type UpdateModuleVersionStatusByIDRow struct {
	ModuleVersionID pgtype.Text        `json:"module_version_id"`
	Version         pgtype.Text        `json:"version"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	Status          pgtype.Text        `json:"status"`
	StatusError     pgtype.Text        `json:"status_error"`
	ModuleID        pgtype.Text        `json:"module_id"`
}

// UpdateModuleVersionStatusByID implements Querier.UpdateModuleVersionStatusByID.
func (q *DBQuerier) UpdateModuleVersionStatusByID(ctx context.Context, params UpdateModuleVersionStatusByIDParams) (UpdateModuleVersionStatusByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateModuleVersionStatusByID")
	row := q.conn.QueryRow(ctx, updateModuleVersionStatusByIDSQL, params.Status, params.StatusError, params.ModuleVersionID)
	var item UpdateModuleVersionStatusByIDRow
	if err := row.Scan(&item.ModuleVersionID, &item.Version, &item.CreatedAt, &item.UpdatedAt, &item.Status, &item.StatusError, &item.ModuleID); err != nil {
		return item, fmt.Errorf("query UpdateModuleVersionStatusByID: %w", err)
	}
	return item, nil
}

// UpdateModuleVersionStatusByIDBatch implements Querier.UpdateModuleVersionStatusByIDBatch.
func (q *DBQuerier) UpdateModuleVersionStatusByIDBatch(batch genericBatch, params UpdateModuleVersionStatusByIDParams) {
	batch.Queue(updateModuleVersionStatusByIDSQL, params.Status, params.StatusError, params.ModuleVersionID)
}

// UpdateModuleVersionStatusByIDScan implements Querier.UpdateModuleVersionStatusByIDScan.
func (q *DBQuerier) UpdateModuleVersionStatusByIDScan(results pgx.BatchResults) (UpdateModuleVersionStatusByIDRow, error) {
	row := results.QueryRow()
	var item UpdateModuleVersionStatusByIDRow
	if err := row.Scan(&item.ModuleVersionID, &item.Version, &item.CreatedAt, &item.UpdatedAt, &item.Status, &item.StatusError, &item.ModuleID); err != nil {
		return item, fmt.Errorf("scan UpdateModuleVersionStatusByIDBatch row: %w", err)
	}
	return item, nil
}

const deleteModuleByIDSQL = `DELETE
FROM modules
WHERE module_id = $1
RETURNING module_id
;`

// DeleteModuleByID implements Querier.DeleteModuleByID.
func (q *DBQuerier) DeleteModuleByID(ctx context.Context, moduleID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteModuleByID")
	row := q.conn.QueryRow(ctx, deleteModuleByIDSQL, moduleID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteModuleByID: %w", err)
	}
	return item, nil
}

// DeleteModuleByIDBatch implements Querier.DeleteModuleByIDBatch.
func (q *DBQuerier) DeleteModuleByIDBatch(batch genericBatch, moduleID pgtype.Text) {
	batch.Queue(deleteModuleByIDSQL, moduleID)
}

// DeleteModuleByIDScan implements Querier.DeleteModuleByIDScan.
func (q *DBQuerier) DeleteModuleByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DeleteModuleByIDBatch row: %w", err)
	}
	return item, nil
}

const deleteModuleVersionByIDSQL = `DELETE
FROM module_versions
WHERE module_version_id = $1
RETURNING module_version_id
;`

// DeleteModuleVersionByID implements Querier.DeleteModuleVersionByID.
func (q *DBQuerier) DeleteModuleVersionByID(ctx context.Context, moduleVersionID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteModuleVersionByID")
	row := q.conn.QueryRow(ctx, deleteModuleVersionByIDSQL, moduleVersionID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteModuleVersionByID: %w", err)
	}
	return item, nil
}

// DeleteModuleVersionByIDBatch implements Querier.DeleteModuleVersionByIDBatch.
func (q *DBQuerier) DeleteModuleVersionByIDBatch(batch genericBatch, moduleVersionID pgtype.Text) {
	batch.Queue(deleteModuleVersionByIDSQL, moduleVersionID)
}

// DeleteModuleVersionByIDScan implements Querier.DeleteModuleVersionByIDScan.
func (q *DBQuerier) DeleteModuleVersionByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DeleteModuleVersionByIDBatch row: %w", err)
	}
	return item, nil
}
