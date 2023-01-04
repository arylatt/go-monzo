package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
		Use:     "get [transaction-id]",
		Short:   "Get all transactions or a specific transaction by ID",
		PreRunE: transactionsGetPreRunE,
		RunE:    transactionsGetRunE,
		Args:    cobra.MaximumNArgs(1),
	}

	transactionAnnotate = &cobra.Command{
		Use:     "annotate transaction-id key=value...",
		Short:   "Get all transactions or a specific transaction by ID",
		PreRunE: transactionAnnotatePreRunE,
		RunE:    transactionAnnotateRunE,
		Args:    cobra.MinimumNArgs(2),
	}

	ErrTransactionIDNonNil = errors.New("transaction-id argument must be set")

	ErrInvalidMetadataKeyValue = errors.New("metadata key/value invalid. missing '=' character in argument")
)

func init() {
	transactions.PersistentFlags().AddFlagSet(FlagSets["cache"])
	transactions.PersistentFlags().AddFlagSet(FlagSets["expand"])
	transactions.PersistentFlags().AddFlagSet(FlagSets["pagination"])

	root.AddGroup(&cobra.Group{ID: "transactions", Title: "Transactions"})
	root.AddCommand(transactions)

	transactionsGet.Flags().AddFlagSet(FlagSets["account"])

	transactions.AddCommand(transactionsGet)

	transactions.AddCommand(transactionAnnotate)
}

type Transactions map[string]*monzo.TransactionList

func (t Transactions) Save() error {
	return SaveCache(CacheFileTransactions, t)
}

func (t Transactions) Find(accountID, transactionID string) *monzo.TransactionSingle {
	for acc, txns := range t {
		if acc != accountID && accountID != "" {
			continue
		}

		for _, tx := range txns.Transactions {
			if tx.ID == transactionID {
				return &monzo.TransactionSingle{Transaction: tx}
			}
		}
	}

	return nil
}

func (t Transactions) FindMulti(accountID string, page *monzo.Pagination) (filteredTxns *monzo.TransactionList) {
	filteredTxns = &monzo.TransactionList{Transactions: []monzo.Transaction{}}

	for acc, txns := range t {
		if acc != accountID {
			continue
		}

		if page == nil {
			filteredTxns.Transactions = txns.Transactions
			return
		}

		for _, tx := range txns.Transactions {
			if page.Limit != 0 && len(filteredTxns.Transactions) == page.Limit {
				return
			}

			if tx.CreatedTime().After(page.SinceTime()) && tx.CreatedTime().Before(page.BeforeTime()) {
				filteredTxns.Transactions = append(filteredTxns.Transactions, tx)
			}
		}

		return
	}

	return
}

func (t Transactions) Upsert(accountID string, transaction monzo.Transaction) {
	defer t.Save()

	if t[accountID] == nil {
		t[accountID] = &monzo.TransactionList{
			Transactions: []monzo.Transaction{
				transaction,
			},
		}

		return
	}

	for i, tx := range t[accountID].Transactions {
		if tx.ID == transaction.ID {
			t[accountID].Transactions[i] = transaction
			return
		}
	}

	t[accountID].Transactions = append(t[accountID].Transactions, transaction)
}

func (t Transactions) UpsertMulti(accountID string, transactions monzo.TransactionList) {
	for _, tx := range transactions.Transactions {
		t.Upsert(accountID, tx)
	}
}

func transactionsGetPreRunE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && viper.GetString("account-id") == "" {
		return errors.New("--account-id flag is required to list transactions")
	}

	return nil
}

func transactionsGetRunE(cmd *cobra.Command, args []string) (err error) {
	transactions := Transactions{}
	LoadCache(CacheFileTransactions, &transactions)

	accountID, expandMerchants, noCache := viper.GetString("account-id"), viper.GetBool("expand-merchants"), viper.GetBool("no-cache")
	page := BuildPagination()

	if len(args) == 0 && (len(transactions) == 0 || noCache) {
		liveTxs, err := _client.Transactions.List(accountID, expandMerchants, page)
		if err != nil {
			return err
		}

		transactions.UpsertMulti(accountID, *liveTxs)
	}

	if len(args) == 0 {
		data, err := json.MarshalIndent(transactions.FindMulti(accountID, page), "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", data)
		return nil
	}

	if len(transactions) == 0 || noCache {
		tx, err := _client.Transactions.Get(args[0], expandMerchants)
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(tx, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", data)

		transactions.Upsert(tx.Transaction.AccountID, tx.Transaction)

		return nil
	}

	txSingle := transactions.Find(accountID, args[0])

	data, err := json.MarshalIndent(txSingle, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", data)

	return
}

func transactionAnnotatePreRunE(cmd *cobra.Command, args []string) (err error) {
	if strings.TrimSpace(args[0]) == "" || strings.Contains(args[0], "=") {
		return ErrTransactionIDNonNil
	}

	return
}

func transactionAnnotateRunE(cmd *cobra.Command, args []string) (err error) {
	transactions := Transactions{}
	LoadCache(CacheFileTransactions, &transactions)

	metadata := map[string]string{}

	for _, argPairStr := range args[1:] {
		sepIndex := strings.Index(argPairStr, "=")
		if sepIndex == -1 {
			return fmt.Errorf("%w, '%s'", ErrInvalidMetadataKeyValue, argPairStr)
		}

		metadata[argPairStr[0:sepIndex]] = argPairStr[sepIndex+1:]
	}

	tx, err := _client.Transactions.Annotate(args[0], metadata)
	if err != nil {
		return
	}

	transactions.Upsert(tx.Transaction.AccountID, tx.Transaction)

	data, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", data)

	return
}
