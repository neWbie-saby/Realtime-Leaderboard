-- +goose Up
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    match_type TEXT NOT NULL,
    match_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE matches;