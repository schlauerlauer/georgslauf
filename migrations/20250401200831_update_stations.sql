-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_stations" table
CREATE TABLE new_stations (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), created_by integer NULL, updated_at integer NOT NULL DEFAULT (unixepoch()), updated_by integer NULL, name text NOT NULL, position_id integer NULL, size integer NOT NULL DEFAULT 0, vegan integer NOT NULL DEFAULT 0, tribe_id integer NOT NULL, category_id integer NULL, image_id text NULL, description text NULL, requirements text NULL, CONSTRAINT tribe_id FOREIGN KEY (tribe_id) REFERENCES tribes (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT image_id FOREIGN KEY (image_id) REFERENCES images (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT created_by FOREIGN KEY (created_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT category_id FOREIGN KEY (category_id) REFERENCES station_categories (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT position_id FOREIGN KEY (position_id) REFERENCES station_positions (id) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Copy rows from old table "stations" to new temporary table "new_stations"
INSERT INTO new_stations (id, created_at, created_by, updated_at, updated_by, name, position_id, size, vegan, tribe_id, category_id, image_id, description, requirements) SELECT id, IFNULL(created_at, (unixepoch())) AS created_at, created_by, IFNULL(updated_at, (unixepoch())) AS updated_at, updated_by, name, position_id, size, vegan, tribe_id, category_id, image_id, description, requirements FROM stations;
-- Drop "stations" table after copying rows
DROP TABLE stations;
-- Rename temporary table "new_stations" to "stations"
ALTER TABLE new_stations RENAME TO stations;
-- Create index "idx_stations_tribe" to table: "stations"
CREATE INDEX idx_stations_tribe ON stations (tribe_id);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
