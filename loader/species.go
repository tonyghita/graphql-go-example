package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"
)

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

func LoadManySpecies(ctx context.Context, urls ...string) (SpeciesResults, error) {
	var results []SpeciesResult

	ldr, err := extract(ctx, speciesLoaderKey)
	if err != nil {
		return results, err
	}

	data, errs := ldr.LoadMany(ctx, copyStrings(urls))()
	results = make([]SpeciesResult, 0, len(urls))

	for i, d := range data {
		var e error
		if errs != nil {
			e = errs[i]
		}

		species, ok := d.(swapi.Species)
		if !ok && e == nil {
			err = errors.UnexpectedResponse
		}

		results = append(results, SpeciesResult{Species: species, Error: e})
	}

	return results, nil
}

type SpeciesResult struct {
	Species swapi.Species
	Error   error
}

type SpeciesResults []SpeciesResult

func (results SpeciesResults) WithoutErrors() []swapi.Species {
	species := make([]swapi.Species, 0, len(results))

	for _, r := range results {
		if r.Error != nil {
			continue
		}

		species = append(species, r.Species)
	}

	return species
}

func PrimeSpecies(ctx context.Context, page swapi.SpeciesPage) error {
	ldr, err := extract(ctx, speciesLoaderKey)
	if err != nil {
		return err
	}

	for _, s := range page.Species {
		ldr.Prime(ctx, s.URL, s)
	}

	return nil
}

type speciesGetter interface {
	Species(ctx context.Context, url string) (swapi.Species, error)
}

type speciesLoader struct {
	get speciesGetter
}

func newSpeciesLoader(client speciesGetter) dataloader.BatchFunc {
	return speciesLoader{get: client}.loadBatch
}

func (ldr speciesLoader) loadBatch(ctx context.Context, urls []interface{}) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, value := range urls {
		go func(i int, v interface{}) {
			defer wg.Done()

			url, ok := v.(string)
			if !ok {
				results[i] = &dataloader.Result{Error: errors.WrongKeyType(url, v)}
				return
			}

			data, err := ldr.get.Species(ctx, url)
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, value)
	}

	wg.Wait()

	return results
}
