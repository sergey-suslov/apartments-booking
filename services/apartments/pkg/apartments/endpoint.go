package apartments

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Errorer interface {
	Error() error
}

type getApartmentsRequest struct {
	City   City `json:"city"`
	Limit  int  `json:"limit"`
	Offset int  `json:"offset"`
}

type getApartmentsResponse struct {
	Apartments []Apartment `json:"apartments"`
	Err        error       `json:"error,omitempty"`
}

func (g getApartmentsResponse) Error() error {
	return g.Err
}

func makeGetApartmentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getApartmentsRequest)
		apartments, err := s.GetApartments(ctx, req.City, req.Limit, req.Offset)
		return getApartmentsResponse{
			Apartments: apartments,
			Err:        err,
		}, nil
	}
}

type getApartmentByIDRequest struct {
	ApartmentID string `json:"apartmentId"`
}

type getApartmentByIDResponse struct {
	Apartment *Apartment `json:"apartment"`
	Err       error      `json:"err,omitempty"`
}

func makeGetApartmentByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getApartmentByIDRequest)
		apartment, err := s.GetApartmentByID(ctx, req.ApartmentID)
		return getApartmentByIDResponse{
			Apartment: apartment,
			Err:       err,
		}, nil
	}
}
