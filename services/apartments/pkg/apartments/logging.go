package apartments

import (
	"context"
	"go.uber.org/zap"
	"time"
)

type loggingService struct {
	logger *zap.Logger
	Service
}

func NewLoggingService(logger *zap.Logger, service Service) Service {
	return &loggingService{logger: logger, Service: service}
}

func (s *loggingService) GetApartments(ctx context.Context, city City, limit, offset int) (a []Apartment, err error) {
	defer func(begin time.Time) {
		s.logger.Debug("calling GetApartments",
			zap.Duration("took", time.Since(begin)),
			zap.Error(err),
		)
	}(time.Now())
	return s.Service.GetApartments(ctx, city, limit, offset)
}
