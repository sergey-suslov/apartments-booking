package nats_tracing

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	kitnats "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation"
	"net/http"
)

type TracerOption func(o *tracerOptions)

type tracerOptions struct {
	tags           map[string]string
	name           string
	logger         log.Logger
	propagate      bool
	requestSampler func(r *http.Request) bool
}

func NATSPublisherTrace(tracer *zipkin.Tracer, options ...TracerOption) kitnats.PublisherOption {
	config := tracerOptions{
		tags:      make(map[string]string),
		name:      "",
		logger:    log.NewNopLogger(),
		propagate: true,
	}

	for _, option := range options {
		option(&config)
	}

	publisherBefore := kitnats.PublisherBefore(func(ctx context.Context, msg *nats.Msg) context.Context {
		var (
			spanContext model.SpanContext
			name        string
		)

		if config.name != "" {
			name = config.name
		} else {
			name = msg.Subject
		}

		if parent := zipkin.SpanFromContext(ctx); parent != nil {
			spanContext = parent.Context()
		}

		span := tracer.StartSpan(
			name,
			zipkin.Kind(model.Client),
			zipkin.Tags(config.tags),
			zipkin.Parent(spanContext),
			zipkin.FlushOnFinish(false),
		)

		if config.propagate {
			if err := InjectNATS(msg)(span.Context()); err != nil {
				_ = config.logger.Log("err", err)
			}
		}

		return zipkin.NewContext(ctx, span)
	})

	publisherAfter := kitnats.PublisherAfter(func(ctx context.Context, msg *nats.Msg) context.Context {
		// TODO trace errors somehow
		if span := zipkin.SpanFromContext(ctx); span != nil {
			span.Finish()
		}

		return ctx
	})

	return func(publisher *kitnats.Publisher) {
		publisherBefore(publisher)
		publisherAfter(publisher)
	}
}

var ErrEmptyContext = errors.New("empty request context")

type natsMessageWithContext struct {
	Sc   model.SpanContext `json:"sc"`
	Data interface{}       `json:"data"`
}

// InjectGRPC will inject a span.Context into NATS message.
func InjectNATS(msg *nats.Msg) propagation.Injector {
	return func(sc model.SpanContext) error {
		if (model.SpanContext{}) == sc {
			return ErrEmptyContext
		}

		if !sc.TraceID.Empty() && sc.ID > 0 {
			messageWithContext := natsMessageWithContext{
				Sc:   sc,
				Data: msg.Data,
			}
			marshalledMessage, err := json.Marshal(&messageWithContext)
			if err != nil {
				return err
			}
			msg.Data = marshalledMessage
		}

		return nil
	}
}
