package bus

import (
	connection "Redbus_backend/Config"
	Generic "Redbus_backend/Generic"
	filtersearchedbuses "Redbus_backend/Helpers/FilterSearchedBuses"
	getbusdetail "Redbus_backend/Helpers/GetBusDetail"
	seating "Redbus_backend/Helpers/SeatingArrangement"
	"strconv"

	minorhelpers "Redbus_backend/Helpers/SmallFunctionalities"
	models "Redbus_backend/Models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	// "time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddBus(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//data is array of bytes
		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error with data retrieval")
		}
		asString := string(data)

		var bus models.Bus
		json.Unmarshal([]byte(asString), &bus)

		seatingArrangement := bus.BusType
		bus.Seats = seating.ArrangingSeats(seatingArrangement)
		bus.TotalSeats = len(bus.Seats)
		bus.AvailableSeats = bus.TotalSeats

		collection := connection.ConnectDB("Buses")
		result, err := collection.InsertOne(context.TODO(), bus)
		if err != nil {
			fmt.Println("Error in inserting bus")
		}

		newBusID := result.InsertedID.(primitive.ObjectID)

		locationCollection := connection.ConnectDB("Locations")

		for _, location := range bus.Stops {
			// arrivalDate := ""
			arrivalTime := ""
			// departureDate := ""
			departureTime := ""
			if location.ArrivalTime != "" {
				arrival, _ := time.Parse("15:04", location.ArrivalTime)
				arrivalTime = arrival.Format("15:04")
			}
			if location.DepartureTime != "" {
				departure, _ := time.Parse("15:04", location.DepartureTime)
				departureTime = departure.Format("15:04")
			}
			weekendBool := false
			if bus.Frequency == "Weekends" {
				weekendBool = true
			}
			filter := bson.M{"location": location.Location}
			update := bson.M{
				"$push": bson.M{
					"buses":            newBusID,
					"seats":            [][]primitive.ObjectID{},
					"arrivalTimings":   arrivalTime,
					"departureTimings": departureTime,
					"isWeekend":        weekendBool,
				},
			}
			upsert := true

			// Perform the update operation
			_, err = locationCollection.UpdateOne(
				context.TODO(),
				filter,
				update,
				&options.UpdateOptions{Upsert: &upsert},
			)
			if err != nil {
				log.Fatalf("Failed to update or insert document: %v", err)
			}
		}

		json.NewEncoder(w).Encode(bus)
	}
}

