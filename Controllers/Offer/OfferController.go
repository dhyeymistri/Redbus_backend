package offer

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
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func AddOffer(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//data is array of bytes
		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error with data retrieval")
		}
		asString := string(data)
		var offerDetails models.Offer
		json.Unmarshal([]byte(asString), &offerDetails)

		offerCollection := connection.ConnectDB("Offers")
		filter := bson.M{
			"oCode": offerDetails.OfferCode,
		}
		count, err := offerCollection.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		if count != 0 {
			json.NewEncoder(w).Encode("This offer code already exists")
			return
		}

		timeTime, _ := time.Parse("2006-01-02", offerDetails.Validity)
		timeTime = timeTime.Add(23 * time.Hour)
		timeTime = timeTime.Add(59 * time.Minute)
		timeTime = timeTime.Add(59 * time.Second)
		offerDetails.Validity = timeTime.Format("2006-01-02 15:04:05")

		_, err = offerCollection.InsertOne(context.TODO(), offerDetails)
		if err != nil {
			log.Fatal("Unable to add new offer", err.Error())
		}
		json.NewEncoder(w).Encode("Offer added. Offer code: " + offerDetails.OfferCode)
	}
}

func GetOffers(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	offerCollection := connection.ConnectDB("Offers")

	//getting current time and likewise filtering
	currentTime := time.Now()
	strCurrentTime := currentTime.Format("2006-01-02 15:04:05")

	filter := bson.M{"validity": bson.M{"$gt": strCurrentTime}}

	cursor, err := offerCollection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(context.TODO())

	var res []models.Offer
	if err = cursor.All(context.Background(), &res); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(res)
}

func ApplyOffer(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//data is array of bytes
		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error with data retrieval")
		}
		asString := string(data)
		var offerDetails map[string]interface{}
		// rr := strings.NewReader(asString)

		json.Unmarshal([]byte(asString), &offerDetails)

		offerCollection := connection.ConnectDB("Offers")
		baseFare := int(offerDetails["cartValue"].(float64))
		offerCode := offerDetails["offerCode"].(string)
		offerFilter := bson.M{"oCode": offerCode}

		var offer models.Offer
		err = offerCollection.FindOne(context.TODO(), offerFilter).Decode(&offer)
		if err != nil {
			json.NewEncoder(w).Encode("This offer does not exist")
			return
			// log.Fatal(err)
		}

		currentTime := time.Now()
		validityTime, _ := time.Parse("2006-01-02 15:04:05", offer.Validity)
		validityDate := validityTime.Format("2006-01-02")
		if currentTime.After(validityTime) {
			json.NewEncoder(w).Encode("This offer is not valid anymore. Validity expired at " + validityDate)
			json.NewEncoder(w).Encode(map[string]interface{}{"baseFare": baseFare, "totalPayableAmount": (baseFare * 105 / 100), "discountedAmount": 0, "gst": (baseFare * 5 / 100)})
			return
		}
		if baseFare < offer.MinOrderVal {
			json.NewEncoder(w).Encode("Cart value is below the minimum cart value for this offer, i.e. Rs." + strconv.Itoa(offer.MinOrderVal))
			json.NewEncoder(w).Encode(map[string]interface{}{"baseFare": baseFare, "totalPayableAmount": (baseFare * 105 / 100), "discountedAmount": 0, "gst": (baseFare * 5 / 100)})
			return
		}
		gst := baseFare * 5 / 100
		gstIncluded := gst + baseFare
		discountedAmount := min(offer.MaxDiscount, (gstIncluded * offer.Discount / 100))

		json.NewEncoder(w).Encode(map[string]interface{}{"baseFare": baseFare, "totalPayableAmount": gstIncluded - discountedAmount, "discountedAmount": discountedAmount, "gst": gst})
	}
}
