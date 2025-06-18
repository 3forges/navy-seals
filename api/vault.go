package api

import (
	"fmt"
	"navy-seals/client"
	http "net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// unsealKey represents data about a record unsealKey.
type VaultInitParams struct {
	UnsealKeysNb       int `json:"keys_nb"`
	UnsealKeysTreshold int `json:"keys_threshold"`
}

type VaultStatus struct {
	Sealed bool `json:"sealed"`
}

var (
	clientPool = client.NewVaultClientPool(7, 150*time.Second)
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
	vc, err := clientPool.BorrowVaultClient()
	//vc, err := client.GetVaultClient()
	if err != nil {
		fmt.Printf(" [GetVaultStatus(c *gin.Context)] - Error getting vault client: %v", err)
	}
	vaultStatusResponse, err := vc.GetClient().Sys().Health()
	if err != nil {
		// logger.Fatalf("error creating vault client: %v", err)
		fmt.Printf(" [GetVaultStatus(c *gin.Context)] - Error getting vault health: %v", err)
	}

	if err != nil {
		// logger.Fatalf("error creating vault client: %v", err)
		fmt.Printf("Error getting vault health: %v", err)
	} else {
		fmt.Printf("🚑 Here is the vault health [response.Initialized] : %v", vaultStatusResponse.Initialized)
		fmt.Println()
		fmt.Printf("🚑 Here is the vault health [response.ClusterName] : %v", vaultStatusResponse.ClusterName)
		fmt.Println()
		fmt.Printf("🚑 Here is the vault health [response.Sealed] : %v", vaultStatusResponse.Sealed)
		fmt.Println()
		fmt.Printf("🚑 Here is the vault health [response.Version] : %v", vaultStatusResponse.Version)
		fmt.Println()
		fmt.Printf("🚑 Here is the vault health [response.ServerTimeUTC] : %v", vaultStatusResponse.ServerTimeUTC)
		fmt.Println()
		fmt.Printf("🚑 Here is the vault health [response.Standby] : %v", vaultStatusResponse.Standby)
	}
	clientPool.ReleaseVaultClient(vc)
	// fetchedVaultStatus any := nil
	// c.IndentedJSON(http.StatusOK, exampleVaultStatus)
	c.IndentedJSON(http.StatusOK, vaultStatusResponse)

}

/**
 * Vault Init
 **/
func InitVault(c *gin.Context) {
	keys_nb := c.Param("keys_nb")
	keys_threshold := c.Param("keys_threshold")

	// initializes the vault
	c.IndentedJSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Initialized OpenBAO Vault with keys_nb = [%v] and keys_threshold = [%v]", keys_nb, keys_threshold)})
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
