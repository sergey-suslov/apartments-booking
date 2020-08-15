package booking

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

func (s *loggingService) GetReservations(ctx context.Context, apartmentId string, start, end time.Time) (out []Reservation, err error) {
	defer func(begin time.Time) {
		s.logger.Debug("calling GetReservations",
			zap.Duration("took", time.Since(begin)),
			zap.Int("returned reservations", len(out)),
			zap.Error(err),
		)
	}(time.Now())
	return s.Service.GetReservations(ctx, apartmentId, start, end)
}

func (s *loggingService) BookApartment(ctx context.Context, userId, apartmentId string, start, end time.Time) (out *Reservation, err error) {
	defer func(begin time.Time) {
		s.logger.Debug("calling BookApartment",
			zap.Duration("took", time.Since(begin)),
			zap.Any("returned reservation", out),
			zap.Error(err),
		)
	}(time.Now())
	return s.Service.BookApartment(ctx, userId, apartmentId, start, end)
}
