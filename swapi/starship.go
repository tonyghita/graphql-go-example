package swapi

import (
	"context"
	"net/url"
)

type Starship struct {
	MGLT                 string   `json:"MGLT"`
	CargoCapacity        string   `json:"cargo_capacity"`
	Consumables          string   `json:"consumables"`
	CostInCredits        string   `json:"cost_in_credits"`
	CreatedAt            string   `json:"created"`
	Crew                 string   `json:"crew"`
	EditedAt             string   `json:"edited"`
	FilmURLs             []string `json:"films"`
	HyperdriveRating     string   `json:"hyperdrive_rating"`
	Length               string   `json:"length"`
	Manufacturer         string   `json:"manufacturer"`
	MaxAtmospheringSpeed string   `json:"max_atmosphering_speed"`
	Model                string   `json:"model"`
	Name                 string   `json:"name"`
	Passengers           string   `json:"passengers"`
	PilotURLs            []string `json:"pilots"`
	StarshipClass        string   `json:"starship_class"`
	URL                  string   `json:"url"`
}

type StarshipPage struct {
	Count     int64      `json:"count"`
	Starships []Starship `json:"results"`
}

func (c *Client) Starship(ctx context.Context, url string) (Starship, error) {
	r, err := c.NewRequest(ctx, url)
	if err != nil {
		return Starship{}, err
	}

	var s Starship
	if _, err := c.Do(r, &s); err != nil {
		return Starship{}, err
	}

	return s, nil
}

func (c *Client) SearchStarships(ctx context.Context, name string) (StarshipPage, error) {
	q := url.Values{"search": {name}}
	r, err := c.NewRequest(ctx, "/starships?"+q.Encode())
	if err != nil {
		return StarshipPage{}, err
	}

	var sp StarshipPage
	if _, err := c.Do(r, &sp); err != nil {
		return StarshipPage{}, err
	}
	return sp, nil
}
