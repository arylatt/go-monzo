package monzo

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// A pot is a place to keep some money separate from the main spending account.
type PotsService service

var (
	// ErrPotInvalidID is returned if a null/empty Pot ID is supplied.
	ErrPotInvalidID = errors.New("pot id cannot be empty")

	// ErrPotInvalidSourceAccountID is returned if a null/empty Source Account ID is supplied.
	ErrPotInvalidSourceAccountID = errors.New("source account id cannot be empty")

	// ErrPotInvalidDepositAmount is returned if a zero/negative deposit amount is supplied.
	ErrPotInvalidDepositAmount = errors.New("deposit amount must be a positive number")

	// ErrPotInvalidWithdrawAmount is returned if a zero/negative withdrawal amount is supplied.
	ErrPotInvalidWithdrawAmount = errors.New("withdraw amount must be a positive number")

	// ErrPotInvalidDedupeID is returned if a null/empty deduplication ID is supplied.
	ErrPotInvalidDedupeID = errors.New("dedupe id must not be empty")

	// ErrPotClientNil is returned if Pot object has a nil client configured (e.g. if the Pot object was manually created).
	ErrPotClientNil = errors.New("pot client is not configured, use pots client instead")
)

// Pot represents a pot object provided by the Monzo API.
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

// PotsList represents the response from the Monzo API for a list of pots.
type PotsList struct {
	Pots []Pot `json:"pots"`
}

// setClient is an internal helper to ensure all pots have the client attached to them for later usage.
func (p *PotsList) setClient(c *Client) {
	if p == nil {
		return
	}

	for i := range p.Pots {
		p.Pots[i].client = c
	}
}

// Returns a list of pots owned by the currently authorised user that are associated with the specified account.
func (s *PotsService) List(accountID string) (list *PotsList, err error) {
	list = &PotsList{}
	u := fmt.Sprintf("/pots?%s", url.Values{"current_account_id": []string{accountID}}.Encode())

	resp, err := s.client.Get(u, nil)
	err = ParseResponse(resp, err, list)

	list.setClient(s.client)

	return
}

// WARNING: Undocumented Monzo API - may be subject to change!
//
// Returns a specific pot owned by the currently authorised user with the given pot ID.
func (s *PotsService) Get(potID string) (pot *Pot, err error) {
	pot = &Pot{}
	u := fmt.Sprintf("/pots/%s", potID)

	resp, err := s.client.Get(u, nil)
	err = ParseResponse(resp, err, pot)

	pot.client = s.client

	return
}

// Move money from an account owned by the currently authorised user into one of their pots.
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

// Move money from an account owned by the currently authorised user into one of their pots.
//
// Pot.Deposit is a convenience method. It is the same as calling Pots.Deposit(pot.ID, sourceAccountID, amount, dedupeID).
func (p Pot) Deposit(sourceAccountID string, amount int64, dedupeID string) (*Pot, error) {
	if p.client == nil {
		return nil, ErrPotClientNil
	}

	return p.client.Pots.Deposit(p.ID, sourceAccountID, amount, dedupeID)
}

// Move money from a pot owned by the currently authorised user into one of their accounts.
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

// Move money from a pot owned by the currently authorised user into one of their accounts.
//
// Pot.Withdraw is a convenience method. It is the same as calling Pots.Withdraw(pot.ID, destinationAccountID, amount, dedupeID).
func (p Pot) Withdraw(destinationAccountID string, amount int64, dedupeID string) (*Pot, error) {
	if p.client == nil {
		return nil, ErrPotClientNil
	}

	return p.client.Pots.Withdraw(p.ID, destinationAccountID, amount, dedupeID)
}
