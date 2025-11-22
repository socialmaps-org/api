package main

import (
	"log/slog"
	"net/http"

	env "github.com/caarlos0/env/v11"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/method"
	"codeberg.org/socialmaps/api/internal/web"
)

func main() {
	var envvars struct {
		OAuth2IntrospectURL string `env:"OAUTH2_INTROSPECT_URL"`
		OAuth2ClientID      string `env:"OAUTH2_CLIENT_ID" envDefault:"org.socialmaps.api"`
		OAuth2ClientSecret  string `env:"OAUTH2_CLIENT_SECRET"`
		DatabaseDSN         string `env:"DATABASE_DSN" envDefault:"socialmaps-api.sqlite3"`
		ListenAddr          string `env:"LISTEN_ADDR" envDefault:"127.0.0.1:8080"`
		OverpassEndpoint    string `env:"OVERPASS_ENDPOINT" envDefault:"https://overpass-api.de/api/interpreter"`
	}
	err := env.Parse(&envvars)
	if err != nil {
		panic(err)
	}

	authr := web.NewAuthenticator(
		envvars.OAuth2IntrospectURL,
		envvars.OAuth2ClientID,
		envvars.OAuth2ClientSecret,
	)

	db := database.Open(envvars.DatabaseDSN)
	defer db.Close()

	mux := method.Mux(authr, db, envvars.OverpassEndpoint)

	slog.Info("LISTENING", "listen_addr", envvars.ListenAddr)
	err = http.ListenAndServe(envvars.ListenAddr, mux)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
