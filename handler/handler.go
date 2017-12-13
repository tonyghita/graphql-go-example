// Package handler contains the type which knows how to parse GraphQL HTTP requests.
package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	graphql "github.com/neelance/graphql-go"

	"github.com/tonyghita/graphql-go-example/loader"
)

// Logger defines an interface with a single method.
type Logger interface {
	Printf(fmt string, values ...interface{})
}

// A Request respresents an HTTP request to the GraphQL endpoint.
// A request can have a single query or a batch of requests with one or more queries.
// It is important to distinguish between a single query request and a batch request with a single query.
// The shape of the response will differ in both cases.
type Request struct {
	queries []Query
	isBatch bool
}

// A Query represents a single GraphQL query.
type Query struct {
	OpName    string                 `json:"operationName"`
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// The GraphQL handler handles GraphQL API requests over HTTP.
type GraphQL struct {
	Schema  *graphql.Schema
	Loaders loader.Collection
	Logger  Logger
}

func (h GraphQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate the request.
	if ok := isSupported(r.Method); !ok {
		respond(w, errorJSON("only POST or GET requests are supported"), http.StatusMethodNotAllowed)
		return
	}

	req, err := parse(r)
	if err != nil {
		respond(w, errorJSON(err.Error()), http.StatusBadRequest)
		return
	}

	n := len(req.queries)
	if n == 0 {
		respond(w, errorJSON("no queries to execute"), http.StatusBadRequest)
		return
	}

	// NOTE: User authentication should happen here, if needed.
	//
	// Authentication determines who the request originated from.
	// Authorization business logic (in services, not this API) will rely on this authentication data.
	//
	// We don't need it for this example application, but a typical production application would
	// perform request authentication.
	//
	// The result of authentication should probably be placed on the request context so it can be
	// passed to resolvers and loaders.

	// Execute the request.
	var (
		ctx       = h.Loaders.Attach(r.Context())
		responses = make([]*graphql.Response, n)
		wg        sync.WaitGroup
	)

	wg.Add(n)

	for i, q := range req.queries {
		go func(i int, q Query) {
			res := h.Schema.Exec(ctx, q.Query, q.OpName, q.Variables)
			// TODO: "Massage" response errors here.
			responses[i] = res
			wg.Done()
		}(i, q)
	}

	wg.Wait()

	// Marshal the response to JSON.
	var resp []byte
	if req.isBatch {
		resp, err = json.Marshal(responses)
	} else if len(responses) > 0 {
		resp, err = json.Marshal(responses[0])
	}

	if err != nil {
		respond(w, errorJSON("server error"), http.StatusInternalServerError)
		return
	}

	respond(w, resp, http.StatusOK)
}

func respond(w http.ResponseWriter, body []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, body)
}

func isSupported(method string) bool {
	return method == "POST" || method == "GET"
}

func errorJSON(msg string) []byte {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, `{"error": "%s"}`, msg)
	return buf.Bytes()
}
