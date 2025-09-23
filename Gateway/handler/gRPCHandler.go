package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	pb "local/common/proto/blogpb"

	"google.golang.org/grpc"
)

type GrpcHandler struct{}

func (h *GrpcHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("blog-backend:50051", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "cannot connect to blog service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewBlogServiceClient(conn)

	// Maksimalna veličina forme (npr. 20MB)
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	creatorIDStr := r.FormValue("creatorID")
	username := r.FormValue("username")

	if title == "" || description == "" || creatorIDStr == "" || username == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	creatorID, err := strconv.Atoi(creatorIDStr)
	if err != nil {
		http.Error(w, "invalid creatorID", http.StatusBadRequest)
		return
	}

	// Parsiranje više fajlova iz polja "images"
	var imagesBytes [][]byte
	files := r.MultipartForm.File["images"]
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			http.Error(w, "cannot open uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "cannot read uploaded file", http.StatusBadRequest)
			return
		}

		imagesBytes = append(imagesBytes, data)
	}

	// Pravimo gRPC request
	req := &pb.CreateBlogRequest{
		Title:       title,
		Description: description,
		CreatorId:   int32(creatorID),
		Username:    username,
		Images:      imagesBytes, // sada je slice of bytes
	}

	resp, err := client.CreateBlog(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *GrpcHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/blogs/")

	conn, err := grpc.Dial("blog-backend:50051", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "cannot connect to blog service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewBlogServiceClient(conn)

	resp, err := client.GetBlog(context.Background(), &pb.GetBlogRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *GrpcHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("blog-backend:50051", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "cannot connect to blog service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewBlogServiceClient(conn)

	resp, err := client.GetAllBlogs(context.Background(), &pb.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *GrpcHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/blogs/")

	conn, err := grpc.Dial("blog-backend:50051", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "cannot connect to blog service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewBlogServiceClient(conn)

	var req pb.UpdateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Id = id

	resp, err := client.UpdateBlog(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}
