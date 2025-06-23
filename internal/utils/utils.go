package utils

import (
	"fmt"
	"math/rand"
	"strconv"
)

func ParseID(s string) (int64, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %v", err)
	}
	return id, nil
}

func CreateShuffledDeck(jokerCount int, total int) []string {
	if jokerCount >= total || jokerCount < 0 || total <= 0 {
		panic("invalid deck configuration")
	}

	deck := make([]string, total)
	for i := 0; i < total; i++ {
		if i < jokerCount {
			deck[i] = "joker"
		} else {
			deck[i] = "safe"
		}
	}

	// 隨機洗牌
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}
