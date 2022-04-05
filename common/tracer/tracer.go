package tracer

import (
	"context"
	"fmt"

	sentry "github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"

	"github.com/emarcey/data-vault/common"
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

type DataDogOpts struct {
	AgentAddr string `yaml:"agentAddr"`
	AgentPort string `yaml:"agentPort"`
}

type TracerOpts struct {
	TracerType  string      `yaml:"tracerType"`
	SentryOpts  SentryOpts  `yaml:"sentryOpts"`
	DataDogOpts DataDogOpts `yaml:"dataDogOpts"`
}

func NewTracerCreator(ctx context.Context, logger *logrus.Logger, opts TracerOpts) (TracerCreator, error) {
	switch opts.TracerType {
	case "sentry":
		tracer, err := NewSentryTracerMaker(sentry.ClientOptions{
			Dsn: opts.SentryOpts.DSN,
		})
		if err != nil {
			return nil, common.NewInitializationError("tracer", err.Error())
		}
		return tracer, nil
	case "local":
		return NewLocalTracerMaker(logger), nil
	case "datadog":
		return NewDataDogTracerMaker(ctx, opts.DataDogOpts)
	case "noop", "":
		return NewNoOpTracerMaker(), nil
	default:
		return nil, common.NewInitializationError("tracer", "Invalid tracer type: %v", opts.TracerType)
	}
}
