package booking

import (
	connection "Redbus_backend/Config"
	Generic "Redbus_backend/Generic"
	auth "Redbus_backend/Helpers/Auth"
	getbusdetail "Redbus_backend/Helpers/GetBusDetail"
	minorhelpers "Redbus_backend/Helpers/SmallFunctionalities"
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

func BookSeat(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error with data retrieval")
		}
		asString := string(data)
		var bookingDetails models.BookSeatDoc

		json.Unmarshal([]byte(asString), &bookingDetails)
		booking := bookingDetails.Booking
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

		bookingCollection := connection.ConnectDB("Bookings")
		bookingsFilter := bson.M{"busID": busID, "travelStartDate": travelStartDate, "travelEndDate": travelEndDate}

		userCollection := connection.ConnectDB("Users")
		filter := bson.M{"_id": objUserID}
		totalPayableAmount := bookingDetails.TotalPayableAmount
		if bookingDetails.DiscountedAmount == 0 {
			bookingDetails.GST = int(float64(bookingDetails.BaseFare*5)) / 100
			totalPayableAmount = int(float64(bookingDetails.BaseFare*105)) / 100
		}
		var user models.User
		err = userCollection.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		if user.WalletBalance < totalPayableAmount {
			json.NewEncoder(w).Encode("You don't have sufficient balance to book the tickets")
			return
		}
		bus := getbusdetail.GetBusDetail(busID)

		var foundBooking models.Booking
		rowSlice := []int{}
		colSlice := []int{}
		deckSlice := []bool{}
		err = bookingCollection.FindOne(context.TODO(), bookingsFilter).Decode(&foundBooking)
		if err != nil {
			seats := bus.Seats
			count := 0
			for idx := range seats {
				for outerIndex := range arrSeatID {
					if seats[idx].SeatID == arrSeatID[outerIndex] {
						if arrPassengerGender[outerIndex] == "F" && seats[idx].OnlyMale {
							json.NewEncoder(w).Encode("This is an only male seat and female passenger cannot book this for safety purpose")
							return
						}
						if arrPassengerGender[outerIndex] == "M" && seats[idx].OnlyFemale {
							json.NewEncoder(w).Encode("This is an only female seat and a male passenger cannot book this for safety purpose")
							return
						}
						seats[idx].PassengerAge = int(arrPassengerAge[outerIndex])
						seats[idx].PassengerGender = arrPassengerGender[outerIndex]
						seats[idx].PassengerName = arrPassengerName[outerIndex]
						seats[idx].SeatAvailibility = false
						rowSlice = append(rowSlice, seats[idx].Row)
						colSlice = append(colSlice, seats[idx].Column)
						deckSlice = append(deckSlice, seats[idx].IsUpperDeck)
						if arrPassengerGender[outerIndex] == "F" {
							seats[idx].OnlyFemale = true
							seats[idx].OnlyMale = false
							// groupNo = booking.Bus.Seats[idx].Group
						} else {
							seats[idx].OnlyFemale = false
							seats[idx].OnlyMale = true
							// groupNo = booking.Bus.Seats[idx].Group
						}
					}
				}
				if seats[idx].SeatAvailibility {
					count++
				}
			}

			booking.AvailableSeats = count

			//gender verification-------------------------------------------------
			for idx := range rowSlice {
				row := rowSlice[idx]
				col := colSlice[idx]
				deck := deckSlice[idx]
				for seatIndex := range seats {
					if (seats[seatIndex].Row == row-1 || seats[seatIndex].Row == row+1) && seats[seatIndex].Column == col && seats[seatIndex].IsUpperDeck == deck && seats[seatIndex].SeatAvailibility {
						if arrPassengerGender[idx] == "F" {
							seats[seatIndex].OnlyFemale = true
						} else {
							seats[seatIndex].OnlyMale = true
						}
					}
				}
			}

			update := bson.M{"$set": bson.M{"availableSeats": booking.AvailableSeats, "travelStartTime": booking.TravelStartTime, "travelEndTime": booking.TravelEndTime, "travelStartLocation": booking.TravelStartLocation, "travelEndLocation": booking.TravelEndLocation}}

			// Update the BookSeatDoc if it exists, otherwise insert a new BookSeatDoc
			opts := options.Update().SetUpsert(true)
			result, err := bookingCollection.UpdateOne(context.TODO(), bookingsFilter, update, opts)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result.UpsertedID)

			timeDiff, _ := minorhelpers.TimeDifference(booking.TravelStartTime, booking.TravelEndTime)
			timeDiffInMinutes := timeDiff.Minutes()
			seats = minorhelpers.AllotSeatPrices(bus, timeDiffInMinutes)
			bookedSeats := models.BookedSeats{
				BookingID: result.UpsertedID.(primitive.ObjectID),
				Seats:     seats,
			}
			filter = bson.M{"bookingID": bookedSeats.BookingID}

			seatsCollection := connection.ConnectDB("BookedSeats")
			update = bson.M{
				"$set": bson.M{
					"seats": bookedSeats.Seats,
				},
			}
			opts = options.Update().SetUpsert(true)
			_, err = seatsCollection.UpdateOne(context.TODO(), filter, update, opts)
			if err != nil {
				log.Fatal(err)
			}

			//ticket------------------------------------------------
			ticket := minorhelpers.MakeTicket(user, booking, bus, bookingDetails)
			ticketCollection := connection.ConnectDB("Tickets")
			insertedResult, err := ticketCollection.InsertOne(context.TODO(), ticket)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(insertedResult.InsertedID)
		} else {
			bookingID := foundBooking.ID
			filter = bson.M{"bookingID": bookingID}
			seatsCollection := connection.ConnectDB("BookedSeats")
			var bookedSeats models.BookedSeats
			err = seatsCollection.FindOne(context.TODO(), filter).Decode(&bookedSeats)
			if err != nil {
				log.Fatal(err)
			}
			seats := bookedSeats.Seats
			count := 0
			for idx := range seats {
				for outerIndex := range arrSeatID {
					if seats[idx].SeatID == arrSeatID[outerIndex] {
						if !seats[idx].SeatAvailibility {
							json.NewEncoder(w).Encode("This seat has been booked already")
							return
						}
						if arrPassengerGender[outerIndex] == "F" && seats[idx].OnlyMale {
							json.NewEncoder(w).Encode("This is an only male seat and female passenger cannot book this for safety purpose")
							return
						}
						if arrPassengerGender[outerIndex] == "M" && seats[idx].OnlyFemale {
							json.NewEncoder(w).Encode("This is an only female seat and a male passenger cannot book this for safety purpose")
							return
						}
						seats[idx].PassengerAge = int(arrPassengerAge[outerIndex])
						seats[idx].PassengerGender = arrPassengerGender[outerIndex]
						seats[idx].PassengerName = arrPassengerName[outerIndex]
						seats[idx].SeatAvailibility = false
						rowSlice = append(rowSlice, seats[idx].Row)
						colSlice = append(colSlice, seats[idx].Column)
						deckSlice = append(deckSlice, seats[idx].IsUpperDeck)
					}
				}
				if seats[idx].SeatAvailibility {
					count++
				}
			}

			foundBooking.AvailableSeats = count

			for idx := range rowSlice {
				row := rowSlice[idx]
				col := colSlice[idx]
				deck := deckSlice[idx]
				for seatIndex := range seats {
					if (seats[seatIndex].Row == row-1 || seats[seatIndex].Row == row+1) && seats[seatIndex].Column == col && seats[seatIndex].IsUpperDeck == deck && seats[seatIndex].SeatAvailibility {
						if arrPassengerGender[idx] == "F" {
							seats[seatIndex].OnlyFemale = true
						} else {
							seats[seatIndex].OnlyMale = true
						}
					}
				}
			}

			update := bson.M{"$set": bson.M{"availableSeats": foundBooking.AvailableSeats}}

			// Update the BookSeatDoc if it exists, otherwise insert a new BookSeatDoc
			opts := options.Update().SetUpsert(true)
			_, err := bookingCollection.UpdateOne(context.TODO(), bookingsFilter, update, opts)
			if err != nil {
				log.Fatal(err)
			}

			update = bson.M{
				"$set": bson.M{
					"seats": seats,
				},
			}
			opts = options.Update().SetUpsert(true)
			_, err = seatsCollection.UpdateOne(context.TODO(), filter, update, opts)
			if err != nil {
				log.Fatal(err)
			}

			//function that assigns values to tickets
			ticket := minorhelpers.MakeTicket(user, foundBooking, bus, bookingDetails)
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

// made a handler to get the seat data from either the booking, if present in the collection
// or from the bus data itself to reduce redundant code in this file

func ViewSeats(w http.ResponseWriter, r *http.Request) {
	handlerData, ok := r.Context().Value("handlerData").(models.HandlerData)
	if !ok {
		http.Error(w, "Error retrieving handler data", http.StatusInternalServerError)
		return
	}
	if handlerData.Error != nil {
		http.Error(w, handlerData.Error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(handlerData.Seats)
}

func SelectSeat(w http.ResponseWriter, r *http.Request) {
	handlerData, ok := r.Context().Value("handlerData").(models.HandlerData)
	if !ok {
		http.Error(w, "Error retrieving handler data", http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)
	seatID := params["seatID"]
	objSeatID, _ := primitive.ObjectIDFromHex(seatID)

	for idx := range handlerData.Seats {
		if handlerData.Seats[idx].SeatID == objSeatID {
			if handlerData.Seats[idx].SeatAvailibility {
				if handlerData.Seats[idx].IsSelected {
					json.NewEncoder(w).Encode(-handlerData.Seats[idx].Cost)
					return
				} else {
					json.NewEncoder(w).Encode(handlerData.Seats[idx].Cost)
					return
				}
			} else {
				json.NewEncoder(w).Encode("This seat is booked and not available")
				return
			}
		}
	}
	http.Error(w, "Seat not found", http.StatusNotFound)
}
