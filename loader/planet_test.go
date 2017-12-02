package loader_test

import (
	"context"
	"testing"

	"github.com/tonyghita/graphql-go-example/loader"
)

func TestPlanetLoader(t *testing.T) {
	t.Run("LoadPlanet", func(t *testing.T) {
		// TODO: implement.

		ctx, url := context.Background(), "https://swapi.co/api/planet/1"
		loader.LoadPlanet(ctx, url)
	})
}
