package vault

import (
	"context"
	"fmt"
	"slices"

	"github.com/rs/zerolog/log"

	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type Vaulter interface {
	addEnginesToList(ctx context.Context, client *vault.Client) ([]string, error)
	createEngines(ctx context.Context, client *vault.Client, secret *secret) (string, error)
	getSecretEngines(ctx context.Context, client *vault.Client) map[string]interface{}
	hydrateNewSecretsStruct(ctx context.Context, c *vault.Client, s []*secret, secretMap map[string]secretMap)
	InitVaultClient(token string, url string) (context.Context, *vault.Client, error)
	writeSecret(ctx context.Context, client *vault.Client, path string, data map[string]interface{}) error
}

func (v *AcmeVault) addEnginesToList(ctx context.Context, client *vault.Client) ([]string, error) {
	engineSlice := []string{}
	for eng := range v.getSecretEngines(ctx, client) {
		engineSlice = append(engineSlice, eng)
	}
	return engineSlice, nil
}

func (v *AcmeVault) createEngines(ctx context.Context, client *vault.Client, secret *secret) (string, error) {
	_, err := client.System.MountsEnableSecretsEngine(ctx, secret.engine, schema.MountsEnableSecretsEngineRequest{Type: "kv-v2"})
	if err != nil {
		log.Warn().Err(err).Msg("err in createEngines")
		return "", err
	}
	return "Processed: " + secret.engine, nil
}

func (v *AcmeVault) getSecretEngines(ctx context.Context, client *vault.Client) map[string]interface{} {
	engs, err := client.System.MountsListSecretsEngines(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	return engs.Data
}

func (v *AcmeVault) writeSecret(ctx context.Context, client *vault.Client, path string, data map[string]interface{}) error {
	if _, err := client.Write(ctx, path, map[string]interface{}{"data": data}); err != nil {
		log.Fatal().Err(err).Msg("")
		return err
	}
	return nil
}

func ReadSecret(ctx context.Context, c *vault.Client, path string, secret string) string {
	response, err := c.Read(ctx, path)
	if err != nil {
		return "error reading secret"
	}
	data, ok := response.Data["data"].(map[string]interface{})
	if !ok {
		return fmt.Sprintf("response data for %s secret not ok", secret)
	}
	return data[secret].(string)
}

func createDataInVault(ctx context.Context, client *vault.Client, v Vaulter, s []*secret) error {
	for _, secret := range s {
		engines, err := v.addEnginesToList(ctx, client)
		if err != nil {
			return err
		}
		if slices.Contains(engines, secret.engine+"/") {
			for _, kv := range secret.keys {
				path := fmt.Sprintf("%v/data/%v", secret.engine, kv.path)
				if err := v.writeSecret(ctx, client, path, kv.data); err != nil {
					log.Warn().Err(err)
				} else {
					log.Info().Msgf("Secrets in: %q written", path)
				}
			}
		} else {
			log.Info().Msgf("Engine: %q does not exist.", secret.engine)
		}
	}
	return nil
}

func (v *AcmeVault) hydrateNewSecretsStruct(ctx context.Context, c *vault.Client, s []*secret, secretMap map[string]secretMap) {
	for _, secret := range s {
		for _, kv := range secret.keys {
			for key := range kv.data {
				sm := secretMap[key]
				if sm.path != "" {
					value := ReadSecret(ctx, c, sm.path, sm.secret)
					kv.data[key] = value
				}
			}
		}
	}
}

func (s AcmeVault) InitVaultClient(token string, url string) (context.Context, *vault.Client, error) {
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

func RunVault(v Vaulter, c VaultConfig) (string, error) {
	ctx, client, err := v.InitVaultClient(c.Token, c.URL)

	if err != nil {
		return "", err
	}

	var s []*secret
	if c.Legacy {
		s = initLegacySecrets()
	} else {
		s = initNewSecrets()
	}

	if !c.Legacy && c.Copy {
		sm := initSecretMap()
		v.hydrateNewSecretsStruct(ctx, client, s, sm)
	}

	for _, secret := range s {
		e, err := v.createEngines(ctx, client, secret)
		if err != nil {
			return "", err
		}
		log.Info().Msg(e)
	}
	err = createDataInVault(ctx, client, v, s)
	if err != nil {
		return "", err
	}
	return "Vault complete", nil

}
