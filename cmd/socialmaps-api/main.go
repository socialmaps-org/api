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
		DatabaseDSN      string `env:"DATABASE_DSN" envDefault:"socialmaps-api.sqlite3"`
		ListenAddr       string `env:"LISTEN_ADDR" envDefault:"127.0.0.1:8080"`
		OverpassEndpoint string `env:"OVERPASS_ENDPOINT" envDefault:"https://overpass-api.de/api/interpreter"`
	}
	err := env.Parse(&envvars)
	if err != nil {
		panic(err)
	}

	db := database.Open(envvars.DatabaseDSN)
	defer db.Close()

	c := method.Common{DB: db}

	http.Handle(
		"GET     /v1/places/lookup",
		web.MethodHandler(&method.LookupPlace{
			Common:           c,
			OverpassEndpoint: envvars.OverpassEndpoint,
		}),
	)
	http.Handle(
		"GET     /v1/places/{place_id}",
		web.MethodHandler(&method.RetrievePlace{
			Common: c,
		}),
	)
	http.Handle(
		"GET     /v1/places/{place_id}/reviews",
		web.MethodHandler(&method.ListReviews{
			Common: c,
		}),
	)
	http.Handle(
		"POST    /v1/places/{place_id}/reviews",
		web.MethodHandler(&method.CreateReview{
			Common: c,
		}),
	)
	http.Handle(
		"PUT     /v1/reviews/{review_id}",
		web.MethodHandler(&method.UpdateReview{
			Common: c,
		}),
	)
	http.Handle(
		"DELETE  /v1/reviews/{review_id}",
		web.MethodHandler(&method.DeleteReview{
			Common: c,
		}),
	)
	http.Handle(
		"PUT     /v1/reviews/{review_id}/like",
		web.MethodHandler(&method.LikeReview{
			Common: c,
		}),
	)
	http.Handle(
		"PUT     /v1/reviews/{review_id}/unlike",
		web.MethodHandler(&method.UnlikeReview{
			Common: c,
		}),
	)

	slog.Info("LISTENING", "listen_addr", envvars.ListenAddr)
	err = http.ListenAndServe(envvars.ListenAddr, nil)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
