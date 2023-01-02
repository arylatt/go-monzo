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
	Version          = "dev"
	DefaultUserAgent = "go-monzo/" + Version + " (https://github.com/arylatt/go-monzo)"
	BaseURL          = "https://api.monzo.com"
)

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

type service struct {
	client *Client
}

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

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	parsedResp := CheckResponse(resp)

	if parsedResp != nil {
		return resp, parsedResp
	}

	return resp, err
}

func encodeBody(body interface{}) (io.Reader, string, error) {
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

func (c *Client) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	return c.NewRequestWithContext(context.Background(), method, url, body)
}

func (c *Client) NewRequestWithContext(ctx context.Context, method, urlStr string, body interface{}) (req *http.Request, err error) {
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

func (c *Client) doQuick(method, url string, body any) (resp *http.Response, err error) {
	req, err := c.NewRequest(method, url, body)
	if err != nil {
		return
	}

	return c.Do(req)
}

// Post sends a HTTP POST request
func (c *Client) Post(url string, body interface{}) (resp *http.Response, err error) {
	return c.doQuick(http.MethodPost, url, body)
}

// Put sends a HTTP PUT request
func (c *Client) Put(url string, body interface{}) (resp *http.Response, err error) {
	return c.doQuick(http.MethodPut, url, body)
}

// Patch sends a HTTP PATCH request
func (c *Client) Patch(url string, body interface{}) (resp *http.Response, err error) {
	return c.doQuick(http.MethodPatch, url, body)
}

// Get sends a HTTP Get request
func (c *Client) Get(url string, body interface{}) (resp *http.Response, err error) {
	return c.doQuick(http.MethodGet, url, body)
}

// LogOut revokes the access and refresh token. A new OAuth2Client will need to be created
func (c *Client) LogOut() (err error) {
	_, err = c.Post("/oauth2/logout", nil)
	return
}

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

type Whoami struct {
	Authenticated bool   `json:"authenticated"`
	ClientID      string `json:"client_id"`
	UserID        string `json:"user_id"`
}

func (c *Client) Whoami() (who *Whoami, err error) {
	who = &Whoami{}
	resp, err := c.Get("/ping/whoami", nil)
	err = ParseResponse(resp, err, who)
	return
}

func (c *Client) Token() (*oauth2.Token, error) {
	switch t := c.client.Transport.(type) {
	case *oauth2.Transport:
		return t.Source.Token()
	}

	return nil, errors.New("could not access token from transport")
}

func (c *Client) RefreshToken() (err error) {
	err = c.RefreshTokenOnNextRequest()
	if err != nil {
		return
	}

	_, err = c.Whoami()

	return
}

func (c *Client) RefreshTokenOnNextRequest() (err error) {
	token, err := c.Token()
	if err != nil {
		return
	}

	token.Expiry = time.Unix(1, 0)

	return
}
