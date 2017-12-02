package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"
)

type planetGetter interface {
	Planet(ctx context.Context, url string) (swapi.Planet, error)
}

type planetLoader struct {
	get planetGetter
}

func newPlanetLoader(client planetGetter) dataloader.BatchFunc {
	return planetLoader{get: client}.loadBatch
}

func LoadPlanet(ctx context.Context, url string) (swapi.Planet, error) {
	var planet swapi.Planet

	ldr, err := extract(ctx, planetLoaderKey)
	if err != nil {
		return planet, err
	}

	data, err := ldr.Load(ctx, url)()
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
	ldr, err := extract(ctx, planetLoaderKey)
	if err != nil {
		return []swapi.Planet{}, err
	}

	data, loadErrs := ldr.LoadMany(ctx, urls)()

	var (
		planets = make([]swapi.Planet, 0, len(data))
		errs    = make(errors.Errors, 0, len(loadErrs))
	)

	for i := range urls {
		d, err := data[i], loadErrs[i]
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		planet, ok := d.(swapi.Planet)
		if !ok && err == nil {
			errs = append(errs, errors.WithIndex(errors.UnexpectedResponse, i))
		}

		planets = append(planets, planet)
	}

	return planets, errs.Err()
}

func PrimePlanets(ctx context.Context, page swapi.PlanetPage) error {
	ldr, err := extract(ctx, planetLoaderKey)
	if err != nil {
		return err
	}

	for _, p := range page.Planets {
		ldr.Prime(p.URL, p)
	}
	return nil
}

func (ldr planetLoader) loadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(ctx context.Context, i int, url string) {
			data, err := ldr.get.Planet(ctx, url)
			results[i] = &dataloader.Result{Data: data, Error: err}
			wg.Done()
		}(ctx, i, url)
	}

	wg.Wait()

	return results
}
