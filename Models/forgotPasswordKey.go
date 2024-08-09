package models

import "time"

type ForgotPasswordKeyDetails struct {
	UserID    string    `bson:"userID"`
	Key       string    `bson:"key"`
	CreatedAt time.Time `bson:"createdAt"`
	ExpiresAt time.Time `bson:"expiresAt"`
}
