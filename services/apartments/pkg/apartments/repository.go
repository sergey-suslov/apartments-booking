package apartments

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const apartmentCollectionName = "apartments"

const maxApartmentLimit = 100

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{db: db}
}

func (r *repository) GetApartmentsByCity(ctx context.Context, city City, limit, offset int) ([]Apartment, error) {
	if limit > maxApartmentLimit {
		limit = maxApartmentLimit
	}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.db.Collection(apartmentCollectionName).Find(ctx, bson.D{{"city", city}}, opts)
	if err != nil {
		return nil, DatabaseError
	}
	apartments := make([]Apartment, 0, limit)
	err = cursor.All(ctx, &apartments)
	if err != nil {
		return nil, DatabaseError
	}
	return apartments, nil
}
