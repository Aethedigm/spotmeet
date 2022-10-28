DROP TABLE IF EXISTS artists cascade;

CREATE TABLE artists (
    id SERIAL PRIMARY KEY,
    spotify_id VARCHAR(255) NOT NULL,
    artist_name VARCHAR(255) NOT NULL
);