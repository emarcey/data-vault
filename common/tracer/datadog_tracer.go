package tracer

import (
	"context"
	"fmt"

	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type DataDogTracer struct {
	ctx  context.Context
	span ddtracer.Span
	opts []ddtracer.FinishOption
}

func (t *DataDogTracer) Close() {
	t.span.Finish(t.opts...)
}

func (t *DataDogTracer) Context() context.Context {
	return t.ctx
}

func (t *DataDogTracer) AddBreadcrumb(data map[string]interface{}) {
	for key, value := range data {
		t.span.SetTag(key, value)
	}
}

func (t *DataDogTracer) CaptureException(err error) {
	t.opts = append(t.opts, ddtracer.WithError(err))
}

func (t *DataDogTracer) StartChild(operation string) Tracer {
	return NewDataDogTracer(t.Context(), operation, ddtracer.ChildOf(t.span.Context()))
}

func NewDataDogTracer(ctx context.Context, operation string, spanOpts ...ddtracer.StartSpanOption) Tracer {
	span := ddtracer.StartSpan(operation, spanOpts...)

	return &DataDogTracer{ctx: ctx, span: span}
}

func NewDataDogTracerMaker(ctx context.Context, opts DataDogOpts) (TracerCreator, error) {
	ddtracer.Start(ddtracer.WithAgentAddr(fmt.Sprintf("%v:%v", opts.AgentAddr, opts.AgentPort)))

	go func() {
		select {
		case <-ctx.Done():
			ddtracer.Stop()
			return
		}
	}()

	return func(ctx context.Context, operation string) Tracer {
		return NewDataDogTracer(ctx, operation)
	}, nil
}
