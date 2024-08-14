package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Seat struct {
	SeatID           primitive.ObjectID `json:"_id" bson:"_id"`
	SeatType         string             `json:"seatType" bson:"seatType"`
	OnlyFemale       bool               `json:"onlyFemale" bson:"onlyFemale"`
	OnlyMale         bool               `json:"onlyMale" bson:"onlyMale"`
	PassengerGender  string             `json:"passengerGender" bson:"passengerGender"`
	SeatAvailibility bool               `json:"availability" bson:"availability"`
	PassengerName    string             `json:"passengerName" bson:"passengerName"`
	PassengerAge     int                `json:"passengerAge" bson:"passengerAge"`
	Row              int                `json:"row" bson:"row"`
	Column           int                `json:"column" bson:"column"`
	Group            string             `json:"group" bson:"group"`
	IsUpperDeck      bool               `json:"isUpperDeck" bson:"isUpperDeck"`
	Class            string             `json:"class" bson:"class"`
	Cost             int                `json:"cost" bson:"cost"`
	IsSelected       bool               `json:"isSelected" bson:"isSelected"`
}
