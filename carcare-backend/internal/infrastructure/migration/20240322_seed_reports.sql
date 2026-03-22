-- Seed: Insert example report
INSERT INTO reports (type, created_at, data) VALUES (
    'maintenance',
    NOW(),
    '{"car_id": 1, "description": "Annual maintenance completed", "cost": 120.50}'
);
