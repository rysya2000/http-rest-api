package apiserver

import (
	"io"
	"log"
	"net/http"
	"os"
)

type APIServer struct {
	config   *Config
	errorLog *log.Logger
	infoLog  *log.Logger
	router   *http.ServeMux
}

// возможна ошибка в наименовании зависимостей с Большой буквы

func New(config *Config) *APIServer {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return &APIServer{
		config:   config,
		errorLog: errorLog,
		infoLog:  infoLog,
		router:   http.NewServeMux(),
	}
}

func (s *APIServer) Start() error {
	s.configureRouter()

	s.infoLog.Printf("Server: http://localhost%v", s.config.BindAddr)

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

// func (s *APIServer) configureLogger() error {}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}
