package apartments

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

type InstrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(requestCount metrics.Counter, requestLatency metrics.Histogram, service Service) *InstrumentingService {
	return &InstrumentingService{requestCount: requestCount, requestLatency: requestLatency, Service: service}
}

func (i *InstrumentingService) GetApartments(ctx context.Context, city City, limit, offset int) ([]Apartment, error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "GetApartments").Add(1)
		i.requestLatency.With("method", "GetApartments").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.Service.GetApartments(ctx, city, limit, offset)
}

func (i *InstrumentingService) GetApartmentByID(ctx context.Context, apartmentID string) (a *Apartment, err error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "GetApartmentByID").Add(1)
		i.requestLatency.With("method", "GetApartmentByID").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.Service.GetApartmentByID(ctx, apartmentID)
}
