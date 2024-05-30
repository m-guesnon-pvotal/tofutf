-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS run_triggers (
    run_trigger_id text PRIMARY KEY,
    created_at        TIMESTAMPTZ NOT NULL,
    sourceable_id   TEXT REFERENCES workspaces ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    workspace_id TEXT REFERENCES workspaces ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    UNIQUE(sourceable_id, workspace_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE run_triggers;

-- +goose StatementEnd
