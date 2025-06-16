package api

const (
	minimumNodes            = 3
	tofu_secrets_dir        = "./.tofu_secrets"
	unseal_keys_secrets_dir = tofu_secrets_dir + "/.unseal_keys"
	root_token_file         = tofu_secrets_dir + "/.root_token"
	qrcodes_prefix          = ""
	// qrcodes_prefix          = "https://kairos/qrcode/"
)

// Config is a combo of the flags passed to the cli and the configuration file (if used).
type NavySealsConfig struct {
	Version                   bool `short:"v" long:"version" description:"display the version of navy-seals and exit"`
	Debug                     bool `short:"D" long:"debug" description:"enable debugging (extra logging)"`
	DefaultUnsealkeysNb       int  `env:"DEFAULT_UNSEAL_KEYS_NB"            long:"default-unseal-keys-nb" description:"The default number of shared keys to init the openbao vault with (integer, maximum 256)" yaml:"default_unseal_keys_nb"`
	DefaultUnsealKeysTreshold int  `env:"DEFAULT_UNSEAL_KEYS_TRESHOLD"            long:"default-unseal-keys-treshold" description:"The default treshold of shared keys to init the openbao vault with (integer, maximum 256)" yaml:"default_unseal_keys_treshold"`
	Seal                      bool `short:"s" long:"seal" description:"seal the OpenBAO vault and exit"`
	UnSeal                    bool `short:"u" long:"unseal" description:"Unseal the OpenBAO vault from th QR codes found inside the 'tofu_secrets_dir' Folder, and exit"`
	Status                    bool `short:"t" long:"status" description:"Show the Status of the OpenBAO vault, and exit"`
	// ConfigPath         string   `env:"CONFIG_PATH" short:"c" long:"config" description:"path to configuration file" value-name:"PATH"`
	BindAddress string `env:"BIND_ADDRESS"            long:"bind" short:"b" description:"bind address" yaml:"bind"`
	Port        int    `env:"PORT"          short:"p"  long:"port" description:"port number (integer, maximum 65535)" yaml:"port"`
	Log         struct {
		Path   string `env:"LOG_PATH"  long:"path"    description:"path to log output to" value-name:"PATH"`
		Quiet  bool   `env:"LOG_QUIET" long:"quiet"   description:"disable logging to stdout (also: see levels)"`
		Level  string `env:"LOG_LEVEL" long:"level"   default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal"  description:"logging level"`
		JSON   bool   `env:"LOG_JSON"   long:"json"   description:"output logs in JSON format"`
		Pretty bool   `env:"LOG_PRETTY" long:"pretty" description:"output logs in a pretty colored format (cannot be easily parsed)"`
	} `group:"Logging Options" namespace:"log"`

	Environment string `env:"ENVIRONMENT" long:"environment" description:"environment this cluster relates to (for logging)" yaml:"environment"`

	TLSSkipVerify bool `env:"TLS_SKIP_VERIFY"   long:"tls-skip-verify"      description:"disables tls certificate validation: DO NOT DO THIS" yaml:"tls_skip_verify"`
}

var (
	ApiConfig = &NavySealsConfig{}

	// logger log.Interface
)
