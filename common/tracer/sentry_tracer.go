package tracer

import (
	"context"
	"time"

	sentry "github.com/getsentry/sentry-go"
)

type SentryTracer struct {
	span *sentry.Span
	hub  *sentry.Hub
}

func (t *SentryTracer) Close() {
	t.span.Finish()
}

func (t *SentryTracer) Context() context.Context {
	return t.span.Context()
}

func (t *SentryTracer) AddBreadcrumb(data map[string]interface{}) {
	t.hub.AddBreadcrumb(&sentry.Breadcrumb{
		Data:      data,
		Timestamp: time.Now(),
	}, nil)
}

func (t *SentryTracer) CaptureException(err error) {
	t.hub.CaptureException(err)
}

func (t *SentryTracer) StartChild(operation string) Tracer {
	return NewSentryTracer(t.span.Context(), operation)
}

func NewSentryTracer(ctx context.Context, operation string) Tracer {
	hub := sentry.CurrentHub().Clone()
	ctx = sentry.SetHubOnContext(ctx, hub)

	spanOpts := []sentry.SpanOption{}
	tx := sentry.TransactionFromContext(ctx)
	if tx != nil {
		spanOpts = append(spanOpts, sentry.TransactionName(operation))
	}

	span := sentry.StartSpan(ctx, operation, spanOpts...)

	return &SentryTracer{span: span, hub: hub}
}

func NewSentryTracerMaker(opts sentry.ClientOptions) (TracerCreator, error) {
	err := sentry.Init(opts)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, operation string) Tracer {
		return NewSentryTracer(ctx, operation)
	}, nil
}
