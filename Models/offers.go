package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Offer struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OfferCode   string             `json:"oCode" bson:"oCode"`
	Description string             `json:"description" bson:"description"`
	Validity    string             `json:"validity" bson:"validity"`
	MinOrderVal int                `json:"minOrderVal" bson:"minOrderVal"`
	MaxDiscount int                `json:"maxDiscount" bson:"maxDiscount"`
	Discount    int                `json:"discount" bson:"discount"`
}
