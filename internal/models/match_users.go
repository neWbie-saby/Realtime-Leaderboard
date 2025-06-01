package models

import (
	"github.com/neWbie-saby/leaderboard/internal/database"
)

func DatabaseUserScoreToUserScore(dbUserScore database.GetMatchUserScoresRow) UserScore {
	return UserScore{
		UserID: int(dbUserScore.UserID),
		Score:  int(dbUserScore.Score),
	}
}

func DatabaseUserScoresToUserScores(dbUserScores []database.GetMatchUserScoresRow) []UserScore {
	scores := []UserScore{}

	for _, score := range dbUserScores {
		scores = append(scores, DatabaseUserScoreToUserScore(score))
	}

	return scores
}
