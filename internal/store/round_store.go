package store

import (
	"context"

	"github.com/y3933y3933/joker/internal/db/sqlc"
)

type Round struct {
	ID               int64    `json:"id"`
	GameID           int64    `json:"gameID"`
	QuestionID       *int64   `json:"questionID,omitempty"` // 尚未選題前為 nil
	Answer           *string  `json:"answer,omitempty"`     // 尚未回答前為 nil
	QuestionPlayerID int64    `json:"questionerID"`
	AnswerPlayerID   int64    `json:"answererID"`
	IsJoker          bool     `json:"isJoker"`
	Status           string   `json:"status"`
	Deck             []string `json:"-"`
}

const (
	RoundStatusWaitingForQuestion = "waiting_for_question"
	RoundStatusWaitingForAnswer   = "waiting_for_answer"
	RoundStatusWaitingForDraw     = "waiting_for_draw"
	RoundStatusRevealed           = "revealed"
	RoundStatusDone               = "done"
)

type PostgresRoundStore struct {
	queries *sqlc.Queries
}

func NewPostgresRoundStore(queries *sqlc.Queries) *PostgresRoundStore {
	return &PostgresRoundStore{queries: queries}
}

type RoundStore interface {
	Create(ctx context.Context, round *Round) (*Round, error)
}

func (pg *PostgresRoundStore) Create(ctx context.Context, round *Round) (*Round, error) {
	arg := sqlc.CreateRoundParams{
		GameID:           round.GameID,
		QuestionID:       toPgInt8(round.QuestionID),
		Answer:           toPgText(round.Answer),
		QuestionPlayerID: round.QuestionPlayerID,
		AnswerPlayerID:   round.AnswerPlayerID,
		IsJoker:          toPgBool(&round.IsJoker),
		Status:           string(round.Status),
		Deck:             round.Deck,
	}

	res, err := pg.queries.CreateRound(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &Round{
		ID:               res.ID,
		GameID:           res.GameID,
		QuestionPlayerID: res.QuestionPlayerID,
		AnswerPlayerID:   res.AnswerPlayerID,
		Status:           res.Status,
		QuestionID:       &res.QuestionID.Int64,
		Answer:           fromPgText(res.Answer),
		IsJoker:          fromPgBool(res.IsJoker),
		Deck:             res.Deck,
	}, nil
}
