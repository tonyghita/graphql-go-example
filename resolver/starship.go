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

// The StarshipResolver contains the data required to resolve the Starship type.
type StarshipResolver struct {
	ship swapi.Starship
}

type NewStarshipArgs struct {
	Ship swapi.Starship
	URL  string
}

type NewStarshipsArgs struct {
	Page swapi.StarshipPage
	URLs []string
}

func NewStarship(ctx context.Context, args NewStarshipArgs) (*StarshipResolver, error) {
	var ship swapi.Starship
	var err error

	switch {
	case args.Ship.URL != "":
		ship = args.Ship
	case args.URL != "":
		ship, err = loader.LoadStarship(ctx, args.URL)
	default:
		err = errors.UnableToResolve
	}

	if err != nil {
		return nil, err
	}

	return &StarshipResolver{ship: ship}, nil
}

func NewStarships(ctx context.Context, args NewStarshipsArgs) (*[]*StarshipResolver, error) {
	err := loader.PrimeStarships(ctx, args.Page)
	if err != nil {
		return nil, err
	}

	results, err := loader.LoadStarships(ctx, append(args.URLs, args.Page.URLs()...))
	if err != nil {
		return nil, err
	}

	var ships = results.WithoutErrors()
	var resolvers = make([]*StarshipResolver, 0, len(ships))
	var errs errors.Errors

	for i, ship := range ships {
		resolver, err := NewStarship(ctx, NewStarshipArgs{Ship: ship})
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		resolvers = append(resolvers, resolver)
	}

	return &resolvers, errs.Err()
}

// ID resolves ...
func (r *StarshipResolver) ID() graphql.ID {
	return extractID(r.ship.URL)
}

// Name resolves ...
func (r *StarshipResolver) Name() string {
	return r.ship.Name
}

// Model resolves ...
func (r *StarshipResolver) Model() string {
	return r.ship.Model
}

// Class resolves ...
func (r *StarshipResolver) Class() string {
	return r.ship.StarshipClass
}

// Manufacturers resolves ...
func (r *StarshipResolver) Manufacturers() []string {
	return strings.Split(r.ship.Manufacturer, ",")
}

// Cost resolves ...
func (r *StarshipResolver) Cost() (int32, error) {
	c, err := strconv.ParseInt(r.ship.CostInCredits, 10, 32)
	if err != nil {
		return 0, errors.UnableToResolve
	}

	return int32(c), nil
}

// Length resolves ...
func (r *StarshipResolver) Length(args LengthUnitArgs) (float64, error) {
	unit, err := ToLengthUnit(args.Unit)
	if err != nil {
		return 0.0, err
	}

	l, err := strconv.ParseFloat(r.ship.Length, 64)
	if err != nil {
		return 0.0, err
	}

	return ConvertLength(l, Meter, unit), nil
}

// CrewSize resolves ...
func (r *StarshipResolver) CrewSize() (int32, error) {
	s, err := strconv.ParseInt(r.ship.Crew, 10, 32)
	if err != nil {
		return 0, errors.UnableToResolve
	}

	return int32(s), nil
}

// PassengerCapacity resolves ...
func (r *StarshipResolver) PassengerCapacity() (int32, error) {
	c, err := strconv.ParseInt(r.ship.CargoCapacity, 10, 32)
	if err != nil {
		return 0, errors.UnableToResolve
	}

	return int32(c), nil
}

// MaxAtmosphericSpeed resolves ...
func (r *StarshipResolver) MaxAtmosphericSpeed() (*int32, error) {
	if r.ship.MaxAtmospheringSpeed == "" {
		return nil, nil
	}

	s, err := strconv.ParseInt(r.ship.MaxAtmospheringSpeed, 10, 32)
	if err != nil {
		return nil, errors.UnableToResolve
	}

	i := int32(s)
	return &i, nil
}

// HyperdriveRating resolves ...
func (r *StarshipResolver) HyperdriveRating() (*float64, error) {
	if r.ship.HyperdriveRating == "" {
		return nil, nil
	}

	f, err := strconv.ParseFloat(r.ship.HyperdriveRating, 64)
	if err != nil {
		return nil, errors.UnableToResolve
	}

	return &f, nil
}

// MaxMegalightsPerHour ...
func (r *StarshipResolver) MaxMegalightsPerHour() (int32, error) {
	i, err := strconv.ParseInt(r.ship.MGLT, 10, 32)
	if err != nil {
		return 0, errors.UnableToResolve
	}

	return int32(i), nil
}

// CargoCapacity resolves ...
func (r *StarshipResolver) CargoCapacity(args LengthUnitArgs) (float64, error) {
	f, err := strconv.ParseFloat(r.ship.CargoCapacity, 64)
	if err != nil {
		return 0.0, errors.UnableToResolve
	}

	return f, nil
}

// ConsumablesDuration resolves ...
func (r *StarshipResolver) ConsumablesDuration() string {
	return r.ship.Consumables
}

// Films resolves ...
func (r *StarshipResolver) Films(ctx context.Context) (*[]*FilmResolver, error) {
	return NewFilms(ctx, NewFilmsArgs{URLs: r.ship.FilmURLs})
}

// Pilots resolves ...
func (r *StarshipResolver) Pilots(ctx context.Context) (*[]*PersonResolver, error) {
	return NewPeople(ctx, NewPeopleArgs{URLs: r.ship.PilotURLs})
}

// CreatedAt resolves ...
func (r *StarshipResolver) CreatedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, r.ship.CreatedAt)
	if err != nil {
		return graphql.Time{}, errors.UnableToResolve
	}

	return graphql.Time{Time: t}, nil
}

// EditedAt resolves ...
func (r *StarshipResolver) EditedAt() (*graphql.Time, error) {
	if r.ship.EditedAt == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, r.ship.EditedAt)
	if err != nil {
		return nil, errors.UnableToResolve
	}

	return &graphql.Time{Time: t}, nil
}
