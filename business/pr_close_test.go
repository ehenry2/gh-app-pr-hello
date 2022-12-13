package business

import (
	"context"
	"github.com/google/go-github/v47/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"net/http"
	"testing"
)

func mockedGithubClient(msg string) *github.Client {
	endpoint := mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber
	expected := github.IssueComment{Body: &msg}
	mockedHttpClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(endpoint, expected),
	)
	return github.NewClient(mockedHttpClient)
}

func mockedErrorGithubClient(msg string) *github.Client {
	endpoint := mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber
	httpClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatchHandler(
			endpoint,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					"github fail whale",
				)
			}),
		),
	)
	return github.NewClient(httpClient)
}

func intRef(i int) *int {
	return &i
}

func stringRef(s string) *string {
	return &s
}

func TestPRCloseHandler_Handle(t *testing.T) {
	type args struct {
		ctx    context.Context
		client *github.Client
		event  github.PullRequestEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "standard comment",
			args: args{
				ctx:    context.Background(),
				client: mockedGithubClient("your site has been cleaned up"),
				event: github.PullRequestEvent{
					Number: intRef(10),
					Repo: &github.Repository{
						Owner: &github.User{Login: stringRef("foo")},
						Name:  stringRef("bar"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "github client error",
			args: args{
				ctx:    context.Background(),
				client: mockedErrorGithubClient("your site has been cleaned up"),
				event: github.PullRequestEvent{
					Number: intRef(10),
					Repo: &github.Repository{
						Owner: &github.User{Login: stringRef("foo")},
						Name:  stringRef("bar"),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PRCloseHandler{}
			if err := h.Handle(tt.args.ctx, tt.args.client, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
