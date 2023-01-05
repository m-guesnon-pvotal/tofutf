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
    organization_id
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
	ID             pgtype.Text
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
	Name           pgtype.Text
	Provider       pgtype.Text
	Status         pgtype.Text
	OrganizationID pgtype.Text
}

// InsertModule implements Querier.InsertModule.
func (q *DBQuerier) InsertModule(ctx context.Context, params InsertModuleParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertModule")
	cmdTag, err := q.conn.Exec(ctx, insertModuleSQL, params.ID, params.CreatedAt, params.UpdatedAt, params.Name, params.Provider, params.Status, params.OrganizationID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertModule: %w", err)
	}
	return cmdTag, err
}

// InsertModuleBatch implements Querier.InsertModuleBatch.
func (q *DBQuerier) InsertModuleBatch(batch genericBatch, params InsertModuleParams) {
	batch.Queue(insertModuleSQL, params.ID, params.CreatedAt, params.UpdatedAt, params.Name, params.Provider, params.Status, params.OrganizationID)
}

// InsertModuleScan implements Querier.InsertModuleScan.
func (q *DBQuerier) InsertModuleScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertModuleBatch: %w", err)
	}
	return cmdTag, err
}

const insertModuleRepoSQL = `INSERT INTO module_repos (
    webhook_id,
    vcs_provider_id,
    module_id
) VALUES (
    $1,
    $2,
    $3
);`

type InsertModuleRepoParams struct {
	WebhookID     pgtype.UUID
	VCSProviderID pgtype.Text
	ModuleID      pgtype.Text
}

// InsertModuleRepo implements Querier.InsertModuleRepo.
func (q *DBQuerier) InsertModuleRepo(ctx context.Context, params InsertModuleRepoParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertModuleRepo")
	cmdTag, err := q.conn.Exec(ctx, insertModuleRepoSQL, params.WebhookID, params.VCSProviderID, params.ModuleID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertModuleRepo: %w", err)
	}
	return cmdTag, err
}

// InsertModuleRepoBatch implements Querier.InsertModuleRepoBatch.
func (q *DBQuerier) InsertModuleRepoBatch(batch genericBatch, params InsertModuleRepoParams) {
	batch.Queue(insertModuleRepoSQL, params.WebhookID, params.VCSProviderID, params.ModuleID)
}

// InsertModuleRepoScan implements Querier.InsertModuleRepoScan.
func (q *DBQuerier) InsertModuleRepoScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertModuleRepoBatch: %w", err)
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
	ModuleVersionID pgtype.Text
	Version         pgtype.Text
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	ModuleID        pgtype.Text
	Status          pgtype.Text
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
    (o.*)::"organizations" AS organization,
    (r.*)::"module_repos" AS module_repo,
    (h.*)::"webhooks" AS webhook,
    (
        SELECT array_agg(v.*) AS versions
        FROM module_versions v
        WHERE v.module_id = m.module_id
    ) AS versions
