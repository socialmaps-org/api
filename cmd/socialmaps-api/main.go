package main

import (
	"context"
	"log/slog"
	"net/http"

	_ "golang.org/x/crypto/x509roots/fallback"

	"golang.socialmaps.org/api/internal/database"
	"golang.socialmaps.org/api/internal/env"
	"golang.socialmaps.org/api/internal/method"
	"golang.socialmaps.org/api/internal/mistral"
	"golang.socialmaps.org/api/internal/model"
	"golang.socialmaps.org/api/internal/moderation"
	"golang.socialmaps.org/api/internal/web"
)

func main() {
	db := database.Open(env.Var.DatabaseDSN)
	defer db.Close()
	qs := model.New(db)

	// Start the background jobs
	ctx := context.Background()
	mod := moderation.NewMistralLarge2512v1(mistral.NewClient(env.Var.MistralSecret))
	go moderation.Process(ctx, qs, mod)

	// Run the API server
	authr := web.NewAuthenticator(
		env.Var.OAuth2IntrospectURL,
		env.Var.OAuth2ClientID,
		env.Var.OAuth2ClientSecret,
	)

	mux := method.Mux(authr, qs)

	slog.Info("LISTENING", "listen_addr", env.Var.ListenAddr)
	err := http.ListenAndServe(env.Var.ListenAddr, method.CORSMiddleware(mux))
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
