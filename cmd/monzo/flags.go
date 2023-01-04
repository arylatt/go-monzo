package main

import (
	"strings"

	"github.com/arylatt/go-monzo"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	FlagSets = map[string]*pflag.FlagSet{}
)

func init() {
	FlagSets["pagination"] = pflag.NewFlagSet("pagination", pflag.ContinueOnError)
	FlagSets["pagination"].StringP("before", "b", "", "Pagination - return results before this date/time")
	FlagSets["pagination"].StringP("since", "s", "", "Pagination - return results since this date/time")
	FlagSets["pagination"].IntP("limit", "l", 0, "Pagination - return at most this many results")

	FlagSets["login"] = pflag.NewFlagSet("login", pflag.ContinueOnError)
	FlagSets["login"].StringP("token", "t", "", "Authenticate with static access token")
	FlagSets["login"].StringP("client-id", "c", "", "Authenticate with Client ID")
	FlagSets["login"].StringP("client-secret", "s", "", "Authenticate with Client Secret")

	FlagSets["account"] = pflag.NewFlagSet("account", pflag.ContinueOnError)
	FlagSets["account"].StringP("account-id", "a", "", "Account ID to list transactions for")

	FlagSets["cache"] = pflag.NewFlagSet("cache", pflag.ContinueOnError)
	FlagSets["cache"].Bool("no-cache", false, "Bypass transactions cache and force call to API")

	FlagSets["expand"] = pflag.NewFlagSet("expand", pflag.ContinueOnError)
	FlagSets["expand"].Bool("expand-merchants", false, "Fetch expanded Merchants data")

	for _, fs := range FlagSets {
		viper.BindPFlags(fs)
	}
}

func BuildPagination() *monzo.Pagination {
	p := &monzo.Pagination{
		Limit:  viper.GetInt("limit"),
		Since:  viper.GetString("since"),
		Before: viper.GetString("before"),
	}

	if p.Limit == 0 && strings.TrimSpace(p.Since) == "" && strings.TrimSpace(p.Before) == "" {
		return nil
	}

	return p
}
