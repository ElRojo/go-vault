package utility

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type vaultRequestTest struct {
	CopyLegacy *bool  `json:"copyLegacy" validate:"required"`
	URL        string `json:"vaultUrl" validate:"required"`
	UseLegacy  *bool  `json:"useLegacy" validate:"required"`
	VaultToken string `json:"vaultToken" validate:"required"`
}

func TestValidateRequestFields(t *testing.T) {
	testVaultReq := &vaultRequestTest{}
	req := `{
		"useLegacy":  true,
		"copyLegacy": true,
		"vaultToken": "dev-only-token",
		"vaultUrl":   "http://vault:8200"
	}`
	if err := json.NewDecoder(strings.NewReader(req)).Decode(&testVaultReq); err != nil {
		log.Error().Err(err)
	}
	err := ValidateRequestFields(testVaultReq)
	if ok := assert.NoError(t, err); !ok {
		log.Error().Err(err)
	}
}

func TestValidateRequestFieldsFail(t *testing.T) {
	testVaultReq := &vaultRequestTest{}
	req := `{}`
	if err := json.NewDecoder(strings.NewReader(req)).Decode(&testVaultReq); err != nil {
		log.Error().Err(err)
	}
	err := ValidateRequestFields(testVaultReq)
	if ok := assert.Error(t, err, err); !ok {
		log.Error().Msg("Fields did not fail and should!")
	}
}
