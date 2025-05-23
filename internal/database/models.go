// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package database

import (
	"database/sql"
)

type Match struct {
	ID        int32
	MatchType string
	MatchDate sql.NullTime
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

type MatchUser struct {
	MatchID int32
	UserID  int32
	Score   int32
}

type MatchWinner struct {
	MatchID int32
	UserID  int32
}

type User struct {
	ID           int32
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
}
