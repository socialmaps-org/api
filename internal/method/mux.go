package method

import (
	"database/sql"
	"net/http"

	"codeberg.org/socialmaps/api/internal/web"
)

func Mux(authr web.Authenticator, db *sql.DB, overpassEndpoint string) *http.ServeMux {
	c := Common{DB: db}

	mux := http.NewServeMux()
	mux.Handle(
		"GET     /v1/places/lookup",
		web.MethodHandler(&LookupPlace{
			Common:           c,
			OverpassEndpoint: overpassEndpoint,
		}),
	)
	mux.Handle(
		"GET     /v1/places/{place_id}",
		web.MethodHandler(&RetrievePlace{
			Common: c,
		}),
	)
	mux.Handle(
		"GET     /v1/places/{place_id}/reviews",
		web.MethodHandler(&ListReviews{
			Common: c,
		}),
	)
	mux.Handle(
		"POST    /v1/places/{place_id}/reviews",
		web.AuthMiddleware(db, authr, "review",
			web.MethodHandler(&CreateReview{
				Common: c,
			}),
		),
	)
	mux.Handle(
		"PUT     /v1/reviews/{review_id}",
		web.AuthMiddleware(db, authr, "review",
			web.MethodHandler(&UpdateReview{
				Common: c,
			}),
		),
	)
	mux.Handle(
		"DELETE  /v1/reviews/{review_id}",
		web.AuthMiddleware(db, authr, "review",
			web.MethodHandler(&DeleteReview{
				Common: c,
			}),
		),
	)
	mux.Handle(
		"PUT     /v1/reviews/{review_id}/like",
		web.AuthMiddleware(db, authr, "review",
			web.MethodHandler(&LikeReview{
				Common: c,
			}),
		),
	)
	mux.Handle(
		"PUT     /v1/reviews/{review_id}/unlike",
		web.AuthMiddleware(db, authr, "review",
			web.MethodHandler(&UnlikeReview{
				Common: c,
			}),
		),
	)

	return mux
}
