package business

import (
	"context"
	"github.com/google/go-github/v47/github"
)

type PROpenHandler struct{}

func (h *PROpenHandler) Handle(ctx context.Context, client *github.Client, event github.PullRequestEvent) error {
	msg := "preview your site at: http://example.com/site"
	repo := event.GetRepo()
	repoName := repo.GetName()
	prNum := event.GetNumber()
	repoOwner := repo.GetOwner().GetLogin()
	comment := &github.IssueComment{
		Body: &msg,
	}
	_, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, comment)
	return err
}
