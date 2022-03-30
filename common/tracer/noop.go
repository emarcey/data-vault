package tracer

import (
	"context"
)

type NoOpTracer struct {
}

func (t *NoOpTracer) Close() {
	return
}

func (t *NoOpTracer) Context() context.Context {
	return context.Background()
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

func NewNoOpTracer() Tracer {
	return &NoOpTracer{}
}

func NewNoOpTracerMaker() TracerCreator {
	return func(_ context.Context, _ string) Tracer {
		return NewNoOpTracer()
	}
}
