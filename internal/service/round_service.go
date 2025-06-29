package service

import (
	"context"

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

	// ✅ 按加入順序取前兩位
	questioner := players[0]
	answerer := players[1]

	round := &store.Round{
		GameID:           game.ID,
		QuestionPlayerID: questioner.ID,
		AnswerPlayerID:   answerer.ID,
		Status:           store.RoundStatusWaitingForQuestion,
		Deck:             generateDeck(DECK_LENGTH), // 自行實作
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
