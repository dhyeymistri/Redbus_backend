package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Review struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Rating     int                `json:"rating" bson:"rating"`
	ReviewText string             `json:"reviewText,omitempty" bson:"reviewText,omitempty"`
	BusID      primitive.ObjectID `json:"busID" bson:"busID"`
	CustomerID primitive.ObjectID `json:"customerID" bson:"customerID"`
}
