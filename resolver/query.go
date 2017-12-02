package resolver

import (
	"context"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"
)

// The QueryResolver is the entry point for all top-level read operations.
type QueryResolver struct {
	client *swapi.Client
}

func NewRoot(client *swapi.Client) (*QueryResolver, error) {
	if client == nil {
		return nil, errors.UnableToResolve
	}

	return &QueryResolver{client: client}, nil
}

// FilmsQueryArgs are the arguments for the "films" query.
type FilmsQueryArgs struct {
	// Title of the film. When nil, all films are fetched.
	Title *string
}

// Films resolves a list of films. If no arguments are provided, all films are fetched.
func (r QueryResolver) Films(ctx context.Context, args FilmsQueryArgs) (*[]*FilmResolver, error) {
	page, err := r.client.SearchFilms(ctx, strValue(args.Title))
	if err != nil {
		return nil, err
	}

	return NewFilms(ctx, NewFilmsArgs{Page: page})
}

// PeopleQueryArgs are the arguments for the "people" query.
type PeopleQueryArgs struct {
	// Name of the person. When nil, all people are fetched.
	Name *string
}

// People resolves a list of people. If no arguments are provided, all people are fetched.
func (r QueryResolver) People(ctx context.Context, args PeopleQueryArgs) (*[]*PersonResolver, error) {
	page, err := r.client.SearchPerson(ctx, strValue(args.Name))
	if err != nil {
		return nil, err
	}

	return NewPeople(ctx, NewPeopleArgs{Page: page})
}

// PlanetsQueryArgs are the arguments for the "planets" query.
type PlanetsQueryArgs struct {
	// Name of the planet. When nil, all planets are fetched.
	Name *string
}

// Planets resolves a list of planets. If no arguments are provided, all planets are fetched.
func (r QueryResolver) Planets(ctx context.Context, args PlanetsQueryArgs) (*[]*PlanetResolver, error) {
	page, err := r.client.SearchPlanets(ctx, strValue(args.Name))
	if err != nil {
		return nil, err
	}

	return NewPlanets(ctx, NewPlanetsArgs{Page: page})
}

// SpeciesQueryArgs are the arguments for the "species" query.
type SpeciesQueryArgs struct {
	// Name of the species. When nil, all planets are fetched.
	Name *string
}

// Species resolves a list of species. If no arguments are provided, all species are fetched.
func (r QueryResolver) Species(ctx context.Context, args SpeciesQueryArgs) (*[]*SpeciesResolver, error) {
	page, err := r.client.SearchSpecies(ctx, strValue(args.Name))
	if err != nil {
		return nil, err
	}

	return NewSpeciesList(ctx, NewSpeciesListArgs{Page: page})
}

type StarshipsQueryArgs struct {
	NameOrModel *string
}

func (r QueryResolver) Starships(ctx context.Context, args StarshipsQueryArgs) (*[]*StarshipResolver, error) {
	page, err := r.client.SearchStarships(ctx, strValue(args.NameOrModel))
	if err != nil {
		return nil, err
	}

	return NewStarships(ctx, NewStarshipsArgs{Page: page})
}

type VehiclesQueryArgs struct {
	NameOrModel *string
}

func (r QueryResolver) Vehicles(ctx context.Context, args VehiclesQueryArgs) (*[]*VehicleResolver, error) {
	page, err := r.client.SearchVehicles(ctx, strValue(args.NameOrModel))
	if err != nil {
		return nil, err
	}

	return NewVehicles(ctx, NewVehiclesArgs{Page: page})
}
