package store

import (
	"context"

	"github.com/nurovic/hmall/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection

func init() {
	// MongoDB client al ve koleksiyonu başlat
	client := GetMongoClient()
	orderCollection = client.Database("hmall").Collection("orders")
}

// Siparişi MongoDB'ye ekleyen fonksiyon
func CreateOrder(order models.Order) error {
	_, err := orderCollection.InsertOne(context.Background(), order)
	return err
}

// Siparişi ID ile getiren fonksiyon
func GetOrderByID(id string) (*models.Order, error) {
	var order models.Order
	err := orderCollection.FindOne(context.Background(), bson.M{"id": id}).Decode(&order)
	return &order, err
}
