package main

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func generatePassword(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune("abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()_+-=[]{}|;:,<.>/?")
	password := ""
	for i := 0; i < length; i++ {
		password += string([]rune(chars)[rand.Intn(len(chars))])
	}
	return password
}

func generateUuid() string {
	uuid := uuid.New().String()
	return uuid
}
