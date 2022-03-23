package dependencies

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"emarcey/data-vault/common"
	"emarcey/data-vault/common/logger"
	"emarcey/data-vault/common/tracer"
	"emarcey/data-vault/database"
	"emarcey/data-vault/dependencies/secrets"
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
	Logger         *logrus.Logger
	Tracer         tracer.TracerCreator
	SecretsManager secrets.SecretsManager
	Database       *database.DatabaseEngine
	AuthUsers      map[string]*common.User
	AccessTokens   map[string]*common.AccessToken
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
	tracer, err := tracer.NewTracerCreator(logger, opts.TracerOpts)
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

	authUsers, err := database.SelectUsersForAuth(ctx, db)
	if err != nil {
		return nil, err
	}

	accessTokens, err := database.SelectAccessTokensForAuth(ctx, db)
	if err != nil {
		return nil, err
	}
	deps := &Dependencies{
		Logger:         logger,
		Tracer:         tracer,
		SecretsManager: secretsManager,
		Database:       db,
		AuthUsers:      authUsers,
		AccessTokens:   accessTokens,
		ServerConfigs:  opts.ServerConfigs,
	}

	timer := time.NewTicker(time.Duration(opts.ServerConfigs.DataRefreshSeconds) * time.Second)
	go func() {
		for true {
			select {
			case <-ctx.Done():
				logger.Warn("Context cancelled. No more refreshes.")
				return
			case <-timer.C:
				authUsers, err := database.SelectUsersForAuth(ctx, db)
				if err != nil {
					logger.Errorf("Error in SelectUsersForAuth refresh: %v", err)
				}
				deps.AuthUsers = authUsers
				accessTokens, err := database.SelectAccessTokensForAuth(ctx, db)
				if err != nil {
					logger.Errorf("Error in SelectAccessTokensForAuth refresh: %v", err)
				}
				deps.AccessTokens = accessTokens
			}
		}
	}()
	return deps, nil
}
