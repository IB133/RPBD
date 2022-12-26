package app

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/IB133/RPBD/final_project/internal/db"
	"github.com/IB133/RPBD/final_project/internal/models"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	router *httprouter.Router
	store  *db.Service
	ctx    context.Context
}

func newServer(db *db.Service, ctx context.Context) *server {

	s := &server{
		router: httprouter.New(),
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
	s.router.POST("/sign-up", s.handleUsersCreate)
	s.router.POST("/sign-in", s.handleUsersAuth)
	s.router.POST("/addcomment", s.handlerAddComment)
	s.router.GET("/commentlist/:id", s.handlerGetCommentList)
	s.router.POST("/user/addnews", s.handlerAddNewsByUser)
	s.router.POST("/moder/addnews", s.handlerAddNewsByModer)
	s.router.PUT("/moder/acceptnews", s.handlerCheckoutNews)
	s.router.DELETE("/moder/deletecoms/:id", s.handlerDeleteComms)
	s.router.GET("/moder/unpostednews", s.handlerUnpostedNews)
	s.router.GET("/user/commentslist/:id", s.handlerUserCommsList)
}

func (s *server) handleUsersCreate(w http.ResponseWriter, r *http.Request, pm httprouter.Params) {
	r.ParseForm()
	u := &models.User{
		Email:    r.Form.Get("email"),
		Login:    r.Form.Get("login"),
		Password: r.Form.Get("password"),
	}
	if err := s.store.Registration(s.ctx, u); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	u.ClearPassword()
	s.respond(w, r, http.StatusCreated, u)
}

func (s *server) handleUsersAuth(w http.ResponseWriter, r *http.Request, pm httprouter.Params) {
	r.ParseForm()
	u := &models.User{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}
	//Не доконца понял как использовать jwt токен, но вроде его нужно хранить в куки, обещаю разобраться.
	token, err := s.store.AuthAndGenerateToken(s.ctx, u)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusCreated, map[string]interface{}{
		"token": token,
	})
}

func (s *server) handlerAddComment(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	uid, _ := strconv.Atoi(r.Form.Get("user_id"))
	nid, _ := strconv.Atoi(r.Form.Get("news_id"))
	rid, _ := strconv.Atoi(r.Form.Get("reply_id"))
	c := &models.Comments{
		User_id:    uid,
		News_id:    nid,
		Reply_to:   rid,
		Created_at: r.Form.Get("created_at"),
		Content:    r.Form.Get("content"),
	}
	if err := s.store.AddComment(s.ctx, c); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusCreated, c)
}

func (s *server) handlerGetCommentList(w http.ResponseWriter, r *http.Request, pm httprouter.Params) {
	uid, _ := strconv.Atoi(pm.ByName("id"))
	mp, err := s.store.GetComments(s.ctx, uid)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusCreated, mp)
}

func (s *server) handlerAddNewsByUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	uid, _ := strconv.Atoi(r.Form.Get("user_id"))
	n := &models.News{
		User_id: uid,
		Title:   r.Form.Get("title"),
		Theme:   r.Form.Get("theme"),
		Content: r.Form.Get("content"),
	}
	if err := s.store.AddNewsByUser(s.ctx, n); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusCreated, n)
}

func (s *server) handlerAddNewsByModer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	uid, _ := strconv.Atoi(r.Form.Get("moder_id"))
	n := &models.News{
		Moder_id:  uid,
		Title:     r.Form.Get("title"),
		Theme:     r.Form.Get("theme"),
		Content:   r.Form.Get("content"),
		Posted_at: r.Form.Get("posted_at"),
	}
	if err := s.store.AddNewsByModer(s.ctx, n); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusCreated, n)
}

func (s *server) handlerCheckoutNews(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	uid, _ := strconv.Atoi(r.Form.Get("moder_id"))
	nid, _ := strconv.Atoi(r.Form.Get("id"))
	p, _ := strconv.ParseBool(r.Form.Get("posted"))
	n := &models.News{
		Moder_id:   uid,
		Posted:     p,
		Posted_at:  r.Form.Get("posted_at"),
		Moder_comm: r.Form.Get("moder_comm"),
		Id:         nid,
	}
	if err := s.store.CheckoutNews(s.ctx, n); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusCreated, n)
}

func (s *server) handlerDeleteComms(w http.ResponseWriter, r *http.Request, pm httprouter.Params) {
	cid, _ := strconv.Atoi(pm.ByName("id"))
	err := s.store.DeleteComm(s.ctx, cid)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusCreated, "Succeful delete")
}

func (s *server) handlerUnpostedNews(w http.ResponseWriter, r *http.Request, pm httprouter.Params) {
	mp, err := s.store.GetUnpostedNewsList(s.ctx)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusCreated, mp)
}

func (s *server) handlerUserCommsList(w http.ResponseWriter, r *http.Request, pm httprouter.Params) {
	uid, _ := strconv.Atoi(pm.ByName("id"))
	mp, err := s.store.GetUserCommentsList(s.ctx, uid)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusCreated, mp)
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
