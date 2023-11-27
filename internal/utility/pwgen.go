package utility

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

func GeneratePassword(length int) string {
	chars := []rune("abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()_+-=[]{}|;:,<.>/?")
	fmt.Print(len(chars))
	password := make([]rune, length)
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return ""
		}
		password[i] = chars[randomIndex.Int64()]
	}
	return string(password)
}

func GenerateUUID() string {
	uuid := uuid.New().String()
	return uuid
}
