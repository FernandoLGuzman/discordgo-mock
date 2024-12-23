// Package mockrest provides functionality for creating an http.RoundTripper
// that can be used with an *http.Client in a *discordgo.Session to handle
// Discord REST API endpoints and maintain state with a given *discordgo.State.
package mockrest

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
)

const (
	resourceUsers     = "users"
	resourceUserIDKey = "userID"
	resourceUserID    = "{" + resourceUserIDKey + "}"

	resourceRoles     = "roles"
	resourceRoleIDKey = "roleID"
	resourceRoleID    = "{" + resourceRoleIDKey + "}"

	resourceChannels     = "channels"
	resourceChannelIDKey = "channelID"
	resourceChannelID    = "{" + resourceChannelIDKey + "}"
	resourceMessages     = "messages"
	resourceInvites      = "invites"
	resourceMessageIDKey = "messageID"
	resourceMessageID    = "{" + resourceMessageIDKey + "}"

	resourceGuilds     = "guilds"
	resourceGuildIDKey = "guildID"
	resourceGuildID    = "{" + resourceGuildIDKey + "}"

	resourceMembers = "members"

	resourceInteractions        = "interactions"
	resourceInteractionIDKey    = "interactionsID"
	resourceInteractionID       = "{" + resourceInteractionIDKey + "}"
	resourceInteractionTokenKey = "interactionKey"
	resourceInteractionToken    = "{" + resourceInteractionTokenKey + "}"

	resourceWebhooks            = "webhooks"
	resourceWebhookIDKey        = "webhookID"
	resourceWebhookID           = "{" + resourceWebhookIDKey + "}"
	resourceWebhookTokenKey     = "webhookKey"
	resourceWebhookToken        = "{" + resourceWebhookTokenKey + "}"
	resourceWebhookMessageIDKey = "webhookMessageIDKey"
	resourceWebhookMessageID    = "{" + resourceWebhookMessageIDKey + "}"
)

// RoundTripper satisfies http.RoundTripper and handles requests using its
// internal *discordgo.State.
type RoundTripper struct {
	router *mux.Router
	state  *discordgo.State
	// map of token id to interaction
	interactions map[string]*discordgo.Interaction
}

// RoundTrip performs the round trip by routing the request to an appropriate
// handler and returns the result.
func (roundTripper *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	responseRecorder := httptest.NewRecorder()
	roundTripper.router.ServeHTTP(responseRecorder, req)

	return responseRecorder.Result(), nil
}

// NewTransport returns a new http.RoundTripper that handles requests with
// respect to the given state.
func NewTransport(state *discordgo.State) http.RoundTripper {
	router := mux.NewRouter()

	roundTripper := &RoundTripper{
		router:       router,
		state:        state,
		interactions: make(map[string]*discordgo.Interaction),
	}

	apiVersion := "/api/v" + discordgo.APIVersion

	roundTripper.addHandlersGuilds(apiVersion)
	roundTripper.addHandlersChannels(apiVersion)
	roundTripper.addHandlersUsers(apiVersion)
	roundTripper.addHandlersInteraction(apiVersion)
	roundTripper.addHandlersWebhooks(apiVersion)

	return roundTripper
}

func sendError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

func sendJSON(w http.ResponseWriter, object interface{}) {
	respBody, err := json.Marshal(object)
	if err != nil {
		sendError(w, err)

		return
	}

	_, _ = w.Write(respBody)
}

func randString() string {
	const (
		letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		stringLength = 10
	)

	b := make([]byte, stringLength)

	for i := range b {
		//nolint:gosec // random number does not need to be cryptographically secure
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}
