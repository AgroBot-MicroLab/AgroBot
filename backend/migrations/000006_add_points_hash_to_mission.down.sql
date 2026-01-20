DROP INDEX IF EXISTS mission_points_hash_uq;

ALTER TABLE mission
DROP COLUMN points_hash;
