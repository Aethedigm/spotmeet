ALTER TABLE links 
ADD COLUMN song_id INTEGER REFERENCES songs(id) ON DELETE CASCADE;