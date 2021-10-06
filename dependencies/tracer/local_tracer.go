package tracer

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type LocalTracer struct {
	logger    *logrus.Logger
	ctx       context.Context
	operation string
}

func (t *LocalTracer) Close() {
	return
}

func (t *LocalTracer) Context() context.Context {
	return t.ctx
}

func (t *LocalTracer) AddBreadcrumb(data map[string]interface{}) {
	t.logger.Debugf("Tracing Op: %s. Adding breadcrumb data %v", t.operation, data)
}

func (t *LocalTracer) CaptureException(err error) {
	t.logger.Errorf("Tracing Op: %s. Error captured: %s", t.operation, err.Error())
}

func (t *LocalTracer) StartChild(operation string) Tracer {
	return NewLocalTracer(t.logger, t.Context(), fmt.Sprintf("%s - %s", t.operation, operation))
}

func NewLocalTracer(logger *logrus.Logger, ctx context.Context, operation string) Tracer {
	return &LocalTracer{logger: logger, ctx: ctx, operation: operation}
}

func NewLocalTracerMaker(logger *logrus.Logger) TracerCreator {
	return func(ctx context.Context, operation string) Tracer {
		return NewLocalTracer(logger, ctx, operation)
	}
}
