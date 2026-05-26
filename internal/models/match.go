package models

import (
	"time"

	"github.com/neWbie-saby/leaderboard/internal/database"
)

func DatabaseCreatedMatchToMatch(dbMatch database.Match) CreateMatchOutput {
	return CreateMatchOutput{
		MatchID: dbMatch.ID,
	}
}

func DatabasePlayerMatchToMatch(dbPlayerMatch database.GetUserPlayerMatchesRow) PlayerMatch {
	var matchDate time.Time

	if dbPlayerMatch.MatchDate.Valid {
		matchDate = dbPlayerMatch.MatchDate.Time
	}

	return PlayerMatch{
		MatchID:   dbPlayerMatch.ID,
		MatchType: dbPlayerMatch.MatchType,
		MatchDate: matchDate,
	}
}

func DatabasePlayerMatchesToMatches(dbPlayerMatches []database.GetUserPlayerMatchesRow) []PlayerMatch {
	matches := []PlayerMatch{}

	for _, match := range dbPlayerMatches {
		matches = append(matches, DatabasePlayerMatchToMatch(match))
	}

	return matches
}

func DatabasePlayerMatchScoreToScore(dbPlayerScore database.GetUserPlayerMatchScoreRow) PlayerMatchScore {
	var matchDate time.Time

	if dbPlayerScore.MatchDate.Valid {
		matchDate = dbPlayerScore.MatchDate.Time
	}

	return PlayerMatchScore{
		MatchID:    dbPlayerScore.ID,
		MatchType:  dbPlayerScore.MatchType,
		MatchDate:  matchDate,
		MatchScore: int(dbPlayerScore.Score),
	}
}
