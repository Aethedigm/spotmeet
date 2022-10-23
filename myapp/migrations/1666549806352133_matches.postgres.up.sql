drop table if exists matches cascade;

CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    user_a_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_b_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    percent_match float,
    artist_id integer,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    expiry timestamp without time zone NOT NULL DEFAULT now() + interval '5 day'
);