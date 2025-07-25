// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Feedback struct {
	ID           int64
	Type         string
	Content      string
	CreatedAt    pgtype.Timestamptz
	ReviewStatus string
}

type Game struct {
	ID        int64
	Code      string
	Status    string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type Player struct {
	ID       int64
	GameID   int64
	Nickname string
	IsHost   pgtype.Bool
	JoinedAt pgtype.Timestamptz
	Status   string
}

type Question struct {
	ID        int64
	Level     string
	Content   string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type Round struct {
	ID               int64
	GameID           int64
	QuestionID       pgtype.Int8
	Answer           pgtype.Text
	QuestionPlayerID int64
	AnswerPlayerID   int64
	IsJoker          pgtype.Bool
	Status           string
	CreatedAt        pgtype.Timestamptz
	Deck             []string
}

type User struct {
	ID           int64
	Username     string
	PasswordHash []byte
	CreatedAt    pgtype.Timestamptz
}
