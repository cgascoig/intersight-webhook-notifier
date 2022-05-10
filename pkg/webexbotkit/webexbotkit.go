package webexbotkit

import (
	"encoding/json"
	"net/http"
	"strings"

	wxt "github.com/jbogarin/go-cisco-webex-teams/sdk"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	authToken      string
	messageHandler MessageHandler
	client         *wxt.Client
	me             *wxt.Person
}

type MessageHandler func(msg Message)

type Message struct {
	Bot     *Bot
	RoomID  string
	Message string
}

func NewWebExBot(authToken string, messageHandler MessageHandler) *Bot {
	client := wxt.NewClient()
	client.SetAuthToken(authToken)

	return &Bot{
		authToken:      authToken,
		messageHandler: messageHandler,
		client:         client,
	}
}

func (b *Bot) HttpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Bot HTTP handler called")

		if r.Method != "POST" {
			return
		}

		var webhook wxt.WebhookRequest

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&webhook)
		if err != nil {
			log.Errorf("Error decoding webhook body: %v", err)
		}

		if webhook.Resource != "messages" {
			log.Errorf("Webhook resource is not 'messages'")
			return
		}

		if webhook.Event != "created" {
			log.Errorf("Webhook event is not 'created'")
			return
		}

		if b.isMessageFromSelf(&webhook.Data) {
			log.Info("Ignoring message from myself")
			return
		}

		msg, _, err := b.client.Messages.GetMessage(webhook.Data.ID)
		if err != nil {
			log.Errorf("Error getting message: %v", err)
			return
		}

		msgText := trimMessageText(msg.Text, b.me.DisplayName)

		// b.messageHandler(b, msg.RoomID, msg.Text)
		b.messageHandler(Message{
			Bot:     b,
			RoomID:  msg.RoomID,
			Message: msgText,
		})
	}
}

func trimMessageText(msgText, displayName string) string {
	msgText = strings.TrimSpace(msgText)
	displayNameTokens := strings.Split(displayName, " ")

	for i := len(displayNameTokens); i >= 0; i-- {
		msgText = strings.TrimPrefix(msgText, strings.Join(displayNameTokens[0:i], " "))
	}

	return strings.TrimSpace(msgText)
}

func (b *Bot) isMessageFromSelf(whData *wxt.WebhookRequestData) bool {
	if b.me == nil {
		log.Info("Getting my own details")
		me, _, err := b.client.People.GetMe()
		if err != nil {
			// If there is an error getting my own details, return true so that we don't respond
			log.Errorf("Error getting my own details: %v", err)
			return true
		}
		b.me = me
		log.Infof("My personID is %s and personEmails are %v", b.me.ID, b.me.Emails)
	}

	if whData.PersonID == b.me.ID {
		return true
	}

	return false
}

func (b *Bot) SendMessageToRoomID(roomID string, message string) {
	msgReq := &wxt.MessageCreateRequest{
		RoomID:   roomID,
		Markdown: message,
	}

	_, _, err := b.client.Messages.CreateMessage(msgReq)
	if err != nil {
		log.Errorf("Error posting message: %v", err)
	}
}
