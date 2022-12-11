package internal

import (
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func RegisterGithubWebhookDispatcher(config *githubapp.Config) error {
	log.Info().Msg("registering route: github webhook dispatcher")
	cc, err := githubapp.NewDefaultCachingClientCreator(
		*config,
		githubapp.WithClientMiddleware(
			githubapp.ClientLogging(zerolog.InfoLevel)),
		githubapp.WithClientTimeout(3*time.Second))
	if err != nil {
		return err
	}
	prHandler := PRHandler{ClientCreator: cc}
	dispatcher := githubapp.NewDefaultEventDispatcher(*config, &prHandler)
	http.Handle("/default/api/github/hook", dispatcher)
	return nil
}
