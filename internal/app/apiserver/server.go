package apiserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/rysya2000/http-rest-api/internal/app/model"
	"github.com/rysya2000/http-rest-api/internal/app/store"
	uuid "github.com/satori/go.uuid"
)

const (
	ctxKeyUser  ctxKey = iota
	sessionName        = "restapi"
)

type ctxKey int8

type server struct {
	router       *pathResolver
	errorLog     *log.Logger
	infoLog      *log.Logger
	store        store.Store
	sessionStore redis.Client
}

func newServer(store store.Store, sessionStore redis.Client) *server {
	s := &server{
		router:       newPathResolver(),
		errorLog:     log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:      log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Add("GET /private/whoami", s.authenticateUser(s.handleWhoami()))

	s.router.Add("POST /users", s.handleUsersCreate())
	s.router.Add("POST /sessions", s.handleSessionsCreate())
}

func (s *server) authenticateUser(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val, err := s.sessionStore.Get("user_id").Result()
		if err != nil {
			s.errorLog.Println(err)
			s.error(w, r, http.StatusUnauthorized, store.ErrNotAuthenticated)
			return
		}
		id, err := strconv.Atoi(val)
		if err != nil {
			s.errorLog.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u, err := s.store.User().Find(id)
		if err != nil {
			s.errorLog.Println(err)
			s.error(w, r, http.StatusUnauthorized, store.ErrNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, store.ErrIncorrectEmailOrPassword)
			return
		}
		err = s.sessionStore.Set("user_id", u.ID, time.Duration(time.Second*1800)).Err()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.CookieSet(w, r, u.ID)

		s.respond(w, r, http.StatusOK, nil)
	}
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

func (s *server) CookieSet(w http.ResponseWriter, r *http.Request, nameid int) {
	u := uuid.NewV4()

	http.SetCookie(w, &http.Cookie{
		Name:   sessionName,
		Value:  u.String(),
		MaxAge: 1800,
	})
}
