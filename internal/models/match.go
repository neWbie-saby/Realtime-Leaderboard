package models

import (
	"github.com/neWbie-saby/leaderboard/internal/database"
)

func DatabaseCreatedMatchToMatch(dbMatch database.Match) CreateMatchOutput {
	return CreateMatchOutput{
		MatchID: dbMatch.ID,
	}
}
