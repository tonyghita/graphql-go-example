package schema_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonyghita/graphql-go-example/schema"
)

func TestString(t *testing.T) {
	s, err := schema.String()

	require.NoError(t, err)
	require.NotEmpty(t, s)
}
