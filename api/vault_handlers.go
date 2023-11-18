package api

import (
	"encoding/json"
	"go-vault/controllers/vault"
	"go-vault/internal/utility"
	"net/http"
)

func (s *APIServer) handleVault(w http.ResponseWriter, r *http.Request) error {
	enableCors(w, s.CORS)
	if r.Method == "POST" {
		return s.handleRunVault(w, r)
	}
	return WriteJSON(w, http.StatusMethodNotAllowed, APIError{Error: r.Method + "not allowed"})
}

func (s *APIServer) handleRunVault(w http.ResponseWriter, r *http.Request) error {
	vaultReq := &VaultRequest{}
	if err := json.NewDecoder(r.Body).Decode(&vaultReq); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: "invalid JSON format"})
	}
	if err := utility.ValidateRequestFields(vaultReq); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}
	v, err := vault.RunVault(&vault.AcmeVault{}, vault.VaultConfig{
		Copy:   *vaultReq.CopyLegacy,
		Legacy: *vaultReq.UseLegacy,
		Token:  vaultReq.VaultToken,
		URL:    vaultReq.URL,
	})
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, v)
}