FROM modules m
JOIN organizations o USING (organization_id)
LEFT JOIN (module_repos r JOIN webhooks h USING (webhook_id)) USING (module_id)
WHERE o.name = $1
;`

type ListModulesByOrganizationRow struct {
	ModuleID     pgtype.Text        `json:"module_id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	Name         pgtype.Text        `json:"name"`
	Provider     pgtype.Text        `json:"provider"`
	Status       pgtype.Text        `json:"status"`
	Organization *Organizations     `json:"organization"`
	ModuleRepo   *ModuleRepos       `json:"module_repo"`
	Webhook      *Webhooks          `json:"webhook"`
	Versions     []ModuleVersions   `json:"versions"`
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
	organizationRow := q.types.newOrganizations()
	moduleRepoRow := q.types.newModuleRepos()
	webhookRow := q.types.newWebhooks()
	versionsArray := q.types.newModuleVersionsArray()
	for rows.Next() {
		var item ListModulesByOrganizationRow
		if err := rows.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, organizationRow, moduleRepoRow, webhookRow, versionsArray); err != nil {
			return nil, fmt.Errorf("scan ListModulesByOrganization row: %w", err)
		}
		if err := organizationRow.AssignTo(&item.Organization); err != nil {
			return nil, fmt.Errorf("assign ListModulesByOrganization row: %w", err)
		}
		if err := moduleRepoRow.AssignTo(&item.ModuleRepo); err != nil {
			return nil, fmt.Errorf("assign ListModulesByOrganization row: %w", err)
		}
		if err := webhookRow.AssignTo(&item.Webhook); err != nil {
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
	organizationRow := q.types.newOrganizations()
	moduleRepoRow := q.types.newModuleRepos()
	webhookRow := q.types.newWebhooks()
	versionsArray := q.types.newModuleVersionsArray()
	for rows.Next() {
		var item ListModulesByOrganizationRow
		if err := rows.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, organizationRow, moduleRepoRow, webhookRow, versionsArray); err != nil {
			return nil, fmt.Errorf("scan ListModulesByOrganizationBatch row: %w", err)
		}
		if err := organizationRow.AssignTo(&item.Organization); err != nil {
			return nil, fmt.Errorf("assign ListModulesByOrganization row: %w", err)
		}
		if err := moduleRepoRow.AssignTo(&item.ModuleRepo); err != nil {
			return nil, fmt.Errorf("assign ListModulesByOrganization row: %w", err)
		}
		if err := webhookRow.AssignTo(&item.Webhook); err != nil {
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
    (o.*)::"organizations" AS organization,
    (r.*)::"module_repos" AS module_repo,
    (h.*)::"webhooks" AS webhook,
    (
        SELECT array_agg(v.*) AS versions
        FROM module_versions v
        WHERE v.module_id = m.module_id
    ) AS versions
FROM modules m
JOIN organizations o USING (organization_id)
LEFT JOIN (module_repos r JOIN webhooks h USING (webhook_id)) USING (module_id)
WHERE o.name = $1
AND   m.name = $2
AND   m.provider = $3
;`

type FindModuleByNameParams struct {
	OrganizatonName pgtype.Text
	Name            pgtype.Text
	Provider        pgtype.Text
}

type FindModuleByNameRow struct {
	ModuleID     pgtype.Text        `json:"module_id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	Name         pgtype.Text        `json:"name"`
	Provider     pgtype.Text        `json:"provider"`
	Status       pgtype.Text        `json:"status"`
	Organization *Organizations     `json:"organization"`
	ModuleRepo   *ModuleRepos       `json:"module_repo"`
	Webhook      *Webhooks          `json:"webhook"`
	Versions     []ModuleVersions   `json:"versions"`
}

// FindModuleByName implements Querier.FindModuleByName.
func (q *DBQuerier) FindModuleByName(ctx context.Context, params FindModuleByNameParams) (FindModuleByNameRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindModuleByName")
	row := q.conn.QueryRow(ctx, findModuleByNameSQL, params.OrganizatonName, params.Name, params.Provider)
	var item FindModuleByNameRow
	organizationRow := q.types.newOrganizations()
	moduleRepoRow := q.types.newModuleRepos()
	webhookRow := q.types.newWebhooks()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, organizationRow, moduleRepoRow, webhookRow, versionsArray); err != nil {
		return item, fmt.Errorf("query FindModuleByName: %w", err)
	}
	if err := organizationRow.AssignTo(&item.Organization); err != nil {
		return item, fmt.Errorf("assign FindModuleByName row: %w", err)
	}
	if err := moduleRepoRow.AssignTo(&item.ModuleRepo); err != nil {
		return item, fmt.Errorf("assign FindModuleByName row: %w", err)
	}
	if err := webhookRow.AssignTo(&item.Webhook); err != nil {
		return item, fmt.Errorf("assign FindModuleByName row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByName row: %w", err)
	}
	return item, nil
}

// FindModuleByNameBatch implements Querier.FindModuleByNameBatch.
func (q *DBQuerier) FindModuleByNameBatch(batch genericBatch, params FindModuleByNameParams) {
	batch.Queue(findModuleByNameSQL, params.OrganizatonName, params.Name, params.Provider)
}

