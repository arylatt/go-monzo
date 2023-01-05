package monzo

import (
	"net/url"
)

// Accounts represent a store of funds, and have a list of transactions.
type AccountsService service

type AccountType string

const (
	// AccountTypeUKRetail is the single owner UK Current Account type.
	AccountTypeUKRetail AccountType = "uk_retail"

	// AccountTypeUKRetailJoint is the joint owner UK Current Account type.
	AccountTypeUKRetailJoint AccountType = "uk_retail_joint"
)

// AccountOwner represents the inner account owner data provided by the Monzo API.
type AccountOwner struct {
	UserID             string `json:"user_id"`
	PreferredName      string `json:"preferred_name"`
	PreferredFirstName string `json:"preferred_first_name"`
}

// PaymentDetailsLocaleUK represents the inner payment details account data for UK accounts provided by the Monzo API.
type PaymentDetailsLocaleUK struct {
	AccountNumber string `json:"account_number"`
	SortCode      string `json:"sort_code"`
}

// PaymentDetails represents the inner payment details account data provided by the Monzo API.
type PaymentDetails struct {
	LocaleUK PaymentDetailsLocaleUK `json:"locale_uk"`
}

// Account represents the account data provided by the Monzo API.
type Account struct {
	ID             string         `json:"id"`
	Description    string         `json:"description"`
	Created        string         `json:"created"`
	Closed         bool           `json:"closed"`
	Type           AccountType    `json:"type"`
	CountryCode    string         `json:"country_code"`
	Owners         []AccountOwner `json:"owners"`
	AccountNumber  string         `json:"account_number"`
	SortCode       string         `json:"sort_code"`
	PaymentDetails PaymentDetails `json:"payment_details"`
}

// AccountsList represents the response from the Monzo API for a list of accounts.
type AccountsList struct {
	Accounts []Account `json:"accounts"`
}

// Returns a list of accounts owned by the currently authorised user.
//
// To filter by either single or joint current account, add accountType argument. Valid accountTypes are AccountTypeUKRetail, AccountTypeUKRetailJoint.
func (s *AccountsService) List(accountType ...AccountType) (list *AccountsList, err error) {
	list = &AccountsList{}
	u := "/accounts"

	if len(accountType) > 0 {
		u += "?" + url.Values{
			"account_type": []string{string(accountType[0])},
		}.Encode()
	}

	resp, err := s.client.Get(u, nil)
	err = ParseResponse(resp, err, list)

	return
}
