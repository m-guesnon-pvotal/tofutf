-- name: InsertRunTrigger :exec
INSERT INTO run_triggers (
    run_trigger_id,
    created_at,
    sourceable_id,
    workspace_id
) VALUES (
    pggen.arg('run_trigger_id'),
    pggen.arg('created_at'),
    pggen.arg('sourceable_id'),
    pggen.arg('workspace_id')
);

-- name: FindRunTriggersBySourceableID :many
SELECT
    run_triggers.run_trigger_id,
    run_triggers.created_at,
    run_triggers.sourceable_id,
    run_triggers.workspace_id
FROM run_triggers
WHERE run_triggers.sourceable_id = pggen.arg('sourceable_id')
ORDER BY run_triggers.created_at DESC
;

-- name: FindRunTriggersByWorkspaceID :many
SELECT
    run_triggers.run_trigger_id,
    run_triggers.created_at,
    run_triggers.sourceable_id,
    run_triggers.workspace_id
FROM run_triggers
WHERE run_triggers.workspace_id = pggen.arg('workspace_id')
ORDER BY run_triggers.created_at DESC
;


-- name: FindRunTriggerByID :one
SELECT
    run_triggers.run_trigger_id,
    run_triggers.created_at,
    run_triggers.sourceable_id,
    run_triggers.workspace_id
FROM run_triggers
WHERE run_triggers.run_trigger_id = pggen.arg('run_trigger_id')
;

-- name: DeleteRunTriggerByID :one
DELETE
FROM run_triggers
WHERE run_trigger_id = pggen.arg('run_trigger_id')
RETURNING run_trigger_id
;
