ALTER TABLE images
DROP CONSTRAINT fk_images_mission;

ALTER TABLE images
DROP COLUMN mission_id;
