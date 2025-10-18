package mongodb

import "context"

// UserMockCollectionConfig will implement mongodb collection functions
type HealthMockCollectionConfig struct {
	Error error
}

// CreateIndexes will create mongodb indexes
func (c *HealthMockCollectionConfig) Ping(ctx context.Context) (bool, error) {
	return true, nil
}
