package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Card defines a user credit card
// swagger:model
type Card struct {
	// swagger:ignore
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// example: 5f4e76699c362be701856be6
	OwnerID primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	// example: My Platinum Card
	Alias string `json:"alias" bson:"alias"`
	// example: VISA
	Network string `json:"network" bson:"network"`
	// example: #ffffff
	Color string `json:"color" bson:"color"`
	// example: 1234
	LastDigits int32 `json:"last_digits" bson:"last_digits"`
	// swagger:ignore
	CreatedAt primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
