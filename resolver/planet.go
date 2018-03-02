package resolver

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/loader"
	"github.com/tonyghita/graphql-go-example/swapi"

	graphql "github.com/graph-gophers/graphql-go"
)

// PlanetResolver resolves the Planet type.
type PlanetResolver struct {
	planet swapi.Planet
}

type NewPlanetArgs struct {
	Planet swapi.Planet
	URL    string
}

type NewPlanetsArgs struct {
	Page swapi.PlanetPage
	URLs []string
}

func NewPlanet(ctx context.Context, args NewPlanetArgs) (*PlanetResolver, error) {
	var planet swapi.Planet
	var err error

	switch {
	case args.Planet.URL != "":
		planet = args.Planet
	case args.URL != "":
		planet, err = loader.LoadPlanet(ctx, args.URL)
	default:
		err = errors.UnableToResolve
	}

	if err != nil {
		return nil, err
	}

	return &PlanetResolver{planet: planet}, nil
}

func NewPlanets(ctx context.Context, args NewPlanetsArgs) (*[]*PlanetResolver, error) {
	err := loader.PrimePlanets(ctx, args.Page)
	if err != nil {
		return nil, err
	}

	results, err := loader.LoadPlanets(ctx, append(args.URLs, args.Page.URLs()...))
	if err != nil {
		return nil, err
	}

	var planets = results.WithoutErrors()
	var resolvers = make([]*PlanetResolver, 0, len(planets))
	var errs errors.Errors

	for i, planet := range planets {
		r, err := NewPlanet(ctx, NewPlanetArgs{Planet: planet})
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}
		resolvers = append(resolvers, r)
	}

	return &resolvers, errs.Err()
}

// ID resolves ..
func (r *PlanetResolver) ID() graphql.ID {
	return extractID(r.planet.URL)
}

// Name resolves ...
func (r *PlanetResolver) Name() string {
	return r.planet.Name
}

// Diameter resolves ...
func (r *PlanetResolver) Diameter(args LengthUnitArgs) (float64, error) {
	return strconv.ParseFloat(r.planet.Diameter, 64)
}

// RotationPeriod resolves ...
func (r *PlanetResolver) RotationPeriod() (float64, error) {
	return strconv.ParseFloat(r.planet.RotationPeriod, 64)
}

// OrbitalPeriod resolves ...
func (r *PlanetResolver) OrbitalPeriod() (float64, error) {
	return strconv.ParseFloat(r.planet.OrbitalPeriod, 64)
}

// Gravity resolves ...
func (r *PlanetResolver) Gravity() (float64, error) {
	return strconv.ParseFloat(r.planet.Gravity, 64)
}

// Population resolves ...
func (r *PlanetResolver) Population() (int32, error) {
	p, err := strconv.ParseInt(r.planet.Population, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(p), nil
}

// Climates resolves ...
func (r *PlanetResolver) Climates() []string {
	return strings.Split(r.planet.Climate, ",")
}

// Terrain resolves ...
func (r *PlanetResolver) Terrains() []string {
	return strings.Split(r.planet.Terrain, ",")
}

// SurfaceWaterPercentage resolves ...
func (r *PlanetResolver) SurfaceWaterPercentage() (float64, error) {
	return strconv.ParseFloat(r.planet.SurfaceWater, 64)
}

// Residents resolves ...
func (r *PlanetResolver) Residents(ctx context.Context) (*[]*PersonResolver, error) {
	return NewPeople(ctx, NewPeopleArgs{URLs: r.planet.ResidentURLs})
}

// Films resolves ...
func (r *PlanetResolver) Films(ctx context.Context) (*[]*FilmResolver, error) {
	return NewFilms(ctx, NewFilmsArgs{URLs: r.planet.FilmURLs})
}

// CreatedAt resolves ...
func (r *PlanetResolver) CreatedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, r.planet.CreatedAt)
	if err != nil {
		return graphql.Time{}, err
	}

	return graphql.Time{Time: t}, nil
}

// EditedAt resolves ...
func (r *PlanetResolver) EditedAt() (*graphql.Time, error) {
	if r.planet.EditedAt == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, r.planet.EditedAt)
	if err != nil {
		return nil, err
	}

	return &graphql.Time{Time: t}, nil
}
