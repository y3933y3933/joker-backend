package errx

import "errors"

var (
	ErrGenerateCode       = errors.New("failed to generate unique game code")
	ErrGameNotFound       = errors.New("game not found")
	ErrInvalidGameStatus  = errors.New("invalid game status")
	ErrNotEnoughPlayers   = errors.New("not enough players")
	ErrRoundNotFound      = errors.New("round not found")
	ErrInvalidStatus      = errors.New("invalid round status")
	ErrForbidden          = errors.New("you are not allowed to perform this action")
	ErrPlayerNotFound     = errors.New("player not found")
	ErrDuplicateNickname  = errors.New("nickname already taken")
	ErrGameAlreadyStarted = errors.New("cannot leave, game has already started")
)
