package mockrest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ewohltman/discordgo-mock/mockconstants"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
)

func (roundTripper *RoundTripper) addHandlersWebhooks(apiVersion string) {
	pathWebhooks := fmt.Sprintf("%s/%s", apiVersion, resourceWebhooks)

	subrouter := roundTripper.router.PathPrefix(pathWebhooks).Subrouter()

	pathWebhooksIDCallback := fmt.Sprintf("/%s/%s/messages/%s", resourceWebhookID, resourceWebhookToken, resourceWebhookMessageID)

	postHandlers := subrouter.Methods(http.MethodPatch).Subrouter()
	postHandlers.HandleFunc(pathWebhooksIDCallback, roundTripper.webhooksMessagesResponsePatch)
}

func (roundTripper *RoundTripper) webhooksMessagesResponsePatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// webhookID := vars[resourceWebhookIDKey]
	// webhookToken := vars[resourceWebhookTokenKey]
	webhookMessageID := vars[resourceWebhookMessageIDKey]

	webhook := &discordgo.WebhookEdit{}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&webhook)
	if err != nil {
		log.Println(err)
	}

	channel, err := roundTripper.state.Channel(mockconstants.TestChannel)
	if err != nil {
		sendError(w, err)

		return
	}

	message := &discordgo.Message{
		ID:         webhookMessageID,
		ChannelID:  mockconstants.TestChannel,
		GuildID:    mockconstants.TestGuild,
		Content:    *webhook.Content,
		Embeds:     *webhook.Embeds,
		Components: *webhook.Components,
	}

	channel.LastMessageID = message.ID
	channel.MessageCount++
	channel.Messages = append(channel.Messages, message)

	err = roundTripper.state.MessageAdd(message)
	if err != nil {
		sendError(w, err)

		return
	}

	sendJSON(w, message)
}
