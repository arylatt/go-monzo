package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var whoami = &cobra.Command{
	Use:   "whoami",
	Short: "Check auth status",
	RunE:  whoamiRunE,
}

func init() {
	root.AddCommand(whoami)
}

func whoamiRunE(cmd *cobra.Command, args []string) (err error) {
	who, err := _client.Whoami()
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
