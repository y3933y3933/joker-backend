package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/y3933y3933/joker/internal/db/sqlc"
	"github.com/y3933y3933/joker/internal/utils/errx"
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

type RoundWithQuestion struct {
	ID               int64
	GameID           int64
	QuestionID       *int64
	Answer           *string
	QuestionPlayerID int64
	AnswerPlayerID   int64
	Status           string
	Deck             []string
	QuestionContent  string
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
	SetRoundQuestion(ctx context.Context, roundID int64, questionID int64) error
	GetRoundByID(ctx context.Context, roundID int64) (*Round, error)
	GetRoundWithQuestion(ctx context.Context, id int64) (*RoundWithQuestion, error)
	UpdateAnswer(ctx context.Context, roundID int64, answer string, status string) error
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

func (pg *PostgresRoundStore) SetRoundQuestion(ctx context.Context, roundID int64, questionID int64) error {
	return pg.queries.SetRoundQuestion(ctx, sqlc.SetRoundQuestionParams{
		QuestionID: toPgInt8(&questionID),
		ID:         roundID,
	})
}

func (s *PostgresRoundStore) GetRoundByID(ctx context.Context, roundID int64) (*Round, error) {
	res, err := s.queries.GetRoundByID(ctx, roundID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrRoundNotFound
		}
		return nil, err
	}

	return &Round{
		ID:               res.ID,
		GameID:           res.GameID,
		QuestionID:       fromPgInt8(res.QuestionID),
		Answer:           fromPgText(res.Answer),
		QuestionPlayerID: res.QuestionPlayerID,
		AnswerPlayerID:   res.AnswerPlayerID,
		IsJoker:          res.IsJoker.Bool,
		Status:           res.Status,
		Deck:             res.Deck,
	}, nil
}

func (pg *PostgresRoundStore) GetRoundWithQuestion(ctx context.Context, id int64) (*RoundWithQuestion, error) {
	res, err := pg.queries.GetRoundWithQuestion(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrRoundNotFound
		}
		return nil, err
	}

	return &RoundWithQuestion{
		ID:               res.ID,
		GameID:           res.GameID,
		QuestionID:       fromPgInt8(res.QuestionID),
		Answer:           fromPgText(res.Answer),
		QuestionPlayerID: res.QuestionPlayerID,
		AnswerPlayerID:   res.AnswerPlayerID,
		Status:           res.Status,
		Deck:             res.Deck,
		QuestionContent:  res.QuestionContent,
	}, nil
}

func (pg *PostgresRoundStore) UpdateAnswer(ctx context.Context, roundID int64, answer string, status string) error {
	args := sqlc.UpdateAnswerParams{
		ID:     roundID,
		Answer: toPgText(&answer),
		Status: status,
	}
	err := pg.queries.UpdateAnswer(ctx, args)
	return err
}
