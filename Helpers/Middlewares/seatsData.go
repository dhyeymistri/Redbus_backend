package middlewares

import (
	connection "Redbus_backend/Config"
	Generic "Redbus_backend/Generic"
	getbusdetail "Redbus_backend/Helpers/GetBusDetail"
	minorhelpers "Redbus_backend/Helpers/SmallFunctionalities"
	models "Redbus_backend/Models"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Middleware to fetch common data and pass it to handlers
func FetchCommonData(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Generic.SetupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")

		// Read request body
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		asString := string(data)

		// Parse JSON body
		var searchDetails map[string]interface{}
		json.Unmarshal([]byte(asString), &searchDetails)
		startDestination, _ := searchDetails["fromLocation"].(string)
		finalDestination := searchDetails["toLocation"].(string)
		travelDate := searchDetails["travelDate"].(string)

		params := mux.Vars(r)
		strBusID := params["busID"]
		busID, _ := primitive.ObjectIDFromHex(strBusID)

		filter := bson.M{
			"travelStartDate":     travelDate,
			"busID":               busID,
			"travelStartLocation": startDestination,
		}

		var booking models.Booking
		bookingCollection := connection.ConnectDB("Bookings")
		err = bookingCollection.FindOne(context.TODO(), filter).Decode(&booking)

		var handlerData models.HandlerData

		if err != nil {
			// Fetch bus details if no booking is found
			bus := getbusdetail.GetBusDetail(busID)
			busStops := bus.Stops
			var travelStartTime, travelEndTime string
			for idx := range busStops {
				if busStops[idx].Location == startDestination {
					travelStartTime = busStops[idx].DepartureTime
				} else if busStops[idx].Location == finalDestination {
					travelEndTime = busStops[idx].DepartureTime
				}
			}
			timeDiff, _ := minorhelpers.TimeDifference(travelStartTime, travelEndTime)
			timeDiffInMinutes := timeDiff.Minutes()
			seats := minorhelpers.AllotSeatPrices(bus, timeDiffInMinutes)
			handlerData.Seats = seats
			handlerData.Bus = bus
		} else {
			// Fetch booked seats if booking is found
			bookingID := booking.ID
			filter = bson.M{"bookingID": bookingID}
			seatCollection := connection.ConnectDB("BookedSeats")
			var bookedSeats models.BookedSeats
			err = seatCollection.FindOne(context.TODO(), filter).Decode(&bookedSeats)
			if err != nil {
				http.Error(w, "Error retrieving booked seats", http.StatusInternalServerError)
				return
			}
			handlerData.Seats = bookedSeats.Seats
		}
		ctx := context.WithValue(r.Context(), "handlerData", handlerData)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
