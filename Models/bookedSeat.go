package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookedSeats struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BookingID primitive.ObjectID `json:"bookingID" bson:"bookingID"`
	Seats     []Seat             `json:"seats" bson:"seats"`
}
