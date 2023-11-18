package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type engineDataMock struct {
	engine map[string]interface{}
}
type responseMock[T any] struct {
	data T
}
type mockVault struct{}

var resp = responseMock[engineDataMock]{
	data: engineDataMock{
		engine: map[string]interface{}{
			"test-engine/":   "",
			"test-engine-1/": "",
			"test-engine-2/": "",
			"my-folder/":     "",
		},
	},
}

func (v *mockVault) addEnginesToList(ctx context.Context, client *vault.Client) ([]string, error) {
	engineSlice := []string{}
	for eng := range v.getSecretEngines(ctx, client) {
		engineSlice = append(engineSlice, eng)
	}
	return engineSlice, nil
}

func (v *mockVault) createEngines(ctx context.Context, client *vault.Client, secret *secret) (string, error) {
	return "Processed: " + secret.engine, nil
}

func (v *mockVault) getSecretEngines(ctx context.Context, client *vault.Client) map[string]interface{} {
	return resp.data.engine
}

func (v *mockVault) writeSecret(ctx context.Context, client *vault.Client, path string, data map[string]interface{}) error {
	for i, v := range data {
		fmt.Printf("Writing: %v:%v \n", i, v)
	}
	return nil
}

func (v *mockVault) InitVaultClient(token string, url string) (context.Context, *vault.Client, error) {
	var ctx = context.Background()

	client, err := vault.New(
		vault.WithAddress(url),
		vault.WithRequestTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatal().Err(err)
		return nil, nil, err
	}
	err = client.SetToken(token)
	if err != nil {
		log.Fatal().Err(err)
		return nil, nil, err
	}

	return ctx, client, nil
}

func (v *mockVault) hydrateNewSecretsStruct(ctx context.Context, c *vault.Client, s []*secret, secretMap map[string]secretMap) {
	for _, secret := range s {
		for _, kv := range secret.keys {
			for key := range kv.data {
				if kv.data[key] == "" {
					sm := secretMap[key]
					if sm.path != "" {
						value := "copiedSecret"
						kv.data[key] = value
					}
				}
			}
		}
	}
}

func TestVaultLegacy(t *testing.T) {
	v := &mockVault{}
	c := VaultConfig{
		Copy:   false,
		Legacy: true,
		Token:  "",
		URL:    "",
	}
	r, err := RunVault(v, c)
	if err != nil {
		log.Err(err)
	}
	log.Info().Msg(r)
}

func TestVaultNew(t *testing.T) {
	v := &mockVault{}
	c := VaultConfig{
		Copy:   false,
		Legacy: false,
		Token:  "",
		URL:    "",
	}
	r, err := RunVault(v, c)
	if err != nil {
		log.Err(err)
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
	r, err := RunVault(v, c)
	if err != nil {
		log.Err(err)
	}
	log.Info().Msg(r)
}

func TestLegacyVaultConfig(t *testing.T) {
	secrets := initLegacySecrets()
	for _, s := range secrets {
		ok := assert.IsType(t, &secret{}, s)
		if !ok {
			log.Warn().Msg("Secret stored in" + s.engine + " is malformed")
		}
	}
}

func TestFlatVaultConfig(t *testing.T) {
	secrets := initNewSecrets()
	for _, s := range secrets {
		ok := assert.IsType(t, &secret{}, s)
		if !ok {
			log.Warn().Msg("Secret stored in" + s.engine + " is malformed")
		}
	}
}
