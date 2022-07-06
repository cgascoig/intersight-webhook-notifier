package iswbx

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/cgascoig/intersight-webhook-notifier/pkg/storage"
	"github.com/cgascoig/intersight-webhook-notifier/pkg/webexbotkit"
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
		// tokens := strings.Split(msg.Message, " ")
		// if len(tokens) == 1 && tokens[0] == "setup" {
		// 	s.setupMessage(msg.RoomID)
		// }

		// Always just resond with the setup message
		s.setupMessage(msg.RoomID)
	}
}

const setupMessage = `Hi, I'm the Intersight Notification Bot. 

To get setup, in Intersight:

* Go to Settings -> Webhooks -> Add Webhook. 
	* Use %s/is/%s as the Payload URL
	* Enter any string in the Secret field (it is not used at this time)
	* Add the Event subscriptions that you are interested in. Currently this bot supports cond.Alarm object types for the event subscriptions - other event types will still be forward to WebEx but in raw form. 

Disclaimer: This bot is an example of how Intersight Webhooks can be used to enable notifications of Intersight events. You are free to use it but you acknowledge that the bot is not supported by Cisco TAC, has no SLA and may stop working at any time.

Bot build details: %s
`

func (s *Server) setupMessage(roomID string) {
	s.bot.SendMessageToRoomID(roomID, fmt.Sprintf(setupMessage, s.serviceURL, roomID, getVersionString()))
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

		if isWorkflowUpdateCleanup(wh, time.Now()) {
			logrus.WithField("body", wh).Debug("Ingoring workflow update that appears to be cleanup")
			return
		}

		msg := webhookToMessage(wh)

		if msg == "" {
			return
		}

		s.bot.SendMessageToRoomID(vars["subscription"], msg)
	}
}

func isWorkflowUpdateCleanup(wh map[string]interface{}, tm time.Time) bool {
	eventObjectType, ok := wh["EventObjectType"]
	if !ok {
		logrus.Error("Event has no object type")
		return false
	}

	operation, ok := wh["Operation"]
	if !ok {
		logrus.Error("Webhook missing operation")
		return false
	}

	switch eventObjectType {
	case "workflow.WorkflowInfo":
		switch operation {
		case "Modified":
			if eventIntf, ok := wh["Event"]; ok {
				if event, ok := eventIntf.(map[string]interface{}); ok {
					if cleanupIntf, ok := event["CleanupTime"]; ok {
						if cleanup, ok := cleanupIntf.(string); ok {
							cleanupTime, err := time.Parse("2006-01-02T15:04:05.000Z", cleanup)
							if err != nil {
								return false
							}

							if tm.After(cleanupTime.Add(time.Minute * -10)) {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func (s *Server) Run() error {
	logrus.Infof("intersight-webhook-notifier starting, version: %s", getVersionString())

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

func getVersionString() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "<error getting build info>"
	}

	var vcsCommit string
	var vcsDirty bool
	var foundVcsCommit, foundVcsDirty bool
	for _, setting := range bi.Settings {
		if setting.Key == "vcs.revision" {
			vcsCommit = setting.Value
			foundVcsCommit = true
		}
		if setting.Key == "vcs.modified" {
			switch setting.Value {
			case "true":
				vcsDirty = true
				foundVcsDirty = true
			case "false":
				vcsDirty = false
				foundVcsDirty = true
			default:
				foundVcsDirty = false
			}
		}
	}

	if foundVcsCommit == false || foundVcsDirty == false {
		return "<build info not found>"
	}

	var dirtyStr string
	if vcsDirty {
		dirtyStr = "(dirty)"
	}

	return fmt.Sprintf("Commit: %s%s", vcsCommit, dirtyStr)
}
