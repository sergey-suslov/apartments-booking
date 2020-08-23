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
	"github.com/openzipkin/zipkin-go"
	"github.com/sergey-suslov/go-kit-nats-zipkin-tracing/natszipkin"
	"github.com/sony/gobreaker"

	"net/http"
)

const queueName = "apartments"
const getApartmentByIDSubject = "apartments.getApartmentById"

func MakeHTTPHandler(s Service, logger kitlog.Logger) http.Handler {
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

func MakeNatsHandler(s Service, nc *nats.Conn, tracer *zipkin.Tracer) {
	apartmentByIDEndpoint := makeGetApartmentByIDEndpoint(s)
	subscriber := kitnats.NewSubscriber(
		apartmentByIDEndpoint,
		decodeGetApartmentByIDRequest,
		kitnats.EncodeJSONResponse,
		natszipkin.NATSSubscriberTrace(tracer, natszipkin.Name("get apartments by id")),
	)
	_, err := nc.QueueSubscribe(getApartmentByIDSubject, queueName, subscriber.ServeMsg(nc))
	if err != nil {
		panic(err)
	}
}

func decodeGetApartmentByIDRequest(_ context.Context, msg *nats.Msg) (request interface{}, err error) {
	var getApartmentByIDRequest getApartmentByIDRequest
	err = json.Unmarshal(msg.Data, &getApartmentByIDRequest)
	if err != nil {
		return nil, err
	}

	return getApartmentByIDRequest, nil
}
