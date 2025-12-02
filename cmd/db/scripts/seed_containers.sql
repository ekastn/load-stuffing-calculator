-- Seed standard containers
INSERT INTO containers (name, inner_length_mm, inner_width_mm, inner_height_mm, max_weight_kg, description)
SELECT '20ft Standard', 5898.00, 2352.00, 2393.00, 28200.00, 'Standard 20ft container'
WHERE NOT EXISTS (SELECT 1 FROM containers WHERE name = '20ft Standard');

INSERT INTO containers (name, inner_length_mm, inner_width_mm, inner_height_mm, max_weight_kg, description)
SELECT '40ft Standard', 12032.00, 2352.00, 2393.00, 26600.00, 'Standard 40ft container'
WHERE NOT EXISTS (SELECT 1 FROM containers WHERE name = '40ft Standard');

INSERT INTO containers (name, inner_length_mm, inner_width_mm, inner_height_mm, max_weight_kg, description)
SELECT '40ft High Cube', 12032.00, 2352.00, 2698.00, 26460.00, 'High Cube 40ft container'
WHERE NOT EXISTS (SELECT 1 FROM containers WHERE name = '40ft High Cube');
