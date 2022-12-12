package internal

import (
	"context"
	"fmt"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig_ToGithubAppConfig(t *testing.T) {
	webSecret := "secret"
	privKey := "supersecret"
	endpoint := "http://example.com/api"
	integId := int64(10)
	type fields struct {
		IntegrationID    int64
		WebhookSecret    string
		PrivateKeyBytes  []byte
		GithubV3Endpoint string
		PrivateKey       string
	}
	tests := []struct {
		name   string
		fields fields
		want   *githubapp.Config
	}{
		{
			name: "valid conversion",
			fields: fields{
				IntegrationID:    integId,
				WebhookSecret:    webSecret,
				PrivateKey:       privKey,
				PrivateKeyBytes:  []byte(privKey),
				GithubV3Endpoint: endpoint,
			},
			want: &githubapp.Config{
				V3APIURL: endpoint,
				App: struct {
					IntegrationID int64  `yaml:"integration_id" json:"integrationId"`
					WebhookSecret string `yaml:"webhook_secret" json:"webhookSecret"`
					PrivateKey    string `yaml:"private_key" json:"privateKey"`
				}{
					IntegrationID: integId,
					WebhookSecret: webSecret,
					PrivateKey:    privKey,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				IntegrationID:    tt.fields.IntegrationID,
				WebhookSecret:    tt.fields.WebhookSecret,
				PrivateKeyBytes:  tt.fields.PrivateKeyBytes,
				GithubV3Endpoint: tt.fields.GithubV3Endpoint,
				PrivateKey:       tt.fields.PrivateKey,
			}
			assert.Equalf(t, tt.want, c.ToGithubAppConfig(), "ToGithubAppConfig()")
		})
	}
}

func TestNewConfig(t *testing.T) {
	type args struct {
		ctx context.Context
		env map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid config",
			args: args{
				ctx: context.Background(),
				env: map[string]string{
					"GITHUB_INTEGRATION_ID": "10",
					"GITHUB_WEBHOOK_SECRET": "webhook",
					"GITHUB_PRIVATE_KEY":    "c2VjcmV0",
					"GITHUB_V3_ENDPOINT":    "http://example.com/api",
				},
			},
			want: &Config{
				IntegrationID:    10,
				WebhookSecret:    "webhook",
				PrivateKeyBytes:  []byte("c2VjcmV0"),
				GithubV3Endpoint: "http://example.com/api",
				PrivateKey:       "secret",
			},
			wantErr: assert.NoError,
		},
		{
			name: "key is not base64 encoded",
			args: args{
				ctx: context.Background(),
				env: map[string]string{
					"GITHUB_INTEGRATION_ID": "10",
					"GITHUB_WEBHOOK_SECRET": "webhook",
					"GITHUB_PRIVATE_KEY":    "foobarbaz",
					"GITHUB_V3_ENDPOINT":    "http://example.com/api",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "missing required config value in environment",
			args: args{
				ctx: context.Background(),
				env: map[string]string{
					"GITHUB_INTEGRATION_ID": "10",
					"GITHUB_WEBHOOK_SECRET": "webhook",
					"GITHUB_V3_ENDPOINT":    "http://example.com/api",
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		// set env vars.
		for k, v := range tt.args.env {
			os.Setenv(k, v)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("NewConfig(%v)", tt.args.ctx)) {
				return
			} else if tt.want == nil {
				return
			}
			assert.Equalf(t, tt.want, got, "NewConfig(%v)", tt.args.ctx)
		})
		// unset env vars.
		for k, _ := range tt.args.env {
			os.Unsetenv(k)
		}
	}
}
