package main

type kv struct {
	data map[string]interface{}
	path string
}

type secret struct {
	engine string
	keys   []kv
}

var (
	key string = generateUuid()

	newSecrets = []*secret{
		{
			engine: "my-folder",
			keys: []kv{
				{
					data: map[string]interface{}{
						"MDM_API_KEY":              "",
						"RINGCENTRAL_SECRET":       "",
						"CLL_API_KEY":              "",
						"MSCOTT_PRIVATE_KEY":       "",
						"USER_SERVICE_ENCODED_KEY": "",
						"ECHO_SERVICE_KEY":         "",
						"PAM_RSA":                  "",
						"APPLE_DEV_PW":             "",
						"AIK_API_KEY":              "",
						"SUPER_SECRET_KEY":         "",
						"HAM_SECRET_API_KEY":       "",
						"MSC_PRIVATE_KEY":          "",
						"SECRET_PATH":              "",
						"VB_SECRET_API_KEY":        "",
						"VAR_SECRET":               "",
						"TEST_SECRET_KEY1":         "",
						"PARC_GPS_KEY":             "",
						"MSSQL_CONN_STRING":        "",
						"OAL_USERNAME":             "",
						"LOR_EM_IPSUM":             "",
						"AWS_CREDENTIALS":          "",
						"UOL_TAL_URL":              "",
						"WPI_WAP_KEY":              "",
						"CHAMP_014":                "",
						"DRIC_FL_STRING":           "",
					},
					path: "sec-engine",
				},
			},
		},
	}

	legacySecrets = []*secret{
		{
			engine: "test-engine",
			keys: []kv{
				{
					data: map[string]interface{}{
						"username_for_oal": "",
						"secrt_key":        key,
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
						"MSSQL_CONN_STRING":   "",
						"jamf_key":            generatePassword(35),
						"ringcentral_secret":  "",
						"encoded_service_key": "",
						"test_sec_key":        "",
						"MSC_PRIVATE_KEY":     "",
						"wpi_wap_key":         "",
						"ch_014_key":          "",
						"drr-fl-string":       generatePassword(25),
					},
					path: "supersecrets",
				},
			},
		},
	}
)
