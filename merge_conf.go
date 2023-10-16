package main

type secretMap struct {
	secret string
	path   string
}

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
	"USER_SERVICE_ENCODED_KEY": {
		path:   "/supersecrets/data/supersecrets",
		secret: "encoded_service_key",
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
		secret: "secrt_key",
	},
	"HAM_SECRET_API_KEY": {
		path:   "/test-engine/data/local/accts",
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
