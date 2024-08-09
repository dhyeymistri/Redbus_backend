package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bus struct {
	ID                   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OperatorName         string             `json:"operatorName" bson:"operatorName"`
	ModelDetails         string             `json:"modelDetails" bson:"modelDetails"`
	TotalSeats           int                `json:"totalSeats" bson:"totalSeats"`
	ImgPath              []string           `json:"imgPath" bson:"imgPath"`
	Amenities            []string           `json:"amenities" bson:"amenities"`
	AverageRating        float64            `json:"avgRating" bson:"avgRating"`
	NumberOfReviews      int                `json:"numberOfReviews" bson:"numberOfReviews"`
	LiveTracking         bool               `json:"liveTracking" bson:"liveTracking"`
	IsAcAvailable        bool               `json:"isAcAvailable" bson:"isAcAvailable"`
	BusType              string             `json:"busType" bson:"busType"`
	Stops                []Stop             `json:"stops" bson:"stops"` //stops are in order
	Seats                []Seat             `json:"seats" bson:"seats"`
	Frequency            string             `json:"frequency" bson:"frequency"` //it is either daily or weekends
	AvailableSeats       int                `json:"seatAvailability" bson:"seatAvailability"`
	SleeperCostPerMinute float64            `json:"sleeperCost" bson:"sleeperCost"`
	SeaterCostPerMinute  float64            `json:"seaterCost" bson:"seaterCost"`
}
