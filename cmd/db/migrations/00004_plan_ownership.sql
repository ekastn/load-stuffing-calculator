-- +goose Up
-- +goose StatementBegin
ALTER TABLE load_plans
    ADD COLUMN created_by_type TEXT NOT NULL DEFAULT 'user',
    ADD COLUMN created_by_id UUID;

-- Backfill legacy plans as system-owned.
-- Note: this UUID is a fixed sentinel (no corresponding user row required).
UPDATE load_plans
SET created_by_type = 'system',
    created_by_id = '11111111-1111-1111-1111-111111111111'
WHERE created_by_id IS NULL;

ALTER TABLE load_plans
    ALTER COLUMN created_by_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_load_plans_created_by ON load_plans (created_by_type, created_by_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_load_plans_created_by;

ALTER TABLE load_plans
    DROP COLUMN IF EXISTS created_by_id,
    DROP COLUMN IF EXISTS created_by_type;
-- +goose StatementEnd
