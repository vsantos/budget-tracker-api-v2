package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Spend defines a user spend to be added to Balance
// swagger:model
type Spend struct {
	// swagger:ignore
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	// example: fixed
	Type string `json:"type" bson:"type"`
	// example: guitar lessons
	Description string `json:"description" bson:"description"`
	// example: 12.90
	Cost float64 `json:"cost" bson:"cost"`
	// example: debit: true
	PaymentMethod PaymentMethod `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	// example: "categories": ["personal development"]
	Categories []string `json:"categories,omitempty" bson:"categories,omitempty"`
	// swagger:ignore
	CreatedAt primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// PaymentMethod defines which payment method was used for a certain spend
// swagger:model
type PaymentMethod struct {
	Credit      Card `json:"credit,omitempty" bson:"credit,omitempty"`
	Debit       bool `json:"debit,omitempty" bson:"debit,omitempty"`
	PaymentSlip bool `json:"payment_slip,omitempty" bson:"payment_slip,omitempty"`
}
