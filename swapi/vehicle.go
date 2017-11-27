package swapi

import (
	"context"
	"net/url"
)

type Vehicle struct {
	CargoCapacity        string   `json:"cargo_capacity"`         // Maximum integer kilogram capacity.
	Consumables          string   `json:"consumables"`            // Maximum length of time this vehicle can provide consumables for its crew without resupply.
	CostInCredits        string   `json:"cost_in_credits"`        // Integer galactic credits.
	CreatedAt            string   `json:"created"`                // ISO8601 date
	Crew                 string   `json:"crew"`                   // Integer number of crew required.
	EditedAt             string   `json:"edited"`                 // ISO8601 date
	FilmURLs             []string `json:"films"`                  // URLs for films this vehicle has appeared in.
	Length               string   `json:"length"`                 // Float meters.
	Manufacturer         string   `json:"manufacturer"`           // Comma-separated if more than one.
	MaxAtmospheringSpeed string   `json:"max_atmosphering_speed"` // An integer.
	Model                string   `json:"model"`                  // The official name, such as "All-Terrain Attack Transport".
	Name                 string   `json:"name"`                   // The common name, such as "Sand Crawler" or "Speeder bike".
	Passengers           string   `json:"passengers"`             // Integer number of non-essential people this vehicle can transport.
	PilotURLs            []string `json:"pilots"`                 // URLs for people who piloted this vehicle.
	URL                  string   `json:"url"`                    // URL for this resource.
	VehicleClass         string   `json:"vehicle_class"`          // Class of vehicle, such as "Wheeled" or "Repulsorcraft".
}

type VehiclePage struct {
	Count    int64     `json:"count"`
	Vehicles []Vehicle `json:"results"`
}

func (c *Client) Vehicle(ctx context.Context, url string) (Vehicle, error) {
	r, err := c.NewRequest(ctx, url)
	if err != nil {
		return Vehicle{}, err
	}

	var v Vehicle
	if _, err := c.Do(r, &v); err != nil {
		return Vehicle{}, err
	}

	return v, nil
}

func (c *Client) SearchVehicles(ctx context.Context, name string) (VehiclePage, error) {
	q := url.Values{"search": {name}}
	r, err := c.NewRequest(ctx, "/vehicles?"+q.Encode())
	if err != nil {
		return VehiclePage{}, err
	}

	var vp VehiclePage
	if _, err := c.Do(r, &vp); err != nil {
		return VehiclePage{}, err
	}

	return vp, nil
}
