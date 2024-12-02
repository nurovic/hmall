package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
