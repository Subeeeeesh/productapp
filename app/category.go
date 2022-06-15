package app

import (
	"context"

	"encoding/json"

	"fmt"

	"log"

	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct for storing data

type category struct {
	CId string `json:"cid"`

	Cname string `json:"cname"`

	Cdesc string `json:"cdesc"`

	Ccreatedby string `json:"ccreatedby"`

	Cmodifiedby string `json:"cmodifiedby"`

	Cstatus bool `json:"cstatus"`
}

type ResponseError struct {
	ErrorMessage string `json:"errormessage"`

	StatusCode int `json:"statuscode"`

	Status bool `json:"status"`

	CustomMessage string `json:"custommmessage"`
}

type Responses struct {

	//ErrorMessage  string `json:"error message"`

	StatusCode int `json:"statuscode"`

	Status bool `json:"status"`

	CustomMessage string `json:"custommmessage"`

	Response primitive.M
}

type Response struct {

	//ErrorMessage  string `json:"error message"`

	StatusCode int `json:"statuscode"`

	Status bool `json:"status"`

	CustomMessage string `json:"custommmessage"`

	Response []primitive.M
}

var categoryCollection = db().Database("ProductApp").Collection("Category") // get collection "users" from db() which returns *mongo.Client

// Create Profile or Signup

func CreateCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // for adding Content-type

	var cat1 category

	err := json.NewDecoder(r.Body).Decode(&cat1) // storing in cat1  variable of type category

	if err != nil {

		fmt.Print(err)

	}

	insertResult, err1 := categoryCollection.InsertOne(context.TODO(), cat1)

	if err1 != nil {

		log.Fatal(err)

	}

	fmt.Println("Inserted a single document: ", insertResult)

	json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document

}

// Get Profile of a particular User by Name

func GetCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"]

	var body category

	e := json.NewDecoder(r.Body).Decode(&body)

	if e != nil {

		fmt.Print(e)

	}

	var result primitive.M //  an unordered representation of a BSON document which is a Map

	err := categoryCollection.FindOne(context.TODO(), bson.D{{"cid", params}}).Decode(&result)

	if err != nil {

		fmt.Println(err)

	}
	if result == nil {

		msg := ResponseError{

			ErrorMessage: "no data found",

			StatusCode: 200,

			Status: false,

			CustomMessage: "id not exist",
		}

		json.NewEncoder(w).Encode(msg)

	} else {

		msg := Responses{

			StatusCode: 200,

			Status: true,

			CustomMessage: "success",

			Response: result}

		json.NewEncoder(w).Encode(msg)

	}

}

//Update Profile of User

func UpdateCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	validate := validator.New()

	type updateBody struct {
		CId         string `json:"cid" validate:"required,alphanum,min=4,max=10"` //value that has to be matched
		Cname       string `json:"cname" validate:"required,min=3,max=20"`        // value that has to be modified
		Cdesc       string `json:"cdesc" validate:"required,min=5,max=100"`       // value that has to be modified
		Cmodifiedby string `json:"cmodifiedby" validate:"required,min=3,max=20"`  // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
		w.WriteHeader(400)
	}

	errv := validate.Struct(body) // update struct validation

	if errv != nil {
		fmt.Println(errv)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Id not found")
	} else {
		filter := bson.D{{"cid", body.CId}} // converting value to BSON type
		after := options.After              // for returning updated document
		returnOpt := options.FindOneAndUpdateOptions{

			ReturnDocument: &after,
		}
		update := bson.D{{"$set", bson.D{{"cname", body.Cname}, {"cdesc", body.Cdesc}, {"cmodifiedby", body.Cmodifiedby}}}}
		updateResult := categoryCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

		var result primitive.M
		_ = updateResult.Decode(&result)

		json.NewEncoder(w).Encode(result)
		w.WriteHeader(200)
	}

}

// Update Category Status

func UpdateCategoryStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		CId string `json:"cid"` //value that has to be matched

		Cstatus bool `json:"cstatus"` // value that has to be modified

	}

	var body updateBody

	e := json.NewDecoder(r.Body).Decode(&body)

	if e != nil {

		fmt.Print(e)

	}

	filter := bson.D{{"cid", body.CId}} // converting value to BSON type

	after := options.After // for returning updated document

	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}

	update := bson.D{{"$set", bson.D{{"cstatus", body.Cstatus}}}}

	updateResult := categoryCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M

	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)

}

//Delete Profile of User

func DeleteCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)["id"] //get Parameter value as string

	// _id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID

	// if err != nil {

	//  fmt.Printf(err.Error())

	// }

	opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase

	res, err := categoryCollection.DeleteOne(context.TODO(), bson.D{{"cid", params}}, opts)

	if err != nil {

		log.Fatal(err)

	}

	fmt.Printf("deleted %v documents\n", res.DeletedCount)

	json.NewEncoder(w).Encode(res.DeletedCount) // return number of documents deleted

}

func GetAllCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                       //slice for multiple documents
	cur, err := categoryCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {

		fmt.Println(err)
		w.WriteHeader(400)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(500)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	fmt.Println("res", results)
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted

	w.WriteHeader(200)

	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted

	if results == nil {

		msg := ResponseError{

			ErrorMessage: "no data found",

			StatusCode: 200,

			Status: false,

			CustomMessage: "no data",
		}

		json.NewEncoder(w).Encode(msg)

	} else {

		msg := Response{

			StatusCode: 200,

			Status: true,

			CustomMessage: "success",

			Response: results}

		json.NewEncoder(w).Encode(msg)

	}

}
