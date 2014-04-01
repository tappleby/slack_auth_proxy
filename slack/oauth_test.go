package slack_test

import (
	"github.com/tappleby/slack-auth-proxy/slack"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOAuthClient(t *testing.T) {
	client := slack.NewOAuthClient("foo", "bar", "baz")

	assert.Equal(t, client.ClientId, "foo")
	assert.Equal(t, client.ClientSecret, "bar")
	assert.Equal(t, client.RedirectUri, "baz")
}

