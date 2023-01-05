package monzo

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWebhooksRegister(t *testing.T) {
	type inputs struct {
		accountID string
		url       string
	}

	type expected struct {
		webhook *WebhookSingle
		err     error
	}

	tests := []struct {
		inputs   inputs
		expected expected
	}{
		{
			inputs:   inputs{"", ""},
			expected: expected{nil, ErrWebhookInvalidAccountID},
		},
		{
			inputs:   inputs{"1", ""},
			expected: expected{nil, ErrWebhookInvalidURL},
		},
		{
			inputs: inputs{"1", "1"},
			expected: expected{&WebhookSingle{
				Webhook: Webhook{
					AccountID: "1",
					URL:       "1",
					ID:        "abc",
				},
			}, nil},
		},
	}

	for _, test := range tests {
		c := MockRequest(test.expected.webhook, func(args mock.Arguments) {
			req := args.Get(0).(*http.Request)

			assert.Equal(t, "/webhooks", req.URL.Path)

			req.ParseForm()
			assert.Equal(t, test.inputs.accountID, req.Form.Get("account_id"))
			assert.Equal(t, test.inputs.url, req.Form.Get("url"))
		})

		test.expected.webhook.setClient(c)

		webhook, err := c.Webhooks.Register(test.inputs.accountID, test.inputs.url)

		assert.Equal(t, test.expected.webhook, webhook)
		assert.Equal(t, test.expected.err, err)
	}
}

func TestWebhooksList(t *testing.T) {
	expected := &WebhookList{
		Webhooks: []Webhook{
			{
				AccountID: "1",
				ID:        "1_1",
				URL:       "1_1_1",
			},
			{
				AccountID: "1",
				ID:        "1_2",
				URL:       "1_2_1",
			},
		},
	}

	c := MockRequest(expected, func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)

		assert.Equal(t, expected.Webhooks[0].AccountID, req.URL.Query().Get("account_id"))
	})

	expected.setClient(c)

	actual, err := c.Webhooks.List(expected.Webhooks[0].AccountID)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestWebhookDelete(t *testing.T) {
	w := Webhook{
		ID: "a",
	}

	c := MockRequest(nil, func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)

		assert.Equal(t, fmt.Sprintf("/webhooks/%s", w.ID), req.URL.Path)
	})

	w.client = c

	err := w.Delete()

	assert.NoError(t, err)
}
