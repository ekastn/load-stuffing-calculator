-- Seed sample products
INSERT INTO products (name, length_mm, width_mm, height_mm, weight_kg, color_hex)
SELECT 'Small Box', 300.00, 200.00, 200.00, 5.00, '#FF5733'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Small Box');

INSERT INTO products (name, length_mm, width_mm, height_mm, weight_kg, color_hex)
SELECT 'Medium Box', 500.00, 400.00, 400.00, 15.00, '#33FF57'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Medium Box');

INSERT INTO products (name, length_mm, width_mm, height_mm, weight_kg, color_hex)
SELECT 'Large Crate', 1000.00, 800.00, 800.00, 100.00, '#3357FF'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Large Crate');

INSERT INTO products (name, length_mm, width_mm, height_mm, weight_kg, color_hex)
SELECT 'Euro Pallet', 1200.00, 800.00, 144.00, 25.00, '#C70039'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Euro Pallet');
