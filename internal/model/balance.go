package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Balance struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID      primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	Balance      BalanceInfo        `json:"balance" bson:"balance"`
	BalanceMonth time.Month         `json:"month,omitempty" bson:"month,omitempty"`
	BalanceYear  int                `json:"year,omitempty" bson:"year,omitempty"`
	Transactions []*Transaction     `json:"transactions,omitempty" bson:"transactions,omitempty"`
	// example: 2025-09-21T20:58:16.53Z
	CreatedAt primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type BalanceInfo struct {
	Total    float32
	Currency string
}
