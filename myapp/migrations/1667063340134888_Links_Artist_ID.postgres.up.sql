ALTER TABLE links 
ADD COLUMN artist_id INTEGER REFERENCES artists(id) ON DELETE CASCADE;