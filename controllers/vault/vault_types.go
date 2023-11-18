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

type kv struct {
	data map[string]interface{}
	path string
}

type secret struct {
	engine string
	keys   []kv
}
