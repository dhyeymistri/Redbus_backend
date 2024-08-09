package review

import (
	connection "Redbus_backend/Config"
	Generic "Redbus_backend/Generic"
	models "Redbus_backend/Models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddReview(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//data is array of bytes
		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error with data retrieval")
		}
		asString := string(data)

		var review models.Review
		json.Unmarshal([]byte(asString), &review)

		var params = mux.Vars(r)

		customerID := params["customerID"]
		busID := params["busID"]
		objIDCustomer, _ := primitive.ObjectIDFromHex(customerID)
		objIDBus, _ := primitive.ObjectIDFromHex(busID)
		review.BusID = objIDBus
		review.CustomerID = objIDCustomer

		var ticket models.Ticket
		ticketCollection := connection.ConnectDB("Tickets")
		filter := bson.M{"customerID": objIDCustomer, "busID": objIDBus}
		err = ticketCollection.FindOne(context.TODO(), filter).Decode(&ticket)
		if err != nil {
			json.NewEncoder(w).Encode("You cannot rate a bus that you have not traveled on!")
			return
		}
		currentTime := time.Now().Format("15:04")
		currentDate := time.Now().Format("2006-01-02")
		if currentDate < ticket.DropDate || (currentDate == ticket.DropDate && currentTime < ticket.DropTime) {
			json.NewEncoder(w).Encode("You can rate the bus only after you have completed the journey")
			return
		}

		//also check that a person can only rate it after the journey has finished

		reviewCollection := connection.ConnectDB("Reviews")
		_, err = reviewCollection.InsertOne(context.TODO(), review)
		if err != nil {
			log.Fatal(err)
		}

		var bus models.Bus
		busCollection := connection.ConnectDB("Buses")
		busFilter := bson.M{"_id": objIDBus}
		err = busCollection.FindOne(context.TODO(), busFilter).Decode(&bus)
		if err != nil {
			log.Fatal(err)
		}
		numberOfReviews := bus.NumberOfReviews
		avgRating := bus.AverageRating
		newRating := (float64(numberOfReviews)*avgRating + float64(review.Rating)) / (float64(numberOfReviews) + 1)
		update := bson.M{
			"$set": bson.M{
				"numberOfReviews": numberOfReviews + 1,
				"avgRating":       newRating,
			},
		}
		result, err := busCollection.UpdateByID(context.TODO(), objIDBus, update)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(result)
	}
}

func GetReviewsByBusID(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	var params = mux.Vars(r)
	busID := params["busID"]
	objBusID, _ := primitive.ObjectIDFromHex(busID)

	var reviews []models.Review

	reviewCollection := connection.ConnectDB("Reviews")
	filter := bson.M{"busID": objBusID}
	cursor, err := reviewCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.TODO()) {
		var review models.Review
		err = cursor.Decode(&review)
		if err != nil {
			log.Fatal(err)
		}
		reviews = append(reviews, review)
	}
	json.NewEncoder(w).Encode(reviews)
}
