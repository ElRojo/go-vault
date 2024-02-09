package vault

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type vaultFunc func() (*vault.Response[map[string]interface{}], error)

type Vaulter interface {
	createEngines(ctx context.Context, client *vault.Client, secret *Secret) (string, error)
	getSecretEngine(ctx context.Context, client *vault.Client) ([]string, error)
	hydrateNewSecretsStruct(ctx context.Context, c *vault.Client, s []*Secret, secretMap map[string]secretMap) error
	writeSecret(ctx context.Context, client *vault.Client, path string, data map[string]interface{}) error
}

func InitVaultClient(token string, url string) (context.Context, *vault.Client, error) {
	ctx := context.Background()

	client, err := vault.New(
		vault.WithAddress(url),
		vault.WithRequestTimeout(10*time.Second),
	)
	if err != nil {
		log.Warn().Err(err)
		return nil, nil, fmt.Errorf("vault client initialization failed: %w", err)

	}
	err = client.SetToken(token)
	if err != nil {
		log.Warn().Err(err)
		return nil, nil, fmt.Errorf("setting vault token failed: %w", err)
	}

	return ctx, client, nil
}

func InitVault(ctx context.Context, client *vault.Client, v Vaulter, s []*Secret, c VaultConfig) (string, error) {
	if !c.Legacy && c.Copy {
		sm := initSecretMap()
		v.hydrateNewSecretsStruct(ctx, client, s, sm)
	}

	err := CreateDataInVault(ctx, client, v, s)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func (v *AcmeVault) createEngines(ctx context.Context, client *vault.Client, secret *Secret) (string, error) {
	_, err := client.System.MountsEnableSecretsEngine(ctx, secret.Engine, schema.MountsEnableSecretsEngineRequest{Type: "kv-v2"})
	if err != nil {
		return "", err
	}
	return "processed: " + secret.Engine, nil
}

func (v *AcmeVault) getSecretEngine(ctx context.Context, client *vault.Client) ([]string, error) {
	engineSlice := []string{}
	engs, err := client.System.MountsListSecretsEngines(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("")
		return nil, err
	}
	for eng := range engs.Data {
		engineSlice = append(engineSlice, eng)
	}
	return engineSlice, nil
}

func (v *AcmeVault) writeSecret(ctx context.Context, client *vault.Client, path string, data map[string]interface{}) error {
	log.Debug().Msgf("writing secrets to: %v", path)
	if err := retryVaultFunc(func() (*vault.Response[map[string]interface{}], error) {
		return client.Write(ctx, path, map[string]interface{}{"data": data})
	}, 3, 2); err != nil {
		return err
	}
	return nil
}

func ReadSecret(ctx context.Context, c *vault.Client, path string, secret string) (string, error) {
	response, err := c.Read(ctx, path)
	if err != nil {
		return "", err
	}
	data, ok := response.Data["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("response data for %s secret not ok", secret)
	}
	return data[secret].(string), nil
}

func CreateDataInVault(ctx context.Context, client *vault.Client, v Vaulter, s []*Secret) error {
	var wg sync.WaitGroup

	engines, err := v.getSecretEngine(ctx, client)
	if err != nil {
		return err
	}

	for _, secret := range s {
		if !slices.Contains(engines, secret.Engine+"/") {
			wg.Add(1)
			go func(secret *Secret) {
				defer wg.Done()
				_, err := v.createEngines(ctx, client, secret)
				if err != nil {
					log.Warn().Err(err).Msg("create engines error")
				}
			}(secret)
		}
		for _, kv := range secret.KV {
			wg.Add(1)
			go func(secret *Secret, kv struct {
				Data map[string]interface{}
				Path string
			}) {
				defer wg.Done()
				path := fmt.Sprintf("%v/data/%v", secret.Engine, kv.Path)
				if err := v.writeSecret(ctx, client, path, kv.Data); err != nil {
					log.Warn().Err(err)
				} else {
					log.Info().Msgf("secrets in: %q written", path)
				}
			}(secret, kv)
		}
	}

	wg.Wait()

	return nil
}

func (v *AcmeVault) hydrateNewSecretsStruct(ctx context.Context, c *vault.Client, s []*Secret, secretMap map[string]secretMap) error {
	for _, secret := range s {
		for _, kv := range secret.KV {
			for key := range kv.Data {
				sm := secretMap[key]
				if sm.path != "" {
					value, err := ReadSecret(ctx, c, sm.path, sm.secret)
					if err != nil {
						return err
					}
					kv.Data[key] = value
				}
			}
		}
	}
	return nil
}

func retryVaultFunc(fn vaultFunc, attempts, sleep int) error {
	var err error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Warn().Msgf("retrying after error: %v", err)
			time.Sleep(time.Duration(sleep) * time.Second)
			sleep *= 2
		}
		_, err = fn()

		if err == nil {
			return nil
		}
		if !vault.IsErrorStatus(err, http.StatusBadRequest) {
			return err
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
