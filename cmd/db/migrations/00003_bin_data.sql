-- +goose Up
-- +goose StatementBegin
CREATE TABLE load_plans (
    plan_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Metadata
    plan_code VARCHAR(50) NOT NULL, 
    status VARCHAR(20) DEFAULT 'DRAFT', -- DRAFT, CALCULATED
    
    cont_label VARCHAR(100), -- e.g. "Truck Box A"
    length_mm NUMERIC(10,2) NOT NULL,
    width_mm NUMERIC(10,2) NOT NULL,
    height_mm NUMERIC(10,2) NOT NULL,
    max_weight_kg NUMERIC(10,2) NOT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE load_items (
    item_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id UUID REFERENCES load_plans(plan_id) ON DELETE CASCADE,
    
    item_label VARCHAR(150), -- e.g. "Kardus Indomie"
    length_mm NUMERIC(10,2) NOT NULL,
    width_mm NUMERIC(10,2) NOT NULL,
    height_mm NUMERIC(10,2) NOT NULL,
    weight_kg NUMERIC(10,2) NOT NULL,
    
    quantity INTEGER NOT NULL DEFAULT 1,
    
    -- Constraints untuk Algoritma
    allow_rotation BOOLEAN DEFAULT TRUE, 
    color_hex VARCHAR(7) DEFAULT '#3498db'
);

CREATE TABLE plan_results (
    result_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id UUID REFERENCES load_plans(plan_id) ON DELETE CASCADE,
    
    -- Metrics
    total_loaded_weight_kg NUMERIC(10,2),
    volume_utilization_pct NUMERIC(5,2), -- Persentase kepenuhan (0-100%)
    
    -- Validasi Algoritma
    is_feasible BOOLEAN DEFAULT TRUE, -- False jika ada barang sisa (gak muat)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE plan_placements (
    placement_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    result_id UUID REFERENCES plan_results(result_id) ON DELETE CASCADE,
    
    item_id UUID REFERENCES load_items(item_id),
    
    pos_x NUMERIC(10,2) NOT NULL,
    pos_y NUMERIC(10,2) NOT NULL,
    pos_z NUMERIC(10,2) NOT NULL,
    
    -- ORIENTASI BARANG
    -- 0: WxHxD (Default)
    -- 1: HxWxD (Diputar 90 derajat, dst...)
    rotation_code INTEGER DEFAULT 0,
    
    step_number INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE plan_placements;
DROP TABLE plan_results;
DROP TABLE load_items;
DROP TABLE load_plans;
-- +goose StatementEnd
