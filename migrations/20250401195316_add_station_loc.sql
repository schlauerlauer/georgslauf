-- Create "station_positions" table
CREATE TABLE station_positions (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, name text NOT NULL);
-- Create index "idx_station_positions_name" to table: "station_positions"
CREATE UNIQUE INDEX idx_station_positions_name ON station_positions (name);

DROP INDEX idx_stations_abbr;
ALTER TABLE stations DROP COLUMN abbr;
ALTER TABLE stations DROP COLUMN pref_loc;
ALTER TABLE stations ADD COLUMN position_id integer NULL;
