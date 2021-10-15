package dependencies

import (
	"context"
	"io/ioutil"

	sentry "github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"emarcey/data-vault/common"
	"emarcey/data-vault/common/logger"
	"emarcey/data-vault/common/tracer"
	"emarcey/data-vault/db"
	"emarcey/data-vault/dependencies/secrets"
)

type DependenciesInitOpts struct {
	LoggerType         string                     `yaml:"loggerType"`
	SecretsManagerOpts secrets.SecretsManagerOpts `yaml:"secretsManagerOpts"`
	DatabaseOpts       db.DatabaseOpts            `yaml:"databaseOpts"`
	Env                string                     `yaml:"env"`
}
type Dependencies struct {
	Logger         *logrus.Logger
	Tracer         tracer.TracerCreator
	SecretsManager secrets.SecretsManager
	Database       *db.Database
}

func ReadOpts(filename string) (DependenciesInitOpts, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return DependenciesInitOpts{}, common.NewInitializationError("read-options", "Unable to read server options for file, %s, with error: %v", filename, err)
	}
	var opts DependenciesInitOpts
	err = yaml.Unmarshal(raw, &opts)
	return opts, nil
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
	db, err := db.NewDatabase(logger, tracer, opts.DatabaseOpts)
	if err != nil {
		return nil, err
	}
	return &Dependencies{
		Logger:         logger,
		Tracer:         tracer,
		SecretsManager: secretsManager,
		Database:       db,
	}, nil
}
