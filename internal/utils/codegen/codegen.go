package codegen

import "math/rand"

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
