-- Create index "idx_stations_position" to table: "stations"
CREATE UNIQUE INDEX idx_stations_position ON stations (position_id) WHERE position_id is not null;
