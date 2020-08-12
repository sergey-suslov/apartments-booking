package booking

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type Errorer interface {
	Error() error
}

type getReservationsRequest struct {
	ApartmentId string    `json:"apartmentId"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
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
		apartments, err := s.GetReservations(ctx, req.ApartmentId, req.Start, req.End)
		return getReservationsResponse{
			Apartments: apartments,
			Err:        err,
		}, nil
	}
}
