ALTER TABLE mission
ADD COLUMN points_hash TEXT;

CREATE UNIQUE INDEX mission_points_hash_uq
ON mission(points_hash);

ALTER TABLE mission
ALTER COLUMN points_hash SET NOT NULL;

