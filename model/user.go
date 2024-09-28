package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User defines a user struct
type User struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Login string             `json:"login,omitempty" bson:"login"`
	Name  string             `json:"name,omitempty" bson:"name"`
	Email string             `json:"email,omitempty" bson:"email"`
}
