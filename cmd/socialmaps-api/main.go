package main

import (
	"context"
	"log/slog"
	"net/http"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/env"
	"codeberg.org/socialmaps/api/internal/method"
	"codeberg.org/socialmaps/api/internal/mistral"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/moderation"
	"codeberg.org/socialmaps/api/internal/web"
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
