package dependencies

import (
	"context"

	sentry "github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"

	"emarcey/data-vault/dependencies/logger"
	"emarcey/data-vault/dependencies/secrets"
	"emarcey/data-vault/dependencies/tracer"
)

type DependenciesInitOpts struct {
	LoggerType         string
	SecretsManagerOpts secrets.NewSecretsManagerOpts
	Env                string
}
type Dependencies struct {
	Logger         *logrus.Logger
	Tracer         tracer.TracerCreator
	SecretsManager secrets.SecretsManager
}

func MakeDependencies(ctx context.Context, opts DependenciesInitOpts) (*Dependencies, error) {
	logger, err := logger.MakeLogger(opts.LoggerType, opts.Env)
	if err != nil {
		return nil, err
	}

	tracerOpts := tracer.NewTracerCreatorOpts{
		Env:        opts.Env,
		Logger:     logger,
		SentryOpts: sentry.ClientOptions{},
	}
	tracer, err := tracer.NewTracerCreator(tracerOpts)
	if err != nil {
		return nil, err
	}

	secretsManager, err := secrets.NewSecretsManager(ctx, opts.SecretsManagerOpts)
	if err != nil {
		return nil, err
	}
	return &Dependencies{
		Logger:         logger,
		Tracer:         tracer,
		SecretsManager: secretsManager,
	}, nil
}