func GetSearchedBus(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	page, _ := strconv.Atoi(params["page"])
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	asString := string(data)
	var searchDetails map[string]interface{}
	json.Unmarshal([]byte(asString), &searchDetails)
	startDestination, _ := searchDetails["fromLocation"].(string)
	finalDestination := searchDetails["toLocation"].(string)
	travelDate := searchDetails["travelDate"].(string)
	filters := map[string]interface{}{}
	value, exists := searchDetails["filters"]
	if exists {
		filters = value.(map[string]interface{})
	} else {
		fmt.Print("No filtering")
	}

	dateForm, _ := time.Parse("2006-01-02", travelDate)

	currentDate := time.Now().Format("2006-01-02")
	if currentDate > travelDate {
		json.NewEncoder(w).Encode("Choose a future date to travel")
		return
	}

	locationCollection := connection.ConnectDB("Locations")

	var startOfRoute models.Route
	var endOfRoute models.Route

	startFilter := bson.M{"location": startDestination}
	endFilter := bson.M{"location": finalDestination}

	errr := locationCollection.FindOne(context.TODO(), startFilter).Decode(&startOfRoute)
	if errr != nil {
		json.NewEncoder(w).Encode("Redbus does not serve in " + startDestination)
		return
	}
	errr = locationCollection.FindOne(context.TODO(), endFilter).Decode(&endOfRoute)
	if errr != nil {
		json.NewEncoder(w).Encode("Redbus does not serve in " + finalDestination)
		return
	}
	// commonBuses := minorhelpers.FindCommonElements(startOfRoute.Buses, endOfRoute.Buses)
	var weekendFilteredBusID []primitive.ObjectID
	var weekendFilteredStartTime []string
	// var weekendFilteredEndTime []string
	fmt.Println(len(endOfRoute.Buses))
	// currentTime := time.Now().Format("15:04")
	for index, checkBool := range startOfRoute.IsWeekend {
		if dateForm.Weekday().String() == "Sunday" || dateForm.Weekday().String() == "Saturday" {
			weekendFilteredStartTime = append(weekendFilteredStartTime, startOfRoute.DepartureTime[index])
			weekendFilteredBusID = append(weekendFilteredBusID, startOfRoute.Buses[index])
		} else {
			if !checkBool {
				weekendFilteredStartTime = append(weekendFilteredStartTime, startOfRoute.DepartureTime[index])
				weekendFilteredBusID = append(weekendFilteredBusID, startOfRoute.Buses[index])
			}
		}
	}

	// var routeFilteredBusID []primitive.ObjectID
	var searchedBusResult []models.Booking
	for index, busID := range weekendFilteredBusID {
		bus := getbusdetail.GetBusDetail(busID)
		busStops := bus.Stops
		flag := false
		for _, obj := range busStops {
			if obj.Location == finalDestination {
				if flag {
					//only then is the result shown
					//checking if any booking is already present for this bus
					//on that day for the same destination
					//if present, we will show it from that bookings so that we have updated seats
					//if not present, we will make a temporary booking that will be shown with total seats equal to available seats
					bookingCollection := connection.ConnectDB("Bookings")
					bookingsFilter := bson.M{"busID": busID, "travelStartDate": travelDate, "travelStartLocation": startDestination}
					var booking models.Booking
					e := bookingCollection.FindOne(context.TODO(), bookingsFilter).Decode(&booking)
					if e != nil {
						booking.BusID = bus.ID
						booking.TravelStartDate = travelDate
						booking.TravelStartLocation = startDestination
						booking.TravelEndLocation = finalDestination
						booking.TravelStartTime = weekendFilteredStartTime[index]
						booking.TravelEndTime = obj.ArrivalTime

						if minorhelpers.HasDateChanged(weekendFilteredStartTime[index], obj.ArrivalTime) {
							booking.TravelEndDate = dateForm.Add(time.Hour * 24).Format("2006-01-02")
						} else {
							booking.TravelEndDate = travelDate
						}
						searchedBusResult = append(searchedBusResult, booking)
						fmt.Println("New booking added")
					} else {
						searchedBusResult = append(searchedBusResult, booking)
						fmt.Println("booking present already")
					}
				}
				break
			}
			if obj.Location == startDestination {
				flag = true
			}
		}
	}
	//filtering ------- make a function
	if len(filters) != 0 {
		searchedBusResult = filtersearchedbuses.Filtering(filters, searchedBusResult)
	}

	//pagination
	startIndex := (page) * 10
	var paginatedResult []models.Booking
	if startIndex < len(searchedBusResult) {
		endIndex := startIndex + 10
		if endIndex > len(searchedBusResult) {
			endIndex = len(searchedBusResult)
		}
		paginatedResult = searchedBusResult[startIndex:endIndex]
	} else {
		json.NewEncoder(w).Encode("No results found")
		return
	}
	json.NewEncoder(w).Encode(paginatedResult)
}

func GetBusDetail(ID primitive.ObjectID) models.Bus {
	collection := connection.ConnectDB("Buses")
	var bus models.Bus
	filter := bson.M{"_id": ID}
	err := collection.FindOne(context.TODO(), filter).Decode(&bus)

	if err != nil {
		fmt.Println("Error retreiving bus")
	}
	return bus
}

func GetBusByID(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	collection := connection.ConnectDB("Buses")

	var bus models.Bus
	var params = mux.Vars(r)

	id := params["busID"]
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	err := collection.FindOne(context.TODO(), filter).Decode(&bus)

	if err != nil {
		connection.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(models.Bus(bus))
}
