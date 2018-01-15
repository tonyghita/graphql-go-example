package handler_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tonyghita/graphql-go-example/handler"
)

func TestGraphiQLHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method string
		status int
	}{
		{http.MethodGet, http.StatusOK},
		{http.MethodDelete, http.StatusMethodNotAllowed},
		{http.MethodPost, http.StatusMethodNotAllowed},
		{http.MethodPut, http.StatusMethodNotAllowed},
		{http.MethodPatch, http.StatusMethodNotAllowed},
	}

	graphiql := handler.GraphiQL{}
	ts := httptest.NewServer(graphiql)
	defer ts.Close()

	for _, test := range tests {
		test := test

		req, err := http.NewRequest(test.method, ts.URL, nil)
		fatalIfErr(t, err)

		resp, err := ts.Client().Do(req)
		fatalIfErr(t, err)

		if resp.StatusCode != test.status {
			t.Errorf("wrong response code: expected %q, got %q", test.status, resp.StatusCode)
		}

		b, err := ioutil.ReadAll(resp.Body)
		fatalIfErr(t, err)

		if len(b) == 0 {
			t.Error("expected response to have a body")
		}
	}
}

func fatalIfErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
