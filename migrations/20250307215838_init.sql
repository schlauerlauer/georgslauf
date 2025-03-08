-- Create "schedule" table
CREATE TABLE schedule (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, start integer NOT NULL, end integer NULL, name text NOT NULL);
-- Create index "idx_schedule_start" to table: "schedule"
CREATE INDEX idx_schedule_start ON schedule (start);
-- Create "tribes" table
CREATE TABLE tribes (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), updated_at integer NOT NULL DEFAULT (unixepoch()), name text NOT NULL, short text NULL, dpsg text NULL, image_id text NULL, email_domain text NULL, stavo_email text NULL, CONSTRAINT image_id FOREIGN KEY (image_id) REFERENCES images (id) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "idx_tribes_name" to table: "tribes"
CREATE UNIQUE INDEX idx_tribes_name ON tribes (name);
-- Create "station_categories" table
CREATE TABLE station_categories (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, name text NOT NULL, max integer NOT NULL DEFAULT 0);
-- Create index "idx_station_categories_name" to table: "station_categories"
CREATE UNIQUE INDEX idx_station_categories_name ON station_categories (name);
-- Create "stations" table
CREATE TABLE stations (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), created_by integer NULL, updated_at integer NOT NULL DEFAULT (unixepoch()), updated_by integer NULL, name text NOT NULL, abbr text NULL, size integer NOT NULL DEFAULT 0, tribe_id integer NOT NULL, category_id integer NULL, lati real NULL, long real NULL, image_id text NULL, description text NULL, requirements text NULL, CONSTRAINT tribe_id FOREIGN KEY (tribe_id) REFERENCES tribes (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT image_id FOREIGN KEY (image_id) REFERENCES images (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT created_by FOREIGN KEY (created_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT category_id FOREIGN KEY (category_id) REFERENCES station_categories (id) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "idx_stations_abbr" to table: "stations"
CREATE UNIQUE INDEX idx_stations_abbr ON stations (abbr) WHERE abbr is not null;
-- Create "groups" table
CREATE TABLE groups (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), created_by integer NULL, updated_at integer NOT NULL DEFAULT (unixepoch()), updated_by integer NULL, name text NOT NULL, abbr text NULL, size integer NOT NULL DEFAULT 0, comment text NULL, grouping integer NOT NULL, tribe_id integer NOT NULL, image_id text NULL, CONSTRAINT tribe_id FOREIGN KEY (tribe_id) REFERENCES tribes (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT image_id FOREIGN KEY (image_id) REFERENCES images (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT created_by FOREIGN KEY (created_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "idx_groups_abbr" to table: "groups"
CREATE UNIQUE INDEX idx_groups_abbr ON groups (abbr) WHERE abbr is not null;
-- Create "images" table
CREATE TABLE images (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), created_by integer NULL, filepath text NOT NULL, tribe_id integer NULL, station_id integer NULL, CONSTRAINT created_by FOREIGN KEY (created_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT tribe_id FOREIGN KEY (tribe_id) REFERENCES tribes (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT station_id FOREIGN KEY (station_id) REFERENCES stations (id) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_image_filepath" to table: "images"
CREATE UNIQUE INDEX idx_image_filepath ON images (filepath);
-- Create "points_to_stations" table
CREATE TABLE points_to_stations (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), created_by integer NULL, updated_at integer NOT NULL DEFAULT (unixepoch()), updated_by integer NULL, group_id integer NOT NULL, station_id integer NOT NULL, CONSTRAINT group_id FOREIGN KEY (group_id) REFERENCES groups (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT station_id FOREIGN KEY (station_id) REFERENCES stations (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT created_by FOREIGN KEY (created_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "idx_pts" to table: "points_to_stations"
CREATE UNIQUE INDEX idx_pts ON points_to_stations (group_id, station_id);
-- Create "points_to_groups" table
CREATE TABLE points_to_groups (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), created_by integer NULL, updated_at integer NOT NULL DEFAULT (unixepoch()), updated_by integer NULL, station_id integer NOT NULL, group_id integer NOT NULL, CONSTRAINT station_id FOREIGN KEY (station_id) REFERENCES stations (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT group_id FOREIGN KEY (group_id) REFERENCES groups (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT created_by FOREIGN KEY (created_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "idx_ptg" to table: "points_to_groups"
CREATE UNIQUE INDEX idx_ptg ON points_to_groups (station_id, group_id);
-- Create "group_roles" table
CREATE TABLE group_roles (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), created_by integer NULL, updated_at integer NOT NULL DEFAULT (unixepoch()), updated_by integer NULL, user_id integer NOT NULL, group_id integer NOT NULL, group_role integer NOT NULL, CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT group_id FOREIGN KEY (group_id) REFERENCES groups (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT created_by FOREIGN KEY (created_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "idx_group_roles_user" to table: "group_roles"
CREATE INDEX idx_group_roles_user ON group_roles (user_id);
-- Create index "idx_group_roles_user_group" to table: "group_roles"
CREATE UNIQUE INDEX idx_group_roles_user_group ON group_roles (user_id, group_id);
-- Create "station_roles" table
CREATE TABLE station_roles (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), created_by integer NULL, updated_at integer NOT NULL DEFAULT (unixepoch()), updated_by integer NULL, user_id integer NOT NULL, station_id integer NOT NULL, station_role integer NOT NULL, CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT station_id FOREIGN KEY (station_id) REFERENCES stations (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT created_by FOREIGN KEY (created_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "idx_station_roles_user" to table: "station_roles"
CREATE INDEX idx_station_roles_user ON station_roles (user_id);
-- Create index "idx_station_roles_user_station" to table: "station_roles"
CREATE UNIQUE INDEX idx_station_roles_user_station ON station_roles (user_id, station_id);
-- Create "tribe_roles" table
CREATE TABLE tribe_roles (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), created_by integer NULL, updated_at integer NOT NULL DEFAULT (unixepoch()), updated_by integer NULL, user_id integer NOT NULL, tribe_id integer NOT NULL, tribe_role integer NOT NULL, accepted_at integer NULL, CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT tribe_id FOREIGN KEY (tribe_id) REFERENCES tribes (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT created_by FOREIGN KEY (created_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "idx_tribe_roles_user" to table: "tribe_roles"
CREATE INDEX idx_tribe_roles_user ON tribe_roles (user_id);
-- Create index "idx_tribe_roles_user_tribe" to table: "tribe_roles"
CREATE UNIQUE INDEX idx_tribe_roles_user_tribe ON tribe_roles (user_id, tribe_id);
-- Create "users" table
CREATE TABLE users (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, ext_id text NULL, username text NOT NULL, email text NOT NULL, last_login integer NOT NULL DEFAULT (unixepoch()), created_at integer NOT NULL DEFAULT (unixepoch()), role integer NOT NULL DEFAULT 0, firstname text NOT NULL, lastname text NOT NULL);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX idx_users_email ON users (email);
-- Create index "idx_users_ext_id" to table: "users"
CREATE UNIQUE INDEX idx_users_ext_id ON users (ext_id) WHERE ext_id is not null;
-- Create "tribe_icons" table
CREATE TABLE tribe_icons (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), created_by integer NULL, image blob NOT NULL, CONSTRAINT id FOREIGN KEY (id) REFERENCES tribes (id) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT created_by FOREIGN KEY (created_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create "user_icons" table
CREATE TABLE user_icons (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at integer NOT NULL DEFAULT (unixepoch()), image blob NOT NULL, CONSTRAINT id FOREIGN KEY (id) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "settings" table
CREATE TABLE settings (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, updated_at integer NOT NULL DEFAULT (unixepoch()), updated_by integer NULL, data blob NOT NULL, CONSTRAINT updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE SET NULL);
