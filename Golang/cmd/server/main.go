package main

import (
	"fmt"
	"golang/db"
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

	return webServer(db, UrlCache)
}

func webServer(db *db.Database, urlCache *url.UrlCache) error {
	trace := util.CreateErrorContext("webServer")

	base, baseErr := util.EnvAsResult("HTTP_BASE")
	port, portErr := util.EnvAsResult("HTTP_PORT")
	timeoutTime, timeoutTimeErr := util.EnvAsIntegerResult("TIMEOUT_TIME")

	if err := trace.Join(portErr, baseErr, timeoutTimeErr); err != nil {
		return err
	}

	router := mux.NewRouter().UseEncodedPath()
	sub := router.PathPrefix(base).Subrouter()

	sub.HandleFunc("/urls", url.GetUrls(urlCache)).Methods("GET")
	sub.HandleFunc("/urls", url.CreateUrl(urlCache, db)).Methods("POST")
	sub.HandleFunc("/urls/{id}", url.GetUrl(urlCache)).Methods("GET")
	sub.HandleFunc("/{id}", url.GetOriginalUrl(urlCache, db)).Methods("GET")
	sub.HandleFunc("/{id}", url.DeleteUrl(urlCache, db)).Methods("DELETE")

	fmt.Println("Starting ENCURTADOR at port", port)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: time.Duration(timeoutTime) * time.Second,
		Handler:           router,
	}

	return trace.Apply(server.ListenAndServe())
}
