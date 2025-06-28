package errx

import "errors"

// 遊戲相關錯誤
var (
	ErrGenerateCode = errors.New("failed to generate unique game code")
)
