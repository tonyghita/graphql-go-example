package loader

import (
	"context"
	"sync"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"

	"github.com/nicksrandall/dataloader"
)

type StarshipGetter interface {
	Starship(ctx context.Context, url string) (swapi.Starship, error)
}

type StarshipLoader struct {
	get StarshipGetter
}

func NewStarshipLoader(client StarshipGetter) dataloader.BatchFunc {
	return StarshipLoader{get: client}.loadBatch
}

// LoadStarship ...
func LoadStarship(ctx context.Context, url string) (swapi.Starship, error) {
	l, err := Extract(ctx, StarshipsByURLs)
	if err != nil {
		return swapi.Starship{}, err
	}

	loadFn := l.Load(ctx, url)
	data, err := loadFn()
	if err != nil {
		return swapi.Starship{}, err
	}

	ship, ok := data.(swapi.Starship)
	if !ok {
		return swapi.Starship{}, errors.UnexpectedResponse
	}

	return ship, nil
}

// LoadBatch ...
func (loader StarshipLoader) loadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(ctx context.Context, url string, i int) {
			data, err := loader.get.Starship(ctx, url)
			results[i] = &dataloader.Result{Data: data, Error: err}
			wg.Done()
		}(ctx, url, i)
	}

	wg.Wait()

	return results
}
