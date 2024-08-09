package ticket

import (
	connection "Redbus_backend/Config"
	Generic "Redbus_backend/Generic"
	models "Redbus_backend/Models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CancelTicket(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "DELETE" {
		w.Header().Set("Content-Type", "application/json")

		param := mux.Vars(r)
		ticketID := param["ticketID"]
		objTicketID, _ := primitive.ObjectIDFromHex(ticketID)

		ticketCollection := connection.ConnectDB("Tickets")
		filter := bson.M{"_id": objTicketID}
		var ticket models.Ticket
		err := ticketCollection.FindOne(context.TODO(), filter).Decode(&ticket)
		if err != nil {
			log.Fatal(err)
		}

		busID := ticket.BusID
		travelStartDate := ticket.PickDate
		travelStartTime := ticket.PickTime
		arrSeatID := ticket.SeatIDs

		rowSlice := []int{}
		colSlice := []int{}
		deckSlice := []bool{}

		bookingCollection := connection.ConnectDB("Bookings")
		var booking models.Booking

		filter = bson.M{"busID": busID, "travelStartDate": travelStartDate, "travelStartTime": travelStartTime}
		err = bookingCollection.FindOne(context.TODO(), filter).Decode(&booking)
		if err != nil {
			log.Fatal(err)
		}
		bookingID := booking.ID
		for _, seatID := range arrSeatID {
			objSeatID, _ := primitive.ObjectIDFromHex(seatID)
			for index := range booking.Bus.Seats {
				if booking.Bus.Seats[index].SeatID == objSeatID {
					booking.Bus.Seats[index].OnlyFemale = false
					booking.Bus.Seats[index].OnlyMale = false
					booking.Bus.Seats[index].SeatAvailibility = true
					booking.Bus.Seats[index].PassengerGender = ""
					booking.Bus.Seats[index].PassengerName = ""
					booking.Bus.Seats[index].PassengerAge = 0
					rowSlice = append(rowSlice, booking.Bus.Seats[index].Row)
					colSlice = append(colSlice, booking.Bus.Seats[index].Column)
					deckSlice = append(deckSlice, booking.Bus.Seats[index].IsUpperDeck)
				}
			}
		}

		////------------------------------------------------------
		booking.AvailableSeats += len(arrSeatID)

		//making only female and only male seats available for all
		for idx := range rowSlice {
			row := rowSlice[idx]
			col := colSlice[idx]
			deck := deckSlice[idx]
			for seatIndex := range booking.Bus.Seats {
				if (booking.Bus.Seats[seatIndex].Row == row-1 || booking.Bus.Seats[seatIndex].Row == row+1) && booking.Bus.Seats[seatIndex].Column == col && booking.Bus.Seats[seatIndex].IsUpperDeck == deck && booking.Bus.Seats[seatIndex].SeatAvailibility {
					booking.Bus.Seats[seatIndex].OnlyFemale = false
					booking.Bus.Seats[seatIndex].OnlyMale = false
				}
			}
		}

		if booking.AvailableSeats == booking.Bus.TotalSeats {
			filter = bson.M{"_id": bookingID}
			_, err := bookingCollection.DeleteOne(context.TODO(), filter)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			filter = bson.M{"_id": bookingID}
			update := bson.M{
				"$set": bson.M{
					"bus":            booking.Bus,
					"availableSeats": booking.AvailableSeats,
				},
			}
			_, err := bookingCollection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Fatal(err)
			}
		}

		filter = bson.M{"_id": ticket.ID}
		_, err = ticketCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}

		userCollection := connection.ConnectDB("Users")
		baseFare := ticket.BaseFare
		currentTime := time.Now()
		busStartTimeStr := ticket.PickDate + " " + ticket.PickTime
		busStartTime, _ := time.Parse("2006-01-02 15:04", busStartTimeStr)
		var refundAmount int
		if currentTime.After(busStartTime) {
			json.NewEncoder(w).Encode("You cannot cancel ticket after the travel has begun")
			return
		}
		duration := busStartTime.Sub(currentTime)
		if duration < 12*time.Hour {
			refundAmount = 0
		} else if duration < 24*time.Hour {
			refundAmount = int(float64(baseFare) * 0.25)
		} else if duration < 48*time.Hour {
			refundAmount = int(float64(baseFare) * 0.5)
		} else {
			refundAmount = baseFare
		}
		update := bson.M{
			"$inc": bson.M{
				"walletBalance": refundAmount,
			},
		}
		_, err = userCollection.UpdateByID(context.TODO(), ticket.CustomerID, update)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode("Ticket has been cancelled. Your refund amount is " + strconv.Itoa(refundAmount))
	}
}

func GetTicketByUserID(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	collection := connection.ConnectDB("Tickets")

	var arrTickets []models.Ticket
	var params = mux.Vars(r)

	id := params["userID"]
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"customerID": objID}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	if err = cursor.All(context.Background(), &arrTickets); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(arrTickets)
}
