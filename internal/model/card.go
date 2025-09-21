package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Card defines a user credit card
// swagger:model
type Card struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID    primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty" example:"5f4e76699c362be701856be6"`
	Alias      string             `json:"alias" bson:"alias" example:"My Platinum"`
	Type       string             `json:"type" bson:"type" example:"debit"`
	Network    string             `json:"network" bson:"network" example:"VISA"`
	Bank       string             `json:"bank" bson:"bank" example:"Wells Fargo"`
	Color      string             `json:"color" bson:"color" example:"#ffffff"`
	LastDigits int32              `json:"last_digits" bson:"last_digits" example:"7041"`
	CreatedAt  primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty" swaggerignore:"true"`
}
