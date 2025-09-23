package handler

import (
	"blogs/model"
	"blogs/service"
	"context"
	"local/common/proto/blogpb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogGrpcHandler struct {
	Service *service.BlogService
	blogpb.UnimplementedBlogServiceServer
}

// CreateBlog implements BlogService.CreateBlog
func (h *BlogGrpcHandler) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.BlogResponse, error) {
	blog := &model.Blog{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   int(req.CreatorId),
		Username:    req.Username,
		Images:      req.Images,
	}

	if err := h.Service.CreateBlog(blog); err != nil {
		return nil, err
	}

	return &blogpb.BlogResponse{
		Id:          blog.ID.Hex(),
		Title:       blog.Title,
		Description: blog.Description,
		Username:    blog.Username,
		CreatorID:   int32(blog.CreatorID),
		Images:      blog.Images,
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

	return &blogpb.BlogResponse{
		Id:          blog.ID.Hex(),
		Title:       blog.Title,
		Description: blog.Description,
		Username:    blog.Username,
		CreatorID:   int32(blog.CreatorID),
		Images:      blog.Images,
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
		respBlogs = append(respBlogs, &blogpb.BlogResponse{
			Id:          b.ID.Hex(),
			Title:       b.Title,
			Description: b.Description,
			Username:    b.Username,
			CreatorID:   int32(b.CreatorID),
			Images:      b.Images,
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

	return &blogpb.BlogResponse{
		Id:          existingBlog.ID.Hex(),
		Title:       existingBlog.Title,
		Description: existingBlog.Description,
		Username:    existingBlog.Username,
		CreatorID:   int32(existingBlog.CreatorID),
		Images:      existingBlog.Images,
	}, nil
}
