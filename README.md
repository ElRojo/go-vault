# go-vault
# go-vault

- [Quickstart](#quickstart)
- [config.go](#config)
- [vault_config.go](#vault-config)
  - [secretsMap](#secrets-map)
- [pwgen.go](#pwgen)
- [vault_test.go](#vault-test)
- [vault.go](#vault)
- [Utility](#utility)

---

This project was used previously as a way to easily stand up multiple Vault instances with the same footprint. I also migrated from an older vault architecture with duplicate entries to a simpler flat-file approach.
I utilized the [vault-client-go](https://pkg.go.dev/github.com/hashicorp/vault-client-go@v0.4.1) package.

This was necessary due to a lack of access to Terraform.

# Quickstart

Before you begin, make sure you have [Docker](https://www.docker.com/get-started/) installed and running.

### Demoing the Application

1. Clone the repository and then `cd` into the project directory.

2. Run the following commands in your terminal inside the project directory:

```bash
docker-compose up
```

3. After the containers have started, access the HashiCorp Vault front-end at [http://localhost:8200](http://localhost:8200) and the API is served on [http://localhost:4269](http://localhost:4269)

4. Log in using **Method: Token** with the following credentials: `dev-only-token`

5. Send a POST request to [http://localhost:4269/vault](http://localhost:4269/vault) with the following JSON object to test. See the [Request Documentation](#request-documentation) for finer details.

```JSON
{
  "copyLegacy": false,
  "useLegacy": false,
  "vaultToken": "dev-only-token",
  "vaultUrl": "http://vault:8200"
}
```

6. Refresh your browser to view the updated secrets engine

7. Exit and kill the containers when done with `CTRL+C`

## Request Documentation

<details>
 <summary><b> POST http://localhost:4269/vault </b></summary>

### Vault Request Object

| property     | type   | value example                   | required | purpose                                                                                                                                   |
| ------------ | ------ | ------------------------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| `copyLegacy` | bool   | `true` / `false`                | Y        | If set to `true` and `useLegacy` is set to `false`, this will copy legacy secrets architecture and place them into the flat architecture. |
| `vaultUrl`   | string | `http://hashicorpVaultUrl:8200` | Y        | The URL of the HashiCorp Vault instance.                                                                                                  |
| `useLegacy`  | bool   | `true` / `false`                | Y        | If set to `true`, this builds secrets using the legacy architecture.                                                                      |
| `vaultToken` | string | `dev-only-token`                | Y        | Token to auth with HashiCorp Vault instance.                                                                                              |

### Vault Request Struct

```go
type VaultRequest struct {
	CopyLegacy bool   `json:"copyLegacy"validate:"required"`
	URL        string `json:"vaultUrl"validate:"required"`
	UseLegacy  bool   `json:"useLegacy"validate:"required"`
	VaultToken string `json:"vaultToken"validate:"required"`
}
```

### Example Vault Request Object

```json
{
  "useLegacy": true,
  "copyLegacy": true,
  "vaultToken": "dev-only-token",
  "vaultUrl": "http://vault:8200"
}
```

</details>

## Vault Config

Vault is based on CRUD operations and as such has decided that all data needs to be created (or updated) at once by passing in a map of `string:string` (more precisely, `map[string]interface{}`).

I wanted to package as much information together as I could so I bundled all of the data into a `kv` struct which holds the arrays of k/v pairs themselves and the path inside the engine where these k/v pairs should live.

Further, I needed to iterate over `engines` (folders in Vault-speak) and place secrets in different paths inside the same engine. Thus was born the `secret` struct.

```go
types.go

type kv struct {
	data map[string]interface{}
	path string
}

type secret struct {
	engine string
	keys   []kv
}
```

```go
vault_config.go

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

### Secrets Map

secretMap came along a while after I had built out the project. My Vault instances had many duplicates and no real organization. The secret names were also confusing/unclear and this caused even more duplicates in vault. I decided to migrate to a flattened structure. With this, I wanted to keep the old structure in-tact in case any old systems were using them and I also didn't want to have to copy-paste information by-hand.

To handle this I built the `hydrateNewSecretsStruct()` function. This would take the `newSecrets` struct and fill in the values from the vault instance and then push the hydrated `newSecrets` into vault, saving hours of work. The structure on this one is pretty simple as it adds an extra layer to the `Secrets{}`:

```go
types.go

type secretMap struct {
	secret string
	path   string
}
```

This function takes the path for the given secret, and then searches the `newSecrets{}` for a matching key. When a matching key is found, it places the secret gathered in as the value.

## Vault

Given the information above I hope that the vault.go file is self explanatory. These are all the functions necessary to authenticate with vault and then read / write secrets as necessary.

### Vault Test

I wanted to make sure I could test many of these functions without needing to make any real API calls. I decided to mock many of the calls and built an interface to utilize dependency injection.

## Utility

### Pwgen

This was built with the help of a medium article. It takes a string and creates a slice of runes and randomly selects a rune, converts it back to a string, and adds it to the return value. GenerateUUID uses Google's UUID generator. This can be helpful when setting up a fresh instance and setting some random passwords.

### Validate Request Fields

Using the go-playground/validator package this validates the vault request object. Errors are returned based on incompatibilities or missing properties.
