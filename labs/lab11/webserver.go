package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongodb
const (
	mongodbEndpoint = "mongodb://192.168.64.2:31846" // Find this from the Mongo container.
)

type dollars float32

// Database collection entries.
type Inventory struct {
	ID    primitive.ObjectID `bson:"_id"`
	Item  string             `bson:"item"`
	Price dollars            `bson:"price,truncate"`
}

// Holds the connection and collection to the database.
type database struct {
	client *mongo.Client
	col    *mongo.Collection
	ctx    context.Context
}

func main() {
	// Creating a database object.
	var db database
	var err error

	// create a mongo client
	db.client, err = mongo.NewClient(
		options.Client().ApplyURI(mongodbEndpoint),
	)
	checkError(err)

	// Connect to mongo
	db.ctx = context.Background()
	err = db.client.Connect(db.ctx)

	// Disconnect
	defer db.client.Disconnect(db.ctx)

	// select collection from database
	db.col = db.client.Database("store").Collection("inventory")

	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

// Responds with the items in the map and their prices
func (db database) list(w http.ResponseWriter, req *http.Request) {
	// Check to see if the request is a GET
	if req.Method == "GET" {
		// send a list of the items in the database
		// Retrieving the list from the database

		cur, err := db.col.Find(context.TODO(), bson.D{})

		// Checking for errors.
		if err != nil {
			fmt.Fprintf(w, "Error: %s\n", err)
			return
		}

		for cur.Next(context.TODO()) {
			//Create a value into which the single document can be decoded
			var elem Inventory
			err := cur.Decode(&elem)

			if err != nil {
				fmt.Fprintf(w, "Error: %s\n", err)
				return
			}

			// Sending back the items.
			fmt.Fprintf(w, "%s: %.2f\n", elem.Item, elem.Price)

		}

		//Close the cursor once finished
		cur.Close(context.TODO())
	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a GET then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

// Responds with the price of the requested item
func (db database) price(w http.ResponseWriter, req *http.Request) {
	// Check to see if the request is a GET
	if req.Method == "GET" {
		item := req.URL.Query().Get("item")

		// Send the price of the requested item if it exists. If not then send an error not found
		filter := bson.M{"item": item}
		var elem Inventory

		if err := db.col.FindOne(db.ctx, filter).Decode(&elem); err != nil {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "(Error) No such item: %q\n", item)
			return
		}

		fmt.Fprintf(w, "%s: %.2f\n", elem.Item, elem.Price)

	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a GET then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

// Creates an element in the db map
func (db database) create(w http.ResponseWriter, req *http.Request) {
	// Check to see if the request is a POST and not a GET
	if req.Method == "POST" {
		item := req.URL.Query().Get("item")
		price, _ := strconv.ParseFloat(req.URL.Query().Get("price"), 32)

		// Check to see if the item is in the database.
		filter := bson.M{"item": item}
		var elem Inventory

		if err := db.col.FindOne(db.ctx, filter).Decode(&elem); err == nil {
			fmt.Fprintf(w, "(Error) Item already exists: %q\n", item)
			return
		}

		// Adding the item to the database
		_, err := db.col.InsertOne(db.ctx, &Inventory{
			ID:    primitive.NewObjectID(),
			Item:  item,
			Price: dollars(price),
		})

		// Checking for error.
		if err != nil {
			fmt.Fprintf(w, "Error: %s\n", err)
			return
		} else {
			fmt.Fprintf(w, "Item added without errors \n")
		}

	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a POST then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

// Updates an item in the db map
func (db database) update(w http.ResponseWriter, req *http.Request) {
	// Check to see if the request is a POST and not a GET
	if req.Method == "POST" {
		item := req.URL.Query().Get("item")
		price, _ := strconv.ParseFloat(req.URL.Query().Get("price"), 32)

		// Check to see if the item  exists.
		filter := bson.M{"item": item}
		var elem Inventory

		if err := db.col.FindOne(db.ctx, filter).Decode(&elem); err != nil {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "(Error) Item does not exists: %q\n", item)
			return
		}

		update := bson.M{"$set": bson.M{"price": dollars(price)}}

		if _, err := db.col.UpdateOne(db.ctx, filter, update); err != nil {
			fmt.Fprintf(w, "(Error) Update request failed: %s\n", err)
			return
		} else {
			fmt.Fprintf(w, "Item updated without errors \n")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a POST then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

// Deletes an item in the db map
func (db database) delete(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		item := req.URL.Query().Get("item")

		// If the item exists then delete it
		// Check to see if the item  exists.
		filter := bson.M{"item": item}
		var elem Inventory

		if err := db.col.FindOne(db.ctx, filter).Decode(&elem); err != nil {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "(Error) Item does not exists: %q\n", item)
			return
		}

		if _, err := db.col.DeleteOne(db.ctx, filter); err != nil {
			fmt.Fprintf(w, "(Error) delete request failed: %s\n", err)
			return
		} else {
			fmt.Fprintf(w, "Item deleted without errors \n")
		}

	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a POST then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
