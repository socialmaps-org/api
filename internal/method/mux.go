package method

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/multiline"
	"codeberg.org/socialmaps/api/internal/web"
)

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func Mux(authr web.Authenticator, qs *model.Queries) *http.ServeMux {
	c := Common{QS: qs}

	mux := http.NewServeMux()

	config := huma.DefaultConfig("Social Maps API", "0.20251213.0")
	// Unset CreateHooks in the default config to stop using
	// `SchemaLinkTransformer` that adds `$schema` fields to responses.
	// See <https://github.com/danielgtaylor/huma/issues/428>
	config.CreateHooks = nil
	config.Components.Schemas = huma.NewMapRegistry(
		"#/components/schemas/",
		func(t reflect.Type, hint string) string {
			if strings.HasPrefix(hint, "override:") {
				return strings.TrimPrefix(hint, "override:")
			}
			return huma.DefaultSchemaNamer(t, hint)
		},
	)
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
		Description: multiline.Dedent(`
			Look up a Place by its name and geo-coordinates.

			We require clients to look up Places by their "name + location" because OpenStreetMap IDs are [not permanent](https://wiki.openstreetmap.org/wiki/Permanent_ID). You can read more about it in the [docs](https://docs.socialmaps.org/https://docs.socialmaps.org/lookups/).
		`),
		Tags: []string{"Places"},
	}, (&LookupPlace{
		Common: c,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID: "query_places",
		Method:      http.MethodGet,
		Path:        "/v1/places",
		Summary:     "Query Places",
		Description: multiline.Dedent(`
			Query Places in a bounding-box filtered by their OpenStreetMap [tags](https://wiki.openstreetmap.org/wiki/Tags).
		`),
		Tags: []string{"Places"},
	}, (&QueryPlaces{
		Common: c,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID: "retrieve_place",
		Method:      http.MethodGet,
		Path:        "/v1/places/{place_id}",
		Summary:     "Retrieve a Place",
		Description: multiline.Dedent(`
			Retrieve a Place by ID.
		`),
		Tags: []string{"Places"},
	}, (&RetrievePlace{
		Common: c,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID: "list_reviews_of_place",
		Method:      http.MethodGet,
		Path:        "/v1/places/{place_id}/reviews",
		Summary:     "List Reviews of a Place",
		Description: multiline.Dedent(`
			List Reviews of a Place by Place ID.
		`),
		Tags: []string{"Reviews"},
	}, (&ListReviews{
		Common: c,
	}).Execute)

	huma.Register(api, huma.Operation{
		OperationID: "create_review",
		Method:      http.MethodPost,
		Path:        "/v1/places/{place_id}/reviews",
		Summary:     "Create a Review",
		Description: multiline.Dedent(`
			Create a Review of a Place.
		`),
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
		Description: multiline.Dedent(`
			Update a Review.

			You cannot update a Review more than an hour after creating it. You can, however, delete a Review any time and create a new one.
		`),
		Tags: []string{"Reviews"},
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
		OperationID: "delete_review",
		Method:      http.MethodDelete,
		Path:        "/v1/reviews/{review_id}",
		Summary:     "Delete a Review",
		Description: multiline.Dedent(`
			Delete a Review.
		`),
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
		OperationID: "like_review",
		Method:      http.MethodPut,
		Path:        "/v1/reviews/{review_id}/like",
		Summary:     "Like a Review",
		Description: multiline.Dedent(`
			Like a **Review**.

			A **User** cannot like their own **Review**.
		`),
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
		OperationID: "unlike_review",
		Method:      http.MethodPut,
		Path:        "/v1/reviews/{review_id}/unlike",
		Summary:     "Unlike a Review",
		Description: multiline.Dedent(`
			Unlike a Review.
		`),
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
