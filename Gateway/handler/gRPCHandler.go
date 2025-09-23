package handler

import (
	"context"
	"encoding/json"
	"net/http"
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

	var req pb.CreateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := client.CreateBlog(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
