package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	redis "gopkg.in/redis.v5"

	"github.com/info344-s17/challenges-leedann/apiserver/handlers"
	"github.com/info344-s17/challenges-leedann/apiserver/middleware"
	"github.com/info344-s17/challenges-leedann/apiserver/models/users"
	"github.com/info344-s17/challenges-leedann/apiserver/sessions"
	_ "github.com/lib/pq"
)

const defaultPort = "443"

const (
	apiRoot    = "/v1/"
	apiSummary = apiRoot + "summary"
	pgPort     = 5432
	usr        = "users"
	sess       = "sessions"
	sessme     = "sessions/mine"
	usrme      = "users/me"
)

//main is the main entry point for this program
func main() {
	//read and use the following environment variables
	//when initializing and starting your web server
	// PORT - port number to listen on for HTTP requests (if not set, use defaultPort)
	// HOST - host address to respond to (if not set, leave empty, which means any host)

	PORT := os.Getenv("PORT")
	if len(PORT) == 0 {
		PORT = defaultPort
	}
	HOST := os.Getenv("HOST")

	certPath := os.Getenv("TLSCERT")
	keyPath := os.Getenv("TLSKEY")

	SESSIONKEY := os.Getenv("SESSIONKEY")
	REDISADDR := os.Getenv("REDISADDR")
	DBADDR := os.Getenv("DBADDR")

	client := redis.NewClient(&redis.Options{
		Addr:     REDISADDR,
		Password: "",
		DB:       0,
	})

	datasrcName := fmt.Sprintf("user=pgstest dbname=pgstest sslmode=disable host=%s port=%d", DBADDR, pgPort)
	pgstore, err := sql.Open("postgres", datasrcName)
	if err != nil {
		log.Fatalf("error starting db: %v", err)
	}
	store := &users.PGStore{
		DB: pgstore,
	}
	//Pings the DB-- establishes a connection to the db
	err = pgstore.Ping()
	if err != nil {
		log.Fatalf("error pinging db %v", err)
	}
	redisStore := sessions.NewRedisStore(client, time.Hour*3600)

	ctx := &handlers.Context{
		SessionKey:   SESSIONKEY,
		SessionStore: redisStore,
		UserStore:    store,
	}
	mux := http.NewServeMux()
	mux.HandleFunc(apiRoot+usr, ctx.UserHandler)
	mux.HandleFunc(apiRoot+sess, ctx.SessionsHandler)
	mux.HandleFunc(apiRoot+sessme, ctx.SessionsMineHandler)
	mux.HandleFunc(apiRoot+usrme, ctx.UsersMeHandler)
	mux.HandleFunc(apiSummary, handlers.SummaryHandler)
	mux.Handle(apiRoot, middleware.Adapt(mux, middleware.CORS("", "", "", "")))

	//add your handlers.SummaryHandler function as a handler
	//for the apiSummary route
	//HINT: https://golang.org/pkg/net/http/#HandleFunc

	//start your web server and use log.Fatal() to log
	//any errors that occur if the server can't start
	//HINT: https://golang.org/pkg/net/http/#ListenAndServe
	addr := HOST + ":" + PORT
	fmt.Printf("listening at %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, certPath, keyPath, mux))

}
