# go-vault-example

- [config.go](#config)
- [envs.go](#envs)
- [merge_conf.go](#merge-config)
- [pwgen.go](#pwgen)
- [vault_test.go](#vault-test)
- [vault.go](#vault)

---

This poject was used previously as a way to easily stand up multiple Vault instances with the same footprint. I also migrated from an older vault architecture with duplicate entries to a simpler flat-file approach.
I utilized the [vault-client-go](https://pkg.go.dev/github.com/hashicorp/vault-client-go@v0.4.1) package.

This was necesary due to a lack of access to Terraform.

---

## Config

Vault is based on CRUD operations and as such has decided that all data needs to be created (or updated) at once by passing in a map of `string:string` (more precisely, `map[string]interface{}`).

I wanted to package as much information together as I could so I bundled all of the data into a `kv` struct which holds the arrays of k/v pairs themselves and the path inside the engine where these k/v pairs should live.

Further, I needed to iterate over `engines` (folders in Vault-speak) and place secrets in different paths inside the same engine. Thus was born the `secret` struct.

```go
type kv struct {
	data map[string]interface{}
	path string
}

type secret struct {
	engine string
	keys   []kv
}

var sampleSecret = []*secret{
	{
		engine: "my-engine",
		keys: []kv{
			{
				data: map[string]interface{}{
					"myKey":  "myValue",
					"myKey2": "myValue2",
				},
				path: "my-path-1",
			},
		},
	},
}
```

## Envs

A fair question is "why don't you just use the built-in `os.Getenv()`?" And frankly, I don't have a technical reason. The main reason I used the [envs](github.com/caarlos0/env/v8) package was to make the environment variables declarative and to have some built-in defaults (instead of if statements or switch cases). I want other engineers to be able to jump in and see what environment variables do what and feel comfortable adding or modifying variables if necessary. Also, the `env` package does work with types which comes in handy for the two boolean variables. So, there we go; there's a technical reason.

## Merge Config

Merge config came along a while after I had built out the project. My Vault instances had many duplicates and no real organization. The secret names were also confusing/unclear and this caused even more duplicates in vault. I decided to migrate to a flattened structure. With this, I wanted to keep the old structure in-tact in case any old systems were using them and I also didn't want to have to copy-paste information by-hand.

To handle this, I created the `merge_conf.go` file and the `hydrateNewSecretsStruct()` function. This would take the `newSecrets` struct and fill in the values from the vault instance and then push the hydrated `newSecrets` into vault, saving hours of work. The structure on this one is pretty simple as it adds an extra layer to the `secrets{}`:

```go
type secretMap struct {
	secret string
	path   string
}
```
This function takes the path for the given secret, and then searches the `newSecrets{}` for a matching key. When a matching key is found, it places the secret gathered in as the value.

## Pwgen

This was built with the help of a medium article and uses Google's UUID generator. This can be helpful when setting up a fresh instance and setting some random passwords.

## Vault Test

I wanted to make sure I could test many of these functions without needing to make any real API calls. I decided to mock many of the calls and built an interface to utilize dependency injection.

## Vault

Given the information above I hope that the vault.go file is self explanatory. These are all the functions necessary to authenticate with vault and then read / write secrets as necessary. I decided to export the `ReadSecrets()` and `InitVaultClient()` functions in the case that this project were to expand into mulitple packages.