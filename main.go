package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	// vapi "github.com/hashicorp/vault/api"

	flags "github.com/jessevdk/go-flags"
	_ "github.com/joho/godotenv/autoload"
	vapi "github.com/openbao/openbao/api/v2"

	// QR CODES GENERATOR

	qrcode "github.com/yeqown/go-qrcode/v2"
	standard "github.com/yeqown/go-qrcode/writer/standard"

	// QR CODES SCANNER
	gozxing "github.com/makiuchi-d/gozxing"
	gozxing_qrcode "github.com/makiuchi-d/gozxing/qrcode"
)

var (
	version = "master"
	commit  = "latest"
	date    = "-"
)

const (
	defaultCheckInterval    = 30 * time.Second
	defaultTimeout          = 15 * time.Second
	configRefreshInterval   = 15 * time.Second
	minimumNodes            = 3
	tofu_secrets_dir        = "./.tofu_secrets"
	unseal_keys_secrets_dir = tofu_secrets_dir + "/.unseal_keys"
	root_token_file         = tofu_secrets_dir + "/.root_token"
	qrcodes_prefix          = ""
	// qrcodes_prefix          = "https://kairos/qrcode/"
)

// Config is a combo of the flags passed to the cli and the configuration file (if used).
type Config struct {
	Version bool `short:"v" long:"version" description:"display the version of navy-seals and exit"`
	Debug   bool `short:"D" long:"debug" description:"enable debugging (extra logging)"`
	Init    bool `short:"i" long:"init" description:"Init the OpenBAO vault and stores the unseal keys as QR codes inside the 'unseal_keys_secrets_dir', also saves the root token in a file of path 'root_token_file', and exits"`
	// UnsealkeysNb       int  `env:"UNSEAL_KEYS_NB"            long:"unseal-keys-nb" description:"number of shared keys to init the openbao vault with (integer, maximum 256)" yaml:"unseal_keys_nb"`
	// UnsealKeysTreshold int  `env:"UNSEAL_KEYS_TRESHOLD"            long:"unseal-keys-treshold" description:"treshold of shared keys to init the openbao vault with (integer, maximum 256)" yaml:"unseal_keys_treshold"`
	UnsealkeysNb       int  `env:"UNSEAL_KEYS_NB"            long:"unseal-keys-nb" description:"number of shared keys to init the openbao vault with (integer, maximum 256)" yaml:"unseal_keys_nb"`
	UnsealKeysTreshold int  `env:"UNSEAL_KEYS_TRESHOLD"            long:"unseal-keys-treshold" description:"treshold of shared keys to init the openbao vault with (integer, maximum 256)" yaml:"unseal_keys_treshold"`
	Seal               bool `short:"s" long:"seal" description:"seal the OpenBAO vault and exit"`
	UnSeal             bool `short:"u" long:"unseal" description:"Unseal the OpenBAO vault from th QR codes found inside the 'tofu_secrets_dir' Folder, and exit"`
	Status             bool `short:"t" long:"status" description:"Show the Status of the OpenBAO vault, and exit"`
	// ConfigPath         string   `env:"CONFIG_PATH" short:"c" long:"config" description:"path to configuration file" value-name:"PATH"`

	Log struct {
		Path   string `env:"LOG_PATH"  long:"path"    description:"path to log output to" value-name:"PATH"`
		Quiet  bool   `env:"LOG_QUIET" long:"quiet"   description:"disable logging to stdout (also: see levels)"`
		Level  string `env:"LOG_LEVEL" long:"level"   default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal"  description:"logging level"`
		JSON   bool   `env:"LOG_JSON"   long:"json"   description:"output logs in JSON format"`
		Pretty bool   `env:"LOG_PRETTY" long:"pretty" description:"output logs in a pretty colored format (cannot be easily parsed)"`
	} `group:"Logging Options" namespace:"log"`

	Environment string `env:"ENVIRONMENT" long:"environment" description:"environment this cluster relates to (for logging)" yaml:"environment"`

	CheckInterval    time.Duration `env:"CHECK_INTERVAL"     long:"check-interval" description:"frequency of sealed checks against nodes" yaml:"check_interval"`
	MaxCheckInterval time.Duration `env:"MAX_CHECK_INTERVAL" long:"max-check-interval" description:"max time that navy-seals will wait for an unseal check/attempt" yaml:"max_check_interval"`

	AllowSingleNode bool     `env:"ALLOW_SINGLE_NODE" long:"allow-single-node"    description:"allow navy-seals to run on a single node" yaml:"allow_single_node" hidden:"true"`
	Nodes           []string `env:"NODES"             long:"nodes" env-delim:","  description:"nodes to connect/provide tokens to (can be provided multiple times & uses comma-separated string for environment variable)" yaml:"vault_nodes"`
	TLSSkipVerify   bool     `env:"TLS_SKIP_VERIFY"   long:"tls-skip-verify"      description:"disables tls certificate validation: DO NOT DO THIS" yaml:"tls_skip_verify"`

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

func loadRootToken() string {
	var valueOfRootToken string
	readFile, err := os.Open(root_token_file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
		valueOfRootToken = fileScanner.Text()
	}

	readFile.Close()
	// os.Setenv("BAO_TOKEN", valueOfRootToken) // This does not work because it would only work for children processes, not the currently running one
	// ~/.bao-token

	// persist the root token
	// initResponse.RootToken
	fmt.Println(fmt.Sprintf("🏈❕ Navy-Seals ❕🏈📣 - [func loadRootToken()] - Succcessfully wrote OpenBAO vault root token to [%v] secret file", root_token_file))

	return valueOfRootToken
}

func sealVault(addr string) {
	client := newVault(addr)
	client.SetToken(loadRootToken())
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

func listUnsealTokensQRcodesPaths() ([]string, error) {
	// var imagePaths []string
	imagePaths := make([]string, 0)
	var errListingUnSealTokens error = nil
	//unsealTokens := make([]string, len(imagePaths))
	//for i := 0; i < len(imagePaths); i++ {}

	dir, err := os.Open(unseal_keys_secrets_dir)
	if err != nil {
		log.Fatal(err)
		fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - [func listUnsealTokensQRcodesPaths()] - ERROR trying to open tofu_secrets_dir=[%v] : %v", unseal_keys_secrets_dir, err))
		os.Exit(1)
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		//fmt.Println(err)
		fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - [func listUnsealTokensQRcodesPaths()] - ERROR trying to read tofu_secrets_dir=[%v] : %v", unseal_keys_secrets_dir, err))
	}

	for _, f := range files {
		if !f.IsDir() {
			fmt.Println(fmt.Sprintf("🏈 Navy-Seals 🏈📣 - [func listUnsealTokensQRcodesPaths()] - found Unseal Token QRcodes file path : %v", unseal_keys_secrets_dir+"/"+f.Name()))
			imagePaths = append(imagePaths, unseal_keys_secrets_dir+"/"+f.Name())
		}
	}
	return imagePaths, errListingUnSealTokens
}

func scanUnsealTokensQRcodes() ([]string, error) {

	unSealTokensListPaths, errListingUnSealTokens := listUnsealTokensQRcodesPaths()

	if errListingUnSealTokens != nil {
		log.Fatal(errListingUnSealTokens)
		fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - [func scanUnsealTokensQRcodes()] - ERROR trying to list Unseal Token QRcodes inside the unseal_keys_secrets_dir=[%v] folder : %v", unseal_keys_secrets_dir, errListingUnSealTokens))
		os.Exit(1)
	}
	fmt.Println(fmt.Sprintf("🏈 Navy-Seals 🏈📣 🔑🔑🔑 - [func scanUnsealTokensQRcodes()] - list of found Unseal Token QRcodes file path : %v", unSealTokensListPaths))

	/*
		imagePaths := []string{
			".sunseal_keys/unseal_key_0.jpeg",
			".sunseal_keys/unseal_key_1.jpeg",
			".sunseal_keys/unseal_key_2.jpeg",
			".sunseal_keys/unseal_key_3.jpeg",
			".sunseal_keys/unseal_key_4.jpeg",
			".sunseal_keys/unseal_key_5.jpeg",
			".sunseal_keys/unseal_key_6.jpeg",
			".sunseal_keys/unseal_key_7.jpeg",
			".sunseal_keys/unseal_key_8.jpeg",
			".sunseal_keys/unseal_key_9.jpeg",
			".sunseal_keys/unseal_key_10.jpeg",
			".sunseal_keys/unseal_key_11.jpeg",
			".sunseal_keys/unseal_key_12.jpeg",
			".sunseal_keys/unseal_key_13.jpeg",
			".sunseal_keys/unseal_key_14.jpeg",
			".sunseal_keys/unseal_key_15.jpeg",
			".sunseal_keys/unseal_key_16.jpeg",
			".sunseal_keys/unseal_key_17.jpeg",
		}

	*/
	// var unsealTokens []string
	unsealTokens := make([]string, len(unSealTokensListPaths))
	for i := 0; i < len(unSealTokensListPaths); i++ {

		fmt.Println(fmt.Sprintf("loading unseal Keys - scanning QR code [%v]", unSealTokensListPaths[i]))
		// open and decode image file
		file, _ := os.Open(unSealTokensListPaths[i])
		img, _, _ := image.Decode(file)

		// prepare BinaryBitmap

		bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

		// decode image
		qrReader := gozxing_qrcode.NewQRCodeReader()
		result, _ := qrReader.Decode(bmp, nil)
		fmt.Println(fmt.Sprintf("loading unseal Keys - scanning QR code [%v], gave the following result : [%v]", unSealTokensListPaths[i], result.GetText()))
		//fmt.Println(result)
		// unsealTokens[i] = result.GetText()
		unsealTokens[i] = strings.Replace(result.GetText(), qrcodes_prefix, "", -1)

	}
	return unsealTokens, nil
}

/**
 * The [unsealVault] method will pick up the value of the unseal tokens
 * by scanning the QR CODES image files in
 * the [unseal_keys_secrets_dir] folder
 *
 *
 **/
func unsealVault(addr string) bool {
	var wasVaultSuccessfullyUnsealed bool = false
	scannedUnsealTokens, errorScanningQrCodes := scanUnsealTokensQRcodes()

	if errorScanningQrCodes != nil {
		fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - ERROR trying to scan Unseal Token QRcodes : %v", errorScanningQrCodes))
	}
	fmt.Printf("🏈💪 Navy-Seals 💪🏈📣 - [func unsealVault(addr string) bool] - scannedUnsealTokens : %v", scannedUnsealTokens)
	client := newVault(addr)

	var numberOfUnsealTokens = len(scannedUnsealTokens)
	fmt.Printf("🏈💪 Navy-Seals 💪🏈📣 - [func unsealVault(addr string) bool] - numberOfUnsealTokens : %v", numberOfUnsealTokens)
	if numberOfUnsealTokens >= 1 {
		fmt.Printf("🏈💪 Navy-Seals 💪🏈📣 - [func unsealVault(addr string) bool] - conf.UnsealTokens[0] : %v", scannedUnsealTokens[0])
	}
	for i, token := range scannedUnsealTokens {
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

	return wasVaultSuccessfullyUnsealed
}

/**
 * Awesome that method works great!
 * We will need to check about GPG encrypting the shared keys
 **/
func initVault(addr string) (*vapi.InitResponse, error) {
	client := newVault(addr)
	var initResponse *vapi.InitResponse = nil
	var initErr error = nil

	statusResponse, statusErr := getVaultStatus(addr)
	if statusErr != nil {
		fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - ERROR querying vault status key on %v: %v", addr, statusResponse))
		os.Exit(7)
	}

	if !statusResponse.Initialized {
		var initRequest = new(vapi.InitRequest)
		initRequest.SecretShares = conf.UnsealkeysNb
		initRequest.SecretThreshold = conf.UnsealKeysTreshold

		initResponse, initErr = client.Sys().Init(initRequest)
		//var wasVaultSuccessfullyInitialized bool
		if initErr == nil {
			fmt.Printf("🏈💪 Navy-Seals 💪🏈📣 - [func initVault(addr string)] - successfully init vault : %v", initResponse)
			// clean the secrets dir
			os.RemoveAll(tofu_secrets_dir)

			var generatedKeys []string = initResponse.KeysB64

			for i := 0; i < len(generatedKeys); i++ {
				generateQRCode(generatedKeys[i], fmt.Sprintf("unseal_key_%v", i))
			}
			// persist the root token
			// initResponse.RootToken
			rootTokenFile, rootTokenFileCreationErr := os.Create(root_token_file)
			if rootTokenFileCreationErr != nil {
				fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - ERROR - [func initVault(addr string)] - An error occured while trying to create [%v] secret file, the error was: %v", root_token_file, rootTokenFileCreationErr))
				os.Exit(11)
			}
			whatheheck, writeToRootTokenFileErr := rootTokenFile.WriteString(initResponse.RootToken)
			if writeToRootTokenFileErr != nil {
				fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - ERROR - [func initVault(addr string)] - An error occured while trying to write OpenBAO vault to [%v] secret file, the error was: %v", root_token_file, writeToRootTokenFileErr))
				// rootTokenFile.Sync()
				os.Exit(11)
			}
			rootTokenFile.Sync()
			fmt.Println(fmt.Sprintf("🏈❕ Navy-Seals ❕🏈📣 - ERROR - [func initVault(addr string)] - Succcessfully wrote OpenBAO vault root token to [%v] secret file", root_token_file))
			fmt.Println(fmt.Sprintf("🏈❕ Navy-Seals ❕🏈📣 - After writing Root Token to secret file, 'whatheheck' is: %v", whatheheck))
			// os.Setenv("BAO_TOKEN", initResponse.RootToken) // This does not work because it would only work for children processes, not the currently running one

		} else {
			fmt.Printf("🏈💥 Navy-Seals 💥🏈📣 - ERROR - [func initVault(addr string)] - failed to init vault : %v, %v", initResponse, initErr)
		}
	} else {
		fmt.Println(fmt.Errorf("🏈❕ Navy-Seals ❕🏈📣 - WARNING vault already initialized, cancel initializing it: %v", addr))
		// os.Exit(11)
	}

	return initResponse, initErr
}

func generateQRCode(unsealKey_B64 string, key_name string) (string, error) {
	var imgFilePath string = fmt.Sprintf("%v/%v.jpeg", unseal_keys_secrets_dir, key_name)
	qrc, err := qrcode.New(fmt.Sprintf(qrcodes_prefix+"%v", unsealKey_B64))
	// qrc, err := qrcode.New(unsealKey_B64)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return imgFilePath, err
	}
	// os.MkdirAll("/tmp/", FileMode)
	os.MkdirAll(unseal_keys_secrets_dir+"/", os.ModePerm)
	w, errCreatingFile := standard.New(imgFilePath)
	if errCreatingFile != nil {
		fmt.Printf("standard.New failed: %v", errCreatingFile)
		return imgFilePath, errCreatingFile
	}
	// save file
	if errSavingImgToFile := qrc.Save(w); errSavingImgToFile != nil {
		fmt.Printf("could not save image: %v", errSavingImgToFile)
		return imgFilePath, errSavingImgToFile
	}

	return imgFilePath, nil
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
	if conf.Init {
		initVault("http://192.168.1.16:8200")
		showVaultStatus("http://192.168.1.16:8200")
		os.Exit(0)
	}
	if conf.Status {
		showVaultStatus("http://192.168.1.16:8200")
		os.Exit(0)
	}

	if conf.Seal {
		sealVault("http://192.168.1.16:8200")
		fmt.Println()
		fmt.Println(" 🏈💪 Navy-Seals 💪🏈📣 :  NOW VAULT IS SEALED ❗")
		fmt.Println()
		showVaultStatus("http://192.168.1.16:8200")
		os.Exit(0)
	}
	if conf.UnSeal {
		// unsealVault("http://192.168.1.16:8200")
		unsealVault("http://192.168.1.16:8200")
		fmt.Println()
		fmt.Println(" 🏈💪 Navy-Seals 💪🏈📣 :  NOW VAULT IS UNSEALED ❗")
		fmt.Println()
		showVaultStatus("http://192.168.1.16:8200")
		os.Exit(0)
	}

}
