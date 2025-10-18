package method

import (
	"context"
	"net/http"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"codeberg.org/socialmaps/api/internal/web"
)

type RetrievePlace struct {
	Common
}

type retrievePlaceArgs struct {
	PlaceID string `path:"place_id"`
}

func (m *RetrievePlace) Execute(ctx context.Context, args *retrievePlaceArgs) *web.Response {
	plcM := model.LoadPlaceByID(ctx, m.DB, args.PlaceID)
	if plcM == nil {
		return web.NewEmptyResponse(http.StatusNotFound)
	}

	plcR := render.Place(plcM)

	return web.NewJSONResponse(http.StatusOK, plcR)
}

func (m *RetrievePlace) Validate(args *retrievePlaceArgs) *web.Response {
	if !model.IsValidPlaceID(args.PlaceID) {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: resource.ErrMsgInvalidPlaceID,
		})
	}

	return nil
}

func (m *RetrievePlace) NewArgs() *retrievePlaceArgs {
	return &retrievePlaceArgs{}
}
