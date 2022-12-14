DROP TABLE IF EXISTS recovery_emails CASCADE;

CREATE TABLE recovery_emails (
    id SERIAL PRIMARY KEY,
    userid int NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);