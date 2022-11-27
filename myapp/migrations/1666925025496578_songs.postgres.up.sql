DROP TABLE IF EXISTS songs cascade;

CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    spotify_id VARCHAR(255) NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    artist_name VARCHAR(255) NOT NULL,
    loudness float NOT NULL,
    tempo float NOT NULL,
    time_sig int NOT NULL,
    acousticness float NOT NULL,
    danceability float NOT NULL,
    energy float NOT NULL,
    instrumentalness float NOT NULL,
    mode int NOT NULL,
    speechiness float NOT NULL,
    valence float NOT NULL
);