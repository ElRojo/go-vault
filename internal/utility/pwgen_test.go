package utility

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	var (
		pw1 string = GeneratePassword(10)
		pw2 string = GeneratePassword(10)
	)
	ok := assert.NotEqual(t, pw1, pw2)
	if !ok {
		log.Error().Msg("Passwords are the same!")
	}
}
