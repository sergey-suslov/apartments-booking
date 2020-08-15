package booking

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/circuitbreaker"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/sony/gobreaker"

	"net/http"
)

func MakeHttpHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getApartmentsEndpoint := makeGetApartmentsEndpoint(s)
	getApartmentsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getApartmentsEndpoint)
	getReservationsHandler := kithttp.NewServer(getApartmentsEndpoint, decodeGetApartmentsRequest, encodeResponse, opts...)

	bookApartmentEndpoint := makeBookApartmentEndpoint(s)
	bookApartmentEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(bookApartmentEndpoint)
	bookApartmentHandler := kithttp.NewServer(bookApartmentEndpoint, DefaultRequestDecoder(decodeBookApartmentRequest), encodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/reservations", getReservationsHandler).Methods("GET")
	r.Handle("/reservations", bookApartmentHandler).Methods("POST")

	return r
}

func decodeGetApartmentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getReservationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeBookApartmentRequest(r *http.Request) (UserClaimable, error) {
	var req bookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
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
	case wrongIdFormat:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

type UserClaimable interface {
	SetUserClaim(claim UserClaim)
}

func DefaultRequestDecoder(decoder func(r *http.Request) (UserClaimable, error)) func(_ context.Context, r *http.Request) (interface{}, error) {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		userClaim, err := GetUserClaimFromRequest(r)
		if err != nil {
			return nil, err
		}

		request, err := decoder(r)
		if err != nil {
			return nil, err
		}
		request.SetUserClaim(*userClaim)
		return request, nil
	}
}
