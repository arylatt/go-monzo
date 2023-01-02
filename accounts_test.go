package monzo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountsList(t *testing.T) {
	expected := &AccountsList{
		Accounts: []Account{{
			ID:          "acc_00009237aqC8c5umZmrRdh",
			Description: "Peter Pan's Account",
			Created:     "2015-11-13T12:17:42Z",
		}},
	}

	c := MockRequest(expected, nil)

	accs, err := c.Accounts.List()

	assert.NoError(t, err)
	assert.Equal(t, expected, accs)
}

func TestAccountsListFilter(t *testing.T) {
	expected := &AccountsList{
		Accounts: []Account{{
			ID:          "acc_00009237aqC8c5umZmrRdh",
			Description: "Peter Pan's Account",
			Created:     "2015-11-13T12:17:42Z",
		}},
	}

	c := MockRequest(expected, func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)

		assert.Equal(t, string(AccountTypeUKRetailJoint), req.URL.Query().Get("account_type"))
	})

	_, err := c.Accounts.List(AccountTypeUKRetailJoint)

	assert.NoError(t, err)
}
