package store

import (
	"context"
	"errors"
	"log"

	"github.com/nurovic/hmall/models"
	"github.com/nurovic/hmall/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func init() {
	client := GetMongoClient()
	if client == nil {
		log.Fatal("MongoDB istemcisi başlatılamadı")
	}
	userCollection = client.Database("hmall").Collection("users")
}

func CreateUser(ctx context.Context, user models.User) error {
	hashPwd, _ := pkg.HashPassword(user.Password)
	user.Password = string(hashPwd)
	_, err := userCollection.InsertOne(ctx, user)
	return err
}

func GetUserByID(ctx context.Context,id string) (*models.User, error) {
	var user models.User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	return &user, err
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("kullanıcı bulunamadı")
		}
		return nil, errors.New("veritabanı hatası")
	}

	return &user, nil
}