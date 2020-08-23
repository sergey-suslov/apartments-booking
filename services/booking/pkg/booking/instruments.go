package booking

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

func (i *InstrumentingService) GetReservations(ctx context.Context, apartmentID string, start, end time.Time) (out []Reservation, err error) { //nolint:lll
	defer func(begin time.Time) {
		i.requestCount.With("method", "GetReservations").Add(1)
		i.requestLatency.With("method", "GetReservations").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.Service.GetReservations(ctx, apartmentID, start, end)
}

func (i *InstrumentingService) BookApartment(ctx context.Context, userID, apartmentID string, start, end time.Time) (out *Reservation, err error) { //nolint:lll
	defer func(begin time.Time) {
		i.requestCount.With("method", "BookApartment").Add(1)
		i.requestLatency.With("method", "BookApartment").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.Service.BookApartment(ctx, userID, apartmentID, start, end)
}
