ALTER TABLE stations DROP COLUMN lati;
ALTER TABLE stations DROP COLUMN long;
ALTER TABLE stations ADD COLUMN pref_loc text NULL;