package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthorModel struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type NewsModel struct {
	Id              primitive.ObjectID         `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Title           string                     `json:"title,omitempty" validate:"required"`
	Thumbnail       string                     `json:"thumbnail,omitempty" validate:"required"`
	Content         string                     `json:"content,omitempty" validate:"required"`
	Views           int                        `json:"views,omitempty" validate:"required"`
	CreatedOn       int                        `json:"created_on,omitempty" bson:"created_on,omitempty" validate:"required"`
	UpdatedOn       int                        `json:"updated_on,omitempty" bson:"updated_on,omitempty" validate:"required"`
	Tags            []string                   `json:"tags,omitempty" validate:"required"`
	Influencers     []string                   `json:"influencers,omitempty"`
	InfluencersData []InfluencerSmallDataModel `json:"influencers_data,omitempty"`
	Slug            string                     `json:"slug,omitempty"`
	AuthorID        string                     `json:"author_id,omitempty" bson:"author_id,omitempty"`
	Author          *AuthorModel               `json:"author,omitempty" bson:"-"`
}

type PayloadNews struct {
	Title       string   `json:"title,omitempty"`
	Content     string   `json:"content,omitempty"`
	Thumbnail   string   `json:"thumbnail,omitempty"`
	Influencers []string `json:"influencers,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Lang        string   `json:"lang,omitempty"`
}
