-- +goose Up
CREATE TYPE match_role AS ENUM ('player', 'scorer', 'official');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    role match_role NOT NULL DEFAULT 'player'
);

-- +goose Down
DROP TABLE users;

DROP TYPE match_role;