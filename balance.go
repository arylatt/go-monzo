package monzo

import (
	"fmt"
	"net/url"
)

type BalanceService service

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

func (s *BalanceService) Get(accountID string) (bal *Balance, err error) {
	bal = &Balance{}
	u := fmt.Sprintf("/balance?%s", url.Values{"account_id": []string{accountID}}.Encode())

	resp, err := s.client.Get(u, nil)
	err = ParseResponse(resp, err, bal)

	return
}
