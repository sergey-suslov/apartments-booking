package apartments

import (
	"context"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"strconv"
	"time"
)

const zipkinTracerKey = "zipkinTracerKey"

type tracingService struct {
	zipkinTracer *zipkin.Tracer
	s            Service
}

func NewTracingService(zipkinTracer *zipkin.Tracer, s Service) *tracingService {
	return &tracingService{zipkinTracer: zipkinTracer, s: s}
}

func (t *tracingService) GetApartments(ctx context.Context, city City, limit, offset int) ([]Apartment, error) {
	return t.s.GetApartments(ctx, city, limit, offset)
}

func (t *tracingService) GetApartmentById(ctx context.Context, apartmentId string) (*Apartment, error) {
	span, spanCtx := t.startSpanFromContext(ctx, "get apartment by id")
	span.Tag("method", "GetApartmentById")
	span.Annotate(time.Now(), "start")

	apartment, err := t.s.GetApartmentById(spanCtx, apartmentId)
	span.Tag("error", strconv.FormatBool(err != nil))
	span.Annotate(time.Now(), "finish")
	span.Finish()
	return apartment, err
}

func (t *tracingService) startSpanFromContext(ctx context.Context, name string) (zipkin.Span, context.Context) {
	rawSpan := ctx.Value(SpanCtxKey)
	var spanCtx model.SpanContext
	if rawSpan != nil {
		spanCtx = rawSpan.(model.SpanContext)
	}
	span, newCtx := t.zipkinTracer.StartSpanFromContext(ctx, name, zipkin.Parent(spanCtx))
	return span, context.WithValue(newCtx, zipkinTracerKey, t.zipkinTracer)
}

func StartSpanFromContext(ctx context.Context, name string) (zipkin.Span, context.Context) {
	value := ctx.Value(zipkinTracerKey)
	if value == nil {
		return nil, ctx
	}

	tracer := value.(*zipkin.Tracer)
	span, spanCtx := tracer.StartSpanFromContext(ctx, name)
	return span, spanCtx
}
