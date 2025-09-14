package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"tours/model"
	"tours/service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KeyPointHandler struct {
	Service *service.KeyPointService
}

// Response structure for keypoint creation
type KeyPointResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Order       int      `json:"order"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Images      []string `json:"images"`
	TourID      string   `json:"tourId,omitempty"`
}

func (h *KeyPointHandler) CreateKeyPoint(w http.ResponseWriter, r *http.Request) {
	// Panic recovery to never return empty reply
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("❌ Panic recovered:", rec)
			//w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("%v", rec),
			})
		}
	}()

	log.Println("➡️ CreateKeyPoint called, method:", r.Method)

	// CORS headers
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse multipart form (50 MB)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		log.Println("❌ Error parsing form:", err)
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Read form values
	name := r.FormValue("name")
	latStr := r.FormValue("latitude")
	lonStr := r.FormValue("longitude")
	description := r.FormValue("description")
	orderStr := r.FormValue("order")
	tourIdStr := r.FormValue("tourId")

	if name == "" || latStr == "" || lonStr == "" {
		http.Error(w, "Missing required fields: name, latitude, longitude", http.StatusBadRequest)
		return
	}

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	var tourId primitive.ObjectID
	if tourIdStr != "" {
		tourId, err = primitive.ObjectIDFromHex(tourIdStr)
		if err != nil {
			http.Error(w, "Invalid tourId", http.StatusBadRequest)
			return
		}
	}

	// Determine order
	order := 1
	if orderStr != "" {
		order, _ = strconv.Atoi(orderStr)
	} else if !tourId.IsZero() {
		order, _ = h.Service.GetNextOrderForTour(tourId)
	}

	// Handle files if present
	imagePaths := []string{}
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		files := r.MultipartForm.File["images"]
		log.Println("Files detected:", len(files))
		if len(files) > 0 {
			os.MkdirAll("uploads/keypoints", os.ModePerm)
			for _, fh := range files {
				file, err := fh.Open()
				if err != nil {
					log.Println("❌ Cannot open file:", err)
					continue
				}
				defer file.Close()

				ext := strings.ToLower(filepath.Ext(fh.Filename))
				if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
					continue
				}

				uniqueFilename := fmt.Sprintf("%d_%s", primitive.NewObjectID().Timestamp().Unix(), fh.Filename)
				savePath := filepath.Join("uploads/keypoints", uniqueFilename)

				dst, err := os.Create(savePath)
				if err != nil {
					log.Println("❌ Cannot create file:", err)
					continue
				}
				defer dst.Close()

				if _, err := io.Copy(dst, file); err != nil {
					log.Println("❌ Error saving file:", err)
					continue
				}

				imagePaths = append(imagePaths, "/uploads/keypoints/"+uniqueFilename)
			}
		}
	}

	// Create keypoint object
	keypoint := model.KeyPoint{
		ID:          primitive.NewObjectID(),
		TourID:      tourId,
		Latitude:    latitude,
		Longitude:   longitude,
		Name:        name,
		Description: description,
		Order:       order,
		Images:      imagePaths,
	}

	if err := h.Service.Create(&keypoint); err != nil {
		log.Println("❌ Failed to save keypoint:", err)
		http.Error(w, "Failed to save keypoint: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          keypoint.ID.Hex(),
		"name":        keypoint.Name,
		"description": keypoint.Description,
		"order":       keypoint.Order,
		"latitude":    keypoint.Latitude,
		"longitude":   keypoint.Longitude,
		"images":      keypoint.Images,
		"tourId": func() string {
			if !keypoint.TourID.IsZero() {
				return keypoint.TourID.Hex()
			} else {
				return ""
			}
		}(),
	})

	log.Println("✅ Keypoint created successfully:", keypoint.ID.Hex())
}

// Update keypoint with tourId and reorder - used after tour creation
func (h *KeyPointHandler) UpdateKeyPointTourId(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	keypointID := vars["id"]

	objID, err := primitive.ObjectIDFromHex(keypointID)
	if err != nil {
		http.Error(w, "Invalid keypoint ID", http.StatusBadRequest)
		return
	}

	var requestData struct {
		TourID primitive.ObjectID `json:"tourId"`
		Order  int                `json:"order"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateTourIdAndOrder(objID, requestData.TourID, requestData.Order); err != nil {
		http.Error(w, "Failed to update keypoint", http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Keypoint updated successfully",
	})
}

