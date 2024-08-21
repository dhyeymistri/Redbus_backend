package wallet

import (
	connection "Redbus_backend/Config"
	Generic "Redbus_backend/Generic"
	models "Redbus_backend/Models"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddToWallet(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		asString := string(data)
		var walletDetails map[string]interface{}

		json.Unmarshal([]byte(asString), &walletDetails)
		addedMoney := walletDetails["moneyToAdd"].(float64)
		if int(addedMoney) > 50000 {
			json.NewEncoder(w).Encode("You cannot add more than Rs.50000 at a time")
			return
		}
		var user models.User
		var params = mux.Vars(r)

		id := params["userID"]
		objID, _ := primitive.ObjectIDFromHex(id)
		filter := bson.M{"_id": objID}
		userCollection := connection.ConnectDB("Users")
		err = userCollection.FindOne(context.TODO(), filter).Decode(&user)

		if err != nil {
			connection.GetError(err, w)
			return
		}
		if user.WalletBalance+int(addedMoney) > 200000 {
			json.NewEncoder(w).Encode("Wallet balance cannot be more than Rs.2,00,000. Available balance: Rs." + strconv.Itoa(user.WalletBalance))
			return
		}

		update := bson.M{"$inc": bson.M{"walletBalance": int(addedMoney)}}
		_, err = userCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode("Rs. " + strconv.Itoa(int(addedMoney)) + " added to your wallet")
	}
}

func WithdrawFromWallet(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		asString := string(data)
		var walletDetails map[string]interface{}

		json.Unmarshal([]byte(asString), &walletDetails)
		toWithdrawMoney := walletDetails["moneyToWithdraw"].(float64)

		var user models.User
		var params = mux.Vars(r)
		id := params["userID"]
		objID, _ := primitive.ObjectIDFromHex(id)
		filter := bson.M{"_id": objID}
		userCollection := connection.ConnectDB("Users")
		err = userCollection.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			connection.GetError(err, w)
			return
		}
		if user.WalletBalance < int(toWithdrawMoney) {
			json.NewEncoder(w).Encode("Insufficient balance. Available balance: Rs." + strconv.Itoa(user.WalletBalance))
			return
		}
		update := bson.M{"$inc": bson.M{"walletBalance": -int(toWithdrawMoney)}}
		_, err = userCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode("Money withdrawn. Available balance:" + strconv.Itoa(user.WalletBalance-int(toWithdrawMoney)))
	}
}

func GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userID := params["userID"]
	objUserID, _ := primitive.ObjectIDFromHex(userID)

	userCollection := connection.ConnectDB("Users")

	var user models.User
	filter := bson.M{"_id": objUserID}
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode("You have Rs. " + strconv.Itoa(user.WalletBalance) + " in your wallet")
}
