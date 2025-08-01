-- +goose Up
CREATE TABLE processing_locks (
    id BIGSERIAL PRIMARY KEY,
    resource_type VARCHAR(50) NOT NULL,
    resource_id VARCHAR(255) NOT NULL,
    locked_by VARCHAR(255) NOT NULL,
    locked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    released_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unique_active_processing_lock
ON processing_locks (resource_type, resource_id)
WHERE released_at IS NULL;

CREATE INDEX idx_processing_locks_expires_at ON processing_locks (expires_at);

-- +goose Down
DROP TABLE processing_locks;