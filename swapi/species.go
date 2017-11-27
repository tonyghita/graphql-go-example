package swapi

import (
	"context"
	"net/url"
)

// Species ...
type Species struct {
	AverageHeight   string   `json:"average_height"`   // String float value, in centimeters.
	AverageLifespan string   `json:"average_lifespan"` // String float value.
	Classification  string   `json:"classification"`   //
	CreatedAt       string   `json:"created"`          // ISO8601 format.
	Designation     string   `json:"designation"`      //
	EditedAt        string   `json:"edited"`           // ISO8601 format.
	EyeColors       string   `json:"eye_colors"`       // Comma-separated values, empty when no eyes.
	FilmURLs        []string `json:"films"`            //
	HairColors      string   `json:"hair_colors"`      // Comma-separated values, empty when no eyes.
	HomeworldURL    string   `json:"homeworld"`        //
	Language        string   `json:"language"`         //
	Name            string   `json:"name"`             //
	PeopleURLs      []string `json:"people"`           //
	SkinColors      string   `json:"skin_colors"`      //
	URL             string   `json:"url"`              //
}

// SpeciesPage ...
type SpeciesPage struct {
	Count   int64     `json:"count"`
	Species []Species `json:"results"`
}

func (p SpeciesPage) URLs() []string {
	urls := make([]string, len(p.Species))
	for i, f := range p.Species {
		urls[i] = f.URL
	}
	return urls
}

// Species ...
func (c *Client) Species(ctx context.Context, url string) (Species, error) {
	r, err := c.NewRequest(ctx, url)
	if err != nil {
		return Species{}, err
	}

	var s Species
	if _, err := c.Do(r, &s); err != nil {
		return Species{}, err
	}

	return s, nil
}

// SearchSpecies ...
func (c *Client) SearchSpecies(ctx context.Context, name string) (SpeciesPage, error) {
	q := url.Values{"search": {name}}
	r, err := c.NewRequest(ctx, "/species?"+q.Encode())
	if err != nil {
		return SpeciesPage{}, err
	}

	var sp SpeciesPage
	if _, err := c.Do(r, &sp); err != nil {
		return SpeciesPage{}, err
	}

	return sp, nil
}
