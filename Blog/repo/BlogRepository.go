package repo

import (
	"blogs/database"
	"blogs/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogRepository struct{}

func (r *BlogRepository) FindAll() ([]model.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.BlogCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blogs []model.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *BlogRepository) FindById(id primitive.ObjectID) (*model.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var blog model.Blog
	err := database.BlogCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&blog)

	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) Create(blog *model.Blog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set creation time
	blog.CreatedAt = time.Now()

	_, err := database.BlogCollection.InsertOne(ctx, blog)
	return err
}

func (r *BlogRepository) FindByCreatorID(creatorID int) ([]model.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.BlogCollection.Find(ctx, bson.M{"creatorID": creatorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blogs []model.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}
	return blogs, nil
}
func (r *BlogRepository) UpdateBlog(blogID primitive.ObjectID, updatedBlog *model.Blog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":       updatedBlog.Title,
			"description": updatedBlog.Description,
			"images":      updatedBlog.Images,
		},
	}

	_, err := database.BlogCollection.UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		update,
	)

	return err
}

func (r *BlogRepository) AddImageToBlog(blogID primitive.ObjectID, imageURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$push": bson.M{
			"images": imageURL,
		},
	}

	_, err := database.BlogCollection.UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		update,
	)

	return err
}

func (r *BlogRepository) RemoveImageFromBlog(blogID primitive.ObjectID, imageURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{
			"images": imageURL,
		},
	}

	_, err := database.BlogCollection.UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		update,
	)

	return err
}
