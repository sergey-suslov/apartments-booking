package booking

import (
	"context"
	"encoding/json"
	"errors"

	natstransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/openzipkin/zipkin-go"
	"github.com/sergey-suslov/go-kit-nats-zipkin-tracing/natszipkin"
)

const getApartmentByIDSubject = "apartments.getApartmentById"

var ErrColdNotGetResponseFromApartment = errors.New("could not get response from the apartment service, wrong response format")

type ApartmentsRepositoryNATS struct {
	nc     *nats.Conn
	tracer *zipkin.Tracer
}

func NewApartmentsRepository(nc *nats.Conn, tracer *zipkin.Tracer) *ApartmentsRepositoryNATS {
	return &ApartmentsRepositoryNATS{nc: nc, tracer: tracer}
}

type getApartmentByIDRequest struct {
	ApartmentID string `json:"apartmentId"`
}

type getApartmentByIDResponse struct {
	Apartment Apartment `json:"apartment"`
	Err       error     `json:"err,omitempty"`
}

func (a *ApartmentsRepositoryNATS) GetApartmentByID(ctx context.Context, apartmentID string) (*Apartment, error) {
	publisher := natstransport.NewPublisher(
		a.nc,
		getApartmentByIDSubject,
		natstransport.EncodeJSONRequest,
		decodeGetApartmentByID,
		natszipkin.NATSPublisherTrace(a.tracer, natszipkin.Name("book an apartment")),
	)
	res, err := publisher.Endpoint()(ctx, getApartmentByIDRequest{ApartmentID: apartmentID})
	if err != nil {
		return nil, err
	}
	response, ok := res.(getApartmentByIDResponse)
	if !ok {
		return nil, ErrColdNotGetResponseFromApartment
	}
	if response.Apartment.ID == "" {
		return nil, ErrNoApartmentWithGivenID
	}
	return &response.Apartment, nil
}

func decodeGetApartmentByID(_ context.Context, msg *nats.Msg) (response interface{}, err error) {
	var res getApartmentByIDResponse
	err = json.Unmarshal(msg.Data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
