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

var productCollection *mongo.Collection

func init() {
	client := GetMongoClient()
	if client == nil {
		log.Fatal("MongoDB istemcisi başlatılamadı")
	}
	productCollection = client.Database("hmall").Collection("products")
}

func CreateProduct(ctx context.Context, product models.Product) error {
	_, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		return errors.New("ürün eklenirken bir hata oluştu")
	}

	return nil
}

func GetProductByID(ctx context.Context, id string) (*models.Product, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("geçersiz ID formatı")
	}

	var product models.Product
	err = productCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ürün bulunamadı")
		}
		return nil, errors.New("ürün getirilirken bir hata oluştu")
	}

	return &product, nil
}
