package swapi_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/tonyghita/graphql-go-example/swapi"
)

func TestClient(t *testing.T) {
	t.Run("NewClient", func(t *testing.T) {
		client := swapi.NewClient(http.DefaultClient)
		if client == nil {
			t.Error("swapi.NewClient(http.DefaultClient): expected non-null result")
		}
	})

	t.Run("NewRequest", func(t *testing.T) {
		client := swapi.NewClient(http.DefaultClient)
		if client == nil {
			t.Error("swapi.NewClient(http.DefaultClient): expected non-null result")
		}

		cases := []struct {
			input string
		}{
			{"https://swapi.dev/api/planets/1/"},
			{"/planets/1/"},
		}

		ctx := context.Background()
		for _, c := range cases {
			r, err := client.NewRequest(ctx, c.input)
			if err != nil {
				t.Fatalf("client.NewRequest(ctx, %q): unexpected error %v", c.input, err)
			}

			expected := "https://swapi.dev/api/planets/1/"
			actual := r.URL.String()
			if actual != expected {
				t.Fatalf("client.NewRequest(ctx, %q): wanted %q, got %q", c.input, expected, actual)
			}
		}
	})
}
