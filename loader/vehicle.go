package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"
	"github.com/tonyghita/graphql-go-example/errors"
	"github.com/tonyghita/graphql-go-example/swapi"
)

type vehicleGetter interface {
	Vehicle(ctx context.Context, url string) (swapi.Vehicle, error)
}

type VehicleLoader struct {
	get vehicleGetter
}

func newVehicleLoader(client vehicleGetter) dataloader.BatchFunc {
	return VehicleLoader{get: client}.loadBatch
}

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

func LoadVehicles(ctx context.Context, urls []string) ([]swapi.Vehicle, error) {
	ldr, err := extract(ctx, vehicleLoaderKey)
	if err != nil {
		return []swapi.Vehicle{}, err
	}

	data, loadErrs := ldr.LoadMany(ctx, urls)()

	var (
		vehicles = make([]swapi.Vehicle, 0, len(data))
		errs     = make(errors.Errors, 0, len(loadErrs))
	)

	for i := range urls {
		d, err := data[i], loadErrs[i]
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		vehicle, ok := d.(swapi.Vehicle)
		if !ok && err == nil {
			errs = append(errs, errors.WithIndex(errors.UnexpectedResponse, i))
		}

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, errs.Err()
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
