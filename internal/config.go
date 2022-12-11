package internal

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/sethvargo/go-envconfig"
	"io/ioutil"
)

type GithubAuthConfig struct {
	IntegrationID int64  `yaml:"integration_id" json:"integrationId"`
	WebhookSecret string `yaml:"webhook_secret" json:"webhookSecret"`
	PrivateKey    string `yaml:"private_key" json:"privateKey"`
}

type Config struct {
	IntegrationID    int64  `env:"GITHUB_INTEGRATION_ID,required"`
	WebhookSecret    string `env:"GITHUB_WEBHOOK_SECRET,required"`
	PrivateKeyBytes  []byte `env:"GITHUB_PRIVATE_KEY,required"`
	GithubV3Endpoint string `env:"GITHUB_V3_ENDPOINT,required"`
	PrivateKey       string
}

func (c *Config) ToGithubAppConfig() *githubapp.Config {
	return &githubapp.Config{
		V3APIURL: c.GithubV3Endpoint,
		App: GithubAuthConfig{
			IntegrationID: c.IntegrationID,
			WebhookSecret: c.WebhookSecret,
			PrivateKey:    c.PrivateKey,
		},
	}
}

func NewConfig(ctx context.Context) (*Config, error) {
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		return &config, err
	}
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(config.PrivateKeyBytes))
	b, err := ioutil.ReadAll(dec)
	if err != nil {
		return &config, err
	}
	config.PrivateKey = string(b)
	return &config, err
}
