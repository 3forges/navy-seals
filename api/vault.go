package api

import (
	"fmt"
	"navy-seals/config"
	http "net/http"

	"github.com/gin-gonic/gin"
)

type VaultStatus struct {
	Sealed bool `json:"sealed"`
}

var (
	// defaultVaultInitParams = &VaultInitParams{
	// 	UnsealKeysNb:       37,
	// 	UnsealKeysTreshold: 23,
	// }
	exampleVaultStatus = &VaultStatus{
		Sealed: false,
	}
)

// var (
// 	defaultVaultInitParams  = VaultInitParams{
//         UnsealKeysNb: 37,
//         UnsealKeysTreshold: 23,
// 	}
// )

/*
var (
	unsealKeys = []unsealKey{
		{ID: "1", Name: "TinTin", OnwerID: "John Coltrane", SecretValue: "cy5ORmNpR1RFTUNrNHJaQUlRUG5ZakZ4b29Jek1DQndZZEJpWQo="},
		{ID: "2", Name: "Haddock", OnwerID: "Gerry Mulligan", SecretValue: "aHZzLm1ZNDh0VzVWVk1COFVweE80dFBYdUY1MAo="},
		{ID: "3", Name: "Tournesol", OnwerID: "Sarah Vaughan", SecretValue: "aHZzLkpJeUVnWm9ROUljcDYzQ2ZZOVYxNXFtRwo="},
	}
)

*/

/**
 * ---
 * If the Function name does not start with a
 * Capital Letter, it can't be imported and
 * used from the [main.go]
 * ---
 * getUnsealKeys responds with the list of all unsealKeys as JSON.
 **/
func GetVaultStatus(c *gin.Context) {
	vaultStatusResponse, err := FetchVaultStatus()
	if err != nil {
		// logger.Fatalf("error creating vault client: %v", err)
		fmt.Printf(" [GetVaultStatus(c *gin.Context)] - Error getting vault status: %v", err)
	}
	c.IndentedJSON(http.StatusOK, vaultStatusResponse)

}

/**
 * Vault Inii
 * ENDPOINT /vault-init
 * -
 * Payload should be of the following form:
 * {
 *   keys_nb: 73,
 *   keys_treshold: 35
 * }
 **/
func InitVault(ctx *gin.Context) {
	// keys_nb := ctx.Param("keys_nb")
	// keys_threshold := ctx.Param("keys_threshold")
	var initVaultParams VaultInitParams

	if err := ctx.ShouldBindJSON(&initVaultParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("🏈💥 Navy-Seals 💥🏈📣 - ERROR initializing vault %v, the payload you provided is not of the expected form : %v", config.ApiConfig.VaultAddress, err.Error())})
		return
	}

	fmt.Printf("🏈 Navy-Seals 🏈📣 - [InitVault] - received JSON PAYLOAD IS: %v", initVaultParams)
	//fmt.Printf("🏈 Navy-Seals 🏈📣 - [InitVault] - received JSON PAYLOAD IS: %v", initVaultParams)
	fmt.Printf("🏈 Navy-Seals 🏈📣 - [InitVault(c *gin.Context)] - received JSON PAYLOAD IS: params.UnsealKeysNb = %v // params.UnsealKeysTreshold = %v", initVaultParams.UnsealKeysNb, initVaultParams.UnsealKeysTreshold)

	initResponse, err := ExecuteInitVault(initVaultParams)

	if err != nil {
		fmt.Println(fmt.Errorf("🏈💥 Navy-Seals 💥🏈📣 - ERROR initializing vault %v: %v", config.ApiConfig.VaultAddress, err))
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("🏈💥 Navy-Seals 💥🏈📣 - ERROR initializing vault %v: %v", config.ApiConfig.VaultAddress, err)})
		return
	} else {
		ctx.IndentedJSON(http.StatusCreated, initResponse)
		// return
	}
}

// This will be a POST?PUT?
func SealVault(c *gin.Context) {
	// Here I seal the vault, and I return the vault status in the gin context
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "OpenBAO Vault now sealed"})
}

// This will be a POST?PUT?
func UnsealVault(c *gin.Context) {
	// Here I Unseal the vault, and I return the vault status in the gin context
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "OpenBAO Vault now unsealed"})
}
