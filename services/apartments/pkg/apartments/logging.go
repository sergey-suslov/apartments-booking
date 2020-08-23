package apartments

import (
	"context"
	"time"

	"go.uber.org/zap"
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

func (s *loggingService) GetApartmentByID(ctx context.Context, apartmentID string) (a *Apartment, err error) {
	defer func(begin time.Time) {
		s.logger.Debug("calling GetApartmentByID",
			zap.Duration("took", time.Since(begin)),
			zap.String("requested apartmentID", apartmentID),
			zap.Bool("is apartment found", a != nil),
			zap.Error(err),
		)
	}(time.Now())
	return s.Service.GetApartmentByID(ctx, apartmentID)
}
