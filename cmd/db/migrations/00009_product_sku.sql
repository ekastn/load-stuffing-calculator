-- +goose Up
-- +goose StatementBegin
ALTER TABLE products ADD COLUMN sku VARCHAR(100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN IF EXISTS sku;
-- +goose StatementEnd
