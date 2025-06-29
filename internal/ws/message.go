package ws

import "encoding/json"

type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

const (
	MsgTypePlayerJoined    = "player_joined"
	MsgTypeGameStarted     = "game_started"
	MsgTypeRoundQuestion   = "round_question"
	MsgTypeAnswerTime      = "answer_time"
	MsgTypeAnswerSubmitted = "answer_submitted"
	MsgTypeJokerRevealed   = "joker_revealed"
	MsgTypePlayerSafe      = "player_safe"
	MsgTypeGameEnded       = "game_ended"
)

type PlayerJoinedPayload struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	IsHost   bool   `json:"isHost"`
}

type GameStartedPayload struct {
	RoundID          int64 `json:"roundID"`
	QuestionPlayerID int64 `json:"questionPlayerID"`
	AnswererID       int64 `json:"answererID"`
}

// NewWSMessage creates a new WSMessage from any data struct or map.
func NewWSMessage(msgType string, data any) (WSMessage, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return WSMessage{}, err
	}

	return WSMessage{
		Type: msgType,
		Data: b,
	}, nil
}
