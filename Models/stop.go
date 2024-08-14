package models

type Stop struct {
	Location      string `json:"location" bson:"location"`
	ArrivalTime   string `json:"arrivalTime" bson:"arrivalTime"`
	DepartureTime string `json:"departureTime" bson:"departureTime"`
}
