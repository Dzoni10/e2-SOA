package handler

import (
	"encoding/json"
	"net/http"
	"tours/model"
	"tours/service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewHandler struct {
	Service *service.ReviewService
}

func (h *ReviewHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	var review model.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.AddReview(&review); err != nil {
		http.Error(w, "Failed to add review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) GetReviewsForTour(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	reviews, err := h.Service.GetReviewsForTour(objID)
	if err != nil {
		http.Error(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reviews)
}
