package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Transaction defines a user transaction to be added to a posterior Balance
type Transaction struct {
	// swagger:ignore
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BalanceID primitive.ObjectID `json:"balance_id,omitempty" bson:"balance_id,omitempty"`
	OwnerID   primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	// example: income/expense
	Type string `json:"type" bson:"type"`
	// example: guitar lessons
	Description string `json:"description" bson:"description"`
	// example: 12.90
	Amount float64 `json:"amount" bson:"amount"`
	// example: BRL
	Currency string `json:"currency" bson:"currency"`
	// example: Credit
	PaymentMethod PaymentMethod `json:"payment_method" bson:"payment_method"`
	// example: 2025-09-16T17:33:10.64Z
	TransactionDate primitive.DateTime `json:"transaction_date,omitempty" bson:"transaction_date,omitempty"`
	// example: "categories": ["personal development"]
	Categories []string `json:"categories,omitempty" bson:"categories,omitempty"`
	// example: 2025-09-21T20:58:16.53Z
	CreatedAt primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// PaymentMethod defines which payment method was used for a certain transaction
type PaymentMethod struct {
	Credit      Card `json:"credit" bson:"credit"`
	Debit       Card `json:"debit" bson:"debit"`
	Pix         bool `json:"pix,omitempty" bson:"pix,omitempty"`
	PaymentSlip bool `json:"payment_slip,omitempty" bson:"payment_slip,omitempty"`
}
