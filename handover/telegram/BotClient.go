package handover_telegram

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	uuid "github.com/gofrs/uuid/v5"
)

type BotClient struct {
	ID     string
	ChatID int64            // the Telegram Chat ID of the Chat between the user and the Telegram Bot
	bot    *tgbotapi.BotAPI // bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
}

func GenerateBotUserUniqueID() (string, error) {
	// Create a Version 4 UUID.
	u2, err := uuid.NewV7()
	if err != nil {
		log.Fatalf("failed to generate UUID V7 : %v", err)
	}
	log.Printf("generated Version 7 UUID %v", u2)
	return u2.String(), err
}

// NewService - our constructor function
func NewBotClient() (*BotClient, error) {
	var err error
	var createdClient *BotClient
	fmt.Printf("Creating new Telegram Bot client")
	// first, generate the UUID
	uuid, err := GenerateBotUserUniqueID()
	if err != nil {
		return nil, err
	}
	// second, create the Bot instance
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	createdClient = &BotClient{
		ID:     uuid,
		ChatID: -1, // The ChatID is initialized to an empty string
		bot:    bot,
	}
	// handles other potentially more complex setup logic
	// for our component, there could be calls to downstream
	// dependencies to check connections etc that could return
	// errors
	return createdClient, nil
}

func (c *BotClient) GetBotUserUniqueID() string {
	return c.ID
}

/**
 * This method will get the ChatID
 **/

func (c *BotClient) WaitForUserToSendUUIDtoBot() bool {
	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := c.bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		// Now that we know we've gotten a new message, we can:
		// 1./ We'll take the Chat ID and set this instance ChatID Filed, to keep it
		c.ChatID = update.Message.Chat.ID
		// 2./ We will construct and send a reply to confirm the Bot can snd message to the specific user
		//
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Par Toutatis! Le bot est prêt à envoyer le QRcode !!! ")

		// We'll also say that this message is a reply to the previous message.
		// For any other specifications than Chat ID or Text, you'll need to
		// set fields on the `MessageConfig`.
		msg.ReplyToMessageID = update.Message.MessageID
		// Okay, we're sending our message off! We don't care about the message
		// we just sent, so we'll discard it.
		if _, err := c.bot.Send(msg); err != nil {
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			panic(err)
		}
		updates.Clear()
	}

	return true
}

func (c *BotClient) SendQRCode(pathToQRcodeFile string) error { // (tgbotapi.Message, error)
	// photoBytes, err := ioutil.ReadFile(pathToQRcodeFile) //("/your/local/path/to/picture.png")
	// if err != nil {
	//   panic(err)
	// }
	// --- // -
	// same as deprecated ioutil.ReadFile

	/*
		photoBytes, err := os.ReadFile(pathToQRcodeFile)
		if err != nil {
			panic(err)
		}
		photoFileBytes := tgbotapi.FileBytes{
			Name:  "YourUnsealKey",
			Bytes: photoBytes,
		}

	*/

	// files []interface{}
	// chatID := 12345678
	// message, err := moi.bot.Send(tgbotapi.NewPhotoUpload(int64(chatID), photoFileBytes))
	// message, err := c.bot.Send(tgbotapi.NewMediaGroup(int64(c.ChatID), photoFileBytes))
	// message, err := c.bot.Send(tgbotapi.NewMediaGroup(int64(c.ChatID), photoFileBytes))
	//.NewChatPhoto(int64(chatID), photoFileBytes))

	msg := tgbotapi.NewPhoto(int64(c.ChatID), tgbotapi.FilePath("tests/image.jpg"))
	msg.Caption = "Your Unseal Key"
	_, err := c.bot.Send(msg)

	if err != nil {
		// t.Error(err)
		fmt.Printf(" ERROR !! [SendQRCode(pathToQRcodeFile string)] ERROR SENDING THE QR CODE : %v", err)
	}
	return err
}

func (c *BotClient) SendTestMessage(pathToQRcodeFile string) error { // (tgbotapi.Message, error)
	msg := tgbotapi.NewMessage(int64(c.ChatID), "A test message from the test library in telegram-bot-api")
	msg.ParseMode = tgbotapi.ModeMarkdown
	_, err := c.bot.Send(msg)
	return err
}

/**
 * newclient, err = NewBotClient()
 **/
