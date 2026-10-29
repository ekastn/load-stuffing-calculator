-- +goose Up
-- +goose StatementBegin
ALTER TABLE load_plans
ADD COLUMN title VARCHAR(150),
ADD COLUMN notes TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE load_plans
DROP COLUMN IF EXISTS notes,
DROP COLUMN IF EXISTS title;
-- +goose StatementEnd
