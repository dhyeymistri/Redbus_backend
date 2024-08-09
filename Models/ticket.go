package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID         primitive.ObjectID `json:"customerID" bson:"customerID"`
	CustomerName       string             `json:"customerName" bson:"customerName"`
	Email              string             `json:"email" bson:"email"`
	TotalPassenger     int                `json:"totalPassenger" bson:"totalPassenger"`
	BusID              primitive.ObjectID `json:"busID" bson:"busID"`
	BusName            string             `json:"busName" bson:"busName"`
	PickupAddress      string             `json:"pickupAddress" bson:"pickupAddress"`
	DropAddress        string             `json:"dropAddress" bson:"dropAddress"`
	PickDate           string             `json:"pickDate" bson:"pickDate"`
	DropDate           string             `json:"dropDate" bson:"dropDate"`
	PickTime           string             `json:"pickTime" bson:"pickTime"`
	DropTime           string             `json:"dropTime" bson:"dropTime"`
	SeatIDs            []string           `json:"seatIDs" bson:"seatIDs"`
	PassengerNames     []string           `json:"passengerNames" bson:"passengerNames"`
	PassengerGenders   []string           `json:"passengerGenders" bson:"passengerGenders"`
	PassengerAges      []int              `json:"passengerAges" bson:"passengerAges"`
	BaseFare           int                `json:"baseFare" bson:"baseFare"`
	DiscountedAmount   int                `json:"discountedAmount" bson:"discountedAmount"`
	GST                int                `json:"gst" bson:"gst"`
	TotalPayableAmount int                `json:"totalPayableAmount" bson:"totalPayableAmount"`
}
