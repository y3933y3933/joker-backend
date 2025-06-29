package errx

import "errors"

// 遊戲相關錯誤
var (
	ErrGenerateCode      = errors.New("failed to generate unique game code")
	ErrGameNotFound      = errors.New("game not found")
	ErrInvalidGameStatus = errors.New("invalid game status")
	ErrNotEnoughPlayers  = errors.New("not enough players")
)

// 回合相關錯誤
var (
	ErrRoundNotFound = errors.New("round not found")
	ErrInvalidStatus = errors.New("invalid round status")
	ErrForbidden     = errors.New("you are not allowed to perform this action")
)
