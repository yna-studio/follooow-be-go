package repositories

import (
	"context"
	"follooow-be/configs"
	"follooow-be/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
}

// function to create new gallery
// auto update updated_on on related influncers
func CreateGallery(ctx context.Context, params CreateGalleryParams) (*mongo.InsertOneResult, error) {
	// get now times
	now := time.Now().UnixNano() / int64(time.Millisecond)
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
		{"author_id", params.AuthorID},
	}

	// insert data to database
	result, err := GalleryCollections.InsertOne(ctx, newData)
	if err != nil {
		// stop process if error
		return result, err
	} else {
		// update influencers updated_on
		err = InfluencersUpdateOnToNow(ctx, params.Influencers)
		return result, err
	}

}
