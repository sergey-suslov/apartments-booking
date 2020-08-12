package booking

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const reservationCollectionName = "reservations"

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{db: db}
}

func (r *repository) GetReservationsBetween(ctx context.Context, apartmentId string, start, end time.Time) ([]Reservation, error) {
	opts := options.Find().SetSort(bson.D{{"created", 1}})
	cursor, err := r.db.Collection(reservationCollectionName).Find(
		ctx,
		bson.D{
			{"apartmentId", apartmentId},
			{"start", bson.D{{"$gte", start}}},
			{"end", bson.D{{"$lte", end}}},
		}, opts)
	if err != nil {
		return nil, DatabaseError
	}
	reservations := make([]Reservation, 0, 40)
	err = cursor.All(ctx, &reservations)
	if err != nil {
		return nil, DatabaseError
	}
	return reservations, nil
}
