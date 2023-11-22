package api

import (
	"net/http"
)

// common

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIServer struct {
	ListenerAddress string
	CORS            string
}

type APIError struct {
	Error string
}

type APIResponse struct {
	Success any
}

// vault
type VaultAuth struct {
	URL        string `json:"vaultUrl" validate:"required"`
	VaultToken string `json:"vaultToken" validate:"required"`
}

type VaultRequest struct {
	Auth       VaultAuth `json:"authentication" validate:"required"`
	CopyLegacy *bool     `json:"copyLegacy" validate:"required"`
	UseLegacy  *bool     `json:"useLegacy" validate:"required"`
}

type KV struct {
	Data map[string]interface{} `json:"data" validate:"required"`
	Path string                 `json:"path" validate:"required"`
}

type Secret struct {
	Engine string `json:"engine"`
	Keys   []KV   `json:"kv"`
}

type VaultSecret struct {
	Auth   VaultAuth `json:"authentication" validate:"required"`
	Secret []Secret  `json:"secret"`
}

type VaultRead struct {
	Auth   VaultAuth `json:"authentication" validate:"required"`
	Engine string    `json:"engine"`
	Path   string    `json:"path" validate:"required"`
	Key    string    `json:"key" validate:"required"`
}
