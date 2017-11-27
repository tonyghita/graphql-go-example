package loader

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nicksrandall/dataloader"
	"github.com/tonyghita/graphql-go-example/swapi"
)

// Loader ...
type Loader interface {
	LoadBatch(context.Context, []string) []*dataloader.Result
}

// TODO: describe why this type has been aliased to an unexported type.
type key int

const (
	FilmsByURLs key = iota
	PeopleByURLs
	PlanetsByURLs
	SpeciesByURLs
	StarshipsByURLs
	VehiclesByURLs
)

var keysToStr = map[key]string{
	FilmsByURLs:     "FilmsByURLs",
	PeopleByURLs:    "PeopleByURLs",
	PlanetsByURLs:   "PlanetsByURLs",
	SpeciesByURLs:   "SpeciesByURLs",
	StarshipsByURLs: "StarshipsByURLs",
	VehiclesByURLs:  "VehiclesByURLs",
}

// TODO: inject clients from somewhere else.
var keysToBatchFn = map[key]Loader{
	FilmsByURLs:     &FilmLoader{swapi.NewClient(http.DefaultClient)},
	PeopleByURLs:    &PersonLoader{swapi.NewClient(http.DefaultClient)},
	SpeciesByURLs:   &SpeciesLoader{swapi.NewClient(http.DefaultClient)},
	StarshipsByURLs: &StarshipLoader{swapi.NewClient(http.DefaultClient)},
}

// Extract ...
func Extract(ctx context.Context, k key) (*dataloader.Loader, error) {
	l, ok := ctx.Value(k).(*dataloader.Loader)
	if !ok {
		return nil, fmt.Errorf("unable to find %q loader on the request context", k)
	}

	return l, nil
}

// Initialize ...
// TODO: inject the swapi.Client instance here.
func Initialize(ctx context.Context) context.Context {
	for k, l := range keysToBatchFn {
		ctx = context.WithValue(ctx, k, dataloader.NewBatchedLoader(l.LoadBatch))
	}

	return ctx
}

// String ...
func (k key) String() string {
	s, ok := keysToStr[k]
	if !ok {
		return "unknown"
	}
	return s
}
