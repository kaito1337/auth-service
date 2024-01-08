package mongodb

import (
	"auth-backend/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	cfg    *config.MongoDBConnectionConfig
}

func NewMongoDB(cfg *config.MongoDBConnectionConfig) *MongoDB {
	return &MongoDB{
		client: nil,
		cfg:    cfg,
	}
}

func (m *MongoDB) Connect() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?retryWrites=true&w=majority&authSource=admin", m.cfg.Username, m.cfg.Password, m.cfg.Host, m.cfg.Port, m.cfg.Database)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return err
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return nil
}

func (m *MongoDB) Release() {
	if m.client != nil {
		m.client.Disconnect(context.TODO())
	}
}

func (m *MongoDB) GetDB() *mongo.Database {
	return m.client.Database(m.cfg.Database)
}
