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
		respond(w, http.StatusMethodNotAllowed, errorMessage("only POST or GET requests are supported"))
		return
	}

	req, err := parse(r)
	if err != nil {
		respond(w, http.StatusBadRequest, errorMessage(err.Error()))
		return
	}

	n := len(req.queries)
	if n == 0 {
		respond(w, http.StatusBadRequest, errorMessage("no queries to execute"))
		return
	}

	// TODO: User authentication should happen here, if needed.

	// Execute the request.
	var (
		ctx       = h.Loaders.Attach(r.Context())
		responses = make([]*graphql.Response, n)
		wg        sync.WaitGroup
	)

	wg.Add(n)

	for i, q := range req.queries {
		go func(i int, q Query) {
			r := h.Schema.Exec(ctx, q.Query, q.OpName, q.Variables)
			responses[i] = r
			wg.Done()
		}(i, q)
	}

	wg.Wait()

	// Build the response.
	// TODO: Massage errors before returning to API consumers.
	var b []byte
	if req.isBatch {
		b, err = json.Marshal(responses)
	} else {
		if len(responses) > 0 {
			b, err = json.Marshal(responses[0])
		}
	}

	if err != nil {
		respond(w, http.StatusInternalServerError, errorMessage("server error"))
		return
	}

	respond(w, http.StatusOK, b)
}

func respond(w http.ResponseWriter, status int, body []byte) {
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func isSupported(method string) bool {
	return method == "POST" || method == "GET"
}

func errorMessage(msg string) []byte {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, `{"error": "%s"}`, msg)
	return buf.Bytes()
}
