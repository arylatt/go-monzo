package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	balance = &cobra.Command{
		Use:     "balance --account-id",
		Short:   "Show balance",
		GroupID: "balance",
		RunE:    balanceRunE,
	}
)

func init() {
	balance.Flags().AddFlagSet(FlagSets["account"])

	root.AddGroup(&cobra.Group{ID: "balance", Title: "Balance"})
	root.AddCommand(balance)
}

func balanceRunE(cmd *cobra.Command, args []string) (err error) {
	balance, err := _client.Balance.Get(viper.GetString("account-id"))
	if err != nil {
		return
	}

	data, err := json.MarshalIndent(balance, "", "  ")
	if err != nil {
		return
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", data)
	return
}
