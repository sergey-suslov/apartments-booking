package apartments

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DatabaseError = errors.New("error requesting data from db")

type City string

type Apartment struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Title   string             `json:"title"`
	Address string             `json:"address"`
	Owner   string             `json:"owner"`
}

type Service interface {
	GetApartments(ctx context.Context, city City, limit, offset int) ([]Apartment, error)
}

type Repository interface {
	GetApartmentsByCity(ctx context.Context, city City, limit, offset int) ([]Apartment, error)
}

type service struct {
	ar Repository
}

func NewService(ar Repository) *service {
	return &service{ar: ar}
}

func (s *service) GetApartments(ctx context.Context, city City, limit, offset int) ([]Apartment, error) {
	panic("implement me")
}
