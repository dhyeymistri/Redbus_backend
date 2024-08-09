package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Bus
	BusID               primitive.ObjectID `json:"busID" bson:"busID"`
	TravelStartDate     string             `json:"travelStartDate" bson:"travelStartDate"`
	TravelEndDate       string             `json:"travelEndDate" bson:"travelEndDate"`
	TravelStartTime     string             `json:"travelStartTime" bson:"travelStartTime"`
	TravelEndTime       string             `json:"travelEndTime" bson:"travelEndTime"`
	TravelStartLocation string             `json:"travelStartLocation" bson:"travelStartLocation"`
	TravelEndLocation   string             `json:"travelEndLocation" bson:"travelEndLocation"`
}
