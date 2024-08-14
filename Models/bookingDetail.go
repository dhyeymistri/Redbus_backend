package models

type BookSeatDoc struct {
	Booking            Booking  `json:"booking" bson:"booking"`
	SeatIDs            []string `json:"seatIDs" bson:"seatIDs"`
	PassengerNames     []string `json:"passengerNames" bson:"passengerNames"`
	PassengerGenders   []string `json:"passengerGenders" bson:"passengerGenders"`
	PassengerAges      []int    `json:"passengerAges" bson:"passengerAges"`
	BaseFare           int      `json:"baseFare" bson:"baseFare"`
	DiscountedAmount   int      `json:"discountedAmount" bson:"discountedAmount"`
	GST                int      `json:"gst" bson:"gst"`
	TotalPayableAmount int      `json:"totalPayableAmount" bson:"totalPayableAmount"`
}
