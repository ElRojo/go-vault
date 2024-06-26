package vault

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault-client-go"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type engineDataMock struct {
	data map[string]interface{}
}
type mockVault struct{}

var resp = engineDataMock{
	data: map[string]interface{}{
		"test-engine/":   "",
		"test-engine-1/": "",
		"test-engine-2/": "",
		"my-folder/":     "",
	},
}

func (v *mockVault) createEngines(ctx context.Context, client *vault.Client, secret *Secret) (string, error) {
	return "Processed: " + secret.Engine, nil
}

func (v *mockVault) getSecretEngine(ctx context.Context, client *vault.Client) ([]string, error) {
	engineSlice := []string{}
	for eng := range resp.data {
		engineSlice = append(engineSlice, eng)
	}
	return engineSlice, nil
}

func (v *mockVault) writeSecret(ctx context.Context, client *vault.Client, path string, data map[string]interface{}) error {
	for i, v := range data {
		fmt.Printf("Writing: %v:%v \n", i, v)
	}
	return nil
}

func (v *mockVault) hydrateNewSecretsStruct(ctx context.Context, c *vault.Client, s []*Secret, secretMap map[string]secretMap) error {
	for _, secret := range s {
		for _, kv := range secret.KV {
			for key := range kv.Data {
				if kv.Data[key] == "" {
					sm := secretMap[key]
					if sm.path != "" {
						value := "copiedSecret"
						kv.Data[key] = value
					}
				}
			}
		}
	}
	return nil
}

func TestVaultLegacy(t *testing.T) {
	v := &mockVault{}
	c := VaultConfig{
		Copy:   false,
		Legacy: true,
		Token:  "",
		URL:    "",
	}
	secrets := InitLegacySecrets()
	ctx, client, err := InitVaultClient(c.Token, c.URL)
	if err != nil {
		return
	}
	r, err := InitVault(ctx, client, v, secrets, c)
	if err != nil {
		log.Err(err).Msg("")
	}
	log.Info().Msg(r)
}

func TestVaultNew(t *testing.T) {
	v := &mockVault{}
	c := VaultConfig{
		Copy:   false,
		Legacy: true,
		Token:  "",
		URL:    "",
	}
	secrets := InitNewSecrets()
	ctx, client, err := InitVaultClient(c.Token, c.URL)
	if err != nil {
		return
	}
	r, err := InitVault(ctx, client, v, secrets, c)
	if err != nil {
		log.Err(err).Msg("")
	}
	log.Info().Msg(r)
}

func TestVaultNewWithCopy(t *testing.T) {
	v := &mockVault{}
	c := VaultConfig{
		Copy:   true,
		Legacy: false,
		Token:  "",
		URL:    "",
	}
	secrets := InitNewSecrets()
	ctx, client, err := InitVaultClient(c.Token, c.URL)
	if err != nil {
		log.Err(err).Msg("")
	}
	r, err := InitVault(ctx, client, v, secrets, c)
	if err != nil {
		log.Err(err).Msg("")
	}
	log.Info().Msg(r)
}

func TestLegacyVaultConfig(t *testing.T) {
	secrets := InitLegacySecrets()
	for _, s := range secrets {
		ok := assert.IsType(t, &Secret{}, s)
		if !ok {
			log.Warn().Msg("Secret stored in" + s.Engine + " is malformed")
		}
	}
}

func TestFlatVaultConfig(t *testing.T) {
	secrets := InitNewSecrets()
	for _, s := range secrets {
		ok := assert.IsType(t, &Secret{}, s)
		if !ok {
			log.Warn().Msg("Secret stored in" + s.Engine + " is malformed")
		}
	}
}
