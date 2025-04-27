-- name: CreateMatch :one
INSERT INTO matches (match_type) VALUES ($1) RETURNING *;

-- name: GetMatchById :one
SELECT * FROM matches WHERE id = $1;