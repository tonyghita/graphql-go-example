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

type NewPersonArgs struct {
	Person swapi.Person
	URL    string
}

type NewPeopleArgs struct {
	Page swapi.PersonPage
	URLs []string
}

func NewPerson(ctx context.Context, args NewPersonArgs) (*PersonResolver, error) {
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

	if err != nil {
		return nil, err
	}

	return &PersonResolver{person: person}, nil
}

func NewPeople(ctx context.Context, args NewPeopleArgs) (*[]*PersonResolver, error) {
	err := loader.PrimePeople(ctx, args.Page)
	if err != nil {
		return nil, err
	}

	results, err := loader.LoadPeople(ctx, append(args.URLs, args.Page.URLs()...))
	if err != nil {
		return nil, err
	}

	var people = results.WithoutErrors()
	var resolvers = make([]*PersonResolver, 0, len(people))
	var errs errors.Errors

	for i, person := range people {
		r, err := NewPerson(ctx, NewPersonArgs{Person: person})
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}
		resolvers = append(resolvers, r)
	}

	return &resolvers, errs.Err()
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
func (r *PersonResolver) Films(ctx context.Context) (*[]*FilmResolver, error) {
	return NewFilms(ctx, NewFilmsArgs{URLs: r.person.FilmURLs})
}

// Species resolves ...
func (r *PersonResolver) Species(ctx context.Context) (*[]*SpeciesResolver, error) {
	return nil, nil
}

// Vehicles resolves ...
func (r *PersonResolver) Vehicles(ctx context.Context) (*[]*VehicleResolver, error) {
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
