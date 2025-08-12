package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"tours/model"
	"tours/service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourHandler struct {
	Service *service.TourService
}

func (h *TourHandler) GetAllTours(w http.ResponseWriter, r *http.Request) {
	tours, err := h.Service.GetAllTours()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch tours"})
		return
	}

	json.NewEncoder(w).Encode(tours)
}

func (h *TourHandler) GetTour(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	log.Printf("User with id %s", idStr)

	objID, err := primitive.ObjectIDFromHex(idStr)

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tour, err := h.Service.GetTour(objID)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tour not found"})
		return
	}

	json.NewEncoder(w).Encode(tour)

}

func (h *TourHandler) CreateTour(w http.ResponseWriter, r *http.Request) {

	///ISpis u konzoli vrednosti unosa fronta
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Request body:", string(bodyBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	/////
	var tour model.Tour

	err := json.NewDecoder(r.Body).Decode(&tour)
	if err != nil {
		log.Println("Error decoding tour: ", err)
		http.Error(w, "Invalid request body: ", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateTour(&tour); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tour)
}
