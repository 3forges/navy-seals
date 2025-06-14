package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	// vapi "github.com/hashicorp/vault/api"

	flags "github.com/jessevdk/go-flags"
	_ "github.com/joho/godotenv/autoload"
	vapi "github.com/openbao/openbao/api/v2"
)

var (
	version = "master"
	commit  = "latest"
	date    = "-"
)

const (
	defaultCheckInterval  = 30 * time.Second
	defaultTimeout        = 15 * time.Second
	configRefreshInterval = 15 * time.Second
	minimumNodes          = 3
)

// Config is a combo of the flags passed to the cli and the configuration file (if used).
type Config struct {
	Version    bool   `short:"v" long:"version" description:"display the version of vault-unseal and exit"`
	Debug      bool   `short:"D" long:"debug" description:"enable debugging (extra logging)"`
	ConfigPath string `env:"CONFIG_PATH" short:"c" long:"config" description:"path to configuration file" value-name:"PATH"`

	Log struct {
		Path   string `env:"LOG_PATH"  long:"path"    description:"path to log output to" value-name:"PATH"`
		Quiet  bool   `env:"LOG_QUIET" long:"quiet"   description:"disable logging to stdout (also: see levels)"`
		Level  string `env:"LOG_LEVEL" long:"level"   default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal"  description:"logging level"`
		JSON   bool   `env:"LOG_JSON"   long:"json"   description:"output logs in JSON format"`
		Pretty bool   `env:"LOG_PRETTY" long:"pretty" description:"output logs in a pretty colored format (cannot be easily parsed)"`
	} `group:"Logging Options" namespace:"log"`

	Environment string `env:"ENVIRONMENT" long:"environment" description:"environment this cluster relates to (for logging)" yaml:"environment"`

	CheckInterval    time.Duration `env:"CHECK_INTERVAL"     long:"check-interval" description:"frequency of sealed checks against nodes" yaml:"check_interval"`
	MaxCheckInterval time.Duration `env:"MAX_CHECK_INTERVAL" long:"max-check-interval" description:"max time that vault-unseal will wait for an unseal check/attempt" yaml:"max_check_interval"`

	AllowSingleNode    bool     `env:"ALLOW_SINGLE_NODE" long:"allow-single-node"    description:"allow vault-unseal to run on a single node" yaml:"allow_single_node" hidden:"true"`
	Nodes              []string `env:"NODES"             long:"nodes" env-delim:","  description:"nodes to connect/provide tokens to (can be provided multiple times & uses comma-separated string for environment variable)" yaml:"vault_nodes"`
	TLSSkipVerify      bool     `env:"TLS_SKIP_VERIFY"   long:"tls-skip-verify"      description:"disables tls certificate validation: DO NOT DO THIS" yaml:"tls_skip_verify"`
	UnsealTokens       []string `env:"UNSEAL_TOKENS"            long:"unseal-tokens" env-delim:"," description:"tokens to provide to nodes (can be provided multiple times & uses comma-separated string for environment variable)" yaml:"unseal_tokens"`
	UnsealkeysNb       int      `env:"UNSEAL_KEYS_NB"            long:"unseal-keys-nb" description:"number of shared keys to init the openbao vault with (integer, maximum 256)" yaml:"unseal_keys_nb"`
	UnsealKeysTreshold int      `env:"UNSEAL_KEYS_TRESHOLD"            long:"unseal-keys-treshold" description:"treshold of shared keys to init the openbao vault with (integer, maximum 256)" yaml:"unseal_keys_treshold"`

	NotifyMaxElapsed time.Duration `env:"NOTIFY_MAX_ELAPSED" long:"notify-max-elapsed" description:"max time before the notification can be queued before it is sent" yaml:"notify_max_elapsed"`
	NotifyQueueDelay time.Duration `env:"NOTIFY_QUEUE_DELAY" long:"notify-queue-delay" description:"time we queue the notification to allow as many notifications to be sent in one go (e.g. if no notification within X time, send all notifications)" yaml:"notify_queue_delay"`

	Email struct {
		Enabled       bool     `env:"EMAIL_ENABLED"         long:"enabled"         description:"enables email support" yaml:"enabled"`
		Hostname      string   `env:"EMAIL_HOSTNAME"        long:"hostname"        description:"hostname of mail server" yaml:"hostname"`
		Port          int      `env:"EMAIL_PORT"            long:"port"            description:"port of mail server" yaml:"port"`
		Username      string   `env:"EMAIL_USERNAME"        long:"username"        description:"username to authenticate to mail server" yaml:"username"`
		Password      string   `env:"EMAIL_PASSWORD"        long:"password"        description:"password to authenticate to mail server" yaml:"password"`
		FromAddr      string   `env:"EMAIL_FROM_ADDR"       long:"from-addr"       description:"address to use as 'From'" yaml:"from_addr"`
		SendAddrs     []string `env:"EMAIL_SEND_ADDRS"      long:"send-addrs"      description:"addresses to send notifications to" yaml:"send_addrs"`
		TLSSkipVerify bool     `env:"EMAIL_TLS_SKIP_VERIFY" long:"tls-skip-verify" description:"skip SMTP TLS certificate validation" yaml:"tls_skip_verify"`
		MandatoryTLS  bool     `env:"EMAIL_MANDATORY_TLS"   long:"mandatory-tls"   description:"require TLS for SMTP connections. Defaults to opportunistic." yaml:"mandatory_tls"`
	} `group:"Email Options" namespace:"email" yaml:"email"`

	lastModifiedCheck time.Time
}

