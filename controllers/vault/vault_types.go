package vault

type AcmeVault struct{}

type secretMap struct {
	secret string
	path   string
}

type VaultConfig struct {
	Copy   bool
	Legacy bool
	Token  string
	URL    string
}

type Secret struct {
	Engine string
	KV     []struct {
		Data map[string]interface{}
		Path string
	}
}
