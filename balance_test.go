package monzo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBalanceGet(t *testing.T) {
	expected := &Balance{
		Balance:      5000,
		TotalBalance: 6000,
		Currency:     "GBP",
		SpendToday:   0,
	}

	accountID := "1234"

	c := MockRequest(expected, func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)

		assert.Equal(t, accountID, req.URL.Query().Get("account_id"))
	})

	bal, err := c.Balance.Get(accountID)

	assert.NoError(t, err)
	assert.Equal(t, expected, bal)
}
