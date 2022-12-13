package business

import (
	"context"
	"github.com/google/go-github/v47/github"
	"testing"
)

func TestPROpenHandler_Handle(t *testing.T) {
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
			name: "valid request",
			args: args{
				client: mockedGithubClient("preview your site at: http://example.com/site"),
			},
			wantErr: true,
		},
		{
			name: "github client error",
			args: args{
				client: mockedErrorGithubClient(""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PROpenHandler{}
			if err := h.Handle(tt.args.ctx, tt.args.client, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
