package monzo

import "golang.org/x/oauth2"

// OAuth2Endpoint returns the default OAuth2 endpoint for the Monzo API.
var OAuth2Endpoint = oauth2.Endpoint{
	AuthURL:   "https://auth.monzo.com/",
	TokenURL:  "https://api.monzo.com/oauth2/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

// OAuth2Config creates an OAuth2 config based on the provided client ID, client secret, and redirect URL.
func OAuth2Config(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     OAuth2Endpoint,
		RedirectURL:  redirectURL,
	}
}
