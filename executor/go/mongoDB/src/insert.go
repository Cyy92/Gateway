package src

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertContainerInfo(client *mongo.Client, database string, collection string, document map[string]interface{}) (string, error) {
	col := client.Database(database).Collection(collection)
	fmt.Println("Current Collection name: ", col.Name())

	// Insert user's json input
	result, err := col.InsertOne(context.TODO(), document)
	if err != nil {
		log.Fatalln(err)
	}

	// Generate reg num with object id and current timestamp ///////
	timestamp := time.Now()
	year := strconv.Itoa(timestamp.Year())
	month := strconv.Itoa(int(timestamp.Month()))
	date := strconv.Itoa(timestamp.Day())
	prefix := fmt.Sprintf("%s%02s%02s", year, month, date)
	reg_num := prefix + result.InsertedID.(primitive.ObjectID).Hex()
	////////////////////////////////////////////////////////////////

	// Update container preinfo with reg num /////////////////////////////////////////////
	// Find document with id and update document with reg num ////////////////////////////
	var updatedDoc bson.M
	filter := bson.D{{"_id", result.InsertedID}}
	fupopts := options.FindOneAndUpdate().SetUpsert(true)

	// Setting new document with reg num & timestamp, then execute mongodb's update method
	document["reg_num"] = reg_num
	document["timestamp"] = result.InsertedID.(primitive.ObjectID).Timestamp().String()
	for _, value := range document["obj_data"].([]interface{}) {
		value.(map[string]interface{})["path"] = "/minio/" + reg_num + "/" + value.(map[string]interface{})["obj_name"].(string)
	}
	phase := bson.D{{"$set", document}}
	err = col.FindOneAndUpdate(context.TODO(), filter, phase, fupopts).Decode(&updatedDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "DB Error", err
		}
		log.Fatal(err)
	}
	/////////////////////////////////////////////////////////////////////////////////////

	fmt.Printf("updated document %v\n", updatedDoc)

	/////////////////////////////////////////////////////////////////////////////////////
	// Then, insert reg num(meta key) and create empty document for multi meta keys
	/////////////////////////////////////////////////////////////////////////////////////
	col = client.Database(database).Collection("mdManager")
	fmt.Println("Current Collection name: ", col.Name())

	keys := make(map[string]interface{})
	keys["reg_num"] = reg_num
	keys["multi_meta_keys"] = []map[string]interface{}{{}}

	result, err = col.InsertOne(context.TODO(), keys)
	if err != nil {
		log.Fatalln(err)
	}
	/////////////////////////////////////////////////////////////////////////////////////

	fmt.Printf("insert meta data (%v) to mdManager, successfully\n", result.InsertedID)

	return reg_num, nil
}

func InsertWorkingSetInfo(client *mongo.Client, database string, collection string, document map[string]interface{}) error {
	col := client.Database(database).Collection(collection)
	fmt.Println("Current Collection name: ", col.Name())

	// Insert user's json input
	result, err := col.InsertOne(context.TODO(), document)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("insert meta data (%v) to workingsetInfo, successfully\n", result.InsertedID)

	return nil
}
