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

// PlanetResolver resolves the Planet type.
type PlanetResolver struct {
	planet swapi.Planet
}

type NewPlanetArgs struct {
	URL string
}

type NewPlanetsArgs struct {
	URLs []string
}

func NewPlanet(ctx context.Context, args NewPlanetArgs) (*PlanetResolver, error) {
	planet, err := loader.LoadPlanet(ctx, args.URL)
	if err != nil {
		return nil, err
	}

	return &PlanetResolver{planet: planet}, nil
}

func NewPlanets(ctx context.Context, args NewPlanetsArgs) ([]*PlanetResolver, error) {
	// TODO: validate arguments
	// Request a batch to avoid sequentially loading when creating individual resolvers.
	loader.LoadPlanets(ctx, args.URLs)

	var resolvers = make([]*PlanetResolver, len(args.URLs))
	var errs errors.Errors

	for i, url := range args.URLs {
		r, err := NewPlanet(ctx, NewPlanetArgs{URL: url})
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}
		resolvers = append(resolvers, r)
	}

	return resolvers, errs.Err()
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
func (r *PlanetResolver) Residents(ctx context.Context) ([]*PersonResolver, error) {
	return NewPeople(ctx, NewPeopleArgs{URLs: r.planet.ResidentURLs})
}

// Films resolves ...
func (r *PlanetResolver) Films(ctx context.Context) ([]*FilmResolver, error) {
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
