package resolver

import (
	"context"

	"github.com/tonyghita/graphql-go-example/swapi"
)

// The QueryResolver is the entry point for all top-level read operations.
type QueryResolver struct {
	Client *swapi.Client
}

// FilmsQueryArgs are the arguments for the "films" query.
type FilmsQueryArgs struct {
	// Title of the film. When nil, all films are fetched.
	Title *string
}

// Films resolves a list of films. If no arguments are provided, all films are fetched.
func (r QueryResolver) Films(ctx context.Context, args FilmsQueryArgs) ([]*FilmResolver, error) {
	var title string
	if !nilOrEmpty(args.Title) {
		title = *args.Title
	}

	page, err := r.Client.SearchFilms(ctx, title)
	if err != nil {
		return []*FilmResolver{}, err
	}

	return NewFilms(ctx, NewFilmsArgs{Page: page})
}

// PeopleQueryArgs are the arguments for the "people" query.
type PeopleQueryArgs struct {
	// Name of the person. When nil, all people are fetched.
	Name *string
}

// People resolves a list of people. If no arguments are provided, all people are fetched.
func (r QueryResolver) People(ctx context.Context, args PeopleQueryArgs) ([]*PersonResolver, error) {
	var name string
	if !nilOrEmpty(args.Name) {
		name = *args.Name
	}

	page, err := r.Client.SearchPerson(ctx, name)
	if err != nil {
		return []*PersonResolver{}, err
	}

	return NewPeople(ctx, NewPeopleArgs{Page: page})
}

// PlanetsQueryArgs are the arguments for the "planets" query.
type PlanetsQueryArgs struct {
	// Name of the planet. When nil, all planets are fetched.
	Name *string
}

// Planets resolves a list of planets. If no arguments are provided, all planets are fetched.
func (r QueryResolver) Planets(ctx context.Context, args PlanetsQueryArgs) ([]*PlanetResolver, error) {
	var name string
	if !nilOrEmpty(args.Name) {
		name = *args.Name
	}

	page, err := r.Client.SearchPlanets(ctx, name)
	if err != nil {
		return []*PlanetResolver{}, err
	}

	return NewPlanets(ctx, NewPlanetsArgs{Page: page})
}

// SpeciesQueryArgs are the arguments for the "species" query.
type SpeciesQueryArgs struct {
	// Name of the species. When nil, all planets are fetched.
	Name *string
}

// Species resolves a list of species. If no arguments are provided, all species are fetched.
func (r QueryResolver) Species(ctx context.Context, args SpeciesQueryArgs) ([]*SpeciesResolver, error) {
	var name string
	if !nilOrEmpty(args.Name) {
		name = *args.Name
	}

	page, err := r.Client.SearchSpecies(ctx, name)
	if err != nil {
		return []*SpeciesResolver{}, err
	}

	return NewSpeciesList(ctx, NewSpeciesListArgs{Page: page})
}

type StarshipsQueryArgs struct {
	Name  *string
	Model *string
}

func (r QueryResolver) Starships(ctx context.Context, args StarshipsQueryArgs) ([]*StarshipResolver, error) {
	var name string
	if !nilOrEmpty(args.Name) {
		name = *args.Name
	}

	page, err := r.Client.SearchStarships(ctx, name)
	if err != nil {
		return []*StarshipResolver{}, err
	}

	return NewStarships(ctx, NewStarshipsArgs{Page: page})
}

type VehiclesQueryArgs struct {
	Name  *string
	Model *string
}

func (r QueryResolver) Vehicles(ctx context.Context, args VehiclesQueryArgs) ([]*VehicleResolver, error) {
	var name string
	if !nilOrEmpty(args.Name) {
		name = *args.Name
	}

	page, err := r.Client.SearchVehicles(ctx, name)
	if err != nil {
		return []*VehicleResolver{}, nil
	}

	return NewVehicles(ctx, NewVehiclesArgs{Page: page})
}
