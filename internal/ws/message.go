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
	MsgNextRoundStarted    = "next_round_started"
	MsgPlayerLeft          = "player_left"
	MsgHostTransferred     = "host_transferred"
	MsgTypeRoundSkipped    = "round_skipped"
)

type PlayerJoinedPayload struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	IsHost   bool   `json:"isHost"`
}

type RoundStartedPayload struct {
	RoundID          int64 `json:"roundID"`
	QuestionPlayerID int64 `json:"questionPlayerID"`
	AnswererID       int64 `json:"answererID"`
}

type JokerRevealedPayload struct {
	Level   string `json:"level"`
	Content string `json:"content"`
}

type PlayerLeftPayload struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

type HostTransferredPayload struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

type AnswerSubmittedPayload struct {
	Answer string `json:"answer"`
}

type RoundSkippedPayload struct {
	Reason string `json:"reason"`
	RoundStartedPayload
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
