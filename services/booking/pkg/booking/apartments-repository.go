package booking

import (
	nats_tracing "booking/pkg/nats-tracing"
	"context"
	"encoding/json"
	"errors"
	natstransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/openzipkin/zipkin-go"
)

const getApartmentByIdSubject = "apartments.getApartmentById"

var coldNotGetResponseFromApartment = errors.New("could not get response from the apartment service, wrong response format")

type apartmentsRepository struct {
	nc     *nats.Conn
	tracer *zipkin.Tracer
}

func NewApartmentsRepository(nc *nats.Conn, tracer *zipkin.Tracer) *apartmentsRepository {
	return &apartmentsRepository{nc: nc, tracer: tracer}
}

type getApartmentByIdRequest struct {
	ApartmentId string `json:"apartmentId"`
}

type getApartmentByIdResponse struct {
	Apartment Apartment `json:"apartment"`
	Err       error     `json:"err,omitempty"`
}

func (a *apartmentsRepository) GetApartmentById(ctx context.Context, apartmentId string) (*Apartment, error) {
	publisher := natstransport.NewPublisher(a.nc, getApartmentByIdSubject, natstransport.EncodeJSONRequest, decodeGetApartmentById, nats_tracing.NATSPublisherTrace(a.tracer, nats_tracing.SetName("book an apartment")))
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
