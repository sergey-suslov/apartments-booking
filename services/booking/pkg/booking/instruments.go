package booking

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(requestCount metrics.Counter, requestLatency metrics.Histogram, service Service) *instrumentingService {
	return &instrumentingService{requestCount: requestCount, requestLatency: requestLatency, Service: service}
}

func (i *instrumentingService) GetReservations(ctx context.Context, apartmentId string, start, end time.Time) (out []Reservation, err error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "GetReservations").Add(1)
		i.requestLatency.With("method", "GetReservations").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.Service.GetReservations(ctx, apartmentId, start, end)
}

func (i *instrumentingService) BookApartment(ctx context.Context, userId, apartmentId string, start, end time.Time) (out *Reservation, err error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "BookApartment").Add(1)
		i.requestLatency.With("method", "BookApartment").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.Service.BookApartment(ctx, userId, apartmentId, start, end)
}
