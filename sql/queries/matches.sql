-- name: CreateMatch :one
INSERT INTO matches (match_type) VALUES ($1) RETURNING *;

-- name: GetMatchById :one
SELECT * FROM matches WHERE id = $1;

-- name: GetMatchIdsForAnalysis :many
SELECT 
    m.id 
FROM 
    matches m
LEFT JOIN 
    match_insights mi ON m.id = mi.match_id
WHERE 
    mi.match_id IS NULL OR mi.processed_at < NOW() - INTERVAL '2 minute'
LIMIT $1;