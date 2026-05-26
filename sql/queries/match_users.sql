-- name: AddUserToMatch :execrows
INSERT INTO match_users (match_id, user_id, score)
SELECT 
    $1, $2, $3
FROM users
WHERE 
    users.id = $4
AND users.role IN ('scorer', 'official');

-- name: MatchUserExists :one
SELECT EXISTS (
    SELECT 1 FROM match_users WHERE user_id = $1 AND match_id = $2
);

-- name: UpdateUserScoreOnMatch :execrows
UPDATE match_users
SET score = $1
WHERE user_id = $2 AND match_id = $3
AND EXISTS (
    SELECT 1
    FROM users
    WHERE users.id = $4
    AND users.role IN ('scorer', 'official')
);

-- name: GetMatchUserScores :many
SELECT user_id, score FROM match_users WHERE match_id = $1;

-- name: GetMatchUserScoresAndUserNames :many
SELECT u.id, u.username, mu.score 
FROM match_users mu 
JOIN users u ON u.id = mu.user_id 
WHERE mu.match_id = $1
AND EXISTS(
    SELECT 1
    FROM users req_u
    LEFT JOIN match_users req_mu ON req_u.id = req_mu.user_id AND req_mu.match_id = $1
    WHERE req_u.id = $2
    AND (
        (req_mu.user_id IS NOT NULL AND req_u.role = 'player')
        OR
        (req_mu.user_id IS NULL AND req_u.role IN ('scorer', 'official'))
    )
);