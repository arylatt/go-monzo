package monzo

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPotsList(t *testing.T) {
	expected := &PotsList{
		Pots: []Pot{{
			ID:               "pot_0000778xxfgh4iu8z83nWb",
			Name:             "Savings",
			Style:            "beach_ball",
			Balance:          133700,
			Currency:         "GBP",
			Created:          "2017-11-09T12:30:53.695Z",
			Updated:          "2017-11-09T12:30:53.695Z",
			Deleted:          false,
			CurrentAccountID: "1234",
		}},
	}

	c := MockRequest(expected, func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)

		assert.Equal(t, expected.Pots[0].CurrentAccountID, req.URL.Query().Get("current_account_id"))
	})

	expected.setClient(c)

	pots, err := c.Pots.List(expected.Pots[0].CurrentAccountID)

	assert.NoError(t, err)
	assert.Equal(t, expected, pots)
}

func TestPotDeposit(t *testing.T) {
	type inputs struct {
		potID           string
		sourceAccountID string
		amount          int64
		dedupeID        string
	}

	type expected struct {
		pot *Pot
		err error
	}

	tests := []struct {
		inputs   inputs
		expected expected
	}{
		{
			inputs:   inputs{"", "", 0, ""},
			expected: expected{nil, ErrPotInvalidID},
		},
		{
			inputs:   inputs{"1234", "", 0, ""},
			expected: expected{nil, ErrPotInvalidSourceAccountID},
		},
		{
			inputs:   inputs{"1234", "5678", 0, ""},
			expected: expected{nil, ErrPotInvalidDepositAmount},
		},
		{
			inputs:   inputs{"1234", "5678", 23, ""},
			expected: expected{nil, ErrPotInvalidDedupeID},
		},
		{
			inputs:   inputs{"1234", "5678", 23, "a"},
			expected: expected{&Pot{Balance: 23}, nil},
		},
	}

	for _, test := range tests {
		c := MockRequest(test.expected.pot, func(args mock.Arguments) {
			req := args.Get(0).(*http.Request)

			assert.Equal(t, fmt.Sprintf("/pots/%s/deposit", test.inputs.potID), req.URL.Path)

			req.ParseForm()
			assert.Equal(t, test.inputs.sourceAccountID, req.Form.Get("source_account_id"))
			assert.Equal(t, strconv.FormatInt(test.inputs.amount, 10), req.Form.Get("amount"))
			assert.Equal(t, test.inputs.dedupeID, req.Form.Get("dedupe_id"))
		})

		pot, err := c.Pots.Deposit(test.inputs.potID, test.inputs.sourceAccountID, test.inputs.amount, test.inputs.dedupeID)

		assert.Equal(t, test.expected.pot, pot)
		assert.Equal(t, test.expected.err, err)
	}
}

func TestPotWithdraw(t *testing.T) {
	type inputs struct {
		potID                string
		destinationAccountID string
		amount               int64
		dedupeID             string
	}

	type expected struct {
		pot *Pot
		err error
	}

	tests := []struct {
		inputs   inputs
		expected expected
	}{
		{
			inputs:   inputs{"", "", 0, ""},
			expected: expected{nil, ErrPotInvalidID},
		},
		{
			inputs:   inputs{"1234", "", 0, ""},
			expected: expected{nil, ErrPotInvalidSourceAccountID},
		},
		{
			inputs:   inputs{"1234", "5678", 0, ""},
			expected: expected{nil, ErrPotInvalidWithdrawAmount},
		},
		{
			inputs:   inputs{"1234", "5678", 23, ""},
			expected: expected{nil, ErrPotInvalidDedupeID},
		},
		{
			inputs:   inputs{"1234", "5678", 23, "a"},
			expected: expected{&Pot{Balance: 0}, nil},
		},
	}

	for _, test := range tests {
		c := MockRequest(test.expected.pot, func(args mock.Arguments) {
			req := args.Get(0).(*http.Request)

			assert.Equal(t, fmt.Sprintf("/pots/%s/withdraw", test.inputs.potID), req.URL.Path)

			req.ParseForm()
			assert.Equal(t, test.inputs.destinationAccountID, req.Form.Get("destination_account_id"))
			assert.Equal(t, strconv.FormatInt(test.inputs.amount, 10), req.Form.Get("amount"))
			assert.Equal(t, test.inputs.dedupeID, req.Form.Get("dedupe_id"))
		})

		pot, err := c.Pots.Withdraw(test.inputs.potID, test.inputs.destinationAccountID, test.inputs.amount, test.inputs.dedupeID)

		assert.Equal(t, test.expected.pot, pot)
		assert.Equal(t, test.expected.err, err)
	}
}
