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
 * Wait for the user to send the uuid and we can get the ChatID we need to send the QRcode only to him
 **/
func WaitUserMessageByID(c *gin.Context) {
	id := c.Param("id")
	// fmt.Println(fmt.Sprintf(" [WaitUserMessageByID(c *gin.Context)] - DEBUG JBL the telegram BOT ID is : [%v]", id))
	fmt.Printf(" [WaitUserMessageByID(c *gin.Context)] - DEBUG JBL the telegram BOT ID is : [%v]", id)
	fmt.Println("")
	// Now we need to retrieve
	client, err := GetBotClientByID(id)

	fmt.Printf(" [WaitUserMessageByID(c *gin.Context)] - DEBUG JBL the retrieved telegram client has this ID: [%v]", client.ID)
	fmt.Println("")
	fmt.Printf(" [WaitUserMessageByID(c *gin.Context)] - DEBUG JBL the retrieved telegram client has this ChatID: [%v]", client.ChatID)
	fmt.Println("")
	if err != nil {
		fmt.Println(fmt.Sprintf("DEBUG JBL getting the bot client encountered an error: %v", err))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("BotClient of ID [%v] not found", id)})
	} else {
		fmt.Println(fmt.Sprintf(" [WaitUserMessageByID(c *gin.Context)] - JUST BEFORE call to client.WaitForUserToSendUUIDtoBot() - DEBUG JBL the retrieved telegram client has this ID: [%v]", client.ID))
		fmt.Println(fmt.Sprintf(" [WaitUserMessageByID(c *gin.Context)] - JUST BEFORE call to client.WaitForUserToSendUUIDtoBot() - DEBUG JBL the retrieved telegram client has this ChatID: [%v]", client.ChatID))
		client.WaitForUserToSendUUIDtoBot()
		err := client.SendTestMessage("Par Toutatis! Les Gaulois!")
		fmt.Println(fmt.Sprintf(" DEBUG JBL - client.SendTestMessage returned this error:", err))
		client.SendQRCode("./.tofu_secrets/.unseal_keys/unseal_key_1.jpeg")
		c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("BotClient of ID [%v] is now Ready to send the Unseal Key QR code!  found", id)})
	}

}

// GetBotClientByID retrieves from pool the BotClient whose ID value matches the id
// provided parameter, then returns that BotClient as a response, or an error, if there is no BotClient with a matching ID in the pool.
func GetBotClientByID(id string) (*handover_telegram.BotClient, error) {

	var err error

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, client := range botClients {
		if client.ID == id {
			fmt.Printf("[GetBotClientByID(id string)] - BotClient of ID [%v] was successfully found", id)
			return client, nil
		}
	}

	err = errors.New(fmt.Sprintf("[GetBotClientByID(id string)] - BotClient of ID [%v] not found", id))
	return &handover_telegram.BotClient{}, err
}
