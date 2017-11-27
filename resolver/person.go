package resolver

import (
	"context"
	"strconv"
	"time"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/loader"
	"github.com/tonyghita/graphql-go-example/swapi"

	graphql "github.com/neelance/graphql-go"
)

// PersonResolver contains the person data to resolve.
type PersonResolver struct {
	person swapi.Person
}

// NewPersonArgs ...
type NewPersonArgs struct {
	Person swapi.Person
	URL    string
}

type NewPeopleArgs struct {
	URLs []string
}

// Load acts on the contained arguments (args) to load the specified person.
// TODO: not sure if this extra abstraction helps or hurts... the constructor code is nicer.
func (args NewPersonArgs) Load(ctx context.Context) (swapi.Person, error) {
	var person swapi.Person
	var err error

	switch {
	case args.Person.URL != "":
		person = args.Person
	case args.URL != "":
		person, err = loader.LoadPerson(ctx, args.URL)
	default:
		err = errors.UnableToResolve
	}

	return person, err
}

// NewPerson ...
func NewPerson(ctx context.Context, args NewPersonArgs) (*PersonResolver, error) {
	person, err := args.Load(ctx)
	if err != nil {
		return nil, err
	}

	return &PersonResolver{person: person}, nil
}

func NewPeople(ctx context.Context, args NewPeopleArgs) ([]*PersonResolver, error) {
	l, err := loader.Extract(ctx, loader.PeopleByURLs)
	if err != nil {
		return []*PersonResolver{}, nil
	}

	// Batch requests for people to avoid loading data in serial.
	l.LoadMany(ctx, args.URLs)

	var resolvers = make([]*PersonResolver, len(args.URLs))
	var errs errors.Errors

	for i, url := range args.URLs {
		r, err := NewPerson(ctx, NewPersonArgs{URL: url})
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}
		resolvers = append(resolvers, r)
	}

	return resolvers, errs.Err()
}

// ID resolves ...
func (r *PersonResolver) ID() graphql.ID {
	return extractID(r.person.URL)
}

// Name resolves ...
func (r *PersonResolver) Name() string {
	return r.person.Name
}

// BirthYear resolves ...
func (r *PersonResolver) BirthYear() string {
	return r.person.BirthYear
}

// EyeColor resolves ...
func (r *PersonResolver) EyeColor() *string {
	return nullableStr(r.person.EyeColor)
}

// Gender resolves ...
func (r *PersonResolver) Gender() *string {
	return nullableStr(r.person.Gender)
}

// HairColor resolves ...
func (r *PersonResolver) HairColor() *string {
	return nullableStr(r.person.HairColor)
}

// Height resolves ...
func (r *PersonResolver) Height(args LengthUnitArgs) (float64, error) {
	unit, err := ToLengthUnit(args.Unit)
	if err != nil {
		return 0.0, err
	}

	h, err := strconv.ParseFloat(r.person.Height, 64)
	if err != nil {
		return 0.0, err
	}

	return ConvertLength(h, Meter, unit), nil
}

// Mass resolves ...
func (r *PersonResolver) Mass(args MassUnitArgs) (float64, error) {
	unit, err := ToMassUnit(args.Unit)
	if err != nil {
		return 0.0, err
	}

	m, err := strconv.ParseFloat(r.person.Mass, 64)
	if err != nil {
		return 0.0, err
	}

	return ConvertMass(m, Kilogram, unit), nil
}

// SkinColor resolves ...
func (r *PersonResolver) SkinColor() *string {
	return nullableStr(r.person.SkinColor)
}

// Homeworld resolves ...
func (r *PersonResolver) Homeworld(ctx context.Context) (*PlanetResolver, error) {
	return nil, nil
}

// Films resolves ...
func (r *PersonResolver) Films(ctx context.Context) ([]*FilmResolver, error) {
	return NewFilms(ctx, NewFilmsArgs{URLs: r.person.FilmURLs})
}

// Species resolves ...
func (r *PersonResolver) Species(ctx context.Context) ([]*SpeciesResolver, error) {
	return nil, nil
}

// Vehicles resolves ...
func (r *PersonResolver) Vehicles(ctx context.Context) ([]*VehicleResolver, error) {
	return NewVehicles(ctx, NewVehiclesArgs{URLs: r.person.VehicleURLs})
}

// CreatedAt resolves ...
func (r *PersonResolver) CreatedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, r.person.CreatedAt)
	if err != nil {
		return graphql.Time{}, errors.UnableToResolve
	}

	return graphql.Time{Time: t}, nil
}

// EditedAt resolves ...
func (r *PersonResolver) EditedAt() (*graphql.Time, error) {
	if r.person.EditedAt == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, r.person.EditedAt)
	if err != nil {
		return nil, errors.UnableToResolve
	}

	return &graphql.Time{Time: t}, nil
}
