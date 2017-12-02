package loader

import (
	"context"
	"sync"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"

	"github.com/nicksrandall/dataloader"
)

type starshipGetter interface {
	Starship(ctx context.Context, url string) (swapi.Starship, error)
}

type StarshipLoader struct {
	get starshipGetter
}

func newStarshipLoader(client starshipGetter) dataloader.BatchFunc {
	return StarshipLoader{get: client}.loadBatch
}

// LoadStarship ...
func LoadStarship(ctx context.Context, url string) (swapi.Starship, error) {
	ldr, err := extract(ctx, starshipLoaderKey)
	if err != nil {
		return swapi.Starship{}, err
	}

	data, err := ldr.Load(ctx, url)()
	if err != nil {
		return swapi.Starship{}, err
	}

	ship, ok := data.(swapi.Starship)
	if !ok {
		return swapi.Starship{}, errors.UnexpectedResponse
	}

	return ship, nil
}

func LoadStarships(ctx context.Context, urls []string) ([]swapi.Starship, error) {
	ldr, err := extract(ctx, starshipLoaderKey)
	if err != nil {
		return []swapi.Starship{}, err
	}

	data, loadErrs := ldr.LoadMany(ctx, urls)()

	var (
		ships = make([]swapi.Starship, 0, len(data))
		errs  = make(errors.Errors, 0, len(loadErrs))
	)

	for i := range urls {
		d, err := data[i], loadErrs[i]
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		ship, ok := d.(swapi.Starship)
		if !ok && err == nil {
			errs = append(errs, errors.WithIndex(errors.UnexpectedResponse, i))
		}

		ships = append(ships, ship)
	}

	return ships, errs.Err()
}

func PrimeStarships(ctx context.Context, page swapi.StarshipPage) error {
	ldr, err := extract(ctx, starshipLoaderKey)
	if err != nil {
		return err
	}

	for _, s := range page.Starships {
		ldr.Prime(s.URL, s)
	}

	return nil
}

func (ldr StarshipLoader) loadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(ctx context.Context, url string, i int) {
			data, err := ldr.get.Starship(ctx, url)
			results[i] = &dataloader.Result{Data: data, Error: err}
			wg.Done()
		}(ctx, url, i)
	}

	wg.Wait()

	return results
}
