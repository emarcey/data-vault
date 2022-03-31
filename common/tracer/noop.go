package tracer

import (
	"context"
)

type NoOpTracer struct {
	ctx context.Context
}

func (t *NoOpTracer) Close() {
	return
}

func (t *NoOpTracer) Context() context.Context {
	return t.ctx
}

func (t *NoOpTracer) AddBreadcrumb(_ map[string]interface{}) {
	return
}

func (t *NoOpTracer) CaptureException(_ error) {
	return
}

func (t *NoOpTracer) StartChild(_ string) Tracer {
	return t
}

func NewNoOpTracer(ctx context.Context) Tracer {
	return &NoOpTracer{
		ctx: ctx,
	}
}

func NewNoOpTracerMaker() TracerCreator {
	return func(ctx context.Context, _ string) Tracer {
		return NewNoOpTracer(ctx)
	}
}
