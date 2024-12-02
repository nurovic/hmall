package store

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	clientOnce sync.Once
)

// MongoDB bağlantısını başlatır ve döner
func GetMongoClient() *mongo.Client {
	clientOnce.Do(func() {
		var err error

		// .env dosyasını yükle
		err = godotenv.Load()
		if err != nil {
			log.Fatalf(".env dosyası yüklenemedi: %v", err)
		}

		// MONGO_DB çevresel değişkenini oku
		mongoURI := os.Getenv("MONGO_DB")
		if mongoURI == "" {
			log.Fatalf("MONGO_DB çevresel değişkeni tanımlanmamış.")
		}

		// MongoDB bağlantı ayarlarını yap
		clientOptions := options.Client().ApplyURI(mongoURI)
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatalf("MongoDB bağlantısı başlatılamadı: %v", err)
		}

		// Bağlantıyı doğrulamak için ping atıyoruz
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatalf("MongoDB bağlantısı doğrulanamadı: %v", err)
		}

		log.Println("MongoDB bağlantısı başarılı!")
	})

	return client
}
