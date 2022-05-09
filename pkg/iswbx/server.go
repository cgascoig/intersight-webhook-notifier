package iswbx

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/cgascoig/intersight-webex/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	store *storage.Store
}

func NewServer() *Server {
	return &Server{
		store: storage.NewStore(context.Background()),
	}
}

func (s *Server) IntersightHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		body, _ := ioutil.ReadAll(r.Body)
		logrus.WithField("body", string(body)).WithField("subscription", vars["subscription"]).Debug("IntersightHandler handling request for subscription")
	}
}

func (s *Server) Run() error {
	r := mux.NewRouter()
	r.HandleFunc("/is/{subscription}", s.IntersightHandler())

	return http.ListenAndServe(":8080", r)
}
