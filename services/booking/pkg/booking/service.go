package booking

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

const MaxReservationQueryingTimespan = 3600 * 24 * 40
const MaxReservationTimespan = 3600 * 24 * 30

var ErrRequestingDatabase = errors.New("error requesting data from db")
var ErrTooWideTimeSpan = errors.New("too wide time span")
var ErrReservationDurationLimitExceeded = errors.New("reservation duration limit exceeded")
var ErrCouldNotGetApartment = errors.New("error with requesting apartment")
var ErrNoApartmentWithGivenID = errors.New("no apartment with given id")

type City string

type Reservation struct {
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	ApartmentID string              `json:"apartmentId"`
	UserID      string              `json:"userId"`
	Start       primitive.Timestamp `json:"start"`
	End         primitive.Timestamp `json:"end"`
	Created     primitive.Timestamp `json:"created"`
}

func NewReservation(apartmentID, userID string, start, end time.Time) *Reservation {
	return &Reservation{
		ApartmentID: apartmentID,
		UserID:      userID,
		Start:       TimeToTimestamp(start),
		End:         TimeToTimestamp(end),
		Created:     TimeToTimestamp(time.Now()),
	}
}

type Apartment struct {
	ID      string `json:"_id"`
	Title   string `json:"title"`
	Address string `json:"address"`
	Owner   string `json:"owner"`
	City    string `json:"city"`
}

type Service interface {
	GetReservations(ctx context.Context, apartmentID string, start, end time.Time) (out []Reservation, err error)
	BookApartment(ctx context.Context, userID, apartmentID string, start, end time.Time) (out *Reservation, err error)
}

type Repository interface {
	GetReservationsBetween(ctx context.Context, apartmentID string, start, end time.Time) ([]Reservation, error)
	MakeReservation(ctx context.Context, reservation *Reservation) (*Reservation, error)
}

type ApartmentsRepository interface {
	GetApartmentByID(ctx context.Context, apartmentID string) (*Apartment, error)
}

type service struct {
	r      Repository
	ar     ApartmentsRepository
	logger *zap.Logger
}

func NewService(r Repository, ar ApartmentsRepository, logger *zap.Logger) Service {
	return &service{r: r, ar: ar, logger: logger}
}

func (s *service) GetReservations(ctx context.Context, apartmentID string, start, end time.Time) ([]Reservation, error) {
	if end.Sub(start).Seconds() > MaxReservationQueryingTimespan {
		return nil, ErrTooWideTimeSpan
	}
	return s.r.GetReservationsBetween(ctx, apartmentID, start, end)
}

func (s *service) BookApartment(ctx context.Context, userID, apartmentID string, start, end time.Time) (*Reservation, error) {
	if end.Sub(start).Seconds() > MaxReservationTimespan {
		return nil, ErrReservationDurationLimitExceeded
	}

	_, err := s.ar.GetApartmentByID(ctx, apartmentID)
	if err != nil {
		s.logger.Error("error getting apartment from apartments service", zap.Error(err))
		return nil, ErrCouldNotGetApartment
	}

	return s.r.MakeReservation(ctx, NewReservation(apartmentID, userID, start, end))
}

func TimeToTimestamp(t time.Time) primitive.Timestamp {
	return primitive.Timestamp{
		T: uint32(t.Unix()),
	}
}
