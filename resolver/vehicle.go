package resolver

import (
	"context"
	"strconv"
	"strings"
	"time"

	graphql "github.com/neelance/graphql-go"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/loader"
	"github.com/tonyghita/graphql-go-example/swapi"
)

// The VehicleResolver resolves the vehicle type.
type VehicleResolver struct {
	vehicle swapi.Vehicle
}

type NewVehicleArgs struct {
	Vehicle swapi.Vehicle
	URL     string
}

type NewVehiclesArgs struct {
	Page swapi.VehiclePage
	URLs []string
}

func NewVehicle(ctx context.Context, args NewVehicleArgs) (*VehicleResolver, error) {
	var vehicle swapi.Vehicle
	var err error

	switch {
	case args.Vehicle.URL != "":
		vehicle = args.Vehicle
	case args.URL != "":
		vehicle, err = loader.LoadVehicle(ctx, args.URL)
	default:
		err = errors.UnableToResolve
	}
	if err != nil {
		return nil, err
	}

	return &VehicleResolver{vehicle: vehicle}, nil
}

func NewVehicles(ctx context.Context, args NewVehiclesArgs) (*[]*VehicleResolver, error) {
	err := loader.PrimeVehicles(ctx, args.Page)
	if err != nil {
		return nil, err
	}

	results, err := loader.LoadVehicles(ctx, append(args.URLs, args.Page.URLs()...))
	if err != nil {
		return nil, err
	}

	var vehicles = results.WithoutErrors()
	var resolvers = make([]*VehicleResolver, 0, len(vehicles))
	var errs errors.Errors

	for i, v := range vehicles {
		resolver, err := NewVehicle(ctx, NewVehicleArgs{Vehicle: v})
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		resolvers = append(resolvers, resolver)
	}

	return &resolvers, errs.Err()
}

// ID resolves ...
func (r *VehicleResolver) ID() graphql.ID {
	return extractID(r.vehicle.URL)
}

// Name resolves ...
func (r *VehicleResolver) Name() string {
	return r.vehicle.Name
}

// Model resolves ...
func (r *VehicleResolver) Model() string {
	return r.vehicle.Model
}

// Class resolves ...
func (r *VehicleResolver) Class() string {
	return r.vehicle.VehicleClass
}

// Manufacturers resolves ...
func (r *VehicleResolver) Manufacturers() []string {
	return strings.Split(r.vehicle.Manufacturer, ",")
}

// Length resolves ...
func (r *VehicleResolver) Length(args LengthUnitArgs) (float64, error) {
	unit, err := ToLengthUnit(args.Unit)
	if err != nil {
		return 0.0, err
	}

	l, err := strconv.ParseFloat(r.vehicle.Length, 64)
	if err != nil {
		return 0.0, err
	}

	return ConvertLength(l, Meter, unit), nil
}

// Cost resolves ...
func (r *VehicleResolver) Cost() (int32, error) {
	c, err := strconv.ParseInt(r.vehicle.CostInCredits, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(c), nil
}

// CrewSize resolves ...
func (r *VehicleResolver) CrewSize() (int32, error) {
	s, err := strconv.ParseInt(r.vehicle.Crew, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(s), nil
}

// PassengerCapacity resolves ...
func (r *VehicleResolver) PassengerCapacity() (int32, error) {
	c, err := strconv.ParseInt(r.vehicle.Passengers, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(c), nil
}

// MaxAtmosphericSpeed resolves ...
func (r *VehicleResolver) MaxAtmosphericSpeed() (float64, error) {
	return strconv.ParseFloat(r.vehicle.MaxAtmospheringSpeed, 64)
}

// CargoCapacity resolves ...
func (r *VehicleResolver) CargoCapacity(args MassUnitArgs) (float64, error) {
	c, err := strconv.ParseFloat(r.vehicle.CargoCapacity, 64)
	if err != nil {
		return 0.0, err
	}

	u, err := ToMassUnit(args.Unit)
	if err != nil {
		return 0.0, err
	}

	return ConvertMass(c, Kilogram, u), nil
}

// ConsumablesDuration resolves ...
func (r *VehicleResolver) ConsumablesDuration() string {
	return r.vehicle.Consumables
}

// Films resolves ...
func (r *VehicleResolver) Films(ctx context.Context) (*[]*FilmResolver, error) {
	return NewFilms(ctx, NewFilmsArgs{URLs: r.vehicle.FilmURLs})
}

// Pilots resolves ...
func (r *VehicleResolver) Pilots(ctx context.Context) (*[]*PersonResolver, error) {
	return NewPeople(ctx, NewPeopleArgs{URLs: r.vehicle.PilotURLs})
}

// CreatedAt resolves ...
func (r *VehicleResolver) CreatedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, r.vehicle.CreatedAt)
	if err != nil {
		return graphql.Time{}, errors.UnableToResolve
	}

	return graphql.Time{Time: t}, nil
}

// EditedAt resolves ...
func (r *VehicleResolver) EditedAt() (*graphql.Time, error) {
	if r.vehicle.EditedAt == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, r.vehicle.EditedAt)
	if err != nil {
		return nil, errors.UnableToResolve
	}

	return &graphql.Time{Time: t}, nil
}
