package handler

import (
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
)

type DeleteReview struct {
	Handler
}

func (h *DeleteReview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reviewID := r.PathValue("review_id")
	model.DeleteReview(r.Context(), h.DB, reviewID)
	w.WriteHeader(http.StatusNoContent)
}
