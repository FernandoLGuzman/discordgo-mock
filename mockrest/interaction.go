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

func (roundTripper *RoundTripper) addHandlersInteraction(apiVersion string) {
	pathInteractions := fmt.Sprintf("%s/%s", apiVersion, resourceInteractions)

	subrouter := roundTripper.router.PathPrefix(pathInteractions).Subrouter()

	pathInteractionIDCallback := fmt.Sprintf("/%s/%s/callback", resourceInteractionID, resourceInteractionToken)

	postHandlers := subrouter.Methods(http.MethodPost).Subrouter()
	postHandlers.HandleFunc(pathInteractionIDCallback, roundTripper.interactionCallbackResponse)
}

func (roundTripper *RoundTripper) interactionCallbackResponse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	interactionID := vars[resourceInteractionIDKey]

	interaction := &discordgo.InteractionResponse{}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&interaction)
	if err != nil {
		log.Println(err)
	}

	channel, err := roundTripper.state.Channel(mockconstants.TestChannel)
	if err != nil {
		sendError(w, err)

		return
	}

	if interaction.Data == nil {
		interaction.Data = &discordgo.InteractionResponseData{}
	}

	message := &discordgo.Message{
		ID:         interactionID,
		ChannelID:  mockconstants.TestChannel,
		GuildID:    mockconstants.TestGuild,
		Content:    interaction.Data.Content,
		TTS:        interaction.Data.TTS,
		Embeds:     interaction.Data.Embeds,
		Components: interaction.Data.Components,
		Flags:      interaction.Data.Flags,
	}

	channel.LastMessageID = message.ID
	channel.MessageCount++
	channel.Messages = append(channel.Messages, message)

	err = roundTripper.state.MessageAdd(message)
	if err != nil {
		sendError(w, err)

		return
	}

	w.WriteHeader(http.StatusNoContent)
	sendJSON(w, nil)
}
