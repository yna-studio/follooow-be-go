package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type LabelModel struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
    Label string             `bson:"label" json:"label"`
    Total int64              `bson:"total" json:"total"`
}
