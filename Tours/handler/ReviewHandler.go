package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tours/model"
	"tours/service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewHandler struct {
	Service *service.ReviewService
}

// Kreiranje recenzije sa slikama
func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	// Limit 50 MB
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Učitaj osnovne podatke
	tourIdStr := r.FormValue("tourId")
	userIdStr := r.FormValue("userId")
	username := r.FormValue("username")
	ratingStr := r.FormValue("rating")
	comment := r.FormValue("comment")

	tourId, err := primitive.ObjectIDFromHex(tourIdStr)
	if err != nil {
		http.Error(w, "Invalid tourId", http.StatusBadRequest)
		return
	}
	var userId int
	fmt.Sscanf(userIdStr, "%d", &userId)
	var rating int
	fmt.Sscanf(ratingStr, "%d", &rating)

	// Obrada fajlova
	files := r.MultipartForm.File["images"]
	var imagePaths []string

	os.MkdirAll("uploads/reviews", os.ModePerm)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// samo dozvoljeni tipovi
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			continue
		}

		savePath := filepath.Join("uploads/reviews", fileHeader.Filename)
		dst, err := os.Create(savePath)
		if err != nil {
			continue
		}
		defer dst.Close()

		io.Copy(dst, file)

		imagePaths = append(imagePaths, "/uploads/reviews/"+fileHeader.Filename)
	}

	review := model.Review{
		ID:        primitive.NewObjectID(),
		TourID:    tourId,
		UserID:    userId,
		Username:  username,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now(),
		VisitedAt: time.Now(),
		Images:    imagePaths,
	}

	if err := h.Service.Create(&review); err != nil {
		http.Error(w, "Failed to save review", http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// Dodavanje slike u postojeći review
func (h *ReviewHandler) AddImageToReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewID := vars["id"]

	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(10 << 20)
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		http.Error(w, "No image uploaded", http.StatusBadRequest)
		return
	}

	var imagePaths []string
	os.MkdirAll("uploads/reviews", os.ModePerm)

	for _, fileHeader := range files {
		file, _ := fileHeader.Open()
		defer file.Close()

		savePath := filepath.Join("uploads/reviews", fileHeader.Filename)
		dst, _ := os.Create(savePath)
		defer dst.Close()
		io.Copy(dst, file)

		imagePaths = append(imagePaths, "/uploads/reviews/"+fileHeader.Filename)
	}

	if err := h.Service.AddImages(objID, imagePaths); err != nil {
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Image(s) added successfully",
		"images":  imagePaths,
	})
}

// Brisanje slike iz review-a
func (h *ReviewHandler) RemoveImageFromReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewID := vars["id"]
	filename := r.URL.Query().Get("filename")

	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	// obriši iz baze
	if err := h.Service.RemoveImage(objID, "/uploads/reviews/"+filename); err != nil {
		http.Error(w, "Failed to remove image", http.StatusInternalServerError)
		return
	}

	// obriši sa diska
	os.Remove(filepath.Join("uploads/reviews", filename))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Image removed successfully",
	})
}

// Serve image
func (h *ReviewHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	http.ServeFile(w, r, filepath.Join("uploads/reviews", filename))
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
