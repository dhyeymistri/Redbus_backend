package models

type Stop struct {
	Location      string `bson:"location"`
	ArrivalTime   string `bson:"arrivalTime"`
	DepartureTime string `bson:"departureTime"`
}
