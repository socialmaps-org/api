package method

import (
	"context"
	"net/http"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"codeberg.org/socialmaps/api/internal/web"
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
	usr := web.GetAuthUser(ctx)

	plc := model.LoadPlaceByID(ctx, m.DB, args.PlaceID)
	if plc == nil {
		return web.NewEmptyResponse(http.StatusNotFound)
	}

	rvwM := model.CreateReview(ctx, m.DB, args.PlaceID, usr.ID, args.Liked, args.Comment)

	rvwR := render.Review(rvwM)

	return web.NewJSONResponse(http.StatusAccepted, rvwR)
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
