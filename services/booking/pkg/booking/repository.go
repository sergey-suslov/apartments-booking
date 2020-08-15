package booking

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const reservationCollectionName = "reservations"

var wrongIdFormat = errors.New("wrong id format")

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{db: db}
}

func (r *repository) GetReservationsBetween(ctx context.Context, apartmentId string, start, end time.Time) ([]Reservation, error) {
	opts := options.Find().SetSort(bson.D{{"created", 1}})
	objectID, err := primitive.ObjectIDFromHex(apartmentId)
	if err != nil {
		return nil, wrongIdFormat
	}
	cursor, err := r.db.Collection(reservationCollectionName).Find(
		ctx,
		bson.D{
			{"apartmentId", objectID},
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

func (r *repository) MakeReservation(ctx context.Context, reservation *Reservation) (*Reservation, error) {
	result, err := r.db.Collection(reservationCollectionName).InsertOne(
		ctx,
		bson.M{
			"apartmentId": reservation.ApartmentId,
			"userId":      reservation.UserId,
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
