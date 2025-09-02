package handler

import (
	"blogs/model"
	"blogs/service"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogHandler struct {
	Service *service.BlogService
}

func createImagesDir() error {
	imagesDir := "./uploads/images"
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return os.MkdirAll(imagesDir, 0755)
	}
	return nil
}

func (h *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.Service.GetAllBlogs()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch blogs"})
		return
	}

	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	log.Printf("Blog with id %s", idStr)

	objID, err := primitive.ObjectIDFromHex(idStr)

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	blog, err := h.Service.GetBlog(objID)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Blog not found"})
		return
	}

	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	if err := createImagesDir(); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Parse multipart form (max 50MB za više slika)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	var blog model.Blog

	blog.Title = r.FormValue("title")
	blog.Description = r.FormValue("description")

	creatorIDStr := r.FormValue("creatorID")
	if creatorIDStr == "" {
		http.Error(w, "Creator ID is required", http.StatusBadRequest)
		return
	}

	creatorID, err := strconv.Atoi(creatorIDStr)
	if err != nil {
		http.Error(w, "Invalid creator ID format", http.StatusBadRequest)
		return
	}
	blog.CreatorID = creatorID

	// Validacija osnovnih podataka
	if blog.Title == "" || blog.Description == "" {
		http.Error(w, "Title and description are required", http.StatusBadRequest)
		return
	}

	// Generiranje ID-a za blog (potrebno za naziv slika)
	blog.ID = primitive.NewObjectID()

	// Procesiranje slika ako postoje
	var imageURLs []string

	// Dobijanje svih fajlova sa ključem "images"
	if files := r.MultipartForm.File["images"]; len(files) > 0 {
		log.Printf("Processing %d images", len(files))

		for i, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				log.Printf("Error opening file %d: %v", i, err)
				continue
			}

			allowedTypes := map[string]bool{
				"image/jpeg": true,
				"image/jpg":  true,
				"image/png":  true,
				"image/gif":  true,
				"image/webp": true,
			}

			contentType := fileHeader.Header.Get("Content-Type")
			if !allowedTypes[contentType] {
				file.Close()
				log.Printf("Invalid file type: %s", contentType)
				continue
			}

			// Generiranje jedinstvenog imena fajla
			timestamp := time.Now().Unix()
			fileExt := filepath.Ext(fileHeader.Filename)
			fileName := fmt.Sprintf("%s_%d_%d%s", blog.ID.Hex(), timestamp, i, fileExt)

			// Kreiranje putanje do fajla
			filePath := filepath.Join("./uploads/images", fileName)

			// Kreiranje fajla na serveru
			dst, err := os.Create(filePath)
			if err != nil {
				file.Close()
				log.Printf("Unable to create file %s: %v", fileName, err)
				continue
			}

			// Kopiranje sadržaja uploaded fajla u novi fajl
			if _, err := io.Copy(dst, file); err != nil {
				file.Close()
				dst.Close()
				os.Remove(filePath) // Obriši neuspešno kreiran fajl
				log.Printf("Unable to save file %s: %v", fileName, err)
				continue
			}

			file.Close()
			dst.Close()

			// Dodavanje URL-a slike u niz
			imageURL := fmt.Sprintf("/uploads/images/%s", fileName)
			imageURLs = append(imageURLs, imageURL)
		}
	}

	blog.Images = imageURLs

	if err := h.Service.CreateBlog(&blog); err != nil {
		// Ako ne možemo da kreiramo blog, obriši sve uploadovane slike
		for _, imageURL := range imageURLs {
			filename := filepath.Base(imageURL)
			filePath := filepath.Join("./uploads/images", filename)
			os.Remove(filePath)
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create blog"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

// Serviranje statičkih fajlova (slika)
func (h *BlogHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	// Validacija filename-a (sprečavanje directory traversal)
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("./uploads/images", filename)

	// Proverava da li fajl postoji
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// Serviranje fajla
	http.ServeFile(w, r, filePath)
}

func (h *BlogHandler) GetBlogsByCreator(w http.ResponseWriter, r *http.Request) {
	creatorIDStr := mux.Vars(r)["creatorId"]
	creatorID, err := strconv.Atoi(creatorIDStr)

	if err != nil {
		http.Error(w, "Invalid creator ID format", http.StatusBadRequest)
		return
	}

	blogs, err := h.Service.GetBlogsByCreator(creatorID)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch blogs"})
		return
	}

	json.NewEncoder(w).Encode(blogs)
}
func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Proverava da li blog postoji i da li korisnik ima pravo da ga menja
	existingBlog, err := h.Service.GetBlog(objID)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	var updateData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatorID   int    `json:"creatorID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Proverava da li je korisnik vlasnik bloga
	if existingBlog.CreatorID != updateData.CreatorID {
		http.Error(w, "You can only edit your own blogs", http.StatusForbidden)
		return
	}

	// Ažuriranje bloga
	updatedBlog := &model.Blog{
		ID:          existingBlog.ID,
		Title:       updateData.Title,
		Description: updateData.Description,
		Images:      existingBlog.Images, // Zadržava postojeće slike
		CreatorID:   existingBlog.CreatorID,
		CreatedAt:   existingBlog.CreatedAt,
	}

	if err := h.Service.UpdateBlog(objID, updatedBlog); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update blog"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBlog)
}

// Dodavanje slike u postojeći blog
func (h *BlogHandler) AddImageToBlog(w http.ResponseWriter, r *http.Request) {
	if err := createImagesDir(); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	blogIDStr := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		http.Error(w, "Invalid blog ID format", http.StatusBadRequest)
		return
	}

	// Proverava da li blog postoji
	existingBlog, err := h.Service.GetBlog(objID)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Proverava vlasništvo bloga
	creatorIDStr := r.FormValue("creatorID")
	creatorID, err := strconv.Atoi(creatorIDStr)
	if err != nil || existingBlog.CreatorID != creatorID {
		http.Error(w, "You can only edit your own blogs", http.StatusForbidden)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validacija tipa fajla
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		http.Error(w, "Only image files (JPEG, PNG, GIF, WEBP) are allowed", http.StatusBadRequest)
		return
	}

	// Generiranje jedinstvenog imena fajla
	timestamp := time.Now().Unix()
	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s_%d%s", blogIDStr, timestamp, fileExt)

	filePath := filepath.Join("./uploads/images", fileName)

	// Kreiranje fajla na serveru
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("/uploads/images/%s", fileName)

	// Dodavanje slike u blog
	if err := h.Service.AddImageToBlog(objID, imageURL); err != nil {
		os.Remove(filePath)
		http.Error(w, "Failed to update blog with image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"imageUrl": imageURL,
		"message":  "Image added successfully",
	})
}

// Uklanjanje slike iz bloga
func (h *BlogHandler) RemoveImageFromBlog(w http.ResponseWriter, r *http.Request) {
	blogIDStr := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		http.Error(w, "Invalid blog ID format", http.StatusBadRequest)
		return
	}

	var requestData struct {
		ImageURL  string `json:"imageUrl"`
		CreatorID int    `json:"creatorID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Proverava da li blog postoji i vlasništvo
	existingBlog, err := h.Service.GetBlog(objID)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	if existingBlog.CreatorID != requestData.CreatorID {
		http.Error(w, "You can only edit your own blogs", http.StatusForbidden)
		return
	}

	// Uklanjanje slike iz baze
	if err := h.Service.RemoveImageFromBlog(objID, requestData.ImageURL); err != nil {
		http.Error(w, "Failed to remove image from blog", http.StatusInternalServerError)
		return
	}

	// Brisanje fajla sa file sistema
	filename := filepath.Base(requestData.ImageURL)
	filePath := filepath.Join("./uploads/images", filename)

	if err := os.Remove(filePath); err != nil {
		log.Printf("Warning: Failed to delete image file %s: %v", filePath, err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Image removed successfully",
	})
}
