package monzo

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	baseClient := http.DefaultClient

	client := New(baseClient)

	assert.Equal(t, baseClient, client.client)
	assert.Equal(t, DefaultUserAgent, client.UserAgent)
	assert.Equal(t, client.common.client, client)
	assert.Equal(t, (*AccountsService)(&client.common), client.Accounts)
	assert.Equal(t, (*BalanceService)(&client.common), client.Balance)
	assert.Equal(t, (*PotsService)(&client.common), client.Pots)
	assert.Equal(t, (*TransactionsService)(&client.common), client.Transactions)
	assert.Equal(t, (*FeedService)(&client.common), client.Feed)
	assert.Equal(t, (*AttachmentsService)(&client.common), client.Attachments)
	assert.Equal(t, (*ReceiptsService)(&client.common), client.Receipts)
	assert.Equal(t, (*WebhooksService)(&client.common), client.Webhooks)
}

type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}

func MockRequest(expected any, run func(args mock.Arguments)) *Client {
	rt := &MockRoundTripper{}

	bytes, _ := json.Marshal(expected)
	reader := strings.NewReader(string(bytes))

	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       io.NopCloser(reader),
	}

	rt.On("RoundTrip", mock.AnythingOfType("*http.Request")).Run(run).Return(resp, nil)

	c := New(http.DefaultClient)
	c.client.Transport = rt

	return c
}
