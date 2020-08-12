package apartments

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

func (i *instrumentingService) GetApartments(ctx context.Context, city City, limit, offset int) ([]Apartment, error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "GetApartments").Add(1)
		i.requestLatency.With("method", "GetApartments").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.Service.GetApartments(ctx, city, limit, offset)
}
