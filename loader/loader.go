package loader

import (
	"context"
	"fmt"

	"github.com/nicksrandall/dataloader"

	"github.com/tonyghita/graphql-go-example/swapi"
)

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

func Initialize(ctx context.Context, client *swapi.Client) context.Context {
	m := map[key]dataloader.BatchFunc{
		FilmsByURLs:     NewFilmLoader(client),
		PeopleByURLs:    NewPersonLoader(client),
		PlanetsByURLs:   NewPlanetLoader(client),
		SpeciesByURLs:   NewSpeciesLoader(client),
		StarshipsByURLs: NewStarshipLoader(client),
		VehiclesByURLs:  NewVehicleLoader(client),
	}

	for k, batchFunc := range m {
		ctx = context.WithValue(ctx, k, dataloader.NewBatchedLoader(batchFunc))
	}

	return ctx
}

func Extract(ctx context.Context, k key) (*dataloader.Loader, error) {
	l, ok := ctx.Value(k).(*dataloader.Loader)
	if !ok {
		return nil, fmt.Errorf("unable to find %q loader on the request context", k)
	}

	return l, nil
}

func (k key) String() string {
	s, ok := keysToStr[k]
	if !ok {
		return "unknown"
	}
	return s
}
