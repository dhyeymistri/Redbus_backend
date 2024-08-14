package getbusdetail

import (
	connection "Redbus_backend/Config"
	models "Redbus_backend/Models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
