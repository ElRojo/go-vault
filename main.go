package main

import (
	"fmt"
)

func main() {
	var s []*secret
	v := &acmeVault{}
	if !getEnv().vaultConfig {
		s = newSecrets
	} else {
		s = legacySecrets
	}
	c := Config{
		Secrets: s,
	}
	cp := getEnv().vaultCopy
	fmt.Println(runVault(v, c, cp))
}
