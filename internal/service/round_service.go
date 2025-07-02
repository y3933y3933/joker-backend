package service

import (
	"context"
	"errors"

	"math/rand"

	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/errx"
)

const DECK_LENGTH = 3
const MIN_GAME_NUM = 3

type RoundService struct {
	roundStore  store.RoundStore
	playerStore store.PlayerStore
	gameStore   store.GameStore
}

func NewRoundService(roundStore store.RoundStore, playerStore store.PlayerStore, gameStore store.GameStore) *RoundService {
	return &RoundService{
		roundStore:  roundStore,
		playerStore: playerStore,
		gameStore:   gameStore,
	}
}

func (s *RoundService) StartGame(ctx context.Context, game *store.Game) (*store.Round, error) {
	if game.Status != store.GameStatusWaiting {
		return nil, errx.ErrInvalidGameStatus
	}
	players, err := s.playerStore.FindPlayersByGameID(ctx, game.ID)
	if err != nil {
		return nil, err
	}

	if len(players) < MIN_GAME_NUM {
		return nil, errx.ErrNotEnoughPlayers
	}

	round, err := s.generateRound(ctx, game.ID, players)
	if err != nil {
		return nil, err
	}

	created, err := s.roundStore.Create(ctx, round)
	if err != nil {
		return nil, err
	}

	// 更新遊戲狀態為 playing
	err = s.gameStore.UpdateStatus(ctx, game.ID, store.GameStatusPlaying)
	if err != nil {
		return nil, err
	}

	return created, nil

}

func generateDeck(n int) []string {
	if n < 1 {
		return []string{}
	}

	deck := make([]string, n)

	// 隨機放一張 joker
	jokerIndex := rand.Intn(n)
	for i := 0; i < n; i++ {
		if i == jokerIndex {
			deck[i] = "joker"
		} else {
			deck[i] = "safe"
		}
	}

	// 打亂牌組
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}

func (s *RoundService) SubmitQuestion(ctx context.Context, roundID int64, questionID int64, playerID int64) error {
	round, err := s.roundStore.GetRoundByID(ctx, roundID)
	if err != nil {
		return err
	}

	if round.QuestionPlayerID != playerID {
		return errx.ErrForbidden
	}

	if round.Status != store.RoundStatusWaitingForQuestion {
		return errx.ErrInvalidStatus
	}

	return s.roundStore.SetRoundQuestion(ctx, roundID, questionID)
}

func (s *RoundService) GetRoundWithQuestion(ctx context.Context, roundID int64) (*store.RoundWithQuestion, error) {
	round, err := s.roundStore.GetRoundWithQuestion(ctx, roundID)
	if err != nil {
		return nil, err
	}
	return round, nil
}

func (s *RoundService) SubmitAnswer(ctx context.Context, roundID int64, answer string, playerID int64) error {
	round, err := s.roundStore.GetRoundByID(ctx, roundID)
	if err != nil {
		return err
	}

	// 驗證身份與狀態
	if round.AnswerPlayerID != playerID {
		return errx.ErrForbidden
	}
	if round.Status != store.RoundStatusWaitingForAnswer {
		return errx.ErrInvalidStatus
	}

	// 更新回答與狀態
	err = s.roundStore.UpdateAnswer(ctx, roundID, answer, store.RoundStatusWaitingForDraw)
	if err != nil {
		return err
	}
	return nil
}

func (s *RoundService) DrawCard(ctx context.Context, roundID, playerID int64, index int) (*store.RoundWithQuestion, error) {
	round, err := s.roundStore.GetRoundByID(ctx, roundID)
	if err != nil {
		return nil, err
	}

	if round.Status != store.RoundStatusWaitingForDraw {
		return nil, errx.ErrInvalidStatus
	}
	if round.AnswerPlayerID != playerID {
		return nil, errx.ErrForbidden
	}
	if index < 0 || index >= len(round.Deck) {
		return nil, errors.New("invalid card index")
	}

	card := round.Deck[index]
	isJoker := card == "joker"

	var newStatus string
	if isJoker {
		newStatus = store.RoundStatusRevealed
	} else {
		newStatus = store.RoundStatusDone
	}

	err = s.roundStore.UpdateDrawResult(ctx, roundID, isJoker, newStatus)
	if err != nil {
		return nil, err
	}

	return s.roundStore.GetRoundWithQuestion(ctx, roundID)
}

func (s *RoundService) CreateNextRound(ctx context.Context, game *store.Game) (*store.Round, error) {
	players, err := s.playerStore.FindPlayersByGameID(ctx, game.ID)
	if err != nil {
		return nil, err
	}
	if len(players) < MIN_GAME_NUM {
		return nil, errx.ErrNotEnoughPlayers
	}

	round, err := s.generateRound(ctx, game.ID, players)
	if err != nil {
		return nil, err
	}

	created, err := s.roundStore.Create(ctx, round)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *RoundService) generateRound(ctx context.Context, gameID int64, players []*store.Player) (*store.Round, error) {
	// 找出上一輪
	lastRound, err := s.roundStore.FindLastRoundByGameID(ctx, gameID)
	if err != nil && !errors.Is(err, errx.ErrRoundNotFound) {
		return nil, err
	}

	var questioner, answerer *store.Player

	if lastRound == nil {
		// 第一輪：取前兩位
		questioner, answerer = players[0], players[1]
	} else {
		questioner, answerer = getNextPair(players, lastRound.QuestionPlayerID)
	}

	round := &store.Round{
		GameID:           gameID,
		QuestionPlayerID: questioner.ID,
		AnswerPlayerID:   answerer.ID,
		Status:           store.RoundStatusWaitingForQuestion,
		Deck:             generateDeck(DECK_LENGTH),
	}

	return round, nil
}

func getNextPair(players []*store.Player, lastQuestionerID int64) (questioner, answerer *store.Player) {
	n := len(players)
	var qIndex int

	// 找出上一輪 questioner 的位置
	for i, p := range players {
		if p.ID == lastQuestionerID {
			qIndex = (i + 1) % n // 下一位成為新的 questioner
			break
		}
	}

	aIndex := (qIndex + 1) % n // 下一位成為回答者
	return players[qIndex], players[aIndex]
}

func (s *RoundService) SkipRound(ctx context.Context, game *store.Game, roundID int64) error {
	err := s.roundStore.UpdateRoundStatus(ctx, roundID, store.RoundStatusDone)
	if err != nil {
		return err
	}

	_, err = s.CreateNextRound(ctx, game)
	if err != nil {
		return err
	}

	return nil

}
