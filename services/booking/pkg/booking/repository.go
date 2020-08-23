package booking

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const reservationCollectionName = "reservations"

var ErrWrongIDFormat = errors.New("wrong id format")

type MongoReservationsRepository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *MongoReservationsRepository {
	return &MongoReservationsRepository{db: db}
}

func (r *MongoReservationsRepository) GetReservationsBetween(ctx context.Context, apartmentID string, start, end time.Time) ([]Reservation, error) { //nolint:lll
	opts := options.Find().SetSort(bson.D{primitive.E{
		Key: "created", Value: 1,
	}})
	objectID, err := primitive.ObjectIDFromHex(apartmentID)
	if err != nil {
		return nil, ErrWrongIDFormat
	}
	cursor, err := r.db.Collection(reservationCollectionName).Find(
		ctx,
		bson.D{
			primitive.E{Key: "apartmentID", Value: objectID},
			primitive.E{Key: "start", Value: bson.D{{"$gte", start}}},
			primitive.E{Key: "end", Value: bson.D{{"$lte", end}}},
		}, opts)
	if err != nil {
		return nil, ErrRequestingDatabase
	}
	reservations := make([]Reservation, 0, 40)
	err = cursor.All(ctx, &reservations)
	if err != nil {
		return nil, ErrRequestingDatabase
	}
	return reservations, nil
}

func (r *MongoReservationsRepository) MakeReservation(ctx context.Context, reservation *Reservation) (*Reservation, error) {
	result, err := r.db.Collection(reservationCollectionName).InsertOne(
		ctx,
		bson.M{
			"apartmentId": reservation.ApartmentID,
			"userId":      reservation.UserID,
			"start":       reservation.Start,
			"end":         reservation.End,
			"created":     reservation.Created},
	)
	if err != nil {
		return nil, err
	}
	reservation.ID = result.InsertedID.(primitive.ObjectID)
	return reservation, nil
}