var (
	conf = &Config{CheckInterval: defaultCheckInterval}

	// logger log.Interface
)

func newVault(addr string) (vault *vapi.Client) {
	var err error

	vconfig := vapi.DefaultConfig()
	vconfig.Address = addr
	vconfig.MaxRetries = 0
	vconfig.Timeout = defaultTimeout

	if err = vconfig.ConfigureTLS(&vapi.TLSConfig{Insecure: conf.TLSSkipVerify}); err != nil {
		// logger.WithError(err).Fatal("error initializing tls config")
		fmt.Printf("error initializing tls config %v", err)
	}

	if vault, err = vapi.NewClient(vconfig); err != nil {
		// logger.Fatalf("error creating vault client: %v", err)
		fmt.Printf("Error creating vault client: %v", err)
	}

	return vault
}

func getVaultStatus(addr string) (*vapi.HealthResponse, error) {
	client := newVault(addr)

	toReturn, err := client.Sys().Health()
	if err != nil {
		// logger.Fatalf("error creating vault client: %v", err)
		fmt.Printf("Error getting vault health: %v", err)
	}
	return toReturn, err

}

func sealVault(addr string) {
	client := newVault(addr)

	// toReturn, err := client.Sys().EnableAuth()
	err := client.Sys().Seal()
	if err != nil {
		// logger.Fatalf("error creating vault client: %v", err)
		fmt.Printf("Error sealing vault: %v", err)
	}
	fmt.Printf("The %v vault was successfully sealed!", addr)
}
func showVaultStatus(addr string) {
	response, err := getVaultStatus(addr)
	if err != nil {
		// logger.Fatalf("error creating vault client: %v", err)
		fmt.Printf("Error getting vault health: %v", err)
	} else {
		fmt.Printf("🚑 Here is the vault health [response.Initialized] : %v", response.Initialized)
		fmt.Println()
		fmt.Printf("🚑 Here is the vault health [response.ClusterName] : %v", response.ClusterName)
		fmt.Println()
		fmt.Printf("🚑 Here is the vault health [response.Sealed] : %v", response.Sealed)
		fmt.Println()
		fmt.Printf("🚑 Here is the vault health [response.Version] : %v", response.Version)
		fmt.Println()
		fmt.Printf("🚑 Here is the vault health [response.ServerTimeUTC] : %v", response.ServerTimeUTC)
		fmt.Println()
		fmt.Printf("🚑 Here is the vault health [response.Standby] : %v", response.Standby)
	}
}

/**
 * The unsealVault method will pick up the value of the unseal tokens
 * from the UNSEAL_TOKENS variable in the ./.env file, (thanks to github.com/joho/godotenv/autoload), which value is a comma-separated list of tokens.
 * e.g. with 17 Unseal Tokens : see ./.env
 *
 **/
