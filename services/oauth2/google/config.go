package google

import (
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthProvider struct {
	*oauth2.Config
	ClientRedirectURL string
	OAuthStateString  string
}

type ConfigOptions struct {
	ClientID          string
	ClientSecret      string
	RedirectURL       string
	ClientRedirectURL string
	OAuthStateString  string
}

var OAuthConfig *OAuthProvider

func Config(config *ConfigOptions) (*OAuthProvider, error) {
	if OAuthConfig != nil {
		return OAuthConfig, nil
	}

	if config.ClientID == "" {
		return nil, errors.New("client_id is required")
	}

	if config.ClientSecret == "" {
		return nil, errors.New("client_secret is required")
	}

	if config.RedirectURL == "" {
		return nil, errors.New("redirect_url is required")
	}

	if config.ClientRedirectURL == "" {
		return nil, errors.New("client_redirect_url is required")
	}

	OAuthConfig = &OAuthProvider{
		Config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		ClientRedirectURL: config.ClientRedirectURL,
		OAuthStateString:  config.OAuthStateString,
	}

	return OAuthConfig, nil
}
