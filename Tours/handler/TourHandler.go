package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"tours/metrics"
	"tours/model"
	"tours/service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type TourHandler struct {
	Service *service.TourService
}

type StatusChangeRequest struct {
	Status    model.Status `json:"status"`
	CreatorID int          `json:"creatorId"`
}

func (h *TourHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	tourIDHex := mux.Vars(r)["id"]
	tourID, err := primitive.ObjectIDFromHex(tourIDHex)
	if err != nil {
		http.Error(w, "Invalid tour ID format", http.StatusBadRequest)
		return
	}

	var request StatusChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingTour, err := h.Service.GetTour(tourID)
	if err != nil {
		http.Error(w, "Tour not found", http.StatusNotFound)
		return
	}

	if existingTour.CreatorID != request.CreatorID {
		http.Error(w, "Only tour author can change status", http.StatusForbidden)
		return
	}

	if request.Status == model.Published {
		if existingTour.Name == "" || existingTour.Description == "" {
			http.Error(w, "Tour must have name and description", http.StatusBadRequest)
			return
		}

		keypoints, err := h.Service.KeyPointRepo.FindByTourID(tourID)
		if err != nil || len(keypoints) < 2 {
			http.Error(w, "Tour must have at least 2 keypoints to be published", http.StatusBadRequest)
			return
		}
	}

	if err := h.Service.UpdateTourStatus(tourID, request.Status); err != nil {
		log.Println("Update status failed:", err)
		http.Error(w, "Failed to update tour status", http.StatusInternalServerError)
		return
	}

	// Increment metrics based on status
	metrics.ToursCreated.WithLabelValues(request.Status.String()).Inc()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "status updated",
		"newStatus": request.Status.String(),
	})
}

func (h *TourHandler) GetTourStatusInfo(w http.ResponseWriter, r *http.Request) {
	tourIDHex := mux.Vars(r)["id"]
	creatorIDStr := r.URL.Query().Get("creatorId")

	tourID, err := primitive.ObjectIDFromHex(tourIDHex)
	if err != nil {
		http.Error(w, "Invalid tour ID format", http.StatusBadRequest)
		return
	}

	tour, err := h.Service.GetTour(tourID)
	if err != nil {
		http.Error(w, "Tour not found", http.StatusNotFound)
		return
	}

	creatorID, _ := strconv.Atoi(creatorIDStr)
	canEdit := tour.CreatorID == creatorID

	response := map[string]interface{}{
		"currentStatus": tour.Status.String(),
		"canEdit":       canEdit,
		"publishedAt":   tour.PublishedAt,
		"archivedAt":    tour.ArchivedAt,
		"createdAt":     tour.CreatedAt,
		"updatedAt":     tour.UpdatedAt,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *TourHandler) GetAllTours(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	tours, err := h.Service.GetAllTours()

	//w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch tours"})
	} else {
		json.NewEncoder(w).Encode(tours)
	}

	duration := time.Since(start).Seconds()
	metrics.RequestDuration.WithLabelValues("/tours/all").Observe(duration)
	metrics.RequestCount.WithLabelValues("/tours/all", r.Method, fmt.Sprint(status)).Inc()
}

func (h *TourHandler) GetTour(w http.ResponseWriter, r *http.Request) {

	tr := otel.Tracer("tour-service")

	_, span := tr.Start(r.Context(), "GetTourHandler")
	defer span.End()

	idStr := mux.Vars(r)["id"]
	log.Printf("User with id %s", idStr)

	objID, err := primitive.ObjectIDFromHex(idStr)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid id")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tour, err := h.Service.GetTour(objID)

	//w.Header().Set("Content-Type", "application/json")

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "not found")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tour not found"})
		return
	}

	span.AddEvent("Tour fetched from DB")
	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tour)
}

/*
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
*/
func (h *TourHandler) CreateTour(w http.ResponseWriter, r *http.Request) {

	tr := otel.Tracer("tour-service")
	_, span := tr.Start(r.Context(), "CreateTourHandler")
	defer span.End()

	// Optional: log request body
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Request body:", string(bodyBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var tour model.Tour

	if err := json.NewDecoder(r.Body).Decode(&tour); err != nil {
		log.Println("Error decoding tour: ", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// ✅ Ensure tour has a proper ObjectID
	if tour.ID.IsZero() {
		tour.ID = primitive.NewObjectID()
	}
	now := time.Now()
	tour.CreatedAt = now
	tour.UpdatedAt = now

	if err := h.Service.CreateTour(&tour); err != nil {
		log.Println("Failed to save tour:", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to save")
		http.Error(w, "Failed to save tour", http.StatusInternalServerError)
		return
	}

	// Increment tour creation metric
	metrics.ToursCreated.WithLabelValues("draft").Inc()
	span.AddEvent("Tour created successfully")

	// Return the real MongoDB ID only
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id": tour.ID.Hex(),
	})
}

func (h *TourHandler) UpdateLength(w http.ResponseWriter, r *http.Request) {
	tourIDHex := mux.Vars(r)["id"]
	tourID, err := primitive.ObjectIDFromHex(tourIDHex)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	log.Println("Updating length for tour", tourID.Hex())

	if err := h.Service.UpdateTourLength(tourID); err != nil {
		log.Println("Update length failed:", err)
		http.Error(w, "Failed to update tour length", http.StatusInternalServerError)
		return
	}

	// DODAJ - vrati ažurirane metrics
	updatedTour, err := h.Service.GetTour(tourID)
	if err != nil {
		http.Error(w, "Failed to get updated tour", http.StatusInternalServerError)
		return
	}

	log.Println("Length updated successfully")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "length updated",
		"tourLength":  updatedTour.TourLength,
		"walkingTime": updatedTour.WalkingTime,
		"bicycleTime": updatedTour.BicycleTime,
		"carTime":     updatedTour.CarTime,
	})
}
func (h *TourHandler) CalculateMetrics(w http.ResponseWriter, r *http.Request) {
	var keypoints []model.KeyPoint

	if err := json.NewDecoder(r.Body).Decode(&keypoints); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	length, walkingTime, bicycleTime, carTime := h.Service.CalculateTourMetrics(keypoints)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tourLength":  length,
		"walkingTime": walkingTime,
		"bicycleTime": bicycleTime,
		"carTime":     carTime,
	})
}

func (h *TourHandler) UpdateCost(w http.ResponseWriter, r *http.Request) {
	tourIDHex := mux.Vars(r)["id"]
	tourID, err := primitive.ObjectIDFromHex(tourIDHex)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Cost float64 `json:"cost"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateTourCost(tourID, request.Cost); err != nil {
		http.Error(w, "Failed to update cost", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "cost updated"})
}
