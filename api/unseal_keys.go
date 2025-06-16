package api

import (
	"github.com/gin-gonic/gin"
	http "net/http"
)

// unsealKey represents data about a record unsealKey.
type unsealKey struct {
    ID     string  `json:"id"`
    Name  string  `json:"name"`
    OnwerID string  `json:"onwer_id"`
    // SecretValue: The value of the unseal key, base64 encoded, and preferably GPG encrypted.
    SecretValue string `json:"secret_value"`
}


var (
	unsealKeys = []unsealKey{
	}
)

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
func GetUnsealKeys(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, unsealKeys)
}

func AddUnsealKey(c *gin.Context) {
    var newUnsealKey unsealKey

    // Call BindJSON to bind the received JSON to
    // newUnsealKey.
    if err := c.BindJSON(&newUnsealKey); err != nil {
        return
    }

    // Add the new unsealKey to the slice.
    unsealKeys = append(unsealKeys, newUnsealKey)
    c.IndentedJSON(http.StatusCreated, newUnsealKey)
}



// getUnsealKeyByID locates the unsealKey whose ID value matches the id
// parameter sent by the client, then returns that unsealKey as a response.
func GetUnsealKeyByID(c *gin.Context) {
    id := c.Param("id")

    // Loop over the list of unsealKeys, looking for
    // an unsealKey whose ID value matches the parameter.
    for _, a := range unsealKeys {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "unsealKey not found"})
}


// getUnsealKeyByID locates the unsealKey whose ID value matches the id
// parameter sent by the client, then returns that unsealKey as a response.
func GetUnsealKeyQRcode(c *gin.Context) {
    id := c.Param("id")

    // Loop over the list of unsealKeys, looking for
    // an unsealKey whose ID value matches the parameter.
    for _, a := range unsealKeys {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "unsealKey not found"})
}