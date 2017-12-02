package resolver

import (
	"strconv"
	"strings"

	graphql "github.com/neelance/graphql-go"
)

func extractID(url string) graphql.ID {
	parts := strings.Split(url, "/")
	n := len(parts)

	if n < 5 {
		// A full URL must have at least 5 '/'.
		// http://swapi.co/api/{resource}/{id}/
		//      xx        x   x          x = 5
		return graphql.ID("")
	}

	id := parts[n-1]
	if id == "" {
		id = parts[n-2]
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return graphql.ID("")
	}

	return graphql.ID(id)
}

func strValue(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return *ptr
}

func nullableStr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}
