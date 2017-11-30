package resolver

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/loader"
	"github.com/tonyghita/graphql-go-example/swapi"

	graphql "github.com/neelance/graphql-go"
)

// The SpeciesResolver resolves the Species type.
type SpeciesResolver struct {
	species swapi.Species
}

// NewSpeciesArgs ...
type NewSpeciesArgs struct {
	Species swapi.Species
	URL     string
}

// NewSpeciesListArgs ...
type NewSpeciesListArgs struct {
	Page swapi.SpeciesPage
	URLs []string
}

// NewSpecies ...
func NewSpecies(ctx context.Context, args NewSpeciesArgs) (*SpeciesResolver, error) {
	var species swapi.Species
	var err error

	switch {
	case args.Species.URL != "":
		species = args.Species
	case args.URL != "":
		species, err = loader.LoadSpecies(ctx, args.URL)
	default:
		err = errors.UnableToResolve
	}

	if err != nil {
		return nil, err
	}

	return &SpeciesResolver{species: species}, nil
}

// NewSpeciesList ...
func NewSpeciesList(ctx context.Context, args NewSpeciesListArgs) ([]*SpeciesResolver, error) {
	loader.PrimeSpecies(ctx, args.Page)

	species, err := loader.LoadManySpecies(ctx, append(args.URLs, args.Page.URLs()...)...)
	if err != nil {
		return []*SpeciesResolver{}, err
	}

	var resolvers = make([]*SpeciesResolver, len(species))
	var errs errors.Errors

	for i, sp := range species {
		resolver, err := NewSpecies(ctx, NewSpeciesArgs{Species: sp})
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		resolvers[i] = resolver
	}

	return resolvers, errs.Err()
}

// ID resolves this species unique identifier.
func (r *SpeciesResolver) ID() graphql.ID {
	return extractID(r.species.URL)
}

// Name resolves the name of the species.
func (r *SpeciesResolver) Name() string {
	return r.species.Name
}

// Classification resolves the classification of this species, such as "mammal" or "reptile".
func (r *SpeciesResolver) Classification() string {
	return r.species.Classification
}

// Designation ...
func (r *SpeciesResolver) Designation() string {
	return r.species.Designation
}

// AverageHeight ...
func (r *SpeciesResolver) AverageHeight(args LengthUnitArgs) (float64, error) {
	unit, err := ToLengthUnit(args.Unit)
	if err != nil {
		return 0.0, err
	}

	h, err := strconv.ParseFloat(r.species.AverageHeight, 64)
	if err != nil {
		return 0.0, err
	}

	return ConvertLength(h, Centimeter, unit), nil
}

// AverageLifespan ...
func (r *SpeciesResolver) AverageLifespan() (float64, error) {
	years, err := strconv.ParseFloat(r.species.AverageLifespan, 64)
	if err != nil {
		return 0.0, err
	}

	return years, nil
}

// EyeColors ...
func (r *SpeciesResolver) EyeColors() []string {
	return strings.Split(r.species.EyeColors, ", ")
}

// HairColors ...
func (r *SpeciesResolver) HairColors() []string {
	return strings.Split(r.species.HairColors, ", ")
}

// SkinColors ...
func (r *SpeciesResolver) SkinColors() []string {
	return strings.Split(r.species.SkinColors, ", ")
}

// Language ...
func (r *SpeciesResolver) Language() string {
	return r.species.Language
}

// Homeworld ...
func (r *SpeciesResolver) Homeworld(ctx context.Context) (*PlanetResolver, error) {
	return nil, nil
}

// Characters ...
func (r *SpeciesResolver) Characters(ctx context.Context) ([]*PersonResolver, error) {
	return nil, nil
}

// Films ...
func (r *SpeciesResolver) Films(ctx context.Context) ([]*FilmResolver, error) {
	return NewFilms(ctx, NewFilmsArgs{URLs: r.species.FilmURLs})
}

// CreatedAt ...
func (r *SpeciesResolver) CreatedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, r.species.CreatedAt)
	return graphql.Time{Time: t}, err
}

// EditedAt ...
func (r *SpeciesResolver) EditedAt() (*graphql.Time, error) {
	if r.species.EditedAt == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, r.species.EditedAt)
	if err != nil {
		return nil, err
	}

	return &graphql.Time{Time: t}, nil
}
