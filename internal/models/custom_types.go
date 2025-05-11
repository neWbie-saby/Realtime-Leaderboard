package models

import "time"

type UserScore struct {
	UserID int `json:"user_id"`
	Score  int `json:"score"`
}

type RegisterInput struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type User struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
}

type UserNameInput struct {
	UserName string `json:"username"`
}

type UserNameOutput struct {
	ID       int32  `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type CreateMatchInput struct {
	MatchType string `json:"match_type"`
}

type CreateMatchOutput struct {
	MatchID int32 `json:"match_id"`
}

type ScoresInput struct {
	Scores []UserScore `json:"scores"`
}
