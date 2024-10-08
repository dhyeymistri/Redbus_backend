package user

import (
	connection "Redbus_backend/Config"
	Generic "Redbus_backend/Generic"
	auth "Redbus_backend/Helpers/Auth"
	helper "Redbus_backend/Helpers/EncryptDecrypt"
	"log"
	"os"
	"time"

	stringgenerator "Redbus_backend/Helpers/RandomStringGenerator"
	models "Redbus_backend/Models"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"

	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ////////////////////////////////////////////////////////////////////////////////////////
// HELPER FUNCTIONS
// ////////////////////////////////////////////////////////////////////////////////////////
func ifEmailExists(email string) string {
	collection := connection.ConnectDB("Users")
	var user models.User
	filter := bson.M{"email": email}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return ""
	}
	return user.ID.Hex()
}

func createTTLIndex(collection *mongo.Collection) error {
	// Create a TTL index on the expireAt field with an expiration time of 10 minutes (600 seconds)
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "expiresAt", Value: 1},
		},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func UpdatePassword(userID string, newPassword string) {
	userCollection := connection.ConnectDB("Users")
	encryptSecret := os.Getenv("ENCRYPTION_KEY")
	var user models.User
	objID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": objID}
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		fmt.Println("User not found!")
	}
	encryptedPass := helper.Encrypt([]byte(newPassword), encryptSecret)
	filter = bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"encryptedPassword": hex.EncodeToString(encryptedPass),
		},
	}
	_, err = userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("Unable to update")
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
// API FUNCTIONS
// ////////////////////////////////////////////////////////////////////////////////////////
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	encryptSecret := os.Getenv("ENCRYPTION_KEY")

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error with data retrieval")
		}
		asString := string(data)

		var user models.User
		json.Unmarshal([]byte(asString), &user)

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			json.NewEncoder(w).Encode("There is some discrepency in input data, please check" + err.Error())
			return
		}

		userEmail := user.Email

		collection := connection.ConnectDB("Users")
		filter := bson.M{
			"email": userEmail,
		}
		count, err := collection.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		if count != 0 {
			json.NewEncoder(w).Encode("This email id is already registered.")
			return
		}

		//Encrypting password
		old_pass := user.EncryptedPassword
		encryptedPass := helper.Encrypt([]byte(old_pass), encryptSecret)
		user.EncryptedPassword = hex.EncodeToString(encryptedPass)
		user.ID = primitive.NewObjectID()

		_, err = collection.InsertOne(context.TODO(), user) //

		if err != nil {
			connection.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode("New user successfully registered")
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		encryptSecret := os.Getenv("ENCRYPTION_KEY")

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		asString := string(data)

		var loginDetails map[string]interface{}
		json.Unmarshal([]byte(asString), &loginDetails)
		email, ok := loginDetails["email"].(string)
		if !ok {
			json.NewEncoder(w).Encode("Email id not found")
			return
		}
		password, ok := loginDetails["password"].(string)
		if !ok {
			json.NewEncoder(w).Encode("Password not found")
			return
		}

		collection := connection.ConnectDB("Users")
		filter := bson.M{"email": email}
		var user models.User
		errr := collection.FindOne(context.TODO(), filter).Decode(&user)
		if errr != nil {
			json.NewEncoder(w).Encode("invalid email")
			return
		}
		hex_pass, _ := hex.DecodeString(user.EncryptedPassword)
		decryptedValue := helper.Decrypt(hex_pass, encryptSecret)
		if password == string(decryptedValue) {
			//the user is verified and therefore a JWT token is created and stored in cookies
			fmt.Println(user.ID)
			token, e := auth.CreateToken(user.ID.Hex(), email, user.Role)
			if e != nil {
				fmt.Println(e)
			}
			cookie := &http.Cookie{
				Name:     "token",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			json.NewEncoder(w).Encode("User logged in")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("incorrect password")
		}
	}
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	collection := connection.ConnectDB("Users")

	var user models.User
	var params = mux.Vars(r)

	id := params["userID"]
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		json.NewEncoder(w).Encode("User not found, check user ID")
		return
	}

	json.NewEncoder(w).Encode(models.User(user))
}

func VerifyEmailAndSendKey(w http.ResponseWriter, r *http.Request) {
	//will input a json with just email
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		asString := string(data)
		var emailDetails map[string]interface{}
		json.Unmarshal([]byte(asString), &emailDetails)

		email, ok := emailDetails["email"].(string)
		if !ok {
			fmt.Println(" this email is not a string")
			return
		}

		userID := ifEmailExists(email)
		if userID != "" {
			randomString := stringgenerator.String(32) //making a random string of length 32
			collection := connection.ConnectDB("ForgotPasswordKeys")

			err := createTTLIndex(collection)
			if err != nil {
				log.Fatalf("Failed to create TTL index: %v", err)
			}

			expirationTime := time.Now().Add(10 * time.Minute)
			detailsToBeSent := models.ForgotPasswordKeyDetails{
				UserID:    userID,
				Key:       randomString,
				CreatedAt: time.Now(),
				ExpiresAt: expirationTime,
			}
			_, err = collection.InsertOne(context.Background(), detailsToBeSent)
			if err != nil {
				log.Fatalf("Failed to insert document: %v", err)
			}
			json.NewEncoder(w).Encode(detailsToBeSent)
			//this userID and Key will be sent in the link via email to the user on which he will click
			//upon clicking he will have to enter password and it will be confirmed in the frontend
			//this password, userID and Key will be sent to backend upon submission
			//where the key will be checked
			//if the key is present, the password will change and the key will be deleted
		} else {
			json.NewEncoder(w).Encode("This email does not exist!")
		}
	}
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//body will have userID, key and new password
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		asString := string(data)
		var updateDetails map[string]interface{}
		json.Unmarshal([]byte(asString), &updateDetails)
		userID, ok := updateDetails["userID"].(string)
		if !ok {
			fmt.Println("userID is not a string")
			return
		}
		newPassword, ok := updateDetails["newPassword"].(string)
		if !ok {
			fmt.Println("newPassword is not a string")
			return
		}
		key, ok := updateDetails["key"].(string)
		if !ok {
			fmt.Println("key is not a string")
			return
		}

		collection := connection.ConnectDB("ForgotPasswordKeys")
		var forgotPasswordDetail models.ForgotPasswordKeyDetails
		filter := bson.M{"userID": userID}
		err = collection.FindOne(context.TODO(), filter).Decode(&forgotPasswordDetail)
		if err != nil {
			json.NewEncoder(w).Encode("The URL has expired! or the userID is wrong")
			return
		}
		if forgotPasswordDetail.Key == key {
			UpdatePassword(userID, newPassword)
			_, err = collection.DeleteOne(context.TODO(), filter)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode("Password updated")
		} else {
			json.NewEncoder(w).Encode("Your key does not match")
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		// SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode("user logged out")
}
