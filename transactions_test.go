package monzo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionsListNoExpand(t *testing.T) {
	expected := &TransactionList{
		[]Transaction{
			{
				ID: "test",
				Merchant: Merchant{
					ID: "test_merchant",
				},
			},
		},
	}

	payload := map[string]interface{}{
		"transactions": []map[string]interface{}{{
			"id":       expected.Transactions[0].ID,
			"merchant": expected.Transactions[0].Merchant.ID,
		}},
	}

	c := MockRequest(payload, nil)

	expected.setClient(c)

	list, err := c.Transactions.List("test", false, nil)

	assert.Equal(t, expected, list)
	assert.NoError(t, err)
}

func TestTransactionsList(t *testing.T) {
	expected := &TransactionList{
		[]Transaction{
			{
				Amount:      -510,
				Created:     "2015-08-22T12:20:18Z",
				Currency:    "GBP",
				Description: "THE DE BEAUVOIR DELI C LONDON        GBR",
				ID:          "tx_00008zIcpb1TB4yeIFXMzx",
				Merchant: Merchant{
					ID: "merch_00008zIcpbAKe8shBxXUtl",
				},
				Metadata: map[string]string{},
				Notes:    "Salmon sandwich 🍞",
				IsLoad:   false,
				Settled:  "2015-08-23T12:20:18Z",
				Category: "eating_out",
			},
			{
				Amount:      -679,
				Created:     "2015-08-23T16:15:03Z",
				Currency:    "GBP",
				Description: "VUE BSL LTD            ISLINGTON     GBR",
				ID:          "tx_00008zL2INM3xZ41THuRF3",
				Merchant: Merchant{
					ID: "merch_00008z6uFVhVBcaZzSQwCX",
				},
				Metadata: map[string]string{},
				Notes:    "",
				IsLoad:   false,
				Settled:  "2015-08-24T16:15:03Z",
				Category: "eating_out",
			},
		},
	}

	page := &Pagination{Limit: 2}

	c := MockRequest(expected, func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)

		assert.Equal(t, "merchant", req.URL.Query().Get("expand[]"))
		assert.Equal(t, strconv.Itoa(page.Limit), req.URL.Query().Get("limit"))
	})

	expected.setClient(c)

	list, err := c.Transactions.List("test", true, page)

	assert.Equal(t, expected, list)
	assert.NoError(t, err)
}

func TestTransactionsGetNoExpand(t *testing.T) {
	expected := &TransactionSingle{
		Transaction{
			ID: "test",
			Merchant: Merchant{
				ID: "test_merchant",
			},
		},
	}

	payload := map[string]interface{}{
		"transaction": map[string]interface{}{
			"id":       expected.Transaction.ID,
			"merchant": expected.Transaction.Merchant.ID,
		},
	}

	c := MockRequest(payload, nil)

	expected.setClient(c)

	list, err := c.Transactions.Get("test", false)

	assert.Equal(t, expected, list)
	assert.NoError(t, err)
}

func TestTransactionsGet(t *testing.T) {
	expected := &TransactionSingle{
		Transaction{
			Amount:      -510,
			Created:     "2015-08-22T12:20:18Z",
			Currency:    "GBP",
			Description: "THE DE BEAUVOIR DELI C LONDON        GBR",
			ID:          "tx_00008zIcpb1TB4yeIFXMzx",
			Merchant: Merchant{
				Address: MerchantAddress{
					Address:   "98 Southgate Road",
					City:      "London",
					Country:   "GB",
					Latitude:  51.54151,
					Longitude: -0.08482400000002599,
					Postcode:  "N1 3JD",
					Region:    "Greater London",
				},
				Created:  "2015-08-22T12:20:18Z",
				GroupID:  "grp_00008zIcpbBOaAr7TTP3sv",
				ID:       "merch_00008zIcpbAKe8shBxXUtl",
				Logo:     "https://pbs.twimg.com/profile_images/527043602623389696/68_SgUWJ.jpeg",
				Emoji:    "🍞",
				Name:     "The De Beauvoir Deli Co.",
				Category: "eating_out",
			},
			Metadata: map[string]string{},
			Notes:    "Salmon sandwich 🍞",
			IsLoad:   false,
			Settled:  "2015-08-23T12:20:18Z",
		},
	}

	c := MockRequest(expected, func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)

		assert.Equal(t, "merchant", req.URL.Query().Get("expand[]"))
		assert.Equal(t, fmt.Sprintf("/transactions/%s", expected.Transaction.ID), req.URL.Path)
	})

	expected.setClient(c)

	tx, err := c.Transactions.Get(expected.Transaction.ID, true)

	assert.Equal(t, expected, tx)
	assert.NoError(t, err)
}

func TestTransactionsAnnotate(t *testing.T) {
	expected := &TransactionSingle{
		Transaction{
			ID: "tx_00008zL2INM3xZ41THuRF3",
			Merchant: Merchant{
				ID: "merch_00008z6uFVhVBcaZzSQwCX",
			},
			Metadata: map[string]string{
				"foo": "bar",
			},
		},
	}

	payload := map[string]interface{}{
		"transaction": map[string]interface{}{
			"id":       expected.Transaction.ID,
			"merchant": expected.Transaction.Merchant.ID,
			"metadata": expected.Transaction.Metadata,
		},
	}

	c := MockRequest(payload, func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)

		params := map[string]interface{}{}

		assert.NoError(t, json.NewDecoder(req.Body).Decode(&params))
		assert.Equal(t, expected.Transaction.Metadata["foo"], params["metadata"].(map[string]interface{})["foo"].(string))
		assert.Equal(t, fmt.Sprintf("/transactions/%s", expected.Transaction.ID), req.URL.Path)
	})

	expected.setClient(c)

	tx, err := c.Transactions.Annotate(expected.Transaction.ID, expected.Transaction.Metadata)

	assert.Equal(t, expected, tx)
	assert.NoError(t, err)
}
