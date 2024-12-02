package store

import (
	"context"
	"errors"
	"log"

	"github.com/nurovic/hmall/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Product koleksiyonu
var productCollection *mongo.Collection

func init() {
	// MongoDB client al ve koleksiyonu başlat
	client := GetMongoClient()
	if client == nil {
		log.Fatal("MongoDB istemcisi başlatılamadı")
	}
	productCollection = client.Database("hmall").Collection("products")
}

// Ürünü MongoDB'ye ekleyen fonksiyon
func CreateProduct(product models.Product) error {
	_, err := productCollection.InsertOne(context.Background(), product)
	return err
}

// Ürünü ID ile getiren fonksiyon
func GetProductByID(id string) (*models.Product, error) {
	// ID'yi ObjectId'ye dönüştür
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("geçersiz ID formatı")
	}
	var product models.Product
	err = productCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ürün bulunamadı")
		}
		return nil, err
	}

	return &product, nil
}