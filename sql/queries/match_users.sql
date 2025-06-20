-- name: AddUserToMatch :exec
INSERT INTO match_users (match_id, user_id, score)
VALUES ($1, $2, $3);

-- name: MatchUserExists :one
SELECT EXISTS (
    SELECT 1 FROM match_users WHERE user_id = $1 AND match_id = $2
);

-- name: UpdateUserScoreOnMatch :exec
UPDATE match_users
SET score = $1
WHERE user_id = $2 AND match_id = $3;

-- name: GetMatchUserScores :many
SELECT user_id, score FROM match_users WHERE match_id = $1;

-- name: GetMatchUserScoresAndUserNames :many
SELECT u.id, u.username, mu.score 
FROM match_users mu 
JOIN users u ON u.id = mu.user_id 
WHERE mu.match_id = $1;