// FindModuleByNameScan implements Querier.FindModuleByNameScan.
func (q *DBQuerier) FindModuleByNameScan(results pgx.BatchResults) (FindModuleByNameRow, error) {
	row := results.QueryRow()
	var item FindModuleByNameRow
	organizationRow := q.types.newOrganizations()
	moduleRepoRow := q.types.newModuleRepos()
	webhookRow := q.types.newWebhooks()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, organizationRow, moduleRepoRow, webhookRow, versionsArray); err != nil {
		return item, fmt.Errorf("scan FindModuleByNameBatch row: %w", err)
	}
	if err := organizationRow.AssignTo(&item.Organization); err != nil {
		return item, fmt.Errorf("assign FindModuleByName row: %w", err)
	}
	if err := moduleRepoRow.AssignTo(&item.ModuleRepo); err != nil {
		return item, fmt.Errorf("assign FindModuleByName row: %w", err)
	}
	if err := webhookRow.AssignTo(&item.Webhook); err != nil {
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
    (o.*)::"organizations" AS organization,
    (r.*)::"module_repos" AS module_repo,
    (h.*)::"webhooks" AS webhook,
    (
        SELECT array_agg(v.*) AS versions
        FROM module_versions v
        WHERE v.module_id = m.module_id
    ) AS versions
FROM modules m
JOIN organizations o USING (organization_id)
LEFT JOIN (module_repos r JOIN webhooks h USING (webhook_id)) USING (module_id)
WHERE m.module_id = $1
;`

type FindModuleByIDRow struct {
	ModuleID     pgtype.Text        `json:"module_id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	Name         pgtype.Text        `json:"name"`
	Provider     pgtype.Text        `json:"provider"`
	Status       pgtype.Text        `json:"status"`
	Organization *Organizations     `json:"organization"`
	ModuleRepo   *ModuleRepos       `json:"module_repo"`
	Webhook      *Webhooks          `json:"webhook"`
	Versions     []ModuleVersions   `json:"versions"`
}

// FindModuleByID implements Querier.FindModuleByID.
func (q *DBQuerier) FindModuleByID(ctx context.Context, id pgtype.Text) (FindModuleByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindModuleByID")
	row := q.conn.QueryRow(ctx, findModuleByIDSQL, id)
	var item FindModuleByIDRow
	organizationRow := q.types.newOrganizations()
	moduleRepoRow := q.types.newModuleRepos()
	webhookRow := q.types.newWebhooks()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, organizationRow, moduleRepoRow, webhookRow, versionsArray); err != nil {
		return item, fmt.Errorf("query FindModuleByID: %w", err)
	}
	if err := organizationRow.AssignTo(&item.Organization); err != nil {
		return item, fmt.Errorf("assign FindModuleByID row: %w", err)
	}
	if err := moduleRepoRow.AssignTo(&item.ModuleRepo); err != nil {
		return item, fmt.Errorf("assign FindModuleByID row: %w", err)
	}
	if err := webhookRow.AssignTo(&item.Webhook); err != nil {
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
	organizationRow := q.types.newOrganizations()
	moduleRepoRow := q.types.newModuleRepos()
	webhookRow := q.types.newWebhooks()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, organizationRow, moduleRepoRow, webhookRow, versionsArray); err != nil {
		return item, fmt.Errorf("scan FindModuleByIDBatch row: %w", err)
	}
	if err := organizationRow.AssignTo(&item.Organization); err != nil {
		return item, fmt.Errorf("assign FindModuleByID row: %w", err)
	}
	if err := moduleRepoRow.AssignTo(&item.ModuleRepo); err != nil {
		return item, fmt.Errorf("assign FindModuleByID row: %w", err)
	}
	if err := webhookRow.AssignTo(&item.Webhook); err != nil {
		return item, fmt.Errorf("assign FindModuleByID row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByID row: %w", err)
	}
	return item, nil
}

