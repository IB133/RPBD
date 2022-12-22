package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/IB133/RPBD/final_project/internal/db"
	"github.com/IB133/RPBD/final_project/internal/models"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *httprouter.Router
	logger *logrus.Logger
	store  *db.Service
	ctx    context.Context
}

func newServer(db *db.Service, ctx context.Context) *server {

	s := &server{
		router: httprouter.New(),
		logger: logrus.New(),
		store:  db,
		ctx:    ctx,
	}
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.POST("/users", s.handleUsersCreate)
}

func (s *server) handleUsersCreate(w http.ResponseWriter, r *http.Request, pm httprouter.Params) {
	// type request struct {
	// 	Login    string `json:"login"`
	// 	Email    string `json:"email"`
	// 	Password string `json:"password"`
	// }
	//req := &request{}
	// if err := json.NewDecoder(r.Body).Decode(req); err != nil {
	// 	s.error(w, r, http.StatusBadRequest, err)
	// 	return
	// }
	r.ParseForm()
	u := &models.User{
		Email:    r.Form.Get("email"),
		Login:    r.Form.Get("login"),
		Password: r.Form.Get("password"),
	}
	log.Printf("%s", u.Email)
	if err := s.store.Registration(s.ctx, u); err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	u.ClearPassword()
	s.respond(w, r, http.StatusCreated, u)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
