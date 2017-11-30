package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"
)

type PersonGetter interface {
	Person(ctx context.Context, url string) (swapi.Person, error)
}

// PersonLoader contains the RPC client necessary to load people.
type PersonLoader struct {
	get PersonGetter
}

func NewPersonLoader(client PersonGetter) dataloader.BatchFunc {
	return PersonLoader{get: client}.loadBatch
}

// LoadPerson loads a person resource from the SWAPI API URL.
func LoadPerson(ctx context.Context, url string) (swapi.Person, error) {
	var person swapi.Person

	l, err := Extract(ctx, PeopleByURLs)
	if err != nil {
		return person, err
	}

	data, err := l.Load(ctx, url)()
	if err != nil {
		return person, err
	}

	person, ok := data.(swapi.Person)
	if !ok {
		return person, errors.UnexpectedResponse
	}

	return person, nil
}

func (l PersonLoader) loadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(ctx context.Context, url string, i int) {
			data, err := l.get.Person(ctx, url)
			results[i] = &dataloader.Result{Data: data, Error: err}
			wg.Done()
		}(ctx, url, i)
	}

	wg.Wait()

	return results
}
