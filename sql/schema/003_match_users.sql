-- +goose Up
CREATE TABLE match_users (
    match_id INTEGER REFERENCES matches(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    score INTEGER NOT NULL,
    PRIMARY KEY (match_id, user_id)
);

-- +goose Down
DROP TABLE match_users;