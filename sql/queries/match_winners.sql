-- name: AddWinnerOfMatch :execrows
INSERT INTO match_winners (match_id, user_id)
SELECT 
    $1, $2
FROM users
WHERE
    users.id = $3
AND users.role IN ('scorer', 'official');

-- name: GetMatchWinners :many
SELECT user_id FROM match_winners WHERE match_id = $1;

-- name: DeleteMatchWinners :execrows
DELETE FROM match_winners 
WHERE match_id = $1
AND EXISTS (
    SELECT 1
    FROM users
    WHERE users.id = $2
    AND users.role IN ('scorer', 'official')
);