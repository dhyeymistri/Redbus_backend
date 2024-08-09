package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Route struct {
	Location      string               `bson:"location"`
	Buses         []primitive.ObjectID `bson:"buses"`
	ArrivalTime   []string             `bson:"arrivalTimings"`
	DepartureTime []string             `bson:"departureTimings"`
	IsWeekend     []bool               `bson:"isWeekend"`
}
