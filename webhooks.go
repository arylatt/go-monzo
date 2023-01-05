package monzo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Webhooks allow your application to receive real-time, push notification of events in an account.
type WebhooksService service

var (
	// ErrWebhookInvalidAccountID is returned if a null/empty Account ID is supplied.
	ErrWebhookInvalidAccountID = errors.New("account id cannot be empty")

	// ErrWebhookInvalidURL is returned if a null/empty URL is supplied.
	ErrWebhookInvalidURL = errors.New("url cannot be empty")

	// ErrWebhookClientNil is returned if Webhook object has a nil client configured (e.g. if the Webhook object was manually created).
	ErrWebhookClientNil = errors.New("webhook client is not configured, use webhooks client instead")
)

// WebhookPayload represents the data that Monzo will send to registered webhook URLs.
type WebhookPayload struct {
	Type string      `json:"type"`
	Data Transaction `json:"data"`
}

// Webhook represents a webhook object provided by the Monzo API.
type Webhook struct {
	AccountID string `json:"account_id"`
	ID        string `json:"id"`
	URL       string `json:"url"`

	client *Client
}

// WebhookSingle represents the response from the Monzo API for a single webhook.
type WebhookSingle struct {
	Webhook Webhook `json:"webhook"`
}

// setClient is an internal helper to ensure the single webhook has the client attached to it for later usage.
func (w *WebhookSingle) setClient(c *Client) {
	if w == nil {
		return
	}

	w.Webhook.client = c
}

// WebhookList represents the response from the Monzo API for a list of webhooks.
type WebhookList struct {
	Webhooks []Webhook `json:"webhooks"`
}

// setClient is an internal helper to ensure all webhooks have the client attached to them for later usage.
func (w *WebhookList) setClient(c *Client) {
	if w == nil {
		return
	}

	for i := range w.Webhooks {
		w.Webhooks[i].client = c
	}
}

// Register creates a webhook entry that Monzo will call.
//
// Each time an event occurs, Monzo will make a POST call to the URL provided. If the call fails, Monzo will retry up to a maximum of 5 attempts, with exponential backoff.
func (s *WebhooksService) Register(accountID, webhookURL string) (w *WebhookSingle, err error) {
	w = &WebhookSingle{}

	if strings.TrimSpace(accountID) == "" {
		return nil, ErrWebhookInvalidAccountID
	}

	if strings.TrimSpace(webhookURL) == "" {
		return nil, ErrWebhookInvalidURL
	}

	params := url.Values{
		"account_id": []string{accountID},
		"url":        []string{webhookURL},
	}

	resp, err := s.client.Post("/webhooks", params)
	err = ParseResponse(resp, err, w)

	w.setClient(s.client)

	return
}

// List the webhooks your application has registered on an account.
func (s *WebhooksService) List(accountID string) (w *WebhookList, err error) {
	w = &WebhookList{}
	u := fmt.Sprintf("/webhooks?%s", url.Values{"account_id": []string{accountID}}.Encode())

	resp, err := s.client.Get(u, nil)
	err = ParseResponse(resp, err, w)

	w.setClient(s.client)

	return
}

// Delete removes a webhook from the account.
//
// When you delete a webhook, Monzo will no longer send notifications to it.
func (s *WebhooksService) Delete(webhookID string) (err error) {
	u := fmt.Sprintf("/webhooks/%s", webhookID)

	_, err = s.client.Delete(u)

	return
}

// Delete removes a webhook from the account.
//
// When you delete a webhook, Monzo will no longer send notifications to it.
//
// Webhook.Delete is a convenience method. It is the same as calling Webhooks.Delete(webhook.ID).
func (w Webhook) Delete() error {
	if w.client == nil {
		return ErrWebhookClientNil
	}

	return w.client.Webhooks.Delete(w.ID)
}

// WebhookPayloadHandler returns a HTTP HandlerFunc that can be used to receive the payloads that Monzo sends.
//
// If valid, the webhook payload will already be parsed into the payload argument of the handler function.
//
// The request Body will be reset for further manual processing as desired.
func WebhookPayloadHandler(handler func(rw http.ResponseWriter, r *http.Request, payload *WebhookPayload)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &WebhookPayload{}

		data, _ := io.ReadAll(r.Body)
		if data != nil {
			json.Unmarshal(data, payload)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(data))

		handler(w, r, payload)
	}
}
