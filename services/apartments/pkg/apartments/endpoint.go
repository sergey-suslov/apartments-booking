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

type getApartmentByIdRequest struct {
	ApartmentId string `json:"apartmentId"`
}

type getApartmentByIdResponse struct {
	Apartment *Apartment `json:"apartment"`
	Err       error      `json:"err,omitempty"`
}

func makeGetApartmentByIdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getApartmentByIdRequest)
		apartment, err := s.GetApartmentById(ctx, req.ApartmentId)
		return getApartmentByIdResponse{
			Apartment: apartment,
			Err:       err,
		}, nil
	}
}
