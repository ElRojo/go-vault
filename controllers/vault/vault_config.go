package vault

import (
	"go-vault/internal/utility"
)

func initNewSecrets() []*secret {
	var key string = utility.GenerateUUID()

	var newSecrets = []*secret{
		{
			engine: "my-folder",
			keys: []kv{
				{
					data: map[string]interface{}{
						"MDM_API_KEY":            "",
						"RINGCENTRAL_SECRET":     "",
						"CLL_API_KEY":            "",
						"MSCOTT_PRIVATE_KEY":     "",
						"USER_vault_ENCODED_KEY": "",
						"ECHO_vault_KEY":         "",
						"PAM_RSA":                "",
						"APPLE_DEV_PW":           "",
						"AIK_API_KEY":            key,
						"SUPER_SECRET_KEY":       "",
						"HAM_SECRET_API_KEY":     "",
						"MSC_PRIVATE_KEY":        "",
						"SECRET_PATH":            "",
						"VB_SECRET_API_KEY":      "",
						"VAR_SECRET":             "",
						"TEST_SECRET_KEY1":       "",
						"PARC_GPS_KEY":           "",
						"MSSQL_CONN_STRING":      utility.GeneratePassword(25),
						"OAL_USERNAME":           utility.GeneratePassword(25),
						"LOR_EM_IPSUM":           utility.GeneratePassword(25),
						"AWS_CREDENTIALS":        "",
						"UOL_TAL_URL":            "",
						"WPI_WAP_KEY":            "",
						"CHAMP_014":              "",
						"DRIC_FL_STRING":         "",
					},
					path: "sec-engine",
				},
			},
		},
	}
	return newSecrets
}

func initLegacySecrets() []*secret {

	var key = utility.GenerateUUID()
	var legacySecrets = []*secret{
		{
			engine: "test-engine",
			keys: []kv{
				{
					data: map[string]interface{}{
						"username_for_oal": "",
						"secret_key":       key,
					},
					path: "colpuls",
				},
				{
					data: map[string]interface{}{
						"apple_dev_pw": "",
						"aik_key":      key,
					},
					path: "apple/developer",
				},
				{
					data: map[string]interface{}{
						"lorem_key": key,
					},
					path: "IPSU",
				},
				{
					data: map[string]interface{}{
						"cll_api_key": "",
						"mscott_key":  "",
					},
					path: "office",
				},
			},
		},
		{
			engine: "test-engine-1",
			keys: []kv{
				{
					data: map[string]interface{}{
						"var_sec": "",
					},
					path: "vars",
				},
				{
					data: map[string]interface{}{
						"private_key": "",
					},
					path: "IPSU",
				},
				{
					data: map[string]interface{}{
						"ham_sec_api_key": "",
					},
					path: "local/accts",
				},
				{
					data: map[string]interface{}{
						"api_key": "",
					},
					path: "jil",
				},
				{
					data: map[string]interface{}{
						"path": "",
					},
					path: "lwe",
				},
				{
					data: map[string]interface{}{
						"api_key": "",
					},
					path: "qet",
				},
				{
					data: map[string]interface{}{
						"api_key": "",
					},
					path: "pwo",
				},
				{
					data: map[string]interface{}{
						"api_key": "",
					},
					path: "office",
				},
			},
		},
		{
			engine: "test-engine-2",
			keys: []kv{
				{
					data: map[string]interface{}{
						"conn_string": "",
						"tal_url":     "",
					},
					path: "test-engine-2",
				},
			},
		},
		{
			engine: "postgres",
			keys: []kv{
				{
					data: map[string]interface{}{
						"conn_string": "",
						"url":         "",
					},
					path: "vars",
				},
			},
		},
		{
			engine: "prc_gps",
			keys: []kv{
				{
					data: map[string]interface{}{
						"gps_key": "",
					},
					path: "office",
				},
			},
		},
		{
			engine: "aws",
			keys: []kv{
				{data: map[string]interface{}{
					"aws-creds": "",
				},
					path: "creds",
				},
			},
		},
		{
			engine: "supersecrets",
			keys: []kv{
				{
					data: map[string]interface{}{
						"MSSQL_CONN_STRING":  "",
						"jamf_key":           utility.GeneratePassword(35),
						"ringcentral_secret": "",
						"encoded_vault_key":  "",
						"test_sec_key":       "",
						"MSC_PRIVATE_KEY":    "",
						"wpi_wap_key":        "",
						"ch_014_key":         "",
						"drr-fl-string":      utility.GeneratePassword(25),
					},
					path: "supersecrets",
				},
			},
		},
	}
	return legacySecrets
}

func initSecretMap() map[string]secretMap {

	var secretsMap = map[string]secretMap{
		"MDM_API_KEY": {
			path:   "/supersecrets/data/supersecrets",
			secret: "jamf_key",
		},
		"RINGCENTRAL_SECRET": {
			path:   "/supersecrets/data/supersecrets",
			secret: "ringcentral_secret",
		},
		"CLL_API_KEY": {
			path:   "/test-engine/data/office",
			secret: "cll_api_key",
		},
		"MSCOTT_PRIVATE_KEY": {
			path:   "/test-engine/data/office",
			secret: "mscott_key",
		},
		"USER_vault_ENCODED_KEY": {
			path:   "/supersecrets/data/supersecrets",
			secret: "encoded_vault_key",
		},
		"APPLE_DEV_PW": {
			path:   "/test-engine/data/apple/developer",
			secret: "apple_dev_pw",
		},
		"AIK_API_KEY": {
			path:   "/test-engine/data/apple/developer",
			secret: "aik_key",
		},
		"SUPER_SECRET_KEY": {
			path:   "/test-engine/data/colpuls",
			secret: "secret_key",
		},
		"HAM_SECRET_API_KEY": {
			path:   "/test-engine-1/data/local/accts",
			secret: "ham_sec_api_key",
		},
		"MSC_PRIVATE_KEY": {
			path:   "/supersecrets/data/supersecrets",
			secret: "MSC_PRIVATE_KEY",
		},
		"SECRET_PATH": {
			path:   "/test-engine-1/data/lwe",
			secret: "path",
		},
		"VAR_SECRET": {
			path:   "/test-engine-1/data/vars",
			secret: "var_sec",
		},
		"TEST_SECRET_KEY1": {
			path:   "/supersecrets/data/supersecrets",
			secret: "test_sec_key",
		},
		"PARC_GPS_KEY": {
			path:   "/prc_gps/data/office",
			secret: "gps_key",
		},
		"MSSQL_CONN_STRING": {
			path:   "/supersecrets/data/supersecrets",
			secret: "MSSQL_CONN_STRING",
		},
		"OAL_USERNAME": {
			path:   "/test-engine/data/colpuls",
			secret: "username_for_oal",
		},
		"LOR_EM_IPSUM": {
			path:   "/test-engine/data/IPSU",
			secret: "lorem_key",
		},
		"AWS_CREDENTIALS": {
			path:   "/aws/data/creds",
			secret: "aws-creds",
		},
		"UOL_TAL_URL": {
			path:   "/test-engine-2/data/test-engine-2",
			secret: "tal_url",
		},
		"WPI_WAP_KEY": {
			path:   "/supersecrets/data/supersecrets",
			secret: "wpi_wap_key",
		},
		"CHAMP_014": {
			path:   "/supersecrets/data/supersecrets",
			secret: "ch_014_key",
		},
		"DRIC_FL_STRING": {
			path:   "/supersecrets/data/supersecrets",
			secret: "drr-fl-string",
		},
	}
	return secretsMap
}
