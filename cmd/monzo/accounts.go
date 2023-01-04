package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/arylatt/go-monzo"
	"github.com/spf13/cobra"
)

var (
	accountsValidArgs = []string{
		"",
		string(monzo.AccountTypeUKRetail),
		string(monzo.AccountTypeUKRetailJoint),
	}

	accounts = &cobra.Command{
		Use:       "accounts [account-type]",
		Short:     "List accounts",
		GroupID:   "accounts",
		PreRunE:   accountsPreRunE,
		RunE:      accountsRunE,
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: accountsValidArgs,
	}

	ErrAccountTypeInvalid = fmt.Errorf("account type invalid. valid types [%s]", strings.Join(accountsValidArgs[1:], ", "))
)

func init() {
	root.AddGroup(&cobra.Group{ID: "accounts", Title: "Accounts"})
	root.AddCommand(accounts)
}

func accountsPreRunE(cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return nil
	}

	for _, arg := range cmd.ValidArgs[1:] {
		if args[0] == arg {
			return
		}
	}

	return ErrAccountTypeInvalid
}

func accountsRunE(cmd *cobra.Command, args []string) (err error) {
	accTypes := []monzo.AccountType{}

	if len(args) != 0 {
		accTypes = append(accTypes, monzo.AccountType(args[0]))
	}

	who, err := _client.Accounts.List(accTypes...)
	if err != nil {
		return
	}

	data, err := json.MarshalIndent(who, "", "  ")
	if err != nil {
		return
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", data)
	return
}
