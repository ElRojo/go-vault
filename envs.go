package main

import (
	"fmt"

	"github.com/caarlos0/env/v8"
)

type config struct {
	env         string `env:"ENV" envDefault:"test"`
	vaultConfig bool   `env:"USE_LEGACY_BOOL" envDefault:"false"`
	vaultCopy   bool   `env:"VAULT_COPY" envDefault:"false"`
	vaultNs     string `env:"VAULT_NAMESPACE" envDefault:"test"`
	vaultTkn    string `env:"VAULT_TOKEN" envDefault:"test"`
	vaultUrl    string `env:"VAULT_URL" envDefault:"test"`
}

func getEnv() config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return cfg
}
