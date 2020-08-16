package booking

import (
	"context"
	"github.com/openzipkin/zipkin-go"
	"strconv"
	"time"
)

type tracingService struct {
	zipkinTracer *zipkin.Tracer
	s            Service
}

func NewTracingService(zipkinTracer *zipkin.Tracer, s Service) *tracingService {
	return &tracingService{zipkinTracer: zipkinTracer, s: s}
}

func (t *tracingService) GetReservations(ctx context.Context, apartmentId string, start, end time.Time) (out []Reservation, err error) {
	return t.s.GetReservations(ctx, apartmentId, start, end)
}

func (t *tracingService) BookApartment(ctx context.Context, userId, apartmentId string, start, end time.Time) (out *Reservation, err error) {
	span, spanCtx := t.zipkinTracer.StartSpanFromContext(ctx, "book apartment")
	span.Tag("method", "BookApartment")
	span.Annotate(time.Now(), "start")

	apartment, err := t.s.BookApartment(context.WithValue(spanCtx, "spanContext", span.Context()), userId, apartmentId, start, end)
	span.Tag("error", strconv.FormatBool(err != nil))
	span.Annotate(time.Now(), "finish")
	span.Finish()
	return apartment, err
}
