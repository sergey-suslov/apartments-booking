package apartments

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const apartmentCollectionName = "apartments"

const maxApartmentLimit = 100

type MongoRepositoryApartments struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *MongoRepositoryApartments {
	return &MongoRepositoryApartments{db: db}
}

func (r *MongoRepositoryApartments) GetApartmentsByCity(ctx context.Context, city City, limit, offset int) ([]Apartment, error) {
	if limit > maxApartmentLimit {
		limit = maxApartmentLimit
	}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.db.Collection(apartmentCollectionName).Find(ctx, bson.D{{"city", city}}, opts)
	if err != nil {
		return nil, ErrDatabase
	}
	apartments := make([]Apartment, 0, limit)
	err = cursor.All(ctx, &apartments)
	if err != nil {
		return nil, ErrDatabase
	}
	return apartments, nil
}

func (r *MongoRepositoryApartments) GetApartmentByID(ctx context.Context, apartmentID string) (a *Apartment, err error) {
	objectID, err := primitive.ObjectIDFromHex(apartmentID)
	if err != nil {
		return nil, err
	}

	result := r.db.Collection(apartmentCollectionName).FindOne(ctx, bson.D{{"_id", objectID}})
	var apartment Apartment
	err = result.Decode(&apartment)
	if err != nil {
		return nil, err
	}
	return &apartment, nil
}
