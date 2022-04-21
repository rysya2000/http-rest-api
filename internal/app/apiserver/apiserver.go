package apiserver

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/go-redis/redis"
	"github.com/rysya2000/http-rest-api/internal/app/store/sqlstore"
)

// routers
type pathResolver struct {
	Handlers map[string]http.HandlerFunc
}

func newPathResolver() *pathResolver {
	return &pathResolver{make(map[string]http.HandlerFunc)}
}

func (p *pathResolver) Add(path string, handler http.HandlerFunc) {
	p.Handlers[path] = handler
}

func (p *pathResolver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	check := r.Method + " " + r.URL.Path

	for pattern, handlerFunc := range p.Handlers {
		ok, err := path.Match(pattern, check)
		if ok && err == nil {
			handlerFunc(w, r)
			return
		} else if err != nil {
			fmt.Fprint(w, err)
		}
	}

	http.NotFound(w, r)
}

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + config.RedisLocalhost,
		Password: "",
		DB:       0,
	})
	srv := newServer(store, *sessionStore)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	query, err := ioutil.ReadFile("./migrations/up.sql")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(string(query)); err != nil {
		return nil, err
	}

	return db, nil
}
