package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User defines a user struct
type User struct {
	// swagger:ignore
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// example: vsantos
	Login string `json:"login" bson:"login"`
	// example: Victor
	Firstname string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	// example: Santos
	Lastname string `json:"lastname,omitempty" bson:"lastname,omitempty"`
	// example: vsantos.py@gmail.com
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	// example: myplaintextpassword
	Password string `json:"password,omitempty" bson:"salted_password,omitempty"`
	// swagger:ignore
	CreatedAt primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
