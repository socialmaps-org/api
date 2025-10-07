package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/resource"
)

type UpdateReview struct {
	Handler
}

type UpdateReviewRequest struct {
	Liked   bool   `json:"liked"`
	Comment string `json:"comment"`
}

// 1 hour
const MAX_DELAY_IN_SECS = 1 * 60 * 60

func (h *UpdateReview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reviewID := r.PathValue("review_id")

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		panic("wrong content type")
	}

	var req UpdateReviewRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	rvwM := model.LoadReview(r.Context(), h.DB, reviewID)
	if rvwM == nil {
		panic("review not found")
	}

	// TODO
	if rvwM.UserID != "usr_foo" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not your review"))
		return
	}

	if time.Now().Unix()-rvwM.Created > MAX_DELAY_IN_SECS {
		w.Write([]byte("too late to update your review"))
		return
	}

	rvwM = model.UpdateReview(r.Context(), h.DB, reviewID, req.Liked, req.Comment)

	rvwR := resource.Review{
		ID:      rvwM.ID,
		Created: rvwM.Created,
		Place: resource.Place{
			ID: rvwM.PlaceID,
		},
		User: resource.User{
			ID: rvwM.UserID,
		},
		Liked:   rvwM.Liked,
		Comment: rvwM.Comment,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rvwR)
	if err != nil {
		panic(err)
	}
}
