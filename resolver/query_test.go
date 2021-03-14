package resolver_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonyghita/graphql-go-example/resolver"
	"github.com/tonyghita/graphql-go-example/schema"

	graphql "github.com/graph-gophers/graphql-go"
)

func TestResolversSatisfySchema(t *testing.T) {
	s, err := schema.String()
	require.NoError(t, err)
	require.NotEmpty(t, s)

	rootResolver := &resolver.QueryResolver{}

	_, err = graphql.ParseSchema(s, rootResolver)
	require.NoError(t, err)
}
