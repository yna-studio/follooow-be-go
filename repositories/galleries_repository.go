package repositories

import (
	"context"
	"fmt"
	"follooow-be/configs"
	"follooow-be/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var GalleryCollections *mongo.Collection = configs.GetCollection(configs.DB, "galleries")

// types
type CreateGalleryParams struct {
	Title       string
	Description string
	Images      []models.ImageModel
	Influencers []string
	Lang        string
	Slug        string
	AuthorID    string
	Tags        []string
}

// function to create new gallery
// auto update updated_on on related influncers
func CreateGallery(ctx context.Context, params CreateGalleryParams) (*mongo.InsertOneResult, error) {
	// get now times
	now := time.Now().UnixNano() / int64(time.Millisecond)

	// Convert author_id string to ObjectId
	var authorObjID primitive.ObjectID
	if params.AuthorID != "" {
		var err error
		authorObjID, err = primitive.ObjectIDFromHex(params.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("invalid author_id format: %w", err)
		}
	}

	// ref: https://stackoverflow.com/a/8689281/2780875
	newData := bson.D{
		{"title", params.Title},
		{"description", params.Description},
		{"slug", params.Slug},
		{"views", 1},
		{"updated_on", now},
		{"created_on", now},
		{"lang", params.Lang},
		{"images", params.Images},
		{"influencers", params.Influencers},
		{"author_id", authorObjID},
		{"tags", params.Tags},
	}

	// insert data to database
	result, err := GalleryCollections.InsertOne(ctx, newData)
	if err != nil {
		return result, err
	}

	// update influencers updated_on as a best-effort secondary action.
	// This should not fail the gallery creation if a related influencer ID is stale/invalid.
	if len(params.Influencers) > 0 {
		updateCtx, updateCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer updateCancel()
		if updateErr := InfluencersUpdateOnToNow(updateCtx, params.Influencers); updateErr != nil {
			fmt.Printf("warning: failed to update influencer timestamps after gallery insert: %v\n", updateErr)
		}
	}

	return result, nil

}
