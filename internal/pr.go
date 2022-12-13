package internal

import (
	"context"
	"encoding/json"
	"github.com/ehenry2/gh-app-pr-hello/business"
	"github.com/google/go-github/v47/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog/log"
)

const (
	OpenedAction = "opened"
	ClosedAction = "closed"
)

type PRHandler struct {
	ClientCreator githubapp.ClientCreator
	OpenHandler   *business.PROpenHandler
	CloseHandler  *business.PRCloseHandler
}

func (h *PRHandler) Handles() []string {
	return []string{"pull_request"}
}

func (h *PRHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	log.Info().Msg("handling github event")
	log.Info().Msg(string(payload))
	// parse json body of event.
	var event github.PullRequestEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Err(err).Msg("failed to decode json")
		return err
	}

	// create github api client to use in posting the comment.
	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := h.ClientCreator.NewInstallationClient(installationID)
	if err != nil {
		log.Err(err).Msg("failed to create installation client")
		return err
	}

	// handle the event.
	switch *event.Action {
	case OpenedAction:
		return h.OpenHandler.Handle(ctx, client, event)
	case ClosedAction:
		return h.CloseHandler.Handle(ctx, client, event)
	}

	return err
}