// Bulk update keypoints with tourId - used after tour creation
func (h *KeyPointHandler) BulkUpdateKeyPointsTourId(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var requestData struct {
		KeyPointIds []string `json:"keypointIds"`
		TourID      string   `json:"tourId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Convert string IDs to ObjectIDs
	var keypointObjIds []primitive.ObjectID
	for _, idStr := range requestData.KeyPointIds {
		objId, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid keypoint ID: %s", idStr), http.StatusBadRequest)
			return
		}
		keypointObjIds = append(keypointObjIds, objId)
	}

	tourObjId, err := primitive.ObjectIDFromHex(requestData.TourID)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.BulkUpdateTourId(keypointObjIds, tourObjId); err != nil {
		http.Error(w, "Failed to update keypoints: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Keypoints updated successfully",
	})
}

// Add image to existing keypoint
func (h *KeyPointHandler) AddImageToKeyPoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keypointID := vars["id"]

	objID, err := primitive.ObjectIDFromHex(keypointID)
	if err != nil {
		http.Error(w, "Invalid keypoint ID", http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(10 << 20)
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		http.Error(w, "No image uploaded", http.StatusBadRequest)
		return
	}

	var imagePaths []string
	os.MkdirAll("uploads/keypoints", os.ModePerm)

	for _, fileHeader := range files {
		file, _ := fileHeader.Open()
		defer file.Close()

		// Generate unique filename
		uniqueFilename := fmt.Sprintf("%d_%s", primitive.NewObjectID().Timestamp().Unix(), fileHeader.Filename)
		savePath := filepath.Join("uploads/keypoints", uniqueFilename)

		dst, _ := os.Create(savePath)
		defer dst.Close()
		io.Copy(dst, file)

		imagePaths = append(imagePaths, "/uploads/keypoints/"+uniqueFilename)
	}

	if err := h.Service.AddImages(objID, imagePaths); err != nil {
		http.Error(w, "Failed to update keypoint", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Image(s) added successfully",
		"images":  imagePaths,
	})
}

// Remove image from keypoint
func (h *KeyPointHandler) RemoveImageFromKeyPoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keypointID := vars["id"]
	filename := r.URL.Query().Get("filename")

	objID, err := primitive.ObjectIDFromHex(keypointID)
	if err != nil {
		http.Error(w, "Invalid keypoint ID", http.StatusBadRequest)
		return
	}

	// Remove from database
	if err := h.Service.RemoveImage(objID, "/uploads/keypoints/"+filename); err != nil {
		http.Error(w, "Failed to remove image", http.StatusInternalServerError)
		return
	}

	// Remove from disk
	os.Remove(filepath.Join("uploads/keypoints", filename))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Image removed successfully",
	})
}

// Serve keypoint image
func (h *KeyPointHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	http.ServeFile(w, r, filepath.Join("uploads/keypoints", filename))
}

// Get all keypoints for a tour (ordered)
func (h *KeyPointHandler) GetKeyPointsForTour(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	keypoints, err := h.Service.GetKeyPointsForTour(objID)
	if err != nil {
		http.Error(w, "Failed to fetch keypoints", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(keypoints)
}

// Update keypoint order
func (h *KeyPointHandler) UpdateKeyPointOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keypointID := vars["id"]

	objID, err := primitive.ObjectIDFromHex(keypointID)
	if err != nil {
		http.Error(w, "Invalid keypoint ID", http.StatusBadRequest)
		return
	}

	var requestData struct {
		Order int `json:"order"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateOrder(objID, requestData.Order); err != nil {
		http.Error(w, "Failed to update keypoint order", http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Keypoint order updated successfully",
	})
}

// Bulk reorder keypoints
func (h *KeyPointHandler) BulkReorderKeyPoints(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		KeyPoints []struct {
			ID    primitive.ObjectID `json:"id"`
			Order int                `json:"order"`
		} `json:"keypoints"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.BulkReorder(requestData.KeyPoints); err != nil {
		http.Error(w, "Failed to reorder keypoints", http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Keypoints reordered successfully",
	})
}

// Delete keypoint
func (h *KeyPointHandler) DeleteKeyPoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keypointID := vars["id"]

	objID, err := primitive.ObjectIDFromHex(keypointID)
	if err != nil {
		http.Error(w, "Invalid keypoint ID", http.StatusBadRequest)
		return
	}

	// Get keypoint to delete associated images
	keypoint, err := h.Service.GetByID(objID)
	if err != nil {
		http.Error(w, "Keypoint not found", http.StatusNotFound)
		return
	}

	// Delete from database
	if err := h.Service.Delete(objID); err != nil {
		http.Error(w, "Failed to delete keypoint", http.StatusInternalServerError)
		return
	}

	// Delete associated images from disk
	for _, imagePath := range keypoint.Images {
		filename := filepath.Base(imagePath)
		os.Remove(filepath.Join("uploads/keypoints", filename))
	}

	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Keypoint deleted successfully",
	})
}

func (h *KeyPointHandler) UpdateKeypoint(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	kpID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid keypoint ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	desc := r.FormValue("description")
	lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	lon, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)

	// Handle image uploads
	var imagePaths []string
	files := r.MultipartForm.File["images"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save file locally (adjust path if needed)
		path := fmt.Sprintf("uploads/keypoints/%s", fileHeader.Filename)
		dst, err := os.Create(path)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		imagePaths = append(imagePaths, path)
	}

	update := bson.M{
		"$set": bson.M{
			"name":        name,
			"description": desc,
			"latitude":    lat,
			"longitude":   lon,
			"images":      imagePaths, // overwrite images
		},
	}

	if err := h.Service.UpdateKeypoint(kpID, update); err != nil {
		http.Error(w, "Failed to update keypoint", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}
