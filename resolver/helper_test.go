package resolver

import (
	"testing"

	graphql "github.com/graph-gophers/graphql-go"
)

func TestParseID(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"https://swapi.co/api/planets/1/", "1"},
		{"https://swapi.co/api/planets/313/", "313"},
		{"https://swapi.co/api/people/4", "4"},
		{"https://swapi.co/api/vehicles/", ""},
		{"https://swapi.co/api/vehicles", ""},
		{"foo bar baz 1234", ""},
	}

	for _, c := range cases {
		expected, actual := graphql.ID(c.expected), extractID(c.input)
		if expected != actual {
			t.Errorf("parseID(%q): wanted %q, got %q", c.input, c.expected, actual)
		}
	}
}

var _benchParseIDResult graphql.ID

func BenchmarkParseID(b *testing.B) {
	url := "https://swapi.co/api/planets/313/"
	for n := 0; n < b.N; n++ {
		_benchParseIDResult = extractID(url)
	}
}
