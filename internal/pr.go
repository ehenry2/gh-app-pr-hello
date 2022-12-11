package internal

import (
	"context"
	"encoding/json"
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

	msg := ""
	switch *event.Action {
	case OpenedAction:
		msg = "creating a new site: url is: https://example.com/random1234"
	case ClosedAction:
		msg = "cleaning up resources..."
	}
	installationID := githubapp.GetInstallationIDFromEvent(&event)
	repo := event.GetRepo()
	repoName := repo.GetName()
	prNum := event.GetNumber()
	repoOwner := repo.GetOwner().GetLogin()
	comment := &github.IssueComment{
		Body: &msg,
	}

	// create github api client to use in posting the comment.
	client, err := h.ClientCreator.NewInstallationClient(installationID)
	if err != nil {
		log.Err(err).Msg("failed to create installation client")
		return err
	}
	
	_, _, err = client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, comment)
	if err != nil {
		log.Err(err).Msg("failed to create comment")
	}
	return err
}
