package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"

	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type vaulter interface {
	hydrateNewSecretsStruct(ctx context.Context, c *vault.Client)
	addEnginesToList(ctx context.Context, client *vault.Client) []string
	createEngines(ctx context.Context, client *vault.Client, secret *secret) string
	writeSecret(ctx context.Context, client *vault.Client, path string, data map[string]interface{}) error
	getSecretEngines(ctx context.Context, client *vault.Client) map[string]interface{}
}

type acmeVault struct{}

type Config struct {
	Secrets []*secret
}

func (v *acmeVault) addEnginesToList(ctx context.Context, client *vault.Client) []string {
	engineSlice := []string{}
	for eng := range v.getSecretEngines(ctx, client) {
		engineSlice = append(engineSlice, eng)
	}
	return engineSlice
}

func (v *acmeVault) createEngines(ctx context.Context, client *vault.Client, secret *secret) string {
	_, err := client.System.MountsEnableSecretsEngine(ctx, secret.engine, schema.MountsEnableSecretsEngineRequest{Type: "kv-v2"})
	if err != nil {
		log.Warn().Err(err).Msg("err in createEngines")
	}
	return "Processed: " + secret.engine
}

func (v *acmeVault) getSecretEngines(ctx context.Context, client *vault.Client) map[string]interface{} {
	engs, err := client.System.MountsListSecretsEngines(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	return engs.Data
}

func (v *acmeVault) writeSecret(ctx context.Context, client *vault.Client, path string, data map[string]interface{}) error {
	if _, err := client.Write(ctx, path, map[string]interface{}{"data": data}); err != nil {
		log.Fatal().Err(err).Msg("")
		return err
	}
	return nil
}

func ReadSecret(ctx context.Context, c *vault.Client, path string, secret string) (string, error) {
	response, err := c.Read(ctx, path)
	if err != nil {
		log.Warn().Err(err).Msg(path)
		return "er", err
	}

	data, _ := response.Data["data"].(map[string]interface{})
	return data[secret].(string), nil
}

func createDataInVault(ctx context.Context, client *vault.Client, v vaulter, c Config) {
	for _, secret := range c.Secrets {
		if slices.Contains(v.addEnginesToList(ctx, client), secret.engine+"/") {
			for _, kv := range secret.keys {
				path := fmt.Sprintf("/%v/data/%v", secret.engine, kv.path)
				err := v.writeSecret(ctx, client, path, kv.data)
				if err != nil {
					log.Warn().Err(err).Msg("")
				} else {
					log.Info().Msgf("Secrets in: " + path + " written")
				}
			}
		} else {
			log.Info().Msgf("engine: %q does not exist.", secret.engine)
		}
	}
}

func (v *acmeVault) hydrateNewSecretsStruct(ctx context.Context, c *vault.Client) {
	for _, secret := range newSecrets {
		for _, kv := range secret.keys {
			for key := range kv.data {
				if kv.data[key] == "" {
					sm := secretsMap[key]
					if sm.path != "" {
						value, err := ReadSecret(ctx, c, sm.path, sm.secret)
						if err != nil {
							log.Err(err).Msg("")
						}
						kv.data[key] = value
					}
				}
			}
		}
	}
}

func InitVaultClient() (context.Context, *vault.Client) {
	var (
		ctx       = context.Background()
		url       = getEnv().vaultUrl
		namespace = fmt.Sprintf("/admin/%s/%s", strings.ToLower(getEnv().vaultNs), getEnv().env)
		token     = getEnv().vaultTkn
	)

	client, err := vault.New(
		vault.WithAddress(url),
		vault.WithRequestTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	err = client.SetToken(token)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	err = client.SetNamespace(namespace)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	return ctx, client
}

func runVault(v vaulter, c Config, copy bool) string {
	ctx, client := InitVaultClient()
	s := getEnv().vaultConfig
	if copy && !s {
		v.hydrateNewSecretsStruct(ctx, client)
	}
	for _, secret := range c.Secrets {
		log.Info().Msg(v.createEngines(ctx, client, secret))
	}
	createDataInVault(ctx, client, v, c)

	return "Vault complete"

}
