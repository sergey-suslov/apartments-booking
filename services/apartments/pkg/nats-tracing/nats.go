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

func NATSSubscriberTrace(tracer *zipkin.Tracer, options ...TracerOption) kitnats.SubscriberOption {
	config := tracerOptions{
		tags:      make(map[string]string),
		name:      "",
		logger:    log.NewNopLogger(),
		propagate: true,
	}

	for _, option := range options {
		option(&config)
	}

	subscriberBefore := kitnats.SubscriberBefore(func(ctx context.Context, msg *nats.Msg) context.Context {
		var (
			spanContext model.SpanContext
			name        string
		)

		if config.name != "" {
			name = config.name
		} else {
			name = msg.Subject
		}

		if config.propagate {
			spanContext = tracer.Extract(ExtractNATS(msg))
			if spanContext.Err != nil {
				_ = config.logger.Log("err", spanContext.Err)
			}
		}

		span := tracer.StartSpan(
			name,
			zipkin.Kind(model.Server),
			zipkin.Tags(config.tags),
			zipkin.Parent(spanContext),
			zipkin.FlushOnFinish(false),
		)

		return zipkin.NewContext(ctx, span)
	})

	subscriberAfter := kitnats.SubscriberAfter(
		func(ctx context.Context, conn *nats.Conn) context.Context {
			// TODO trace errors somehow
			if span := zipkin.SpanFromContext(ctx); span != nil {
				span.Finish()
			}

			return ctx
		})

	finalizer := kitnats.SubscriberFinalizer(func(ctx context.Context, msg *nats.Msg) {
		if span := zipkin.SpanFromContext(ctx); span != nil {
			span.Finish()
			span.Flush()
		}
	})

	return func(subscriber *kitnats.Subscriber) {
		subscriberBefore(subscriber)
		subscriberAfter(subscriber)
		finalizer(subscriber)
	}
}

var ErrEmptyContext = errors.New("empty request context")

type natsMessageWithContext struct {
	Sc   model.SpanContext `json:"sc"`
	Data []byte            `json:"data"`
}

// ExtractNATS will extract a span.Context from a NATS message.
func ExtractNATS(msg *nats.Msg) propagation.Extractor {
	return func() (*model.SpanContext, error) {
		var payload natsMessageWithContext
		err := json.Unmarshal(msg.Data, &payload)
		// not natsMessageWithContext
		if err != nil {
			return nil, nil
		}

		msg.Data = payload.Data

		if (model.SpanContext{}) == payload.Sc {
			return nil, ErrEmptyContext
		}

		if payload.Sc.TraceID.Empty() {
			return nil, ErrEmptyContext
		}

		return &payload.Sc, nil
	}

}

func SetName(name string) TracerOption {
	return func(o *tracerOptions) {
		o.name = name
	}
}
