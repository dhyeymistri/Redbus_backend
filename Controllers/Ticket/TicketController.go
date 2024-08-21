package ticket

import (
	connection "Redbus_backend/Config"
	Generic "Redbus_backend/Generic"
	getbusdetail "Redbus_backend/Helpers/GetBusDetail"
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
			json.NewEncoder(w).Encode("This ticket does not exist")
			return
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
		bus := getbusdetail.GetBusDetail(booking.BusID)

		var bookedSeats models.BookedSeats
		seatsCollection := connection.ConnectDB("BookedSeats")
		filter = bson.M{"bookingID": bookingID}
		err = seatsCollection.FindOne(context.TODO(), filter).Decode(&bookedSeats)
		if err != nil {
			log.Fatal(err)
		}
		seats := bookedSeats.Seats
		for _, seatID := range arrSeatID {
			objSeatID, _ := primitive.ObjectIDFromHex(seatID)
			for index := range seats {
				if seats[index].SeatID == objSeatID {
					seats[index].OnlyFemale = false
					seats[index].OnlyMale = false
					seats[index].SeatAvailibility = true
					seats[index].PassengerGender = ""
					seats[index].PassengerName = ""
					seats[index].PassengerAge = 0
					rowSlice = append(rowSlice, seats[index].Row)
					colSlice = append(colSlice, seats[index].Column)
					deckSlice = append(deckSlice, seats[index].IsUpperDeck)
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
			for seatIndex := range seats {
				if (seats[seatIndex].Row == row-1 || seats[seatIndex].Row == row+1) && seats[seatIndex].Column == col && seats[seatIndex].IsUpperDeck == deck && seats[seatIndex].SeatAvailibility {
					seats[seatIndex].OnlyFemale = false
					seats[seatIndex].OnlyMale = false
				}
			}
		}

		if booking.AvailableSeats == bus.TotalSeats {
			filter = bson.M{"bookingID": bookingID}
			_, err := seatsCollection.DeleteOne(context.TODO(), filter)
			if err != nil {
				log.Fatal(err)
			}
			filter = bson.M{"_id": bookingID}
			_, err = bookingCollection.DeleteOne(context.TODO(), filter)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			filter = bson.M{"bookingID": bookingID}
			update := bson.M{
				"$set": bson.M{
					"seats": seats,
				},
			}
			_, err := seatsCollection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Fatal(err)
			}

			filter = bson.M{"_id": bookingID}
			update = bson.M{
				"$set": bson.M{
					"availableSeats": booking.AvailableSeats,
				},
			}
			_, err = bookingCollection.UpdateOne(context.TODO(), filter, update)
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
		currentTime := time.Now()
		busStartTimeStr := ticket.PickDate + " " + ticket.PickTime
		busStartTime, _ := time.Parse("2006-01-02 15:04", busStartTimeStr)
		var refundAmount int
		if currentTime.After(busStartTime) {
			json.NewEncoder(w).Encode("You cannot cancel ticket after the travel has begun")
			return
		}
		duration := busStartTime.Sub(currentTime)
		consideredRefund := ticket.TotalPayableAmount - ticket.GST
		if duration < 12*time.Hour {
			refundAmount = 0
		} else if duration < 24*time.Hour {
			refundAmount = int(float64(consideredRefund) * 0.25)
		} else if duration < 48*time.Hour {
			refundAmount = int(float64(consideredRefund) * 0.5)
		} else {
			refundAmount = consideredRefund
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
	if len(arrTickets) == 0 {
		json.NewEncoder(w).Encode("No booked tickets by this user or the user ID is not correct, check it")
		return
	}

	json.NewEncoder(w).Encode(arrTickets)
}
