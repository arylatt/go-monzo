package monzo

import (
	"net/url"
)

type AccountsService service

type AccountType string

const (
	AccountTypeUKRetail      AccountType = "uk_retail"
	AccountTypeUKRetailJoint AccountType = "uk_retail_joint"
)

type AccountOwner struct {
	UserID             string `json:"user_id"`
	PreferredName      string `json:"preferred_name"`
	PreferredFirstName string `json:"preferred_first_name"`
}

type PaymentDetailsLocaleUK struct {
	AccountNumber string `json:"account_number"`
	SortCode      string `json:"sort_code"`
}

type PaymentDetails struct {
	LocaleUK PaymentDetailsLocaleUK `json:"locale_uk"`
}

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

type AccountsList struct {
	Accounts []Account `json:"accounts"`
}

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
