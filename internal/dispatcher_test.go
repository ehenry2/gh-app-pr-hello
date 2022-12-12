package internal

import (
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterGithubWebhookDispatcher(t *testing.T) {
	config := &githubapp.Config{
		V3APIURL: "",
		App: struct {
			IntegrationID int64  `yaml:"integration_id" json:"integrationId"`
			WebhookSecret string `yaml:"webhook_secret" json:"webhookSecret"`
			PrivateKey    string `yaml:"private_key" json:"privateKey"`
		}{
			IntegrationID: 10,
			WebhookSecret: "secret",
			PrivateKey:    "pem",
		},
	}
	err := RegisterGithubWebhookDispatcher(config)
	assert.NoError(t, err)
}
