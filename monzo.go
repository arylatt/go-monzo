package monzo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	// DefaultUserAgent is the user-agent string that will be sent to the server, unless it is overridden on the Client.
	DefaultUserAgent = "go-monzo/1.0 (https://github.com/arylatt/go-monzo)"

	// BaseURL is the default Monzo API URL.
	BaseURL = "https://api.monzo.com"
)

// Client is the Monzo API client.
//
// Client contains modifiable fields: BaseURL for changing where requests are sent, and UserAgent for changing the user-agent string sent to the server.
//
// The various API endpoints are accessed through the different Service fields (e.g. Accounts, Balance, Pots, etc...),
// based on the Monzo API Reference - https://docs.monzo.com/.
type Client struct {
	client *http.Client

	BaseURL   *url.URL
	UserAgent string

	common service

	Accounts     *AccountsService
	Balance      *BalanceService
	Pots         *PotsService
	Transactions *TransactionsService
	Feed         *FeedService
	Attachments  *AttachmentsService
	Receipts     *ReceiptsService
	Webhooks     *WebhooksService
}

// Internal struct to provide the different API services with the common client.
type service struct {
	client *Client
}

// New creates a new Monzo API client based on the parent HTTP client.
//
// The parent HTTP client should contain a transport capable of authorizing requests, either via static access token, or OAuth2.
func New(client *http.Client) (c *Client) {
	baseURL, _ := url.Parse(BaseURL)

	c = &Client{
		client: client,

		BaseURL:   baseURL,
		UserAgent: DefaultUserAgent,
	}

	c.common.client = c
	c.Accounts = (*AccountsService)(&c.common)
	c.Balance = (*BalanceService)(&c.common)
	c.Pots = (*PotsService)(&c.common)
	c.Transactions = (*TransactionsService)(&c.common)
	c.Feed = (*FeedService)(&c.common)
	c.Attachments = (*AttachmentsService)(&c.common)
	c.Receipts = (*ReceiptsService)(&c.common)
	c.Webhooks = (*WebhooksService)(&c.common)

	return
}

// Do sends a request to the server and attempts to parse the response data for a Monzo API error.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	parsedResp := CheckResponse(resp)

	if parsedResp != nil {
		return resp, parsedResp
	}

	return resp, err
}

// Internal helper to encode request body data and provide the appropriate content type string.
func encodeBody(body any) (io.Reader, string, error) {
	if body == nil {
		return nil, "", nil
	}

	switch body := body.(type) {
	case url.Values:
		return strings.NewReader(body.Encode()), "application/x-www-form-urlencoded", nil
	default:
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)

		enc.SetEscapeHTML(false)

		return buf, "application/json", enc.Encode(body)
	}
}

// NewRequest creates a new HTTP request to be sent to the Monzo API with a background context.
func (c *Client) NewRequest(method, url string, body any) (*http.Request, error) {
	return c.NewRequestWithContext(context.Background(), method, url, body)
}

// NewRequestWithContext creates a new HTTP request to be sent to the Monzo API with the provided context.
//
// The request body is encoded and attached to the request, as well as the appropriate user-agent, content type, and accept request headers.
func (c *Client) NewRequestWithContext(ctx context.Context, method, urlStr string, body any) (req *http.Request, err error) {
	bodyBuf, contentType, err := encodeBody(body)
	if err != nil {
		return
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return
	}

	req, err = http.NewRequestWithContext(ctx, method, u.String(), bodyBuf)
	if err != nil {
		return
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	req.Header.Set("Accept", "application/json")

	return
}

// Internal helper to call NewRequest and Do and return the results.
func (c *Client) doQuick(method, url string, body any) (resp *http.Response, err error) {
	req, err := c.NewRequest(method, url, body)
	if err != nil {
		return
	}

	return c.Do(req)
}

// Post sends a HTTP POST request.
func (c *Client) Post(url string, body any) (resp *http.Response, err error) {
	return c.doQuick(http.MethodPost, url, body)
}

// Put sends a HTTP PUT request.
func (c *Client) Put(url string, body any) (resp *http.Response, err error) {
	return c.doQuick(http.MethodPut, url, body)
}

// Patch sends a HTTP PATCH request.
func (c *Client) Patch(url string, body any) (resp *http.Response, err error) {
	return c.doQuick(http.MethodPatch, url, body)
}

// Get sends a HTTP GET request.
func (c *Client) Get(url string, body any) (resp *http.Response, err error) {
	return c.doQuick(http.MethodGet, url, body)
}

// Delete sends a HTTP DELETE request.
func (c *Client) Delete(url string) (resp *http.Response, err error) {
	return c.doQuick(http.MethodDelete, url, nil)
}

// LogOut revokes the access and refresh token. A new OAuth2Client will need to be created.
func (c *Client) LogOut() (err error) {
	_, err = c.Post("/oauth2/logout", nil)
	return
}

// ParseResponse attempts to decode the HTTP response body into the provided structure.
func ParseResponse(resp *http.Response, errIn error, v any) (err error) {
	err = errIn

	if err != nil {
		return
	}

	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF {
			err = nil
		}
	}

	return
}

// Pot represents data about the currently authenticated user provided by the Monzo API.
type Whoami struct {
	Authenticated bool   `json:"authenticated"`
	ClientID      string `json:"client_id"`
	UserID        string `json:"user_id"`
}

// Returns information about the current access token.
func (c *Client) Whoami() (who *Whoami, err error) {
	who = &Whoami{}
	resp, err := c.Get("/ping/whoami", nil)
	err = ParseResponse(resp, err, who)
	return
}

// Returns the OAuth2 token being currently used.
func (c *Client) Token() (*oauth2.Token, error) {
	switch t := c.client.Transport.(type) {
	case *oauth2.Transport:
		return t.Source.Token()
	}

	return nil, errors.New("could not access token from transport")
}

// RefreshToken updates the expiry time on the OAuth2 token to be in the past, and then calls Whoami to
// force the OAuth2 transport to refresh the token.
func (c *Client) RefreshToken() (err error) {
	err = c.RefreshTokenOnNextRequest()
	if err != nil {
		return
	}

	_, err = c.Whoami()

	return
}

// RefreshTokenOnNextRequest updates the expiry time on the OAuth2 token to be in the past, so that the next API call will
// force the OAuth2 transport to refresh the token.
func (c *Client) RefreshTokenOnNextRequest() (err error) {
	token, err := c.Token()
	if err != nil {
		return
	}

	token.Expiry = time.Unix(1, 0)

	return
}
