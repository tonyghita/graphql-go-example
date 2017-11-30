package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"
	"github.com/tonyghita/graphql-go-example/swapi"
)

type PlanetGetter interface {
	Planet(ctx context.Context, url string) (swapi.Planet, error)
}

type PlanetLoader struct {
	get PlanetGetter
}

func NewPlanetLoader(client PlanetGetter) dataloader.BatchFunc {
	return PlanetLoader{get: client}.loadBatch
}

func (loader PlanetLoader) loadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(ctx context.Context, i int, url string) {
			data, err := loader.get.Planet(ctx, url)
			results[i] = &dataloader.Result{Data: data, Error: err}
			wg.Done()
		}(ctx, i, url)
	}

	wg.Wait()

	return results
}
