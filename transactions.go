package monzo

import (
	"errors"
	"fmt"
	"net/url"
)

type TransactionsService service

var (
	ErrTransactionClientNil = errors.New("transaction client is not configured, use transactions client instead")
)

type MerchantAddress struct {
	ShortFormatted string  `json:"short_formatted"`
	City           string  `json:"city"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	ZoomLevel      int64   `json:"zoom_level"`
	Approximate    bool    `json:"approximate"`
	Formatted      string  `json:"formatted"`
	Address        string  `json:"address"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Postcode       string  `json:"postcode"`
}

type Merchant struct {
	ID              string            `json:"id"`
	GroupID         string            `json:"group_id"`
	Name            string            `json:"name"`
	Logo            string            `json:"logo"`
	Emoji           string            `json:"emoji"`
	Category        string            `json:"category"`
	Created         string            `json:"created"`
	Online          bool              `json:"online"`
	ATM             bool              `json:"atm"`
	Address         MerchantAddress   `json:"address"`
	DisableFeedback bool              `json:"disable_feedback"`
	SuggestedTags   string            `json:"suggested_tags"`
	Metadata        map[string]string `json:"metadata"`
}

type Transaction struct {
	AccountID                            string            `json:"account_id"`
	Amount                               int64             `json:"amount"`
	AmountIsPending                      bool              `json:"amount_is_pending"`
	ATMFeesDetailed                      interface{}       `json:"atm_fees_detailed"`
	Attachments                          interface{}       `json:"attachments"`
	CanAddToTab                          bool              `json:"can_add_to_tab"`
	CanBeExcludedFromBreakdown           bool              `json:"can_be_excluded_from_breakdown"`
	CanBeMadeSubscription                bool              `json:"can_be_made_subscription"`
	CanMatchTransactionsInCategorization bool              `json:"can_match_transactions_in_categorization"`
	CanSplitTheBill                      bool              `json:"can_split_the_bill"`
	Categories                           map[string]int64  `json:"categories"`
	Category                             string            `json:"category"`
	Counterparty                         interface{}       `json:"counterparty"`
	Created                              string            `json:"created"`
	Currency                             string            `json:"currency"`
	DedupeID                             string            `json:"dedupe_id"`
	Description                          string            `json:"description"`
	Fees                                 interface{}       `json:"fees"`
	ID                                   string            `json:"id"`
	IncludeInSpending                    bool              `json:"include_in_spending"`
	International                        interface{}       `json:"international"`
	IsLoad                               bool              `json:"is_load"`
	Labels                               interface{}       `json:"labels"`
	LocalAmount                          int64             `json:"local_amount"`
	LocalCurrency                        string            `json:"local_currency"`
	Merchant                             Merchant          `json:"merchant"`
	Metadata                             map[string]string `json:"metadata"`
	Notes                                string            `json:"notes"`
	Originator                           bool              `json:"originator"`
	ParentAccountID                      string            `json:"parent_account_id"`
	Scheme                               string            `json:"scheme"`
	Settled                              string            `json:"settled"`
	Updated                              string            `json:"updated"`
	UserID                               string            `json:"user_id"`

	client *Client
}

type TransactionList struct {
	Transactions []Transaction `json:"transactions"`
}

func (t *TransactionList) setClient(c *Client) {
	if t == nil {
		return
	}

	updatedTransactions := []Transaction{}

	for _, tx := range t.Transactions {
		tx.client = c
		updatedTransactions = append(updatedTransactions, tx)
	}

	t.Transactions = updatedTransactions
}

type TransactionSingle struct {
	Transaction Transaction `json:"transaction"`
}

func (t *TransactionSingle) setClient(c *Client) {
	if t == nil {
		return
	}

	t.Transaction.client = c
}

type transactionStringMerchant struct {
	Transaction

	Merchant string `json:"merchant"`
}

func (t transactionStringMerchant) Expand() *Transaction {
	t.Transaction.Merchant.ID = t.Merchant

	return &t.Transaction
}

type transactionStringMerchantList struct {
	Transactions []transactionStringMerchant `json:"transactions"`
}

func (t transactionStringMerchantList) Expand() *TransactionList {
	tl := &TransactionList{}

	for _, tx := range t.Transactions {
		tl.Transactions = append(tl.Transactions, *tx.Expand())
	}

	return tl
}

type transactionStringMerchantSingle struct {
	Transaction transactionStringMerchant `json:"transaction"`
}

func (t transactionStringMerchantSingle) Expand() *TransactionSingle {
	t.Transaction.Transaction.Merchant.ID = t.Transaction.Merchant

	return &TransactionSingle{Transaction: t.Transaction.Transaction}
}

func (s *TransactionsService) List(accountID string, expandMerchant bool, paging *Pagination) (list *TransactionList, err error) {
	var out any

	params := url.Values{
		"account_id": []string{accountID},
	}

	out = &transactionStringMerchantList{}

	if expandMerchant {
		params.Add("expand[]", "merchant")
		out = &TransactionList{}
	}

	if paging != nil {
		params = paging.Values(params)
	}

	u := fmt.Sprintf("/transactions?%s", params.Encode())

	resp, err := s.client.Get(u, nil)
	err = ParseResponse(resp, err, out)

	switch t := out.(type) {
	case *transactionStringMerchantList:
		list = t.Expand()
	case *TransactionList:
		list = t
	}

	list.setClient(s.client)

	return
}

func (s *TransactionsService) Get(transactionID string, expandMerchant bool) (tx *TransactionSingle, err error) {
	var out any

	params := url.Values{}
	out = &transactionStringMerchantSingle{}

	if expandMerchant {
		params.Add("expand[]", "merchant")
		out = &TransactionSingle{}
	}

	u := fmt.Sprintf("/transactions/%s?%s", transactionID, params.Encode())

	resp, err := s.client.Get(u, nil)
	err = ParseResponse(resp, err, out)

	switch t := out.(type) {
	case *transactionStringMerchantSingle:
		tx = t.Expand()
	case *TransactionSingle:
		tx = t
	}

	tx.setClient(s.client)

	return
}

func (s *TransactionsService) Annotate(transactionID string, metadata map[string]string) (tx *TransactionSingle, err error) {
	u := fmt.Sprintf("/transactions/%s", transactionID)
	body := url.Values{}
	out := &transactionStringMerchantSingle{}

	for k, v := range metadata {
		body.Add(fmt.Sprintf("metadata[%s]", k), v)
	}

	resp, err := s.client.Patch(u, body)
	err = ParseResponse(resp, err, out)

	tx = out.Expand()

	tx.setClient(s.client)

	return
}

func (t *Transaction) Annotate(metadata map[string]string) (*TransactionSingle, error) {
	if t.client == nil {
		return nil, ErrTransactionClientNil
	}

	return t.client.Transactions.Annotate(t.ID, metadata)
}
