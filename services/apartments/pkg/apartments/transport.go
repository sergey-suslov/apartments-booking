package apartments

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/circuitbreaker"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	kitnats "github.com/go-kit/kit/transport/nats"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/sony/gobreaker"

	"net/http"
)

const queueName = "apartments"
const getApartmentByIdSubject = "apartments.getApartmentById"

const SpanCtxKey = "SpanCtxKey"

func MakeHttpHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	endpoint := makeGetApartmentsEndpoint(s)
	endpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(endpoint)
	getApartmentsHandler := kithttp.NewServer(endpoint, decodeGetApartmentsRequest, encodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/apartments", getApartmentsHandler).Methods("GET")

	return r
}

func decodeGetApartmentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getApartmentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(Errorer); ok && e.Error() != nil {
		encodeError(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func MakeNatsHandler(s Service, nc *nats.Conn) {
	apartmentByIdEndpoint := makeGetApartmentByIdEndpoint(s)
	subscriber := kitnats.NewSubscriber(
		apartmentByIdEndpoint,
		decodeGetApartmentByIdRequest,
		kitnats.EncodeJSONResponse,

		// turn on zipkin context parsing
		kitnats.SubscriberBefore(DecodeSpanContext),
	)
	_, err := nc.QueueSubscribe(getApartmentByIdSubject, queueName, subscriber.ServeMsg(nc))
	if err != nil {
		panic(err)
	}
}

// DecodeSpanContext parse zipkin context from NATS message. Use as NewSubscriber option, for ex. nats.SubscriberBefore(DecodeSpanContext)
func DecodeSpanContext(ctx context.Context, msg *nats.Msg) context.Context {
	var payload natsPayload
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ctx
	}

	marshaled, _ := json.Marshal(payload.Data)
	msg.Data = marshaled

	if payload.SpanContext.ID != 0 {
		return context.WithValue(ctx, SpanCtxKey, payload.SpanContext)
	}
	return ctx
}

func decodeGetApartmentByIdRequest(_ context.Context, msg *nats.Msg) (request interface{}, err error) {
	var getApartmentByIdRequest getApartmentByIdRequest
	err = json.Unmarshal(msg.Data, &getApartmentByIdRequest)
	if err != nil {
		return nil, err
	}

	return getApartmentByIdRequest, nil
}
