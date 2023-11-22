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

type KV struct {
	Data map[string]interface{}
	Path string
}

type Secret struct {
	Engine string
	Keys   []KV
}
