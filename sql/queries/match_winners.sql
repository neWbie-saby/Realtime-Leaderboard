-- name: AddWinnerOfMatch :exec
INSERT INTO match_winners (match_id, user_id)
VALUES ($1, $2);

-- name: GetMatchWinners :many
SELECT user_id FROM match_winners WHERE match_id = $1;

-- name: DeleteMatchWinners :exec
DELETE FROM match_winners WHERE match_id = $1;