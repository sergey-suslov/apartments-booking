package booking

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var DatabaseError = errors.New("error requesting data from db")
var TooWideTimeSpan = errors.New("too wide time span")

type City string

type Reservation struct {
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	ApartmentId string              `json:"apartmentId"`
	UserId      string              `json:"userId"`
	Start       primitive.Timestamp `json:"start"`
	End         primitive.Timestamp `json:"end"`
	Created     primitive.Timestamp `json:"created"`
}

type Service interface {
	GetReservations(ctx context.Context, apartmentId string, start, end time.Time) (out []Reservation, err error)
}

type Repository interface {
	GetReservationsBetween(ctx context.Context, apartmentId string, start, end time.Time) ([]Reservation, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

func (s *service) GetReservations(ctx context.Context, apartmentId string, start, end time.Time) ([]Reservation, error) {
	if end.Sub(start).Seconds() > 3600*24*40 {
		return nil, TooWideTimeSpan
	}
	return s.r.GetReservationsBetween(ctx, apartmentId, start, end)
}
