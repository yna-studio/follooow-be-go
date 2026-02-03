package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Username  string              `json:"username,omitempty" validate:"required"`
	Password  string              `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
	CreatedAt int64               `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt int64               `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type CreateUserModel struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type UserResponse struct {
	ID       primitive.ObjectID `json:"id,omitempty"`
	Username string              `json:"username,omitempty"`
	CreatedAt int64              `json:"created_at,omitempty"`
	UpdatedAt int64              `json:"updated_at,omitempty"`
}
