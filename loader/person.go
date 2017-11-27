package loader

import (
	"context"
	"sync"

	"github.com/tonyghita/graphql-go-example/swapi"

	"github.com/nicksrandall/dataloader"
)

// PersonLoader contains the RPC client necessary to load people.
type PersonLoader struct {
	client *swapi.Client
}

// LoadBatch ...
func (l PersonLoader) LoadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(ctx context.Context, url string, i int) {
			data, err := l.client.Person(ctx, url)
			results[i] = &dataloader.Result{Data: data, Error: err}
			wg.Done()
		}(ctx, url, i)
	}

	wg.Wait()
	return results
}

// LoadPerson loads a person resource from the SWAPI API URL.
func LoadPerson(ctx context.Context, url string) (swapi.Person, error) {
	l, err := Extract(ctx, PeopleByURLs)
	if err != nil {
		return swapi.Person{}, err
	}

	loadFn := l.Load(ctx, url)
	data, err := loadFn()
	if err != nil {
		return swapi.Person{}, err
	}

	return data.(swapi.Person), nil
}
