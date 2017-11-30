package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"

	"github.com/tonyghita/graphql-go-example/errors"
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

func LoadPlanet(ctx context.Context, url string) (swapi.Planet, error) {
	var planet swapi.Planet

	l, err := Extract(ctx, PlanetsByURLs)
	if err != nil {
		return planet, err
	}

	data, err := l.Load(ctx, url)()
	if err != nil {
		return planet, err
	}

	planet, ok := data.(swapi.Planet)
	if !ok {
		return planet, errors.UnexpectedResponse
	}

	return planet, nil
}

func LoadPlanets(ctx context.Context, urls []string) ([]swapi.Planet, error) {
	// TODO: implement
	return []swapi.Planet{}, nil
}

func PrimePlanets(ctx context.Context, page swapi.PlanetPage) error {
	l, err := Extract(ctx, PlanetsByURLs)
	if err != nil {
		return err
	}

	for _, planet := range page.Planets {
		l.Prime(planet.URL, planet)
	}
	return nil
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
