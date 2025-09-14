package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"strconv"
	"tours/model"
	"tours/service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type PositionHandler struct {
	Service *service.PositionService
}

func (h *PositionHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var userID int
	fmt.Sscanf(id, "%d", &userID)

	var position model.Position

	// Read JSON body instead of form values
	if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//userID, err := primitive.ObjectIDFromHex(id)
	//if &userID != nil {
	//http.Error(w, "Invalid user ID", http.StatusBadRequest)
	//return
	//}
	//lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	//lon, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)

	update := bson.M{
		"$set": bson.M{
			"latitude":  position.Latitude,
			"longitude": position.Longitude,
		},
	}

	if err := h.Service.UpdatePosition(userID, update); err != nil {
		http.Error(w, "Failed to update position", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *PositionHandler) InitializePosition(w http.ResponseWriter, r *http.Request) {
	var position model.Position

	// Read JSON body instead of form values
	if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.Service.InitializePosition(&position); err != nil {
		http.Error(w, "Failed to save position", http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(position)
}

func (h *PositionHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	log.Printf("User with id %s", idStr)

	var userId int
	fmt.Sscanf(idStr, "%d", &userId)

	position, err := h.Service.GetPosition(userId)

	//w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	json.NewEncoder(w).Encode(position)
}
