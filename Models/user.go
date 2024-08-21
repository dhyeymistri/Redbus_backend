package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName" validate:"required"`
	LastName          string             `bson:"lastName" json:"lastName" validate:"required"`
	Age               int                `bson:"age" json:"age" validate:"required"`
	DOB               string             `bson:"dob" json:"dob" validate:"required"`
	Email             string             `bson:"email" json:"email" validate:"required,email"`
	Gender            string             `bson:"gender" json:"gender" validate:"required"`
	Role              string             `bson:"role" json:"role" validate:"required"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"encryptedPassword" validate:"required"`
	ProfilePicPath    string             `bson:"profilePicPath,omniempty" json:"profilePicPath"`
	WalletBalance     int                `bson:"walletBalance" json:"walletBalance"`
}
