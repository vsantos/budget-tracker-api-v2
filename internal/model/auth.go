package model

// JWTUSer defines a minimum set of user for posterior authentication
type JWTUser struct {
	// ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// example: vsantos
	Login string `json:"login" bson:"login"`
	// example: myplaintextpassword
	Password string `json:"password" bson:"password"`
}

// JWTResponse returns as HTTP response the user details (to be used along with the generated JWT token)
type JWTResponse struct {
	Type         string `json:"type"`
	RefreshToken string `json:"refresh"`
	AccessToken  string `json:"token"`
	// Details      SanitizedUser `json:"details,omitempty"`
}
