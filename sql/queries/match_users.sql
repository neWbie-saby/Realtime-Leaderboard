-- name: AddUserToMatch :exec
INSERT INTO match_users (match_id, user_id, score)
VALUES ($1, $2, $3);