package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"
)

func LoadFilm(ctx context.Context, url string) (swapi.Film, error) {
	var film swapi.Film

	ldr, err := extract(ctx, filmLoaderKey)
	if err != nil {
		return film, err
	}

	data, err := ldr.Load(ctx, url)()
	if err != nil {
		return film, err
	}

	film, ok := data.(swapi.Film)
	if !ok {
		return film, errors.UnexpectedResponse
	}

	return film, nil
}

func LoadFilms(ctx context.Context, urls []string) ([]swapi.Film, error) {
	ldr, err := extract(ctx, filmLoaderKey)
	if err != nil {
		return []swapi.Film{}, err
	}

	data, loadErrs := ldr.LoadMany(ctx, urls)()

	var (
		films = make([]swapi.Film, 0, len(data))
		errs  = make(errors.Errors, 0, len(loadErrs))
	)

	for i := range urls {
		d, err := data[i], loadErrs[i]
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		film, ok := d.(swapi.Film)
		if !ok && err == nil {
			errs = append(errs, errors.WithIndex(errors.UnexpectedResponse, i))
		}

		films = append(films, film)
	}

	return films, errs.Err()
}

// PrimeFilms ...
func PrimeFilms(ctx context.Context, page swapi.FilmPage) error {
	ldr, err := extract(ctx, filmLoaderKey)
	if err != nil {
		return err
	}

	for _, f := range page.Films {
		ldr.Prime(f.URL, f)
	}

	return nil
}

type filmGetter interface {
	Film(ctx context.Context, url string) (swapi.Film, error)
}

// FilmLoader contains the client required to load film resources.
type filmLoader struct {
	get filmGetter
}

func newFilmLoader(client filmGetter) dataloader.BatchFunc {
	return filmLoader{get: client}.loadBatch
}

func (ldr filmLoader) loadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url string) {
			resp, err := ldr.get.Film(ctx, url)
			results[i] = &dataloader.Result{Data: resp, Error: err}
			wg.Done()
		}(i, url)
	}

	wg.Wait()

	return results
}
