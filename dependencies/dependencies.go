package dependencies

import (
	"context"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/emarcey/data-vault/common"
	"github.com/emarcey/data-vault/common/logger"
	"github.com/emarcey/data-vault/common/tracer"
	"github.com/emarcey/data-vault/database"
	"github.com/emarcey/data-vault/dependencies/secrets"
)

type ServerConfigs struct {
	AccessTokenHours   int `yaml:"accessTokenHours"`
	DataRefreshSeconds int `yaml:"dataRefreshSeconds"`
}

type DependenciesInitOpts struct {
	HttpAddr           string                     `yaml:"httpAddr"`
	LoggerType         string                     `yaml:"loggerType"`
	SecretsManagerOpts secrets.SecretsManagerOpts `yaml:"secretsManagerOpts"`
	DatabaseOpts       database.DatabaseOpts      `yaml:"databaseOpts"`
	TracerOpts         tracer.TracerOpts          `yaml:"tracerOpts"`
	Env                string                     `yaml:"env"`
	Version            string                     `yaml:"version"`
	ServerConfigs      *ServerConfigs             `yaml:"serverConfigs"`
}

type Dependencies struct {
	Env            string
	Logger         *logrus.Logger
	Tracer         tracer.TracerCreator
	SecretsManager secrets.SecretsManager
	Database       *database.DatabaseEngine
	AuthUsers      *UserCache
	AccessTokens   *AccessTokenCache
	ServerConfigs  *ServerConfigs
}

func ReadOpts(filename string) (DependenciesInitOpts, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return DependenciesInitOpts{}, common.NewInitializationError("dependency-options", "Unable to read server options for file, %s, with error: %v", filename, err)
	}
	var opts DependenciesInitOpts
	err = yaml.Unmarshal(raw, &opts)
	if err != nil {
		return DependenciesInitOpts{}, common.NewInitializationError("dependency-options", "Unable to unmarshal file, %s, with error: %v", filename, err)
	}
	return opts, nil
}

func MakeDependencies(ctx context.Context, opts DependenciesInitOpts) (*Dependencies, error) {
	logger, err := logger.MakeLogger(opts.LoggerType, opts.Env)
	if err != nil {
		return nil, err
	}
	tracer, err := tracer.NewTracerCreator(ctx, logger, opts.TracerOpts)
	if err != nil {
		return nil, err
	}

	secretsManager, err := secrets.NewSecretsManager(ctx, opts.SecretsManagerOpts)
	if err != nil {
		return nil, err
	}
	db, err := database.NewDatabase(logger, tracer, opts.DatabaseOpts)
	if err != nil {
		return nil, err
	}

	authUsers, err := NewUserCache(ctx, logger, db, opts.ServerConfigs.DataRefreshSeconds)
	if err != nil {
		return nil, err
	}

	accessTokens, err := NewAccessTokenCache(ctx, logger, db, opts.ServerConfigs.DataRefreshSeconds)
	if err != nil {
		return nil, err
	}
	deps := &Dependencies{
		Env:            opts.Env,
		Logger:         logger,
		Tracer:         tracer,
		SecretsManager: secretsManager,
		Database:       db,
		AuthUsers:      authUsers,
		AccessTokens:   accessTokens,
		ServerConfigs:  opts.ServerConfigs,
	}
	return deps, nil
}
