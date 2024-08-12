package booking

import (
	connection "Redbus_backend/Config"
	Generic "Redbus_backend/Generic"
	auth "Redbus_backend/Helpers/Auth"
	models "Redbus_backend/Models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookedSeats struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BookingID primitive.ObjectID `json:"bookingID" bson:"bookingID"`
	Seats     []models.Seat      `json:"seats" bson:"seats"`
}

type BookSeatDoc struct {
	Booking            models.Booking `json:"booking" bson:"booking"`
	SeatIDs            []string       `json:"seatIDs" bson:"seatIDs"`
	PassengerNames     []string       `json:"passengerNames" bson:"passengerNames"`
	PassengerGenders   []string       `json:"passengerGenders" bson:"passengerGenders"`
	PassengerAges      []int          `json:"passengerAges" bson:"passengerAges"`
	BaseFare           int            `json:"baseFare" bson:"baseFare"`
	DiscountedAmount   int            `json:"discountedAmount" bson:"discountedAmount"`
	GST                int            `json:"gst" bson:"gst"`
	TotalPayableAmount int            `json:"totalPayableAmount" bson:"totalPayableAmount"`
}

func BookSeat(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//data is array of bytes
		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error with data retrieval")
		}
		asString := string(data)
		var bookingDetails BookSeatDoc
		// rr := strings.NewReader(asString)

		json.Unmarshal([]byte(asString), &bookingDetails)
		booking := bookingDetails.Booking
		json.NewEncoder(w).Encode(booking)
		travelStartDate := booking.TravelStartDate
		travelEndDate := booking.TravelEndDate
		arrPassengerName := bookingDetails.PassengerNames
		arrPassengerGender := bookingDetails.PassengerGenders
		arrPassengerAge := bookingDetails.PassengerAges
		busID := booking.BusID
		seatIDString := bookingDetails.SeatIDs
		var arrSeatID []primitive.ObjectID
		for _, seat := range seatIDString {
			seatID, _ := primitive.ObjectIDFromHex(seat)
			arrSeatID = append(arrSeatID, seatID)
		}

		cookie, err := r.Cookie("token")
		if cookie.Value == "" {
			json.NewEncoder(w).Encode("User is not logged in")
		}
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userID := auth.GetUserID(cookie.Value)
		objUserID, _ := primitive.ObjectIDFromHex(userID)
		fmt.Println(objUserID)
		fmt.Println(busID, travelStartDate, travelEndDate, arrPassengerAge, arrPassengerName, arrPassengerGender, arrSeatID)

		// objBusID, _ := primitive.ObjectIDFromHex(busID)
		bookingCollection := connection.ConnectDB("Bookings")
		bookingsFilter := bson.M{"busID": busID, "travelStartDate": travelStartDate, "travelEndDate": travelEndDate}
		// var booking models.Booking

		userCollection := connection.ConnectDB("Users")
		// userID, ok := getUserID(r)
		filter := bson.M{"_id": objUserID}
		totalPayableAmount := bookingDetails.TotalPayableAmount
		var user models.User
		err = userCollection.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		if user.WalletBalance < totalPayableAmount {
			json.NewEncoder(w).Encode("You don't have sufficient balance to book the tickets")
			return
		}

		var foundBooking models.Booking
		rowSlice := []int{}
		colSlice := []int{}
		deckSlice := []bool{}
		err = bookingCollection.FindOne(context.TODO(), bookingsFilter).Decode(&foundBooking)
		if err != nil {
			for outerIndex := range arrSeatID {
				for idx := range booking.Bus.Seats {
					fmt.Println("CHECK")
					if booking.Bus.Seats[idx].SeatID == arrSeatID[outerIndex] {
						if arrPassengerGender[outerIndex] == "F" && booking.Bus.Seats[idx].OnlyMale {
							json.NewEncoder(w).Encode("This is an only male seat and female passenger cannot book this for safety purpose")
							return
						}
						if arrPassengerGender[outerIndex] == "M" && booking.Bus.Seats[idx].OnlyFemale {
							json.NewEncoder(w).Encode("This is an only female seat and a male passenger cannot book this for safety purpose")
							return
						}
						booking.Bus.Seats[idx].PassengerAge = int(arrPassengerAge[outerIndex])
						booking.Bus.Seats[idx].PassengerGender = arrPassengerGender[outerIndex]
						booking.Bus.Seats[idx].PassengerName = arrPassengerName[outerIndex]
						booking.Bus.Seats[idx].SeatAvailibility = false
						rowSlice = append(rowSlice, booking.Bus.Seats[idx].Row)
						colSlice = append(colSlice, booking.Bus.Seats[idx].Column)
						deckSlice = append(deckSlice, booking.Bus.Seats[idx].IsUpperDeck)
						if arrPassengerGender[outerIndex] == "F" {
							booking.Bus.Seats[idx].OnlyFemale = true
							booking.Bus.Seats[idx].OnlyMale = false
							// groupNo = booking.Bus.Seats[idx].Group
						} else {
							booking.Bus.Seats[idx].OnlyFemale = false
							booking.Bus.Seats[idx].OnlyMale = true
							// groupNo = booking.Bus.Seats[idx].Group
						}
					}
				}
			}
			booking.AvailableSeats = booking.AvailableSeats - len(arrPassengerAge)

			//gender verification-------------------------------------------------
			for idx := range rowSlice {
				row := rowSlice[idx]
				col := colSlice[idx]
				deck := deckSlice[idx]
				for seatIndex := range booking.Bus.Seats {
					if (booking.Bus.Seats[seatIndex].Row == row-1 || booking.Bus.Seats[seatIndex].Row == row+1) && booking.Bus.Seats[seatIndex].Column == col && booking.Bus.Seats[seatIndex].IsUpperDeck == deck && booking.Bus.Seats[seatIndex].SeatAvailibility {
						if arrPassengerGender[idx] == "F" {
							booking.Bus.Seats[seatIndex].OnlyFemale = true
						} else {
							booking.Bus.Seats[seatIndex].OnlyMale = true
						}
					}
				}
			}

			update := bson.M{"$set": bson.M{"bus": booking.Bus, "availableSeats": booking.AvailableSeats, "travelStartTime": booking.TravelStartTime, "travelEndTime": booking.TravelEndTime, "travelStartLocation": booking.TravelStartLocation, "travelEndLocation": booking.TravelEndLocation}}

			// Update the BookSeatDoc if it exists, otherwise insert a new BookSeatDoc
			opts := options.Update().SetUpsert(true)
			result, err := bookingCollection.UpdateOne(context.TODO(), bookingsFilter, update, opts)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result.UpsertedID)
			bookedSeats := BookedSeats{
				BookingID: result.UpsertedID.(primitive.ObjectID),
				Seats:     booking.Bus.Seats,
			}
			filter = bson.M{"bookingID": bookedSeats.BookingID}

			seatsCollection := connection.ConnectDB("BookedSeats")
			update = bson.M{
				"$set": bson.M{
					"seats": bookedSeats.Seats,
				},
			}
			opts = options.Update().SetUpsert(true)
			result, err = seatsCollection.UpdateOne(context.TODO(), filter, update, opts)
			if err != nil {
				log.Fatal(err)
			}

			//ticket------------------------------------------------
			var ticket models.Ticket

			userCollection := connection.ConnectDB("Users")
			var user models.User

			err = userCollection.FindOne(context.TODO(), bson.M{"_id": objUserID}).Decode(&user)
			if err != nil {
				log.Fatal(err)
			}
			ticket.CustomerID = objUserID
			ticket.CustomerName = user.FirstName + " " + user.LastName
			ticket.BusID = booking.BusID
			ticket.BusName = booking.Bus.OperatorName
			ticket.DropAddress = booking.TravelEndLocation
			ticket.DropDate = booking.TravelEndDate
			ticket.Email = user.Email
			ticket.DropTime = booking.TravelEndTime
			ticket.PickTime = booking.TravelStartTime
			ticket.PickDate = booking.TravelStartDate
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

			ticketCollection := connection.ConnectDB("Tickets")
			insertedResult, err := ticketCollection.InsertOne(context.TODO(), ticket)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(insertedResult.InsertedID)
		} else {
			myBus := foundBooking.Bus
			for outerIndex := range arrSeatID {
				for idx := range myBus.Seats {
					if myBus.Seats[idx].SeatID == arrSeatID[outerIndex] {
						if !myBus.Seats[idx].SeatAvailibility {
							json.NewEncoder(w).Encode("This seat has been booked already")
							return
						}
						if arrPassengerGender[outerIndex] == "F" && myBus.Seats[idx].OnlyMale {
							json.NewEncoder(w).Encode("This is an only male seat and female passenger cannot book this for safety purpose")
							return
						}
						if arrPassengerGender[outerIndex] == "M" && myBus.Seats[idx].OnlyFemale {
							json.NewEncoder(w).Encode("This is an only female seat and a male passenger cannot book this for safety purpose")
							return
						}
						myBus.Seats[idx].PassengerAge = int(arrPassengerAge[outerIndex])
						myBus.Seats[idx].PassengerGender = arrPassengerGender[outerIndex]
						myBus.Seats[idx].PassengerName = arrPassengerName[outerIndex]
						myBus.Seats[idx].SeatAvailibility = false
						rowSlice = append(rowSlice, booking.Bus.Seats[idx].Row)
						colSlice = append(colSlice, booking.Bus.Seats[idx].Column)
						deckSlice = append(deckSlice, booking.Bus.Seats[idx].IsUpperDeck)
					}
				}
			}
			foundBooking.AvailableSeats = foundBooking.AvailableSeats - len(arrPassengerAge)

			for idx := range rowSlice {
				row := rowSlice[idx]
				col := colSlice[idx]
				deck := deckSlice[idx]
				for seatIndex := range foundBooking.Bus.Seats {
					if (foundBooking.Bus.Seats[seatIndex].Row == row-1 || foundBooking.Bus.Seats[seatIndex].Row == row+1) && foundBooking.Bus.Seats[seatIndex].Column == col && foundBooking.Bus.Seats[seatIndex].IsUpperDeck == deck && foundBooking.Bus.Seats[seatIndex].SeatAvailibility {
						if arrPassengerGender[idx] == "F" {
							foundBooking.Bus.Seats[seatIndex].OnlyFemale = true
						} else {
							foundBooking.Bus.Seats[seatIndex].OnlyMale = true
						}
					}
				}
			}

			update := bson.M{"$set": bson.M{"bus": foundBooking.Bus, "availableSeats": foundBooking.AvailableSeats}}

			// Update the BookSeatDoc if it exists, otherwise insert a new BookSeatDoc
			opts := options.Update().SetUpsert(true)
			_, err := bookingCollection.UpdateOne(context.TODO(), bookingsFilter, update, opts)
			if err != nil {
				log.Fatal(err)
			}

			//ticket------------------------------------------------
			var ticket models.Ticket
			ticket.CustomerID = objUserID
			userCollection := connection.ConnectDB("Users")
			var user models.User
			err = userCollection.FindOne(context.TODO(), bson.M{"_id": objUserID}).Decode(&user)
			if err != nil {
				log.Fatal(err)
			}
			ticket.CustomerName = user.FirstName + " " + user.LastName
			ticket.BusID = foundBooking.BusID
			ticket.BusName = foundBooking.Bus.OperatorName
			ticket.DropAddress = foundBooking.TravelEndLocation
			ticket.DropDate = foundBooking.TravelEndDate
			ticket.Email = user.Email
			ticket.PickDate = foundBooking.TravelStartDate
			ticket.DropTime = foundBooking.TravelEndTime
			ticket.PickTime = foundBooking.TravelStartTime
			ticket.PickupAddress = foundBooking.TravelStartLocation
			ticket.BaseFare = bookingDetails.BaseFare
			ticket.DiscountedAmount = bookingDetails.DiscountedAmount
			ticket.GST = bookingDetails.GST
			ticket.PassengerAges = bookingDetails.PassengerAges
			ticket.PassengerGenders = bookingDetails.PassengerGenders
			ticket.PassengerNames = bookingDetails.PassengerNames
			ticket.SeatIDs = bookingDetails.SeatIDs
			ticket.TotalPayableAmount = bookingDetails.TotalPayableAmount
			ticket.TotalPassenger = len(bookingDetails.SeatIDs)

			ticketCollection := connection.ConnectDB("Tickets")
			result, err := ticketCollection.InsertOne(context.TODO(), ticket)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(result.InsertedID)
		}

		update := bson.M{
			"$inc": bson.M{
				"walletBalance": -totalPayableAmount,
			},
		}
		_, err = userCollection.UpdateByID(context.TODO(), objUserID, update)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func SelectSeat(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//data is array of bytes
		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error with data retrieval")
		}
		asString := string(data)
		var selectSeatDetails models.Booking

		var params = mux.Vars(r)

		id := params["seatID"]
		objID, _ := primitive.ObjectIDFromHex(id)
		fmt.Println(objID)

		json.Unmarshal([]byte(asString), &selectSeatDetails)

		for idx := range selectSeatDetails.Bus.Seats {
			if selectSeatDetails.Bus.Seats[idx].SeatID == objID {
				if selectSeatDetails.Bus.Seats[idx].SeatAvailibility {
					if selectSeatDetails.Bus.Seats[idx].IsSelected {
						json.NewEncoder(w).Encode(-selectSeatDetails.Bus.Seats[idx].Cost)
						return
					} else {
						json.NewEncoder(w).Encode(selectSeatDetails.Bus.Seats[idx].Cost)
						return
					}
				} else {
					json.NewEncoder(w).Encode("This seat is booked and not available")
					return
				}
			}
		}
	}
}