const findModuleByWebhookIDSQL = `SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    (o.*)::"organizations" AS organization,
    (r.*)::"module_repos" AS module_repo,
    (h.*)::"webhooks" AS webhook,
    (
        SELECT array_agg(v.*) AS versions
        FROM module_versions v
        WHERE v.module_id = m.module_id
    ) AS versions
FROM modules m
JOIN organizations o USING (organization_id)
JOIN (module_repos r JOIN webhooks h USING (webhook_id)) USING (module_id)
WHERE h.webhook_id = $1
;`

type FindModuleByWebhookIDRow struct {
	ModuleID     pgtype.Text        `json:"module_id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	Name         pgtype.Text        `json:"name"`
	Provider     pgtype.Text        `json:"provider"`
	Status       pgtype.Text        `json:"status"`
	Organization *Organizations     `json:"organization"`
	ModuleRepo   *ModuleRepos       `json:"module_repo"`
	Webhook      *Webhooks          `json:"webhook"`
	Versions     []ModuleVersions   `json:"versions"`
}

// FindModuleByWebhookID implements Querier.FindModuleByWebhookID.
func (q *DBQuerier) FindModuleByWebhookID(ctx context.Context, webhookID pgtype.UUID) (FindModuleByWebhookIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindModuleByWebhookID")
	row := q.conn.QueryRow(ctx, findModuleByWebhookIDSQL, webhookID)
	var item FindModuleByWebhookIDRow
	organizationRow := q.types.newOrganizations()
	moduleRepoRow := q.types.newModuleRepos()
	webhookRow := q.types.newWebhooks()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, organizationRow, moduleRepoRow, webhookRow, versionsArray); err != nil {
		return item, fmt.Errorf("query FindModuleByWebhookID: %w", err)
	}
	if err := organizationRow.AssignTo(&item.Organization); err != nil {
		return item, fmt.Errorf("assign FindModuleByWebhookID row: %w", err)
	}
	if err := moduleRepoRow.AssignTo(&item.ModuleRepo); err != nil {
		return item, fmt.Errorf("assign FindModuleByWebhookID row: %w", err)
	}
	if err := webhookRow.AssignTo(&item.Webhook); err != nil {
		return item, fmt.Errorf("assign FindModuleByWebhookID row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByWebhookID row: %w", err)
	}
	return item, nil
}

// FindModuleByWebhookIDBatch implements Querier.FindModuleByWebhookIDBatch.
func (q *DBQuerier) FindModuleByWebhookIDBatch(batch genericBatch, webhookID pgtype.UUID) {
	batch.Queue(findModuleByWebhookIDSQL, webhookID)
}

// FindModuleByWebhookIDScan implements Querier.FindModuleByWebhookIDScan.
func (q *DBQuerier) FindModuleByWebhookIDScan(results pgx.BatchResults) (FindModuleByWebhookIDRow, error) {
	row := results.QueryRow()
	var item FindModuleByWebhookIDRow
	organizationRow := q.types.newOrganizations()
	moduleRepoRow := q.types.newModuleRepos()
	webhookRow := q.types.newWebhooks()
	versionsArray := q.types.newModuleVersionsArray()
	if err := row.Scan(&item.ModuleID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Provider, &item.Status, organizationRow, moduleRepoRow, webhookRow, versionsArray); err != nil {
		return item, fmt.Errorf("scan FindModuleByWebhookIDBatch row: %w", err)
	}
	if err := organizationRow.AssignTo(&item.Organization); err != nil {
		return item, fmt.Errorf("assign FindModuleByWebhookID row: %w", err)
	}
	if err := moduleRepoRow.AssignTo(&item.ModuleRepo); err != nil {
		return item, fmt.Errorf("assign FindModuleByWebhookID row: %w", err)
	}
	if err := webhookRow.AssignTo(&item.Webhook); err != nil {
		return item, fmt.Errorf("assign FindModuleByWebhookID row: %w", err)
	}
	if err := versionsArray.AssignTo(&item.Versions); err != nil {
		return item, fmt.Errorf("assign FindModuleByWebhookID row: %w", err)
	}
	return item, nil
}

const updateModuleStatusSQL = `UPDATE modules
SET status = $1
WHERE module_id = $2
RETURNING module_id
;`

// UpdateModuleStatus implements Querier.UpdateModuleStatus.
func (q *DBQuerier) UpdateModuleStatus(ctx context.Context, status pgtype.Text, moduleID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateModuleStatus")
	row := q.conn.QueryRow(ctx, updateModuleStatusSQL, status, moduleID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateModuleStatus: %w", err)
	}
	return item, nil
}

// UpdateModuleStatusBatch implements Querier.UpdateModuleStatusBatch.
func (q *DBQuerier) UpdateModuleStatusBatch(batch genericBatch, status pgtype.Text, moduleID pgtype.Text) {
	batch.Queue(updateModuleStatusSQL, status, moduleID)
}

// UpdateModuleStatusScan implements Querier.UpdateModuleStatusScan.
func (q *DBQuerier) UpdateModuleStatusScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateModuleStatusBatch row: %w", err)
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

const updateModuleVersionStatusSQL = `UPDATE module_versions
SET
    status = $1,
    status_error = $2
WHERE module_version_id = $3
RETURNING *
;`

type UpdateModuleVersionStatusParams struct {
	Status          pgtype.Text
	StatusError     pgtype.Text
	ModuleVersionID pgtype.Text
}

type UpdateModuleVersionStatusRow struct {
	ModuleVersionID pgtype.Text        `json:"module_version_id"`
	Version         pgtype.Text        `json:"version"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	Status          pgtype.Text        `json:"status"`
	StatusError     pgtype.Text        `json:"status_error"`
	ModuleID        pgtype.Text        `json:"module_id"`
}

