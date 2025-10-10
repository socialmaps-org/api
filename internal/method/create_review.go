package method

import (
	"context"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/render"
	"codeberg.org/socialmaps/auth/internal/resource"
	"codeberg.org/socialmaps/auth/internal/web"
)

type CreateReview struct {
	Common
}

type createReviewArgs struct {
	PlaceID string `path:"place_id"`
	Liked   bool   `json:"liked"`
	Comment string `json:"comment"`
}

func (m *CreateReview) Execute(ctx context.Context, args *createReviewArgs) *web.Response {
	rvwM := model.CreateReview(ctx, m.DB, args.PlaceID, "usr_foo", args.Liked, args.Comment)

	rvwR := render.Review(rvwM)

	return web.NewJSONResponse(http.StatusOK, rvwR)
}

func (m *CreateReview) Validate(args *createReviewArgs) *web.Response {
	if !model.IsValidPlaceID(args.PlaceID) {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: resource.ErrMsgInvalidPlaceID,
		})
	}

	return nil
}

func (m *CreateReview) NewArgs() *createReviewArgs {
	return &createReviewArgs{}
}
