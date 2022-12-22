package app

import (
	"net/http"

	"github.com/IB133/RPBD/final_project/internal/db"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *httprouter.Router
	logger *logrus.Logger
	store  *db.Service
}

func newServer(db *db.Service) *server {
	return &server{
		router: httprouter.New(),
		logger: logrus.New(),
		store:  db,
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
