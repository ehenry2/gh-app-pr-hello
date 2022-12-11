package internal

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/sethvargo/go-envconfig"
	"io/ioutil"
)

type Config struct {
	IntegrationID    int64  `env:"GITHUB_INTEGRATION_ID,required"`
	WebhookSecret    string `env:"GITHUB_WEBHOOK_SECRET,required"`
	PrivateKey       []byte `env:"GITHUB_PRIVATE_KEY,required"`
	GithubV3Endpoint string `env:"GITHUB_V3_ENDPOINT,required"`
}

func NewConfig(ctx context.Context) (*Config, error) {
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		return &config, err
	}
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(config.PrivateKey))
	b, err := ioutil.ReadAll(dec)
	if err != nil {
		return &config, err
	}
	config.PrivateKey = b
	return &config, err
}
