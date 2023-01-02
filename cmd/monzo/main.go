package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/arylatt/go-monzo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	Version = "dev"
)

var _client *monzo.Client

var root = &cobra.Command{
	Use:                "monzo",
	Short:              "CLI for interacting with Monzo APIs",
	PersistentPreRunE:  rootPersistentPreRunE,
	PersistentPostRunE: rootPersistentPostRunE,
}

func init() {
	viper.SetEnvPrefix("MONZO")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func main() {
	if root.Execute() != nil {
		os.Exit(1)
	}
}

func rootPersistentPreRunE(cmd *cobra.Command, args []string) (err error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return
	}

	viper.SetDefault("home-dir", path.Join(userHome, "/.monzo/"))

	if err = os.MkdirAll(viper.GetString("home-dir"), os.ModeDir); err != nil {
		return
	}

	if cmd.Name() == "login" || cmd == cmd.Root() {
		return
	}

	token, err := LoadToken()
	if err != nil {
		return fmt.Errorf("not authenticated, try running monzo login - %w", err)
	}

	_client = BuildClient(cmd.Context(), token)

	return
}

func rootPersistentPostRunE(cmd *cobra.Command, args []string) (err error) {
	if _client == nil {
		return
	}

	token, _ := LoadToken()

	token.Token, err = _client.Token()
	if err != nil {
		return
	}

	return token.Save()
}

func BuildClient(ctx context.Context, token *Token) *monzo.Client {
	var ts oauth2.TokenSource

	if token.RefreshToken == "" {
		ts = oauth2.StaticTokenSource(token.Token)
	} else {
		config := oauth2.Config{
			ClientID:     token.ClientID,
			ClientSecret: token.ClientSecret,
			Endpoint:     monzo.OAuth2Endpoint,
		}

		ts = config.TokenSource(ctx, token.Token)
	}

	return monzo.New(oauth2.NewClient(ctx, ts))
}
