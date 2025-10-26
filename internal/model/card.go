package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Card defines a user credit card
// swagger:model
type Card struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID    primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	Alias      string             `json:"alias" bson:"alias"`
	Type       string             `json:"type" bson:"type"`
	Network    string             `json:"network" bson:"network"`
	Bank       string             `json:"bank" bson:"bank"`
	Color      string             `json:"color,omitempty" bson:"color"`
	LastDigits int32              `json:"last_digits" bson:"last_digits"`
	CreatedAt  primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty" swaggerignore:"true"`
}