// UpdateModuleVersionStatus implements Querier.UpdateModuleVersionStatus.
func (q *DBQuerier) UpdateModuleVersionStatus(ctx context.Context, params UpdateModuleVersionStatusParams) (UpdateModuleVersionStatusRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateModuleVersionStatus")
	row := q.conn.QueryRow(ctx, updateModuleVersionStatusSQL, params.Status, params.StatusError, params.ModuleVersionID)
	var item UpdateModuleVersionStatusRow
	if err := row.Scan(&item.ModuleVersionID, &item.Version, &item.CreatedAt, &item.UpdatedAt, &item.Status, &item.StatusError, &item.ModuleID); err != nil {
		return item, fmt.Errorf("query UpdateModuleVersionStatus: %w", err)
	}
	return item, nil
}

// UpdateModuleVersionStatusBatch implements Querier.UpdateModuleVersionStatusBatch.
func (q *DBQuerier) UpdateModuleVersionStatusBatch(batch genericBatch, params UpdateModuleVersionStatusParams) {
	batch.Queue(updateModuleVersionStatusSQL, params.Status, params.StatusError, params.ModuleVersionID)
}

// UpdateModuleVersionStatusScan implements Querier.UpdateModuleVersionStatusScan.
func (q *DBQuerier) UpdateModuleVersionStatusScan(results pgx.BatchResults) (UpdateModuleVersionStatusRow, error) {
	row := results.QueryRow()
	var item UpdateModuleVersionStatusRow
	if err := row.Scan(&item.ModuleVersionID, &item.Version, &item.CreatedAt, &item.UpdatedAt, &item.Status, &item.StatusError, &item.ModuleID); err != nil {
		return item, fmt.Errorf("scan UpdateModuleVersionStatusBatch row: %w", err)
	}
	return item, nil
}

const deleteModuleByIDSQL = `DELETE
FROM modules
WHERE module_id = $1
RETURNING module_id;`

// DeleteModuleByID implements Querier.DeleteModuleByID.
func (q *DBQuerier) DeleteModuleByID(ctx context.Context, id pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteModuleByID")
	row := q.conn.QueryRow(ctx, deleteModuleByIDSQL, id)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteModuleByID: %w", err)
	}
	return item, nil
}

// DeleteModuleByIDBatch implements Querier.DeleteModuleByIDBatch.
func (q *DBQuerier) DeleteModuleByIDBatch(batch genericBatch, id pgtype.Text) {
	batch.Queue(deleteModuleByIDSQL, id)
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