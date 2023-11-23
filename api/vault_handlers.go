package api

import (
	"encoding/json"
	"fmt"
	"go-vault/controllers/vault"
	"go-vault/internal/utility"
	"net/http"
)

func (s *APIServer) handleVault(w http.ResponseWriter, r *http.Request) error {
	enableCors(w, s.CORS)
	switch r.Method {
	case "POST":
		return s.handleInitVault(w, r)
	}
	return WriteJSON(w, http.StatusMethodNotAllowed, APIError{Error: r.Method + "not allowed"})
}

// func (s *APIServer) handleEngine(w http.ResponseWriter, r *http.Request) error {
// 	enableCors(w, s.CORS)
// 	switch r.Method {
// 	case "POST":
// 		return s.handleCreateEngine(w, r)
// 	}

// 	return WriteJSON(w, http.StatusMethodNotAllowed, APIError{Error: r.Method + "not allowed"})
// }

func (s *APIServer) handleSecret(w http.ResponseWriter, r *http.Request) error {
	enableCors(w, s.CORS)
	switch r.Method {
	case "POST":
		return s.handleCreateSecret(w, r)
	case "GET":
		return s.handleGetSecret(w, r)
	}
	return WriteJSON(w, http.StatusMethodNotAllowed, APIError{Error: r.Method + " not allowed"})
}

// func (s *APIServer) handleCreateEngine(w http.ResponseWriter, r *http.Request) error {
// 	var (
// 		vaultReq      = &VaultSecret{}
// 		vaultInstance = &vault.AcmeVault{}
// 	)
// 	if err := json.NewDecoder(r.Body).Decode(&vaultReq); err != nil {
// 		return WriteJSON(w, http.StatusBadRequest, APIError{Error: "invalid JSON format"})
// 	}

// 	if err := utility.ValidateRequestFields(vaultReq); err != nil {
// 		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
// 	}
// 	return WriteJSON(w, http.StatusOK, "ok")
// }

func (s *APIServer) handleCreateSecret(w http.ResponseWriter, r *http.Request) error {
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

	ctx, client, err := vault.InitVaultClient(req.Auth.VaultToken, req.Auth.URL)
	if err != nil {
		return err
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

func (s *APIServer) handleGetSecret(w http.ResponseWriter, r *http.Request) error {
	var req = &VaultRead{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: "invalid JSON format"})
	}

	if err := utility.ValidateRequestFields(req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	ctx, client, err := vault.InitVaultClient(req.Auth.VaultToken, req.Auth.URL)
	if err != nil {
		return err
	}

	var path = fmt.Sprintf("%v/data/%v", req.Engine, req.Path)

	secret, err := vault.ReadSecret(ctx, client, path, req.Key)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, APIResponse{Success: map[string]string{req.Key: secret}})
}

func (s *APIServer) handleInitVault(w http.ResponseWriter, r *http.Request) error {
	var (
		req           = &VaultRequest{}
		vaultInstance = &vault.AcmeVault{}
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: "invalid JSON format"})
	}

	if err := utility.ValidateRequestFields(req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	var secrets = initSecrets(*req.UseLegacy)

	ctx, client, err := vault.InitVaultClient(req.Auth.VaultToken, req.Auth.URL)
	if err != nil {
		return err
	}

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
	var vaultSecrets []*vault.Secret

	for _, secret := range s.Secret {
		vaultSecret := &vault.Secret{
			Engine: secret.Engine,
			Keys:   make([]vault.KV, len(secret.Keys)),
		}

		for i, kv := range secret.Keys {
			vaultSecret.Keys[i] = vault.KV{
				Data: kv.Data,
				Path: kv.Path,
			}
		}
		vaultSecrets = append(vaultSecrets, vaultSecret)
	}
	return vaultSecrets, nil
}
