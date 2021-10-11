package tracer

import (
	"context"

	sentry "github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"

	"emarcey/data-vault/common"
)

type Tracer interface {
	Close()
	Context() context.Context
	AddBreadcrumb(data map[string]interface{})
	CaptureException(err error)
	StartChild(operation string) Tracer
}

type TracerCreator func(ctx context.Context, operation string) Tracer

type NewTracerCreatorOpts struct {
	Env        string
	Logger     *logrus.Logger
	SentryOpts sentry.ClientOptions
}

func NewTracerCreator(opts NewTracerCreatorOpts) (TracerCreator, error) {
	if opts.Env == "local" {
		return NewLocalTracerMaker(opts.Logger), nil
	}
	tracer, err := NewSentryTracerMaker(opts.SentryOpts)
	if err != nil {
		return nil, common.NewInitializationError("tracer", err.Error())
	}
	return tracer, nil
}
