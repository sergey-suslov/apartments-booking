package booking

import (
	"context"
	"encoding/json"
	"errors"
	natstransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
)

const getApartmentByIdSubject = "apartments.getApartmentById"

var coldNotGetResponseFromApartment = errors.New("could not get response from the apartment service, wrong response format")

type apartmentsRepository struct {
	nc *nats.Conn
}

func NewApartmentsRepository(nc *nats.Conn) *apartmentsRepository {
	return &apartmentsRepository{nc: nc}
}

type getApartmentByIdRequest struct {
	ApartmentId string `json:"apartmentId"`
}

type getApartmentByIdResponse struct {
	Apartment Apartment `json:"apartment"`
	Err       error     `json:"err,omitempty"`
}

func (a *apartmentsRepository) GetApartmentById(ctx context.Context, apartmentId string) (*Apartment, error) {
	publisher := natstransport.NewPublisher(a.nc, getApartmentByIdSubject, natstransport.EncodeJSONRequest, decodeGetApartmentById)
	res, err := publisher.Endpoint()(ctx, getApartmentByIdRequest{ApartmentId: apartmentId})
	if err != nil {
		return nil, err
	}
	getApartmentByIdResponse, ok := res.(getApartmentByIdResponse)
	if !ok {
		return nil, coldNotGetResponseFromApartment
	}
	if getApartmentByIdResponse.Apartment.ID == "" {
		return nil, NoApartmentWithGivenId
	}
	return &getApartmentByIdResponse.Apartment, nil
}

func decodeGetApartmentById(_ context.Context, msg *nats.Msg) (response interface{}, err error) {
	var res getApartmentByIdResponse
	err = json.Unmarshal(msg.Data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
