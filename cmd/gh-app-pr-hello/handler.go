package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ehenry2/gh-app-pr-hello/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	// configure logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().
		Msg("lambda function starting. parsing config from environment")

	// load configuration
	ctx := context.Background()
	config, err := internal.NewConfig(ctx)
	if err != nil {
		log.Err(err).Msg("failed to read config")
		os.Exit(1)
	}
	log.Info().Msg("parsed config successfully")
	githubAppConfig := config.ToGithubAppConfig()

	// register routes
	log.Info().Msg("registering routes")
	if err := internal.RegisterGithubWebhookDispatcher(githubAppConfig); err != nil {
		log.Err(err).Msg("failed to load client creator")
		os.Exit(1)
	}
	internal.RegisterHealthCheck()
	log.Info().Msg("routes registered successfully")

	// start the lambda handler
	log.Info().Msg("starting lambda handler")
	lambda.Start(internal.AlbHandler{}.ProxyWithContext)
}
