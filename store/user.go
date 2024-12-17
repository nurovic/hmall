package store

import (
	"context"
	"log"

	"github.com/nurovic/hmall/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User koleksiyonu
var userCollection *mongo.Collection

func init() {
	// MongoDB client al ve koleksiyonu başlat
	client := GetMongoClient()
	if client == nil {
		log.Fatal("MongoDB istemcisi başlatılamadı")
	}
	userCollection = client.Database("hmall").Collection("users")
}

func CreateUser(ctx context.Context, user models.User) error {
	_, err := userCollection.InsertOne(ctx, user)
	return err
}

// Kullanıcıyı ID ile getiren fonksiyon
func GetUserByID(ctx context.Context,id string) (*models.User, error) {
	var user models.User
	// Convert the string id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Search using ObjectID
	err = userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	return &user, err
}