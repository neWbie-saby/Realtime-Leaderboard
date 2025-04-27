-- +goose Up
CREATE TABLE match_winners (
    match_id INTEGER REFERENCES matches(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (match_id, user_id)
);

-- +goose Down
DROP TABLE match_winners;