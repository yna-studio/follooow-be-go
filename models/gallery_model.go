package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageModel struct {
	IsCover   bool   `json:"is_cover, omitempty" validate:"required"`
	Url       string `json:"url, omitempty" validate:"required"`
	Caption   string `json:"caption, omitempty" validate:"required"`
	CreatedOn int    `json:"created_on, omitempty" validate:"required"`
	UpdatedOn int    `json:"updated_on, omitempty" validate:"required"`
}

type GalleryModel struct {
	Id              primitive.ObjectID         `json:"id, omitempty" bson:"_id, omitempty" validate:"required"`
	Title           string                     `json:"title, omitempty" validate:"required"`
	Description     string                     `json:"description, omitempty" validate:"required"`
	Images          []ImageModel               `json:"images, omitempty" validate:"required"`
	CreatedOn       int                        `json:"created_on, omitempty"  bson:"created_on, omitempty"`
	UpdatedOn       int                        `json:"updated_on, omitempty"   bson:"updated_on, omitempty"`
	Influencers     []string                   `json:"influencers,omitempty"  validate:"required"`
	InfluencersData []InfluencerSmallDataModel `json:"influencers_data,omitempty"  validate:"required"`
	Lang            string                     `json:"lang, omitempty"  validate:"required"`
	Views           int                        `json:"views, omitempty"  validate:"required"`
	Slug            string                     `json:"slug, omitempty"  validate:"required"`
	Tags            []string                   `json:"tags,omitempty" bson:"tags,omitempty"`
	AuthorID        string                     `json:"author_id,omitempty" bson:"author_id,omitempty"`
	Author          *AuthorModel               `json:"author,omitempty" bson:"-"`
}

type PayloadGallery struct {
	Title       string       `json:"title, omitempty"`
	Description string       `json:"description, omitempty"`
	Images      []ImageModel `json:"images, omitempty"`
	Influencers []string     `json:"influencers, omitempty"`
	Lang        string       `json:"lang,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
}
