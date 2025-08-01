-- name: AcquireLock :execrows
INSERT INTO processing_locks (resource_type, resource_id, locked_by, expires_at, released_at)
VALUES ($1, $2, $3, NOW() + INTERVAL '5 minute', NULL)
ON CONFLICT (resource_type, resource_id) WHERE released_at IS NULL DO NOTHING;

-- name: ReleaseLock :exec
UPDATE processing_locks 
SET released_at = NOW()
WHERE resource_type = $1
AND resource_id = $2
AND locked_by = $3;

-- name: DeleteExpiredLocks :execrows
DELETE FROM processing_locks
WHERE expires_at <= NOW() AND released_at IS NULL;