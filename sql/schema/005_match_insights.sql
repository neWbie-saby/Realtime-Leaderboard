-- +goose Up
CREATE TABLE match_insights (
    match_id INTEGER PRIMARY KEY REFERENCES matches(id),
    highest_score INTEGER NOT NULL,
    average_score FLOAT NOT NULL,
    total_participants INTEGER NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE match_insights;