package models

import (
	"time"

	"github.com/neWbie-saby/leaderboard/internal/database"
)

func DatabaseUserToUser(dbUser database.User) User {
	var createdAt time.Time
	var updatedAt time.Time

	if dbUser.CreatedAt.Valid {
		createdAt = dbUser.CreatedAt.Time
	}
	if dbUser.UpdatedAt.Valid {
		updatedAt = dbUser.UpdatedAt.Time
	}

	return User{
		ID:        dbUser.ID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		UserName:  dbUser.Username,
		Email:     dbUser.Email,
	}
}
