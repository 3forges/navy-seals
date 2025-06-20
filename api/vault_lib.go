package api

import (
	"fmt"
	"navy-seals/client"
	"navy-seals/config"
	"os"
	"time"

	vapi "github.com/openbao/openbao/api/v2"
)

var (
	clientPool = client.NewVaultClientPool(70, 2*time.Second)
	// defaultVaultInitParams = &VaultInitParams{
	// 	UnsealKeysNb:       37,
	// 	UnsealKeysTreshold: 23,
	// }
)

// unsealKey represents data about a record unsealKey.

//	type VaultInitParams struct {
//		KeysNb        int `json:"KeysNb"`
//		KeysThreshold int `json:"KeysThreshold"`
//	}
type VaultInitParams struct {
	UnsealKeysNb       int `json:"UnsealKeysNb"`
	UnsealKeysTreshold int `json:"UnsealKeysTreshold"`
}
type VaultUnsealParams struct {
	Key string `json:"Key"`
}

// type VaultInitParams struct {
// 	UnsealKeysNb       int
// 	UnsealKeysTreshold int
// }

func FetchVaultStatus() (*vapi.HealthResponse, error) {

	vc, err := clientPool.BorrowVaultClient()
	//vc, err := client.GetVaultClient()
	if err != nil {
		fmt.Printf(" NAVY SEALS [FetchVaultStatus(c *gin.Context)] - Error getting vault client: %v", err)
	}
	vaultStatusResponse, err := vc.GetClient().Sys().Health()
	if err != nil {
		// logger.Fatalf("error creating vault client: %v", err)
		fmt.Printf(" NAVY SEALS [FetchVaultStatus(c *gin.Context)] - Error getting vault health: %v", err)
	}

	if err != nil {
		// logger.Fatalf("error creating vault client: %v", err)
		fmt.Printf("NAVY SEALS [FetchVaultStatus(c *gin.Context)] - Error getting vault health: %v", err)
	} else {
		fmt.Printf("NAVY SEALS [FetchVaultStatus(c *gin.Context)] - 🚑 Here is the vault health [response.Initialized] : %v", vaultStatusResponse.Initialized)
		fmt.Println()
		fmt.Printf("NAVY SEALS [FetchVaultStatus(c *gin.Context)] - 🚑 Here is the vault health [response.ClusterName] : %v", vaultStatusResponse.ClusterName)
		fmt.Println()
		fmt.Printf("NAVY SEALS [FetchVaultStatus(c *gin.Context)] - 🚑 Here is the vault health [response.Sealed] : %v", vaultStatusResponse.Sealed)
		fmt.Println()
		fmt.Printf("NAVY SEALS [FetchVaultStatus(c *gin.Context)] - 🚑 Here is the vault health [response.Version] : %v", vaultStatusResponse.Version)
		fmt.Println()
		fmt.Printf("NAVY SEALS [FetchVaultStatus(c *gin.Context)] - 🚑 Here is the vault health [response.ServerTimeUTC] : %v", vaultStatusResponse.ServerTimeUTC)
		fmt.Println()
		fmt.Printf("NAVY SEALS [FetchVaultStatus(c *gin.Context)] - 🚑 Here is the vault health [response.Standby] : %v", vaultStatusResponse.Standby)
	}
	clientPool.ReleaseVaultClient(vc)
	return vaultStatusResponse, err
}

/**
 * Initializes the vault, and saves Unseal Token as QR codes, and Root Token as a file
 **/
