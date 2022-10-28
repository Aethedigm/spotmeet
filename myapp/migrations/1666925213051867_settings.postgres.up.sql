DROP TABLE IF EXISTS settings cascade;

CREATE TABLE settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    distance NUMERIC NOT NULL,
    looking_for VARCHAR(255) NOT NULL
);