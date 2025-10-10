package main

import (
	"log"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/database"
	"codeberg.org/socialmaps/auth/internal/method"
	"codeberg.org/socialmaps/auth/internal/web"
)

func main() {
	db := database.Open("db.sqlite3")
	defer db.Close()

	c := method.Common{DB: db}

	http.Handle("GET     /v1/places/lookup", web.MethodHandler(&method.LookupPlace{Common: c}))
	http.Handle("GET     /v1/places/{place_id}", web.MethodHandler(&method.RetrievePlace{Common: c}))
	http.Handle("GET     /v1/places/{place_id}/reviews", web.MethodHandler(&method.ListReviews{Common: c}))
	http.Handle("POST    /v1/places/{place_id}/reviews", web.MethodHandler(&method.CreateReview{Common: c}))
	http.Handle("PUT     /v1/reviews/{review_id}", web.MethodHandler(&method.UpdateReview{Common: c}))
	http.Handle("DELETE  /v1/reviews/{review_id}", web.MethodHandler(&method.DeleteReview{Common: c}))
	http.Handle("PUT     /v1/reviews/{review_id}/like", web.MethodHandler(&method.LikeReview{Common: c}))
	http.Handle("PUT     /v1/reviews/{review_id}/unlike", web.MethodHandler(&method.UnlikeReview{Common: c}))

	log.Println("serving...")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
