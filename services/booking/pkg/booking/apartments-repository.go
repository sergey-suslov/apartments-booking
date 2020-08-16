package booking

import (
	"context"
	"encoding/json"
	"errors"
	natstransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/openzipkin/zipkin-go/model"
)

const getApartmentByIdSubject = "apartments.getApartmentById"

var coldNotGetResponseFromApartment = errors.New("could not get response from the apartment service, wrong response format")

type apartmentsRepository struct {
	nc *nats.Conn
}

func NewApartmentsRepository(nc *nats.Conn) *apartmentsRepository {
	return &apartmentsRepository{nc: nc}
}

type natsPayload struct {
	SpanContext model.SpanContext `json:"spanContext"`
	Data        interface{}       `json:"data"`
}

type getApartmentByIdRequest struct {
	ApartmentId string `json:"apartmentId"`
}

type getApartmentByIdResponse struct {
	Apartment Apartment `json:"apartment"`
	Err       error     `json:"err,omitempty"`
}

func (a *apartmentsRepository) GetApartmentById(ctx context.Context, apartmentId string) (*Apartment, error) {
	publisher := natstransport.NewPublisher(a.nc, getApartmentByIdSubject, encodeWithContext, decodeGetApartmentById)
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

func encodeWithContext(ctx context.Context, msg *nats.Msg, request interface{}) error {
	value := ctx.Value("spanContext")
	var spanCtx model.SpanContext
	if value != nil {
		spanCtx = value.(model.SpanContext)
	}
	b, err := json.Marshal(&natsPayload{
		SpanContext: spanCtx,
		Data:        request,
	})
	if err != nil {
		return err
	}

	msg.Data = b

	return nil
}
