package handler

import (
	"encoding/json"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/render"
)

type RetrievePlace struct {
	Handler
}

func (h *RetrievePlace) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	placeID := r.PathValue("place_id")

	plcM := model.LoadPlaceByID(r.Context(), h.DB, placeID)
	if plcM == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	plcR := render.Place(plcM)
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(plcR)
	if err != nil {
		panic(err)
	}
}
