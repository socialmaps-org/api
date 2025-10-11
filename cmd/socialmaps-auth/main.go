package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
	"github.com/ory/fosite/token/jwt"

	envvar "github.com/caarlos0/env/v11"

	"codeberg.org/socialmaps/auth/internal/crypto"
	"codeberg.org/socialmaps/auth/internal/database"
	"codeberg.org/socialmaps/auth/internal/env"
	"codeberg.org/socialmaps/auth/internal/handler"
	"codeberg.org/socialmaps/auth/internal/middleware"
)

func main() {
	var env env.AuthEnv

	err := envvar.ParseWithOptions(&env, envvar.Options{
		RequiredIfNoDef: true,
	})
	if err != nil {
		panic(err)
	}

	if len(env.CookieSecret) < 32 {
		panic("environment variable `COOKIE_SECRET` is too short: must be at least 32 bytes when decoded")
	}

	oauth2ClientConfig := &oauth2.Config{
		ClientID:     env.OSMClientID,
		ClientSecret: env.OSMClientSecret,
		Scopes:       []string{"openid"},
		Endpoint:     endpoints.OpenStreetMap,
		RedirectURL:  fmt.Sprintf("http://%s/auth/callback", env.Host),
	}

	store := storage.NewMemoryStore()

	secret := []byte("ndZC4kpqi6f*!3!v2TDYpfqbAyMATT")
	hashedSecret, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	store.Clients["my-client"] = &fosite.DefaultClient{
		ID:     "my-client",
		Secret: hashedSecret,
		RedirectURIs: []string{
			"http://127.0.0.1:8000/auth/callback",
			"https://app.insomnia.rest/oauth/redirect",
		},
		Scopes: []string{"openid", "offline_access"},
		GrantTypes: []string{
			"authorization_code",
			"refresh_token",
		},
		ResponseTypes: []string{"code"},
	}

	privateKey := crypto.LoadRSAKey(env.Oauth2PrivateKeyFile)

	config := fosite.Config{
		IDTokenIssuer:              "https://auth.socialmaps.org",
		EnforcePKCE:                true,
		AccessTokenLifespan:        env.AccessTokenLifespan,
		RefreshTokenLifespan:       env.RefreshTokenLifespan,
		AuthorizeCodeLifespan:      env.AuthCodeLifespan,
		IDTokenLifespan:            env.IDTokenLifespan,
		RefreshTokenScopes:         []string{"offline_access"},
		GlobalSecret:               env.OAuth2Secret,
		SendDebugMessagesToClients: true,
	}
	oauth2Server := compose.Compose(
		&config,
		store,
		&compose.CommonStrategy{
			CoreStrategy: compose.NewOAuth2HMACStrategy(&config),
			Signer: &jwt.DefaultSigner{GetPrivateKey: func(context.Context) (interface{}, error) {
				return privateKey, nil
			}},
		},
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		compose.OAuth2TokenIntrospectionFactory,
		compose.OAuth2TokenRevocationFactory,
		compose.OAuth2PKCEFactory,
	)

	db := database.Open("db.sqlite3")
	defer db.Close()

	c := handler.Common{
		DB:  db,
		Env: env,
	}

	http.Handle("GET /{$}", middleware.CanonicalLog(&handler.Index{Common: c}))
	http.Handle("GET /login", middleware.CanonicalLog(&handler.Login{Common: c}))
	http.Handle("GET /logout", middleware.CanonicalLog(&handler.Logout{Common: c}))
	// OAuth2 Client
	http.Handle("GET /auth/openstreetmap", middleware.CanonicalLog(&handler.AuthnOSM{Common: c, OAuth2Config: oauth2ClientConfig}))
	http.Handle("GET /auth/callback", middleware.CanonicalLog(&handler.AuthnCallback{Common: c, OAuth2Config: oauth2ClientConfig}))
	// OAuth2 Server
	http.Handle("/oauth2/authorize", middleware.CanonicalLog(&handler.Authorize{Common: c, OAuth2Server: oauth2Server}))
	http.Handle("/oauth2/token", middleware.CanonicalLog(&handler.Token{Common: c, OAuth2Server: oauth2Server}))

	slog.Info("socialmaps-auth serving...")
	err = http.ListenAndServe(env.Host, nil)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
