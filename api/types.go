package api

import (
	"net/http"
)

// common

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIServer struct {
	CORS string
	srv  *http.Server
}

type APIError struct {
	Error string
}

type APIResponse[T any] struct {
	Success T
}

// vault

type VaultInit struct {
	CopyLegacy *bool `json:"copyLegacy" validate:"required"`
	UseLegacy  *bool `json:"useLegacy" validate:"required"`
}

type Secret struct {
	Engine string `json:"engine"`
	KV     []struct {
		Data map[string]interface{} `json:"data" validate:"required"`
		Path string                 `json:"path" validate:"required"`
	} `json:"kv"`
}

type VaultSecret struct {
	Secret []Secret `json:"secret"`
}

type VaultRead struct {
	Engine string `json:"engine" validate:"required"`
	Path   string `json:"path" validate:"required"`
	Key    string `json:"key" validate:"required"`
}
