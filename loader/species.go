package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"
)

type speciesGetter interface {
	Species(ctx context.Context, url string) (swapi.Species, error)
}

type speciesLoader struct {
	get speciesGetter
}

func newSpeciesLoader(client speciesGetter) dataloader.BatchFunc {
	return speciesLoader{get: client}.loadBatch
}

func LoadSpecies(ctx context.Context, url string) (swapi.Species, error) {
	var species swapi.Species

	ldr, err := extract(ctx, speciesLoaderKey)
	if err != nil {
		return species, err
	}

	data, err := ldr.Load(ctx, url)()
	if err != nil {
		return species, err
	}

	species, ok := data.(swapi.Species)
	if !ok {
		return species, errors.UnexpectedResponse
	}

	return species, nil
}

func LoadManySpecies(ctx context.Context, urls ...string) ([]swapi.Species, error) {
	ldr, err := extract(ctx, speciesLoaderKey)
	if err != nil {
		return []swapi.Species{}, err
	}

	data, loadErrs := ldr.LoadMany(ctx, urls)()

	var (
		species = make([]swapi.Species, 0, len(urls))
		errs    = make(errors.Errors, 0, len(loadErrs))
	)

	for i := range urls {
		d, err := data[i], loadErrs[i]
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		sp, ok := d.(swapi.Species)
		if !ok && err == nil {
			// Ensure one error per element in the list.
			errs = append(errs, errors.WithIndex(errors.UnexpectedResponse, i))
		}

		species = append(species, sp)
	}

	return species, errs.Err()
}

func PrimeSpecies(ctx context.Context, page swapi.SpeciesPage) error {
	ldr, err := extract(ctx, speciesLoaderKey)
	if err != nil {
		return err
	}

	for _, s := range page.Species {
		ldr.Prime(s.URL, s)
	}

	return nil
}

func (ldr speciesLoader) loadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url string) {
			sp, err := ldr.get.Species(ctx, url)
			results[i] = &dataloader.Result{Data: sp, Error: err}
			wg.Done()
		}(i, url)
	}

	wg.Wait()

	return results
}