// func unsealVault(addr string) (*vapi.SealStatusResponse, error) {
func unsealVault(addr string) bool {
	client := newVault(addr)
	// toReturn, err := client.Sys().EnableAuth()
	// theUnsealTokens := conf.UnsealTokens
	//response, err := client.Sys().Unseal("uCnIbTAyd3RgYLBv/GneFAqQ0uWEvmEQG1bg15MN3E4=")
	var wasVaultSuccessfullyUnsealed bool

	var numberOfUnsealTokens = len(conf.UnsealTokens)
	fmt.Printf("🏈💪 Navy-Seals 💪🏈📣 - [func unsealVault(addr string) bool] - numberOfUnsealTokens : %v", numberOfUnsealTokens)
	if numberOfUnsealTokens >= 1 {
		fmt.Printf("🏈💪 Navy-Seals 💪🏈📣 - [func unsealVault(addr string) bool] - conf.UnsealTokens[0] : %v", conf.UnsealTokens[0])
	}
	for i, token := range conf.UnsealTokens {
		var resp *vapi.SealStatusResponse
		var err error
		resp, err = client.Sys().Unseal(token)
		if err != nil {
			fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - ERROR using unseal key %d on %v: %w", i+1, addr, err))
			continue
		}
		if !resp.Sealed {
			fmt.Println(fmt.Errorf("✅ 🏈💪 Navy-Seals 💪🏈📣 ✅ - (Vault was sealed) %v now unsealed with tokens", addr))
			wasVaultSuccessfullyUnsealed = true
			break
		}
	}

	// fmt.Printf("The %v vault was successfully unsealed!", addr)
	return wasVaultSuccessfullyUnsealed
}

/**
 * Awesome that method works great!
 * We will need to check about GPG encrypting the shared keys
 **/
func initVault(addr string) (*vapi.InitResponse, error) {
	client := newVault(addr)
	// toReturn, err := client.Sys().EnableAuth()
	// theUnsealTokens := conf.UnsealTokens
	//response, err := client.Sys().Unseal("uCnIbTAyd3RgYLBv/GneFAqQ0uWEvmEQG1bg15MN3E4=")

	var initRequest = new(vapi.InitRequest)
	initRequest.SecretShares = conf.UnsealkeysNb
	initRequest.SecretThreshold = conf.UnsealKeysTreshold
	var response *vapi.InitResponse
	var err error
	response, err = client.Sys().Init(initRequest)
	//var wasVaultSuccessfullyInitialized bool

	if err == nil {

		fmt.Printf("🏈💪 Navy-Seals 💪🏈📣 - [func initVault(addr string)] - successfully init vault : %v", response)
	} else {
		fmt.Printf("🏈💥 Navy-Seals 💥🏈📣 - ERROR - [func initVault(addr string)] - failed to init vault : %v", response)
	}

	return response, err
}

func main() {
	fmt.Println("!! Welcome to (Navy)-Seals !!")
	var err error
	if _, err = flags.Parse(conf); err != nil {
		var ferr *flags.Error
		if errors.As(err, &ferr) && ferr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if conf.Version {
		fmt.Printf("(Navy)-Seals version: %s [%s] (%s, %s), compiled %s\n", version, commit, runtime.GOOS, runtime.GOARCH, date) //nolint:forbidigo
		os.Exit(0)
	}
	// initVault("http://192.168.1.16:8200")
	showVaultStatus("http://192.168.1.16:8200")
	sealVault("http://192.168.1.16:8200")
	fmt.Println()
	fmt.Println(" 🏈💪 Navy-Seals 💪🏈📣 :  NOW VAULT IS SEALED ❗")
	fmt.Println()
	showVaultStatus("http://192.168.1.16:8200")

	unsealVault("http://192.168.1.16:8200")
	fmt.Println()
	fmt.Println(" 🏈💪 Navy-Seals 💪🏈📣 :  NOW VAULT IS UNSEALED ❗")
	fmt.Println()
	showVaultStatus("http://192.168.1.16:8200")

}
