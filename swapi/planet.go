package swapi

import (
	"context"
)

type Planet struct {
	Climate        string   `json:"climate"`         // Comma-separated if more than one.
	CreatedAt      string   `json:"created"`         // ...
	Diameter       string   `json:"diameter"`        // String float value in kilometers.
	EditedAt       string   `json:"edited"`          // ...
	FilmURLs       []string `json:"films"`           // ...
	Gravity        string   `json:"gravity"`         // String float value.
	Name           string   `json:"name"`            // ...
	OrbitalPeriod  string   `json:"orbital_period"`  // String integer value.
	Population     string   `json:"population"`      // String integer value.
	ResidentURLs   []string `json:"residents"`       // ...
	RotationPeriod string   `json:"rotation_period"` // String integer value.
	SurfaceWater   string   `json:"surface_water"`   // String float value.
	Terrain        string   `json:"terrain"`         // Comma-separated if more than one.
	URL            string   `json:"url"`             // ...
}

type PlanetPage struct {
	Count   int64    `json:"count"`
	Planets []Planet `json:"results"`
}

func (p PlanetPage) URLs() []string {
	urls := make([]string, 0, len(p.Planets))
	for _, planet := range p.Planets {
		urls = append(urls, planet.URL)
	}
	return urls
}

func (c *Client) Planet(ctx context.Context, url string) (Planet, error) {
	// TODO: implement
	return Planet{}, nil
}

func (c *Client) SearchPlanets(ctx context.Context, name string) (PlanetPage, error) {
	// TODO: implement
	return PlanetPage{}, nil
}
