package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/arylatt/go-monzo"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var (
	_client *monzo.Client

	root = &cobra.Command{
		Use:                "monzo",
		Short:              "CLI for interacting with Monzo APIs",
		PersistentPreRunE:  rootPersistentPreRunE,
		PersistentPostRunE: rootPersistentPostRunE,
	}

	genDocs = &cobra.Command{
		Use:   "generate-docs [docs-dir]",
		Short: "Generate Markdown docs for CLI",
		RunE:  genDocsRunE,
		Args:  cobra.MaximumNArgs(1),
	}

	version   = "dev"
	userAgent = "monzo-cli/" + version
)

func init() {
	viper.SetEnvPrefix("MONZO")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	root.AddCommand(genDocs)
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

	if err = os.MkdirAll(viper.GetString("home-dir"), os.ModeDir|0600); err != nil {
		return
	}

	if cmd.Name() == "login" || cmd.Name() == "logout" || cmd == cmd.Root() {
		return
	}

	token := &Token{}
	err = LoadCache(CacheFileToken, token)
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

	token := &Token{}
	LoadCache(CacheFileToken, token)

	token.Token, err = _client.Token()
	if err != nil {
		return
	}

	return token.Save()
}

func genDocsRunE(cmd *cobra.Command, args []string) (err error) {
	dirName := "docs/"
	if len(args) > 0 {
		dirName = args[0]
	}

	err = os.MkdirAll(dirName, os.ModeDir)
	if err != nil {
		return
	}

	return doc.GenMarkdownTree(root, dirName)
}

func BuildClient(ctx context.Context, token *Token) (c *monzo.Client) {
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

	c = monzo.New(oauth2.NewClient(ctx, ts))

	c.UserAgent = fmt.Sprintf("%s, %s", userAgent, monzo.DefaultUserAgent)

	return
}
