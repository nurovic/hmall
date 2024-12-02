package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string  `json:"user_id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}
