package randutils

import (
	"math/rand"
	"time"
)

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(n int) string {
	result := make([]byte, n)
	rand.New(rand.NewSource(time.Now().UnixMicro()))

	for i := 0; i < n; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
