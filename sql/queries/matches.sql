-- name: CreateMatch :one
INSERT INTO matches (match_type)
SELECT $2
FROM
    users
WHERE
    users.id = $1 AND users.role = 'official'
LIMIT 1
RETURNING matches.*;

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

-- name: GetUserPlayerMatches :many
SELECT
    m.id, m.match_type, m.match_date
FROM matches m
JOIN
    match_users mu ON mu.match_id = m.id
JOIN
    users u ON u.id = mu.user_id
WHERE
    u.role = 'player'
AND
    mu.user_id = $1;

-- name: GetUserPlayerMatchScore :one
SELECT
    m.id, m.match_type, m.match_date, mu.score
FROM matches m
JOIN
    match_users mu ON mu.match_id = m.id
JOIN
    users u ON u.id = mu.user_id
WHERE
    u.role = 'player'
AND
    u.id = $1
AND
    m.id = $2;
