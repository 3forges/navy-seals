package api

import (
	"errors"
	"fmt"
	handover_telegram "navy-seals/handover/telegram"
	http "net/http"

	"github.com/gin-gonic/gin"
)

var (
	botClients = []*handover_telegram.BotClient{}
)

/**
 * --------------------------------------------------
 * --------------------------------------------------
 * ++++++++++++++++++++++++++++++++++++++++++++++++++
 * --------------------------------------------------
 * --------------------------------------------------
 *                    GENERATE USER ID
 * --------------------------------------------------
 * --------------------------------------------------
 * ++++++++++++++++++++++++++++++++++++++++++++++++++
 * --------------------------------------------------
 * --------------------------------------------------
 **/

/**
 * ---
 * Ok so here is the idea :
 * You will generate a BotClient, with the UUID of the BotClient, you
 * will next send an apI request to trigger the loop which polls the telegram API to check if the UUID was sent
 **/
func GetBotUserUniqueID(c *gin.Context) {
	newclient, err := handover_telegram.NewBotClient()

	if err != nil {
		fmt.Println("DEBUG POINT JBL")
		fmt.Println(fmt.Sprintf("DEBUG POINT JBL the error is : %v", err))
		c.IndentedJSON(http.StatusInternalServerError, err)
	} else {
		botClients = append(botClients, newclient)
		c.IndentedJSON(http.StatusOK, &BotUserUniqueID{
			UniqueId: newclient.GetBotUserUniqueID(),
		})

	}
}

type BotUserUniqueID struct {
	UniqueId string
}

/**
* --------------------------------------------------
 * LIBRARY FUNCTIONS to move to [handover/telegram/bot.go]
 * --------------------------------------------------
 **/

/**
 * --------------------------------------------------
 * --------------------------------------------------
 * ++++++++++++++++++++++++++++++++++++++++++++++++++
 * --------------------------------------------------
 * --------------------------------------------------
 *                    START POLLING
 * --------------------------------------------------
 * --------------------------------------------------
 * ++++++++++++++++++++++++++++++++++++++++++++++++++
 * --------------------------------------------------
 * --------------------------------------------------
 **/

/**
 * Wait for the user to send
 **/
func WaitUserMessageByID(c *gin.Context) {
	id := c.Param("id")
	// Now we need to retrieve
	client, err := GetBotClientByID(id)
	if err != nil {
		fmt.Println(fmt.Sprintf("DEBUG JBL getting the bot client encountered an error: %v", err))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("BotClient of ID [%v] not found", id)})
	} else {
		client.WaitForUserToSendUUIDtoBot()
		client.SendTestMessage("Par Toutatis! Les Gaulois!")
		c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("BotClient of ID [%v] is now Ready to send the Unseal Key QR code!  found", id)})
	}

}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func GetBotClientByID(id string) (*handover_telegram.BotClient, error) {

	var err error

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, client := range botClients {
		if client.ID == id {
			return client, nil
		}
	}

	err = errors.New(fmt.Sprintf("BotClient of ID [%v] not found", id))
	return &handover_telegram.BotClient{}, err
}
