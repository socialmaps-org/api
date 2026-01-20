package method

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/web"
)

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func Mux(authr web.Authenticator, qs *model.Queries, overpassEndpoint string) *http.ServeMux {
	c := Common{QS: qs}

	mux := http.NewServeMux()

	config := huma.DefaultConfig("Social Maps API", "0.20251213.0")
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"oauth2": {
			Type: "oauth2",
			Flows: &huma.OAuthFlows{
				AuthorizationCode: &huma.OAuthFlow{
					AuthorizationURL: "https://auth.socialmaps.org/realms/socialmaps/protocol/openid-connect/auth",
					TokenURL:         "https://auth.socialmaps.org/realms/socialmaps/protocol/openid-connect/token",
				},
			},
		},
	}

	api := humago.New(mux, config)

	huma.Register(api, huma.Operation{
		OperationID: "lookup_place",
		Method:      http.MethodGet,
		Path:        "/v1/places/lookup",
		Summary:     "Lookup a Place",
		Description: "desc",
		Tags:        []string{"Places"},
	}, (&LookupPlace{
		Common:           c,
		OverpassEndpoint: overpassEndpoint,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID: "retrieve_place",
		Method:      http.MethodGet,
		Path:        "/v1/places/{place_id}",
		Summary:     "Retrieve a Place",
		Description: "desc",
		Tags:        []string{"Places"},
	}, (&RetrievePlace{
		Common: c,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID: "list_reviews_of_place",
		Method:      http.MethodGet,
		Path:        "/v1/places/{place_id}/reviews",
		Summary:     "List Reviews of a Place",
		Description: "desc",
		Tags:        []string{"Reviews"},
	}, (&ListReviews{
		Common: c,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID:   "create_review",
		Method:        http.MethodPost,
		Path:          "/v1/places/{place_id}/reviews",
		Summary:       "Create a Review",
		Description:   "desc",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusAccepted,
		Middlewares: huma.Middlewares{
			GetAuthMiddleware(api, authr, qs),
		},
		Security: []map[string][]string{
			{"oauth2": {"reviews:write"}},
		},
	}, (&CreateReview{
		Common: c,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID: "update_review",
		Method:      http.MethodPut,
		Path:        "/v1/reviews/{review_id}",
		Summary:     "Update a Review",
		Description: "desc",
		Tags:        []string{"Reviews"},
		Middlewares: huma.Middlewares{
			GetAuthMiddleware(api, authr, qs),
		},
		Security: []map[string][]string{
			{"oauth2": {"reviews:write"}},
		},
	}, (&UpdateReview{
		Common: c,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID:   "delete_review",
		Method:        http.MethodDelete,
		Path:          "/v1/reviews/{review_id}",
		Summary:       "Delete a Review",
		Description:   "desc",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusNoContent,
		Middlewares: huma.Middlewares{
			GetAuthMiddleware(api, authr, qs),
		},
		Security: []map[string][]string{
			{"oauth2": {"reviews:write"}},
		},
	}, (&DeleteReview{
		Common: c,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID:   "like_review",
		Method:        http.MethodPut,
		Path:          "/v1/reviews/{review_id}/like",
		Summary:       "Like a Review",
		Description:   "desc",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusNoContent,
		Middlewares: huma.Middlewares{
			GetAuthMiddleware(api, authr, qs),
		},
		Security: []map[string][]string{
			{"oauth2": {"reviews:write"}},
		},
	}, (&LikeReview{
		Common: c,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID:   "unlike_review",
		Method:        http.MethodPut,
		Path:          "/v1/reviews/{review_id}/unlike",
		Summary:       "Unlike a Review",
		Description:   "desc",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusNoContent,
		Middlewares: huma.Middlewares{
			GetAuthMiddleware(api, authr, qs),
		},
		Security: []map[string][]string{
			{"oauth2": {"reviews:write"}},
		},
	}, (&UnlikeReview{
		Common: c,
	}).Execute)

	return mux
}
