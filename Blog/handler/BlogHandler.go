package handler

import (
	"blogs/model"
	"blogs/service"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogHandler struct {
	Service *service.BlogService
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
	// Debug: print request body
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Request body:", string(bodyBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var blog model.Blog

	err := json.NewDecoder(r.Body).Decode(&blog)
	if err != nil {
		log.Println("Error decoding blog: ", err)
		http.Error(w, "Invalid request body: ", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateBlog(&blog); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create blog"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
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
