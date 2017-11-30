package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"
)

type SpeciesGetter interface {
	Species(ctx context.Context, url string) (swapi.Species, error)
}

type SpeciesLoader struct {
	get SpeciesGetter
}

func NewSpeciesLoader(client SpeciesGetter) dataloader.BatchFunc {
	return SpeciesLoader{get: client}.loadBatch
}

func (loader SpeciesLoader) loadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url string) {
			sp, err := loader.get.Species(ctx, url)
			results[i] = &dataloader.Result{Data: sp, Error: err}
			wg.Done()
		}(i, url)
	}

	wg.Wait()

	return results
}

func LoadSpecies(ctx context.Context, url string) (swapi.Species, error) {
	var species swapi.Species
	l, err := Extract(ctx, SpeciesByURLs)
	if err != nil {
		return species, err
	}

	data, err := l.Load(ctx, url)()
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
	l, err := Extract(ctx, SpeciesByURLs)
	if err != nil {
		return []swapi.Species{}, err
	}

	data, loadErrors := l.LoadMany(ctx, urls)()
	species := make([]swapi.Species, len(urls))
	errs := make(errors.Errors, 0, len(loadErrors))

	for i := range urls {
		d, err := data[i], loadErrors[i]
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		sp, ok := d.(swapi.Species)
		if !ok && err == nil {
			// Ensure one error per element in the list.
			errs = append(errs, errors.WithIndex(err, i))
		}

		species = append(species, sp)
	}

	return species, errs.Err()
}

func PrimeSpecies(ctx context.Context, page swapi.SpeciesPage) error {
	l, err := Extract(ctx, SpeciesByURLs)
	if err != nil {
		return err
	}

	for _, sp := range page.Species {
		l.Prime(sp.URL, sp)
	}

	return nil
}
