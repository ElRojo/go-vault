package api

import (
	"encoding/json"
	"go-vault/controllers/vault"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var expectedVaultStruct = []*vault.Secret{
	{
		Engine: "apiengine",
		KV: []struct {
			Data map[string]interface{}
			Path string
		}{
			{
				Data: map[string]interface{}{
					"api_key": "testvalue",
					"test":    "testing value2",
				},
				Path: "api-test",
			},
			{
				Data: map[string]interface{}{
					"another-Key": "val",
					"test":        "another value",
				},
				Path: "api-test-2",
			},
		},
	},
}
var mockPayload string = `[
	{
	  "engine": "apiengine",
	  "kv": [
		{
		  "path": "api-test",
		  "data": {
			"api_key": "testvalue",
			"test": "testing value2"
		  }
		},
		{
		  "path": "api-test-2",
		  "data": {
			"another-Key": "val",
			"test": "another value"
		  }
		}
	  ]
	}
  ]`

func TestSecretConversion(t *testing.T) {
	var v []*vault.Secret
	var secret []Secret
	json.Unmarshal([]byte(mockPayload), &secret)
	var vaultSecret, err = convertVaultSecret(secret)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	ok := assert.IsType(t, v, vaultSecret)
	if !ok {
		log.Error().Msgf("Type of %T does not match %T", vaultSecret, v)
	}
	ok = assert.Equal(t, expectedVaultStruct, vaultSecret)
	if !ok {
		log.Error().Msg("Mismatch in structs")
	}
}
