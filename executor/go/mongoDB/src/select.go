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

func GetWorkingsetInfo(client *mongo.Client, database string, collection string) string {
	// Get collection's all documents
	// then, marshaling result
	////////////////////////////////////////////////////////////////////////////
	col := client.Database(database).Collection(collection)

	opts := options.Find().SetSort(bson.D{{"code", 1}})
	cursor, err := col.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	marshaled, _ := json.Marshal(results)

	mapped := make([]map[string]interface{}, 0)
	json.Unmarshal(marshaled, &mapped)

	////////////////////////////////////////////////////////////////////////////
	// Get workingset info
	////////////////////////////////////////////////////////////////////////////
	wsValue := make([]map[string]interface{}, len(mapped))
	meta := make(map[string]interface{})

	for i, val := range mapped {
		custommap := make(map[string]interface{})
		custommap["code"] = val["code"].(string)
		custommap["name"] = val["name"].(string)
		custommap["trigger_name"] = val["trigger_name"].(string)
		custommap["db_table"] = val["db_table"].(string)
		wsValue[i] = custommap
	}

	meta["workingset_info"] = wsValue
	out, _ := json.Marshal(meta)

	return string(out)
}

func GetRegNum(client *mongo.Client, database string, collection string, code string) string {
	// Get collection's document with code
	// then, marshaling result
	////////////////////////////////////////////////////////////////////////////
	col := client.Database(database).Collection(collection)
	filter := bson.D{{"related_works.code", code}}

	opts := options.Find().SetSort(bson.D{{"reg_num", 1}})
	cursor, err := col.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	marshaled, _ := json.Marshal(results)

	mapped := make([]map[string]interface{}, 0)
	json.Unmarshal(marshaled, &mapped)

	////////////////////////////////////////////////////////////////////////////
	// Get reg num matching the code
	// and return string result
	////////////////////////////////////////////////////////////////////////////
	regValue := make([]map[string]interface{}, len(mapped))
	meta := make(map[string]interface{})

	for i, val := range mapped {
		custommap := make(map[string]interface{})
		custommap["reg_num"] = val["reg_num"].(string)
		custommap["value"] = "---"
		regValue[i] = custommap
	}

	meta["reg_value"] = regValue
	out, _ := json.Marshal(meta)

	return string(out)
}

func GetTriggerName(client *mongo.Client, database string, collection string, code string) string {
	// Get collection's document with code
	// then, marshaling result
	////////////////////////////////////////////////////////////////////////////
	col := client.Database(database).Collection(collection)
	filter := bson.D{{"code", code}}

	var result bson.M
	if err := col.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		log.Fatal(err)
	}

	marshaled, _ := json.Marshal(result)

	mapped := make(map[string]interface{})
	json.Unmarshal(marshaled, &mapped)

	////////////////////////////////////////////////////////////////////////////
	// Return triiger name
	////////////////////////////////////////////////////////////////////////////
	return mapped["trigger_name"].(string)
}

func GetFilePath(client *mongo.Client, database string, collection string, regNum string, objName string) string {
	// Get collection's document with code
	// then, marshaling result
	////////////////////////////////////////////////////////////////////////////
	col := client.Database(database).Collection(collection)

	// Set aggregate query & option
	unwind := bson.D{{"$unwind", "$obj_data"}}
	match := bson.D{{"$match", bson.D{{"reg_num", regNum}, {"obj_data.obj_name", objName}}}}
	opts := options.Aggregate().SetMaxTime(2 * time.Second)

	// then, aggregation
	cursor, ag_err := col.Aggregate(
		context.TODO(),
		mongo.Pipeline{unwind, match},
		opts)
	if ag_err != nil {
		log.Fatal(ag_err)
	}

	// get documen from db
	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	marshaled, _ := json.Marshal(results)

	mapped := make([]map[string]interface{}, 0)
	json.Unmarshal(marshaled, &mapped)

	////////////////////////////////////////////////////////////////////////////
	// Return file path
	////////////////////////////////////////////////////////////////////////////
	return mapped[0]["obj_data"].(map[string]interface{})["path"].(string)
}

func SelectMetaAll(client *mongo.Client, database string, collection string, key string) string {
	col := client.Database(database).Collection(collection)
	fmt.Println("Current Collection name: ", col.Name())
	filter := bson.D{{"reg_num", key}}

	////////////////////////////////////////////////////////////////
	// Firstly, get multi meta keys by reg num
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
	// Secondly, get multi meta data
	////////////////////////////////////////////////////////////////
	multi_meta := make(map[string]interface{})
	con_meta := make(map[string]interface{})
	data := make([]map[string]interface{}, len(m["multi_meta_keys"].([]interface{})))

	con_col := client.Database(database).Collection("containerPreInfo")
	con_filter := bson.D{{"reg_num", key}}

	var con_result bson.M
	if err := con_col.FindOne(context.TODO(), con_filter).Decode(&con_result); err != nil {
		log.Fatal(err)
	}

	con_marshaled, _ := json.Marshal(con_result)
	conmap := make(map[string]interface{})
	json.Unmarshal(con_marshaled, &conmap)

	// Set new map for user return
	con_meta["code"] = "00000"
	con_meta["reg_num"] = conmap["reg_num"]
	con_meta["obj_data"] = conmap["obj_data"]
	con_meta["meta_data"] = conmap["meta_data"]

	data = append(data, con_meta)

	for i, item := range m["multi_meta_keys"].([]interface{}) {
		if item.(map[string]interface{})["code"].(string) == "DA001" {
			col = client.Database(database).Collection("dangerInfo")
			filter = bson.D{{"working_reg_num", item.(map[string]interface{})["working_reg_num"].(string)}}
			//var danger_result bson.M
			if err := col.FindOne(context.TODO(), filter).Decode(result); err != nil {
				log.Fatal(err)
			}
			marshaled, _ = json.Marshal(result)
			mp := make(map[string]interface{})
			json.Unmarshal(marshaled, &mp)

			// Set new map for user return
			ws_meta := make(map[string]interface{})
			ws_meta["code"] = item.(map[string]interface{})["code"]
			ws_meta["reg_num"] = mp["working_reg_num"]
			ws_meta["obj_data"] = mp["obj_data"]
			ws_meta["meta_data"] = mp["meta_data"]

			data[i] = ws_meta
		} else if item.(map[string]interface{})["code"].(string) == "EE332" {
			col = client.Database(database).Collection("drugInfo")
			filter = bson.D{{"working_reg_num", item.(map[string]interface{})["working_reg_num"].(string)}}
			//var danger_result bson.M
			if err := col.FindOne(context.TODO(), filter).Decode(result); err != nil {
				log.Fatal(err)
			}
			marshaled, _ = json.Marshal(result)
			mp := make(map[string]interface{})
			json.Unmarshal(marshaled, &mp)

			// Set new map for user return
			ws_meta := make(map[string]interface{})
			ws_meta["code"] = item.(map[string]interface{})["code"]
			ws_meta["reg_num"] = mp["working_reg_num"]
			ws_meta["obj_data"] = mp["obj_data"]
			ws_meta["meta_data"] = mp["meta_data"]

			data[i] = ws_meta
		}
	}

	// Merge container pre info & working set info
	//multi_meta["multi_meta_data"] = []map[string]interface{}{con_meta, danger_meta}
	multi_meta["multi_meta_data"] = data

	out, _ := json.Marshal(multi_meta)

	return string(out)
}
