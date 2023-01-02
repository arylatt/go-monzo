package monzo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestOAuth2Config(t *testing.T) {
	id, secret, redirect := "id", "secret", "redirect"

	expected := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     OAuth2Endpoint,
		RedirectURL:  redirect,
	}

	assert.Equal(t, expected, OAuth2Config(id, secret, redirect))
}
