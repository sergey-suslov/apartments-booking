package booking

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type Errorer interface {
	Error() error
}

type getReservationsRequest struct {
	UserClaim
	ApartmentID string    `json:"apartmentId"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}

func (c *getReservationsRequest) SetUserClaim(claim *UserClaim) {
	c.UserClaim = *claim
}

type getReservationsResponse struct {
	Apartments []Reservation `json:"reservations"`
	Err        error         `json:"error,omitempty"`
}

func (g getReservationsResponse) Error() error {
	return g.Err
}

func makeGetApartmentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getReservationsRequest)
		apartments, err := s.GetReservations(ctx, req.ApartmentID, req.Start, req.End)
		return getReservationsResponse{
			Apartments: apartments,
			Err:        err,
		}, nil
	}
}

type bookRequest struct {
	UserClaim
	ApartmentID string    `json:"apartmentId"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}

func (c *bookRequest) SetUserClaim(claim *UserClaim) {
	c.UserClaim = *claim
}

type booksResponse struct {
	Reservation *Reservation `json:"reservation"`
	Err         error        `json:"error,omitempty"`
}

func (g booksResponse) Error() error {
	return g.Err
}

func makeBookApartmentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*bookRequest)
		reservation, err := s.BookApartment(ctx, req.ID, req.ApartmentID, req.Start, req.End)
		return booksResponse{Reservation: reservation, Err: err}, nil
	}
}
