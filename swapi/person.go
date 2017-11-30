package swapi

import (
	"context"
	"net/url"
)

type Person struct {
	BirthYear    string   `json:"birth_year"`
	CreatedAt    string   `json:"created"`
	EditedAt     string   `json:"edited"`
	EyeColor     string   `json:"eye_color"`
	FilmURLs     []string `json:"films"`
	Gender       string   `json:"gender"`
	HairColor    string   `json:"hair_color"`
	Height       string   `json:"height"`
	HomeworldURL string   `json:"homeworld"`
	Mass         string   `json:"mass"`
	Name         string   `json:"name"`
	SkinColor    string   `json:"skin_color"`
	SpeciesURLs  []string `json:"species"`
	StarshipURLs []string `json:"starships"`
	URL          string   `json:"url"`
	VehicleURLs  []string `json:"vehicles"`
}

type PersonPage struct {
	Count  int64    `json:"count"`
	People []Person `json:"results"`
}

func (p PersonPage) URLs() []string {
	urls := make([]string, 0, len(p.People))
	for _, person := range p.People {
		urls = append(urls, person.URL)
	}
	return urls
}

func (c *Client) Person(ctx context.Context, url string) (Person, error) {
	r, err := c.NewRequest(ctx, url)
	if err != nil {
		return Person{}, err
	}

	var p Person
	if _, err := c.Do(r, &p); err != nil {
		return Person{}, err
	}

	return p, nil
}

func (c *Client) SearchPerson(ctx context.Context, name string) (PersonPage, error) {
	q := url.Values{"search": {name}}
	r, err := c.NewRequest(ctx, "/people?"+q.Encode())
	if err != nil {
		return PersonPage{}, nil
	}

	var pp PersonPage
	if _, err := c.Do(r, &pp); err != nil {
		return PersonPage{}, nil
	}

	return pp, nil
}
