package handler

import (
	"net/http"
	"strconv"

	"codeberg.org/socialmaps/auth/internal/model"
)

type LikeReview struct {
	Handler
}

func (h *LikeReview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reviewIDs := r.PathValue("review_id")
	reviewID, err := strconv.ParseUint(reviewIDs, 10, 64)
	if err != nil {
		panic(err)
	}

	userID := uint64(1)

	model.LikeReview(h.DB, userID, reviewID)

	w.WriteHeader(http.StatusNoContent)
}
