// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: matches.sql

package database

import (
	"context"
)

const createMatch = `-- name: CreateMatch :one
INSERT INTO matches (match_type) VALUES ($1) RETURNING id, match_type, match_date, created_at, updated_at
`

func (q *Queries) CreateMatch(ctx context.Context, matchType string) (Match, error) {
	row := q.db.QueryRowContext(ctx, createMatch, matchType)
	var i Match
	err := row.Scan(
		&i.ID,
		&i.MatchType,
		&i.MatchDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMatchById = `-- name: GetMatchById :one
SELECT id, match_type, match_date, created_at, updated_at FROM matches WHERE id = $1
`

func (q *Queries) GetMatchById(ctx context.Context, id int32) (Match, error) {
	row := q.db.QueryRowContext(ctx, getMatchById, id)
	var i Match
	err := row.Scan(
		&i.ID,
		&i.MatchType,
		&i.MatchDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
