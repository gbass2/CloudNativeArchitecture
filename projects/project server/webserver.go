package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"math/rand"
	"time"
	"encoding/json"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"project/gameapi"
)

// Mongodb
const (
	mongodbEndpoint = "mongodb://172.17.0.2:27017" // Find this from the Mongo container.
)

type dollars float32

// Database collection entries.
type Inventory struct {
	ID       primitive.ObjectID `bson:"_id"`
	Title    string             `bson:"title"`
	Price    dollars            `bson:"price,truncate"`
	Quantity int					      `bson:"quantity"`
	InStock  bool								`bson:"in_stock"`
	Sku			 int							  `bson:"sku"`
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
	db.col = db.client.Database("gamepub").Collection("inventory")

	mux := http.NewServeMux()
	mux.HandleFunc("/search", db.search)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

// Responds with the price of the requested item
func (db database) search(w http.ResponseWriter, req *http.Request) {
	// Check to see if the request is a GET
	if req.Method == "GET" {
		title := req.URL.Query().Get("title")
		title = strings.ReplaceAll(title, "_", " ") // Replace underscores with spaces.
		title = strings.ToLower(title) // Lowercase the title.
		title = strings.TrimSpace(title) // Remove leading and trailing whitespace.


		// Query the IGDB api.
		var games []map[string]interface{} // Slice of map of the returned data.
		games, err := gameapi.SearchGame(title)

		if err != nil {
			gotJson, _ := json.MarshalIndent(games, "", "    ")
			fmt.Fprintf(w, string(gotJson))
			return
		}

		// Look into local db to get information of the games and if exists then add info to game map.
		for i,game := range games {
				// Send the game info of the requested title if it exists. If not then send an error not found
				filter := bson.M{"title": strings.ToLower(game["name"].(string))}
				var elem Inventory

				if err := db.col.FindOne(db.ctx, filter).Decode(&elem); err != nil {
					games[i]["price"] = nil
					games[i]["quantity"] = nil
					games[i]["in_stock"] = false
					games[i]["sku"] = nil
				} else {
					games[i]["price"] = elem.Price
					games[i]["quantity"] = elem.Quantity
					games[i]["in_stock"] = elem.InStock
					games[i]["sku"] = elem.Sku
				}
		}
		gotJson, _ := json.MarshalIndent(games, "", "    ")
		fmt.Fprintf(w, string(gotJson))

	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a GET then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

// Creates an element in the db map
func (db database) create(w http.ResponseWriter, req *http.Request) {
	// Check to see if the request is a POST and not a GET
	if req.Method == "POST" {
		title := req.URL.Query().Get("title")
		title = strings.ReplaceAll(title, "_", " ") 		// Replace underscores with spaces.
		title = strings.ToLower(title) // Lower case the title.
		title = strings.TrimSpace(title) // Remove leading and trailing whitespace.

		price, _ := strconv.ParseFloat(req.URL.Query().Get("price"), 32)
		quantity, err := strconv.Atoi(req.URL.Query().Get("quantity"))
		inStock := true

		if quantity < 1 {
				inStock = false
		}
		// Check to see if the item is in the database.
		filter := bson.M{"title": title}
		var elem Inventory

		if err := db.col.FindOne(db.ctx, filter).Decode(&elem); err == nil {
			db.update(w,req)
			return
		}

		// Get random number for the sku.
		rand.Seed(time.Now().UnixNano())
		sku := rand.Intn(1000)

		// Adding the item to the database
		_, err = db.col.InsertOne(db.ctx, &Inventory{
			ID:    primitive.NewObjectID(),
			Title: title,
			Price: dollars(price),
			Quantity: quantity,
			InStock: inStock,
			Sku: sku,
		})

		// Checking for error.
		if err != nil {
			fmt.Fprintf(w, "Error: %s\n", err)
			return
		} else {
			fmt.Fprintf(w, "Title added without errors \n")
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
		title := req.URL.Query().Get("title")
		title = strings.ReplaceAll(title, "_", " ") // Replace underscores with spaces.
		title = strings.ToLower(title) // Lowercase the title.
		title = strings.TrimSpace(title) // Remove leading and trailing whitespace.

		price, _ := strconv.ParseFloat(req.URL.Query().Get("price"), 32)
		quantity, _ := strconv.Atoi(req.URL.Query().Get("quantity"))
		inStock := true

		if quantity < 1 {
				inStock = false
		}

		// Check to see if the item  exists.
		filter := bson.M{"title": title}
		var elem Inventory

		if err := db.col.FindOne(db.ctx, filter).Decode(&elem); err != nil {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "(Error) Title does not exists: %q\n", title)
			return
		}

		update := bson.M{"$set": bson.M{"price": dollars(price),"quantity":quantity,"in_stock":inStock}}

		if _, err := db.col.UpdateOne(db.ctx, filter, update); err != nil {
			fmt.Fprintf(w, "(Error) Update request failed: %s\n", err)
			return
		} else {
			fmt.Fprintf(w, "Title updated without errors \n")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest) // If the request is not a POST then respond with a bad request
		fmt.Fprintf(w, "Error: Bad Request\n")
	}
}

// Deletes an item in the db map
func (db database) delete(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		title := req.URL.Query().Get("title")
		title = strings.ReplaceAll(title, "_", " ") // Replace underscores with spaces.
		title = strings.ToLower(title) // Lowercase the title.
		title = strings.TrimSpace(title) // Remove leading and trailing whitespace.

		// If the title exists then delete it
		// Check to see if the title  exists.
		filter := bson.M{"title": title}
		var elem Inventory

		if err := db.col.FindOne(db.ctx, filter).Decode(&elem); err != nil {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "(Error) title does not exists: %q\n", title)
			return
		}

		if _, err := db.col.DeleteOne(db.ctx, filter); err != nil {
			fmt.Fprintf(w, "(Error) delete request failed: %s\n", err)
			return
		} else {
			fmt.Fprintf(w, "title deleted without errors \n")
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
