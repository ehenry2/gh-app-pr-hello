package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ehenry2/gh-app-pr-hello/internal"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().
		Msg("lambda function starting. parsing config from environment")
	ctx := context.Background()
	appConfig, err := internal.NewConfig(ctx)
	if err != nil {
		log.Err(err).Msg("failed to read config")
		os.Exit(1)
	}
	log.Info().Msg("parsed config successfully")

	clientConfig := githubapp.Config{
		V3APIURL: appConfig.GithubV3Endpoint,
		App: struct {
			IntegrationID int64  `yaml:"integration_id" json:"integrationId"`
			WebhookSecret string `yaml:"webhook_secret" json:"webhookSecret"`
			PrivateKey    string `yaml:"private_key" json:"privateKey"`
		}{
			IntegrationID: appConfig.IntegrationID,
			WebhookSecret: appConfig.WebhookSecret,
			PrivateKey:    string(appConfig.PrivateKey),
		},
	}
	log.Info().Msg("registering routes")
	if err := internal.RegisterGithubWebhookDispatcher(&clientConfig); err != nil {
		log.Err(err).Msg("failed to load client creator")
		os.Exit(1)
	}
	internal.RegisterHealthCheck()
	log.Info().Msg("routes registered successfully")
	log.Info().Msg("starting lambda handler")
	lambda.Start(internal.AlbHandler{}.ProxyWithContext)
}
