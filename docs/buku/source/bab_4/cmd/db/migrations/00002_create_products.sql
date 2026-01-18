-- +goose Up
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    label VARCHAR(255) NOT NULL,
    sku VARCHAR(100) UNIQUE NOT NULL,
    length_mm DOUBLE PRECISION NOT NULL,
    width_mm DOUBLE PRECISION NOT NULL,
    height_mm DOUBLE PRECISION NOT NULL,
    weight_kg DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS products;
