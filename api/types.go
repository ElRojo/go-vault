package api

import "net/http"

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIServer struct {
	ListenerAddress string
	CORS            string
}

type APIError struct {
	Error string
}

type VaultRequest struct {
	CopyLegacy *bool  `json:"copyLegacy"validate:"required"`
	URL        string `json:"vaultUrl"validate:"required"`
	UseLegacy  *bool  `json:"useLegacy"validate:"required"`
	VaultToken string `json:"vaultToken"validate:"required"`
}
