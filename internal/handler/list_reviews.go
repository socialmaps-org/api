package handler

import (
	"encoding/json"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/render"
	"codeberg.org/socialmaps/auth/internal/resource"
)

type ListReviews struct {
	Handler
}

func (h *ListReviews) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	placeID := r.PathValue("place_id")

	reviewMs := model.ListLatestReviewsOfPlace(r.Context(), h.DB, placeID)

	var reviewRs []*resource.Review
	for _, rvwM := range reviewMs {
		rvwR := render.Review(rvwM)
		reviewRs = append(reviewRs, rvwR)
	}

	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(reviewRs)
	if err != nil {
		panic(err)
	}
}
