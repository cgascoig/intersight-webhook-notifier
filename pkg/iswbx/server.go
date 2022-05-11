package iswbx

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cgascoig/intersight-webex/pkg/storage"
	"github.com/cgascoig/intersight-webex/pkg/webexbotkit"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/run/v1"
)

type Server struct {
	store *storage.Store

	bot *webexbotkit.Bot

	serviceURL string
}

func NewServer() *Server {
	webexAuthToken := os.Getenv("WEBEX_AUTH_TOKEN")
	if webexAuthToken == "" {
		logrus.Warn("WebEx authentication token (WEBEX_AUTH_TOKEN) is empty")
	}

	s := &Server{
		store: storage.NewStore(context.Background()),
	}

	s.bot = webexbotkit.NewWebExBot(webexAuthToken, s.messageHandler())

	return s
}

func (s *Server) messageHandler() webexbotkit.MessageHandler {
	return func(msg webexbotkit.Message) {
		tokens := strings.Split(msg.Message, " ")
		if len(tokens) == 1 && tokens[0] == "setup" {
			s.setupMessage(msg.RoomID)
		}
	}
}

func (s *Server) setupMessage(roomID string) {
	s.bot.SendMessageToRoomID(roomID, fmt.Sprintf("In Intersight, go to Settings -> Webhooks -> Add Webhook and use '%s/is/%s' as the Payload URL", s.serviceURL, roomID))
}

func (s *Server) IntersightHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("Received Intersight request with invalid body: %v", err)
			return
		}

		wh := map[string]interface{}{}
		err = json.Unmarshal(body, &wh)
		if err != nil {
			logrus.Errorf("Received Intersight request with invalid body: %v", err)
			return
		}

		logrus.WithField("subscription", vars["subscription"]).WithField("body", wh).Debug("IntersightHandler handling request for subscription")

		msg := webhookToMessage(wh)

		if msg == "" {
			return
		}

		s.bot.SendMessageToRoomID(vars["subscription"], msg)
	}
}

func (s *Server) Run() error {
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logrus.Infof("defaulting to port %s", port)
	}

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		logrus.Error("PROJECT_ID not specified")
	}

	serviceName := os.Getenv("K_SERVICE")
	if serviceName == "" {
		logrus.Error("K_SERVICE not specified")
	}

	runClient, err := run.NewService(context.Background())
	if err == nil {
		res, err := runClient.Namespaces.Services.List(fmt.Sprintf("namespaces/%s", projectID)).Do()
		if err == nil {
			for _, item := range res.Items {
				if item.Metadata.Name == serviceName {
					s.serviceURL = item.Status.Url
					break
				}
			}
		}
	}

	logrus.Infof("Got my URL as %s", s.serviceURL)

	r := mux.NewRouter()
	r.HandleFunc("/is/{subscription}", s.IntersightHandler())
	r.HandleFunc("/webex", s.bot.HttpHandler())

	return http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