func ExecuteInitVault(params VaultInitParams) (*vapi.InitResponse, error) {
	fmt.Printf("🏈 Navy-Seals 🏈📣 - [ExecuteInitVault(params VaultInitParams)] - received JSON PAYLOAD IS: params.UnsealKeysNb = %v // params.UnsealKeysTreshold = %v", params.UnsealKeysNb, params.UnsealKeysTreshold)
	statusResponse, statusErr := FetchVaultStatus()
	if statusErr != nil {
		fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - [ExecuteInitVault(params VaultInitParams)] - ERROR querying vault status key on %v: %v", config.ApiConfig.VaultAddress, statusResponse))
		// os.Exit(7)
		return nil, statusErr
	}

	vc, err := clientPool.BorrowVaultClient()
	//vc, err := client.GetVaultClient()
	if err != nil {
		fmt.Printf(" [ExecuteInitVault(params VaultInitParams)] - Error getting vault client: %v", err)
	}

	var initResponse *vapi.InitResponse = nil
	var initErr error = nil

	if !statusResponse.Initialized {
		var initRequest = new(vapi.InitRequest)
		initRequest.SecretShares = params.UnsealKeysNb
		initRequest.SecretThreshold = params.UnsealKeysTreshold

		initResponse, initErr = vc.GetClient().Sys().Init(initRequest)
		//var wasVaultSuccessfullyInitialized bool
		if initErr == nil {
			fmt.Printf("🏈💪 Navy-Seals 💪🏈📣 - [ExecuteInitVault(params VaultInitParams)] - successfully init vault : %v", initResponse)
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
				fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - ERROR - [ExecuteInitVault(params VaultInitParams)] - An error occured while trying to create [%v] secret file, the error was: %v", root_token_file, rootTokenFileCreationErr))
				os.Exit(11)
			}
			whatheheck, writeToRootTokenFileErr := rootTokenFile.WriteString(initResponse.RootToken)
			if writeToRootTokenFileErr != nil {
				fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - ERROR - [ExecuteInitVault(params VaultInitParams)] - An error occured while trying to write OpenBAO vault to [%v] secret file, the error was: %v", root_token_file, writeToRootTokenFileErr))
				// rootTokenFile.Sync()
				// os.Exit(11)
				return nil, writeToRootTokenFileErr
			}
			rootTokenFile.Sync()
			fmt.Println(fmt.Sprintf("🏈❕ Navy-Seals ❕🏈📣 - ERROR - [ExecuteInitVault(params VaultInitParams)] - Succcessfully wrote OpenBAO vault root token to [%v] secret file", root_token_file))
			fmt.Println(fmt.Sprintf("🏈❕ Navy-Seals ❕🏈📣 - After writing Root Token to secret file, 'whatheheck' is: %v", whatheheck))
			// os.Setenv("BAO_TOKEN", initResponse.RootToken) // This does not work because it would only work for children processes, not the currently running one

		} else {
			fmt.Printf("🏈💥 Navy-Seals 💥🏈📣 - ERROR - [ExecuteInitVault(params VaultInitParams)] - failed to init vault : %v, %v", initResponse, initErr)
			return initResponse, initErr
		}
	} else {
		fmt.Printf("🏈❕ Navy-Seals ❕🏈📣 - WARNING vault already initialized, cancel initializing it: %v", config.ApiConfig.VaultAddress)
		// os.Exit(11)
	}
	clientPool.ReleaseVaultClient(vc)
	return initResponse, err
}

/**
 * Unseals the vault, and saves Unseal Token as QR codes, and Root Token as a file
 **/
func ExecuteUnsealVault(params VaultUnsealParams) (*vapi.SealStatusResponse, error) {
	fmt.Printf("🏈 Navy-Seals 🏈📣 - [ExecuteUnsealVault(params VaultUnsealParams)] - received JSON PAYLOAD IS: params.Key = %v", params.Key)
	statusResponse, statusErr := FetchVaultStatus()
	if statusErr != nil {
		fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - [ExecuteUnsealVault(params VaultUnsealParams)] - ERROR querying vault status key on %v: %v", config.ApiConfig.VaultAddress, statusResponse))
		// os.Exit(7)
		return nil, statusErr
	}
	vc, err := clientPool.BorrowVaultClient()
	//vc, err := client.GetVaultClient()
	if err != nil {
		fmt.Printf(" [ExecuteUnsealVault(params VaultUnsealParams)] - Error getting vault client: %v", err)
	}

	var unsealResponse *vapi.SealStatusResponse = nil
	var unsealErr error = nil

	if statusResponse.Sealed {

		unsealResponse, unsealErr = vc.GetClient().Sys().Unseal(params.Key)
		//initResponse, initErr = vc.GetClient().Sys().Init(initRequest)
		//var wasVaultSuccessfullyInitialized bool
		if unsealErr == nil {
			fmt.Printf("🏈💪 Navy-Seals 💪🏈📣 - [ExecuteUnsealVault(params VaultUnsealParams)] - successfully unsealed vault : %v", unsealResponse)
			return unsealResponse, unsealErr
			// clean the secrets dir
			// os.RemoveAll(tofu_secrets_dir)
			// os.Setenv("BAO_TOKEN", initResponse.RootToken) // This does not work because it would only work for children processes, not the currently running one

		} else {
			fmt.Printf("🏈💥 Navy-Seals 💥🏈📣 - ERROR - [ExecuteUnsealVault(params VaultUnsealParams)] - failed to unseal vault : %v, %v", unsealResponse, unsealErr)
			return unsealResponse, unsealErr
		}
	} else {
		fmt.Printf("🏈❕ Navy-Seals ❕🏈📣 - WARNING vault already initialized, cancel initializing it: %v", config.ApiConfig.VaultAddress)
		// os.Exit(11)
	}
	clientPool.ReleaseVaultClient(vc)
	return unsealResponse, unsealErr
}
