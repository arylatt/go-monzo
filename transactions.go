package monzo

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Transactions are movements of funds into or out of an account.
// Negative transactions represent debits (ie. spending money) and positive transactions represent credits (ie. receiving money).
type TransactionsService service

var (
	// ErrTransactionClientNil is returned if Transaction object has a nil client configured
	// (e.g. if the Transaction object was manually created).
	ErrTransactionClientNil = errors.New("transaction client is not configured, use transactions client instead")
)

// MerchantAddress represents the inner merchant address data provided by the Monzo API.
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

// MerchantAddress represents the inner merchant data provided by the Monzo API.
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

// Transaction represents a transaction provided by the Monzo API.
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

// CreatedTime converts the RFC3339 string time into a Time object.
func (t Transaction) CreatedTime() time.Time {
	tme, _ := time.Parse(time.RFC3339Nano, t.Created)
	return tme
}

// UpdatedTime converts the RFC3339 string time into a Time object.
func (t Transaction) UpdatedTime() time.Time {
	tme, _ := time.Parse(time.RFC3339Nano, t.Updated)
	return tme
}

// TransactionList represents the response from the Monzo API for a list of transactions.
type TransactionList struct {
	Transactions []Transaction `json:"transactions"`
}

// setClient is an internal helper to ensure all transactions have the client attached to them for later usage.
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

// TransactionSingle represents the response from the Monzo API for a single transaction.
type TransactionSingle struct {
	Transaction Transaction `json:"transaction"`
}

// setClient is an internal helper to ensure the single transaction has the client attached to it for later usage.
func (t *TransactionSingle) setClient(c *Client) {
	if t == nil {
		return
	}

	t.Transaction.client = c
}

// transactionStringMerchant represents a transaction provided by the Monzo API with non-expanded merchant data.
type transactionStringMerchant struct {
	Transaction

	Merchant string `json:"merchant"`
}

// Expand converts a transaction provided by the Monzo API with non-expanded merchant data to the regular Transaction type.
func (t transactionStringMerchant) Expand() *Transaction {
	t.Transaction.Merchant.ID = t.Merchant

	return &t.Transaction
}

// transactionStringMerchantList represents the response from the Monzo API for a list of transactions with non-expanded merchant data.
type transactionStringMerchantList struct {
	Transactions []transactionStringMerchant `json:"transactions"`
}

// Expand converts a list of transactions provided by the Monzo API with non-expanded merchant data to the regular Transaction type.
func (t transactionStringMerchantList) Expand() *TransactionList {
	tl := &TransactionList{}

	for _, tx := range t.Transactions {
		tl.Transactions = append(tl.Transactions, *tx.Expand())
	}

	return tl
}

// transactionStringMerchantSingle represents the response from the Monzo API for a single transaction with non-expanded merchant data.
type transactionStringMerchantSingle struct {
	Transaction transactionStringMerchant `json:"transaction"`
}

// Expand converts a single transaction provided by the Monzo API with non-expanded merchant data to the regular Transaction type.
func (t transactionStringMerchantSingle) Expand() *TransactionSingle {
	t.Transaction.Transaction.Merchant.ID = t.Transaction.Merchant

	return &TransactionSingle{Transaction: t.Transaction.Transaction}
}

// Returns a list of transactions on the user's account.
//
// IMPORTANT - Strong Customer Authentication:
// After a user has authenticated, your client can fetch all of their transactions, and after 5 minutes, it can only sync the last 90 days of transactions.
// If you need the user’s entire transaction history, you should consider fetching and storing it right after authentication.
//
// Transactions within the last 90 days can be accessed using the paging argument.
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

// Returns an individual transaction, fetched by its id.
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

// You may store your own key-value annotations against a transaction in its metadata.
//
// Note: the Monzo API does not seem to be respecting these requests.
func (s *TransactionsService) Annotate(transactionID string, metadata map[string]string) (tx *TransactionSingle, err error) {
	u := fmt.Sprintf("/transactions/%s", transactionID)
	out := &transactionStringMerchantSingle{}

	body := map[string]interface{}{
		"metadata": metadata,
	}

	resp, err := s.client.Patch(u, body)
	err = ParseResponse(resp, err, out)

	tx = out.Expand()

	tx.setClient(s.client)

	return
}

// You may store your own key-value annotations against a transaction in its metadata.
//
// Note: the Monzo API does not seem to be respecting these requests.
//
// Transaction.Annotate is a convenience method. It is the same as calling Transactions.Annotate(transaction.ID, metadata).
func (t *Transaction) Annotate(metadata map[string]string) (*TransactionSingle, error) {
	if t.client == nil {
		return nil, ErrTransactionClientNil
	}

	return t.client.Transactions.Annotate(t.ID, metadata)
}
