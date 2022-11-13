DROP TABLE IF EXISTS user_music_profile CASCADE;

CREATE TABLE user_music_profile (
    id serial PRIMARY KEY,
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    loudness float NOT NULL,
    tempo float NOT NULL,
    time_sig int NOT NULL,
    updated_at timestamp without time zone NOT NULL DEFAULT now()
)