package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/arylatt/go-monzo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	transactions = &cobra.Command{
		Use:     "transactions",
		Short:   "Get and annotate transactions",
		GroupID: "transactions",
	}

	transactionsGet = &cobra.Command{
		Use:     "get [transaction id]",
		Short:   "Get all transactions or a specific transaction by ID",
		PreRunE: transactionsGetPreRunE,
		RunE:    transactionsGetRunE,
		Args:    cobra.MaximumNArgs(1),
	}
)

func init() {
	transactions.PersistentFlags().Bool("no-cache", false, "Bypass transactions cache and force call to API")
	transactions.PersistentFlags().Bool("expand-merchants", false, "Fetch expanded Merchants data")

	viper.BindPFlags(transactions.PersistentFlags())

	root.AddGroup(&cobra.Group{ID: "transactions", Title: "Transactions"})
	root.AddCommand(transactions)

	transactionsGet.Flags().StringP("account-id", "a", "", "Account ID to list transactions for")

	viper.BindPFlags(transactionsGet.Flags())

	transactions.AddCommand(transactionsGet)
}

type Transactions map[string]*monzo.TransactionList

func (t Transactions) Save() error {
	filePath := path.Join(viper.GetString("home-dir"), "transactions.json")

	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0700)
}

func LoadTransactions() (Transactions, error) {
	filePath := path.Join(viper.GetString("home-dir"), "transactions.json")

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	t := Transactions{}
	err = json.NewDecoder(file).Decode(&t)

	return t, err
}

func transactionsGetPreRunE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && viper.GetString("account-id") == "" {
		return errors.New("--account-id flag is required to list transactions")
	}

	return nil
}

func transactionsGetRunE(cmd *cobra.Command, args []string) (err error) {
	transactions, _ := LoadTransactions()

	if len(args) == 0 && (transactions == nil || viper.GetBool("no-cache")) {
		liveTxs, err := _client.Transactions.List(viper.GetString("account-id"), viper.GetBool("expand-merchants"), nil)
		if err != nil {
			return err
		}

		transactions[viper.GetString("account-id")] = liveTxs

		transactions.Save()
	}

	if len(args) == 0 {
		data, err := json.MarshalIndent(transactions[viper.GetString("account-id")], "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", data)
		return nil
	}

	if transactions == nil || viper.GetBool("no-cache") {
		tx, err := _client.Transactions.Get(args[0], viper.GetBool("expand-merchants"))
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(tx, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", data)

		if transactions == nil {
			transactions = Transactions{}
		}

		if transactions[tx.Transaction.AccountID] == nil {
			transactions[tx.Transaction.AccountID] = &monzo.TransactionList{}
		}

		transactions[tx.Transaction.AccountID].Transactions = append(transactions[tx.Transaction.AccountID].Transactions, tx.Transaction)

		transactions.Save()

		return nil
	}

	for _, txns := range transactions {
		for _, tx := range txns.Transactions {
			if tx.ID == args[0] {
				txSingle := monzo.TransactionSingle{Transaction: tx}

				data, err := json.MarshalIndent(txSingle, "", "  ")
				if err != nil {
					return err
				}

				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", data)
				return nil
			}
		}
	}

	return
}
