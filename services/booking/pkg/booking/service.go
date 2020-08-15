package booking

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"time"
)

const MaxReservationQueryingTimespan = 3600 * 24 * 40
const MaxReservationTimespan = 3600 * 24 * 30

var DatabaseError = errors.New("error requesting data from db")
var TooWideTimeSpan = errors.New("too wide time span")
var ReservationDurationLimitExceeded = errors.New("reservation duration limit exceeded")
var CouldNotGetApartment = errors.New("error with requesting apartment")
var NoApartmentWithGivenId = errors.New("no apartment with given id")

type City string

type Reservation struct {
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	ApartmentId string              `json:"apartmentId"`
	UserId      string              `json:"userId"`
	Start       primitive.Timestamp `json:"start"`
	End         primitive.Timestamp `json:"end"`
	Created     primitive.Timestamp `json:"created"`
}

func NewReservation(apartmentId string, userId string, start time.Time, end time.Time) *Reservation {
	return &Reservation{ApartmentId: apartmentId, UserId: userId, Start: TimeToTimestamp(start), End: TimeToTimestamp(end), Created: TimeToTimestamp(time.Now())}
}

type Apartment struct {
	ID      string `json:"_id"`
	Title   string `json:"title"`
	Address string `json:"address"`
	Owner   string `json:"owner"`
	City    string `json:"city"`
}

type Service interface {
	GetReservations(ctx context.Context, apartmentId string, start, end time.Time) (out []Reservation, err error)
	BookApartment(ctx context.Context, userId, apartmentId string, start, end time.Time) (out *Reservation, err error)
}

type Repository interface {
	GetReservationsBetween(ctx context.Context, apartmentId string, start, end time.Time) ([]Reservation, error)
	MakeReservation(ctx context.Context, reservation *Reservation) (*Reservation, error)
}

type ApartmentsRepository interface {
	GetApartmentById(ctx context.Context, apartmentId string) (*Apartment, error)
}

type service struct {
	r      Repository
	ar     ApartmentsRepository
	logger *zap.Logger
}

func NewService(r Repository, ar ApartmentsRepository, logger *zap.Logger) Service {
	return &service{r: r, ar: ar, logger: logger}
}

func (s *service) GetReservations(ctx context.Context, apartmentId string, start, end time.Time) ([]Reservation, error) {
	if end.Sub(start).Seconds() > MaxReservationQueryingTimespan {
		return nil, TooWideTimeSpan
	}
	return s.r.GetReservationsBetween(ctx, apartmentId, start, end)
}

func (s *service) BookApartment(ctx context.Context, userId, apartmentId string, start, end time.Time) (*Reservation, error) {
	if end.Sub(start).Seconds() > MaxReservationTimespan {
		return nil, ReservationDurationLimitExceeded
	}

	_, err := s.ar.GetApartmentById(ctx, apartmentId)
	if err != nil {
		s.logger.Error("error getting apartment from apartments service", zap.Error(err))
		return nil, CouldNotGetApartment
	}

	return s.r.MakeReservation(ctx, NewReservation(apartmentId, userId, start, end))
}

func TimeToTimestamp(t time.Time) primitive.Timestamp {
	return primitive.Timestamp{
		T: uint32(t.Unix()),
	}
}
