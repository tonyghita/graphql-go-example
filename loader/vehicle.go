package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"
	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"
)

func LoadVehicle(ctx context.Context, url string) (swapi.Vehicle, error) {
	ldr, err := extract(ctx, vehicleLoaderKey)
	if err != nil {
		return swapi.Vehicle{}, err
	}

	data, err := ldr.Load(ctx, url)()
	if err != nil {
		return swapi.Vehicle{}, err
	}

	vehicle, ok := data.(swapi.Vehicle)
	if !ok {
		return swapi.Vehicle{}, errors.UnexpectedResponse
	}

	return vehicle, nil
}

func LoadVehicles(ctx context.Context, urls []string) (VehicleResults, error) {
	var results []VehicleResult

	ldr, err := extract(ctx, vehicleLoaderKey)
	if err != nil {
		return results, err
	}

	data, errs := ldr.LoadMany(ctx, urls)()

	for i, d := range data {
		var e error
		if errs != nil {
			e = errs[i]
		}

		vehicle, ok := d.(swapi.Vehicle)
		if !ok && e == nil {
			e = errors.UnexpectedResponse
		}

		results = append(results, VehicleResult{Vehicle: vehicle, Error: e})
	}

	return results, nil
}

type VehicleResult struct {
	Vehicle swapi.Vehicle
	Error   error
}

type VehicleResults []VehicleResult

func (results VehicleResults) WithoutErrors() []swapi.Vehicle {
	vehicles := make([]swapi.Vehicle, 0, len(results))
	for _, r := range results {
		if r.Error != nil {
			continue
		}
		vehicles = append(vehicles, r.Vehicle)
	}
	return vehicles
}

func PrimeVehicles(ctx context.Context, page swapi.VehiclePage) error {
	ldr, err := extract(ctx, vehicleLoaderKey)
	if err != nil {
		return err
	}

	for _, v := range page.Vehicles {
		ldr.Prime(v.URL, v)
	}
	return nil
}

type vehicleGetter interface {
	Vehicle(ctx context.Context, url string) (swapi.Vehicle, error)
}

type VehicleLoader struct {
	get vehicleGetter
}

func newVehicleLoader(client vehicleGetter) dataloader.BatchFunc {
	return VehicleLoader{get: client}.loadBatch
}

func (ldr VehicleLoader) loadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(ctx context.Context, i int, url string) {
			data, err := ldr.get.Vehicle(ctx, url)
			results[i] = &dataloader.Result{Data: data, Error: err}
			wg.Done()
		}(ctx, i, url)
	}

	wg.Wait()

	return results
}
