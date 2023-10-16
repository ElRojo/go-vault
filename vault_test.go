package main

import (
	"context"
	"fmt"
	"testing"

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

func (v *mockVault) addEnginesToList(ctx context.Context, client *vault.Client) []string {
	engineSlice := []string{}
	for eng := range v.getSecretEngines(ctx, client) {
		engineSlice = append(engineSlice, eng)
	}
	return engineSlice
}

func (v *mockVault) createEngines(ctx context.Context, client *vault.Client, secret *secret) string {
	return "Processed: " + secret.engine
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

func (v *mockVault) hydrateFlatSecretsStruct(ctx context.Context, c *vault.Client) {
	for _, secret := range newSecrets {
		for _, kv := range secret.Keys {
			for key := range kv.data {
				if kv.data[key] == "" {
					tr := secretsMap[key]
					if tr.path != "" {
						value := "secret"
						kv.data[key] = value
					}
				}
			}
		}
	}
}

func TestVaultLegacy(t *testing.T) {
	v := &mockVault{}
	c := &Config{
		Secrets: legacySecrets,
	}
	log.Info().Msg(runVault(v, *c))
}

func TestVaultNew(t *testing.T) {
	v := &mockVault{}
	c := &Config{
		Secrets: newSecrets,
	}
	log.Info().Msg(runVault(v, *c))
}

func TestLegacyVaultConfig(t *testing.T) {
	for _, s := range legacySecrets {
		ok := assert.IsType(t, &secret{}, s)
		if !ok {
			log.Warn().Msg("Secret stored in" + s.engine + " is malformed")
		}
	}
}

func TestFlatVaultConfig(t *testing.T) {
	for _, s := range newSecrets {
		ok := assert.IsType(t, &secret{}, s)
		if !ok {
			log.Warn().Msg("Secret stored in" + s.engine + " is malformed")
		}
	}
}
