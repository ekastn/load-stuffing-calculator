-- +goose Up
-- +goose StatementBegin
CREATE TABLE containers (
    container_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    
    inner_length_mm NUMERIC(10,2) NOT NULL,
    inner_width_mm NUMERIC(10,2) NOT NULL,
    inner_height_mm NUMERIC(10,2) NOT NULL,
    
    max_weight_kg NUMERIC(10,2) NOT NULL,
    
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    product_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    
    length_mm NUMERIC(10,2) NOT NULL,
    width_mm NUMERIC(10,2) NOT NULL,
    height_mm NUMERIC(10,2) NOT NULL,
    weight_kg NUMERIC(10,2) NOT NULL,
    
    color_hex VARCHAR(7) DEFAULT '#cccccc',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS containers;
-- +goose StatementEnd
