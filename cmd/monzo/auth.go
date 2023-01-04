package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/arylatt/go-monzo"
	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var (
	login = &cobra.Command{
		Use:     "login --token | --client-id --client-secret",
		Short:   "Authenticate to the Monzo API",
		GroupID: "auth",
		PreRunE: loginPreRunE,
		RunE:    loginRunE,
	}

	refreshToken = &cobra.Command{
		Use:     "refresh-token",
		Short:   "Force refresh of token if using OAuth2",
		GroupID: "auth",
		PreRunE: refreshTokenPreRunE,
		RunE:    refreshTokenRunE,
	}

	logout = &cobra.Command{
		Use:     "logout",
		Short:   "Delete all cached data",
		GroupID: "auth",
		RunE:    logoutRunE,
	}

	ErrLoginAuthTypesMutuallyExclusive = errors.New("cannot use --token with --client-id and --client-secret")

	ErrLoginAuthTypesOAuth2MissingPart = errors.New("--client-id and --client-secret must both be provided for oauth2")

	ErrLoginAuthTypesMissing = errors.New("--token or --client-id and --client-secret must be supplied")
)

func init() {
	login.Flags().AddFlagSet(FlagSets["login"])

	root.AddGroup(&cobra.Group{ID: "auth", Title: "Auth"})
	root.AddCommand(login)

	root.AddCommand(refreshToken)

	root.AddCommand(logout)
}

func loginPreRunE(cmd *cobra.Command, args []string) error {
	if viper.GetString("token") != "" && (viper.GetString("client-id") != "" || viper.GetString("client-secret") != "") {
		return ErrLoginAuthTypesMutuallyExclusive
	}

	if viper.GetString("token") == "" && (viper.GetString("client-id") == "" || viper.GetString("client-secret") == "") {
		return ErrLoginAuthTypesOAuth2MissingPart
	}

	if viper.GetString("token") == "" && viper.GetString("client-id") == "" && viper.GetString("client-secret") == "" {
		return ErrLoginAuthTypesMissing
	}

	return nil
}

func loginRunE(cmd *cobra.Command, args []string) (err error) {
	token := &Token{}

	if tokenStr := viper.GetString("token"); tokenStr != "" {
		token.Token = &oauth2.Token{AccessToken: tokenStr}
	} else {
		token, err = LoginOAuth2(cmd.Context(), viper.GetString("client-id"), viper.GetString("client-secret"))
	}

	if err != nil {
		return
	}

	c := BuildClient(cmd.Context(), token)
	who, err := c.Whoami()

	if err != nil {
		return
	}

	err = token.Save()
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stdout, "Authenticated to Monzo! User: %s\n\n", who.UserID)
	return nil
}

type Token struct {
	*oauth2.Token

	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

func (t *Token) Save() error {
	return SaveCache(CacheFileToken, t)
}

func LoginOAuth2(ctx context.Context, clientID, clientSecret string) (t *Token, err error) {
	state := uuid.NewString()
	tokenChan := make(chan *Token, 1)

	config := monzo.OAuth2Config(clientID, clientSecret, "http://127.0.0.1:54092/callback")

	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	srv := http.Server{
		Addr:    "127.0.0.1:54092",
		Handler: http.HandlerFunc(CallbackHandler(ctx, state, tokenChan, config)),
	}

	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	err = browser.OpenURL(authURL)
	if err != nil {
		// debug
		fmt.Fprintf(os.Stdout, "Failed to launch browser. Please copy and paste auth URL into browser:\n\n\t%s\n", authURL)
	}

	errChan := make(chan error, 1)

	go func() {
		teardownErr := []string{}
		select {
		case <-ctx.Done():
			teardownErr = append(teardownErr, ctx.Err().Error())
		case t = <-tokenChan:
		}

		if err := srv.Shutdown(ctx); err != nil {
			teardownErr = append(teardownErr, err.Error())
		}

		if t == nil {
			teardownErr = append(teardownErr, "authentication failure")
		}

		if len(teardownErr) != 0 {
			errChan <- errors.New(strings.Join(teardownErr, ", "))
			return
		}

		errChan <- nil
	}()

	err = srv.ListenAndServe()
	if err == http.ErrServerClosed {
		err = <-errChan
	} else {
		err = fmt.Errorf("%w, %w", err, <-errChan)
	}

	return
}

func CallbackHandler(ctx context.Context, state string, tokenChan chan *Token, config *oauth2.Config) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		code, reqState := r.URL.Query().Get("code"), r.URL.Query().Get("state")
		if code == "" || reqState == "" {
			http.Error(rw, "bad request - query parameters should contain code and state", http.StatusBadRequest)
			return
		}

		if state != reqState {
			http.Error(rw, "bad request - state mismatch", http.StatusBadRequest)
			return
		}

		innerToken, err := config.Exchange(ctx, code, oauth2.AccessTypeOffline)
		if err != nil {
			http.Error(rw, "authentication failure", http.StatusInternalServerError)
			tokenChan <- nil
			return
		}

		tokenChan <- &Token{
			Token:        innerToken,
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("authentication successful"))
	}
}

func refreshTokenPreRunE(cmd *cobra.Command, args []string) (err error) {
	token := &Token{}
	err = LoadCache(CacheFileToken, token)
	if err != nil {
		return
	}

	if token.ClientID == "" || token.ClientSecret == "" || token.RefreshToken == "" {
		return errors.New("cannot refresh - missing client id, client secret, or refresh token")
	}

	return nil
}

func refreshTokenRunE(cmd *cobra.Command, args []string) (err error) {
	token := &Token{}
	LoadCache(CacheFileToken, token)

	err = _client.RefreshToken()
	if err != nil {
		return
	}

	token.Token, err = _client.Token()
	if err != nil {
		return
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Token refreshed, new expiry: %s", token.Expiry.String())

	return token.Save()
}

func logoutRunE(cmd *cobra.Command, args []string) error {
	return os.RemoveAll(viper.GetString("home-dir"))
}
