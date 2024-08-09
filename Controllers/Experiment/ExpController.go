package exp

import (
	// connection "Redbus_backend/Config"

	Generic "Redbus_backend/Generic"

	// "context"
	// "log"
	// "time"

	"net/http"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// type Document struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty"`
// 	Name      string             `bson:"name"`
// 	PackageID int                `bson:"packageId"`
// 	Age       int                `bson:"age"`
// 	CreatedAt time.Time          `bson:"createdAt"`
// 	ExpiresAt time.Time          `bson:"expiresAt"`
// }

// func createTTLIndex(collection *mongo.Collection) error {
// 	// Create a TTL index on the CreatedAt field with an expiration time of 10 minutes (600 seconds)
// 	indexModel := mongo.IndexModel{
// 		Keys: bson.D{
// 			{Key: "createdAt", Value: 1},
// 		},
// 		Options: options.Index().SetExpireAfterSeconds(0),
// 	}

// 	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
// 	return err
// }

func TempDoc(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	if r.Method == "POST" {
		// collection := connection.ConnectDB("ForgotPasswordKeys")

		// err := createTTLIndex(collection)
		// if err != nil {
		// 	log.Fatalf("Failed to create TTL index: %v", err)
		// }

		// expirationTime := time.Now().Add(60 * time.Second)

		// doc := Document{
		// 	Name:      "Examp",
		// 	PackageID: 123,
		// 	Age:       25,
		// 	CreatedAt: time.Now(),
		// 	ExpiresAt: expirationTime,
		// }

		// _, err = collection.InsertOne(context.Background(), doc)
		// if err != nil {
		// 	log.Fatalf("Failed to insert document: %v", err)
		// }

		// log.Println("Document inserted successfully")

		//--------------------------------------------------------------------------
		//getting role--------------------------------------------------------------
		//--------------------------------------------------------------------------

		// returnedString := auth.GetRole("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImRoeWV5QGdtYWlsLmNvbSIsImV4cCI6MTcyMjMxODE1NywiaWQiOiI2NmEwZTRhMWNkMDBhYjNkYTViYTc3NzEiLCJyb2xlIjoiYWRtaW4ifQ.f9OzGl93_hBv-JfTGzo4pt9f9ddSYKymXcaQLD86cFM")
		// json.NewEncoder(w).Encode(returnedString)

		// json.NewEncoder(w).Encode("You are admin")

		//--------------------------------------------------------------------------
		//checking the seating arrangements-----------------------------------------
		//--------------------------------------------------------------------------
		// var arrangement []models.Seat
		// arrangement := seating.ArrangingSeats("8SE33SL")
		// json.NewEncoder(w).Encode(arrangement)

		//--------------------------------------------------------------------------
		//how I will send dates as string-------------------------------------------
		//--------------------------------------------------------------------------

		// timeNow := time.Now()
		// timeStr := timeNow.Format("15:04")
		// newTime, _ := time.Parse("15:04", timeStr)
		// fmt.Println(newTime.Format("15:04"))
		// json.NewEncoder(w).Encode(newTime.Day()," ",newTime.Month())

		//--------------------------------------------------------------------------
		//testing the change in dates based on times--------------------------------
		//--------------------------------------------------------------------------
		// fmt.Println(minorhelpers.HasDateChanged("2:00", "23:00"))

		//--------------------------------------------------------------------------
		//Starting with bookings----------------------------------------------------
		//--------------------------------------------------------------------------

		//data is array of bytes
		// data, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	fmt.Println("Error with data retrieval")
		// }
		// asString := string(data)

		// var bus models.Bus
		// json.Unmarshal([]byte(asString), &bus)

		// seatingArrangement := bus.BusType
		// bus.Seats = seating.ArrangingSeats(seatingArrangement)
		// bus.TotalSeats = len(bus.Seats)
		// bus.AvailableSeats = bus.TotalSeats

		// var booking models.Booking
		// booking.Bus = bus
		// booking.TravelEndDate = "2024-08-03"
		// booking.TravelStartDate = "2024-08-02"

		// collection := connection.ConnectDB("Bookings")
		// result, err := collection.InsertOne(context.TODO(), booking)
		// if err != nil {
		// 	fmt.Println("Error in inserting bus")
		// }
		// json.NewEncoder(w).Encode(result)

		// currentTime := time.Now()
		// busStartTimeStr := "2024-08-09" + " " + "13:01"
		// busStartTime, _ := time.Parse("2006-01-02 15:04", busStartTimeStr)
		// json.NewEncoder(w).Encode(busStartTime)

		// refund := 300
		// userCollection := connection.ConnectDB("Users")
		// update := bson.M{
		// 	"$inc": bson.M{
		// 		"walletBalance": -refund,
		// 	},
		// }
		// objID, _ := primitive.ObjectIDFromHex("66a0e4a1cd00ab3da5ba7771")
		// result, err := userCollection.UpdateByID(context.TODO(), objID, update)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// json.NewEncoder(w).Encode(result)

	}
}
