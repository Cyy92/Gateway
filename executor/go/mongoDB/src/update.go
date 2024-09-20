package src

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateMetaByTrigger(client *mongo.Client, database string, collection string, key string, regNum string, code string) error {
	col := client.Database(database).Collection(collection)
	fmt.Println("Current Collection name: ", col.Name())
	filter := bson.D{{"reg_num", key}}

	// Setting update options
	upopts := options.Update().SetUpsert(true)

	////////////////////////////////////////////////////////////////
	// Firstly, check either multi meta keys field is empty or not
	////////////////////////////////////////////////////////////////
	var result bson.M
	if err := col.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		log.Fatal(err)
	}
	// then, marshaling / unmarshaling result
	marshaled, _ := json.Marshal(result)
	m := make(map[string]interface{})
	json.Unmarshal(marshaled, &m)

	////////////////////////////////////////////////////////////////
	// Secondly, update multi meta keys field
	////////////////////////////////////////////////////////////////

	for _, item := range m["multi_meta_keys"].([]interface{}) {
		// if multi meta keys field is empty,
		if item.(map[string]interface{})["code"] == nil || item.(map[string]interface{})["working_reg_num"] == nil {
			update := make(map[string]interface{})
			update["multi_meta_keys"] = []map[string]interface{}{
				{
					"code":            code,
					"working_reg_num": regNum,
				},
			}
			phase := bson.D{{"$set", update}}
			upresult, err := col.UpdateOne(context.TODO(), filter, phase, upopts)
			if err != nil {
				log.Fatalln(err)
			}

			if upresult.MatchedCount != 0 {
				fmt.Println("matched and replaced an existing document")
				return nil
			}

			if upresult.UpsertedCount != 0 {
				fmt.Printf("inserted a new document with ID %v\n", upresult.UpsertedID)
			}
		} else {
			update := make(map[string]interface{})
			update["multi_meta_keys"] = map[string]interface{}{
				"code":            code,
				"working_reg_num": regNum,
			}
			phase := bson.D{{"$push", update}}
			upresult, err := col.UpdateOne(context.TODO(), filter, phase, upopts)
			if err != nil {
				log.Fatalln(err)
			}

			if upresult.MatchedCount != 0 {
				fmt.Println("matched and replaced an existing document")
				return nil
			}

			if upresult.UpsertedCount != 0 {
				fmt.Printf("inserted a new document with ID %v\n", upresult.UpsertedID)
			}

		}
	}

	return nil
}

func UpdateMetaByWorkingSet(client *mongo.Client, database string, collection string, key string, regNum string, code string) error {
	col := client.Database(database).Collection(collection)
	fmt.Println("Current Collection name: ", col.Name())
	filter := bson.D{{"reg_num", key}}

	// Setting update options
	fupopts := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.code": code},
		},
	})

	////////////////////////////////////////////////////////////////
	// update working set reg num
	////////////////////////////////////////////////////////////////
	var updatedDoc bson.M
	up := bson.M{"$set": bson.M{"multi_meta_keys.$[elem].working_reg_num": regNum}}
	err := col.FindOneAndUpdate(context.TODO(), filter, up, fupopts).Decode(&updatedDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Fatal(err)
		}
		log.Fatal(err)
	}
	fmt.Println("update doc %v", updatedDoc)

	return nil
}

func UpdateWorkingSetInfo(client *mongo.Client, database string, collection string, key string, regNum string, update map[string]interface{}) error {
	col := client.Database(database).Collection(collection)
	fmt.Println("Current Collection name: ", col.Name())
	filter := bson.D{{"reg_num", key}}

	// Setting update options
	upopts := options.Update().SetUpsert(true)

	////////////////////////////////////////////////////////////////
	// Firstly, get container num
	////////////////////////////////////////////////////////////////
	var result bson.M
	if err := col.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		log.Fatal(err)
	}
	// then, marshaling / unmarshaling result
	marshaled, _ := json.Marshal(result)
	m := make(map[string]interface{})
	json.Unmarshal(marshaled, &m)

	////////////////////////////////////////////////////////////////
	// Secondly, setting new document
	// then, update database
	////////////////////////////////////////////////////////////////
	update["working_reg_num"] = regNum
	update["container_num"] = m["container_num"].(string)
	for _, value := range update["obj_data"].([]interface{}) {
		value.(map[string]interface{})["path"] = "/minio/" + key + "/" + regNum + "/" + value.(map[string]interface{})["obj_name"].(string)
	}
	unixTime := time.Now()
	update["timestamp"] = time.Unix(unixTime.Unix(), 0).UTC().String()
	phase := bson.D{{"$set", update}}
	upresult, err := col.UpdateOne(context.TODO(), filter, phase, upopts)
	if err != nil {
		log.Fatalln(err)
	}

	if upresult.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return nil
	}

	if upresult.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", upresult.UpsertedID)
	}

	return nil
}
