package minorhelpers

import (
	models "Redbus_backend/Models"
	"fmt"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindCommonElements(arr1, arr2 []primitive.ObjectID) []primitive.ObjectID {
	elementSet := make(map[primitive.ObjectID]struct{})
	for _, elem := range arr1 {
		elementSet[elem] = struct{}{}
	}
	var commonElements []primitive.ObjectID

	for _, elem := range arr2 {
		if _, exists := elementSet[elem]; exists {
			commonElements = append(commonElements, elem)
			// Remove the element from the set to avoid duplicates in the result
			delete(elementSet, elem)
		}
	}

	return commonElements
}

func HasDateChanged(startTimeStr, endTimeStr string) bool {
	// Define the time layout (format)
	layout := "15:04"

	// Parse the start and end times
	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		fmt.Println("Error parsing start time:", err)
		return false
	}

	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		fmt.Println("Error parsing end time:", err)
		return false
	}

	// Assume today’s date for the start and end times
	now := time.Now().Local()
	startTime = time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), startTime.Local().Hour(), time.UTC)
	endTime = time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), endTime.Local().Hour(), time.UTC)

	// If end time is earlier than start time, adjust end time to the next day
	if endTime.Before(startTime) {
		endTime = endTime.Add(24 * time.Hour)
	}

	// Compare the dates
	return endTime.YearDay() != startTime.YearDay()
}

func TimeDifference(t1, t2 string) (time.Duration, error) {
	// Define the time layout for parsing
	const layout = "15:04"

	// Parse the input times
	time1, err := time.Parse(layout, t1)
	if err != nil {
		return 0, err
	}
	time2, err := time.Parse(layout, t2)
	if err != nil {
		return 0, err
	}

	diff := time2.Sub(time1)

	// If the difference is negative, adjust by adding 24 hours
	if diff < 0 {
		diff += 24 * time.Hour
	}

	return diff, nil
}

func MakeTicket(user models.User, booking models.Booking, bus models.Bus, bookingDetails models.BookSeatDoc) models.Ticket {
	var ticket models.Ticket
	ticket.CustomerID = user.ID
	ticket.CustomerName = user.FirstName + " " + user.LastName
	ticket.BusID = booking.BusID
	ticket.BusName = bus.OperatorName
	ticket.DropAddress = booking.TravelEndLocation
	ticket.DropDate = booking.TravelEndDate
	ticket.Email = user.Email
	ticket.PickDate = booking.TravelStartDate
	ticket.DropTime = booking.TravelEndTime
	ticket.PickTime = booking.TravelStartTime
	ticket.PickupAddress = booking.TravelStartLocation
	ticket.BaseFare = bookingDetails.BaseFare
	ticket.DiscountedAmount = bookingDetails.DiscountedAmount
	ticket.GST = bookingDetails.GST
	ticket.PassengerAges = bookingDetails.PassengerAges
	ticket.PassengerGenders = bookingDetails.PassengerGenders
	ticket.PassengerNames = bookingDetails.PassengerNames
	ticket.SeatIDs = bookingDetails.SeatIDs
	ticket.TotalPayableAmount = bookingDetails.TotalPayableAmount
	ticket.TotalPassenger = len(bookingDetails.SeatIDs)

	return ticket
}

func AllotSeatPrices(bus models.Bus, timeDiffInMinutes float64) []models.Seat {
	seats := bus.Seats
	seaterPrice := bus.SeaterCostPerMinute
	sleeperPrice := bus.SleeperCostPerMinute
	for idx, seat := range seats {
		if seat.Class == "C" {
			if seat.SeatType == "SE" {
				seats[idx].Cost = int(timeDiffInMinutes * seaterPrice)
			} else {
				seats[idx].Cost = int(timeDiffInMinutes * sleeperPrice)
			}
		} else if seat.Class == "B" {
			if seat.SeatType == "SE" {
				tmp := timeDiffInMinutes * seaterPrice * 1.2
				seats[idx].Cost = int(tmp)
			} else {
				tmp := timeDiffInMinutes * sleeperPrice * 1.2
				seats[idx].Cost = int(tmp)
			}
		} else if seat.Class == "A" {
			if seat.SeatType == "SE" {
				tmp := timeDiffInMinutes * seaterPrice * 1.5
				seats[idx].Cost = int(tmp)
			} else {
				tmp := timeDiffInMinutes * sleeperPrice * 1.5
				seats[idx].Cost = int(tmp)
			}
		}
	}
	return seats
}

func ExtractFileNames(files []*multipart.FileHeader) []string {
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Filename)
	}
	return fileNames
}
