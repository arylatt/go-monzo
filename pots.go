package monzo

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type PotsService service

var (
	ErrPotInvalidID = errors.New("pot id cannot be empty")

	ErrPotInvalidSourceAccountID = errors.New("source account id cannot be empty")

	ErrPotInvalidDepositAmount = errors.New("deposit amount must be a positive number")

	ErrPotInvalidWithdrawAmount = errors.New("withdraw amount must be a positive number")

	ErrPotInvalidDedupeID = errors.New("dedupe id must not be empty")

	ErrPotClientNil = errors.New("pot client is not configured, use pots client instead")
)

type Pot struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Style             string  `json:"style"`
	Balance           int64   `json:"balance"`
	Currency          string  `json:"currency"`
	Type              string  `json:"type"`
	ProductID         string  `json:"product_id"`
	CurrentAccountID  string  `json:"current_account_id"`
	CoverImageURL     string  `json:"cover_image_url"`
	ISAWrapper        string  `json:"isa_wrapper"`
	RoundUp           bool    `json:"round_up"`
	RoundUpMultiplier float64 `json:"round_up_multiplier"`
	IsTaxPot          bool    `json:"is_tax_pot"`
	Created           string  `json:"created"`
	Updated           string  `json:"updated"`
	Deleted           bool    `json:"deleted"`
	Locked            bool    `json:"locked"`
	AvailableForBills bool    `json:"available_for_bills"`
	HasVirtualCards   bool    `json:"has_virtual_cards"`

	client *Client
}

type PotsList struct {
	Pots []Pot `json:"pots"`
}

func (p *PotsList) setClient(c *Client) {
	if p == nil {
		return
	}

	updatedPots := []Pot{}

	for _, pot := range p.Pots {
		pot.client = c
		updatedPots = append(updatedPots, pot)
	}

	p.Pots = updatedPots
}

func (s *PotsService) List(accountID string) (list *PotsList, err error) {
	list = &PotsList{}
	u := fmt.Sprintf("/pots?%s", url.Values{"current_account_id": []string{accountID}}.Encode())

	resp, err := s.client.Get(u, nil)
	err = ParseResponse(resp, err, list)

	list.setClient(s.client)

	return
}

func (s *PotsService) Deposit(potID, sourceAccountID string, amount int64, dedupeID string) (pot *Pot, err error) {
	pot = &Pot{}

	if potID == "" {
		return nil, ErrPotInvalidID
	}

	if sourceAccountID == "" {
		return nil, ErrPotInvalidSourceAccountID
	}

	if amount <= 0 {
		return nil, ErrPotInvalidDepositAmount
	}

	if dedupeID == "" {
		return nil, ErrPotInvalidDedupeID
	}

	u := fmt.Sprintf("/pots/%s/deposit", potID)

	params := url.Values{
		"source_account_id": []string{sourceAccountID},
		"amount":            []string{strconv.FormatInt(amount, 10)},
		"dedupe_id":         []string{dedupeID},
	}

	resp, err := s.client.Put(u, params)
	err = ParseResponse(resp, err, pot)
	return
}

func (p Pot) Deposit(sourceAccountID string, amount int64, dedupeID string) (*Pot, error) {
	if p.client == nil {
		return nil, ErrPotClientNil
	}

	return p.client.Pots.Deposit(p.ID, sourceAccountID, amount, dedupeID)
}

func (s *PotsService) Withdraw(potID, destinationAccountID string, amount int64, dedupeID string) (pot *Pot, err error) {
	pot = &Pot{}

	if potID == "" {
		return nil, ErrPotInvalidID
	}

	if destinationAccountID == "" {
		return nil, ErrPotInvalidSourceAccountID
	}

	if amount <= 0 {
		return nil, ErrPotInvalidWithdrawAmount
	}

	if dedupeID == "" {
		return nil, ErrPotInvalidDedupeID
	}

	u := fmt.Sprintf("/pots/%s/withdraw", potID)

	params := url.Values{
		"destination_account_id": []string{destinationAccountID},
		"amount":                 []string{strconv.FormatInt(amount, 10)},
		"dedupe_id":              []string{dedupeID},
	}

	resp, err := s.client.Put(u, params)
	err = ParseResponse(resp, err, pot)
	return
}

func (p Pot) Withdraw(destinationAccountID string, amount int64, dedupeID string) (*Pot, error) {
	if p.client == nil {
		return nil, ErrPotClientNil
	}

	return p.client.Pots.Withdraw(p.ID, destinationAccountID, amount, dedupeID)
}
