package handler

import (
	"blogs/model"
	"blogs/service"
	"context"
	"fmt"
	"local/common/proto/blogpb"
	"os"
	"path"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BlogGrpcHandler struct {
	Service *service.BlogService
	blogpb.UnimplementedBlogServiceServer
}

// CreateBlog implements BlogService.CreateBlog
func (h *BlogGrpcHandler) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.BlogResponse, error) {
	now := time.Now()

	// Pretvori bytes u fajlove i sačuvaj URL-ove
	var imageURLs []string
	for i, imgBytes := range req.Images {
		if len(imgBytes) == 0 {
			continue
		}

		// Napravi ime fajla
		fileName := fmt.Sprintf("%s_%d.png", primitive.NewObjectID().Hex(), i)
		filePath := path.Join("uploads", fileName)

		// Sačuvaj fajl na disk
		if err := os.WriteFile(filePath, imgBytes, 0644); err != nil {
			return nil, fmt.Errorf("failed to save image: %w", err)
		}

		// Kreiraj URL (npr. dostupno preko static servera)
		imageURLs = append(imageURLs, "/uploads/"+fileName)
	}

	// Napravi blog model
	blog := &model.Blog{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   int(req.CreatorId),
		Username:    req.Username,
		Images:      imageURLs, // ovde ide URL, ne bytes
		CreatedAt:   now,
		Likes:       []model.Like{},
	}

	if err := h.Service.CreateBlog(blog); err != nil {
		return nil, err
	}

	// Convert likes to protobuf
	var pbLikes []*blogpb.Like
	for _, like := range blog.Likes {
		pbLikes = append(pbLikes, &blogpb.Like{
			UserId:    int32(like.UserID),
			CreatedAt: timestamppb.New(like.CreatedAt),
		})
	}

	return &blogpb.BlogResponse{
		Id:          blog.ID.Hex(),
		Title:       blog.Title,
		Description: blog.Description,
		CreatorID:   int32(blog.CreatorID),
		Username:    blog.Username,
		Images:      blog.Images, // ovo su URL-ovi
		Likes:       pbLikes,
		CreatedAt:   timestamppb.New(blog.CreatedAt),
	}, nil
}

// GetBlog implements BlogService.GetBlog
func (h *BlogGrpcHandler) GetBlog(ctx context.Context, req *blogpb.GetBlogRequest) (*blogpb.BlogResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	blog, err := h.Service.GetBlog(objID)
	if err != nil {
		return nil, err
	}

	// Convert likes to protobuf format
	var pbLikes []*blogpb.Like
	for _, like := range blog.Likes {
		pbLikes = append(pbLikes, &blogpb.Like{
			UserId:    int32(like.UserID),
			CreatedAt: timestamppb.New(like.CreatedAt),
		})
	}

	return &blogpb.BlogResponse{
		Id:          blog.ID.Hex(),
		Title:       blog.Title,
		Description: blog.Description,
		Username:    blog.Username,
		CreatorID:   int32(blog.CreatorID),
		Images:      blog.Images,
		CreatedAt:   timestamppb.New(blog.CreatedAt),
		Likes:       pbLikes,
	}, nil
}

// GetAllBlogs implements BlogService.GetAllBlogs
func (h *BlogGrpcHandler) GetAllBlogs(ctx context.Context, req *blogpb.Empty) (*blogpb.BlogListResponse, error) {
	blogs, err := h.Service.GetAllBlogs()
	if err != nil {
		return nil, err
	}

	var respBlogs []*blogpb.BlogResponse
	for _, b := range blogs {
		// Convert likes to protobuf format
		var pbLikes []*blogpb.Like
		for _, like := range b.Likes {
			pbLikes = append(pbLikes, &blogpb.Like{
				UserId:    int32(like.UserID),
				CreatedAt: timestamppb.New(like.CreatedAt),
			})
		}

		respBlogs = append(respBlogs, &blogpb.BlogResponse{
			Id:          b.ID.Hex(),
			Title:       b.Title,
			Description: b.Description,
			Username:    b.Username,
			CreatorID:   int32(b.CreatorID),
			Images:      b.Images,
			CreatedAt:   timestamppb.New(b.CreatedAt),
			Likes:       pbLikes,
		})
	}

	return &blogpb.BlogListResponse{
		Blogs: respBlogs,
	}, nil
}

// UpdateBlog implements BlogService.UpdateBlog
func (h *BlogGrpcHandler) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.BlogResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	existingBlog, err := h.Service.GetBlog(objID)
	if err != nil {
		return nil, err
	}

	existingBlog.Title = req.Title
	existingBlog.Description = req.Description
	// Optional: verify creatorID matches
	existingBlog.CreatorID = int(req.CreatorID)

	if err := h.Service.UpdateBlog(objID, existingBlog); err != nil {
		return nil, err
	}

	// Convert likes to protobuf format
	var pbLikes []*blogpb.Like
	for _, like := range existingBlog.Likes {
		pbLikes = append(pbLikes, &blogpb.Like{
			UserId:    int32(like.UserID),
			CreatedAt: timestamppb.New(like.CreatedAt),
		})
	}

	return &blogpb.BlogResponse{
		Id:          existingBlog.ID.Hex(),
		Title:       existingBlog.Title,
		Description: existingBlog.Description,
		Username:    existingBlog.Username,
		CreatorID:   int32(existingBlog.CreatorID),
		Images:      existingBlog.Images,
		CreatedAt:   timestamppb.New(existingBlog.CreatedAt),
		Likes:       pbLikes,
	}, nil
}
