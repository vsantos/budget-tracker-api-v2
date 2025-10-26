package mongodb

import (
	"budget-tracker-api-v2/internal/model"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransactionMockCollectionConfig will implement mongodb collection functions
type TransactionMockCollectionConfig struct {
	Error error
}

var (
	objID, balanceID, ownerID primitive.ObjectID
	tDate                     time.Time
)

func init() {
	objID, _ = primitive.ObjectIDFromHex("68fd6b00f4c9e77e59aaf97e")
	balanceID, _ = primitive.ObjectIDFromHex("687baad049572fb8c4e305f3")
	ownerID, _ = primitive.ObjectIDFromHex("66f1cca3c37c733c4ada103d")

	mockedTransactionDate := "2023-10-26 14:30:00"
	layout := "2006-01-02 15:04:05"
	tDate, _ = time.Parse(layout, mockedTransactionDate)
}

// CreateIndexes will create mongodb indexes
func (c *TransactionMockCollectionConfig) CreateIndexes(ctx context.Context, indexes []string) error {
	return nil
}

// InsertOne will insert a document into mongodb
func (c *TransactionMockCollectionConfig) InsertOne(ctx context.Context, emp *model.Transaction) (transaction *model.Transaction, err error) {
	if c.Error != nil {
		return nil, c.Error
	}

	return &model.Transaction{
		ID:              objID,
		BalanceID:       balanceID,
		OwnerID:         ownerID,
		Type:            "income",
		Description:     "My favorite chinese restaurant",
		Amount:          15.3,
		Currency:        "BRL",
		PaymentMethod:   model.PaymentMethod{},
		TransactionDate: primitive.NewDateTimeFromTime(tDate),
		Categories:      []string{"food"},
		CreatedAt:       primitive.NewDateTimeFromTime(tDate),
	}, nil
}

// FindOne will find a document based on ID
func (c *TransactionMockCollectionConfig) FindOne(ctx context.Context, id string) (*model.Transaction, error) {
	receivedID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if receivedID != objID {
		return nil, errors.New("not found")
	}

	return &model.Transaction{
		ID:              objID,
		BalanceID:       balanceID,
		OwnerID:         ownerID,
		Type:            "income",
		Description:     "My favorite chinese restaurant",
		Amount:          15.3,
		Currency:        "BRL",
		PaymentMethod:   model.PaymentMethod{},
		TransactionDate: primitive.NewDateTimeFromTime(tDate),
		Categories:      []string{"food"},
		CreatedAt:       primitive.NewDateTimeFromTime(tDate),
	}, nil
}

// DeleteOne will insert a document into mongodb
func (c *TransactionMockCollectionConfig) DeleteOne(ctx context.Context, id string) (int64, error) {
	return 1, nil
}
