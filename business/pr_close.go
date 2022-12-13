package business

import (
	"context"
	"github.com/google/go-github/v47/github"
)

type PRCloseHandler struct{}

func (h *PRCloseHandler) Handle(ctx context.Context, client *github.Client, event github.PullRequestEvent) error {
	msg := "your site has been cleaned up"
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
