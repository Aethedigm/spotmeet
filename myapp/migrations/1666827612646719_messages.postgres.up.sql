DROP TABLE IF EXISTS messages cascade;

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    match_id integer NOT NULL REFERENCES matches(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at timestamp without time zone NOT NULL DEFAULT now()
);