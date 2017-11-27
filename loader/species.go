package loader

import (
	"context"
	"sync"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"

	"github.com/nicksrandall/dataloader"
)

// SpeciesLoader ...
type SpeciesLoader struct {
	client *swapi.Client
}

// LoadBatch ...
func (l SpeciesLoader) LoadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url string) {
			sp, err := l.client.Species(ctx, url)
			results[i] = &dataloader.Result{Data: sp, Error: err}
			wg.Done()
		}(i, url)
	}

	wg.Wait()

	return results
}

// LoadSpecies ...
func LoadSpecies(ctx context.Context, url string) (swapi.Species, error) {
	l, err := Extract(ctx, SpeciesByURLs)
	if err != nil {
		return swapi.Species{}, err
	}

	loadFn := l.Load(ctx, url)
	data, err := loadFn()
	if err != nil {
		return swapi.Species{}, err
	}

	sp, ok := data.(swapi.Species)
	if !ok {
		return swapi.Species{}, errors.New("unexpected response")
	}

	return sp, nil
}

// LoadManySpecies ...
func LoadManySpecies(ctx context.Context, urls ...string) ([]swapi.Species, error) {
	l, err := Extract(ctx, SpeciesByURLs)
	if err != nil {
		return []swapi.Species{}, err
	}

	loadFn := l.LoadMany(ctx, urls)
	data, loadErrors := loadFn()

	var species = make([]swapi.Species, len(urls))
	var errs errors.Errors

	// TODO: Construct a map of url -> data and range over URL data.
	for i := range urls {
		d, err := data[i], loadErrors[i]
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}
		species = append(species, d.(swapi.Species))
	}

	return species, nil
}

// PrimeSpecies ...
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
