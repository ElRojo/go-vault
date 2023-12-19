package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-vault/controllers/vault"
	"go-vault/internal/utility"
	"net/http"

	vc "github.com/hashicorp/vault-client-go"
)

func (s *APIServer) handleVault(w http.ResponseWriter, r *http.Request) error {
	setCORS(w, s.CORS)
	token := r.Header.Get("x-api-key")
	URL := r.Header.Get("x-vault-url")

	if err := validateVaultHeaders(token, URL); err != nil {
		return err
	}

	ctx, client, err := vault.InitVaultClient(token, URL)
	if err != nil {
		return WriteJSON(w, 500, APIError{Error: "initialization failed"})
	}

	switch r.Method {
	case http.MethodGet:
		return s.handleVaultGETRequests(w, r, ctx, client)
	case http.MethodPost:
		return s.handleVaultPOSTRequests(w, r, ctx, client)
	}
	return WriteJSON(w, http.StatusMethodNotAllowed, APIError{Error: r.Method + " not allowed"})
}

func (s *APIServer) handleVaultGETRequests(w http.ResponseWriter, r *http.Request, ctx context.Context, client *vc.Client) error {
	switch r.URL.Path {
	case "/vault/secret":
		return s.handleGetSecret(w, r, ctx, client)
	default:
		return WriteJSON(w, http.StatusNotFound, APIError{Error: "404 path not found " + r.URL.Path})
	}
}

func (s *APIServer) handleVaultPOSTRequests(w http.ResponseWriter, r *http.Request, ctx context.Context, client *vc.Client) error {
	switch r.URL.Path {
	case "/vault/secret":
		return s.handleCreateSecret(w, r, ctx, client)
	case "/vault/init":
		return s.handleInitVault(w, r, ctx, client)
	}
	return WriteJSON(w, http.StatusNotFound, APIError{Error: "404 path not found " + r.URL.Path})
}

func (s *APIServer) handleCreateSecret(w http.ResponseWriter, r *http.Request, ctx context.Context, client *vc.Client) error {
	var (
		req           VaultSecret
		vaultInstance vault.AcmeVault
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid json format")
	}

	if err := utility.ValidateRequestFields(req); err != nil {
		return err
	}

	secrets, err := convertVaultSecret(req.Secret)
	if err != nil {
		return err
	}

	if err := vault.CreateDataInVault(ctx, client, &vaultInstance, secrets); err != nil {
		return errors.Unwrap(err)
	}

	return WriteJSON(w, http.StatusOK, req)
}

func (s *APIServer) handleGetSecret(w http.ResponseWriter, r *http.Request, ctx context.Context, client *vc.Client) error {
	var req VaultRead

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid json format")
	}

	if err := utility.ValidateRequestFields(req); err != nil {
		return err
	}

	var path = fmt.Sprintf("%v/data/%v", req.Engine, req.Path)

	secret, err := vault.ReadSecret(ctx, client, path, req.Key)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, APIResponse[map[string]string]{Success: map[string]string{req.Key: secret}})
}

func (s *APIServer) handleInitVault(w http.ResponseWriter, r *http.Request, ctx context.Context, client *vc.Client) error {

	var (
		req           VaultInit
		vaultInstance vault.AcmeVault
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid json format")
	}

	if err := utility.ValidateRequestFields(req); err != nil {
		return err
	}

	var secrets = initSecrets(*req.UseLegacy)

	v, err := vault.InitVault(ctx, client, &vaultInstance, secrets, vault.VaultConfig{
		Copy:   *req.CopyLegacy,
		Legacy: *req.UseLegacy,
	})
	if err != nil {
		return errors.Unwrap(err)
	}

	return WriteJSON(w, http.StatusOK, map[string]*VaultInit{v: &req})
}

func initSecrets(legacy bool) []*vault.Secret {
	if legacy {
		return vault.InitLegacySecrets()
	}
	return vault.InitNewSecrets()
}

func convertVaultSecret(s []Secret) ([]*vault.Secret, error) {
	var secretSlice []*vault.Secret

	for _, secret := range s {
		createdSecret := &vault.Secret{
			Engine: secret.Engine,
			KV: make([]struct {
				Data map[string]interface{}
				Path string
			}, len(secret.KV)),
		}
		for i, kv := range secret.KV {
			createdSecret.KV[i] = struct {
				Data map[string]interface{}
				Path string
			}{
				Data: kv.Data,
				Path: kv.Path,
			}
		}
		secretSlice = append(secretSlice, createdSecret)
	}
	return secretSlice, nil
}

func validateVaultHeaders(token string, URL string) error {
	switch {
	case len(token) == 0 && len(URL) == 0:
		return fmt.Errorf("x-api-key, vault-url missing from headers")
	case len(token) == 0:
		return fmt.Errorf("x-api-key missing from headers")
	case len(URL) == 0:
		return fmt.Errorf("x-vault-url missing from headers")
	}
	return nil
}
