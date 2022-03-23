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

type SentryOpts struct {
	DSN string `yaml:"dsn"`
}

type TracerOpts struct {
	TracerType string     `yaml:"tracerType"`
	SentryOpts SentryOpts `yaml:"sentryOpts"`
}

func NewTracerCreator(logger *logrus.Logger, opts TracerOpts) (TracerCreator, error) {
	if opts.TracerType == "sentry" {
		tracer, err := NewSentryTracerMaker(sentry.ClientOptions{
			Dsn: opts.SentryOpts.DSN,
		})
		if err != nil {
			return nil, common.NewInitializationError("tracer", err.Error())
		}
		return tracer, nil
	}

	return NewLocalTracerMaker(logger), nil
}
