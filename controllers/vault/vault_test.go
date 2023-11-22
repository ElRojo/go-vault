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

func (v *mockVault) makeEngineSlice(ctx context.Context, client *vault.Client) ([]string, error) {
	engineSlice := []string{}
	for eng := range v.GetSecretEngine(ctx, client) {
		engineSlice = append(engineSlice, eng)
	}
	return engineSlice, nil
}

func (v *mockVault) createEngines(ctx context.Context, client *vault.Client, secret *Secret) (string, error) {
	return "Processed: " + secret.Engine, nil
}

func (v *mockVault) GetSecretEngine(ctx context.Context, client *vault.Client) map[string]interface{} {
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

func (v *mockVault) hydrateNewSecretsStruct(ctx context.Context, c *vault.Client, s []*Secret, secretMap map[string]secretMap) error {
	for _, secret := range s {
		for _, kv := range secret.Keys {
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
	ctx, client, err := v.InitVaultClient(c.Token, c.URL)
	if err != nil {
		return
	}
	r, err := InitVault(ctx, client, v, secrets, c)
	if err != nil {
		log.Err(err)
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
	ctx, client, err := v.InitVaultClient(c.Token, c.URL)
	if err != nil {
		return
	}
	r, err := InitVault(ctx, client, v, secrets, c)
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
	secrets := InitNewSecrets()
	ctx, client, err := v.InitVaultClient(c.Token, c.URL)
	if err != nil {
		log.Err(err)
	}
	r, err := InitVault(ctx, client, v, secrets, c)
	if err != nil {
		log.Err(err)
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
