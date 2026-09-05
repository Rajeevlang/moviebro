package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	Openidd      *oidc.Provider
	Config       *oauth2.Config
	Oidcverifier *oidc.IDTokenVerifier
)

func InitOpenID(clientID, clientSecret, redirectURL string) error {
	ctx := context.Background()

	var err error
	Openidd, err = oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return fmt.Errorf("failed to initialize oidc provider: %w", err)
	}
	Oidcverifier = Openidd.Verifier(&oidc.Config{ClientID: clientID})

	Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		Endpoint:     Openidd.Endpoint(), // Automatically set by go-oidc
	}

	return nil
}

// generateStateOauthCookie creates a cryptographically secure random string
func GenerateState() (string, error) {
	b := make([]byte, 32) // 32 bytes provides 256 bits of entropy
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	// Encode to URL-safe base64 so it safely travels in HTTP headers and URLs
	return base64.URLEncoding.EncodeToString(b), nil
}
