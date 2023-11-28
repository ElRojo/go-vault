package api

import (
	"context"
	"encoding/json"
	"fmt"
	"go-vault/controllers/vault"
	"go-vault/internal/utility"
	"net/http"

	vc "github.com/hashicorp/vault-client-go"
)

func (s *APIServer) handleVault(w http.ResponseWriter, r *http.Request) error {
	setCORS(w, s.CORS)
	token := r.Header.Get("Api-Key")
	URL := r.Header.Get("Vault-Url")

	if err := validateHeaders(token, URL); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	ctx, client, err := vault.InitVaultClient(token, URL)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	switch r.URL.Path {
	case "/vault/init":
		switch r.Method {
		case http.MethodPost:
			return s.handleInitVault(w, r, ctx, client)
		default:
			return WriteJSON(w, http.StatusMethodNotAllowed, APIError{Error: r.Method + " not allowed"})
		}

	case "/vault/secret":
		switch r.Method {
		case http.MethodPost:
			return s.handleCreateSecret(w, r, ctx, client)
		case http.MethodGet:
			return s.handleGetSecret(w, r, ctx, client)
		default:
			return WriteJSON(w, http.StatusMethodNotAllowed, APIError{Error: r.Method + " not allowed"})
		}
	}
	return WriteJSON(w, http.StatusNotFound, APIError{Error: "404 path not found " + r.URL.Path})
}

func (s *APIServer) handleCreateSecret(w http.ResponseWriter, r *http.Request, ctx context.Context, client *vc.Client) error {
	var (
		req           = &VaultSecret{}
		vaultInstance = &vault.AcmeVault{}
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: "invalid JSON format"})
	}

	if err := utility.ValidateRequestFields(req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	secrets, err := convertSecret(req)
	if err != nil {
		return err
	}

	if err := vault.CreateDataInVault(ctx, client, vaultInstance, secrets); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, req)
}

func (s *APIServer) handleGetSecret(w http.ResponseWriter, r *http.Request, ctx context.Context, client *vc.Client) error {
	req := &VaultRead{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: "invalid JSON format"})
	}

	if err := utility.ValidateRequestFields(req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	var path = fmt.Sprintf("%v/data/%v", req.Engine, req.Path)

	secret, err := vault.ReadSecret(ctx, client, path, req.Key)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, APIResponse{Success: map[string]string{req.Key: secret}})
}

func (s *APIServer) handleInitVault(w http.ResponseWriter, r *http.Request, ctx context.Context, client *vc.Client) error {
	var (
		req           = &VaultInit{}
		vaultInstance = &vault.AcmeVault{}
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: "invalid JSON format"})
	}

	if err := utility.ValidateRequestFields(req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	var secrets = initSecrets(*req.UseLegacy)

	v, err := vault.InitVault(ctx, client, vaultInstance, secrets, vault.VaultConfig{
		Copy:   *req.CopyLegacy,
		Legacy: *req.UseLegacy,
	})
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, v)
}

func initSecrets(legacy bool) []*vault.Secret {
	if legacy {
		return vault.InitLegacySecrets()
	}
	return vault.InitNewSecrets()
}

func convertSecret(s *VaultSecret) ([]*vault.Secret, error) {
	var secretSlice []*vault.Secret

	for _, secret := range s.Secret {
		newObj := &vault.Secret{
			Engine: secret.Engine,
			KV: make([]struct {
				Data map[string]interface{}
				Path string
			}, len(secret.KV)),
		}
		for i, kv := range secret.KV {
			newObj.KV[i] = struct {
				Data map[string]interface{}
				Path string
			}{
				Data: kv.Data,
				Path: kv.Path,
			}
		}
		secretSlice = append(secretSlice, newObj)
	}
	return secretSlice, nil
}

func validateHeaders(token string, URL string) error {
	switch {
	case len(token) == 0 && len(URL) == 0:
		return fmt.Errorf("Api-Key, Vault-Url missing from headers")
	case len(token) == 0:
		return fmt.Errorf("Api-Key missing from headers")
	case len(URL) == 0:
		return fmt.Errorf("Vault-Url missing from headers")
	}
	return nil
}
