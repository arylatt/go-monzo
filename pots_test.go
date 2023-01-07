package monzo

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		potID  string
		values map[string]interface{}
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
			inputs:   inputs{"", map[string]interface{}{"source_account_id": "", "amount": float64(0), "dedupe_id": ""}},
			expected: expected{nil, ErrPotInvalidID},
		},
		{
			inputs:   inputs{"1234", map[string]interface{}{"source_account_id": "", "amount": float64(0), "dedupe_id": ""}},
			expected: expected{nil, ErrPotInvalidSourceAccountID},
		},
		{
			inputs:   inputs{"1234", map[string]interface{}{"source_account_id": "5678", "amount": float64(0), "dedupe_id": ""}},
			expected: expected{nil, ErrPotInvalidDepositAmount},
		},
		{
			inputs:   inputs{"1234", map[string]interface{}{"source_account_id": "5678", "amount": float64(23), "dedupe_id": ""}},
			expected: expected{nil, ErrPotInvalidDedupeID},
		},
		{
			inputs:   inputs{"1234", map[string]interface{}{"source_account_id": "5678", "amount": float64(23), "dedupe_id": "a"}},
			expected: expected{&Pot{Balance: 23}, nil},
		},
	}

	for _, test := range tests {
		c := MockRequest(test.expected.pot, func(args mock.Arguments) {
			req := args.Get(0).(*http.Request)

			assert.Equal(t, fmt.Sprintf("/pots/%s/deposit", test.inputs.potID), req.URL.Path)

			params := &map[string]interface{}{}

			assert.NoError(t, json.NewDecoder(req.Body).Decode(params))
			assert.Equal(t, test.inputs.values, *params)
		})

		pot, err := c.Pots.Deposit(test.inputs.potID, test.inputs.values["source_account_id"].(string), int(test.inputs.values["amount"].(float64)), test.inputs.values["dedupe_id"].(string))

		assert.Equal(t, test.expected.pot, pot)
		assert.Equal(t, test.expected.err, err)
	}
}

func TestPotWithdraw(t *testing.T) {
	type inputs struct {
		potID  string
		values map[string]interface{}
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
			inputs:   inputs{"", map[string]interface{}{"destination_account_id": "", "amount": float64(0), "dedupe_id": ""}},
			expected: expected{nil, ErrPotInvalidID},
		},
		{
			inputs:   inputs{"1234", map[string]interface{}{"destination_account_id": "", "amount": float64(0), "dedupe_id": ""}},
			expected: expected{nil, ErrPotInvalidSourceAccountID},
		},
		{
			inputs:   inputs{"1234", map[string]interface{}{"destination_account_id": "5678", "amount": float64(0), "dedupe_id": ""}},
			expected: expected{nil, ErrPotInvalidWithdrawAmount},
		},
		{
			inputs:   inputs{"1234", map[string]interface{}{"destination_account_id": "5678", "amount": float64(23), "dedupe_id": ""}},
			expected: expected{nil, ErrPotInvalidDedupeID},
		},
		{
			inputs:   inputs{"1234", map[string]interface{}{"destination_account_id": "5678", "amount": float64(23), "dedupe_id": "a"}},
			expected: expected{&Pot{Balance: 0}, nil},
		},
	}

	for _, test := range tests {
		c := MockRequest(test.expected.pot, func(args mock.Arguments) {
			req := args.Get(0).(*http.Request)

			assert.Equal(t, fmt.Sprintf("/pots/%s/withdraw", test.inputs.potID), req.URL.Path)

			params := &map[string]interface{}{}

			assert.NoError(t, json.NewDecoder(req.Body).Decode(params))
			assert.Equal(t, test.inputs.values, *params)
		})

		pot, err := c.Pots.Withdraw(test.inputs.potID, test.inputs.values["destination_account_id"].(string), int(test.inputs.values["amount"].(float64)), test.inputs.values["dedupe_id"].(string))

		assert.Equal(t, test.expected.pot, pot)
		assert.Equal(t, test.expected.err, err)
	}
}
