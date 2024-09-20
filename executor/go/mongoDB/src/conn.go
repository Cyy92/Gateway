package src

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConn(uri string, user string, password string, database string) *mongo.Client {
	credential := options.Credential{
		Username:   user,
		Password:   password,
		AuthSource: database,
	}

	clientOptions := options.Client().ApplyURI(uri).SetAuth(credential)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("MongoDB Connection Made")

	return client
}
