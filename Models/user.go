package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	FirstName         string             `bson:"firstName"`
	LastName          string             `bson:"lastName"`
	Age               int                `bson:"age"`
	DOB               string             `bson:"dob"`
	Email             string             `bson:"email"`
	Gender            string             `bson:"gender"`
	Role              string             `bson:"role"`
	EncryptedPassword string             `bson:"encryptedPassword"`
	ProfilePicPath    string             `bson:"profilePicPath"`
	WalletBalance     int                `bson:"walletBalance" json:"walletBalance"`
}
