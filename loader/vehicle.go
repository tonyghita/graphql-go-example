package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"
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
