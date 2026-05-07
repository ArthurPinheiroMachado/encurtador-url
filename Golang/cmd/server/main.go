package main

import (
	"fmt"
	"golang/db"
	"golang/internal/middlewares"
	"golang/internal/url"
	"golang/internal/util"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	trace := util.CreateErrorContext("main")
	if err := run(); err != nil {
		panic(fmt.Sprintf("\n%s\n", trace.Apply(err)))
	}
}

func run() error {
	trace := util.CreateErrorContext("run")

	user, userErr := util.EnvAsResult("USER")
	pass, passErr := util.EnvAsResult("PASS")

	if err := trace.Join(userErr, passErr); err != nil {
		return trace.Apply(err)
	}

	db, dbErr := db.ConstructDatabase()
	if dbErr != nil {
		return trace.Apply(dbErr)
	}

	if migrationErr := db.Migrate(); migrationErr != nil {
		return trace.Apply(migrationErr)
	}

	UrlCache, urlCacheErr := url.InitUrlsCache(db)
	if urlCacheErr != nil {
		return trace.Apply(urlCacheErr)
	}

	return webServer(user, pass, db, UrlCache)
}

func webServer(user, pass string, db *db.Database, urlCache *url.UrlCache) error {
	trace := util.CreateErrorContext("webServer")

	base, baseErr := util.EnvAsResult("HTTP_BASE")
	port, portErr := util.EnvAsResult("HTTP_PORT")
	timeoutTime, timeoutTimeErr := util.EnvAsIntegerResult("TIMEOUT_TIME")

	if err := trace.Join(portErr, baseErr, timeoutTimeErr); err != nil {
		return err
	}

	router := mux.NewRouter().UseEncodedPath()
	sub := router.PathPrefix(base).Subrouter()

	sub.HandleFunc("/urls", middlewares.Auth(user, pass, url.GetUrls(urlCache))).Methods("GET")
	sub.HandleFunc("/urls", middlewares.Auth(user, pass, url.CreateUrl(urlCache, db))).Methods("POST")
	sub.HandleFunc("/urls/{id}", middlewares.Auth(user, pass, url.GetUrl(urlCache))).Methods("GET")
	sub.HandleFunc("/{id}", middlewares.Auth(user, pass, url.GetOriginalUrl(urlCache, db))).Methods("GET")
	sub.HandleFunc("/{id}", middlewares.Auth(user, pass, url.DeleteUrl(urlCache, db))).Methods("DELETE")

	fmt.Println("Starting ENCURTADOR at port", port)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: time.Duration(timeoutTime) * time.Second,
		Handler:           router,
	}

	return trace.Apply(server.ListenAndServe())
}
