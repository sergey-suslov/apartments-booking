package apartments

import (
	"context"
	"go.uber.org/zap"
	"time"
)

type loggingService struct {
	logger zap.Logger
	Service
}

func NewLoggingService(logger zap.Logger, service Service) *loggingService {
	return &loggingService{logger: logger, Service: service}
}

func (s *loggingService) GetApartments(ctx context.Context, city City, limit, offset int) (a []Apartment, err error) {
	defer func(begin time.Time) {
		s.logger.With(
			zap.Field{Key: "method", String: "GetApartments"},
			zap.Field{Key: "took", Integer: time.Since(begin).Milliseconds()},
			zap.Field{Key: "error", String: err.Error()})
	}(time.Now())
	return s.Service.GetApartments(ctx, city, limit, offset)
}
