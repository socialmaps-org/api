package main

import (
	"log"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/database"
	"codeberg.org/socialmaps/auth/internal/handler"
)

func main() {
	db := database.Open("db.sqlite3")
	defer db.Close()

	h := handler.Handler{DB: db}

	http.Handle("GET     /v1/places/lookup", &handler.LookupPlace{h})
	http.Handle("GET     /v1/places/{place_id}", &handler.RetrievePlace{h})
	http.Handle("GET     /v1/places/{place_id}/reviews", &handler.ListReviews{h})
	http.Handle("POST    /v1/places/{place_id}/reviews", &handler.CreateReview{h})
	http.Handle("PUT     /v1/reviews/{review_id}", &handler.UpdateReview{h})
	http.Handle("DELETE  /v1/reviews/{review_id}", &handler.DeleteReview{h})
	http.Handle("PUT     /v1/reviews/{review_id}/like", &handler.LikeReview{h})
	http.Handle("PUT     /v1/reviews/{review_id}/unlike", &handler.UnlikeReview{h})

	log.Println("serving...")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
