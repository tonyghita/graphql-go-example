package loader

import (
	"context"
	"sync"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"

	"github.com/nicksrandall/dataloader"
)

// FilmLoader contains the client required to load film resources.
type FilmLoader struct {
	client *swapi.Client
}

// LoadBatch ...
func (l FilmLoader) LoadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url string) {
			resp, err := l.client.Film(ctx, url)
			results[i] = &dataloader.Result{Data: resp, Error: err}
			wg.Done()
		}(i, url)
	}

	wg.Wait()

	return results
}

// LoadFilm ...
func LoadFilm(ctx context.Context, url string) (swapi.Film, error) {
	l, err := Extract(ctx, FilmsByURLs)
	if err != nil {
		return swapi.Film{}, err
	}

	loadFn := l.Load(ctx, url)
	data, err := loadFn()
	if err != nil {
		return swapi.Film{}, err
	}

	film, ok := data.(swapi.Film)
	if !ok {
		return swapi.Film{}, errors.UnexpectedResponse
	}

	return film, nil
}

// LoadFilms ...
func LoadFilms(ctx context.Context, urls ...string) ([]swapi.Film, error) {
	l, err := Extract(ctx, FilmsByURLs)
	if err != nil {
		return []swapi.Film{}, err
	}

	loadFn := l.LoadMany(ctx, urls)
	data, _ := loadFn() // TODO: Use these errors instead.

	var films = make([]swapi.Film, len(data))
	var errs errors.Errors

	// TODO: range over URLs instead of data.
	for i, d := range data {
		film, ok := d.(swapi.Film)
		if !ok {
			errs = append(errs, errors.WithIndex(errors.UnexpectedResponse, i))
		}

		films[i] = film
	}

	return films, errs.Err()
}

// PrimeFilms ...
func PrimeFilms(ctx context.Context, page swapi.FilmPage) error {
	l, err := Extract(ctx, FilmsByURLs)
	if err != nil {
		return err
	}

	for _, f := range page.Films {
		l.Prime(f.URL, f)
	}
	return nil
}
