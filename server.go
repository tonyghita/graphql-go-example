// Command graphql-go-example starts an HTTP GraphQL API server which is backed by data
// against the https://swapi.co REST API.
package main

import (
	"log"
	"net/http"
	"time"

	graphql "github.com/graph-gophers/graphql-go"

	"github.com/tonyghita/graphql-go-example/handler"
	"github.com/tonyghita/graphql-go-example/loader"
	"github.com/tonyghita/graphql-go-example/resolver"
	"github.com/tonyghita/graphql-go-example/schema"
	"github.com/tonyghita/graphql-go-example/swapi"
)

func main() {
	// Tweak configuration values here.
	var (
		addr              = ":8000"
		readHeaderTimeout = 1 * time.Second
		writeTimeout      = 10 * time.Second
		idleTimeout       = 90 * time.Second
		maxHeaderBytes    = http.DefaultMaxHeaderBytes
	)

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	c := swapi.NewClient(http.DefaultClient) // TODO: don't use the default client.

	root, err := resolver.NewRoot(c)
	if err != nil {
		log.Fatal(err)
	}

	// Create the request handler; inject dependencies.
	h := handler.GraphQL{
		// Parse and validate schema. Panic if unable to do so.
		Schema:  graphql.MustParseSchema(schema.String(), root),
		Loaders: loader.Initialize(c),
	}

	// Register handlers to routes.
	mux := http.NewServeMux()
	mux.Handle("/", handler.GraphiQL{})
	mux.Handle("/graphql/", h)
	mux.Handle("/graphql", h) // Register without a trailing slash to avoid redirect.

	// Configure the HTTP server.
	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	// Begin listeing for requests.
	log.Printf("Listening for requests on %s", s.Addr)

	if err = s.ListenAndServe(); err != nil {
		log.Println("server.ListenAndServe:", err)
	}

	// TODO: intercept shutdown signals for cleanup of connections.
	log.Println("Shut down.")
}
