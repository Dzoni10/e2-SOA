package handler

import (
	"blogs/logger"
	"blogs/model"
	"blogs/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	logger.Info("Fetching all blogs", logrus.Fields{
		"action": "get_all_blogs",
	})

	blogs, err := h.Service.GetAllBlogs()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		logger.Error("Failed to fetch all blogs", err, logrus.Fields{
			"action": "get_all_blogs",
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch blogs"})
		return
	}

	logger.Info("Successfully fetched all blogs", logrus.Fields{
		"action":     "get_all_blogs",
		"blog_count": len(blogs),
	})

	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	logger.Info("Fetching single blog", logrus.Fields{
		"blog_id": idStr,
		"action":  "get_blog",
	})

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.Error("Invalid blog ID format", err, logrus.Fields{
			"blog_id": idStr,
			"action":  "get_blog",
		})
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	blog, err := h.Service.GetBlog(objID)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		logger.Error("Blog not found", err, logrus.Fields{
			"blog_id": idStr,
			"action":  "get_blog",
		})
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Blog not found"})
		return
	}

	logger.Info("Successfully fetched blog", logrus.Fields{
		"blog_id": idStr,
		"title":   blog.Title,
		"action":  "get_blog",
	})

	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	logger.Info("Starting blog creation", logrus.Fields{
		"action": "create_blog",
	})

	if err := createImagesDir(); err != nil {
		logger.Error("Failed to create upload directory", err, logrus.Fields{
			"action": "create_blog",
		})
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Parse multipart form (max 50MB za više slika)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		logger.Error("Unable to parse multipart form", err, logrus.Fields{
			"action": "create_blog",
		})
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	var blog model.Blog

	blog.Title = r.FormValue("title")
	blog.Description = r.FormValue("description")
	blog.Username = r.FormValue("username")
	creatorIDStr := r.FormValue("creatorID")

	if creatorIDStr == "" {
		logger.Error("Creator ID missing", nil, logrus.Fields{
			"action": "create_blog",
		})
		http.Error(w, "Creator ID is required", http.StatusBadRequest)
		return
	}

	creatorID, err := strconv.Atoi(creatorIDStr)
	if err != nil {
		logger.Error("Invalid creator ID format", err, logrus.Fields{
			"creator_id_str": creatorIDStr,
			"action":         "create_blog",
		})
		http.Error(w, "Invalid creator ID format", http.StatusBadRequest)
		return
	}
	blog.CreatorID = creatorID

	// Validacija osnovnih podataka
	if blog.Title == "" || blog.Description == "" || blog.Username == "" {
		logger.Error("Missing required fields", nil, logrus.Fields{
			"user_id":  creatorID,
			"username": blog.Username,
			"title":    blog.Title,
			"action":   "create_blog",
		})
		http.Error(w, "Title, description, and username are required", http.StatusBadRequest)
		return
	}

	logger.Info("User creating new blog", logrus.Fields{
		"user_id":  creatorID,
		"username": blog.Username,
		"title":    blog.Title,
		"action":   "create_blog",
	})

	// Generiranje ID-a za blog (potrebno za naziv slika)
	blog.ID = primitive.NewObjectID()

	// Procesiranje slika ako postoje
	var imageURLs []string

	// Dobijanje svih fajlova sa ključem "images"
	if files := r.MultipartForm.File["images"]; len(files) > 0 {
		logger.Info("Processing uploaded images", logrus.Fields{
			"user_id":     creatorID,
			"blog_id":     blog.ID.Hex(),
			"image_count": len(files),
			"action":      "create_blog",
		})

		for i, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				logger.Error("Failed to open uploaded file", err, logrus.Fields{
					"user_id":    creatorID,
					"blog_id":    blog.ID.Hex(),
					"file_index": i,
					"filename":   fileHeader.Filename,
					"action":     "create_blog",
				})
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
				logger.Warn("Invalid file type uploaded", logrus.Fields{
					"user_id":      creatorID,
					"blog_id":      blog.ID.Hex(),
					"filename":     fileHeader.Filename,
					"content_type": contentType,
					"action":       "create_blog",
				})
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
				logger.Error("Unable to create file", err, logrus.Fields{
					"user_id":  creatorID,
					"blog_id":  blog.ID.Hex(),
					"filename": fileName,
					"action":   "create_blog",
				})
				continue
			}

			// Kopiranje sadržaja uploaded fajla u novi fajl
			if _, err := io.Copy(dst, file); err != nil {
				file.Close()
				dst.Close()
				os.Remove(filePath) // Obriši neuspešno kreiran fajl
				logger.Error("Unable to save file", err, logrus.Fields{
					"user_id":  creatorID,
					"blog_id":  blog.ID.Hex(),
					"filename": fileName,
					"action":   "create_blog",
				})
				continue
			}

			file.Close()
			dst.Close()

			// Dodavanje URL-a slike u niz
			imageURL := fmt.Sprintf("/uploads/images/%s", fileName)
			imageURLs = append(imageURLs, imageURL)

			logger.Debug("Successfully processed image", logrus.Fields{
				"user_id":   creatorID,
				"blog_id":   blog.ID.Hex(),
				"filename":  fileName,
				"image_url": imageURL,
				"action":    "create_blog",
			})
		}

		logger.Info("Image processing completed", logrus.Fields{
			"user_id":          creatorID,
			"blog_id":          blog.ID.Hex(),
			"processed_images": len(imageURLs),
			"total_files_sent": len(files),
			"action":           "create_blog",
		})
	}

	blog.Images = imageURLs

	if err := h.Service.CreateBlog(&blog); err != nil {
		// Ako ne možemo da kreiramo blog, obriši sve uploadovane slike
		for _, imageURL := range imageURLs {
			filename := filepath.Base(imageURL)
			filePath := filepath.Join("./uploads/images", filename)
			os.Remove(filePath)
		}

		logger.Error("Failed to save blog to database", err, logrus.Fields{
			"user_id": creatorID,
			"blog_id": blog.ID.Hex(),
			"title":   blog.Title,
			"action":  "create_blog",
		})

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create blog"})
		return
	}

	logger.Info("Blog successfully created", logrus.Fields{
		"user_id":     creatorID,
		"username":    blog.Username,
		"blog_id":     blog.ID.Hex(),
		"title":       blog.Title,
		"image_count": len(imageURLs),
		"action":      "create_blog",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

// Serviranje statičkih fajlova (slika)
func (h *BlogHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	if filename == "" {
		logger.Warn("Image request with empty filename", logrus.Fields{
			"action": "serve_image",
		})
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	// Validacija filename-a (sprečavanje directory traversal)
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		logger.Warn("Suspicious filename in image request", logrus.Fields{
			"filename": filename,
			"action":   "serve_image",
		})
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("./uploads/images", filename)

	// Proverava da li fajl postoji
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Warn("Requested image not found", logrus.Fields{
			"filename": filename,
			"action":   "serve_image",
		})
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	logger.Debug("Serving image", logrus.Fields{
		"filename": filename,
		"action":   "serve_image",
	})

	// Serviranje fajla
	http.ServeFile(w, r, filePath)
}

func (h *BlogHandler) GetBlogsByCreator(w http.ResponseWriter, r *http.Request) {
	creatorIDStr := mux.Vars(r)["creatorId"]
	creatorID, err := strconv.Atoi(creatorIDStr)

	if err != nil {
		logger.Error("Invalid creator ID format", err, logrus.Fields{
			"creator_id_str": creatorIDStr,
			"action":         "get_blogs_by_creator",
		})
		http.Error(w, "Invalid creator ID format", http.StatusBadRequest)
		return
	}

	logger.Info("Fetching blogs by creator", logrus.Fields{
		"creator_id": creatorID,
		"action":     "get_blogs_by_creator",
	})

	blogs, err := h.Service.GetBlogsByCreator(creatorID)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		logger.Error("Failed to fetch blogs by creator", err, logrus.Fields{
			"creator_id": creatorID,
			"action":     "get_blogs_by_creator",
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch blogs"})
		return
	}

	logger.Info("Successfully fetched blogs by creator", logrus.Fields{
		"creator_id": creatorID,
		"blog_count": len(blogs),
		"action":     "get_blogs_by_creator",
	})

	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	logger.Info("Starting blog update", logrus.Fields{
		"blog_id": idStr,
		"action":  "update_blog",
	})

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.Error("Invalid blog ID format for update", err, logrus.Fields{
			"blog_id": idStr,
			"action":  "update_blog",
		})
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Proverava da li blog postoji i da li korisnik ima pravo da ga menja
	existingBlog, err := h.Service.GetBlog(objID)
	if err != nil {
		logger.Error("Blog not found for update", err, logrus.Fields{
			"blog_id": idStr,
			"action":  "update_blog",
		})
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	var updateData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatorID   int    `json:"creatorID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		logger.Error("Invalid request body for blog update", err, logrus.Fields{
			"blog_id": idStr,
			"action":  "update_blog",
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Proverava da li je korisnik vlasnik bloga
	if existingBlog.CreatorID != updateData.CreatorID {
		logger.Warn("Unauthorized blog update attempt", logrus.Fields{
			"blog_id":         idStr,
			"blog_creator_id": existingBlog.CreatorID,
			"request_user_id": updateData.CreatorID,
			"action":          "update_blog",
		})
		http.Error(w, "You can only edit your own blogs", http.StatusForbidden)
		return
	}

	logger.Info("User updating their blog", logrus.Fields{
		"user_id":   updateData.CreatorID,
		"blog_id":   idStr,
		"old_title": existingBlog.Title,
		"new_title": updateData.Title,
		"action":    "update_blog",
	})

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
		logger.Error("Failed to update blog", err, logrus.Fields{
			"user_id": updateData.CreatorID,
			"blog_id": idStr,
			"action":  "update_blog",
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update blog"})
		return
	}

	logger.Info("Blog successfully updated", logrus.Fields{
		"user_id": updateData.CreatorID,
		"blog_id": idStr,
		"title":   updateData.Title,
		"action":  "update_blog",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBlog)
}

// Dodavanje slike u postojeći blog
func (h *BlogHandler) AddImageToBlog(w http.ResponseWriter, r *http.Request) {
	blogIDStr := mux.Vars(r)["id"]
	logger.Info("Adding image to blog", logrus.Fields{
		"blog_id": blogIDStr,
		"action":  "add_image_to_blog",
	})

	if err := createImagesDir(); err != nil {
		logger.Error("Failed to create upload directory", err, logrus.Fields{
			"blog_id": blogIDStr,
			"action":  "add_image_to_blog",
		})
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	objID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		logger.Error("Invalid blog ID format", err, logrus.Fields{
			"blog_id": blogIDStr,
			"action":  "add_image_to_blog",
		})
		http.Error(w, "Invalid blog ID format", http.StatusBadRequest)
		return
	}

	// Proverava da li blog postoji
	existingBlog, err := h.Service.GetBlog(objID)
	if err != nil {
		logger.Error("Blog not found for image addition", err, logrus.Fields{
			"blog_id": blogIDStr,
			"action":  "add_image_to_blog",
		})
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logger.Error("Unable to parse multipart form", err, logrus.Fields{
			"blog_id": blogIDStr,
			"action":  "add_image_to_blog",
		})
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Proverava vlasništvo bloga
	creatorIDStr := r.FormValue("creatorID")
	creatorID, err := strconv.Atoi(creatorIDStr)
	if err != nil || existingBlog.CreatorID != creatorID {
		logger.Warn("Unauthorized attempt to add image to blog", logrus.Fields{
			"blog_id":         blogIDStr,
			"blog_creator_id": existingBlog.CreatorID,
			"request_user_id": creatorID,
			"action":          "add_image_to_blog",
		})
		http.Error(w, "You can only edit your own blogs", http.StatusForbidden)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		logger.Error("Unable to get file from form", err, logrus.Fields{
			"user_id": creatorID,
			"blog_id": blogIDStr,
			"action":  "add_image_to_blog",
		})
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	logger.Info("User adding image to blog", logrus.Fields{
		"user_id":  creatorID,
		"blog_id":  blogIDStr,
		"filename": fileHeader.Filename,
		"action":   "add_image_to_blog",
	})

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
		logger.Warn("Invalid file type for image upload", logrus.Fields{
			"user_id":      creatorID,
			"blog_id":      blogIDStr,
			"filename":     fileHeader.Filename,
			"content_type": contentType,
			"action":       "add_image_to_blog",
		})
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
		logger.Error("Unable to create file", err, logrus.Fields{
			"user_id":  creatorID,
			"blog_id":  blogIDStr,
			"filename": fileName,
			"action":   "add_image_to_blog",
		})
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath)
		logger.Error("Unable to save file", err, logrus.Fields{
			"user_id":  creatorID,
			"blog_id":  blogIDStr,
			"filename": fileName,
			"action":   "add_image_to_blog",
		})
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("/uploads/images/%s", fileName)

	// Dodavanje slike u blog
	if err := h.Service.AddImageToBlog(objID, imageURL); err != nil {
		os.Remove(filePath)
		logger.Error("Failed to update blog with image", err, logrus.Fields{
			"user_id":   creatorID,
			"blog_id":   blogIDStr,
			"image_url": imageURL,
			"action":    "add_image_to_blog",
		})
		http.Error(w, "Failed to update blog with image", http.StatusInternalServerError)
		return
	}

	logger.Info("Image successfully added to blog", logrus.Fields{
		"user_id":   creatorID,
		"blog_id":   blogIDStr,
		"filename":  fileName,
		"image_url": imageURL,
		"action":    "add_image_to_blog",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"imageUrl": imageURL,
		"message":  "Image added successfully",
	})
}

// Uklanjanje slike iz bloga
func (h *BlogHandler) RemoveImageFromBlog(w http.ResponseWriter, r *http.Request) {
	blogIDStr := mux.Vars(r)["id"]
	logger.Info("Removing image from blog", logrus.Fields{
		"blog_id": blogIDStr,
		"action":  "remove_image_from_blog",
	})

	objID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		logger.Error("Invalid blog ID format", err, logrus.Fields{
			"blog_id": blogIDStr,
			"action":  "remove_image_from_blog",
		})
		http.Error(w, "Invalid blog ID format", http.StatusBadRequest)
		return
	}

	var requestData struct {
		ImageURL  string `json:"imageUrl"`
		CreatorID int    `json:"creatorID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.Error("Invalid request body", err, logrus.Fields{
			"blog_id": blogIDStr,
			"action":  "remove_image_from_blog",
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Proverava da li blog postoji i vlasništvo
	existingBlog, err := h.Service.GetBlog(objID)
	if err != nil {
		logger.Error("Blog not found for image removal", err, logrus.Fields{
			"blog_id": blogIDStr,
			"user_id": requestData.CreatorID,
			"action":  "remove_image_from_blog",
		})
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	if existingBlog.CreatorID != requestData.CreatorID {
		logger.Warn("Unauthorized attempt to remove image from blog", logrus.Fields{
			"blog_id":         blogIDStr,
			"blog_creator_id": existingBlog.CreatorID,
			"request_user_id": requestData.CreatorID,
			"image_url":       requestData.ImageURL,
			"action":          "remove_image_from_blog",
		})
		http.Error(w, "You can only edit your own blogs", http.StatusForbidden)
		return
	}

	logger.Info("User removing image from blog", logrus.Fields{
		"user_id":   requestData.CreatorID,
		"blog_id":   blogIDStr,
		"image_url": requestData.ImageURL,
		"action":    "remove_image_from_blog",
	})

	// Uklanjanje slike iz baze
	if err := h.Service.RemoveImageFromBlog(objID, requestData.ImageURL); err != nil {
		logger.Error("Failed to remove image from blog in database", err, logrus.Fields{
			"user_id":   requestData.CreatorID,
			"blog_id":   blogIDStr,
			"image_url": requestData.ImageURL,
			"action":    "remove_image_from_blog",
		})
		http.Error(w, "Failed to remove image from blog", http.StatusInternalServerError)
		return
	}

	// Brisanje fajla sa file sistema
	filename := filepath.Base(requestData.ImageURL)
	filePath := filepath.Join("./uploads/images", filename)

	if err := os.Remove(filePath); err != nil {
		logger.Warn("Failed to delete image file from filesystem", logrus.Fields{
			"user_id":   requestData.CreatorID,
			"blog_id":   blogIDStr,
			"image_url": requestData.ImageURL,
			"file_path": filePath,
			"error":     err.Error(),
			"action":    "remove_image_from_blog",
		})
	}

	logger.Info("Image successfully removed from blog", logrus.Fields{
		"user_id":   requestData.CreatorID,
		"blog_id":   blogIDStr,
		"image_url": requestData.ImageURL,
		"action":    "remove_image_from_blog",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Image removed successfully",
	})
}
