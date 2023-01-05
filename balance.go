package monzo

import (
	"fmt"
	"net/url"
)

// Retrieve information about an account's balance.
type BalanceService service

// Balance represents the balance data provided by the Monzo API.
type Balance struct {
	Balance                         int64         `json:"balance"`
	TotalBalance                    int64         `json:"total_balance"`
	BalanceIncludingFlexibleSavings int64         `json:"balance_including_flexible_savings"`
	Currency                        string        `json:"currency"`
	SpendToday                      int64         `json:"spend_today"`
	LocalCurrency                   string        `json:"local_currency"`
	LocalExchangeRate               int64         `json:"local_exchange_rate"`
	LocalSpend                      []interface{} `json:"local_spend"`
}

// Returns balance information for a specific account.
func (s *BalanceService) Get(accountID string) (bal *Balance, err error) {
	bal = &Balance{}
	u := fmt.Sprintf("/balance?%s", url.Values{"account_id": []string{accountID}}.Encode())

	resp, err := s.client.Get(u, nil)
	err = ParseResponse(resp, err, bal)

	return
}
