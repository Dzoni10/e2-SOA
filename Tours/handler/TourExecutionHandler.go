package handler

import (
	"encoding/json"
	"net/http"
	"tours/service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourExecutionHandler struct {
	Service *service.TourExecutionService
}

func (h *TourExecutionHandler) StartTour(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TourId    string  `json:"tourId"`
		TouristId int     `json:"touristId"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	tourID, err := primitive.ObjectIDFromHex(req.TourId)
	if err != nil {
		http.Error(w, "Invalid tourId", http.StatusBadRequest)
		return
	}

	execution, err := h.Service.StartTour(tourID, req.TouristId, req.Latitude, req.Longitude)
	if err != nil {
		http.Error(w, "Failed to start tour", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(execution)
}

func (h *TourExecutionHandler) FinishTour(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := h.Service.FinishTour(id); err != nil {
		http.Error(w, "Failed to finish tour", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TourExecutionHandler) AbandonTour(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := h.Service.AbandonTour(id); err != nil {
		http.Error(w, "Failed to abandon tour", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TourExecutionHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateLocation(id, req.Latitude, req.Longitude); err != nil {
		http.Error(w, "Failed to update location", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TourExecutionHandler) AddCompletedKeyPoint(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		KeyPointId     string `json:"keyPointId"`
		TotalKeyPoints int    `json:"totalKeyPoints"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	keyPointID, err := primitive.ObjectIDFromHex(req.KeyPointId)
	if err != nil {
		http.Error(w, "Invalid keyPointId", http.StatusBadRequest)
		return
	}

	if err := h.Service.AddCompletedKeyPoint(id, keyPointID, req.TotalKeyPoints); err != nil {
		http.Error(w, "Failed to add keypoint", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TourExecutionHandler) GetExecutionById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	objID, err := primitive.ObjectIDFromHex(idStr)

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tour, err := h.Service.GetExecutionById(objID)

	//w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Execution not found"})
		return
	}

	json.NewEncoder(w).Encode(tour)

}

