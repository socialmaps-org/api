package handler

import (
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/render"
	"codeberg.org/socialmaps/auth/internal/resource"
	"codeberg.org/socialmaps/auth/internal/web"
)

type CreateReview struct {
	Handler
}

type createReviewArgs struct {
	PlaceID string `path:"place_id"`
	Liked   bool   `json:"liked"`
	Comment string `json:"comment"`
}

func (h *CreateReview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var args createReviewArgs
	if err := web.Parse(r, &args); err != nil {
		web.JSON(w, http.StatusBadRequest, err)
		return
	}
	if err := h.validate(args); err != nil {
		web.JSON(w, err.StatusCode, err)
		return
	}

	rvwM := model.CreateReview(r.Context(), h.DB, args.PlaceID, "usr_foo", args.Liked, args.Comment)

	rvwR := render.Review(rvwM)

	web.JSON(w, http.StatusOK, rvwR)
}

func (h *CreateReview) validate(args createReviewArgs) *web.Error {
	if !model.IsValidPlaceID(args.PlaceID) {
		return &web.Error{
			StatusCode: http.StatusBadRequest,
			Resource:   resource.Error{Message: "invalid place ID; a place ID starts with a `plc_` prefix"},
		}
	}

	return nil
}
